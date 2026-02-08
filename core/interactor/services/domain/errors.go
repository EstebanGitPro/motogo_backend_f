package domain

import "errors"

// Domain errors - Used as keys for error-to-message mapping
// The actual user-facing messages come from cache/database
// These are only used internally for error identification and logging

// User Management Errors (MOD_U_*)
var (
	ErrDuplicateUser             = errors.New("ERR_DUPLICATE_USER")
	ErrUserCannotSave            = errors.New("ERR_USER_CANNOT_SAVE")
	ErrPersonNotFound            = errors.New("ERR_PERSON_NOT_FOUND")
	ErrUserNotFound              = errors.New("ERR_USER_NOT_FOUND")
	ErrEmailAlreadyVerified      = errors.New("ERR_EMAIL_ALREADY_VERIFIED")
	ErrInvalidTransaction        = errors.New("ERR_INVALID_TRANSACTION")
	ErrGettingUserByEmail        = errors.New("ERR_GETTING_USER_BY_EMAIL")
	ErrNotFoundUserByEmail       = errors.New("ERR_NOT_FOUND_USER_BY_EMAIL")
	ErrUserCannotFound           = errors.New("ERR_USER_CANNOT_FOUND")
	ErrUserCannotGet             = errors.New("ERR_USER_CANNOT_GET")
	ErrorEmailNotVerified        = errors.New("ERR_EMAIL_NOT_VERIFIED")
	ErrVerificationTokenNotFound = errors.New("ERR_VERIFICATION_TOKEN_NOT_FOUND")
	ErrTokenExpired              = errors.New("ERR_TOKEN_EXPIRED")
	ErrTokenAlreadyUsed          = errors.New("ERR_TOKEN_ALREADY_USED")
	ErrInvalidToken              = errors.New("ERR_INVALID_TOKEN")
	ErrRegistrationFailed        = errors.New("ERR_REGISTRATION_FAILED")
	ErrRoleRequired              = errors.New("ERR_ROLE_REQUIRED")
	ErrUserCannotDelete          = errors.New("ERR_USER_CANNOT_DELETE")
	ErrPasswordUpdateFailed      = errors.New("ERR_PASSWORD_UPDATE_FAILED")
	ErrPasswordPolicyViolation   = errors.New("ERR_PASSWORD_POLICY_VIOLATION")
	ErrInvalidCredentials        = errors.New("ERR_INVALID_CREDENTIALS")
	ErrForbidden                 = errors.New("ERR_FORBIDDEN")
)

// Infrastructure Errors (MOD_INFRA_*)
var (
	ErrKeycloakInconsistentState  = errors.New("ERR_KC_INCONSISTENT_STATE")
	ErrKeycloakUserCreationFailed = errors.New("ERR_KC_USER_CREATION_FAILED")
	ErrKeycloakCleanupFailed      = errors.New("ERR_KC_CLEANUP_FAILED")
	// Dependency availability errors
	ErrKeycloakUnavailable = errors.New("ERR_KC_UNAVAILABLE")
	ErrDatabaseUnavailable = errors.New("ERR_DB_UNAVAILABLE")
	// Specific user existence errors
	ErrKeycloakUserExists     = errors.New("ERR_KC_USER_EXISTS")
	ErrDatabaseUserExists     = errors.New("ERR_DB_USER_EXISTS")
	ErrIncompleteRegistration = errors.New("ERR_INCOMPLETE_REGISTRATION")
)

// Request Validation Errors (MOD_V_*)
var (
	ErrInvalidJSONFormat = errors.New("ERR_INVALID_JSON_FORMAT")
	ErrInvalidRequest    = errors.New("ERR_INVALID_REQUEST")
	ErrInvalidID         = errors.New("ERR_INVALID_ID")
	ErrInternalServer    = errors.New("ERR_INTERNAL_SERVER")

	// Schema validation errors
	ErrSchemaBadRequest       = errors.New("ERR_SCHEMA_BAD_REQUEST")
	ErrSchemaInvalidRequest   = errors.New("ERR_SCHEMA_INVALID_REQUEST")
	ErrSchemaReadFailed       = errors.New("ERR_SCHEMA_READ_FAILED")
	ErrSchemaEmpty            = errors.New("ERR_SCHEMA_EMPTY")
	ErrSchemaCompileFailed    = errors.New("ERR_SCHEMA_COMPILE_FAILED")
	ErrSchemaValidationFailed = errors.New("ERR_SCHEMA_VALIDATION_FAILED")
	ErrSchemaBodyReadFailed   = errors.New("ERR_SCHEMA_BODY_READ_FAILED")
	ErrSchemaFieldFormat      = errors.New("ERR_SCHEMA_FIELD_FORMAT")
	ErrSchemaFieldRequired    = errors.New("ERR_SCHEMA_FIELD_REQUIRED")
	ErrSchemaFieldType        = errors.New("ERR_SCHEMA_FIELD_TYPE")
	ErrSchemaMultipleFields   = errors.New("ERR_SCHEMA_MULTIPLE_FIELDS")
)

