package dependency

import (
	"context"
	"database/sql"
	"path/filepath"
	"time"

	"github.com/EstebanGitPro/motogo-backend/config"
	"github.com/EstebanGitPro/motogo-backend/core/interactor"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	messagingCache "github.com/EstebanGitPro/motogo-backend/platform/cache/messaging"
	"github.com/EstebanGitPro/motogo-backend/platform/firebase"
	"github.com/EstebanGitPro/motogo-backend/platform/geocoding"
	"github.com/EstebanGitPro/motogo-backend/platform/identity_provider/keycloak"
	"github.com/EstebanGitPro/motogo-backend/platform/jwt"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/EstebanGitPro/motogo-backend/tools/idencoder"
	"github.com/EstebanGitPro/motogo-backend/tools/utils"

	mysql "github.com/EstebanGitPro/motogo-backend/platform/databases/mysql"

	branchRepo "github.com/EstebanGitPro/motogo-backend/platform/databases/repositories/branch"
	brandRepo "github.com/EstebanGitPro/motogo-backend/platform/databases/repositories/brand"
	completedServiceRepo "github.com/EstebanGitPro/motogo-backend/platform/databases/repositories/completed_service" // HU64
	diagnosticRepo "github.com/EstebanGitPro/motogo-backend/platform/databases/repositories/diagnostic"              // HU11-14
	diagPermRepo "github.com/EstebanGitPro/motogo-backend/platform/databases/repositories/diagnostic_permission"
	evidenceRepo "github.com/EstebanGitPro/motogo-backend/platform/databases/repositories/evidence" // HU16-19
	franchiseRepo "github.com/EstebanGitPro/motogo-backend/platform/databases/repositories/franchise"
	locationRepo "github.com/EstebanGitPro/motogo-backend/platform/databases/repositories/location"
	messageRepo "github.com/EstebanGitPro/motogo-backend/platform/databases/repositories/message"
	motorcycleRepo "github.com/EstebanGitPro/motogo-backend/platform/databases/repositories/motorcycle"
	repo "github.com/EstebanGitPro/motogo-backend/platform/databases/repositories/person"
	ratingRepo "github.com/EstebanGitPro/motogo-backend/platform/databases/repositories/rating" // HU48
	scheduleRepo "github.com/EstebanGitPro/motogo-backend/platform/databases/repositories/schedule"
	scheduleDetailRepo "github.com/EstebanGitPro/motogo-backend/platform/databases/repositories/schedule_detail"
	serviceRepo "github.com/EstebanGitPro/motogo-backend/platform/databases/repositories/service"
)

type Dependencies struct {
	PersonService               input.Service
	PersonRepo                  output.Repository
	KeycloakClient              output.AuthClient
	Interactor                  *interactor.Interactor
	MessageInteractor           *interactor.MessageInteractor
	BranchInteractor            *interactor.BranchInteractor            // HU59
	BrandInteractor             *interactor.BrandInteractor             // Brands catalog
	LocationInteractor          *interactor.LocationInteractor          // Geographic catalogs
	ServiceInteractor           *interactor.ServiceInteractor           // Service catalog (HU63, HU75)
	FranchiseInteractor         *interactor.FranchiseInteractor         // HU26-29
	ScheduleInteractor          *interactor.ScheduleInteractor          // HU30-35
	ScheduleDetailInteractor    *interactor.ScheduleDetailInteractor    // HU6-9
	ScheduleExceptionInteractor *interactor.ScheduleExceptionInteractor // HU20-25
	MotorcycleInteractor        *interactor.MotorcycleInteractor        // HU43-47
	EvidenceInteractor          *interactor.EvidenceInteractor          // HU16-19
	DiagnosticInteractor        *interactor.DiagnosticInteractor        // HU11-14
	CompletedServiceInteractor  *interactor.CompletedServiceInteractor  // HU64
	RatingInteractor            *interactor.RatingInteractor            // HU48
	FirebaseClient              output.CustomTokenProvider              // Firebase Auth
	JWTValidator                output.JWTValidator                     // JWT validation with JWKS
	Config                      *config.Config
	Logger                      logger.Logger
	IDEncoder                   *idencoder.HashidsEncoder
	MessagingCache              *messagingCache.MessageCache
	ResponseHandler             *middleware.ResponseHandler
}

