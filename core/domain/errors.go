package domain

import "errors"

type DomainError struct {
	Code    string
	Message string
	err     error
}

func (e *DomainError) Error() string {
	return e.Message
}

func (e *DomainError) Unwrap() error {
	return e.err
}


func NewDomainError(code, message string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		err:     errors.New(message),
	}
}

// User Management Errors (MOD_U_*)
var (
	ErrDuplicateUser  = NewDomainError("MOD_U_USU_ERR_00001", "Usuario ya existe")
	ErrUserCannotSave = NewDomainError("MOD_U_USU_ERR_00002", "No se puede guardar el usuario")
	ErrPersonNotFound = NewDomainError("MOD_U_USU_ERR_00003", "Persona no encontrada")

	ErrGettingUserByEmail        = NewDomainError("MOD_U_USU_ERR_00004", "Error obteniendo usuario por email")
	ErrNotFoundUserByEmail       = NewDomainError("MOD_U_USU_ERR_00005", "Usuario no encontrado por email")
	ErrUserCannotFound           = NewDomainError("MOD_U_USU_ERR_00006", "Usuario no puede ser encontrado")
	ErrUserCannotGet             = NewDomainError("MOD_U_USU_ERR_00007", "Usuario no puede ser obtenido")
	ErrorEmailNotVerified        = NewDomainError("MOD_U_USU_ERR_00008", "Email no verificado")
	ErrVerificationTokenNotFound = NewDomainError("MOD_U_USU_ERR_00009", "Token de verificación no encontrado")
	ErrTokenExpired              = NewDomainError("MOD_U_USU_ERR_00010", "Token expirado")
	ErrTokenAlreadyUsed          = NewDomainError("MOD_U_USU_ERR_00011", "Token ya utilizado")
	ErrRegistrationFailed        = NewDomainError("MOD_U_USU_ERR_00012", "Error en el proceso de registro")
	ErrRoleRequired              = NewDomainError("MOD_U_USU_ERR_00013", "El rol es requerido")
)

// Request Validation Errors (MOD_V_*)
var (
	ErrInvalidJSONFormat = NewDomainError("MOD_V_VAL_ERR_00001", "Formato JSON inválido")
	ErrInvalidRequest    = NewDomainError("MOD_V_VAL_ERR_00002", "Parámetros de solicitud inválidos")
)

// Authorization Errors (MOD_A_*)
var (
	ErrRoleAssignmentFailed = NewDomainError("MOD_A_AUT_ERR_00001", "Error al asignar rol")
	ErrRoleRemovalFailed    = NewDomainError("MOD_A_AUT_ERR_00002", "Error al remover rol")
	ErrRoleCheckFailed      = NewDomainError("MOD_A_AUT_ERR_00003", "Error al verificar rol")
	ErrGetUserRolesFailed   = NewDomainError("MOD_A_AUT_ERR_00004", "Error al obtener roles del usuario")
)