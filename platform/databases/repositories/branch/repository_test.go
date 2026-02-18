package branch

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/stretchr/testify/assert"
)

// ============================================
// NewRepository Tests
// ============================================

func TestNewRepository_NilDB(t *testing.T) {
	repo, err := NewRepository(nil)

	assert.Nil(t, repo)
	assert.Equal(t, sql.ErrConnDone, err)
}

func TestNewRepository_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Expect all prepared statements
	mock.ExpectPrepare("INSERT INTO branches")
	mock.ExpectPrepare("UPDATE branches")
	mock.ExpectPrepare("DELETE FROM branches")
	mock.ExpectPrepare("SELECT b.id, b.representative_id.*WHERE b.id")
	mock.ExpectPrepare("SELECT id, representative_id.*WHERE franchise_id")
	mock.ExpectPrepare("SELECT b.id, b.representative_id.*WHERE b.representative_id")
	mock.ExpectPrepare("INSERT INTO branch_brands")
	mock.ExpectPrepare("DELETE FROM branch_brands")
	mock.ExpectPrepare("SELECT brand_id FROM branch_brands")
	mock.ExpectPrepare("INSERT INTO branch_displacement_ranges")
	mock.ExpectPrepare("DELETE FROM branch_displacement_ranges")
	mock.ExpectPrepare("SELECT displacement_range FROM branch_displacement_ranges")

	repo, err := NewRepository(db)

	assert.NoError(t, err)
	assert.NotNil(t, repo)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNewRepository_PrepareError_SaveBranch(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO branches").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtSaveBranch")
}

func TestNewRepository_PrepareError_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO branches")
	mock.ExpectPrepare("UPDATE branches").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtUpdateBranch")
}

func TestNewRepository_PrepareError_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO branches")
	mock.ExpectPrepare("UPDATE branches")
	mock.ExpectPrepare("DELETE FROM branches").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtDeleteBranch")
}

func TestNewRepository_PrepareError_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO branches")
	mock.ExpectPrepare("UPDATE branches")
	mock.ExpectPrepare("DELETE FROM branches")
	mock.ExpectPrepare("SELECT.*FROM branches.*LEFT JOIN locations").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtGetBranchByID")
}

func TestNewRepository_PrepareError_GetByFranchiseAndName(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO branches")
	mock.ExpectPrepare("UPDATE branches")
	mock.ExpectPrepare("DELETE FROM branches")
	mock.ExpectPrepare("SELECT.*FROM branches.*LEFT JOIN locations")
	mock.ExpectPrepare("SELECT.*FROM branches.*WHERE franchise_id").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtGetBranchByFranchiseAndName")
}

// ============================================
// BeginTx Tests
// ============================================

func TestBeginTx_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()

	repo := &repository{db: db}

	tx, err := repo.BeginTx(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, tx)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBeginTx_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin().WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}

	tx, err := repo.BeginTx(context.Background())

	assert.Nil(t, tx)
	assert.Equal(t, sql.ErrConnDone, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ============================================
// GetBranchByID Tests
// ============================================

func TestGetBranchByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT b.id, b.representative_id")
	stmt.ExpectQuery().
		WithArgs("not-found").
		WillReturnError(sql.ErrNoRows)

	repo := &repository{db: db}
	repo.stmtGetBranchByID, _ = db.Prepare("SELECT b.id, b.representative_id FROM branches b WHERE b.id = ?")

	branch, err := repo.GetBranchByID(context.Background(), "not-found")

	assert.Nil(t, branch)
	assert.Equal(t, domain.ErrBranchNotFound, err)
}

func TestGetBranchByID_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT b.id, b.representative_id")
	stmt.ExpectQuery().
		WithArgs("branch-error").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetBranchByID, _ = db.Prepare("SELECT b.id, b.representative_id FROM branches b WHERE b.id = ?")

	branch, err := repo.GetBranchByID(context.Background(), "branch-error")

	assert.Nil(t, branch)
	assert.Error(t, err)
}

// ============================================
// HasBranchesByRepresentative Tests
// ============================================

func TestHasBranchesByRepresentative_True(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs("rep-123").
		WillReturnRows(rows)

	repo := &repository{db: db}

	exists, err := repo.HasBranchesByRepresentative(context.Background(), "rep-123")

	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestHasBranchesByRepresentative_False(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"exists"}).AddRow(false)

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs("rep-no-branches").
		WillReturnRows(rows)

	repo := &repository{db: db}

	exists, err := repo.HasBranchesByRepresentative(context.Background(), "rep-no-branches")

	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestHasBranchesByRepresentative_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs("rep-error").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}

	exists, err := repo.HasBranchesByRepresentative(context.Background(), "rep-error")

	assert.False(t, exists)
	assert.Error(t, err)
}