func Init() (*Dependencies, error) {
	log := logger.NewSlogLogger()
	log.Info(logger.LogAppStarting)

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error(logger.LogAppConfigError, "error", err)
		return nil, err
	}
	log.Info(logger.LogAppConfigLoaded)

	middleware.PrometheusInit()
	log.Success(logger.LogPrometheusInitOK)

	db, err := mysql.GetDB(cfg.Database)
	if err != nil {
		log.Error(logger.LogAppDatabaseError, "error", err)
		return nil, err
	}
	log.Success(logger.LogAppDatabaseConnected)

	keycloakClient, err := keycloak.NewClient(&cfg.Keycloak)
	if err != nil {
		log.Error(logger.LogKeycloakClientError, "error", err)
		return nil, err
	}
	log.Success(logger.LogKeycloakClientOK)

	personRepo, err := repo.NewClientRepository(db)
	if err != nil {
		log.Error(logger.LogPersonRepoInitError, "error", err)
		return nil, err
	}

	personService := services.NewService(personRepo, keycloakClient)

	interactorFacade := interactor.NewInteractor(personService)

	encoder, err := idencoder.NewHashidsEncoder(idencoder.Config{
		Secret:    cfg.IDEncoder.Secret,
		MinLength: cfg.IDEncoder.MinLength,
	})
	if err != nil {
		log.Error(logger.LogIDEncoderInitError, "error", err)
		return nil, err
	}
	log.Success(logger.LogIDEncodeOK)

	msgRepo, err := messageRepo.NewMessageRepository(db)
	if err != nil {
		log.Error(logger.LogRepoMsgInitError, "error", err)
		return nil, err
	}
	log.Success(logger.LogDependencyMessageRepoInit)

	refreshInterval := 5 * time.Minute
	messagingCache := messagingCache.NewMessageCache(msgRepo, refreshInterval)

	if err := messagingCache.LoadMessages(context.Background()); err != nil {
		log.Warn(logger.LogMsgCacheLoadError, "error", err)
	}
	log.Success(logger.LogMsgCacheInit, "messages_loaded", messagingCache.MessageCount())

	messagingCache.StartAutoRefresh(context.Background())

	responseHandler := middleware.NewResponseHandler(messagingCache)

	messageService := services.NewMessageService(msgRepo)
	messageInteractor := interactor.NewMessageInteractor(messageService)
	log.Success(logger.LogDependencyMessageIntInit)

	// Initialize all domain repositories and interactors
	repos, err := initRepositories(db, log)
	if err != nil {
		return nil, err
	}

	// Build services and interactors from repositories
	geocodingClient := initGeocodingClient(cfg, log)

	branchService := services.NewBranchService(repos.branch, repos.location, geocodingClient)
	branchInteractor := interactor.NewBranchInteractor(branchService)
	log.Success(logger.LogDepBranchInteractorInitOK)

	brandService := services.NewBrandService(repos.brand)
	brandInteractor := interactor.NewBrandInteractor(brandService)
	log.Success(logger.LogDepBrandInteractorInitOK)

	locationService := services.NewLocationService(repos.location)
	locationInteractor := interactor.NewLocationInteractor(locationService)
	log.Success(logger.LogDepLocationInteractorInitOK)

	serviceCatalogService := services.NewServiceCatalogService(repos.service)
	serviceInteractor := interactor.NewServiceInteractor(serviceCatalogService)
	log.Success(logger.LogDepServiceInteractorInitOK)

	franchiseService := services.NewFranchiseService(repos.franchise, repos.branch)
	franchiseInteractor := interactor.NewFranchiseInteractor(franchiseService)
	log.Success(logger.LogDepFranchiseInteractorInitOK)

	scheduleService := services.NewScheduleService(repos.schedule, repos.branch)
	scheduleInteractor := interactor.NewScheduleInteractor(scheduleService, branchService)
	log.Success(logger.LogDepScheduleIntInitOK)

	scheduleDetailService := services.NewScheduleDetailService(repos.scheduleDetail, repos.schedule)
	scheduleDetailInteractor := interactor.NewScheduleDetailInteractor(scheduleDetailService, scheduleService, branchService)
	log.Success(logger.LogDepScheduleDetailIntInitOK)

	scheduleExceptionInteractor := interactor.NewScheduleExceptionInteractor(scheduleDetailService, scheduleService, branchService)
	log.Success(logger.LogDepScheduleExceptionIntInitOK)

	motorcycleService := services.NewMotorcycleService(repos.motorcycle, repos.diagPerm)
	log.Success(logger.LogDepMotorcycleServiceInitOK)

	motorcycleInteractor := interactor.NewMotorcycleInteractor(motorcycleService)
	log.Success(logger.LogDepMotorcycleInteractorInitOK)

	evidenceService := services.NewEvidenceService(repos.evidence, repos.motorcycle)
	evidenceInteractor := interactor.NewEvidenceInteractor(evidenceService)
	log.Success(logger.LogDepEvidenceInteractorInitOK)

	diagnosticService := services.NewDiagnosticService(repos.diagnostic, repos.motorcycle, repos.branch)
	diagnosticInteractor := interactor.NewDiagnosticInteractor(diagnosticService)
	log.Success(logger.LogDepDiagnosticInteractorInitOK)

	// Completed Service (HU64)
	completedServiceRepository, err := completedServiceRepo.NewRepository(db)
	if err != nil {
		log.Error(logger.LogDepCSRepoInitErr, "error", err)
		return nil, err
	}
	log.Success(logger.LogDepCSRepoInitOK)

	completedServiceService := services.NewCompletedServiceService(completedServiceRepository, repos.diagnostic)
	completedServiceInteractor := interactor.NewCompletedServiceInteractor(completedServiceService)
	log.Success(logger.LogDepCSInteractorInitOK)

	// Rating (RELEASE_14 / HU48)
	ratingRepository, err := ratingRepo.NewRepository(db)
	if err != nil {
		log.Error(logger.LogDepCSRepoInitErr, "error", err)
		return nil, err
	}
	ratingService := services.NewRatingService(ratingRepository, completedServiceRepository)
	ratingInteractor := interactor.NewRatingInteractor(ratingService)
	log.Success(logger.LogDepCSInteractorInitOK)

	// Firebase and storage integration
	firebaseClient := initFirebaseIntegration(cfg, log, motorcycleService, evidenceService, branchInteractor)

	// JWT / JWKS validation
	jwtValidator := initJWTValidator(cfg, log)

	return &Dependencies{
		PersonService:               personService,
		PersonRepo:                  personRepo,
		KeycloakClient:              keycloakClient,
		Interactor:                  interactorFacade,
		MessageInteractor:           messageInteractor,
		BranchInteractor:            branchInteractor,
		BrandInteractor:             brandInteractor,
		LocationInteractor:          locationInteractor,
		ServiceInteractor:           serviceInteractor,
		FranchiseInteractor:         franchiseInteractor,
		ScheduleInteractor:          scheduleInteractor,
		ScheduleDetailInteractor:    scheduleDetailInteractor,
		ScheduleExceptionInteractor: scheduleExceptionInteractor,
		MotorcycleInteractor:        motorcycleInteractor,
		EvidenceInteractor:          evidenceInteractor,         // HU16-19
		DiagnosticInteractor:        diagnosticInteractor,       // HU11-14
		CompletedServiceInteractor:  completedServiceInteractor, // HU64
		RatingInteractor:            ratingInteractor,           // HU48
		FirebaseClient:              firebaseClient,
		JWTValidator:                jwtValidator,
		Config:                      cfg,
		Logger:                      log,
		IDEncoder:                   encoder,
		MessagingCache:              messagingCache,
		ResponseHandler:             responseHandler,
	}, nil
}

