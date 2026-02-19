package rating

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strings"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

const (
	queryRateServiceItem = `
		UPDATE completed_service_items
		SET rating = ?, comment = ?, is_offensive_comment = ?, rated_at = NOW()
		WHERE id = ?
	`

	queryGetItemByID = `
		SELECT id, completed_service_id, service_id,
			rating, comment, rated_at, is_offensive_comment
		FROM completed_service_items
		WHERE id = ?
	`

	queryGetReviewsByServiceID = `
		SELECT
			csi.rating,
			csi.comment,
			csi.rated_at,
			p.first_name,
			p.last_name,
			CONCAT(b.name, ' ', r.model) AS motorcycle_model,
			s.name AS service_name
		FROM completed_service_items csi
		INNER JOIN completed_services cs ON cs.id = csi.completed_service_id
		INNER JOIN motorcycles m ON m.id = cs.motorcycle_id
		INNER JOIN persons p ON p.id = m.owner_id
		INNER JOIN services s ON s.id = csi.service_id
		LEFT JOIN motorcycle_references r ON r.id = m.reference_id
		LEFT JOIN brands b ON r.brand_id = b.id
		WHERE csi.service_id = ?
		  AND csi.rating IS NOT NULL
		  AND csi.is_offensive_comment = 0
		  AND cs.deleted_at IS NULL
		ORDER BY csi.rated_at DESC
	`
)

var log logger.Logger = logger.NewSlogLogger()

// RatingItem represents the database model for a completed service item (rating-specific)
type RatingItem struct {
	ID                 string         `db:"id"`
	CompletedServiceID string         `db:"completed_service_id"`
	ServiceID          string         `db:"service_id"`
	Rating             sql.NullInt32  `db:"rating"`
	Comment            sql.NullString `db:"comment"`
	RatedAt            sql.NullTime   `db:"rated_at"`
	IsOffensiveComment bool           `db:"is_offensive_comment"`
}

// ReviewRow is the scan target for the reviews query
type ReviewRow struct {
	Rating          int
	Comment         sql.NullString
	RatedAt         sql.NullTime
	FirstName       string
	LastName        string
	MotorcycleModel sql.NullString
	ServiceName     string
}

// ToDomain converts the database RatingItem to domain entity
func (i *RatingItem) ToDomain() domain.CompletedServiceItem {
	result := domain.CompletedServiceItem{
		ID:                 i.ID,
		CompletedServiceID: i.CompletedServiceID,
		ServiceID:          i.ServiceID,
		IsOffensiveComment: i.IsOffensiveComment,
	}
	if i.Rating.Valid {
		rating := int(i.Rating.Int32)
		result.Rating = &rating
	}
	if i.Comment.Valid {
		result.Comment = &i.Comment.String
	}
	if i.RatedAt.Valid {
		result.RatedAt = &i.RatedAt.Time
	}
	return result
}

type repository struct {
	db                        *sql.DB
	stmtRateServiceItem       *sql.Stmt
	stmtGetItemByID           *sql.Stmt
	stmtGetReviewsByServiceID *sql.Stmt
}

// NewRepository creates a new rating repository with prepared statements
func NewRepository(db *sql.DB) (output.RatingRepository, error) {
	if db == nil {
		return nil, sql.ErrConnDone
	}

	stmtRateServiceItem, err := db.Prepare(queryRateServiceItem)
	if err != nil {
		log.Error(logger.LogCSRepoPrepareRateItem, err)
		return nil, fmt.Errorf("error preparing stmtRateServiceItem: %w", err)
	}

	stmtGetItemByID, err := db.Prepare(queryGetItemByID)
	if err != nil {
		log.Error(logger.LogCSRepoPrepareGetItemByID, err)
		return nil, fmt.Errorf("error preparing stmtGetItemByID: %w", err)
	}

	stmtGetReviewsByServiceID, err := db.Prepare(queryGetReviewsByServiceID)
	if err != nil {
		log.Error(logger.LogCSRepoPrepareGetReviews, err)
		return nil, fmt.Errorf("error preparing stmtGetReviewsByServiceID: %w", err)
	}

	return &repository{
		db:                        db,
		stmtRateServiceItem:       stmtRateServiceItem,
		stmtGetItemByID:           stmtGetItemByID,
		stmtGetReviewsByServiceID: stmtGetReviewsByServiceID,
	}, nil
}