// ============================================
// ValidateBrands Tests
// ============================================

func TestValidateBrands_EmptySlice(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	err = repo.ValidateBrands(context.Background(), []string{})

	assert.NoError(t, err)
}

func TestValidateBrands_AllFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow("brand-001").
		AddRow("brand-002")

	mock.ExpectQuery("SELECT id FROM brands WHERE id IN").
		WithArgs("brand-001", "brand-002").
		WillReturnRows(rows)

	repo := &repository{db: db}

	err = repo.ValidateBrands(context.Background(), []string{"brand-001", "brand-002"})

	assert.NoError(t, err)
}

func TestValidateBrands_SomeNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow("brand-001")

	mock.ExpectQuery("SELECT id FROM brands WHERE id IN").
		WithArgs("brand-001", "brand-notfound").
		WillReturnRows(rows)

	repo := &repository{db: db}

	err = repo.ValidateBrands(context.Background(), []string{"brand-001", "brand-notfound"})

	assert.Equal(t, domain.ErrBrandNotFound, err)
}

func TestValidateBrands_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT id FROM brands WHERE id IN").
		WithArgs("brand-001").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}

	err = repo.ValidateBrands(context.Background(), []string{"brand-001"})

	assert.Error(t, err)
}

// ============================================
// GetBranchByFranchiseAndName Tests
// ============================================

func TestGetBranchByFranchiseAndName_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"id", "representative_id", "franchise_id", "name", "establishment_type", "profile_image_url", "status",
	}).AddRow("branch-001", "rep-123", "franchise-001", "Sucursal Norte", "WORKSHOP", "http://example.com/img.jpg", "ACTIVE")

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("franchise-001", "Sucursal Norte").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetBranchByFranchiseAndName, _ = db.Prepare("SELECT * FROM branches WHERE franchise_id = ? AND name = ?")

	branch, err := repo.GetBranchByFranchiseAndName(context.Background(), "franchise-001", "Sucursal Norte")

	assert.NoError(t, err)
	assert.NotNil(t, branch)
	assert.Equal(t, "branch-001", branch.ID)
	assert.Equal(t, "Sucursal Norte", branch.Name)
}

func TestGetBranchByFranchiseAndName_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("franchise-001", "No Existe").
		WillReturnError(sql.ErrNoRows)

	repo := &repository{db: db}
	repo.stmtGetBranchByFranchiseAndName, _ = db.Prepare("SELECT * FROM branches WHERE franchise_id = ? AND name = ?")

	branch, err := repo.GetBranchByFranchiseAndName(context.Background(), "franchise-001", "No Existe")

	assert.Nil(t, branch)
	assert.Equal(t, domain.ErrBranchNotFound, err)
}

func TestGetBranchByFranchiseAndName_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("franchise-error", "Test").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetBranchByFranchiseAndName, _ = db.Prepare("SELECT * FROM branches WHERE franchise_id = ? AND name = ?")

	branch, err := repo.GetBranchByFranchiseAndName(context.Background(), "franchise-error", "Test")

	assert.Nil(t, branch)
	assert.Error(t, err)
}

// ============================================
// GetBranchesByRepresentative Tests
// ============================================

func TestGetBranchesByRepresentative_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"id", "representative_id", "franchise_id", "name", "establishment_type", "profile_image_url", "status",
		"location_id", "city_id", "address", "latitude", "longitude", "department_id", "phone_number",
	})

	stmt := mock.ExpectPrepare("SELECT b.id, b.representative_id")
	stmt.ExpectQuery().
		WithArgs("rep-no-branches").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetBranchesByRepresentative, _ = db.Prepare("SELECT b.id, b.representative_id FROM branches WHERE b.representative_id = ?")

	branches, err := repo.GetBranchesByRepresentative(context.Background(), "rep-no-branches")

	assert.NoError(t, err)
	assert.Empty(t, branches)
}

func TestGetBranchesByRepresentative_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT b.id, b.representative_id")
	stmt.ExpectQuery().
		WithArgs("rep-error").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetBranchesByRepresentative, _ = db.Prepare("SELECT b.id, b.representative_id FROM branches WHERE b.representative_id = ?")

	branches, err := repo.GetBranchesByRepresentative(context.Background(), "rep-error")

	assert.Nil(t, branches)
	assert.Error(t, err)
}

