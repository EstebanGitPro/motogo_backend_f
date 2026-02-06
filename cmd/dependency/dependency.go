package dependency

import (
	"context"
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
	evidenceRepo "github.com/EstebanGitPro/motogo-backend/platform/databases/repositories/evidence" // HU16-19
	franchiseRepo "github.com/EstebanGitPro/motogo-backend/platform/databases/repositories/franchise"
	locationRepo "github.com/EstebanGitPro/motogo-backend/platform/databases/repositories/location"
	messageRepo "github.com/EstebanGitPro/motogo-backend/platform/databases/repositories/message"
	motorcycleRepo "github.com/EstebanGitPro/motogo-backend/platform/databases/repositories/motorcycle"
	repo "github.com/EstebanGitPro/motogo-backend/platform/databases/repositories/person"
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
	FirebaseClient              *firebase.Client                        // Firebase Auth
	JWTValidator                *jwt.JWKSValidator                      // JWT validation with JWKS
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

	branchRepository, err := branchRepo.NewRepository(db)
	if err != nil {
		log.Error(logger.LogDepBranchRepoInitErr, "error", err)
		return nil, err
	}
	log.Success(logger.LogDepBranchRepoInitOK)

	timeout := time.Duration(cfg.Geocoding.TimeoutSeconds) * time.Second

	var primaryClient geocoding.Client
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

	var geocodingClient geocoding.Client
	if cfg.Geocoding.FallbackProvider != "" {
		var fallbackClient geocoding.Client
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
			geocodingClient = geocoding.NewFallbackClient(primaryClient, fallbackClient)
			log.Success(logger.LogDepGeocodingClientInitOK,
				"primary", cfg.Geocoding.Provider,
				"fallback", cfg.Geocoding.FallbackProvider)
		} else {
			geocodingClient = primaryClient
			log.Success(logger.LogDepGeocodingClientInitOK, "provider", cfg.Geocoding.Provider)
		}
	} else {
		geocodingClient = primaryClient
		log.Success(logger.LogDepGeocodingClientInitOK, "provider", cfg.Geocoding.Provider)
	}

	locationRepository, err := locationRepo.NewRepository(db)
	if err != nil {
		log.Error(logger.LogDepLocationRepoInitErr, "error", err)
		return nil, err
	}
	log.Success(logger.LogDepLocationRepoInitOK)

	branchService := services.NewBranchService(branchRepository, locationRepository, geocodingClient)
	branchInteractor := interactor.NewBranchInteractor(branchService)
	log.Success(logger.LogDepBranchInteractorInitOK)

	brandRepository, err := brandRepo.NewRepository(db)
	if err != nil {
		log.Error(logger.LogDepBrandRepoInitErr, "error", err)
		return nil, err
	}
	log.Success(logger.LogDepBrandRepoInitOK)

	brandService := services.NewBrandService(brandRepository)
	brandInteractor := interactor.NewBrandInteractor(brandService)
	log.Success(logger.LogDepBrandInteractorInitOK)

	locationService := services.NewLocationService(locationRepository)
	locationInteractor := interactor.NewLocationInteractor(locationService)
	log.Success(logger.LogDepLocationInteractorInitOK)

	serviceRepository, err := serviceRepo.NewRepository(db)
	if err != nil {
		log.Error(logger.LogDepServiceRepoInitErr, "error", err)
		return nil, err
	}
	log.Success(logger.LogDepServiceRepoInitOK)

	serviceCatalogService := services.NewServiceCatalogService(serviceRepository)
	serviceInteractor := interactor.NewServiceInteractor(serviceCatalogService)
	log.Success(logger.LogDepServiceInteractorInitOK)

	franchiseRepository, err := franchiseRepo.NewRepository(db)
	if err != nil {
		log.Error(logger.LogDepFranchiseRepoInitErr, "error", err)
		return nil, err
	}
	log.Success(logger.LogDepFranchiseRepoInitOK)

	franchiseService := services.NewFranchiseService(franchiseRepository)
	franchiseInteractor := interactor.NewFranchiseInteractor(franchiseService, branchService)
	log.Success(logger.LogDepFranchiseInteractorInitOK)

	scheduleRepository, err := scheduleRepo.NewRepository(db)
	if err != nil {
		log.Error(logger.LogDepScheduleRepoInitErr, "error", err)
		return nil, err
	}
	log.Success(logger.LogDepScheduleRepoInitOK)

	scheduleService := services.NewScheduleService(scheduleRepository, branchRepository)
	scheduleInteractor := interactor.NewScheduleInteractor(scheduleService, branchService)
	log.Success(logger.LogDepScheduleIntInitOK)

	scheduleDetailRepository, err := scheduleDetailRepo.NewRepository(db)
	if err != nil {
		log.Error(logger.LogDepScheduleDetailRepoInitErr, "error", err)
		return nil, err
	}
	log.Success(logger.LogDepScheduleDetailRepoInitOK)

	scheduleDetailService := services.NewScheduleDetailService(scheduleDetailRepository, scheduleRepository)
	scheduleDetailInteractor := interactor.NewScheduleDetailInteractor(scheduleDetailService, scheduleService, branchService)
	log.Success(logger.LogDepScheduleDetailIntInitOK)

	scheduleExceptionInteractor := interactor.NewScheduleExceptionInteractor(scheduleDetailService, scheduleService, branchService)
	log.Success(logger.LogDepScheduleExceptionIntInitOK)

	motorcycleRepository, err := motorcycleRepo.NewRepository(db)
	if err != nil {
		log.Error(logger.LogDepMotorcycleRepoInitErr, "error", err)
		return nil, err
	}
	log.Success(logger.LogDepMotorcycleRepoInitOK)

	// Create motorcycle interactor - will connect Firebase Storage later if available
	motorcycleInteractor := interactor.NewMotorcycleInteractor(motorcycleRepository)
	log.Success(logger.LogDepMotorcycleInteractorInitOK)

	// Evidence feature (HU16-19)
	evidenceRepository, err := evidenceRepo.NewRepository(db)
	if err != nil {
		log.Error(logger.LogDepEvidenceRepoInitErr, "error", err)
		return nil, err
	}
	log.Success(logger.LogDepEvidenceRepoInitOK)

	evidenceInteractor := interactor.NewEvidenceInteractor(evidenceRepository, motorcycleRepository)
	log.Success(logger.LogDepEvidenceInteractorInitOK)

	var firebaseClient *firebase.Client
	if cfg.Firebase.CredentialsPath != "" {
		firebaseCredPath := cfg.Firebase.CredentialsPath
		if !filepath.IsAbs(firebaseCredPath) {
			root, rootErr := utils.FindModuleRoot()
			if rootErr == nil {
				firebaseCredPath = filepath.Join(root, firebaseCredPath)
			}
		}
		log.Debug(logger.LogDepFirebaseCredPathResolved, "path", firebaseCredPath)
		firebaseClient, err = firebase.NewClient(firebaseCredPath)
		if err != nil {
			log.Warn(logger.LogDepFirebaseInitSkipped, "error", err)
			// Don't fail startup if Firebase is not configured
		} else {
			log.Success(logger.LogDepFirebaseClientInitOK)
			// Connect Firebase Storage to MotorcycleInteractor for image deletion (HU45)
			motorcycleInteractor.WithStorageClient(firebaseClient)
			log.Success(logger.LogDepMotorcycleInteractorInitOK, "with_storage", true)
			// Connect Firebase Storage to EvidenceInteractor for evidence deletion (HU19)
			evidenceInteractor.WithStorageClient(firebaseClient)
			log.Success(logger.LogDepEvidenceInteractorInitOK, "with_storage", true)
			// Connect Firebase Storage to BranchInteractor for profile image deletion (HU60-61)
			branchInteractor.WithStorageClient(firebaseClient)
			log.Success(logger.LogDepBranchInteractorInitOK, "with_storage", true)
		}
	} else {
		log.Warn(logger.LogDepFirebaseCredNotConfig)
	}

	var jwtValidator *jwt.JWKSValidator
	jwtConfig := jwt.JWKSConfig{
		JWKSURL:         cfg.GetKeycloakJWKSURL(),
		Issuer:          cfg.GetKeycloakIssuerURL(),
		RefreshInterval: 15 * time.Minute, // Refresh keys every 15 minutes
	}
	jwtValidator, err = jwt.NewJWKSValidator(context.Background(), jwtConfig)
	if err != nil {
		log.Warn(logger.LogDepJWKSValidatorInitErr, "error", err)
		jwtValidator = nil
	} else {
		log.Success(logger.LogDepJWKSValidatorInitOK, "jwks_url", jwtConfig.JWKSURL)
	}

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
		EvidenceInteractor:          evidenceInteractor, // HU16-19
		FirebaseClient:              firebaseClient,
		JWTValidator:                jwtValidator,
		Config:                      cfg,
		Logger:                      log,
		IDEncoder:                   encoder,
		MessagingCache:              messagingCache,
		ResponseHandler:             responseHandler,
	}, nil
}
