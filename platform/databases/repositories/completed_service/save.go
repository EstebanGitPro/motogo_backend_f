package completed_service

import (
	"context"
	"fmt"
	"strings"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// Save inserts a new completed service record
func (r *repository) Save(ctx context.Context, tx output.Tx, service *domain.CompletedService) error {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	dbCS := FromDomain(service)
	_, err := sqlTx.ExecContext(ctx, queryInsert,
		dbCS.ID,
		dbCS.BranchID,
		dbCS.MotorcycleID,
		dbCS.DiagnosticID,
		dbCS.RequestDate,
		dbCS.Status,
		dbCS.QuotedPrice,
		dbCS.FinalPrice,
		dbCS.RepresentativeNotes,
	)
	if err != nil {
		log.Error(logger.LogCSRepoSaveError, err)
		return fmt.Errorf("error saving completed service: %w", err)
	}

	return nil
}

// SaveItems inserts completed service items (pivot records)
func (r *repository) SaveItems(ctx context.Context, tx output.Tx, items []domain.CompletedServiceItem) error {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	for _, item := range items {
		dbItem := ItemFromDomain(&item)
		_, err := sqlTx.ExecContext(ctx, queryInsertItem,
			dbItem.ID,
			dbItem.CompletedServiceID,
			dbItem.ServiceID,
		)
		if err != nil {
			log.Error(logger.LogCSRepoSaveItemErr, err)
			return fmt.Errorf("error saving completed service item: %w", err)
		}
	}

	return nil
}

// SaveStatusHistory inserts a status history entry
func (r *repository) SaveStatusHistory(ctx context.Context, tx output.Tx, history *domain.ServiceStatusHistory) error {
	sqlTx, ok := tx.(*common.SQLTx)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	dbH := HistoryFromDomain(history)
	_, err := sqlTx.ExecContext(ctx, queryInsertStatusHistory,
		dbH.ID,
		dbH.CompletedServiceID,
		dbH.PreviousStatus,
		dbH.NewStatus,
		dbH.CreatedBy,
	)
	if err != nil {
		log.Error(logger.LogCSRepoSaveHistoryErr, err)
		return fmt.Errorf("error saving status history: %w", err)
	}

	return nil
}

// ValidateBranchServices validates that all service IDs belong to the branch's active services
func (r *repository) ValidateBranchServices(ctx context.Context, branchID string, serviceIDs []string) error {
	// Build IN clause dynamically
	placeholders := make([]string, len(serviceIDs))
	args := make([]interface{}, 0, len(serviceIDs)+1)
	args = append(args, branchID)

	for i, id := range serviceIDs {
		placeholders[i] = "?"
		args = append(args, id)
	}

	// Safe: fmt.Sprintf only injects hardcoded "?" placeholders; all values are passed as parameterized arguments via args
	query := fmt.Sprintf(queryValidateBranchServices, strings.Join(placeholders, ", "))

	var count int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		log.Error(logger.LogCSRepoValidateSvcErr, err)
		return fmt.Errorf("error validating branch services: %w", err)
	}

	if count != len(serviceIDs) {
		return fmt.Errorf("some services are not associated with this branch or are inactive")
	}

	return nil
}