// domainRepositories holds all initialized domain repositories.
type domainRepositories struct {
	branch         output.BranchRepository
	brand          output.BrandRepository
	location       output.LocationRepository
	service        output.ServiceRepository
	franchise      output.FranchiseRepository
	schedule       output.ScheduleRepository
	scheduleDetail output.ScheduleDetailRepository
	motorcycle     output.MotorcycleRepository
	diagPerm       output.DiagnosticPermissionRepository
	evidence       output.EvidenceRepository
	diagnostic     output.DiagnosticRepository
}

// initRepositories initializes all database repositories.
func initRepositories(db *sql.DB, log logger.Logger) (*domainRepositories, error) {
	branchRepository, err := branchRepo.NewRepository(db)
	if err != nil {
		log.Error(logger.LogDepBranchRepoInitErr, "error", err)
		return nil, err
	}
	log.Success(logger.LogDepBranchRepoInitOK)

	locationRepository, err := locationRepo.NewRepository(db)
	if err != nil {
		log.Error(logger.LogDepLocationRepoInitErr, "error", err)
		return nil, err
	}
	log.Success(logger.LogDepLocationRepoInitOK)

	brandRepository, err := brandRepo.NewRepository(db)
	if err != nil {
		log.Error(logger.LogDepBrandRepoInitErr, "error", err)
		return nil, err
	}
	log.Success(logger.LogDepBrandRepoInitOK)

	serviceRepository, err := serviceRepo.NewRepository(db)
	if err != nil {
		log.Error(logger.LogDepServiceRepoInitErr, "error", err)
		return nil, err
	}
	log.Success(logger.LogDepServiceRepoInitOK)

	franchiseRepository, err := franchiseRepo.NewRepository(db)
	if err != nil {
		log.Error(logger.LogDepFranchiseRepoInitErr, "error", err)
		return nil, err
	}
	log.Success(logger.LogDepFranchiseRepoInitOK)

	scheduleRepository, err := scheduleRepo.NewRepository(db)
	if err != nil {
		log.Error(logger.LogDepScheduleRepoInitErr, "error", err)
		return nil, err
	}
	log.Success(logger.LogDepScheduleRepoInitOK)

	scheduleDetailRepository, err := scheduleDetailRepo.NewRepository(db)
	if err != nil {
		log.Error(logger.LogDepScheduleDetailRepoInitErr, "error", err)
		return nil, err
	}
	log.Success(logger.LogDepScheduleDetailRepoInitOK)

	motorcycleRepository, err := motorcycleRepo.NewRepository(db)
	if err != nil {
		log.Error(logger.LogDepMotorcycleRepoInitErr, "error", err)
		return nil, err
	}
	log.Success(logger.LogDepMotorcycleRepoInitOK)

	diagnosticPermissionRepository, err := diagPermRepo.NewRepository(db)
	if err != nil {
		log.Error(logger.LogDepDiagPermRepoInitErr, "error", err)
		return nil, err
	}
	log.Success(logger.LogDepDiagPermRepoInitOK)

	evidenceRepository, err := evidenceRepo.NewRepository(db)
	if err != nil {
		log.Error(logger.LogDepEvidenceRepoInitErr, "error", err)
		return nil, err
	}
	log.Success(logger.LogDepEvidenceRepoInitOK)

	diagnosticRepository, err := diagnosticRepo.NewRepository(db)
	if err != nil {
		log.Error(logger.LogDepDiagnosticRepoInitErr, "error", err)
		return nil, err
	}
	log.Success(logger.LogDepDiagnosticRepoInitOK)

	return &domainRepositories{
		branch:         branchRepository,
		brand:          brandRepository,
		location:       locationRepository,
		service:        serviceRepository,
		franchise:      franchiseRepository,
		schedule:       scheduleRepository,
		scheduleDetail: scheduleDetailRepository,
		motorcycle:     motorcycleRepository,
		diagPerm:       diagnosticPermissionRepository,
		evidence:       evidenceRepository,
		diagnostic:     diagnosticRepository,
	}, nil
}

