package completed_service

import (
	"database/sql"
	"time"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
)

// CompletedService represents the database model for completed_services table
type CompletedService struct {
	ID                  string          `db:"id"`
	BranchID            string          `db:"branch_id"`
	BranchName          sql.NullString  `db:"branch_name"`
	MotorcycleID        string          `db:"motorcycle_id"`
	DiagnosticID        sql.NullString  `db:"diagnostic_id"`
	RequestDate         time.Time       `db:"request_date"`
	CompletionDate      sql.NullTime    `db:"completion_date"`
	Status              string          `db:"status"`
	QuotedPrice         sql.NullFloat64 `db:"quoted_price"`
	FinalPrice          sql.NullFloat64 `db:"final_price"`
	RepresentativeNotes sql.NullString  `db:"representative_notes"`
	DeletedAt           sql.NullTime    `db:"deleted_at"`
	CreatedAt           time.Time       `db:"created_at"`
	UpdatedAt           time.Time       `db:"updated_at"`
}

// CompletedServiceItem represents the database model for completed_service_items table (pivot)
type CompletedServiceItem struct {
	ID                 string         `db:"id"`
	CompletedServiceID string         `db:"completed_service_id"`
	ServiceID          string         `db:"service_id"`
	ServiceName        sql.NullString `db:"service_name"`
	Rating             sql.NullInt32  `db:"rating"`
	Comment            sql.NullString `db:"comment"`
	RatedAt            sql.NullTime   `db:"rated_at"`
	IsOffensiveComment bool           `db:"is_offensive_comment"`
}

// ServiceStatusHistory represents the database model for service_status_transitions table
type ServiceStatusHistory struct {
	ID                 string         `db:"id"`
	CompletedServiceID string         `db:"completed_service_id"`
	PreviousStatus     sql.NullString `db:"previous_status"`
	NewStatus          string         `db:"new_status"`
	CreatedBy          string         `db:"created_by"`
	CreatedAt          time.Time      `db:"created_at"`
}

// ToDomain converts the database CompletedService model to domain entity
func (cs *CompletedService) ToDomain() domain.CompletedService {
	result := domain.CompletedService{
		ID:           cs.ID,
		BranchID:     cs.BranchID,
		MotorcycleID: cs.MotorcycleID,
		RequestDate:  cs.RequestDate,
		Status:       domain.ServiceStatus(cs.Status),
		CreatedAt:    cs.CreatedAt,
		UpdatedAt:    cs.UpdatedAt,
	}

	if cs.DiagnosticID.Valid {
		result.DiagnosticID = &cs.DiagnosticID.String
	}
	if cs.CompletionDate.Valid {
		result.CompletionDate = &cs.CompletionDate.Time
	}
	if cs.QuotedPrice.Valid {
		result.QuotedPrice = &cs.QuotedPrice.Float64
	}
	if cs.FinalPrice.Valid {
		result.FinalPrice = &cs.FinalPrice.Float64
	}
	if cs.RepresentativeNotes.Valid {
		result.RepresentativeNotes = &cs.RepresentativeNotes.String
	}
	if cs.BranchName.Valid {
		result.BranchName = &cs.BranchName.String
	}
	if cs.DeletedAt.Valid {
		result.DeletedAt = &cs.DeletedAt.Time
	}

	return result
}

// FromDomain converts a domain CompletedService entity to database model
func FromDomain(cs *domain.CompletedService) *CompletedService {
	dbCS := &CompletedService{
		ID:           cs.ID,
		BranchID:     cs.BranchID,
		MotorcycleID: cs.MotorcycleID,
		RequestDate:  cs.RequestDate,
		Status:       string(cs.Status),
	}

	if cs.DiagnosticID != nil {
		dbCS.DiagnosticID = sql.NullString{String: *cs.DiagnosticID, Valid: true}
	}
	if cs.CompletionDate != nil {
		dbCS.CompletionDate = sql.NullTime{Time: *cs.CompletionDate, Valid: true}
	}
	if cs.QuotedPrice != nil {
		dbCS.QuotedPrice = sql.NullFloat64{Float64: *cs.QuotedPrice, Valid: true}
	}
	if cs.FinalPrice != nil {
		dbCS.FinalPrice = sql.NullFloat64{Float64: *cs.FinalPrice, Valid: true}
	}
	if cs.RepresentativeNotes != nil {
		dbCS.RepresentativeNotes = sql.NullString{String: *cs.RepresentativeNotes, Valid: true}
	}

	return dbCS
}

// ItemFromDomain converts a domain CompletedServiceItem to database model
func ItemFromDomain(item *domain.CompletedServiceItem) *CompletedServiceItem {
	return &CompletedServiceItem{
		ID:                 item.ID,
		CompletedServiceID: item.CompletedServiceID,
		ServiceID:          item.ServiceID,
	}
}

// ItemToDomain converts a database CompletedServiceItem to domain entity
func (i *CompletedServiceItem) ItemToDomain() domain.CompletedServiceItem {
	result := domain.CompletedServiceItem{
		ID:                 i.ID,
		CompletedServiceID: i.CompletedServiceID,
		ServiceID:          i.ServiceID,
		IsOffensiveComment: i.IsOffensiveComment,
	}
	if i.ServiceName.Valid {
		result.ServiceName = &i.ServiceName.String
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

// HistoryFromDomain converts a domain ServiceStatusHistory to database model
func HistoryFromDomain(h *domain.ServiceStatusHistory) *ServiceStatusHistory {
	dbH := &ServiceStatusHistory{
		ID:                 h.ID,
		CompletedServiceID: h.CompletedServiceID,
		NewStatus:          string(h.NewStatus),
		CreatedBy:          h.CreatedBy,
	}
	if h.PreviousStatus != nil {
		dbH.PreviousStatus = sql.NullString{String: string(*h.PreviousStatus), Valid: true}
	}
	return dbH
}

// HistoryToDomain converts a database ServiceStatusHistory model to domain entity
func (h *ServiceStatusHistory) HistoryToDomain() domain.ServiceStatusHistory {
	result := domain.ServiceStatusHistory{
		ID:                 h.ID,
		CompletedServiceID: h.CompletedServiceID,
		NewStatus:          domain.ServiceStatus(h.NewStatus),
		CreatedBy:          h.CreatedBy,
		CreatedAt:          h.CreatedAt,
	}
	if h.PreviousStatus.Valid {
		ps := domain.ServiceStatus(h.PreviousStatus.String)
		result.PreviousStatus = &ps
	}
	return result
}
