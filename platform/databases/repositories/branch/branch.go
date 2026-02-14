package branch

import (
	"database/sql"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
)

type Branch struct {
	ID                string         `db:"id"`
	RepresentativeID  string         `db:"representative_id"`
	FranchiseID       sql.NullString `db:"franchise_id"`
	Name              string         `db:"name"`
	EstablishmentType string         `db:"establishment_type"`
	ProfileImageURL   sql.NullString `db:"profile_image_url"`
	Status            string         `db:"status"`
}

func (b *Branch) ToDomain() domain.Branch {
	branch := domain.Branch{
		ID:                b.ID,
		RepresentativeID:  b.RepresentativeID,
		Name:              b.Name,
		EstablishmentType: domain.EstablishmentType(b.EstablishmentType),
		Status:            domain.BranchStatus(b.Status),
	}
	if b.FranchiseID.Valid {
		franchiseID := b.FranchiseID.String
		branch.FranchiseID = &franchiseID
	}
	if b.ProfileImageURL.Valid {
		imageURL := b.ProfileImageURL.String
		branch.ProfileImageURL = &imageURL
	}
	return branch
}

func FromDomain(domainBranch domain.Branch) Branch {
	b := Branch{
		ID:                domainBranch.ID,
		RepresentativeID:  domainBranch.RepresentativeID,
		Name:              domainBranch.Name,
		EstablishmentType: string(domainBranch.EstablishmentType),
		Status:            string(domainBranch.Status),
	}
	if domainBranch.FranchiseID != nil {
		b.FranchiseID = sql.NullString{String: *domainBranch.FranchiseID, Valid: true}
	}
	if domainBranch.ProfileImageURL != nil {
		b.ProfileImageURL = sql.NullString{String: *domainBranch.ProfileImageURL, Valid: true}
	}
	return b
}

type Location struct {
	ID           sql.NullString  `db:"id"`
	DepartmentID sql.NullString  `db:"department_id"`
	CityID       sql.NullString  `db:"city_id"`
	Address      sql.NullString  `db:"address"`
	Latitude     sql.NullFloat64 `db:"latitude"`
	Longitude    sql.NullFloat64 `db:"longitude"`
}

func (l *Location) ToDomain(branchID string) *domain.Location {
	if !l.ID.Valid {
		return nil
	}
	loc := &domain.Location{
		ID:       l.ID.String,
		BranchID: branchID,
	}
	if l.DepartmentID.Valid {
		loc.DepartmentID = l.DepartmentID.String
	}
	if l.CityID.Valid {
		loc.CityID = l.CityID.String
	}
	if l.Address.Valid {
		loc.Address = l.Address.String
	}
	if l.Latitude.Valid {
		lat := l.Latitude.Float64
		loc.Latitude = &lat
	}
	if l.Longitude.Valid {
		lon := l.Longitude.Float64
		loc.Longitude = &lon
	}
	return loc
}