// Authorization Errors (MOD_A_*)
var (
	ErrRoleAssignmentFailed = errors.New("ERR_ROLE_ASSIGNMENT_FAILED")
	ErrRoleRemovalFailed    = errors.New("ERR_ROLE_REMOVAL_FAILED")
	ErrRoleCheckFailed      = errors.New("ERR_ROLE_CHECK_FAILED")
	ErrGetUserRolesFailed   = errors.New("ERR_GET_USER_ROLES_FAILED")
)

// Message Management Errors (MOD_M_*)
var (
	ErrMessageNotFound         = errors.New("ERR_MESSAGE_NOT_FOUND")
	ErrMessageCodeRequired     = errors.New("ERR_MESSAGE_CODE_REQUIRED")
	ErrMessageTypeRequired     = errors.New("ERR_MESSAGE_TYPE_REQUIRED")
	ErrMessageTitleRequired    = errors.New("ERR_MESSAGE_TITLE_REQUIRED")
	ErrMessageContentRequired  = errors.New("ERR_MESSAGE_CONTENT_REQUIRED")
	ErrMessageModuleRequired   = errors.New("ERR_MESSAGE_MODULE_REQUIRED")
	ErrMessageCategoryRequired = errors.New("ERR_MESSAGE_CATEGORY_REQUIRED")
	ErrMessageCodeDuplicate    = errors.New("ERR_MESSAGE_CODE_DUPLICATE")
	ErrMessageCannotSave       = errors.New("ERR_MESSAGE_CANNOT_SAVE")
	ErrMessageCannotUpdate     = errors.New("ERR_MESSAGE_CANNOT_UPDATE")
	ErrMessageCannotDelete     = errors.New("ERR_MESSAGE_CANNOT_DELETE")
	ErrMessageInvalidType      = errors.New("ERR_MESSAGE_INVALID_TYPE")
	ErrMessageListFailed       = errors.New("ERR_MESSAGE_LIST_FAILED")
)

// Branch Management Errors (MOD_B_*) - HU59
var (
	ErrBranchNotFound      = errors.New("ERR_BRANCH_NOT_FOUND")
	ErrBranchCannotSave    = errors.New("ERR_BRANCH_CANNOT_SAVE")
	ErrDuplicateBranchName = errors.New("ERR_DUPLICATE_BRANCH_NAME")
	ErrInvalidBranchType   = errors.New("ERR_INVALID_BRANCH_TYPE")
	ErrBranchCannotUpdate  = errors.New("ERR_BRANCH_CANNOT_UPDATE")
	ErrBranchCannotDelete  = errors.New("ERR_BRANCH_CANNOT_DELETE")
	ErrLocationCannotSave  = errors.New("ERR_LOCATION_CANNOT_SAVE")
	ErrDuplicateAddress    = errors.New("ERR_DUPLICATE_ADDRESS")
	ErrBrandNotFound       = errors.New("ERR_BRAND_NOT_FOUND")
	ErrInvalidImageURL     = errors.New("ERR_INVALID_IMAGE_URL")
)

// Person Deletion Errors (HU53)
var (
	ErrPersonHasBranches = errors.New("ERR_PERSON_HAS_BRANCHES")
)

// Location/Geographic Errors (MOD_L_*)
var (
	ErrDepartmentNotFound  = errors.New("ERR_DEPARTMENT_NOT_FOUND")
	ErrCityNotFound        = errors.New("ERR_CITY_NOT_FOUND")
	ErrCityNotInDepartment = errors.New("ERR_CITY_NOT_IN_DEPARTMENT")
)

// Location Module (MOD_L_*) - Geographic Catalogs
const (
	MsgDepartmentsRetrieved = "MOD_L_DEP_EXI_00001"
	MsgCitiesRetrieved      = "MOD_L_CIT_EXI_00001"
)

// Franchise Errors (MOD_F_*) - HU26-29
var (
	ErrFranchiseNotFound       = errors.New("ERR_FRANCHISE_NOT_FOUND")
	ErrFranchiseDuplicateName  = errors.New("ERR_FRANCHISE_DUPLICATE_NAME")
	ErrFranchiseNoBranches     = errors.New("ERR_FRANCHISE_NO_BRANCHES")
	ErrFranchiseBranchNotOwned = errors.New("ERR_FRANCHISE_BRANCH_NOT_OWNED")
	ErrFranchiseHasBranches    = errors.New("ERR_FRANCHISE_HAS_BRANCHES")
	ErrFranchiseCannotSave     = errors.New("ERR_FRANCHISE_CANNOT_SAVE")
	ErrFranchiseCannotUpdate   = errors.New("ERR_FRANCHISE_CANNOT_UPDATE")
	ErrFranchiseCannotDelete   = errors.New("ERR_FRANCHISE_CANNOT_DELETE")
	ErrFranchiseMinBranches    = errors.New("ERR_FRANCHISE_MIN_BRANCHES")
)

