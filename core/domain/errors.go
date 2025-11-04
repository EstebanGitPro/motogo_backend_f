package domain

import "errors"


// User Management Errors (MOD_U_*)
var (
	ErrDuplicateUser  = errors.New("usuario ya existe")
	ErrUserCannotSave = errors.New("no se puede guardar el usuario")
	ErrPersonNotFound = errors.New("persona no encontrada")

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
)

// Authorization Errors (MOD_A_*)
var (
	ErrRoleAssignmentFailed = errors.New("error al asignar rol")
	ErrRoleRemovalFailed    = errors.New("error al remover rol")
	ErrRoleCheckFailed      = errors.New("error al verificar rol")
	ErrGetUserRolesFailed   = errors.New("error al obtener roles del usuario")
)