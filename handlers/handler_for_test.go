package handlers

import (
	"github.com/EstebanGitPro/motogo-backend/core/interactor"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	messagingCache "github.com/EstebanGitPro/motogo-backend/platform/cache/messaging"
	"github.com/EstebanGitPro/motogo-backend/tools/idencoder"
	"github.com/gin-gonic/gin"
)

// NewForTest creates a handler with interface-based interactors for integration testing.
// This allows passing mock implementations directly without wrapping in concrete types.
func NewForTest(
	brandInteractor input.BrandLister,
	locationInteractor input.LocationInteractorInterface,
	motorcycleInteractor input.MotorcycleInteractorInterface,
	evidenceInteractor input.EvidenceInteractorInterface,
	messageCache *messagingCache.MessageCache,
	encoder *idencoder.HashidsEncoder,
	responseHandler *middleware.ResponseHandler,
) *handler {
	return &handler{
		BrandInteractor:      brandInteractor,
		LocationInteractor:   locationInteractor,
		MotorcycleInteractor: motorcycleInteractor,
		EvidenceInteractor:   evidenceInteractor,
		MessagingCache:       messageCache,
		IDEncoder:            encoder,
		Response:             responseHandler,
	}
}

// NewForTestWithConcrete creates a handler with concrete interactors for integration testing
// of controllers that use concrete types (Branch, Service, Diagnostic, Franchise).
func NewForTestWithConcrete(
	branchInteractor *interactor.BranchInteractor,
	serviceInteractor *interactor.ServiceInteractor,
	diagnosticInteractor *interactor.DiagnosticInteractor,
	franchiseInteractor *interactor.FranchiseInteractor,
	encoder *idencoder.HashidsEncoder,
	responseHandler *middleware.ResponseHandler,
) *handler {
	return &handler{
		BranchInteractor:     branchInteractor,
		ServiceInteractor:    serviceInteractor,
		DiagnosticInteractor: diagnosticInteractor,
		FranchiseInteractor:  franchiseInteractor,
		IDEncoder:            encoder,
		Response:             responseHandler,
	}
}

// NewForTestWithFirebase creates a handler with a CustomTokenProvider for integration testing
// of controllers that require Firebase authentication (GetFirebaseToken).
func NewForTestWithFirebase(
	firebaseClient output.CustomTokenProvider,
	encoder *idencoder.HashidsEncoder,
	responseHandler *middleware.ResponseHandler,
) *handler {
	return &handler{
		FirebaseClient: firebaseClient,
		IDEncoder:      encoder,
		Response:       responseHandler,
	}
}

// NewForTestWithMessage creates a handler with MessageInteractor for integration testing
// of admin message controllers (CreateMessage, UpdateMessage, etc.).
func NewForTestWithMessage(
	messageInteractor *interactor.MessageInteractor,
	encoder *idencoder.HashidsEncoder,
	responseHandler *middleware.ResponseHandler,
) *handler {
	return &handler{
		MessageInteractor: messageInteractor,
		IDEncoder:         encoder,
		Response:          responseHandler,
	}
}

// NewForTestWithPerson creates a handler with Interactor and optional BranchInteractor
// for integration testing of Person controllers (Register, Login, UpdateProfile, etc.).
func NewForTestWithPerson(
	personInteractor *interactor.Interactor,
	branchInteractor *interactor.BranchInteractor,
	encoder *idencoder.HashidsEncoder,
	responseHandler *middleware.ResponseHandler,
) *handler {
	return &handler{
		Interactor:       personInteractor,
		BranchInteractor: branchInteractor,
		IDEncoder:        encoder,
		Response:         responseHandler,
	}
}

// NewForTestWithLookup creates a handler with all interactors needed for
// LookupMotorcycleByPlate integration tests (Branch, Motorcycle, Diagnostic, Evidence).
func NewForTestWithLookup(
	branchInteractor *interactor.BranchInteractor,
	motorcycleInteractor input.MotorcycleInteractorInterface,
	diagnosticInteractor *interactor.DiagnosticInteractor,
	evidenceInteractor input.EvidenceInteractorInterface,
	encoder *idencoder.HashidsEncoder,
	responseHandler *middleware.ResponseHandler,
) *handler {
	return &handler{
		BranchInteractor:     branchInteractor,
		MotorcycleInteractor: motorcycleInteractor,
		DiagnosticInteractor: diagnosticInteractor,
		EvidenceInteractor:   evidenceInteractor,
		IDEncoder:            encoder,
		Response:             responseHandler,
	}
}

// === Exported wrappers for error mapper testing ===

// MapRegisterBranchError exposes mapRegisterBranchError for testing.
func (h *handler) MapRegisterBranchError(c *gin.Context, err error) {
	h.mapRegisterBranchError(c, err)
}

// MapUpdateBranchError exposes mapUpdateBranchError for testing.
func (h *handler) MapUpdateBranchError(c *gin.Context, err error) {
	h.mapUpdateBranchError(c, err)
}

// MapMotorcycleRegError exposes mapMotorcycleRegError for testing.
func (h *handler) MapMotorcycleRegError(c *gin.Context, err error) {
	h.mapMotorcycleRegError(c, err)
}

// MapMotorcycleUpdateError exposes mapMotorcycleUpdateError for testing.
func (h *handler) MapMotorcycleUpdateError(c *gin.Context, err error) {
	h.mapMotorcycleUpdateError(c, err)
}

// MapUpdateProfileError exposes mapUpdateProfileError for testing.
func (h *handler) MapUpdateProfileError(c *gin.Context, err error) {
	h.mapUpdateProfileError(c, err)
}

// MapScheduleUpdateError exposes mapScheduleUpdateError for testing.
func (h *handler) MapScheduleUpdateError(c *gin.Context, err error) {
	h.mapScheduleUpdateError(c, err)
}

// MapDetailCreationError exposes mapDetailCreationError for testing.
func (h *handler) MapDetailCreationError(c *gin.Context, err error) {
	h.mapDetailCreationError(c, err)
}

// MapExceptionCreationError exposes mapExceptionCreationError for testing.
func (h *handler) MapExceptionCreationError(c *gin.Context, err error) {
	h.mapExceptionCreationError(c, err)
}