// ============================================
// MESSAGE CODES - Constants for use in code
// ============================================

// User Module (MOD_U_*)
const (
	MsgUserDuplicate        = "MOD_U_DUP_ERR_00001"
	MsgUserCannotSave       = "MOD_U_SAVE_ERR_00002"
	MsgUserNotFound         = "MOD_U_GET_ERR_00003"
	MsgUserEmailError       = "MOD_U_EMAIL_ERR_00004"
	MsgUserEmailNotFound    = "MOD_U_EMAIL_NF_ERR_00005"
	MsgUserEmailNotVerified = "MOD_U_EMAIL_NV_ERR_00006"
	MsgUserTokenNotFound    = "MOD_U_TOKEN_NF_ERR_00007"
	MsgUserTokenExpired     = "MOD_U_TOKEN_EXP_ERR_00008"
	MsgUserTokenUsed        = "MOD_U_TOKEN_USED_ERR_00009"
	MsgUserRegError         = "MOD_U_REG_ERR_00010"
	MsgUserRoleRequired     = "MOD_U_ROLE_REQ_ERR_00011"
	MsgUserCannotDelete     = "MOD_U_DEL_ERR_00012"

	MsgUserRegistered    = "MOD_U_REG_EXI_00001"
	MsgUserUpdated       = "MOD_U_UPD_EXI_00002"
	MsgUserDeleted       = "MOD_U_DEL_EXI_00003"
	MsgUserEmailVerified = "MOD_U_VER_EXI_00004"
)

// Person Module (MOD_P_*)
const (
	MsgPersonNotFound         = "MOD_P_NOT_FOUND_ERR_00001"
	MsgPersonInvalidTx        = "MOD_P_TRANS_ERR_00002"
	MsgPersonRegistered       = "MOD_P_REG_EXI_00001"
	MsgPersonUpdated          = "MOD_P_UPD_EXI_00002"
	MsgPersonCannotDelete     = "MOD_P_DEL_ERR_00001"          // HU53: Generic delete error
	MsgPersonContactRetrieved = "MOD_P_CONTACT_EXI_00001"      // HU55: Public contact retrieved
	MsgPersonDeleted          = "MOD_P_DEL_EXI_00001"          // HU53: Person deleted success
	MsgPersonHasBranches      = "MOD_P_HAS_BRANCHES_ERR_00001" // HU53: Has branches
)

// Validation Module (MOD_V_*)
const (
	MsgValBadFormat     = "MOD_V_VAL_ERR_00001"
	MsgValInvalidReq    = "MOD_V_VAL_ERR_00002"
	MsgValSchemaRead    = "MOD_V_VAL_ERR_00003"
	MsgValSchemaEmpty   = "MOD_V_VAL_ERR_00004"
	MsgValSchemaCompile = "MOD_V_VAL_ERR_00005"
	MsgValFailed        = "MOD_V_VAL_ERR_00006"
	MsgValBodyRead      = "MOD_V_VAL_ERR_00007"
	MsgValFieldFormat   = "MOD_V_VAL_ERR_00008"
	MsgValFieldRequired = "MOD_V_VAL_ERR_00009"
	MsgValFieldType     = "MOD_V_VAL_ERR_00010"
	MsgValMultiple      = "MOD_V_VAL_ERR_00011"
	MsgValJSONInvalid   = "MOD_V_JSON_ERR_00012"
	MsgValIDInvalid     = "MOD_V_ID_ERR_00013"
)

// Authorization Module (MOD_A_*)
const (
	MsgRoleAssignError = "MOD_A_ROLE_ASSIGN_ERR_00001"
	MsgRoleRemoveError = "MOD_A_ROLE_REMOVE_ERR_00002"
	MsgRoleCheckError  = "MOD_A_ROLE_CHECK_ERR_00003"
	MsgRoleGetError    = "MOD_A_ROLE_GET_ERR_00004"
	MsgRoleAssigned    = "MOD_A_ROLE_ASSIGN_EXI_00001"
	MsgRoleRemoved     = "MOD_A_ROLE_REMOVE_EXI_00002"
)

// General Module (GEN_*)
const (
	MsgServerError   = "GEN_SRV_ERR_00001"
	MsgUnauthorized  = "GEN_AUTH_ERR_00002"
	MsgForbidden     = "GEN_FORBIDDEN_ERR_00003"
	MsgOpSuccess     = "GEN_OPE_EXI_00001"
	MsgInfoProcess   = "GEN_INFO_00001"
	MsgWarningAction = "GEN_WARN_00001"
)

