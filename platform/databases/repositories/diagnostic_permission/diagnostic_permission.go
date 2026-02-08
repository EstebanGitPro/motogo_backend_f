package diagnostic_permission

import (
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
)

// DiagnosticPermission represents the database model for permisos_moto_sede table
type DiagnosticPermission struct {
	ID           string `db:"id"`
	MotorcycleID string `db:"motorcycle_id"`
	BranchID     string `db:"branch_id"`
	Active       bool   `db:"active"`
}

// ToDomain converts the database model to domain entity
func (p *DiagnosticPermission) ToDomain() domain.DiagnosticPermission {
	return domain.DiagnosticPermission{
		ID:           p.ID,
		MotorcycleID: p.MotorcycleID,
		BranchID:     p.BranchID,
		Active:       p.Active,
	}
}

// FromDomain converts a domain entity to database model
func FromDomain(permission *domain.DiagnosticPermission) *DiagnosticPermission {
	return &DiagnosticPermission{
		ID:           permission.ID,
		MotorcycleID: permission.MotorcycleID,
		BranchID:     permission.BranchID,
		Active:       permission.Active,
	}
}
