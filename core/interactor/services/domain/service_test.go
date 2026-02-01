package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ============================================
// ServiceType Tests
// ============================================

func TestServiceType_Constants(t *testing.T) {
	assert.Equal(t, ServiceType("Mantenimiento"), ServiceTypeMaintenance)
	assert.Equal(t, ServiceType("Reparación"), ServiceTypeRepair)
	assert.Equal(t, ServiceType("Llantas"), ServiceTypeTires)
	assert.Equal(t, ServiceType("Diagnóstico"), ServiceTypeDiagnostics)
	assert.Equal(t, ServiceType("Estética"), ServiceTypeAesthetics)
	assert.Equal(t, ServiceType("Accesorios"), ServiceTypeAccessories)
	assert.Equal(t, ServiceType("Eléctrico"), ServiceTypeElectrical)
	assert.Equal(t, ServiceType("Legal"), ServiceTypeLegal)
}

func TestAllServiceTypes(t *testing.T) {
	types := AllServiceTypes()

	assert.Len(t, types, 8)
	assert.Contains(t, types, ServiceTypeMaintenance)
	assert.Contains(t, types, ServiceTypeRepair)
	assert.Contains(t, types, ServiceTypeTires)
	assert.Contains(t, types, ServiceTypeDiagnostics)
	assert.Contains(t, types, ServiceTypeAesthetics)
	assert.Contains(t, types, ServiceTypeAccessories)
	assert.Contains(t, types, ServiceTypeElectrical)
	assert.Contains(t, types, ServiceTypeLegal)
}

func TestIsValidServiceType_Valid(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"Mantenimiento", true},
		{"Reparación", true},
		{"Llantas", true},
		{"Diagnóstico", true},
		{"Estética", true},
		{"Accesorios", true},
		{"Eléctrico", true},
		{"Legal", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := IsValidServiceType(tt.input)
			assert.True(t, result)
		})
	}
}

func TestIsValidServiceType_Invalid(t *testing.T) {
	tests := []string{
		"INVALID",
		"",
		"maintenance",
		"MAINTENANCE",
		"Other",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			result := IsValidServiceType(input)
			assert.False(t, result)
		})
	}
}

// ============================================
// Service Tests
// ============================================

func TestService_Structure(t *testing.T) {
	service := Service{
		ID:          "svc-001",
		Name:        "Cambio de aceite",
		Description: "Servicio básico de cambio de aceite",
		ServiceType: ServiceTypeMaintenance,
		IsActive:    true,
	}

	assert.Equal(t, "svc-001", service.ID)
	assert.Equal(t, "Cambio de aceite", service.Name)
	assert.Equal(t, "Servicio básico de cambio de aceite", service.Description)
	assert.Equal(t, ServiceTypeMaintenance, service.ServiceType)
	assert.True(t, service.IsActive)
}

func TestService_EmptyDescription(t *testing.T) {
	service := Service{
		ID:          "svc-002",
		Name:        "Lavado",
		ServiceType: ServiceTypeAesthetics,
		IsActive:    false,
	}

	assert.Empty(t, service.Description)
	assert.False(t, service.IsActive)
}

// ============================================
// BranchServiceInfo Tests
// ============================================

func TestBranchServiceInfo_Structure(t *testing.T) {
	info := BranchServiceInfo{
		Service: Service{
			ID:          "svc-001",
			Name:        "Revisión",
			ServiceType: ServiceTypeDiagnostics,
			IsActive:    true,
		},
		AddedAt: "2024-01-15T10:30:00Z",
		Active:  true,
	}

	assert.Equal(t, "svc-001", info.Service.ID)
	assert.Equal(t, "Revisión", info.Service.Name)
	assert.Equal(t, "2024-01-15T10:30:00Z", info.AddedAt)
	assert.True(t, info.Active)
}

// ============================================
// GenerateUUID Tests
// ============================================

func TestGenerateUUID_ReturnsValidUUID(t *testing.T) {
	uuid := GenerateUUID()

	assert.NotEmpty(t, uuid)
	assert.Len(t, uuid, 36) // UUID format
}

func TestGenerateUUID_ReturnsUniqueValues(t *testing.T) {
	uuid1 := GenerateUUID()
	uuid2 := GenerateUUID()

	assert.NotEqual(t, uuid1, uuid2)
}
