package branch

import (
	"database/sql"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
)

// Branch represents the database entity for branches
type Branch struct {
	ID                string         `db:"id"`
	RepresentativeID  string         `db:"representative_id"`
	FranchiseID       sql.NullString `db:"franchise_id"`
	Name              string         `db:"name"`
	EstablishmentType string         `db:"establishment_type"`
	ProfileImageURL   sql.NullString `db:"profile_image_url"`
	Status            string         `db:"status"`
}

// ToDomain converts Branch to domain.Branch
func (b *Branch) ToDomain() domain.Branch {
	branch := domain.Branch{
		ID:                b.ID,
		RepresentativeID:  b.RepresentativeID,
		Name:              b.Name,
		EstablishmentType: b.EstablishmentType,
		Status:            b.Status,
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

// FromDomain converts domain.Branch to Branch entity
func FromDomain(domainBranch domain.Branch) Branch {
	b := Branch{
		ID:                domainBranch.ID,
		RepresentativeID:  domainBranch.RepresentativeID,
		Name:              domainBranch.Name,
		EstablishmentType: domainBranch.EstablishmentType,
		Status:            domainBranch.Status,
	}
	if domainBranch.FranchiseID != nil {
		b.FranchiseID = sql.NullString{String: *domainBranch.FranchiseID, Valid: true}
	}
	if domainBranch.ProfileImageURL != nil {
		b.ProfileImageURL = sql.NullString{String: *domainBranch.ProfileImageURL, Valid: true}
	}
	return b
}

// Location represents the database entity for branch locations
type Location struct {
	ID           sql.NullString  `db:"id"`
	DepartmentID sql.NullString  `db:"department_id"`
	CityID       sql.NullString  `db:"city_id"`
	Address      sql.NullString  `db:"address"`
	Latitude     sql.NullFloat64 `db:"latitude"`
	Longitude    sql.NullFloat64 `db:"longitude"`
}

// ToDomain converts Location to domain.Location pointer (nil if no location)
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
