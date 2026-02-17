package interactor

import (
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// deferRollback is a shared helper that handles transaction rollback on error.
// Usage: defer deferRollback(tx, &err, log)
func deferRollback(tx output.Tx, err *error, log logger.Logger) {
	if *err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			log.Error(logger.LogTxRollbackError,
				"rollback_error", rbErr,
				"original_error", *err)
		} else {
			log.Warn(logger.LogTxRollbackOK)
		}
	}
}
