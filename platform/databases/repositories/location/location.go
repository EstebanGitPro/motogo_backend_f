package location

import (
	"database/sql"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
)

// Location represents the database entity for branch locations
type Location struct {
	ID        string          `db:"id"`
	BranchID  string          `db:"branch_id"`
	CityID    string          `db:"city_id"`
	Address   string          `db:"address"`
	Latitude  sql.NullFloat64 `db:"latitude"`
	Longitude sql.NullFloat64 `db:"longitude"`
}

// ToDomain converts Location to domain.Location
func (l *Location) ToDomain() domain.Location {
	loc := domain.Location{
		ID:       l.ID,
		BranchID: l.BranchID,
		CityID:   l.CityID,
		Address:  l.Address,
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

// FromDomainLocation converts domain.Location to Location entity
func FromDomainLocation(domainLoc domain.Location) Location {
	loc := Location{
		ID:       domainLoc.ID,
		BranchID: domainLoc.BranchID,
		CityID:   domainLoc.CityID,
		Address:  domainLoc.Address,
	}
	if domainLoc.Latitude != nil {
		loc.Latitude = sql.NullFloat64{Float64: *domainLoc.Latitude, Valid: true}
	}
	if domainLoc.Longitude != nil {
		loc.Longitude = sql.NullFloat64{Float64: *domainLoc.Longitude, Valid: true}
	}
	return loc
}

// Department represents the database entity for departments
type Department struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

// ToDomain converts Department to domain.Department
func (d *Department) ToDomain() domain.Department {
	return domain.Department{
		ID:   d.ID,
		Name: d.Name,
	}
}

// City represents the database entity for cities
type City struct {
	ID           string `db:"id"`
	Name         string `db:"name"`
	DepartmentID string `db:"department_id"`
}

// ToDomain converts City to domain.City
func (c *City) ToDomain() domain.City {
	return domain.City{
		ID:           c.ID,
		Name:         c.Name,
		DepartmentID: c.DepartmentID,
	}
}
