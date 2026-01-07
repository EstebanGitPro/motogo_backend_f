package brand

import (
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
)

// Brand represents the database entity for motorcycle brands
type Brand struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

// ToDomain converts Brand to domain.Brand
func (b *Brand) ToDomain() domain.Brand {
	return domain.Brand{
		ID:   b.ID,
		Name: b.Name,
	}
}

// FromDomain converts domain.Brand to Brand entity
func FromDomain(domainBrand domain.Brand) Brand {
	return Brand{
		ID:   domainBrand.ID,
		Name: domainBrand.Name,
	}
}