// initFirebaseIntegration sets up Firebase client and wires storage to services that need it.
func initFirebaseIntegration(
	cfg *config.Config,
	log logger.Logger,
	motorcycleService *services.MotorcycleServiceImpl,
	evidenceService *services.EvidenceServiceImpl,
	branchInteractor *interactor.BranchInteractor,
) output.CustomTokenProvider {
	if cfg.Firebase.CredentialsPath == "" {
		log.Warn(logger.LogDepFirebaseCredNotConfig)
		return nil
	}

	firebaseCredPath := cfg.Firebase.CredentialsPath
	if !filepath.IsAbs(firebaseCredPath) {
		root, rootErr := utils.FindModuleRoot()
		if rootErr == nil {
			firebaseCredPath = filepath.Join(root, firebaseCredPath)
		}
	}
	log.Debug(logger.LogDepFirebaseCredPathResolved, "path", firebaseCredPath)

	firebaseClient, err := firebase.NewClient(firebaseCredPath)
	if err != nil {
		log.Warn(logger.LogDepFirebaseInitSkipped, "error", err)
		return nil
	}

	log.Success(logger.LogDepFirebaseClientInitOK)
	motorcycleService.WithStorageClient(firebaseClient)
	log.Success(logger.LogDepMotorcycleServiceInitOK, "with_storage", true)
	evidenceService.WithStorageClient(firebaseClient)
	log.Success(logger.LogDepEvidenceInteractorInitOK, "with_storage", true)
	branchInteractor.WithStorageClient(firebaseClient)
	log.Success(logger.LogDepBranchInteractorInitOK, "with_storage", true)

	return firebaseClient
}