// Message Module (MOD_M_*)
const (
	MsgMessageNotFound       = "MOD_M_NOT_FOUND_ERR_00001"
	MsgMessageCodeRequired   = "MOD_M_CODE_REQ_ERR_00002"
	MsgMessageTypeRequired   = "MOD_M_TYPE_REQ_ERR_00003"
	MsgMessageTitleRequired  = "MOD_M_TITLE_REQ_ERR_00004"
	MsgMessageContentReq     = "MOD_M_CONTENT_REQ_ERR_00005"
	MsgMessageModuleRequired = "MOD_M_MODULE_REQ_ERR_00006"
	MsgMessageCategoryReq    = "MOD_M_CATEGORY_REQ_ERR_00007"
	MsgMessageCodeDuplicate  = "MOD_M_CODE_DUP_ERR_00008"
	MsgMessageSaveError      = "MOD_M_SAVE_ERR_00009"
	MsgMessageUpdateError    = "MOD_M_UPDATE_ERR_00010"
	MsgMessageDeleteError    = "MOD_M_DELETE_ERR_00011"
	MsgMessageInvalidType    = "MOD_M_TYPE_INV_ERR_00012"
	MsgMessageListError      = "MOD_M_LIST_ERR_00013"
	MsgMessageCannotDelete   = "MOD_M_DEL_ERR_00014"

	MsgMessageCreated = "MOD_M_CREATE_EXI_00001"
	MsgMessageUpdated = "MOD_M_UPDATE_EXI_00002"
	MsgMessageDeleted = "MOD_M_DELETE_EXI_00003"
	MsgMessageListed  = "MOD_M_LIST_EXI_00004"
)

// Infrastructure Module (MOD_INFRA_*)
const (
	MsgKeycloakInconsistentState = "MOD_INFRA_KC_INCONSISTENT_ERR_00001"
	MsgKeycloakCreateError       = "MOD_INFRA_KC_CREATE_ERR_00002"
	MsgKeycloakCleanupError      = "MOD_INFRA_KC_CLEANUP_ERR_00003"
	// Dependency availability messages
	MsgKeycloakUnavailable = "MOD_INFRA_KC_UNAVAIL_ERR_00004"
	MsgDatabaseUnavailable = "MOD_INFRA_DB_UNAVAIL_ERR_00005"
	MsgDependencyFailure   = "MOD_INFRA_DEP_FAIL_ERR_00006"
	// Specific user existence messages
	MsgKeycloakUserExists     = "MOD_INFRA_KC_USER_EXISTS_ERR_00007"
	MsgDatabaseUserExists     = "MOD_INFRA_DB_USER_EXISTS_ERR_00008"
	MsgIncompleteRegistration = "MOD_INFRA_INCOMPLETE_REG_ERR_00009"
)

// Keycloak Module (MOD_KC_*) - Email Verification and Password Reset
const (
	// Email Verification
	MsgKCEmailVerified        = "MOD_KC_EMAIL_VERIFIED_EXI_00001"
	MsgKCInvalidToken         = "MOD_KC_INVALID_TOKEN_ERR_00001"
	MsgKCEmailVerifyError     = "MOD_KC_EMAIL_VERIFY_ERROR_ERR_00001"
	MsgKCUserNotFound         = "MOD_KC_USER_NOT_FOUND_ERR_00001"
	MsgKCEmailAlreadyVerified = "MOD_KC_EMAIL_ALREADY_VERIFIED_WARN_00001"
	// Verification Email Sending
	MsgKCVerifEmailSent  = "MOD_KC_VERIF_EMAIL_SENT_EXI_00001"
	MsgKCVerifEmailError = "MOD_KC_VERIF_EMAIL_ERROR_ERR_00001"
	// Password Reset
	MsgKCPwdResetSent  = "MOD_KC_PWD_RESET_SENT_EXI_00001"
	MsgKCPwdResetError = "MOD_KC_PWD_RESET_ERROR_ERR_00001"
	// Auth Profile
	MsgAuthProfileRetrieved = "MOD_AUTH_PROFILE_EXI_00001"
	// Change Password (HU57)
	MsgChangePasswordSuccess        = "MOD_P_CHANGE_EXI_00001"
	MsgChangePasswordInvalidCurrent = "MOD_P_CHANGE_ERR_00001"
	MsgChangePasswordUpdateError    = "MOD_P_CHANGE_ERR_00002"
	MsgChangePasswordPolicyError    = "MOD_P_CHANGE_ERR_00003"
)

