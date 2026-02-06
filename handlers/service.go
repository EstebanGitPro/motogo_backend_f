package handlers

import (
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/tools/idencoder"
)

// ServiceTypeItemResponse represents a single service type
type ServiceTypeItemResponse struct {
	Value string `json:"value"`
}

// ServiceTypeListResponse represents the response for GET /service-types
type ServiceTypeListResponse struct {
	Types []ServiceTypeItemResponse `json:"types"`
	Links []Link                    `json:"_links"`
}

// ServiceItemResponse represents a single service in the catalog
type ServiceItemResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	ServiceType string `json:"service_type"`
	IsActive    bool   `json:"is_active"`
}

// ServiceListResponse represents the response for GET /services
type ServiceListResponse struct {
	Services []ServiceItemResponse `json:"services"`
	Links    []Link                `json:"_links"`
}

// NewServiceTypeListResponse creates a ServiceTypeListResponse from domain service types
func NewServiceTypeListResponse(types []domain.ServiceType, links []Link) ServiceTypeListResponse {
	items := make([]ServiceTypeItemResponse, len(types))
	for i, t := range types {
		items[i] = ServiceTypeItemResponse{
			Value: string(t),
		}
	}
	return ServiceTypeListResponse{
		Types: items,
		Links: links,
	}
}

// NewServiceListResponse creates a ServiceListResponse from domain services (without ID encoding)
func NewServiceListResponse(services []domain.Service, links []Link) ServiceListResponse {
	items := make([]ServiceItemResponse, len(services))
	for i, svc := range services {
		items[i] = ServiceItemResponse{
			ID:          svc.ID,
			Name:        svc.Name,
			Description: svc.Description,
			ServiceType: string(svc.ServiceType),
			IsActive:    svc.IsActive,
		}
	}
	return ServiceListResponse{
		Services: items,
		Links:    links,
	}
}

// NewServiceListResponseWithEncoder creates a ServiceListResponse with encoded IDs
func NewServiceListResponseWithEncoder(services []domain.Service, links []Link, encoder *idencoder.HashidsEncoder) (ServiceListResponse, error) {
	items := make([]ServiceItemResponse, len(services))
	for i, svc := range services {
		encodedID, err := encoder.Encode(svc.ID)
		if err != nil {
			return ServiceListResponse{}, err
		}
		items[i] = ServiceItemResponse{
			ID:          encodedID,
			Name:        svc.Name,
			Description: svc.Description,
			ServiceType: string(svc.ServiceType),
			IsActive:    svc.IsActive,
		}
	}
	return ServiceListResponse{
		Services: items,
		Links:    links,
	}, nil
}

// BranchServiceItemResponse represents a service associated with a branch
type BranchServiceItemResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	ServiceType string `json:"service_type"`
	AddedAt     string `json:"added_at"` // ISO 8601 - when the service was added to this branch
	Active      bool   `json:"active"`
}

// BranchServiceListResponse represents the response for GET /branches/:id/services
type BranchServiceListResponse struct {
	Services []BranchServiceItemResponse `json:"services"`
	Links    []Link                      `json:"_links"`
}

// NewBranchServiceListResponse creates a BranchServiceListResponse from domain BranchServiceInfo
func NewBranchServiceListResponse(services []domain.BranchServiceInfo, links []Link) BranchServiceListResponse {
	items := make([]BranchServiceItemResponse, len(services))
	for i, info := range services {
		items[i] = BranchServiceItemResponse{
			ID:          info.Service.ID,
			Name:        info.Service.Name,
			Description: info.Service.Description,
			ServiceType: string(info.Service.ServiceType),
			AddedAt:     info.AddedAt,
			Active:      info.Active,
		}
	}
	return BranchServiceListResponse{
		Services: items,
		Links:    links,
	}
}

// NewBranchServiceListResponseWithEncoder creates a BranchServiceListResponse with encoded IDs
func NewBranchServiceListResponseWithEncoder(services []domain.BranchServiceInfo, links []Link, encoder *idencoder.HashidsEncoder) (BranchServiceListResponse, error) {
	items := make([]BranchServiceItemResponse, len(services))
	for i, info := range services {
		encodedID, err := encoder.Encode(info.Service.ID)
		if err != nil {
			return BranchServiceListResponse{}, err
		}
		items[i] = BranchServiceItemResponse{
			ID:          encodedID,
			Name:        info.Service.Name,
			Description: info.Service.Description,
			ServiceType: string(info.Service.ServiceType),
			AddedAt:     info.AddedAt,
			Active:      info.Active,
		}
	}
	return BranchServiceListResponse{
		Services: items,
		Links:    links,
	}, nil
}

// UpdateServiceRequest represents the request body for PUT /admin/services/:id (HU68)
type UpdateServiceRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	ServiceType string `json:"service_type" binding:"required"`
	IsActive    *bool  `json:"is_active"` // Optional: activate/deactivate service
}

// Sanitize trims whitespace from all string fields
func (r *UpdateServiceRequest) Sanitize() {
	r.Name = TrimString(r.Name)
	r.Description = TrimString(r.Description)
	r.ServiceType = TrimString(r.ServiceType)
}

// ServiceDetailResponse represents a single service response (HU68)
type ServiceDetailResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	ServiceType string `json:"service_type"`
	IsActive    bool   `json:"is_active"`
	Links       []Link `json:"_links"`
}

// NewServiceDetailResponse creates a ServiceDetailResponse from domain service
func NewServiceDetailResponse(service *domain.Service, links []Link) ServiceDetailResponse {
	return ServiceDetailResponse{
		ID:          service.ID,
		Name:        service.Name,
		Description: service.Description,
		ServiceType: string(service.ServiceType),
		Links:       links,
	}
}

// NewServiceDetailResponseWithEncoder creates a ServiceDetailResponse with encoded ID
func NewServiceDetailResponseWithEncoder(service *domain.Service, links []Link, encoder *idencoder.HashidsEncoder) (ServiceDetailResponse, error) {
	encodedID, err := encoder.Encode(service.ID)
	if err != nil {
		return ServiceDetailResponse{}, err
	}
	return ServiceDetailResponse{
		ID:          encodedID,
		Name:        service.Name,
		Description: service.Description,
		ServiceType: string(service.ServiceType),
		IsActive:    service.IsActive,
		Links:       links,
	}, nil
}
