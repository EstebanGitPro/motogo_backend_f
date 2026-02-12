package handlers

import (
	"github.com/EstebanGitPro/motogo-backend/core/interactor"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	messagingCache "github.com/EstebanGitPro/motogo-backend/platform/cache/messaging"
	"github.com/EstebanGitPro/motogo-backend/tools/idencoder"
)

// NewForTest creates a handler with interface-based interactors for integration testing.
// This allows passing mock implementations directly without wrapping in concrete types.
func NewForTest(
	brandInteractor input.BrandInteractorInterface,
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