// initJWTValidator creates a JWKS-based JWT validator from Keycloak config.
func initJWTValidator(cfg *config.Config, log logger.Logger) output.JWTValidator {
	jwtConfig := jwt.JWKSConfig{
		JWKSURL:         cfg.GetKeycloakJWKSURL(),
		Issuer:          cfg.GetKeycloakIssuerURL(),
		RefreshInterval: 15 * time.Minute,
	}
	jwtValidator, err := jwt.NewJWKSValidator(context.Background(), jwtConfig)
	if err != nil {
		log.Warn(logger.LogDepJWKSValidatorInitErr, "error", err)
		return nil
	}
	log.Success(logger.LogDepJWKSValidatorInitOK, "jwks_url", jwtConfig.JWKSURL)
	return jwtValidator
}

// initGeocodingClient creates a geocoding client based on config, with optional fallback provider.
func initGeocodingClient(cfg *config.Config, log logger.Logger) geocoding.Geocoder {
	timeout := time.Duration(cfg.Geocoding.TimeoutSeconds) * time.Second

	var primaryClient geocoding.Geocoder
	switch cfg.Geocoding.Provider {
	case "google":
		primaryClient = geocoding.NewGoogleMapsClient(
			cfg.Geocoding.APIKey,
			cfg.Geocoding.BaseURL,
			cfg.Geocoding.CountryCode,
			timeout,
			log,
		)
	case "mapbox":
		fallthrough
	default:
		primaryClient = geocoding.NewMapboxClient(
			cfg.Geocoding.APIKey,
			cfg.Geocoding.BaseURL,
			cfg.Geocoding.CountryCode,
			timeout,
		)
	}

	if cfg.Geocoding.FallbackProvider == "" {
		log.Success(logger.LogDepGeocodingClientInitOK, "provider", cfg.Geocoding.Provider)
		return primaryClient
	}

	var fallbackClient geocoding.Geocoder
	switch cfg.Geocoding.FallbackProvider {
	case "mapbox":
		fallbackClient = geocoding.NewMapboxClient(
			cfg.Geocoding.FallbackAPIKey,
			cfg.Geocoding.FallbackBaseURL,
			cfg.Geocoding.CountryCode,
			timeout,
		)
	case "google":
		fallbackClient = geocoding.NewGoogleMapsClient(
			cfg.Geocoding.FallbackAPIKey,
			cfg.Geocoding.FallbackBaseURL,
			cfg.Geocoding.CountryCode,
			timeout,
			log,
		)
	}

	if fallbackClient != nil {
		log.Success(logger.LogDepGeocodingClientInitOK,
			"primary", cfg.Geocoding.Provider,
			"fallback", cfg.Geocoding.FallbackProvider)
		return geocoding.NewFallbackClient(primaryClient, fallbackClient)
	}

	log.Success(logger.LogDepGeocodingClientInitOK, "provider", cfg.Geocoding.Provider)
	return primaryClient
}
