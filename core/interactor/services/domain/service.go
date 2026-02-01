package domain

import (
	uuid "github.com/EstebanGitPro/motogo-backend/tools/utils"
)

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

type Service struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	ServiceType ServiceType `json:"service_type"`
	IsActive    bool        `json:"is_active"`
}

func IsValidServiceType(s string) bool {
	switch ServiceType(s) {
	case ServiceTypeMaintenance, ServiceTypeRepair, ServiceTypeTires,
		ServiceTypeDiagnostics, ServiceTypeAesthetics, ServiceTypeAccessories,
		ServiceTypeElectrical, ServiceTypeLegal:
		return true
	}
	return false
}

type BranchServiceInfo struct {
	Service Service `json:"service"`
	AddedAt string  `json:"added_at"`
	Active  bool    `json:"active"`
}

func GenerateUUID() string {
	return uuid.Generate()
}