// Branch Module (MOD_B_*) - HU59, HU62, HU76
const (
	MsgBranchRegistered    = "MOD_B_REG_EXI_00001"
	MsgBranchCannotSave    = "MOD_B_REG_ERR_00001"
	MsgBranchDuplicateName = "MOD_B_DUP_NAME_ERR_00001"
	MsgBranchInvalidType   = "MOD_B_INVALID_TYPE_ERR_00001"
	MsgBranchNotFound      = "MOD_B_NOT_FOUND_ERR_00001"
	MsgBranchFound         = "MOD_B_GET_EXI_00001"   // HU62: Branch found
	MsgBranchTypesFound    = "MOD_B_TYPES_EXI_00001" // HU76: Branch types catalog
	MsgBranchUpdated       = "MOD_B_UPD_EXI_00001"
	MsgBranchCannotUpdate  = "MOD_B_UPD_ERR_00001" // HU60: Update error
	MsgBranchDeleted       = "MOD_B_DEL_EXI_00001"
	MsgBrandNotFound       = "MOD_B_BRAND_NOT_FOUND_ERR_00001"
	MsgDuplicateAddress    = "MOD_B_DUP_ADDR_ERR_00001"
	MsgBranchListFound     = "MOD_B_LIST_EXI_00001"      // HU62: Branch list found
	MsgBranchCannotDelete  = "MOD_B_DEL_ERR_00001"       // HU61: Delete error
	MsgBranchHasAssoc      = "MOD_B_HAS_ASSOC_ERR_00001" // HU61: Has associations
	MsgBrandsRetrieved     = "MOD_B_BRD_EXI_00001"       // Brand catalog retrieved
	MsgBranchNearbyFound   = "MOD_B_NEARBY_EXI_00001"    // HU89: Nearby branches found
)

// Validation Module - Geolocation (MOD_V_GEO_*)
const (
	MsgValLatitudeRequired  = "MOD_V_GEO_LAT_REQ_ERR_00001"
	MsgValLatitudeInvalid   = "MOD_V_GEO_LAT_INV_ERR_00001"
	MsgValLongitudeRequired = "MOD_V_GEO_LNG_REQ_ERR_00001"
	MsgValLongitudeInvalid  = "MOD_V_GEO_LNG_INV_ERR_00001"
	MsgValRadiusInvalid     = "MOD_V_GEO_RAD_INV_ERR_00001"
)

// Franchise Module (MOD_F_*) - HU26-29
const (
	// Success messages
	MsgFranchiseCreated       = "MOD_F_REG_EXI_00001"
	MsgFranchiseFound         = "MOD_F_GET_EXI_00001"
	MsgFranchisesListed       = "MOD_F_LIST_EXI_00001"
	MsgFranchiseUpdated       = "MOD_F_UPD_EXI_00001"
	MsgFranchiseDeleted       = "MOD_F_DEL_EXI_00001"
	MsgFranchiseBranchAdded   = "MOD_F_BRANCH_ADD_EXI_00001"
	MsgFranchiseBranchRemoved = "MOD_F_BRANCH_REM_EXI_00001"

	// Error messages
	MsgFranchiseNotFound       = "MOD_F_NOT_FOUND_ERR_00001"
	MsgFranchiseDuplicateName  = "MOD_F_DUP_NAME_ERR_00001"
	MsgFranchiseNoBranches     = "MOD_F_NO_BRANCHES_ERR_00001"
	MsgFranchiseBranchNotOwned = "MOD_F_BRANCH_NOT_OWNED_ERR_00001"
	MsgFranchiseHasBranches    = "MOD_F_HAS_BRANCHES_ERR_00001"
	MsgFranchiseMinBranches    = "MOD_F_MIN_BRANCHES_ERR_00001"
)

// Service Catalog Module (MOD_S_*) - HU63, HU68, HU75
const (
	// Success messages
	MsgServiceTypesRetrieved = "MOD_S_TYPES_EXI_00001"       // HU75: Service types catalog
	MsgServicesRetrieved     = "MOD_S_LIST_EXI_00001"        // HU63: Services list
	MsgServiceAssociated     = "MOD_S_ASSOC_EXI_00001"       // Servicios asociados a sede
	MsgServiceDissociated    = "MOD_S_DISSOC_EXI_00001"      // Servicio desasociado de sede
	MsgServiceUpdated        = "MOD_S_UPD_EXI_00001"         // HU68: Service updated (Admin)
	MsgServiceActivated      = "MOD_S_ACTIVATED_EXI_00001"   // HU68: Service activated (Admin)
	MsgServiceDeactivated    = "MOD_S_DEACTIVATED_EXI_00001" // HU68: Service deactivated (Admin)

	// Error messages
	MsgServiceInvalidType       = "MOD_S_INVALID_TYPE_ERR_00001"  // Tipo de servicio inválido
	MsgServiceNotFound          = "MOD_S_NOT_FOUND_ERR_00001"     // Servicio no encontrado (legacy)
	MsgServiceAlreadyAssociated = "MOD_S_ALREADY_ASSOC_ERR_00001" // Servicio ya asociado
	MsgServiceResNotFound       = "MOD_S_RES_ERR_00001"           // HU68: Service not found in catalog
	MsgServiceTypeInvalid       = "MOD_S_TYPE_ERR_00001"          // HU68: Invalid service type
)

