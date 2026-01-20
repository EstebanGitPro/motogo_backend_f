package domain

// ServiceType represents the type of service (user-facing, in Spanish)
type ServiceType string

const (
	ServiceTypeMaintenance ServiceType = "Mantenimiento"
	ServiceTypeRepair      ServiceType = "Reparación"
	ServiceTypeTires       ServiceType = "Llantas"
	ServiceTypeDiagnostics ServiceType = "Diagnóstico"
	ServiceTypeAesthetics  ServiceType = "Estética"
	ServiceTypeAccessories ServiceType = "Accesorios"
	ServiceTypeElectrical  ServiceType = "Eléctrico"
	ServiceTypeLegal       ServiceType = "Legal"
)

// AllServiceTypes returns all available service types
func AllServiceTypes() []ServiceType {
	return []ServiceType{
		ServiceTypeMaintenance,
		ServiceTypeRepair,
		ServiceTypeTires,
		ServiceTypeDiagnostics,
		ServiceTypeAesthetics,
		ServiceTypeAccessories,
		ServiceTypeElectrical,
		ServiceTypeLegal,
	}
}

// Service represents a service from the catalog
type Service struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	ServiceType ServiceType `json:"service_type"`
	IsActive    bool        `json:"is_active"`
}

// IsValidServiceType checks if the provided string is a valid service type
func IsValidServiceType(s string) bool {
	switch ServiceType(s) {
	case ServiceTypeMaintenance, ServiceTypeRepair, ServiceTypeTires,
		ServiceTypeDiagnostics, ServiceTypeAesthetics, ServiceTypeAccessories,
		ServiceTypeElectrical, ServiceTypeLegal:
		return true
	}
	return false
}

// BranchServiceInfo represents a service associated with a specific branch
// Includes when the service was added to that particular branch
type BranchServiceInfo struct {
	Service Service `json:"service"`  // The service from the catalog
	AddedAt string  `json:"added_at"` // When the service was added to this branch (ISO 8601)
	Active  bool    `json:"active"`   // Whether the service is active at this branch
}