// ============================================
// GetBranchesNearby Tests
// ============================================

func TestGetBranchesNearby_Success(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	assert.NoError(t, err)
	defer db.Close()

	// 1. Expect prepared statements for hydration
	stmtBrands := mock.ExpectPrepare("SELECT brand_id FROM branch_brands")
	stmtDR := mock.ExpectPrepare("SELECT displacement_range FROM branch_displacement_ranges")

	// 2. Set up repo with real prepared statements
	repo := &repository{db: db}
	repo.stmtGetBranchBrands, _ = db.Prepare("SELECT brand_id FROM branch_brands WHERE branch_id = ?")
	repo.stmtGetBranchDisplacementRanges, _ = db.Prepare("SELECT displacement_range FROM branch_displacement_ranges WHERE branch_id = ?")

	// 3. Expect main nearby query
	rows := sqlmock.NewRows([]string{
		"id", "name", "establishment_type", "profile_image_url", "status",
		"address", "latitude", "longitude", "city_name", "department_name", "phone_number", "distance_km",
	}).AddRow(
		"branch-001", "Taller Central", "WORKSHOP", "http://example.com/img.jpg", "ACTIVE",
		"Calle 123 #45-67", 4.7110, -74.0721, "Bogotá", "Cundinamarca", "3001234567", 1.5,
	).AddRow(
		"branch-002", "Tienda Norte", "STORE", nil, "ACTIVE",
		"Carrera 15 #80-45", 4.7200, -74.0650, "Bogotá", "Cundinamarca", nil, 2.3,
	)
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	// 4. Expect hydration queries for branch-001
	stmtBrands.ExpectQuery().WithArgs("branch-001").
		WillReturnRows(sqlmock.NewRows([]string{"brand_id"}).AddRow("brand-A").AddRow("brand-B"))
	stmtDR.ExpectQuery().WithArgs("branch-001").
		WillReturnRows(sqlmock.NewRows([]string{"displacement_range"}).AddRow("BAJO").AddRow("MEDIO"))

	// 5. Expect hydration queries for branch-002
	stmtBrands.ExpectQuery().WithArgs("branch-002").
		WillReturnRows(sqlmock.NewRows([]string{"brand_id"}))
	stmtDR.ExpectQuery().WithArgs("branch-002").
		WillReturnRows(sqlmock.NewRows([]string{"displacement_range"}).AddRow("ALTO"))

	branches, err := repo.GetBranchesNearby(
		context.Background(),
		domain.NearbySearchParams{
			Lat: 4.7110, Lng: -74.0721, RadiusKm: 5.0,
			LatMin: 4.6110, LatMax: 4.8110, LngMin: -74.1721, LngMax: -73.9721,
		},
	)

	assert.NoError(t, err)
	assert.Len(t, branches, 2)
	assert.Equal(t, "branch-001", branches[0].ID)
	assert.Equal(t, 1.5, branches[0].DistanceKm)
	assert.NotNil(t, branches[0].ProfileImageURL)
	assert.Nil(t, branches[1].ProfileImageURL)
	// Verify hydrated data
	assert.Equal(t, []string{"brand-A", "brand-B"}, branches[0].Brands)
	assert.Equal(t, []domain.DisplacementRange{"BAJO", "MEDIO"}, branches[0].DisplacementRanges)
	assert.Empty(t, branches[1].Brands)
	assert.Equal(t, []domain.DisplacementRange{"ALTO"}, branches[1].DisplacementRanges)
}

func TestGetBranchesNearby_Empty(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"id", "name", "establishment_type", "profile_image_url", "status",
		"address", "latitude", "longitude", "city_name", "department_name", "phone_number", "distance_km",
	})

	mock.ExpectQuery("SELECT").
		WillReturnRows(rows)

	repo := &repository{db: db}

	branches, err := repo.GetBranchesNearby(
		context.Background(),
		domain.NearbySearchParams{
			Lat: 4.7110, Lng: -74.0721, RadiusKm: 1.0,
			EstablishmentType: "WORKSHOP",
			LatMin:            4.6110, LatMax: 4.8110, LngMin: -74.1721, LngMax: -73.9721,
		},
	)

	assert.NoError(t, err)
	assert.Empty(t, branches)
}

func TestGetBranchesNearby_DBError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}

	branches, err := repo.GetBranchesNearby(
		context.Background(),
		domain.NearbySearchParams{
			Lat: 4.7110, Lng: -74.0721, RadiusKm: 5.0,
			LatMin: 4.6110, LatMax: 4.8110, LngMin: -74.1721, LngMax: -73.9721,
		},
	)

	assert.Nil(t, branches)
	assert.Error(t, err)
}