// BeginTx starts a new database transaction
func (r *repository) BeginTx(ctx context.Context) (output.Tx, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return common.NewSQLTx(tx), nil
}

// RateServiceItem updates the rating fields for a completed service item (RELEASE_14 / HU48)
func (r *repository) RateServiceItem(ctx context.Context, tx output.Tx, itemID string, rating int, comment *string, isOffensive bool) error {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	_, err := sqlTx.ExecContext(ctx, queryRateServiceItem, rating, comment, isOffensive, itemID)
	if err != nil {
		log.Error(logger.LogCSRepoRateItemErr, err)
		return err
	}

	return nil
}

// GetItemByID retrieves a single completed service item by its ID (RELEASE_14 / HU48)
func (r *repository) GetItemByID(ctx context.Context, itemID string) (*domain.CompletedServiceItem, error) {
	var item RatingItem

	err := r.stmtGetItemByID.QueryRowContext(ctx, itemID).Scan(
		&item.ID, &item.CompletedServiceID, &item.ServiceID,
		&item.Rating, &item.Comment, &item.RatedAt, &item.IsOffensiveComment,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrServiceItemNotFound
		}
		log.Error(logger.LogCSRepoGetItemByIDErr, err)
		return nil, err
	}

	result := item.ToDomain()
	return &result, nil
}

// formatReviewerName returns "FirstName L." for privacy
func formatReviewerName(firstName, lastName string) string {
	firstName = strings.TrimSpace(firstName)
	lastName = strings.TrimSpace(lastName)
	if lastName == "" {
		return firstName
	}
	return fmt.Sprintf("%s %s.", firstName, strings.ToUpper(lastName[:1]))
}

// GetReviewsByServiceID retrieves all reviews for a service type, aggregating average and breakdown
func (r *repository) GetReviewsByServiceID(ctx context.Context, serviceID string) (*domain.ServiceReviewSummary, error) {
	rows, err := r.stmtGetReviewsByServiceID.QueryContext(ctx, serviceID)
	if err != nil {
		log.Error(logger.LogCSRepoGetReviewsErr, err)
		return nil, err
	}
	defer rows.Close()

	var reviews []domain.ServiceReview
	var serviceName string
	var totalRating int
	breakdown := map[int]int{1: 0, 2: 0, 3: 0, 4: 0, 5: 0}

	for rows.Next() {
		var row ReviewRow
		if scanErr := rows.Scan(
			&row.Rating, &row.Comment, &row.RatedAt,
			&row.FirstName, &row.LastName,
			&row.MotorcycleModel, &row.ServiceName,
		); scanErr != nil {
			log.Error(logger.LogCSRepoScanReviewErr, scanErr)
			return nil, scanErr
		}

		serviceName = row.ServiceName
		totalRating += row.Rating
		breakdown[row.Rating]++

		review := domain.ServiceReview{
			ReviewerName: formatReviewerName(row.FirstName, row.LastName),
			Rating:       row.Rating,
		}
		if row.Comment.Valid && row.Comment.String != "" {
			review.Comment = &row.Comment.String
		}
		if row.RatedAt.Valid {
			review.RatedAt = row.RatedAt.Time
		}
		if row.MotorcycleModel.Valid && row.MotorcycleModel.String != "" {
			review.MotorcycleModel = &row.MotorcycleModel.String
		}

		reviews = append(reviews, review)
	}

	if err = rows.Err(); err != nil {
		log.Error(logger.LogCSRepoGetReviewsErr, err)
		return nil, err
	}

	totalReviews := len(reviews)
	var averageRating float64
	if totalReviews > 0 {
		averageRating = math.Round(float64(totalRating)/float64(totalReviews)*10) / 10
	}

	return &domain.ServiceReviewSummary{
		ServiceID:     serviceID,
		ServiceName:   serviceName,
		AverageRating: averageRating,
		TotalReviews:  totalReviews,
		Breakdown:     breakdown,
		Reviews:       reviews,
	}, nil
}