// Service Errors (MOD_S_*)
var (
	ErrServiceNotFound          = errors.New("ERR_SERVICE_NOT_FOUND")
	ErrServiceAlreadyAssociated = errors.New("ERR_SERVICE_ALREADY_ASSOCIATED")
)

// Schedule Errors (MOD_H_*) - HU30-35
var (
	ErrScheduleNotFound      = errors.New("ERR_SCHEDULE_NOT_FOUND")
	ErrScheduleAlreadyExists = errors.New("ERR_SCHEDULE_ALREADY_EXISTS")
	ErrInvalidDayOfWeek      = errors.New("ERR_INVALID_DAY_OF_WEEK")
	ErrInvalidTimeFormat     = errors.New("ERR_INVALID_TIME_FORMAT")
	ErrClosingBeforeOpening  = errors.New("ERR_CLOSING_BEFORE_OPENING")
	ErrScheduleInactive      = errors.New("ERR_SCHEDULE_INACTIVE")
)

// Schedule Module (MOD_H_*) - HU30-35, HU10
const (
	// Success messages
	MsgScheduleCreated      = "MOD_H_CREATE_EXI_00001"
	MsgScheduleRetrieved    = "MOD_H_GET_EXI_00001"
	MsgScheduleUpdated      = "MOD_H_UPDATE_EXI_00001"
	MsgScheduleDeleted      = "MOD_H_DELETE_EXI_00001"
	MsgScheduleActivated    = "MOD_H_ACTIV_EXI_00001"
	MsgScheduleDeactivated  = "MOD_H_DEACT_EXI_00001"
	MsgDaysCatalogRetrieved = "MOD_H_DAYS_EXI_00001"

	// Error messages
	MsgScheduleNotFound          = "MOD_H_NOT_FOUND_ERR_00001"
	MsgScheduleAlreadyExists     = "MOD_H_EXISTS_ERR_00001"
	MsgInvalidDayOfWeek          = "MOD_H_DAY_ERR_00001"
	MsgInvalidTimeFormat         = "MOD_H_TIME_ERR_00001"
	MsgClosingBeforeOpening      = "MOD_H_TIME_ORDER_ERR_00001"
	MsgScheduleInactive          = "MOD_H_INACTIVE_ERR_00001"
	MsgScheduleInvalidDateFormat = "MOD_H_DATE_FORMAT_ERR_00001"
	MsgScheduleInvalidDateRange  = "MOD_H_DATE_RANGE_ERR_00001"
)

// Schedule Detail Errors (MOD_HD_*) - HU6-9
var (
	ErrScheduleDetailNotFound         = errors.New("ERR_SCHEDULE_DETAIL_NOT_FOUND")
	ErrScheduleDetailTimeConflict     = errors.New("ERR_SCHEDULE_DETAIL_TIME_CONFLICT")
	ErrScheduleDetailInvalidTime      = errors.New("ERR_SCHEDULE_DETAIL_INVALID_TIME")
	ErrScheduleDetailInvalidDay       = errors.New("ERR_SCHEDULE_DETAIL_INVALID_DAY")
	ErrScheduleDetailDayAlreadyClosed = errors.New("ERR_SCHEDULE_DETAIL_DAY_CLOSED")
	ErrScheduleDetailDayHasSlots      = errors.New("ERR_SCHEDULE_DETAIL_DAY_HAS_SLOTS")
)

// Schedule Detail Module (MOD_HD_*) - HU6-9
const (
	// Success messages
	MsgScheduleDetailCreated   = "MOD_HD_CREATE_EXI_00001"
	MsgScheduleDetailRetrieved = "MOD_HD_GET_EXI_00001"
	MsgScheduleDetailUpdated   = "MOD_HD_UPDATE_EXI_00001"
	MsgScheduleDetailDeleted   = "MOD_HD_DELETE_EXI_00001"
	MsgScheduleDetailsListed   = "MOD_HD_LIST_EXI_00001"

	// Error messages
	MsgScheduleDetailNotFound         = "MOD_HD_NOT_FOUND_ERR_00001"
	MsgScheduleDetailTimeConflict     = "MOD_HD_CONFLICT_ERR_00001"
	MsgScheduleDetailInvalidTime      = "MOD_HD_TIME_ERR_00001"
	MsgScheduleDetailInvalidDay       = "MOD_HD_DAY_ERR_00001"
	MsgScheduleDetailDayAlreadyClosed = "MOD_HD_DAY_CLOSED_ERR_00001"
	MsgScheduleDetailDayHasSlots      = "MOD_HD_DAY_HAS_SLOTS_ERR_00001"
)