// ============================================
// getBranchBrands Tests
// ============================================

func TestGetBranchBrands_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"brand_id"}).
		AddRow("brand-001").
		AddRow("brand-002").
		AddRow("brand-003")

	stmt := mock.ExpectPrepare("SELECT brand_id FROM branch_brands")
	stmt.ExpectQuery().
		WithArgs("branch-001").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetBranchBrands, _ = db.Prepare("SELECT brand_id FROM branch_brands WHERE branch_id = ?")

	brands, err := repo.getBranchBrands(context.Background(), "branch-001")

	assert.NoError(t, err)
	assert.Len(t, brands, 3)
	assert.Equal(t, "brand-001", brands[0])
}

func TestGetBranchBrands_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"brand_id"})

	stmt := mock.ExpectPrepare("SELECT brand_id FROM branch_brands")
	stmt.ExpectQuery().
		WithArgs("branch-no-brands").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetBranchBrands, _ = db.Prepare("SELECT brand_id FROM branch_brands WHERE branch_id = ?")

	brands, err := repo.getBranchBrands(context.Background(), "branch-no-brands")

	assert.NoError(t, err)
	assert.Empty(t, brands)
}

func TestGetBranchBrands_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT brand_id FROM branch_brands")
	stmt.ExpectQuery().
		WithArgs("branch-error").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetBranchBrands, _ = db.Prepare("SELECT brand_id FROM branch_brands WHERE branch_id = ?")

	brands, err := repo.getBranchBrands(context.Background(), "branch-error")

	assert.Nil(t, brands)
	assert.Error(t, err)
}

// ============================================
// Branch Model Tests
// ============================================

func TestBranch_ToDomain(t *testing.T) {
	b := &Branch{
		ID:                "branch-001",
		RepresentativeID:  "rep-123",
		FranchiseID:       sql.NullString{String: "franchise-001", Valid: true},
		Name:              "Sucursal Norte",
		EstablishmentType: "WORKSHOP",
		ProfileImageURL:   sql.NullString{String: "http://example.com/img.jpg", Valid: true},
		Status:            "ACTIVE",
	}

	domain := b.ToDomain()

	assert.Equal(t, "branch-001", domain.ID)
	assert.Equal(t, "rep-123", domain.RepresentativeID)
	assert.NotNil(t, domain.FranchiseID)
	assert.Equal(t, "franchise-001", *domain.FranchiseID)
	assert.NotNil(t, domain.ProfileImageURL)
	assert.Equal(t, "http://example.com/img.jpg", *domain.ProfileImageURL)
}

func TestBranch_ToDomain_NullFields(t *testing.T) {
	b := &Branch{
		ID:                "branch-002",
		RepresentativeID:  "rep-456",
		FranchiseID:       sql.NullString{Valid: false},
		Name:              "Sucursal Sin Franquicia",
		EstablishmentType: "STORE",
		ProfileImageURL:   sql.NullString{Valid: false},
		Status:            "INACTIVE",
	}

	domain := b.ToDomain()

	assert.Equal(t, "branch-002", domain.ID)
	assert.Nil(t, domain.FranchiseID)
	assert.Nil(t, domain.ProfileImageURL)
}

func TestFromDomain(t *testing.T) {
	franchiseID := "franchise-001"
	imageURL := "http://example.com/img.jpg"
	domainBranch := domain.Branch{
		ID:                "branch-001",
		RepresentativeID:  "rep-123",
		FranchiseID:       &franchiseID,
		Name:              "Sucursal Norte",
		EstablishmentType: "WORKSHOP",
		ProfileImageURL:   &imageURL,
		Status:            "ACTIVE",
	}

	b := FromDomain(domainBranch)

	assert.Equal(t, "branch-001", b.ID)
	assert.True(t, b.FranchiseID.Valid)
	assert.Equal(t, "franchise-001", b.FranchiseID.String)
	assert.True(t, b.ProfileImageURL.Valid)
}

func TestFromDomain_NilFields(t *testing.T) {
	domainBranch := domain.Branch{
		ID:                "branch-002",
		RepresentativeID:  "rep-456",
		FranchiseID:       nil,
		Name:              "Sucursal Sin Franquicia",
		EstablishmentType: "STORE",
		ProfileImageURL:   nil,
		Status:            "INACTIVE",
	}

	b := FromDomain(domainBranch)

	assert.Equal(t, "branch-002", b.ID)
	assert.False(t, b.FranchiseID.Valid)
	assert.False(t, b.ProfileImageURL.Valid)
}

