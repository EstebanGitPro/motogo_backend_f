package domain

import "errors"

// User Management Errors (MOD_U_*)
var (
	ErrDuplicateUser      = errors.New("usuario ya existe")
	ErrUserCannotSave     = errors.New("no se puede guardar el usuario")
	ErrPersonNotFound     = errors.New("persona no encontrada")
	ErrInvalidTransaction = errors.New("transacción inválida")

	ErrGettingUserByEmail        = errors.New("error obteniendo usuario por email")
	ErrNotFoundUserByEmail       = errors.New("usuario no encontrado por email")
	ErrUserCannotFound           = errors.New("usuario no puede ser encontrado")
	ErrUserCannotGet             = errors.New("usuario no puede ser obtenido")
	ErrorEmailNotVerified        = errors.New("email no verificado")
	ErrVerificationTokenNotFound = errors.New("token de verificación no encontrado")
	ErrTokenExpired              = errors.New("token expirado")
	ErrTokenAlreadyUsed          = errors.New("token ya utilizado")
	ErrRegistrationFailed        = errors.New("error en el proceso de registro")
	ErrRoleRequired              = errors.New("el rol es requerido")
)

// Request Validation Errors (MOD_V_*)
var (
	ErrInvalidJSONFormat = errors.New("formato JSON inválido")
	ErrInvalidRequest    = errors.New("parámetros de solicitud inválidos")
	ErrInvalidID         = errors.New("ID no válido")
	ErrInternalServer    = errors.New("error interno del servidor")
)

// Authorization Errors (MOD_A_*)
var (
	ErrRoleAssignmentFailed = errors.New("error al asignar rol")
	ErrRoleRemovalFailed    = errors.New("error al remover rol")
	ErrRoleCheckFailed      = errors.New("error al verificar rol")
	ErrGetUserRolesFailed   = errors.New("error al obtener roles del usuario")
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

	MsgUserRegistered    = "MOD_U_REG_EXI_00001"
	MsgUserUpdated       = "MOD_U_UPD_EXI_00002"
	MsgUserDeleted       = "MOD_U_DEL_EXI_00003"
	MsgUserEmailVerified = "MOD_U_VER_EXI_00004"
)

// Person Module (MOD_P_*)
const (
	MsgPersonNotFound   = "MOD_P_NOT_FOUND_ERR_00001"
	MsgPersonInvalidTx  = "MOD_P_TRANS_ERR_00002"
	MsgPersonRegistered = "MOD_P_REG_EXI_00001"
	MsgPersonUpdated    = "MOD_P_UPD_EXI_00002"
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