// Schedule Exception Errors (MOD_EH_*) - HU20-25
var (
	ErrScheduleExceptionNotFound     = errors.New("ERR_SCHEDULE_EXCEPTION_NOT_FOUND")
	ErrScheduleExceptionDateConflict = errors.New("ERR_SCHEDULE_EXCEPTION_DATE_CONFLICT")
	ErrScheduleExceptionDatePast     = errors.New("ERR_SCHEDULE_EXCEPTION_DATE_PAST")
	ErrScheduleExceptionInvalidTime  = errors.New("ERR_SCHEDULE_EXCEPTION_INVALID_TIME")
	ErrScheduleExceptionRedundant    = errors.New("ERR_SCHEDULE_EXCEPTION_REDUNDANT")
)

// Schedule Exception Module (MOD_EH_*) - HU20-25
const (
	// Success messages
	MsgScheduleExceptionCreated     = "MOD_EH_CREATE_EXI_00001"
	MsgScheduleExceptionRetrieved   = "MOD_EH_GET_EXI_00001"
	MsgScheduleExceptionsListed     = "MOD_EH_LIST_EXI_00001"
	MsgScheduleExceptionUpdated     = "MOD_EH_UPDATE_EXI_00001"
	MsgScheduleExceptionDeleted     = "MOD_EH_DELETE_EXI_00001"
	MsgScheduleExceptionActivated   = "MOD_EH_ACTIVATE_EXI_00001"
	MsgScheduleExceptionDeactivated = "MOD_EH_DEACTIVATE_EXI_00001"

	// Error messages
	MsgScheduleExceptionNotFound     = "MOD_EH_NOT_FOUND_ERR_00001"
	MsgScheduleExceptionDateConflict = "MOD_EH_DATE_CONFLICT_ERR_00001"
	MsgScheduleExceptionDatePast     = "MOD_EH_DATE_PAST_ERR_00001"
	MsgScheduleExceptionInvalidTime  = "MOD_EH_TIME_ERR_00001"
	MsgScheduleExceptionRedundant    = "MOD_EH_REDUNDANT_ERR_00001"
)

// Motorcycle Errors (MOD_MOT_*) - HU43-47
var (
	ErrMotorcycleNotFound     = errors.New("ERR_MOTORCYCLE_NOT_FOUND")
	ErrMotorcycleCannotSave   = errors.New("ERR_MOTORCYCLE_CANNOT_SAVE")
	ErrMotorcycleCannotUpdate = errors.New("ERR_MOTORCYCLE_CANNOT_UPDATE")
	ErrMotorcycleCannotDelete = errors.New("ERR_MOTORCYCLE_CANNOT_DELETE")
	ErrDuplicateLicensePlate  = errors.New("ERR_DUPLICATE_LICENSE_PLATE")
	ErrReferenceNotFound      = errors.New("ERR_REFERENCE_NOT_FOUND")
	ErrReferenceRequired      = errors.New("ERR_REFERENCE_REQUIRED")
)

// Motorcycle Module (MOD_MOT_*) - HU43-47
const (
	// Success messages
	MsgMotorcycleCreated          = "MOD_MOT_CREATE_EXI_00001"
	MsgMotorcycleRetrieved        = "MOD_MOT_GET_EXI_00001"
	MsgMotorcycleUpdated          = "MOD_MOT_UPDATE_EXI_00001"
	MsgMotorcycleDeleted          = "MOD_MOT_DELETE_EXI_00001"
	MsgMotorcyclesListed          = "MOD_MOT_LIST_EXI_00001"
	MsgMotorcycleReferencesListed = "MOD_MOT_REF_LIST_EXI_00001"    // HU50
	MsgBrandLinesRetrieved        = "MOD_MOT_BRAND_LINES_EXI_00001" // HU40

	// Profile Image Success Messages (HU36-39)
	MsgProfileImageUpdated = "MOD_MOT_IMG_UPDATE_EXI_00001" // HU36/37
	MsgProfileImageGet     = "MOD_MOT_IMG_GET_EXI_00001"    // HU38
	MsgProfileImageDeleted = "MOD_MOT_IMG_DELETE_EXI_00001" // HU39

	// Error messages
	MsgMotorcycleNotFound          = "MOD_MOT_NOT_FOUND_ERR_00001"
	MsgMotorcycleCannotSave        = "MOD_MOT_CREATE_ERR_00001"
	MsgMotorcycleCannotUpdate      = "MOD_MOT_UPDATE_ERR_00001"
	MsgMotorcycleCannotDelete      = "MOD_MOT_DELETE_ERR_00001"
	MsgDuplicateLicensePlate       = "MOD_MOT_DUP_PLATE_ERR_00001"
	MsgMotorcycleReferenceNotFound = "MOD_MOT_REF_NOT_FOUND_ERR_00001"
	MsgReferenceRequired           = "MOD_MOT_REF_REQ_ERR_00001"
	MsgMissingPlateParam           = "MOD_MOT_PLATE_REQ_ERR_00001"

	// Profile Image Error Messages (HU36-39)
	MsgProfileImageUpdateError = "MOD_MOT_IMG_UPDATE_ERR_00001"
	MsgProfileImageNotFound    = "MOD_MOT_IMG_NOT_FOUND_ERR_00001"
	MsgProfileImageURLRequired = "MOD_MOT_IMG_URL_REQ_ERR_00001"
)