// ============================================
// Location Model Tests
// ============================================

func TestLocation_ToDomain(t *testing.T) {
	lat := 4.7110
	lng := -74.0721
	l := &Location{
		ID:           sql.NullString{String: "loc-001", Valid: true},
		DepartmentID: sql.NullString{String: "dept-001", Valid: true},
		CityID:       sql.NullString{String: "city-001", Valid: true},
		Address:      sql.NullString{String: "Calle 123 #45-67", Valid: true},
		Latitude:     sql.NullFloat64{Float64: lat, Valid: true},
		Longitude:    sql.NullFloat64{Float64: lng, Valid: true},
	}

	domainLoc := l.ToDomain("branch-001")

	assert.NotNil(t, domainLoc)
	assert.Equal(t, "loc-001", domainLoc.ID)
	assert.Equal(t, "branch-001", domainLoc.BranchID)
	assert.Equal(t, "dept-001", domainLoc.DepartmentID)
	assert.Equal(t, "city-001", domainLoc.CityID)
	assert.Equal(t, "Calle 123 #45-67", domainLoc.Address)
	assert.NotNil(t, domainLoc.Latitude)
	assert.Equal(t, lat, *domainLoc.Latitude)
	assert.NotNil(t, domainLoc.Longitude)
	assert.Equal(t, lng, *domainLoc.Longitude)
}

func TestLocation_ToDomain_InvalidID(t *testing.T) {
	l := &Location{
		ID: sql.NullString{Valid: false},
	}

	domainLoc := l.ToDomain("branch-001")

	assert.Nil(t, domainLoc)
}

func TestLocation_ToDomain_PartialData(t *testing.T) {
	l := &Location{
		ID:           sql.NullString{String: "loc-002", Valid: true},
		DepartmentID: sql.NullString{Valid: false},
		CityID:       sql.NullString{Valid: false},
		Address:      sql.NullString{Valid: false},
		Latitude:     sql.NullFloat64{Valid: false},
		Longitude:    sql.NullFloat64{Valid: false},
	}

	domainLoc := l.ToDomain("branch-002")

	assert.NotNil(t, domainLoc)
	assert.Equal(t, "loc-002", domainLoc.ID)
	assert.Empty(t, domainLoc.DepartmentID)
	assert.Empty(t, domainLoc.CityID)
	assert.Nil(t, domainLoc.Latitude)
	assert.Nil(t, domainLoc.Longitude)
}

// ============================================
// SaveBranch Tests
// ============================================

func TestSaveBranch_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO branches").
		WithArgs("branch-123", "rep-123", nil, "Test Branch", "WORKSHOP", nil, "ACTIVE").
		WillReturnResult(sqlmock.NewResult(1, 1))

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	branch := domain.Branch{ID: "branch-123", RepresentativeID: "rep-123", Name: "Test Branch", EstablishmentType: "WORKSHOP", Status: "ACTIVE"}
	err = repo.SaveBranch(context.Background(), sqlTx, branch)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSaveBranch_DuplicateError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO branches").WillReturnError(sql.ErrConnDone)

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	err = repo.SaveBranch(context.Background(), sqlTx, domain.Branch{})
	assert.Equal(t, domain.ErrBranchCannotSave, err)
}

func TestSaveBranch_InvalidTransaction(t *testing.T) {
	repo := &repository{}
	err := repo.SaveBranch(context.Background(), nil, domain.Branch{})
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

// ============================================
// UpdateBranch Tests
// ============================================

func TestUpdateBranch_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE branches").
		WithArgs(nil, "Updated Branch", "STORE", nil, "INACTIVE", "branch-123").
		WillReturnResult(sqlmock.NewResult(0, 1))

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	branch := domain.Branch{ID: "branch-123", Name: "Updated Branch", EstablishmentType: "STORE", Status: "INACTIVE"}
	err = repo.UpdateBranch(context.Background(), sqlTx, branch)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateBranch_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE branches").WillReturnError(sql.ErrConnDone)

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	err = repo.UpdateBranch(context.Background(), sqlTx, domain.Branch{})
	assert.Equal(t, domain.ErrBranchCannotUpdate, err)
}

