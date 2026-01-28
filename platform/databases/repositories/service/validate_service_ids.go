package service

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

func (r *repository) ValidateServiceIDs(ctx context.Context, serviceIDs []string) error {
	if len(serviceIDs) == 0 {
		return nil
	}

	log.Info(logger.LogBranchServicesRepoValidateIDs, "count", len(serviceIDs))

	placeholders := ""
	args := make([]interface{}, len(serviceIDs))
	for i, id := range serviceIDs {
		if i > 0 {
			placeholders += ","
		}
		placeholders += "?"
		args[i] = id
	}

	query := "SELECT COUNT(*) FROM services WHERE id IN (" + placeholders + ")"
	var count int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		log.Error(logger.LogBranchServicesRepoValidateErr, "error", err)
		return err
	}

	if count != len(serviceIDs) {
		log.Warn(logger.LogBranchServicesRepoValidateMiss, "expected", len(serviceIDs), "found", count)
		return domain.ErrServiceNotFound
	}

	return nil
}