// Motorcycle Evidence Errors (MOD_EVD_*) - HU16-19
var (
	ErrEvidenceNotFound      = errors.New("ERR_EVIDENCE_NOT_FOUND")
	ErrEvidenceCannotSave    = errors.New("ERR_EVIDENCE_CANNOT_SAVE")
	ErrEvidenceCannotUpdate  = errors.New("ERR_EVIDENCE_CANNOT_UPDATE")
	ErrEvidenceCannotDelete  = errors.New("ERR_EVIDENCE_CANNOT_DELETE")
	ErrEvidenceLimitExceeded = errors.New("ERR_EVIDENCE_LIMIT_EXCEEDED")
	ErrInvalidEvidenceURL    = errors.New("ERR_INVALID_EVIDENCE_URL")
	ErrInvalidEvidenceAngle  = errors.New("ERR_INVALID_EVIDENCE_ANGLE")
)

// Motorcycle Evidence Module (MOD_EVD_*) - HU16-19
const (
	// Success messages
	MsgEvidenceCreated   = "MOD_EVD_CREATE_EXI_00001"
	MsgEvidenceRetrieved = "MOD_EVD_GET_EXI_00001"
	MsgEvidenceUpdated   = "MOD_EVD_UPDATE_EXI_00001"
	MsgEvidenceDeleted   = "MOD_EVD_DELETE_EXI_00001"
	MsgEvidencesListed   = "MOD_EVD_LIST_EXI_00001"

	// Error messages
	MsgEvidenceNotFound      = "MOD_EVD_NOT_FOUND_ERR_00001"
	MsgEvidenceCannotSave    = "MOD_EVD_CREATE_ERR_00001"
	MsgEvidenceCannotUpdate  = "MOD_EVD_UPDATE_ERR_00001"
	MsgEvidenceCannotDelete  = "MOD_EVD_DELETE_ERR_00001"
	MsgEvidenceLimitExceeded = "MOD_EVD_LIMIT_ERR_00001"
	MsgInvalidEvidenceURL    = "MOD_EVD_URL_ERR_00001"
	MsgInvalidEvidenceAngle  = "MOD_EVD_ANGLE_ERR_00001"
)

// Diagnostic Errors (MOD_DGN_*) - HU11-14
var (
	ErrDiagnosticNotFound     = errors.New("ERR_DIAGNOSTIC_NOT_FOUND")
	ErrDiagnosticCannotSave   = errors.New("ERR_DIAGNOSTIC_CANNOT_SAVE")
	ErrDiagnosticCannotUpdate = errors.New("ERR_DIAGNOSTIC_CANNOT_UPDATE")
	ErrDiagnosticCannotDelete = errors.New("ERR_DIAGNOSTIC_CANNOT_DELETE")
)

// Diagnostic Module (MOD_DGN_*) - HU11-14
const (
	// Success messages
	MsgDiagnosticCreated   = "MOD_DGN_CREATE_EXI_00001"
	MsgDiagnosticRetrieved = "MOD_DGN_GET_EXI_00001"
	MsgDiagnosticUpdated   = "MOD_DGN_UPDATE_EXI_00001"
	MsgDiagnosticDeleted   = "MOD_DGN_DELETE_EXI_00001"
	MsgDiagnosticsListed   = "MOD_DGN_LIST_EXI_00001"

	// Error messages
	MsgDiagnosticNotFound     = "MOD_DGN_NOT_FOUND_ERR_00001"
	MsgDiagnosticCannotSave   = "MOD_DGN_CREATE_ERR_00001"
	MsgDiagnosticCannotUpdate = "MOD_DGN_UPDATE_ERR_00001"
	MsgDiagnosticCannotDelete = "MOD_DGN_DELETE_ERR_00001"
)

// Diagnostic Permission Module (MOD_DGP_*)
const (
	// Success messages
	MsgPermissionGranted = "MOD_DGP_GRANT_EXI_00001"
	MsgPermissionRevoked = "MOD_DGP_REVOKE_EXI_00001"
	MsgPermissionsListed = "MOD_DGP_LIST_EXI_00001"

	// Error messages
	MsgPermissionNotFound     = "MOD_DGP_NOT_FOUND_ERR_00001"
	MsgPermissionCannotSave   = "MOD_DGP_SAVE_ERR_00001"
	MsgPermissionCannotDelete = "MOD_DGP_DELETE_ERR_00001"
)