func TestUpdateBranch_InvalidTransaction(t *testing.T) {
	repo := &repository{}
	err := repo.UpdateBranch(context.Background(), nil, domain.Branch{})
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

// ============================================
// DeleteBranch Tests
// ============================================

func TestDeleteBranch_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM branches").
		WithArgs("branch-123").
		WillReturnResult(sqlmock.NewResult(0, 1))

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	err = repo.DeleteBranch(context.Background(), sqlTx, "branch-123")

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteBranch_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM branches").WillReturnError(sql.ErrConnDone)

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	err = repo.DeleteBranch(context.Background(), sqlTx, "branch-err")
	assert.Equal(t, domain.ErrBranchCannotDelete, err)
}

func TestDeleteBranch_InvalidTransaction(t *testing.T) {
	repo := &repository{}
	err := repo.DeleteBranch(context.Background(), nil, "branch-123")
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

// ============================================
// SaveBranchBrands Tests
// ============================================

func TestSaveBranchBrands_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO branch_brands").
		WithArgs(sqlmock.AnyArg(), "branch-123", "brand-1").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO branch_brands").
		WithArgs(sqlmock.AnyArg(), "branch-123", "brand-2").
		WillReturnResult(sqlmock.NewResult(1, 1))

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	err = repo.SaveBranchBrands(context.Background(), sqlTx, "branch-123", []string{"brand-1", "brand-2"})

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSaveBranchBrands_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO branch_brands").WillReturnError(sql.ErrConnDone)

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	err = repo.SaveBranchBrands(context.Background(), sqlTx, "branch-123", []string{"brand-1"})
	assert.Equal(t, domain.ErrBranchCannotSave, err)
}

func TestSaveBranchBrands_InvalidTransaction(t *testing.T) {
	repo := &repository{}
	err := repo.SaveBranchBrands(context.Background(), nil, "branch-123", []string{"brand-1"})
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

// ============================================
// DeleteBranchBrands Tests
// ============================================

func TestDeleteBranchBrands_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM branch_brands").
		WithArgs("branch-123").
		WillReturnResult(sqlmock.NewResult(0, 3))

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	err = repo.DeleteBranchBrands(context.Background(), sqlTx, "branch-123")

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteBranchBrands_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM branch_brands").WillReturnError(sql.ErrConnDone)

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	err = repo.DeleteBranchBrands(context.Background(), sqlTx, "branch-err")
	assert.Equal(t, domain.ErrBranchCannotDelete, err)
}

func TestDeleteBranchBrands_InvalidTransaction(t *testing.T) {
	repo := &repository{}
	err := repo.DeleteBranchBrands(context.Background(), nil, "branch-123")
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

// ============================================
// SaveBranchDisplacementRanges Tests
// ============================================

func TestSaveBranchDisplacementRanges_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO branch_displacement_ranges").
		WithArgs(sqlmock.AnyArg(), "branch-123", "BAJO").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO branch_displacement_ranges").
		WithArgs(sqlmock.AnyArg(), "branch-123", "MEDIO").
		WillReturnResult(sqlmock.NewResult(1, 1))

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	err = repo.SaveBranchDisplacementRanges(context.Background(), sqlTx, "branch-123", []string{"BAJO", "MEDIO"})

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSaveBranchDisplacementRanges_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO branch_displacement_ranges").WillReturnError(sql.ErrConnDone)

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	err = repo.SaveBranchDisplacementRanges(context.Background(), sqlTx, "branch-123", []string{"BAJO"})
	assert.Equal(t, domain.ErrBranchCannotSave, err)
}

func TestSaveBranchDisplacementRanges_InvalidTransaction(t *testing.T) {
	repo := &repository{}
	err := repo.SaveBranchDisplacementRanges(context.Background(), nil, "branch-123", []string{"BAJO"})
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

func TestSaveBranchDisplacementRanges_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	err = repo.SaveBranchDisplacementRanges(context.Background(), sqlTx, "branch-123", []string{})

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ============================================
// DeleteBranchDisplacementRanges Tests
// ============================================

func TestDeleteBranchDisplacementRanges_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM branch_displacement_ranges").
		WithArgs("branch-123").
		WillReturnResult(sqlmock.NewResult(0, 3))

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	err = repo.DeleteBranchDisplacementRanges(context.Background(), sqlTx, "branch-123")

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteBranchDisplacementRanges_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM branch_displacement_ranges").WillReturnError(sql.ErrConnDone)

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	err = repo.DeleteBranchDisplacementRanges(context.Background(), sqlTx, "branch-err")
	assert.Equal(t, domain.ErrBranchCannotDelete, err)
}

