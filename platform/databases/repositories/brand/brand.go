package brand

import (
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
)

type Brand struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}
func (b *Brand) ToDomain() domain.Brand {
	return domain.Brand{
		ID:   b.ID,
		Name: b.Name,
	}
}

func FromDomain(domainBrand domain.Brand) Brand {
	return Brand{
		ID:   domainBrand.ID,
		Name: domainBrand.Name,
	}
}