func TestDeleteBranchDisplacementRanges_InvalidTransaction(t *testing.T) {
	repo := &repository{}
	err := repo.DeleteBranchDisplacementRanges(context.Background(), nil, "branch-123")
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

// ============================================
// GetBranchByID Success Path Tests
// ============================================

func TestGetBranchByID_Success_WithAllFields(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// 1. Prepare hydration statements FIRST (they are used after main query)
	stmtBrands := mock.ExpectPrepare("SELECT brand_id FROM branch_brands")
	stmtDR := mock.ExpectPrepare("SELECT displacement_range FROM branch_displacement_ranges")

	// 2. Prepare main query
	stmt := mock.ExpectPrepare("SELECT b.id, b.representative_id")

	// 3. Set up repo with real prepared statements
	repo := &repository{db: db}
	repo.stmtGetBranchBrands, _ = db.Prepare("SELECT brand_id FROM branch_brands WHERE branch_id = ?")
	repo.stmtGetBranchDisplacementRanges, _ = db.Prepare("SELECT displacement_range FROM branch_displacement_ranges WHERE branch_id = ?")
	repo.stmtGetBranchByID, _ = db.Prepare("SELECT b.id, b.representative_id FROM branches b WHERE b.id = ?")

	// 4. Main query returns branch with all nullable fields populated
	rows := sqlmock.NewRows([]string{
		"id", "representative_id", "franchise_id", "name", "establishment_type", "profile_image_url", "status",
		"location_id", "city_id", "address", "latitude", "longitude", "department_id", "phone_number",
	}).AddRow(
		"branch-001", "rep-123", "franchise-001", "Taller Central", "WORKSHOP", "http://example.com/img.jpg", "ACTIVE",
		"loc-001", "city-001", "Calle 123 #45-67", 4.7110, -74.0721, "dept-001", "3001234567",
	)
	stmt.ExpectQuery().
		WithArgs("branch-001").
		WillReturnRows(rows)

	// 5. Hydration queries (executed after main query)
	stmtBrands.ExpectQuery().
		WithArgs("branch-001").
		WillReturnRows(sqlmock.NewRows([]string{"brand_id"}).AddRow("brand-A").AddRow("brand-B"))
	stmtDR.ExpectQuery().
		WithArgs("branch-001").
		WillReturnRows(sqlmock.NewRows([]string{"displacement_range"}).AddRow("BAJO").AddRow("MEDIO"))

	branch, err := repo.GetBranchByID(context.Background(), "branch-001")

	assert.NoError(t, err)
	assert.NotNil(t, branch)
	assert.Equal(t, "branch-001", branch.ID)
	assert.Equal(t, "rep-123", branch.RepresentativeID)
	assert.NotNil(t, branch.FranchiseID)
	assert.Equal(t, "franchise-001", *branch.FranchiseID)
	assert.NotNil(t, branch.ProfileImageURL)
	assert.Equal(t, "http://example.com/img.jpg", *branch.ProfileImageURL)
	assert.NotNil(t, branch.RepresentativePhone)
	assert.Equal(t, "3001234567", *branch.RepresentativePhone)
	assert.NotNil(t, branch.Location)
	assert.Equal(t, "loc-001", branch.Location.ID)
	assert.Equal(t, "dept-001", branch.Location.DepartmentID)
	assert.Equal(t, "city-001", branch.Location.CityID)
	assert.NotNil(t, branch.Location.Latitude)
	assert.NotNil(t, branch.Location.Longitude)
	assert.Equal(t, []string{"brand-A", "brand-B"}, branch.Brands)
	assert.Equal(t, []domain.DisplacementRange{"BAJO", "MEDIO"}, branch.DisplacementRanges)
}

func TestGetBranchByID_Success_NullableFieldsNil(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// 1. Prepare hydration statements FIRST
	stmtBrands := mock.ExpectPrepare("SELECT brand_id FROM branch_brands")
	stmtDR := mock.ExpectPrepare("SELECT displacement_range FROM branch_displacement_ranges")

	// 2. Prepare main query
	stmt := mock.ExpectPrepare("SELECT b.id, b.representative_id")

	// 3. Set up repo
	repo := &repository{db: db}
	repo.stmtGetBranchBrands, _ = db.Prepare("SELECT brand_id FROM branch_brands WHERE branch_id = ?")
	repo.stmtGetBranchDisplacementRanges, _ = db.Prepare("SELECT displacement_range FROM branch_displacement_ranges WHERE branch_id = ?")
	repo.stmtGetBranchByID, _ = db.Prepare("SELECT b.id, b.representative_id FROM branches b WHERE b.id = ?")

	// 4. Main query with nullable fields as NULL
	rows := sqlmock.NewRows([]string{
		"id", "representative_id", "franchise_id", "name", "establishment_type", "profile_image_url", "status",
		"location_id", "city_id", "address", "latitude", "longitude", "department_id", "phone_number",
	}).AddRow(
		"branch-002", "rep-456", nil, "Tienda Sin Extra", "STORE", nil, "ACTIVE",
		nil, nil, nil, nil, nil, nil, nil,
	)
	stmt.ExpectQuery().
		WithArgs("branch-002").
		WillReturnRows(rows)

	// 5. Hydration queries (empty results)
	stmtBrands.ExpectQuery().
		WithArgs("branch-002").
		WillReturnRows(sqlmock.NewRows([]string{"brand_id"}))
	stmtDR.ExpectQuery().
		WithArgs("branch-002").
		WillReturnRows(sqlmock.NewRows([]string{"displacement_range"}))

	branch, err := repo.GetBranchByID(context.Background(), "branch-002")

	assert.NoError(t, err)
	assert.NotNil(t, branch)
	assert.Nil(t, branch.FranchiseID)
	assert.Nil(t, branch.ProfileImageURL)
	assert.Nil(t, branch.RepresentativePhone)
	assert.Nil(t, branch.Location)
	assert.Empty(t, branch.Brands)
}

// ============================================
// GetBranchesByRepresentative Success Path Tests
// ============================================

func TestGetBranchesByRepresentative_Success_WithData(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	assert.NoError(t, err)
	defer db.Close()

	// 1. Prepare hydration stmts on repo
	mock.ExpectPrepare("SELECT brand_id FROM branch_brands")
	mock.ExpectPrepare("SELECT displacement_range FROM branch_displacement_ranges")

	repo := &repository{db: db}
	repo.stmtGetBranchBrands, _ = db.Prepare("SELECT brand_id FROM branch_brands WHERE branch_id = ?")
	repo.stmtGetBranchDisplacementRanges, _ = db.Prepare("SELECT displacement_range FROM branch_displacement_ranges WHERE branch_id = ?")

	// 2. Prepare main query
	stmtMain := mock.ExpectPrepare("SELECT b.id, b.representative_id")
	repo.stmtGetBranchesByRepresentative, _ = db.Prepare("SELECT b.id, b.representative_id FROM branches b WHERE b.representative_id = ?")

	// 3. Main query returns 2 branches
	rows := sqlmock.NewRows([]string{
		"id", "representative_id", "franchise_id", "name", "establishment_type", "profile_image_url", "status",
		"location_id", "city_id", "address", "latitude", "longitude", "department_id", "phone_number",
	}).AddRow(
		"branch-001", "rep-123", "franchise-001", "Taller Norte", "WORKSHOP", "http://example.com/img.jpg", "ACTIVE",
		"loc-001", "city-001", "Calle 123", 4.7110, -74.0721, "dept-001", "3001234567",
	).AddRow(
		"branch-002", "rep-123", nil, "Tienda Sur", "STORE", nil, "ACTIVE",
		nil, nil, nil, nil, nil, nil, nil,
	)
	stmtMain.ExpectQuery().WillReturnRows(rows)

	// 4. During rows iteration, Go re-prepares hydration stmts on new connections
	reBrands := mock.ExpectPrepare("SELECT brand_id FROM branch_brands")
	reBrands.ExpectQuery().WillReturnRows(sqlmock.NewRows([]string{"brand_id"}).AddRow("brand-A"))
	reDR := mock.ExpectPrepare("SELECT displacement_range FROM branch_displacement_ranges")
	reDR.ExpectQuery().WillReturnRows(sqlmock.NewRows([]string{"displacement_range"}).AddRow("BAJO"))

	branches, err := repo.GetBranchesByRepresentative(context.Background(), "rep-123")

	assert.NoError(t, err)
	assert.NotEmpty(t, branches)

	// First branch: all fields populated
	assert.Equal(t, "branch-001", branches[0].ID)
	assert.Equal(t, "rep-123", branches[0].RepresentativeID)
	assert.Equal(t, "Taller Norte", branches[0].Name)
	assert.NotNil(t, branches[0].FranchiseID)
	assert.Equal(t, "franchise-001", *branches[0].FranchiseID)
	assert.NotNil(t, branches[0].ProfileImageURL)
	assert.NotNil(t, branches[0].RepresentativePhone)
	assert.NotNil(t, branches[0].Location)
	assert.Equal(t, "loc-001", branches[0].Location.ID)
	assert.Equal(t, []string{"brand-A"}, branches[0].Brands)
	assert.Equal(t, []domain.DisplacementRange{"BAJO"}, branches[0].DisplacementRanges)
}
