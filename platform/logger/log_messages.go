package logger

// ============================================
// APPLICATION LIFECYCLE
// ============================================
const (
	LogAppStarting          = "Iniciando aplicación MotoGo Backend"
	LogAppConfigLoaded      = "Configuración cargada exitosamente"
	LogAppConfigError       = "Error cargando configuración"
	LogAppDatabaseConnected = "Conexión a base de datos establecida"
	LogAppDatabaseError     = "Error conectando a base de datos"
	LogAppServerStarting    = "Servidor iniciando"
	LogAppServerListening   = "Servidor escuchando"
	LogAppShuttingDown      = "Apagando aplicación gracefully"
)

// ============================================
// REGISTRATION / AUTHENTICATION
// ============================================
const (
	LogRegRequestReceived   = "Nueva solicitud de registro recibida"
	LogRegProcessing        = "Procesando registro de usuario"
	LogRegJSONParseError    = "Error parseando JSON del request"
	LogRegProcessError      = "Error en proceso de registro"
	LogRegIDEncodeError     = "Error ofuscando ID de usuario"
	LogRegSuccess           = "Usuario registrado exitosamente"
	LogRegKeycloakSync      = "Sincronizando usuario con Keycloak"
	LogRegKeycloakSyncError = "Error sincronizando con Keycloak"
)

// ============================================
// DATABASE OPERATIONS
// ============================================
const (
	LogDBQueryExecuting      = "Ejecutando query de base de datos"
	LogDBQuerySuccess        = "Query ejecutado exitosamente"
	LogDBQueryError          = "Error ejecutando query"
	LogDBTransactionStart    = "Iniciando transacción de base de datos"
	LogDBTransactionCommit   = "Commit de transacción exitoso"
	LogDBTransactionRollback = "Rollback de transacción"
	LogDBConnectionPoolInfo  = "Información de connection pool"
)

// ============================================
// MESSAGING / CACHE
// ============================================
const (
	LogMsgCacheInit            = "Inicializando cache de mensajes"
	LogMsgCacheLoaded          = "Mensajes del sistema cargados en cache"
	LogMsgCacheLoadError       = "Error cargando mensajes del sistema desde BD"
	LogMsgCacheRefreshStart    = "Iniciando auto-refresh de cache de mensajes"
	LogMsgCacheRefreshDisabled = "Auto-refresh de cache deshabilitado"
	LogMsgCacheRefreshing      = "Auto-refrescando cache de mensajes desde BD"
	LogMsgCacheRefreshOK       = "Cache de mensajes refrescado exitosamente"
	LogMsgCacheRefreshError    = "Error durante auto-refresh de cache"
	LogMsgCacheRefreshStop     = "Deteniendo auto-refresh de cache de mensajes"
	LogMsgNotInCache           = "Mensaje no encontrado en cache, cargando desde BD"
	LogMsgNotInDB              = "Mensaje no encontrado en base de datos"
	LogMsgCachedFromDB         = "Mensaje cargado desde BD y cacheado"
	LogMsgInactive             = "Mensaje encontrado pero está desactivado"
)

// ============================================
// ROUTING / MIDDLEWARE
// ============================================
const (
	LogRouteConfiguring      = "Configurando rutas de la aplicación"
	LogRouteConfigured       = "Rutas configuradas correctamente"
	LogRouteValidatorInit    = "Inicializando validador de schemas"
	LogRouteValidatorOK      = "Validador de schema inicializado"
	LogRouteValidatorError   = "Error creando validador de schema"
	LogMiddlewareErrorCaught = "Error de negocio capturado"
	LogMiddlewareInternalErr = "Error interno del servidor"
)

// ============================================
// VALIDATION
// ============================================
const (
	LogValidationStart   = "Iniciando validación de request"
	LogValidationOK      = "Validación de request exitosa"
	LogValidationFailed  = "Validación de request fallida"
	LogValidationDetails = "Detalles de validación"
)

// ============================================
// KEYCLOAK / EXTERNAL SERVICES
// ============================================
const (
	LogKeycloakClientInit          = "Inicializando cliente Keycloak"
	LogKeycloakClientOK            = "Cliente Keycloak inicializado correctamente"
	LogKeycloakClientError         = "Error inicializando cliente Keycloak"
	LogKeycloakAdminAuth           = "Autenticando admin de Keycloak"
	LogKeycloakAdminAuthError      = "Error autenticando admin de Keycloak"
	LogKeycloakTokenRefresh        = "Refrescando token de admin de Keycloak"
	LogKeycloakTokenRefreshOK      = "Token de admin refrescado exitosamente"
	LogKeycloakTokenRefreshErr     = "Error refrescando token de admin de Keycloak"
	LogKeycloakUserLogin           = "Intentando login de usuario"
	LogKeycloakUserLoginOK         = "Login de usuario exitoso"
	LogKeycloakUserLoginError      = "Error en login de usuario"
	LogKeycloakUserCreate          = "Creando usuario en Keycloak"
	LogKeycloakUserCreateOK        = "Usuario creado en Keycloak"
	LogKeycloakUserCreateError     = "Error creando usuario en Keycloak"
	LogKeycloakUserGet             = "Obteniendo usuario de Keycloak"
	LogKeycloakUserGetError        = "Error obteniendo usuario de Keycloak"
	LogKeycloakUserDelete          = "Eliminando usuario de Keycloak"
	LogKeycloakUserDeleteOK        = "Usuario eliminado de Keycloak"
	LogKeycloakUserDeleteError     = "Error eliminando usuario de Keycloak"
	LogKeycloakPasswordSet         = "Configurando password para usuario"
	LogKeycloakPasswordSetOK       = "Password configurado exitosamente"
	LogKeycloakPasswordSetError    = "Error configurando password"
	LogKeycloakRoleGet             = "Obteniendo rol"
	LogKeycloakRoleGetError        = "Error obteniendo rol"
	LogKeycloakRoleAssign          = "Asignando rol a usuario"
	LogKeycloakRoleAssignOK        = "Rol asignado exitosamente"
	LogKeycloakRoleAssignError     = "Error asignando rol a usuario"
	LogKeycloakUserTokenRefresh    = "Refrescando token de usuario"
	LogKeycloakUserTokenRefreshOK  = "Token de usuario refrescado exitosamente"
	LogKeycloakUserTokenRefreshErr = "Error refrescando token de usuario"
	// Email Verification
	LogKeycloakSendVerificationEmail      = "Enviando email de verificación"
	LogKeycloakSendVerificationEmailOK    = "Email de verificación enviado exitosamente"
	LogKeycloakSendVerificationEmailError = "Error enviando email de verificación"
	// Password Reset
	LogKeycloakSendPasswordReset      = "Enviando email de recuperación de contraseña"
	LogKeycloakSendPasswordResetOK    = "Email de recuperación enviado exitosamente"
	LogKeycloakSendPasswordResetError = "Error enviando email de recuperación"
	// User Search
	LogKeycloakSearchUserByEmail   = "Buscando usuario en Keycloak por email"
	LogKeycloakSearchUserByEmailOK = "Usuario encontrado en Keycloak"
	LogKeycloakUserNotFound        = "Usuario no encontrado en Keycloak"
	// Email Verification (via proxy endpoint)
	LogKeycloakEmailVerify          = "Verificando email de usuario"
	LogKeycloakEmailVerifyOK        = "Email verificado exitosamente"
	LogKeycloakEmailVerifyError     = "Error verificando email"
	LogKeycloakEmailAlreadyVerified = "Email ya estaba verificado"
	LogKeycloakEmailNotVerified     = "Email no verificado - login bloqueado"
	// Resend Verification Email (during login)
	LogKeycloakResendingVerificationEmail   = "Reenviando email de verificación automáticamente"
	LogKeycloakResendVerificationEmailOK    = "Email de verificación reenviado exitosamente"
	LogKeycloakResendVerificationEmailError = "Error reenviando email de verificación"
	// Password Reset Flow
	LogPasswordResetStart          = "Iniciando proceso de recuperación de contraseña"
	LogPasswordResetTokenError     = "Error extrayendo email del token de reset"
	LogPasswordResetEmailExtracted = "Email extraído exitosamente del token"
	LogPasswordResetUserNotFound   = "Usuario no encontrado para reset de contraseña"
	LogPasswordResetUserFound      = "Usuario encontrado para reset de contraseña"
	LogPasswordResetUpdateError    = "Error actualizando contraseña en Keycloak"
	LogPasswordResetSuccess        = "Contraseña actualizada exitosamente"
	// Change Password Flow (HU57)
	LogChangePasswordStart          = "Iniciando proceso de cambio de contraseña"
	LogChangePasswordUserNotFound   = "Usuario no encontrado para cambio de contraseña"
	LogChangePasswordInvalidCurrent = "Contraseña actual incorrecta"
	LogChangePasswordUpdateError    = "Error actualizando contraseña"
	LogChangePasswordSuccess        = "Contraseña cambiada exitosamente"
	// Update Profile Flow (HU52)
	LogUpdateProfileStart            = "Iniciando proceso de actualización de perfil"
	LogUpdateProfileValidation       = "Validando datos del perfil"
	LogUpdateProfileDBSuccess        = "Perfil actualizado en base de datos"
	LogUpdateProfileKeycloakSyncWarn = "Advertencia sincronizando perfil con Keycloak"
	LogUpdateProfileKeycloakSyncOK   = "Perfil sincronizado con Keycloak"
	LogUpdateProfileSuccess          = "Perfil actualizado exitosamente"
	LogUpdateProfileError            = "Error actualizando perfil"
)

// ============================================
// KEYCLOAK AVAILABILITY
// ============================================
const (
	LogKeycloakAvailabilityCheck = "Verificando disponibilidad de Keycloak"
	LogKeycloakAvailable         = "Keycloak disponible y respondiendo"
	LogKeycloakUnavailable       = "Keycloak no disponible"
	LogKeycloakConnectionError   = "Error de conexión con Keycloak"
	LogKeycloakTimeoutError      = "Timeout en conexión con Keycloak"
)

// ============================================
// DATABASE AVAILABILITY
// ============================================
const (
	LogDatabaseAvailabilityCheck = "Verificando disponibilidad de base de datos"
	LogDatabaseAvailable         = "Base de datos disponible y respondiendo"
	LogDatabaseUnavailable       = "Base de datos no disponible"
	LogDatabaseConnectionError   = "Error de conexión con base de datos"
)

// ============================================
// DUAL SYSTEM VALIDATION
// ============================================
const (
	LogDualSystemCheck          = "Validando existencia en ambos sistemas"
	LogUserExistsInBoth         = "Usuario existe en ambos sistemas"
	LogUserExistsOnlyInDB       = "Usuario existe solo en base de datos"
	LogUserExistsOnlyInKeycloak = "Usuario existe solo en Keycloak"
	LogUserNotFoundInEither     = "Usuario no encontrado en ningún sistema"
	LogInconsistentStateDetect  = "Estado inconsistente detectado entre sistemas"
)

// ============================================
// REPOSITORY / MESSAGE REPOSITORY
// ============================================
const (
	LogRepoMsgInit         = "Inicializando repositorio de mensajes"
	LogRepoMsgInitOK       = "Repositorio de mensajes inicializado"
	LogRepoMsgInitError    = "Error inicializando repositorio de mensajes"
	LogPersonRepoInit      = "Inicializando repositorio de personas"
	LogPersonRepoInitOK    = "Repositorio de personas inicializado"
	LogPersonRepoInitError = "Error inicializando repositorio de personas"
	LogIDEncoderInit       = "Inicializando ID encoder"
	LogIDEncoderInitError  = "Error inicializando ID encoder"
)

// ============================================
// HTTP REQUESTS
// ============================================
const (
	LogHTTPRequestIncoming = "Request HTTP entrante"
	LogHTTPResponseSent    = "Respuesta HTTP enviada"
	LogHTTPClientIP        = "IP del cliente"
	LogHTTPMethod          = "Método HTTP"
	LogHTTPPath            = "Path HTTP"
	LogHTTPStatus          = "Status HTTP"
)

// ============================================
// PERSON / USER SERVICES
// ============================================
const (
	LogPersonCreating    = "Creando persona en base de datos"
	LogPersonCreated     = "Persona creada exitosamente"
	LogPersonCreateError = "Error creando persona"
	LogPersonUpdating    = "Actualizando persona"
	LogPersonUpdated     = "Persona actualizada"
	LogPersonUpdateError = "Error actualizando persona"
	LogPersonSearching   = "Buscando persona"
	LogPersonFound       = "Persona encontrada"
	LogPersonNotFound    = "Persona no encontrada"
	LogPersonSearchError = "Error buscando persona"
)

// ============================================
// DEPENDENCY INJECTION
// ============================================
const (
	LogDepInit          = "Inicializando dependencies"
	LogDepInitComplete  = "Dependencies inicializadas completamente"
	LogDepInitError     = "Error inicializando dependencies"
	LogDepWiringService = "Inyectando servicio"
)

// ============================================
// SECURITY / ENCODING
// ============================================
const (
	LogIDEncode               = "Codificando ID"
	LogIDEncodeOK             = "ID codificado exitosamente"
	LogIDEncodeError          = "Error codificando ID"
	LogIDDecode               = "Decodificando ID"
	LogIDDecodeOK             = "ID decodificado exitosamente"
	LogIDDecodeError          = "Error decodificando ID"
	LogIDEncoderInvalidUUID   = "UUID inválido"
	LogIDEncoderHashidsCreate = "Error creando hashids"
	LogIDEncoderEncodingError = "Error encodeando UUID"
	LogIDEncoderEmptyID       = "ID ofuscado no puede estar vacío"
	LogIDEncoderDecodingError = "Error decodeando ID ofuscado"
	LogIDEncoderInvalidFormat = "ID ofuscado tiene formato incorrecto"
	LogIDEncoderUUIDError     = "Error reconstruyendo UUID"
	LogIDEncoderMinLengthWarn = "MinLength es igual a 36, lo cual es el valor por defecto"
)

// ============================================
// UTILS / HELPERS
// ============================================
const (
	LogUtilsModuleRootSearch = "Buscando raíz del módulo"
	LogUtilsModuleRootFound  = "Raíz del módulo encontrada"
	LogUtilsModuleRootError  = "No se pudo encontrar la raíz del módulo"
	LogUtilsCurrentDirError  = "No se pudo determinar el directorio actual"
	LogUtilsPathResolved     = "Ruta resuelta exitosamente"
	LogUtilsPathError        = "Error resolviendo ruta"
)

// ============================================
// MIDDLEWARE / VALIDATORS
// ============================================
const (
	LogMiddlewareValidationStart    = "Iniciando validación de request body"
	LogMiddlewareValidationOK       = "Validación exitosa"
	LogMiddlewareValidationFailed   = "Validación de request fallida"
	LogMiddlewareBodyReadError      = "Error leyendo body del request"
	LogMiddlewareJSONParseError     = "Error parseando JSON del body"
	LogMiddlewareSchemaError        = "Error de validación de schema"
	LogMiddlewareResponseCacheError = "Error obteniendo mensaje de cache"
	LogMiddlewareResponseSuccess    = "Respuesta enviada exitosamente"
	LogMiddlewareNotFound           = "Endpoint no encontrado (404)"
	LogMiddlewareRateLimitHit       = "Rate limit excedido"
	LogMiddlewareTypeCastError      = "Error de conversión de tipo en middleware"
)

// ============================================
// MESSAGES MODULE / SYSTEM MESSAGES
// ============================================
const (
	LogMessageServiceInit       = "Inicializando servicio de mensajes"
	LogMessageServiceInitOK     = "Servicio de mensajes inicializado"
	LogMessageCreate            = "Creando mensaje del sistema"
	LogMessageCreateOK          = "Mensaje creado exitosamente"
	LogMessageCreateError       = "Error creando mensaje"
	LogMessageCreateProcessing  = "Procesando creación de mensaje"
	LogMessageUpdate            = "Actualizando mensaje del sistema"
	LogMessageUpdateOK          = "Mensaje actualizado exitosamente"
	LogMessageUpdateError       = "Error actualizando mensaje"
	LogMessageUpdateProcessing  = "Procesando actualización de mensaje"
	LogMessageDelete            = "Eliminando mensaje del sistema"
	LogMessageDeleteOK          = "Mensaje eliminado exitosamente"
	LogMessageDeleteError       = "Error eliminando mensaje"
	LogMessageDeleteProcessing  = "Procesando eliminación de mensaje"
	LogMessageGet               = "Obteniendo mensaje del sistema"
	LogMessageGetOK             = "Mensaje obtenido exitosamente"
	LogMessageGetError          = "Error obteniendo mensaje"
	LogMessageList              = "Listando mensajes del sistema"
	LogMessageListOK            = "Mensajes listados exitosamente"
	LogMessageListError         = "Error listando mensajes"
	LogMessageValidation        = "Validando mensaje"
	LogMessageValidationOK      = "Mensaje validado exitosamente"
	LogMessageValidationError   = "Error validando mensaje"
	LogMessageCodeDuplicate     = "Código de mensaje duplicado"
	LogMessageTxBegin           = "Iniciando transacción para mensaje"
	LogMessageTxBeginOK         = "Transacción iniciada para mensaje"
	LogMessageTxCommit          = "Confirmando transacción de mensaje"
	LogMessageTxCommitOK        = "Transacción confirmada exitosamente"
	LogMessageTxCommitError     = "Error confirmando transacción"
	LogMessageTxRollback        = "Ejecutando rollback de transacción"
	LogMessageTxRollbackOK      = "Rollback ejecutado exitosamente"
	LogMessageTxRollbackError   = "Error ejecutando rollback"
	LogMessageCacheRefresh      = "Refrescando cache de mensajes"
	LogMessageCacheRefreshOK    = "Cache de mensajes refrescado"
	LogMessageCacheRefreshError = "Error refrescando cache de mensajes"
	LogMessageInvalidID         = "ID inválido"
	LogMessageIDEncodeError     = "Error ofuscando ID"
	LogMessageIDDecodeError     = "Error decodificando ID"
)

// ============================================
// PROMETHEUS / OBSERVABILITY
// ============================================
const (
	LogPrometheusInit          = "Inicializando métricas de Prometheus"
	LogPrometheusInitOK        = "Métricas de Prometheus inicializadas correctamente"
	LogPrometheusInitError     = "Error inicializando métricas de Prometheus"
	LogPrometheusMetricRecord  = "Registrando métrica"
	LogPrometheusMetricError   = "Error registrando métrica"
	LogPrometheusScrapeSuccess = "Scraping de métricas exitoso"
	LogPrometheusScrapeError   = "Error durante scraping de métricas"
)

// ============================================
// PERSON SERVICES
// ============================================
const (
	LogPersonServiceSearchByEmail      = "Buscando persona por email"
	LogPersonServiceSearchByID         = "Buscando persona por ID"
	LogPersonServiceSearchByKeycloakID = "Buscando persona por Keycloak ID"
	LogPersonServiceFoundByEmail       = "Persona encontrada por email"
	LogPersonServiceFoundByID          = "Persona encontrada por ID"
	LogPersonServiceFoundByKeycloakID  = "Persona encontrada por Keycloak ID"
	LogPersonServiceErrorByEmail       = "Error buscando persona por email"
	LogPersonServiceErrorByID          = "Error buscando persona por ID"
	LogPersonServiceErrorByKeycloakID  = "Error buscando persona por Keycloak ID"
	LogPersonServiceValidationStart    = "Iniciando validaciones de registro"

	LogPersonServiceValidationComplete        = "Validaciones de registro completadas"
	LogPersonServiceDuplicateEmail            = "Intento de registro con email duplicado"
	LogPersonServiceSavingToDB                = "Guardando persona en base de datos"
	LogPersonServiceSavedToDB                 = "Persona guardada en base de datos"
	LogPersonServiceSaveError                 = "Error guardando persona en BD"
	LogPersonServiceCreatingKeycloak          = "Creando usuario en Keycloak"
	LogPersonServiceCreatedKeycloak           = "Usuario creado en Keycloak"
	LogPersonServiceKeycloakError             = "Error creando usuario en Keycloak"
	LogPersonServicePasswordSet               = "Configurando password de usuario"
	LogPersonServicePasswordSetOK             = "Password configurado"
	LogPersonServicePasswordError             = "Error configurando password"
	LogPersonServiceRoleAssigning             = "Asignando rol a usuario"
	LogPersonServiceRoleAssigned              = "Rol asignado"
	LogPersonServiceRoleError                 = "Error asignando rol"
	LogPersonServiceKeycloakIDUpdate          = "Actualizando keycloak_user_id en BD"
	LogPersonServiceKeycloakIDUpdated         = "Keycloak_user_id actualizado"
	LogPersonServiceKeycloakIDUpdateError     = "Error actualizando keycloak_user_id"
	LogPersonServiceRollbackPerson            = "Ejecutando rollback: eliminando persona de BD"
	LogPersonServiceRollbackPersonError       = "Error en rollback de persona"
	LogPersonServiceRollbackPersonComplete    = "Rollback de persona completado"
	LogPersonServiceRollbackKeycloak          = "Ejecutando rollback: eliminando usuario de Keycloak"
	LogPersonServiceRollbackKeycloakError     = "Error en rollback de usuario Keycloak"
	LogPersonServiceRollbackKeycloakComplete  = "Rollback de usuario Keycloak completado"
	LogPersonServiceInconsistentStateDetected = "Estado inconsistente detectado entre Keycloak y BD"
	LogPersonServiceCleaningOrphan            = "Limpiando usuario huérfano"
	LogPersonServiceOrphanCleaned             = "Usuario huérfano eliminado exitosamente"
	LogPersonServiceOrphanCleanError          = "Error limpiando usuario huérfano"
)

// ============================================
// PERSON INTERACTOR
// ============================================
const (
	LogPersonInteractorRegStart             = "Iniciando proceso de registro"
	LogPersonInteractorStep1_Error          = "[PASO 1/8] Validaciones fallidas"
	LogPersonInteractorStep1_OK             = "[PASO 1/8] Validaciones completadas"
	LogPersonInteractorIDGenerated          = "ID generado para persona"
	LogPersonInteractorStep15_Error         = "[PASO 1.5/8] Estado inconsistente detectado y limpiado"
	LogPersonInteractorStep15_OK            = "[PASO 1.5/8] Estado consistente verificado"
	LogPersonInteractorStep2_Error          = "[PASO 2/8] Error iniciando transacción"
	LogPersonInteractorStep2_OK             = "[PASO 2/8] Transacción iniciada"
	LogPersonInteractorStep3_Error          = "[PASO 3/8] Error guardando persona"
	LogPersonInteractorStep3_OK             = "[PASO 3/8] Persona guardada en BD"
	LogPersonInteractorStep4_Error          = "[PASO 4/8] Error creando usuario en Keycloak"
	LogPersonInteractorStep4_OK             = "[PASO 4/8] Usuario creado en Keycloak"
	LogPersonInteractorStep5_Error          = "[PASO 5/8] Error configurando password"
	LogPersonInteractorStep5_OK             = "[PASO 5/8] Password configurado"
	LogPersonInteractorStep6_Error          = "[PASO 6/8] Error asignando rol"
	LogPersonInteractorStep6_OK             = "[PASO 6/8] Rol asignado"
	LogPersonInteractorStep7_Error          = "[PASO 7/8] Error actualizando Keycloak ID en BD"
	LogPersonInteractorStep7_OK             = "[PASO 7/8] Keycloak_user_id actualizado en BD"
	LogPersonInteractorCommit_Error         = LogBranchInteractorCommitError
	LogPersonInteractorCommit_OK            = LogMessageTxCommitOK
	LogPersonInteractorRegComplete          = "Registro completado exitosamente"
	LogPersonInteractorRollbackDB_Error     = LogTxRollbackError
	LogPersonInteractorRollbackDB_OK        = LogTxRollbackOK
	LogPersonInteractorRollbackKeycloak_Err = "ROLLBACK KEYCLOAK FALLÓ - ALERTA CRÍTICA"
	LogPersonInteractorRollbackKeycloak_OK  = "Rollback Keycloak ejecutado correctamente"
	LogPersonInteractorIncompleteDetected   = "Registro incompleto detectado"
	LogPersonInteractorCleanup_Error        = "Error limpiando estado inconsistente"
	LogPersonInteractorCleanup_OK           = "Estado inconsistente limpiado exitosamente"
	LogPersonInteractorLoginStart           = "Inicio de sesión en Keycloak"
	LogPersonInteractorLoginError           = "Error iniciando sesión en Keycloak"
	LogPersonInteractorLoginOK              = "Sesión iniciada exitosamente"
	LogPersonInteractorLoginComplete        = "Sesión iniciada exitosamente"
	// Refresh Token / Contact / Delete
	LogPersonInteractorRefreshOK        = "Refresh token completado exitosamente"
	LogPersonInteractorContactGetOK     = "Contacto público obtenido exitosamente"
	LogPersonInteractorKeycloakDeleteOK = "Usuario eliminado de Keycloak"
	LogPersonInteractorPersonDeleteOK   = "Persona eliminada de base de datos"
)

// ============================================
// PERSON SERVICE
// ============================================
const (
	LogPersonServiceRefreshStart            = "RefreshToken llamado"
	LogPersonServiceRefreshError            = "RefreshToken fallido"
	LogPersonServiceRefreshOK               = "RefreshToken completado exitosamente"
	LogPersonServiceEmailExtracted          = "Email extraído del token"
	LogPersonServicePasswordPolicyViolation = "Violación de política de contraseña"
)

// ============================================
// DEPENDENCY INITIALIZATION
// ============================================
const (
	LogDependencyMessageRepoInit = "MessageRepository inicializado"
	LogDependencyMessageIntInit  = "MessageInteractor inicializado"
)

// ============================================
// DATABASE CONNECTION (MySQL)
// ============================================
const (
	LogDBConnecting      = "Conectando a base de datos MySQL"
	LogDBSSLEnabled      = "SSL habilitado para conexión a base de datos"
	LogDBConnectionError = "Error abriendo conexión a base de datos"
	LogDBPoolConfig      = "Configurando pool de conexiones"
	LogDBPinging         = "Verificando conectividad con base de datos (ping)..."
	LogDBPingError       = "Error en ping a base de datos"
	LogDBConnected       = "Conexión a base de datos establecida exitosamente"
)

// ============================================
// MESSAGE INTERACTOR
// ============================================
const (
	// CREATE flow
	LogMessageInteractorCreateStep1Error = "[PASO 1/3] Validación de mensaje fallida"
	LogMessageInteractorCreateStep1OK    = "[PASO 1/3] Validación de mensaje completada"
	LogMessageInteractorCreateStep2Error = "[PASO 2/3] Error iniciando transacción"
	LogMessageInteractorCreateStep2OK    = "[PASO 2/3] Transacción iniciada"
	LogMessageInteractorCreateStep3Error = "[PASO 3/3] Error guardando mensaje"
	LogMessageInteractorCreateStep3OK    = "[PASO 3/3] Mensaje guardado en BD"
	LogMessageInteractorCreateCommitErr  = LogBranchInteractorCommitError
	LogMessageInteractorCreateCommitOK   = LogMessageTxCommitOK
	LogMessageInteractorCreateComplete   = LogMessageCreateOK

	// UPDATE flow
	LogMessageInteractorUpdateStep1Error = "[PASO 1/4] Mensaje no encontrado"
	LogMessageInteractorUpdateStep1OK    = "[PASO 1/4] Mensaje encontrado"
	LogMessageInteractorUpdateStep2Error = "[PASO 2/4] Validación de mensaje fallida"
	LogMessageInteractorUpdateStep2OK    = "[PASO 2/4] Validación de mensaje completada"
	LogMessageInteractorUpdateStep3Error = "[PASO 3/4] Error iniciando transacción"
	LogMessageInteractorUpdateStep3OK    = "[PASO 3/4] Transacción iniciada"
	LogMessageInteractorUpdateStep4Error = "[PASO 4/4] Error actualizando mensaje"
	LogMessageInteractorUpdateStep4OK    = "[PASO 4/4] Mensaje actualizado en BD"
	LogMessageInteractorUpdateCommitErr  = LogBranchInteractorCommitError
	LogMessageInteractorUpdateCommitOK   = LogMessageTxCommitOK
	LogMessageInteractorUpdateComplete   = LogMessageUpdateOK

	// DELETE flow
	LogMessageInteractorDeleteStep1Error = "[PASO 1/3] Mensaje no encontrado"
	LogMessageInteractorDeleteStep1OK    = "[PASO 1/3] Mensaje encontrado"
	LogMessageInteractorDeleteStep2Error = "[PASO 2/3] Error iniciando transacción"
	LogMessageInteractorDeleteStep2OK    = "[PASO 2/3] Transacción iniciada"
	LogMessageInteractorDeleteStep3Error = "[PASO 3/3] Error eliminando mensaje"
	LogMessageInteractorDeleteStep3OK    = "[PASO 3/3] Mensaje eliminado de BD"
	LogMessageInteractorDeleteCommitErr  = LogBranchInteractorCommitError
	LogMessageInteractorDeleteCommitOK   = LogMessageTxCommitOK
	LogMessageInteractorDeleteComplete   = LogMessageDeleteOK

	// Common rollback
	LogMessageInteractorRollbackError = LogTxRollbackError
	LogMessageInteractorRollbackOK    = LogTxRollbackOK
)

// ============================================
// BRANCH INTERACTOR (HU59)
// ============================================
const (
	LogBranchInteractorRegStart        = "Iniciando proceso de registro de sede"
	LogBranchInteractorValidationError = "Error de validación de sede"
	LogBranchInteractorBrandsValidated = "Marcas validadas correctamente"
	LogBranchInteractorIDGenerated     = "ID generado para sede"
	LogBranchInteractorTxError         = "Error iniciando transacción"
	LogBranchInteractorTxStarted       = "Transacción iniciada"
	LogBranchInteractorRegError        = "Error registrando sede"
	LogBranchInteractorRegSaved        = "Sede guardada en BD"
	LogBranchInteractorCommitError     = "COMMIT FALLÓ - ALERTA CRÍTICA"
	LogBranchInteractorRegComplete     = "Sede registrada exitosamente"
	LogBranchInteractorRollbackError   = LogTxRollbackError
	LogBranchInteractorRollbackOK      = LogTxRollbackOK
	LogBranchInteractorGetByID         = "Obteniendo sede por ID"
	LogBranchInteractorGetByIDError    = "Error obteniendo sede por ID"
	LogBranchInteractorGetByIDOK       = "Sede obtenida exitosamente"
	// HU62: List branches by representative
	LogBranchInteractorListByRep    = "Listando sedes por representante"
	LogBranchInteractorListByRepErr = "Error listando sedes por representante"
	LogBranchInteractorListByRepOK  = "Sedes listadas por representante exitosamente"
	// HU60: Update branch
	LogBranchInteractorUpdateStart    = "Starting branch update process"
	LogBranchInteractorOwnershipError = "User is not the owner of this branch"
	LogBranchInteractorUpdateError    = "Error updating branch"
	LogBranchInteractorUpdateComplete = "Branch updated successfully"
	// HU61: Delete branch
	LogBranchInteractorDeleteStart    = "Starting branch delete process"
	LogBranchInteractorDeleteError    = "Error deleting branch"
	LogBranchInteractorDeleteComplete = "Branch deleted successfully"
	LogBranchInteractorHasAssocError  = "Branch has associations that prevent deletion"
	// HU89: Search nearby branches
	LogBranchInteractorNearbyStart    = "Searching branches near location"
	LogBranchInteractorNearbyError    = "Error searching nearby branches"
	LogBranchInteractorNearbyComplete = "Nearby branches search completed"
)

// ============================================
// BRANCH REPOSITORY (HU59)
// ============================================
const (
	LogBranchRepoSaveError        = "Error guardando sede en BD"
	LogBranchRepoUpdateError      = "Error actualizando sede en BD"
	LogBranchRepoDeleteError      = "Error eliminando sede de BD"
	LogBranchRepoGetByIDError     = LogBranchInteractorGetByIDError
	LogBranchRepoGetByNameError   = "Error obteniendo sede por nombre"
	LogBranchRepoGetByRepError    = "Error obteniendo sedes por representante"
	LogBranchRepoScanError        = "Error escaneando fila de sede"
	LogBranchRepoLocationSaveErr  = "Error guardando ubicación"
	LogBranchRepoLocationUpdErr   = "Error actualizando ubicación"
	LogBranchRepoBrandSaveError   = "Error guardando marca de sede"
	LogBranchRepoBrandDelError    = "Error eliminando marcas de sede"
	LogBranchRepoBrandGetError    = "Error obteniendo marcas de sede"
	LogBranchRepoBrandValidateErr = "Error validando marcas"
	// Displacement range pivot
	LogBranchRepoDisplRangeSaveError = "Error guardando rango de cilindraje de sede"
	LogBranchRepoDisplRangeDelError  = "Error eliminando rangos de cilindraje de sede"
	LogBranchRepoDisplRangeGetError  = "Error obteniendo rangos de cilindraje de sede"
)

// ============================================
// BRANCH SERVICE (HU59)
// ============================================
const (
	LogBranchServiceInvalidType  = "Tipo de establecimiento inválido"
	LogBranchServiceDupNameCheck = "Error verificando nombre duplicado"
	LogBranchServiceDupName      = "Nombre de sede duplicado en franquicia"
	LogBranchServiceSaveError    = "Error guardando sede"
	LogBranchServiceLocSaveError = LogBranchRepoLocationSaveErr
	LogBranchServiceBrandSaveErr = "Error guardando marcas"
	LogBranchServiceRegComplete  = "Sede registrada exitosamente"
	LogBranchServiceGetError     = LogBranchInteractorGetByIDError
	// HU61: Delete branch
	LogBranchServiceDelError    = "Error eliminando sede"
	LogBranchServiceDelComplete = "Sede eliminada exitosamente"
)

// ============================================
// BRANCH CONTROLLER (HU59, HU62)
// ============================================
const (
	LogBranchControllerRegRequest    = "Solicitud de registro de sede recibida"
	LogBranchControllerUserAuth      = "Usuario autenticado"
	LogBranchControllerUserUnauth    = "Usuario no autenticado intentando registrar sede"
	LogBranchControllerRoleForbidden = "Usuario sin rol de representante intentando registrar sede"
	LogBranchControllerBindError     = "Error parseando JSON de solicitud"
	LogBranchControllerProcessing    = "Procesando registro de sede"
	LogBranchControllerRegError      = "Error registrando sede"
	LogBranchControllerRegSuccess    = "Sede registrada exitosamente en controller"
	// HU62: Get Branch
	LogBranchControllerGetRequest    = "Solicitud de consulta de sede recibida"
	LogBranchControllerIDDecodeError = "Error decodificando ID de sede"
	LogBranchControllerGetByID       = "Buscando sede por ID"
	LogBranchControllerGetError      = "Error obteniendo sede"
	LogBranchControllerGetSuccess    = "Sede obtenida exitosamente"
	// HU76: Get Branch Types
	LogBranchControllerGetTypes   = "Solicitud de tipos de sede recibida"
	LogBranchControllerGetTypesOK = "Tipos de sede obtenidos exitosamente"
	// HU62: List my branches
	LogBranchControllerListRequest = "Solicitud de listado de sedes recibida"
	LogBranchControllerListError   = "Error listando sedes"
	LogBranchControllerListSuccess = "Sedes listadas exitosamente"
	// HU60: Update branch
	LogBranchControllerUpdateRequest = "Branch update request received"
	LogBranchControllerUpdateError   = "Error updating branch"
	LogBranchControllerUpdateSuccess = "Branch updated successfully"
	// HU61: Delete branch
	LogBranchControllerDeleteRequest = "Branch delete request received"
	LogBranchControllerDeleteError   = "Error deleting branch"
	LogBranchControllerDeleteSuccess = "Branch deleted successfully"
)

// ============================================
// BRAND INTERACTOR
// ============================================
const (
	LogBrandInteractorGetAll      = "Obteniendo lista de marcas"
	LogBrandInteractorGetAllOK    = "Lista de marcas obtenida exitosamente"
	LogBrandInteractorGetAllError = "Error obteniendo lista de marcas"
)

// ============================================
// GEOCODING SERVICE (OpenCage)
// ============================================
const (
	LogGeocodingRequest   = "geocoding_request_initiated"
	LogGeocodingSuccess   = "geocoding_completed_successfully"
	LogGeocodingNoResults = "geocoding_no_results_found"
	LogGeocodingError     = "geocoding_request_failed"
	LogGeocodingSkipped   = "geocoding_skipped_coordinates_present"
	LogGeocodingCityError = "geocoding_city_lookup_failed"
	LogGeocodingRateLimit = "geocoding_rate_limit_info"
)

// ============================================
// GEOCODING CONTROLLER
// ============================================
const (
	LogGeocodingControllerTestRequest = "Solicitud de prueba de geocodificación recibida"
	LogGeocodingControllerTestInvalid = "Solicitud de prueba de geocodificación inválida"
	LogGeocodingControllerTestSuccess = "Prueba de geocodificación exitosa"
	LogGeocodingControllerTestFailed  = "Prueba de geocodificación fallida"
)

// ============================================
// PERSON CONTROLLER
// ============================================
const (
	LogPersonControllerRegComplete      = "Registro de usuario completado exitosamente"
	LogPersonControllerTokenRefreshOK   = "Token de usuario refrescado exitosamente"
	LogPersonControllerProfileGetOK     = "Perfil de usuario obtenido exitosamente"
	LogPersonControllerUserNotInContext = "Usuario autenticado no encontrado en contexto"
	LogPersonControllerIDEncodeError    = "Error codificando ID de usuario"
	LogPersonControllerContactGetOK     = "Información de contacto público obtenida"
	LogPersonControllerAccountDeleteOK  = "Cuenta de usuario eliminada exitosamente"
)

// ============================================
// BRANCH CONTROLLER - NEARBY
// ============================================
const (
	LogBranchControllerNearbyFound = "Sucursales cercanas encontradas"
)

// ============================================
// LOCATION INTERACTOR
// ============================================
const (
	LogLocationInteractorGetDepartments      = "Obteniendo lista de departamentos"
	LogLocationInteractorGetDepartmentsOK    = "Lista de departamentos obtenida exitosamente"
	LogLocationInteractorGetDepartmentsError = "Error obteniendo departamentos"
	LogLocationInteractorGetCities           = "Obteniendo ciudades del departamento"
	LogLocationInteractorGetCitiesOK         = "Lista de ciudades obtenida exitosamente"
	LogLocationInteractorGetCitiesError      = "Error obteniendo ciudades"
)

// ============================================
// LOCATION REPOSITORY
// ============================================
const (
	LogLocationRepoGetDepartmentsError     = LogLocationInteractorGetDepartmentsError
	LogLocationRepoGetDepartmentsScanError = "Error escaneando departamento"
	LogLocationRepoGetDepartmentsIterError = "Error iterando departamentos"
	LogLocationRepoGetCitiesError          = LogLocationInteractorGetCitiesError
	LogLocationRepoGetCitiesScanError      = "Error escaneando ciudad"
	LogLocationRepoGetCitiesIterError      = "Error iterando ciudades"
	LogLocationRepoValidateCityError       = "Error validando ciudad en departamento"
	LogLocationRepoGetDeptByIDError        = "Error obteniendo departamento por ID"
	LogLocationRepoSaveError               = LogBranchRepoLocationSaveErr
	LogLocationRepoUpdateError             = "Error actualizando ubicación"
	LogLocationRepoPrepareError            = "Error preparando statement de ubicación"
)

// ============================================
// FIREBASE CLIENT
// ============================================
const (
	LogFirebaseInitApp          = "Inicializando aplicación Firebase"
	LogFirebaseInitAppError     = "Error inicializando aplicación Firebase"
	LogFirebaseAuthClientError  = "Error obteniendo cliente de autenticación Firebase"
	LogFirebaseInitOK           = "Firebase Admin SDK inicializado exitosamente"
	LogFirebaseTokenCreate      = "Creando token personalizado Firebase"
	LogFirebaseTokenCreateOK    = "Token personalizado Firebase creado"
	LogFirebaseTokenCreateError = "Error creando token personalizado Firebase"
	LogFirebaseTokenClaimsOK    = "Token con claims personalizado creado"
	LogFirebaseTokenClaimsError = "Error creando token con claims"
	// Storage operations (HU39/HU45)
	LogFirebaseStorageClientWarn     = "No se pudo inicializar cliente de Storage Firebase (continuando sin Storage)"
	LogFirebaseStorageNotConfigured  = "Cliente de Storage Firebase no configurado - delete ignorado"
	LogFirebaseStorageParseError     = "Error parseando URL de Firebase Storage"
	LogFirebaseStorageDeleteOK       = "Archivo eliminado de Firebase Storage exitosamente"
	LogFirebaseStorageDeleteError    = "Error eliminando archivo de Firebase Storage"
	LogFirebaseStorageAlreadyDeleted = "Archivo ya eliminado de Firebase Storage (no existe)"
)

// ============================================
// BRAND REPOSITORY
// ============================================
const (
	LogBrandRepoValidateError = "Error validando marcas"
	LogBrandRepoNotFound      = "Marca no encontrada en catálogo"
)

// ============================================
// FRANCHISE REPOSITORY (HU26-29)
// ============================================
const (
	LogFranchiseRepoSaveError       = "Error guardando franquicia en BD"
	LogFranchiseRepoUpdateError     = "Error actualizando franquicia en BD"
	LogFranchiseRepoDeleteError     = "Error eliminando franquicia de BD"
	LogFranchiseRepoGetByIDError    = "Error obteniendo franquicia por ID"
	LogFranchiseRepoGetByNameError  = "Error obteniendo franquicia por nombre"
	LogFranchiseRepoGetByRepError   = "Error obteniendo franquicias por representante"
	LogFranchiseRepoScanError       = "Error escaneando fila de franquicia"
	LogFranchiseRepoCountBranches   = "Error contando sedes de franquicia"
	LogFranchiseRepoAssociateError  = "Error asociando sede a franquicia"
	LogFranchiseRepoDissociateError = "Error disociando sedes de franquicia"
	LogFranchiseRepoPrepareError    = "Error preparando statement de franquicia"
)

// ============================================
// FRANCHISE SERVICE (HU26-29)
// ============================================
const (
	LogFranchiseServiceDupNameCheck = "Error verificando nombre duplicado de franquicia"
	LogFranchiseServiceDupName      = "Nombre de franquicia duplicado"
	LogFranchiseServiceSaveError    = "Error guardando franquicia"
	LogFranchiseServiceGetError     = "Error obteniendo franquicia por ID"
	LogFranchiseServiceUpdateError  = LogFranchiseInteractorUpdateError
	LogFranchiseServiceDeleteError  = LogFranchiseInteractorDeleteError
	LogFranchiseServiceDeleted      = LogFranchiseInteractorDeleteComplete
)

// ============================================
// FRANCHISE INTERACTOR (HU26-29)
// ============================================
const (
	LogFranchiseInteractorCreateStart    = "Iniciando creación de franquicia"
	LogFranchiseInteractorNoBranches     = "Franquicia requiere al menos una sede"
	LogFranchiseInteractorBranchNotOwned = "Sede no pertenece al representante"
	LogFranchiseInteractorTxError        = "Error iniciando transacción para franquicia"
	LogFranchiseInteractorCreateError    = "Error creando franquicia"
	LogFranchiseInteractorCreateComplete = "Franquicia creada exitosamente"
	LogFranchiseInteractorUpdateStart    = "Iniciando actualización de franquicia"
	LogFranchiseInteractorUpdateError    = "Error actualizando franquicia"
	LogFranchiseInteractorUpdateComplete = "Franquicia actualizada exitosamente"
	LogFranchiseInteractorDeleteStart    = "Iniciando eliminación de franquicia"
	LogFranchiseInteractorDeleteError    = "Error eliminando franquicia"
	LogFranchiseInteractorDeleteComplete = "Franquicia eliminada exitosamente"
	LogFranchiseInteractorCommitError    = "COMMIT FALLÓ - Franquicia"
	LogFranchiseInteractorRollbackError  = "ROLLBACK FALLÓ - Franquicia"
	LogFranchiseInteractorRollbackOK     = "Rollback ejecutado correctamente (franquicia)"

	// Franchise Controller
	LogFranchiseControllerRequest          = "Solicitud recibida - Franquicia"
	LogFranchiseControllerProcessing       = "Procesando franquicia"
	LogFranchiseControllerCreateError      = "Error creando franquicia"
	LogFranchiseControllerCreateSuccess    = "Franquicia creada exitosamente"
	LogFranchiseControllerListRequest      = "Solicitud de listado de franquicias"
	LogFranchiseControllerListError        = "Error listando franquicias"
	LogFranchiseControllerListSuccess      = "Franquicias listadas exitosamente"
	LogFranchiseControllerGetError         = "Error obteniendo franquicia"
	LogFranchiseControllerGetSuccess       = "Franquicia obtenida exitosamente"
	LogFranchiseControllerBindError        = "Error parseando JSON de franquicia"
	LogFranchiseControllerIDDecodeError    = "Error decodificando ID de franquicia"
	LogFranchiseControllerUpdateError      = LogFranchiseInteractorUpdateError
	LogFranchiseControllerUpdateSuccess    = "Franquicia actualizada exitosamente"
	LogFranchiseControllerDeleteError      = LogFranchiseInteractorDeleteError
	LogFranchiseControllerDeleteSuccess    = LogFranchiseInteractorDeleteComplete
	LogFranchiseControllerAddBranchRequest = "Solicitud vincular sede a franquicia"
	LogFranchiseControllerAddBranchError   = "Error vinculando sede a franquicia"
	LogFranchiseControllerAddBranchSuccess = "Sede vinculada a franquicia exitosamente"
	LogFranchiseControllerRemBranchRequest = "Solicitud desvincular sede de franquicia"
	LogFranchiseControllerRemBranchError   = "Error desvinculando sede de franquicia"
	LogFranchiseControllerRemBranchSuccess = "Sede desvinculada de franquicia exitosamente"
	LogFranchiseInteractorMinBranches      = "No se puede desvincular la última sede"
)

// ============================================
// DEPENDENCY INITIALIZATION
// ============================================
const (
	// Repository initialization
	LogDepBranchRepoInit             = "Inicializando repositorio de sedes"
	LogDepBranchRepoInitOK           = "Repositorio de sedes inicializado"
	LogDepBranchRepoInitErr          = "Error inicializando repositorio de sedes"
	LogDepLocationRepoInitOK         = "Repositorio de ubicaciones inicializado"
	LogDepLocationRepoInitErr        = "Error inicializando repositorio de ubicaciones"
	LogDepBrandRepoInitOK            = "Repositorio de marcas inicializado"
	LogDepBrandRepoInitErr           = "Error inicializando repositorio de marcas"
	LogDepFranchiseRepoInitOK        = "Repositorio de franquicias inicializado"
	LogDepFranchiseRepoInitErr       = "Error inicializando repositorio de franquicias"
	LogDepScheduleRepoInitErr        = "Error inicializando repositorio de horarios"
	LogDepScheduleRepoInitOK         = "Repositorio de horarios inicializado"
	LogDepScheduleIntInitOK          = "Interactor de horarios inicializado"
	LogDepScheduleDetailRepoInitErr  = "Error inicializando repositorio de detalles de horario"
	LogDepScheduleDetailRepoInitOK   = "Repositorio de detalles de horario inicializado"
	LogDepScheduleDetailIntInitOK    = "Interactor de detalles de horario inicializado"
	LogDepScheduleExceptionIntInitOK = "Interactor de excepciones de horario inicializado"

	// Interactor initialization
	LogDepBranchInteractorInitOK    = "Interactor de sedes inicializado"
	LogDepBrandInteractorInitOK     = "Interactor de marcas inicializado"
	LogDepLocationInteractorInitOK  = "Interactor de ubicaciones inicializado"
	LogDepFranchiseInteractorInitOK = "Interactor de franquicias inicializado"

	// External services
	LogDepGeocodingClientInitOK = "Cliente de geocodificación inicializado"
	LogDepFirebaseClientInitOK  = "Cliente de Firebase inicializado"
	LogDepJWKSValidatorInitOK   = "Validador JWKS inicializado"
	LogDepJWKSValidatorInitErr  = "Error inicializando validador JWKS - usando validación fallback"

	// Firebase initialization
	LogDepFirebaseCredPathResolved = "Ruta de credenciales Firebase resuelta"
	LogDepFirebaseInitSkipped      = "Inicialización de Firebase omitida"
	LogDepFirebaseCredNotConfig    = "Credenciales de Firebase no configuradas - omitiendo inicialización"
)

// ============================================
// SERVICE CATALOG (HU63, HU68, HU75)
// ============================================
const (
	// Service Interactor
	LogServiceInteractorGetTypes       = "Obteniendo tipos de servicio"
	LogServiceInteractorGetTypesOK     = "Tipos de servicio obtenidos exitosamente"
	LogServiceInteractorGetAll         = "Obteniendo catálogo de servicios"
	LogServiceInteractorGetAllOK       = "Catálogo de servicios obtenido exitosamente"
	LogServiceInteractorGetAllError    = "Error obteniendo catálogo de servicios"
	LogServiceInteractorGetByType      = "Obteniendo servicios por tipo"
	LogServiceInteractorGetByTypeOK    = "Servicios por tipo obtenidos exitosamente"
	LogServiceInteractorGetByTypeError = "Error obteniendo servicios por tipo"
	// HU68: Service Update (Admin)
	LogServiceInteractorGetByID       = "Obteniendo servicio por ID"
	LogServiceInteractorGetByIDOK     = "Servicio obtenido por ID exitosamente"
	LogServiceInteractorGetByIDError  = "Error obteniendo servicio por ID"
	LogServiceInteractorUpdate        = "Actualizando servicio"
	LogServiceInteractorUpdateOK      = "Servicio actualizado exitosamente"
	LogServiceInteractorUpdateError   = "Error actualizando servicio"
	LogServiceInteractorRollbackError = "ROLLBACK BD FALLÓ - ALERTA CRÍTICA (servicio)"
	LogServiceInteractorRollbackOK    = "Rollback BD ejecutado correctamente (servicio)"

	// Service Repository
	LogServiceRepoGetAll         = "Consultando todos los servicios desde BD"
	LogServiceRepoGetAllError    = "Error consultando servicios desde BD"
	LogServiceRepoGetByType      = "Consultando servicios por tipo desde BD"
	LogServiceRepoGetByTypeError = "Error consultando servicios por tipo desde BD"
	LogServiceRepoScanError      = "Error escaneando fila de servicio"
	LogServiceRepoPrepareError   = "Error preparando statement de servicio"
	// HU68: Service Update (Admin)
	LogServiceRepoGetByID      = "Consultando servicio por ID desde BD"
	LogServiceRepoGetByIDOK    = "Servicio consultado por ID exitosamente"
	LogServiceRepoGetByIDError = "Error consultando servicio por ID desde BD"
	LogServiceRepoNotFound     = "Servicio no encontrado en BD"
	LogServiceRepoUpdate       = "Actualizando servicio en BD"
	LogServiceRepoUpdateOK     = "Servicio actualizado en BD exitosamente"
	LogServiceRepoUpdateError  = "Error actualizando servicio en BD"

	// Service Controller
	LogServiceControllerGetTypes    = "Solicitud de tipos de servicio recibida"
	LogServiceControllerGetTypesOK  = "Tipos de servicio enviados exitosamente"
	LogServiceControllerGetAll      = "Solicitud de catálogo de servicios recibida"
	LogServiceControllerGetAllOK    = "Catálogo de servicios enviado exitosamente"
	LogServiceControllerGetAllError = "Error obteniendo catálogo de servicios"
	LogServiceControllerInvalidType = "Tipo de servicio inválido recibido"
	// HU68: Service Update (Admin)
	LogServiceControllerUpdate      = "Solicitud de actualización de servicio recibida"
	LogServiceControllerUpdateOK    = "Servicio actualizado exitosamente"
	LogServiceControllerUpdateError = "Error actualizando servicio"

	// Dependency Initialization
	LogDepServiceRepoInitOK       = "Repositorio de servicios inicializado"
	LogDepServiceRepoInitErr      = "Error inicializando repositorio de servicios"
	LogDepServiceInteractorInitOK = "Interactor de servicios inicializado"
)

// ============================================
// BRANCH SERVICES (Service-Branch Association)
// ============================================
const (
	// Branch Services Controller
	LogBranchServicesControllerGet          = "branch_services_controller_get"
	LogBranchServicesControllerGetOK        = "branch_services_controller_get_ok"
	LogBranchServicesControllerGetError     = "branch_services_get_error"
	LogBranchServicesControllerInvalidID    = "branch_services_invalid_id"
	LogBranchServicesControllerAssociate    = "branch_services_associate"
	LogBranchServicesControllerAssociateOK  = "branch_services_associate_ok"
	LogBranchServicesControllerAssociateErr = "branch_services_associate_error"
	LogBranchServicesControllerInvalidBody  = "branch_services_invalid_body"
	LogBranchServicesControllerInvalidSvcID = "branch_services_invalid_service_id"
	LogBranchServicesControllerInvalidSvcs  = "branch_services_invalid_services"
	LogBranchServicesControllerTxError      = "branch_services_tx_error"
	LogBranchServicesControllerCommitError  = "branch_services_commit_error"
	LogBranchServicesControllerDissociate   = "branch_services_dissociate"
	LogBranchServicesControllerDissociateOK = "branch_services_dissociate_ok"
	LogBranchServicesControllerDisassocErr  = "branch_services_dissociate_error"
	LogBranchServicesControllerNotFound     = "branch_services_not_found"

	// Branch Services Repository
	LogBranchServicesRepoGetByBranch    = "service_repo_get_by_branch"
	LogBranchServicesRepoGetByBranchErr = "service_repo_get_by_branch_error"
	LogBranchServicesRepoAssociate      = "service_repo_associate"
	LogBranchServicesRepoAssociateOK    = "service_repo_associate_ok"
	LogBranchServicesRepoAssociateErr   = "service_repo_associate_error"
	LogBranchServicesRepoPrepareErr     = "service_repo_associate_prepare_error"
	LogBranchServicesRepoDissociate     = "service_repo_dissociate"
	LogBranchServicesRepoDissociateOK   = "service_repo_dissociate_ok"
	LogBranchServicesRepoDissociateErr  = "service_repo_dissociate_error"
	LogBranchServicesRepoNotFound       = "service_repo_dissociate_not_found"
	LogBranchServicesRepoValidateIDs    = "service_repo_validate_ids"
	LogBranchServicesRepoValidateErr    = "service_repo_validate_ids_error"
	LogBranchServicesRepoValidateMiss   = "service_repo_validate_ids_mismatch"
	LogBranchServicesRepoCheckAssoc     = "service_repo_check_association"
	LogBranchServicesRepoCheckAssocErr  = "service_repo_check_association_error"

	// Branch Services Interactor
	LogBranchServicesIntGetByBranch    = "service_interactor_get_by_branch"
	LogBranchServicesIntGetByBranchOK  = "service_interactor_get_by_branch_ok"
	LogBranchServicesIntGetByBranchErr = "service_interactor_get_by_branch_error"
)

// ============================================
// FIREBASE CONTROLLER
// ============================================
const (
	LogFirebaseControllerRequest    = "Solicitud de token Firebase recibida"
	LogFirebaseControllerUnauth     = "Solicitud no autenticada para token Firebase"
	LogFirebaseControllerNotConfig  = "Cliente Firebase no configurado"
	LogFirebaseControllerTokenError = "Error generando token Firebase"
	LogFirebaseControllerTokenOK    = "Token Firebase generado exitosamente"
)

// ============================================
// LOCATION CONTROLLER
// ============================================
const (
	LogLocationControllerGetDepts      = "Solicitud de listado de departamentos recibida"
	LogLocationControllerGetDeptsError = LogLocationInteractorGetDepartmentsError
	LogLocationControllerGetDeptsOK    = "Departamentos obtenidos exitosamente"
	LogLocationControllerGetCities     = "Solicitud de listado de ciudades recibida"
	LogLocationControllerGetCitiesErr  = LogLocationInteractorGetCitiesError
	LogLocationControllerGetCitiesOK   = "Ciudades obtenidas exitosamente"
)

// ============================================
// MESSAGE CACHE CONTROLLER
// ============================================
const (
	LogMessageCacheReloadRequest = "Solicitud de recarga de caché de mensajes recibida"
	LogMessageCacheReloadError   = "Error al recargar caché de mensajes"
	LogMessageCacheReloadOK      = "Caché de mensajes recargado exitosamente"
	LogMessageCreatedOK          = LogMessageCreateOK
	LogMessageUpdatedOK          = LogMessageUpdateOK
	LogMessageDeletedOK          = LogMessageDeleteOK
)

// ============================================
// SCHEDULE REPOSITORY (HU30-35)
// ============================================
const (
	LogScheduleRepoPrepareError     = "Error preparando statement de horario"
	LogScheduleRepoSaveError        = "Error guardando horario en BD"
	LogScheduleRepoUpdateError      = "Error actualizando horario en BD"
	LogScheduleRepoDeleteError      = "Error eliminando horario de BD"
	LogScheduleRepoGetByIDError     = "Error obteniendo horario por ID"
	LogScheduleRepoGetByBranchError = "Error obteniendo horario por sede"
	LogScheduleRepoActivateError    = "Error activando/desactivando horario"
)

// ============================================
// SCHEDULE SERVICE (HU30-35)
// ============================================
const (
	LogScheduleServiceCreateStart    = "Iniciando creación de horario"
	LogScheduleServiceBranchNotFound = "Sede no encontrada para horario"
	LogScheduleServiceAlreadyExists  = "La sede ya tiene un horario configurado"
	LogScheduleServiceSaveError      = "Error guardando horario"
	LogScheduleServiceCreateOK       = "Horario creado exitosamente"
	LogScheduleServiceGetError       = "Error obteniendo horario"
	LogScheduleServiceGetOK          = "Horario obtenido exitosamente"
	LogScheduleServiceUpdateError    = "Error actualizando horario"
	LogScheduleServiceUpdateOK       = "Horario actualizado exitosamente"
	LogScheduleServiceDeleteError    = "Error eliminando horario"
	LogScheduleServiceDeleteOK       = "Horario eliminado exitosamente"
	LogScheduleServiceActivateError  = "Error activando/desactivando horario"
	LogScheduleServiceActivateOK     = "Horario activado/desactivado exitosamente"
)

// ============================================
// SCHEDULE INTERACTOR (HU30-35)
// ============================================
const (
	LogScheduleInteractorCreateStart    = "Iniciando proceso de registro de horario"
	LogScheduleInteractorTxError        = "Error iniciando transacción para horario"
	LogScheduleInteractorTxStarted      = "Transacción iniciada para horario"
	LogScheduleInteractorCreateError    = "Error registrando horario"
	LogScheduleInteractorCommitError    = "COMMIT FALLÓ - Horario"
	LogScheduleInteractorCreateComplete = "Horario registrado exitosamente"
	LogScheduleInteractorRollbackError  = "ROLLBACK FALLÓ - Horario"
	LogScheduleInteractorRollbackOK     = "Rollback ejecutado para horario"
	LogScheduleInteractorUpdateStart    = "Iniciando actualización de horario"
	LogScheduleInteractorUpdateComplete = LogScheduleServiceUpdateOK
	LogScheduleInteractorDeleteStart    = "Iniciando eliminación de horario"
	LogScheduleInteractorDeleteComplete = LogScheduleServiceDeleteOK
	LogScheduleInteractorGetError       = LogScheduleServiceGetError
	LogScheduleInteractorGetOK          = LogScheduleServiceGetOK
)

// ============================================
// SCHEDULE CONTROLLER (HU30-35, HU10)
// ============================================
const (
	LogScheduleControllerRequest             = "Solicitud recibida - Horario"
	LogScheduleControllerCreateRequest       = "Solicitud de creación de horario recibida"
	LogScheduleControllerCreateError         = "Error creando horario"
	LogScheduleControllerCreateOK            = "Horario creado exitosamente"
	LogScheduleControllerGetRequest          = "Solicitud de consulta de horario recibida"
	LogScheduleControllerGetError            = LogScheduleServiceGetError
	LogScheduleControllerGetOK               = LogScheduleServiceGetOK
	LogScheduleControllerUpdateRequest       = "Solicitud de actualización de horario recibida"
	LogScheduleControllerUpdateError         = "Error actualizando horario"
	LogScheduleControllerUpdateOK            = LogScheduleServiceUpdateOK
	LogScheduleControllerDeleteRequest       = "Solicitud de eliminación de horario recibida"
	LogScheduleControllerDeleteError         = "Error eliminando horario"
	LogScheduleControllerDeleteOK            = LogScheduleServiceDeleteOK
	LogScheduleControllerActivateReq         = "Solicitud de activación de horario recibida"
	LogScheduleControllerActivateError       = "Error activando horario"
	LogScheduleControllerActivateOK          = "Horario activado exitosamente"
	LogScheduleControllerDeactivateReq       = "Solicitud de desactivación de horario recibida"
	LogScheduleControllerDeactivateErr       = "Error desactivando horario"
	LogScheduleControllerDeactivateOK        = "Horario desactivado exitosamente"
	LogScheduleControllerGetDaysReq          = "Solicitud de catálogo de días recibida"
	LogScheduleControllerGetDaysOK           = "Catálogo de días enviado exitosamente"
	LogScheduleControllerBindError           = "Error parseando JSON de horario"
	LogScheduleControllerIDDecodeError       = "Error decodificando ID de horario"
	LogScheduleControllerDateParseError      = "Error parseando fecha de horario"
	LogScheduleControllerDateValidationError = "Error validando rango de fechas"
)

// ============================================
// SCHEDULE DETAIL SERVICE (HU6-9)
// ============================================
const (
	LogScheduleDetailServiceCreateStart      = "Iniciando creación de detalle horario"
	LogScheduleDetailServiceScheduleNotFound = "Horario base no encontrado para detalle"
	LogScheduleDetailServiceInvalidDay       = "Día de la semana inválido"
	LogScheduleDetailServiceInvalidTime      = "Formato de hora inválido"
	LogScheduleDetailServiceInvalidTimeRange = "Hora de cierre anterior a apertura"
	LogScheduleDetailServiceConflictCheck    = "Error verificando conflicts de horario"
	LogScheduleDetailServiceTimeConflict     = "Conflicto de horario detectado"
	LogScheduleDetailServiceSaveError        = "Error guardando detalle horario"
	LogScheduleDetailServiceCreateOK         = "Detalle horario creado exitosamente"
	LogScheduleDetailServiceGetError         = "Error obteniendo detalle horario"
	LogScheduleDetailServiceGetOK            = "Detalle horario obtenido exitosamente"
	LogScheduleDetailServiceListError        = "Error listando detalles horario"
	LogScheduleDetailServiceListOK           = "Detalles horario listados exitosamente"
	LogScheduleDetailServiceUpdateError      = "Error actualizando detalle horario"
	LogScheduleDetailServiceUpdateOK         = "Detalle horario actualizado exitosamente"
	LogScheduleDetailServiceDeleteError      = "Error eliminando detalle horario"
	LogScheduleDetailServiceDeleteOK         = "Detalle horario eliminado exitosamente"
)

// ============================================
// SCHEDULE DETAIL INTERACTOR (HU6-9)
// ============================================
const (
	LogScheduleDetailInteractorCreateStart    = "Iniciando creación de detalle horario"
	LogScheduleDetailInteractorBranchError    = "Error obteniendo sede para detalle horario"
	LogScheduleDetailInteractorOwnershipError = "Usuario no es dueño de la sede"
	LogScheduleDetailInteractorTxError        = "Error iniciando transacción de detalle horario"
	LogScheduleDetailInteractorCreateError    = "Error creando detalle horario"
	LogScheduleDetailInteractorCommitError    = "COMMIT FALLÓ - Detalle horario"
	LogScheduleDetailInteractorCreateOK       = LogScheduleDetailServiceCreateOK
	LogScheduleDetailInteractorListError      = LogScheduleDetailServiceListError
	LogScheduleDetailInteractorListOK         = LogScheduleDetailServiceListOK
	LogScheduleDetailInteractorRollbackError  = "ROLLBACK BD FALLÓ - ALERTA CRÍTICA (detalle horario)"
	LogScheduleDetailInteractorRollbackOK     = "Rollback BD ejecutado correctamente (detalle horario)"
)

// ============================================
// SCHEDULE DETAIL REPOSITORY (HU6-9)
// ============================================
const (
	LogScheduleDetailRepoPrepareError    = "Error preparando statement de detalle horario"
	LogScheduleDetailRepoSaveError       = "Error guardando detalle horario en BD"
	LogScheduleDetailRepoUpdateError     = "Error actualizando detalle horario en BD"
	LogScheduleDetailRepoDeleteError     = "Error eliminando detalle horario de BD"
	LogScheduleDetailRepoGetByIDError    = "Error obteniendo detalle horario por ID"
	LogScheduleDetailRepoGetBySchedError = "Error obteniendo detalles por horario"
	LogScheduleDetailRepoConflictCheck   = "Error verificando conflicts de horario"
	LogScheduleDetailRepoScanError       = "Error escaneando fila de detalle horario"
)

// ============================================
// SCHEDULE DETAIL CONTROLLER (HU6-9)
// ============================================
const (
	LogScheduleDetailControllerCreateRequest = "Solicitud de creación de detalle horario recibida"
	LogScheduleDetailControllerCreateError   = "Error creando detalle horario"
	LogScheduleDetailControllerCreateOK      = LogScheduleDetailServiceCreateOK
	LogScheduleDetailControllerListRequest   = "Solicitud de listado de detalles horario recibida"
	LogScheduleDetailControllerListError     = LogScheduleDetailServiceListError
	LogScheduleDetailControllerListOK        = LogScheduleDetailServiceListOK
	LogScheduleDetailControllerGetRequest    = "Solicitud de consulta de detalle horario recibida"
	LogScheduleDetailControllerGetError      = "Error obteniendo detalle horario"
	LogScheduleDetailControllerGetOK         = "Detalle horario obtenido exitosamente"
	LogScheduleDetailControllerUpdateRequest = "Solicitud de actualización de detalle horario recibida"
	LogScheduleDetailControllerUpdateError   = "Error actualizando detalle horario"
	LogScheduleDetailControllerUpdateOK      = "Detalle horario actualizado exitosamente"
	LogScheduleDetailControllerDeleteRequest = "Solicitud de eliminación de detalle horario recibida"
	LogScheduleDetailControllerDeleteError   = "Error eliminando detalle horario"
	LogScheduleDetailControllerDeleteOK      = "Detalle horario eliminado exitosamente"
	LogScheduleDetailControllerBindError     = "Error parseando JSON de detalle horario"
	LogScheduleDetailControllerIDDecodeError = "Error decodificando ID de detalle horario"
)

// ============================================
// MOTORCYCLE MODULE (HU43-47)
// ============================================
const (
	// Motorcycle Controller - GET (HU46)
	LogMotorcycleControllerGetRequest      = "Solicitud de consulta de motocicleta recibida"
	LogMotorcycleControllerIDDecodeError   = "Error decodificando ID de motocicleta"
	LogMotorcycleControllerGetByID         = "Buscando motocicleta por ID"
	LogMotorcycleControllerGetError        = "Error obteniendo motocicleta"
	LogMotorcycleControllerGetSuccess      = "Motocicleta consultada exitosamente"
	LogMotorcycleControllerOwnershipDenied = "Acceso denegado - la motocicleta pertenece a otro usuario"

	// Motorcycle Controller - POST (HU43)
	LogMotorcycleControllerRegRequest  = "Solicitud de registro de motocicleta recibida"
	LogMotorcycleControllerRegBody     = "Cuerpo de registro de motocicleta parseado"
	LogMotorcycleControllerAuthError   = "Error de autenticación en registro de motocicleta"
	LogMotorcycleControllerBindError   = "Error parseando cuerpo de solicitud de motocicleta"
	LogMotorcycleControllerRefDecError = "Error decodificando ID de referencia de motocicleta"
	LogMotorcycleControllerRegError    = "Error registrando motocicleta"
	LogMotorcycleControllerIDEncError  = "Error codificando ID de motocicleta"
	LogMotorcycleControllerRegSuccess  = "Motocicleta registrada exitosamente"

	// Motorcycle Controller - LIST (HU47)
	LogMotorcycleControllerListRequest = "Solicitud de listado de motocicletas recibida"
	LogMotorcycleControllerListError   = "Error listando motocicletas"
	LogMotorcycleControllerListSuccess = "Motocicletas listadas exitosamente"

	// Motorcycle Interactor (HU43-47)
	LogMotorcycleInteractorRegStart        = "Registro de motocicleta iniciado"
	LogMotorcycleInteractorValidateRef     = "Validando existencia de referencia de motocicleta"
	LogMotorcycleInteractorRefError        = "Error validando referencia de motocicleta"
	LogMotorcycleInteractorRefNotFound     = "Referencia de motocicleta no encontrada"
	LogMotorcycleInteractorRefRequired     = "El reference_id de motocicleta es requerido"
	LogMotorcycleInteractorCheckPlate      = "Verificando unicidad de placa"
	LogMotorcycleInteractorCheckPlateErr   = "Error verificando unicidad de placa"
	LogMotorcycleInteractorDupPlate        = "Placa duplicada detectada"
	LogMotorcycleInteractorIDGenerated     = "ID de motocicleta generado"
	LogMotorcycleInteractorBeginTxError    = "Error iniciando transacción para motocicleta"
	LogMotorcycleInteractorSaveError       = "Error guardando motocicleta"
	LogMotorcycleInteractorCommitError     = "Error confirmando transacción de motocicleta"
	LogMotorcycleInteractorRollbackError   = "ROLLBACK BD FALLÓ - ALERTA CRÍTICA (motocicleta)"
	LogMotorcycleInteractorRollbackOK      = "Rollback BD ejecutado correctamente (motocicleta)"
	LogMotorcycleInteractorRegSuccess      = "Motocicleta registrada exitosamente"
	LogMotorcycleInteractorGetStart        = "Consulta de motocicleta por ID iniciada"
	LogMotorcycleInteractorGetError        = "Error obteniendo motocicleta por ID"
	LogMotorcycleInteractorGetSuccess      = "Motocicleta obtenida exitosamente"
	LogMotorcycleInteractorGetOwnerStart   = "Consulta de motocicletas por propietario iniciada"
	LogMotorcycleInteractorGetOwnerError   = "Error obteniendo motocicletas por propietario"
	LogMotorcycleInteractorGetOwnerSuccess = "Motocicletas por propietario obtenidas exitosamente"
	LogMotorcycleInteractorGetPlateStart   = "Consulta de motocicleta por placa iniciada"
	LogMotorcycleInteractorGetPlateError   = "Error obteniendo motocicleta por placa"
	LogMotorcycleInteractorGetPlateSuccess = "Motocicleta por placa obtenida exitosamente"

	// Motorcycle Controller - GET BY PLATE (HU47)
	LogMotorcycleControllerPlateRequest = "Solicitud de búsqueda por placa recibida"
	LogMotorcycleControllerPlateError   = "Error buscando motocicleta por placa"
	LogMotorcycleControllerPlateSuccess = "Motocicleta por placa obtenida exitosamente"

	// Motorcycle Repository (HU43-47)
	LogMotorcycleRepoGetByID         = "Repositorio obteniendo motocicleta por ID"
	LogMotorcycleRepoGetByIDError    = "Error de repositorio obteniendo motocicleta por ID"
	LogMotorcycleRepoGetByOwner      = "Repositorio obteniendo motocicletas por propietario"
	LogMotorcycleRepoGetByOwnerError = "Error de repositorio obteniendo motocicletas por propietario"
	LogMotorcycleRepoGetByOwnerScan  = "Error de repositorio escaneando fila de motocicleta"
	LogMotorcycleRepoGetByOwnerIter  = "Error de repositorio iterando filas de motocicletas"
	LogMotorcycleRepoGetByPlate      = "Repositorio obteniendo motocicleta por placa"
	LogMotorcycleRepoGetByPlateError = "Error de repositorio obteniendo motocicleta por placa"
	LogMotorcycleRepoInvalidTx       = "Tipo de transacción inválido en repositorio"
	LogMotorcycleRepoSave            = "Repositorio guardando motocicleta"
	LogMotorcycleRepoSaveError       = "Error de repositorio guardando motocicleta"
	LogMotorcycleRepoSaveSuccess     = "Motocicleta guardada en repositorio"
	LogMotorcycleRepoValidateRef     = "Repositorio validando referencia"
	LogMotorcycleRepoValidateRefErr  = "Error de repositorio validando referencia"
	LogMotorcycleRepoCheckPlate      = "Repositorio verificando placa"
	LogMotorcycleRepoCheckPlateErr   = "Error de repositorio verificando placa"
	LogMotorcycleRepoUpdate          = "Repositorio actualizando motocicleta"
	LogMotorcycleRepoUpdateError     = "Error de repositorio actualizando motocicleta"
	LogMotorcycleRepoUpdateSuccess   = "Motocicleta actualizada en repositorio"

	// Motorcycle Interactor - UPDATE (HU44)
	LogMotorcycleInteractorUpdateStart   = "Actualización de motocicleta iniciada"
	LogMotorcycleInteractorUpdateError   = "Error actualizando motocicleta"
	LogMotorcycleInteractorUpdateSuccess = "Motocicleta actualizada exitosamente"

	// Motorcycle Controller - UPDATE (HU44)
	LogMotorcycleControllerUpdateRequest  = "Solicitud de actualización de motocicleta recibida"
	LogMotorcycleControllerUpdateError    = "Error actualizando motocicleta"
	LogMotorcycleControllerUpdateSuccess  = "Motocicleta actualizada exitosamente"
	LogMotorcycleControllerUpdateDebug    = "Actualizando motocicleta"
	LogMotorcycleControllerOwnershipDebug = "Verificación de propiedad (debug)"
	LogMotorcycleControllerNoAuthUser     = "Usuario no autenticado en contexto"
	LogMotorcycleControllerPlateDebug     = "Buscando motocicleta por placa"
	LogMotorcycleControllerMissingPlate   = "Parámetro de placa faltante"

	// Motorcycle Repository - DELETE (HU45)
	LogMotorcycleRepoDelete        = "Repositorio eliminando motocicleta"
	LogMotorcycleRepoDeleteError   = "Error de repositorio eliminando motocicleta"
	LogMotorcycleRepoDeleteSuccess = "Motocicleta eliminada en repositorio"

	// Motorcycle Interactor - DELETE (HU45)
	LogMotorcycleInteractorDeleteStart   = "Eliminación de motocicleta iniciada"
	LogMotorcycleInteractorDeleteError   = "Error eliminando motocicleta"
	LogMotorcycleInteractorDeleteSuccess = "Motocicleta eliminada exitosamente"

	// Motorcycle Controller - DELETE (HU45)
	LogMotorcycleControllerDeleteRequest = "Solicitud de eliminación de motocicleta recibida"
	LogMotorcycleControllerDeleteError   = "Error eliminando motocicleta"
	LogMotorcycleControllerDeleteSuccess = "Motocicleta eliminada exitosamente"

	// Dependency initialization
	LogDepMotorcycleRepoInitOK       = "Repositorio de motocicletas inicializado"
	LogDepMotorcycleRepoInitErr      = "Error inicializando repositorio de motocicletas"
	LogDepMotorcycleInteractorInitOK = "Interactor de motocicletas inicializado"

	// Evidence dependency initialization (HU16-19)
	LogDepEvidenceRepoInitOK       = "Repositorio de evidencias inicializado"
	LogDepEvidenceRepoInitErr      = "Error inicializando repositorio de evidencias"
	LogDepEvidenceInteractorInitOK = "Interactor de evidencias inicializado"

	// Motorcycle Reference Catalog (HU50)
	LogMotorcycleRepoGetAllRefQuery     = "Repositorio obteniendo todas las referencias de motocicletas"
	LogMotorcycleRepoGetAllRefScanError = "Error de repositorio escaneando fila de referencia"
	LogMotorcycleRepoGetAllRefIterError = "Error de repositorio iterando filas de referencias"
	LogMotorcycleInteractorGetRefsStart = "Consulta de referencias de motocicletas iniciada"
	LogMotorcycleInteractorGetRefsError = "Error obteniendo referencias de motocicletas"
	LogMotorcycleControllerRefsRequest  = "Solicitud de referencias de motocicletas recibida"
	LogMotorcycleControllerRefsError    = "Error obteniendo referencias de motocicletas"
	LogMotorcycleControllerRefsSuccess  = "Referencias de motocicletas obtenidas exitosamente"

	// Motorcycle Brand Lines (HU40)
	LogMotorcycleRepoBrandLinesQuery         = "Repositorio obteniendo referencias por marca"
	LogMotorcycleRepoBrandLinesError         = "Error de repositorio obteniendo referencias por marca"
	LogMotorcycleRepoBrandLinesScanError     = "Error de repositorio escaneando fila de línea de marca"
	LogMotorcycleRepoBrandLinesIterError     = "Error de repositorio iterando filas de líneas de marca"
	LogMotorcycleInteractorBrandLinesStart   = "Consulta de referencias por marca iniciada"
	LogMotorcycleInteractorBrandLinesError   = "Error obteniendo referencias por marca"
	LogMotorcycleInteractorBrandLinesSuccess = "Referencias por marca obtenidas exitosamente"
	LogMotorcycleControllerBrandLinesRequest = "Solicitud de líneas de marca recibida"
	LogMotorcycleControllerBrandLinesError   = "Error obteniendo líneas de marca"
	LogMotorcycleControllerBrandLinesSuccess = "Líneas de marca obtenidas exitosamente"

	// Motorcycle Category Catalog (HU41)
	LogMotorcycleCatRepoQuery          = "Repositorio obteniendo categorías distintas de motocicletas"
	LogMotorcycleCatRepoError          = "Error de repositorio obteniendo categorías de motocicletas"
	LogMotorcycleCatRepoScanError      = "Error de repositorio escaneando fila de categoría"
	LogMotorcycleCatRepoIterError      = "Error de repositorio iterando filas de categorías"
	LogMotorcycleCatLinesRepoQuery     = "Repositorio obteniendo líneas por categoría"
	LogMotorcycleCatLinesRepoError     = "Error de repositorio obteniendo líneas por categoría"
	LogMotorcycleCatLinesRepoScanError = "Error de repositorio escaneando fila de línea de categoría"
	LogMotorcycleCatLinesRepoIterError = "Error de repositorio iterando filas de líneas de categoría"

	LogMotorcycleCatInteractorStart        = "Consulta de categorías de motocicletas iniciada"
	LogMotorcycleCatInteractorError        = "Error obteniendo categorías de motocicletas"
	LogMotorcycleCatInteractorSuccess      = "Categorías de motocicletas obtenidas exitosamente"
	LogMotorcycleCatLinesInteractorStart   = "Consulta de líneas por categoría iniciada"
	LogMotorcycleCatLinesInteractorError   = "Error obteniendo líneas por categoría"
	LogMotorcycleCatLinesInteractorSuccess = "Líneas por categoría obtenidas exitosamente"

	LogMotorcycleCatControllerRequest      = "Solicitud de categorías de motocicletas recibida"
	LogMotorcycleCatControllerError        = "Error obteniendo categorías de motocicletas"
	LogMotorcycleCatControllerSuccess      = "Categorías de motocicletas obtenidas exitosamente"
	LogMotorcycleCatLinesControllerRequest = "Solicitud de líneas por categoría recibida"
	LogMotorcycleCatLinesControllerError   = "Error obteniendo líneas por categoría"
	LogMotorcycleCatLinesControllerSuccess = "Líneas por categoría obtenidas exitosamente"

	// Engine Displacement Ranges (HU49)
	LogMotorcycleDispRepoQuery     = "Repositorio obteniendo rangos de cilindraje distintos"
	LogMotorcycleDispRepoError     = "Error de repositorio obteniendo rangos de cilindraje"
	LogMotorcycleDispRepoScanError = "Error de repositorio escaneando fila de cilindraje"
	LogMotorcycleDispRepoIterError = "Error de repositorio iterando filas de cilindraje"

	LogMotorcycleDispInteractorStart   = "Consulta de rangos de cilindraje iniciada"
	LogMotorcycleDispInteractorError   = "Error obteniendo rangos de cilindraje"
	LogMotorcycleDispInteractorSuccess = "Rangos de cilindraje obtenidos exitosamente"

	LogMotorcycleDispControllerRequest = "Solicitud de rangos de cilindraje recibida"
	LogMotorcycleDispControllerError   = "Error obteniendo rangos de cilindraje"
	LogMotorcycleDispControllerSuccess = "Rangos de cilindraje obtenidos exitosamente"

	// Rating Ranges (HU48)
	LogRatingRangeInteractorStart   = "Consulta de rangos de calificación iniciada"
	LogRatingRangeInteractorError   = "Error obteniendo rangos de calificación"
	LogRatingRangeInteractorSuccess = "Rangos de calificación obtenidos exitosamente"

	LogRatingRangeControllerRequest = "Solicitud de rangos de calificación recibida"
	LogRatingRangeControllerError   = "Error obteniendo rangos de calificación"
	LogRatingRangeControllerSuccess = "Rangos de calificación obtenidos exitosamente"
)

// ============================================
// MOTORCYCLE SERVICE (HU43-47, HU50, HU40)
// ============================================
const (
	LogMotorcycleServiceRefError         = "Error validando referencia de motocicleta en servicio"
	LogMotorcycleServiceCheckPlateErr    = "Error verificando unicidad de placa en servicio"
	LogMotorcycleServiceSaveError        = "Error guardando motocicleta en servicio"
	LogMotorcycleServiceUpdateError      = "Error actualizando motocicleta en servicio"
	LogMotorcycleServiceDeleteStart      = "Eliminación de motocicleta en servicio"
	LogMotorcycleServiceDeleteError      = "Error eliminando motocicleta en servicio"
	LogMotorcycleServiceStorageDeleteErr = "Error eliminando archivo de almacenamiento (continuando)"

	// Dependency initialization
	LogDepMotorcycleServiceInitOK = "Servicio de motocicletas inicializado"
)

// ============================================
// DIAGNOSTIC PERMISSION SERVICE
// ============================================
const (
	LogDiagPermServiceSaveError   = "Error guardando permiso de diagnóstico en servicio"
	LogDiagPermServiceDeleteError = "Error eliminando permiso de diagnóstico en servicio"
)

// ============================================
// MOTORCYCLE EVIDENCE INTERACTOR (HU16-19)
// ============================================
const (
	// Evidence Create (HU16)
	LogEvidenceInteractorCreateStart     = "Creación de evidencia iniciada"
	LogEvidenceInteractorMotorcycleError = "Error obteniendo motocicleta para evidencia"
	LogEvidenceInteractorOwnerError      = "El usuario no es propietario de esta motocicleta"
	LogEvidenceInteractorURLInvalid      = "URL de Firebase Storage inválida"
	LogEvidenceInteractorAngleInvalid    = "Ángulo de evidencia proporcionado inválido"
	LogEvidenceInteractorCountError      = "Error contando evidencias de la motocicleta"
	LogEvidenceInteractorLimitExceeded   = "Límite máximo de evidencias excedido"
	LogEvidenceInteractorIDGenerated     = "ID de evidencia generado"
	LogEvidenceInteractorBeginTxError    = "Error iniciando transacción para evidencia"
	LogEvidenceInteractorSaveError       = "Error guardando evidencia"
	LogEvidenceInteractorCommitError     = "Error confirmando transacción de evidencia"
	LogEvidenceInteractorRollbackError   = "ROLLBACK BD FALLÓ - ALERTA CRÍTICA (evidencia)"
	LogEvidenceInteractorRollbackOK      = "Rollback BD ejecutado correctamente (evidencia)"
	LogEvidenceInteractorCreateSuccess   = "Evidencia creada exitosamente"

	// Evidence Get (HU18)
	LogEvidenceInteractorGetStart   = "Consulta de evidencia por ID iniciada"
	LogEvidenceInteractorGetError   = "Error obteniendo evidencia por ID"
	LogEvidenceInteractorGetSuccess = "Evidencia obtenida exitosamente"

	// Evidence List (HU18)
	LogEvidenceInteractorListStart   = "Listado de evidencias por motocicleta iniciado"
	LogEvidenceInteractorListError   = "Error listando evidencias por motocicleta"
	LogEvidenceInteractorListSuccess = "Evidencias listadas exitosamente"

	// Evidence Update (HU17)
	LogEvidenceInteractorUpdateStart   = "Actualización de evidencia iniciada"
	LogEvidenceInteractorUpdateError   = "Error actualizando evidencia"
	LogEvidenceInteractorUpdateSuccess = "Evidencia actualizada exitosamente"

	// Evidence Delete (HU19)
	LogEvidenceInteractorDeleteStart   = "Eliminación de evidencia iniciada"
	LogEvidenceInteractorDeleteError   = "Error eliminando evidencia"
	LogEvidenceInteractorDeleteSuccess = "Evidencia eliminada exitosamente"

	// LookupEvidence (representative plate lookup - no ownership check)
	LogEvidenceInteractorLookupStart   = "Consulta de evidencias para búsqueda por placa iniciada"
	LogEvidenceInteractorLookupError   = "Error consultando evidencias para búsqueda por placa"
	LogEvidenceInteractorLookupSuccess = "Evidencias para búsqueda por placa obtenidas exitosamente"
)

// ============================================
// MOTORCYCLE EVIDENCE CONTROLLER (HU16-19)
// ============================================
const (
	LogEvidenceControllerCreateRequest = "Solicitud de creación de evidencia recibida"
	LogEvidenceControllerCreateError   = "Error creando evidencia"
	LogEvidenceControllerCreateSuccess = "Evidencia creada exitosamente"
	LogEvidenceControllerListRequest   = "Solicitud de listado de evidencias recibida"
	LogEvidenceControllerListError     = "Error listando evidencias"
	LogEvidenceControllerListSuccess   = "Evidencias listadas exitosamente"
	LogEvidenceControllerGetRequest    = "Solicitud de consulta de evidencia recibida"
	LogEvidenceControllerGetError      = "Error obteniendo evidencia"
	LogEvidenceControllerGetSuccess    = "Evidencia obtenida exitosamente"
	LogEvidenceControllerUpdateRequest = "Solicitud de actualización de evidencia recibida"
	LogEvidenceControllerUpdateError   = LogEvidenceInteractorUpdateError
	LogEvidenceControllerUpdateSuccess = "Evidencia actualizada exitosamente"
	LogEvidenceControllerDeleteRequest = "Solicitud de eliminación de evidencia recibida"
	LogEvidenceControllerDeleteError   = LogEvidenceInteractorDeleteError
	LogEvidenceControllerDeleteSuccess = "Evidencia eliminada exitosamente"
)

// ============================================
// MOTORCYCLE EVIDENCE REPOSITORY (HU16-19)
// ============================================
const (
	LogDatabaseError = "Error de base de datos"
	// Evidence Repository operations
	LogEvidenceRepoSaveError          = "Error guardando evidencia en BD"
	LogEvidenceRepoGetByIDError       = "Error obteniendo evidencia por ID"
	LogEvidenceRepoListByMotoError    = "Error listando evidencias por motocicleta"
	LogEvidenceRepoScanError          = "Error escaneando fila de evidencia"
	LogEvidenceRepoCountError         = "Error contando evidencias"
	LogEvidenceRepoUpdateError        = LogEvidenceInteractorUpdateError
	LogEvidenceRepoDeleteError        = LogEvidenceInteractorDeleteError
	LogEvidenceRepoPrepareInsertError = "Error preparando statement de inserción de evidencia"
	LogEvidenceRepoPrepareUpdateError = "Error preparando statement de actualización de evidencia"
	LogEvidenceRepoPrepareDeleteError = "Error preparando statement de eliminación de evidencia"
	LogEvidenceRepoPrepareGetIDError  = "Error preparando statement de consulta por ID de evidencia"
	LogEvidenceRepoPrepareGetMotoErr  = "Error preparando statement de consulta por motocicleta"
	LogEvidenceRepoPrepareCountError  = "Error preparando statement de conteo de evidencia"
)

// ============================================
// FRANCHISE OPERATIONS (Centralized)
// ============================================
const (
	LogFranchiseNameExists       = "Nombre de franquicia ya existe"
	LogFranchiseCreated          = "Franquicia creada"
	LogFranchiseUpdated          = "Franquicia actualizada"
	LogFranchiseDeleted          = "Franquicia eliminada"
	LogFranchiseBranchAdded      = "Sede agregada a franquicia"
	LogFranchiseBranchRemoved    = "Sede removida de franquicia"
	LogFranchiseCannotRemoveLast = "No se puede remover la última sede de la franquicia"
)

// ============================================
// PERSON OPERATIONS (Centralized)
// ============================================
const (
	LogPersonRefreshTokenStart     = "Inicio de refresh token"
	LogPersonRefreshTokenError     = "Error en refresh token"
	LogPersonGetPublicContactStart = "Obteniendo información de contacto público"
	LogPersonGetPublicContactError = "Error obteniendo persona para contacto público"
	LogPersonDeleteKeycloakStart   = "Eliminando usuario de Keycloak"
	LogPersonDeleteKeycloakError   = "Error eliminando usuario de Keycloak"
	LogPersonDeleteDBStart         = "Eliminando persona de base de datos"
	LogPersonDeleteDBError         = "Error eliminando persona de base de datos"
	LogPersonAuthNotFoundInContext = "Usuario autenticado no encontrado en contexto"
	LogPersonRetrievingProfile     = "Obteniendo perfil de usuario autenticado"
	LogPersonIDEncodeError         = "Error codificando ID de usuario"
	LogPersonDeleteSelfRequest     = "Solicitud de eliminación de cuenta propia recibida"
	LogPersonCheckBranchesError    = "Error verificando sedes del usuario"
	LogPersonHasActiveBranches     = "Usuario tiene sedes activas, no se puede eliminar"
	LogPersonDeleteKeycloakFailed  = "Error eliminando de Keycloak"
	LogPersonDeleteDBFailed        = "Error eliminando de base de datos"
	LogPersonMissingIDInURL        = "ID de persona faltante en URL"
	LogPersonIDDecodeError         = "Error decodificando ID de persona"
	LogPersonGetError              = "Error obteniendo persona"
)

// ============================================
// BRANCH OPERATIONS (Centralized)
// ============================================
const (
	LogBranchGeocodingFailed     = "Paso de geocodificación falló"
	LogBranchGeocodingGenerated  = "Coordenadas de geocodificación generadas"
	LogBranchRefetchFailed       = "Error re-obteniendo sede"
	LogBranchDeleteProcessing    = "Procesando eliminación de sede"
	LogBranchNearbyMissingLat    = "Latitud faltante en búsqueda de sedes cercanas"
	LogBranchNearbyInvalidLat    = "Latitud inválida en búsqueda de sedes cercanas"
	LogBranchNearbyMissingLng    = "Longitud faltante en búsqueda de sedes cercanas"
	LogBranchNearbyInvalidLng    = "Longitud inválida en búsqueda de sedes cercanas"
	LogBranchNearbyInvalidRadius = "Radio inválido en búsqueda de sedes cercanas"
	LogBranchNearbyInvalidType   = "Tipo de establecimiento inválido"
	LogBranchNearbySearch        = "Búsqueda de sedes cercanas"
	LogBranchNearbyError         = "Error en búsqueda de sedes cercanas"
)

// ============================================
// GEOCODING OPERATIONS (Centralized)
// ============================================
const (
	LogGeocodingTestRequest    = "Solicitud de prueba de geocodificación recibida"
	LogGeocodingTestInvalid    = "Solicitud de prueba de geocodificación inválida"
	LogGeocodingTestFailed     = "Prueba de geocodificación falló"
	LogGeocodingPrimaryQuota   = "Cuota del servicio primario de geocodificación excedida"
	LogGeocodingFallbackFailed = "Servicio de respaldo de geocodificación también falló"
	LogGeocodingFallbackOK     = "Geocodificación exitosa con servicio de respaldo"
)

// ============================================
// SCHEDULE DETAIL DEBUG (Centralized)
// ============================================
const (
	LogScheduleDetailDebugQueryParams   = "Parámetros de consulta de verificación de conflicto de excepción"
	LogScheduleDetailDebugQueryResult   = "Resultado de consulta de verificación de conflicto de excepción"
	LogScheduleDetailDebugGetError      = "Error obteniendo excepciones para verificación de conflicto"
	LogScheduleDetailDebugExplicitCheck = "Verificación explícita de conflicto de excepciones"
	LogScheduleDetailDebugCheckOverlap  = "Verificación de solapamiento de excepciones"
)

// ============================================
// PERSON CONTROLLER (Centralized)
// ============================================
const (
	LogPersonControllerRefreshRequest  = "Solicitud de refresh token recibida"
	LogPersonControllerRefreshError    = "Error en refresh token"
	LogPersonControllerGetContactStart = "Obteniendo información de contacto"
)

// ============================================
// MOTORCYCLE CONTROLLER DEBUG
// ============================================
const (
	LogMotorcycleControllerBrandIDDecoded = "ID de marca decodificado"
)

// ============================================
// DIAGNOSTIC REPOSITORY (HU11-14)
// ============================================
const (
	LogDiagnosticRepoSaveError           = "Error guardando diagnóstico en BD"
	LogDiagnosticRepoSaveEvidenceError   = "Error guardando evidencia de diagnóstico en BD"
	LogDiagnosticRepoGetByIDError        = "Error obteniendo diagnóstico por ID"
	LogDiagnosticRepoListByMotoError     = "Error listando diagnósticos por motocicleta"
	LogDiagnosticRepoListEvidenceError   = "Error listando evidencias de diagnóstico"
	LogDiagnosticRepoScanError           = "Error escaneando fila de diagnóstico"
	LogDiagnosticRepoUpdateError         = "Error actualizando diagnóstico en BD"
	LogDiagnosticRepoDeleteError         = "Error eliminando diagnóstico de BD"
	LogDiagnosticRepoPrepareInsertError  = "Error preparando statement de inserción de diagnóstico"
	LogDiagnosticRepoPrepareUpdateError  = "Error preparando statement de actualización de diagnóstico"
	LogDiagnosticRepoPrepareDeleteError  = "Error preparando statement de eliminación de diagnóstico"
	LogDiagnosticRepoPrepareGetIDError   = "Error preparando statement de consulta por ID de diagnóstico"
	LogDiagnosticRepoPrepareGetMotoError = "Error preparando statement de consulta por motocicleta de diagnóstico"
	LogDiagnosticRepoPrepareEvidInsError = "Error preparando statement de inserción de evidencia de diagnóstico"
	LogDiagnosticRepoPrepareEvidGetError = "Error preparando statement de consulta de evidencia de diagnóstico"

	// UPSERT support (diagnostic deduplication by motorcycle+branch)
	LogDiagnosticRepoPrepareGetMotoBranchError = "Error preparando statement de consulta por moto+sede de diagnóstico"
	LogDiagnosticRepoPrepareEvidDelError       = "Error preparando statement de eliminación de evidencias de diagnóstico"
	LogDiagnosticRepoGetByMotoBranchError      = "Error obteniendo diagnóstico por moto+sede"
	LogDiagnosticRepoDeleteEvidenceError       = "Error eliminando evidencias de diagnóstico"
)

// ============================================
// DIAGNOSTIC INTERACTOR (HU11-14)
// ============================================
const (
	// Diagnostic Create (HU11)
	LogDiagnosticInteractorCreateStart   = "Creación de diagnóstico iniciada"
	LogDiagnosticInteractorMotoError     = "Error obteniendo motocicleta para diagnóstico"
	LogDiagnosticInteractorOwnerError    = "El usuario no es propietario de esta motocicleta"
	LogDiagnosticInteractorBranchError   = "Error validando sede para diagnóstico"
	LogDiagnosticInteractorIDGenerated   = "ID de diagnóstico generado"
	LogDiagnosticInteractorBeginTxError  = "Error iniciando transacción para diagnóstico"
	LogDiagnosticInteractorSaveError     = "Error guardando diagnóstico"
	LogDiagnosticInteractorSaveEvidError = "Error guardando evidencia de diagnóstico"
	LogDiagnosticInteractorCommitError   = "Error confirmando transacción de diagnóstico"
	LogDiagnosticInteractorRollbackError = "ROLLBACK BD FALLÓ - ALERTA CRÍTICA (diagnóstico)"
	LogDiagnosticInteractorRollbackOK    = "Rollback BD ejecutado correctamente (diagnóstico)"
	LogDiagnosticInteractorCreateSuccess = "Diagnóstico creado exitosamente"

	// Diagnostic UPSERT (same moto+branch → update)
	LogDiagnosticInteractorExistingFound    = "Diagnóstico existente encontrado para moto+sede, actualizando"
	LogDiagnosticInteractorEvidCleanupError = "Error limpiando evidencias del diagnóstico existente"
	LogDiagnosticInteractorUpsertUpdateErr  = "Error actualizando diagnóstico existente"
	LogDiagnosticInteractorUpsertSuccess    = "Diagnóstico actualizado por UPSERT"

	// Diagnostic Get (HU14)
	LogDiagnosticInteractorGetStart   = "Consulta de diagnóstico por ID iniciada"
	LogDiagnosticInteractorGetError   = "Error obteniendo diagnóstico por ID"
	LogDiagnosticInteractorGetSuccess = "Diagnóstico obtenido exitosamente"

	// Diagnostic List
	LogDiagnosticInteractorListStart   = "Listado de diagnósticos por motocicleta iniciado"
	LogDiagnosticInteractorListError   = "Error listando diagnósticos por motocicleta"
	LogDiagnosticInteractorListSuccess = "Diagnósticos listados exitosamente"

	// Diagnostic Update (HU12)
	LogDiagnosticInteractorUpdateStart   = "Actualización de diagnóstico iniciada"
	LogDiagnosticInteractorUpdateError   = "Error actualizando diagnóstico"
	LogDiagnosticInteractorUpdateSuccess = "Diagnóstico actualizado exitosamente"

	// Diagnostic Delete (HU13)
	LogDiagnosticInteractorDeleteStart   = "Eliminación de diagnóstico iniciada"
	LogDiagnosticInteractorDeleteError   = "Error eliminando diagnóstico"
	LogDiagnosticInteractorDeleteSuccess = "Diagnóstico eliminado exitosamente"

	// Set Solution (representative)
	LogDiagnosticInteractorSetSolutionStart   = "Asignación de solución a diagnóstico iniciada"
	LogDiagnosticInteractorSetSolutionError   = "Error asignando solución a diagnóstico"
	LogDiagnosticInteractorSetSolutionSuccess = "Solución asignada a diagnóstico exitosamente"
)

// ============================================
// DIAGNOSTIC CONTROLLER (HU11-14)
// ============================================
const (
	LogDiagnosticControllerCreateRequest = "Solicitud de creación de diagnóstico recibida"
	LogDiagnosticControllerCreateError   = "Error creando diagnóstico"
	LogDiagnosticControllerCreateSuccess = "Diagnóstico creado exitosamente"
	LogDiagnosticControllerListRequest   = "Solicitud de listado de diagnósticos recibida"
	LogDiagnosticControllerListError     = "Error listando diagnósticos"
	LogDiagnosticControllerListSuccess   = "Diagnósticos listados exitosamente"
	LogDiagnosticControllerGetRequest    = "Solicitud de consulta de diagnóstico recibida"
	LogDiagnosticControllerGetError      = "Error obteniendo diagnóstico"
	LogDiagnosticControllerGetSuccess    = "Diagnóstico obtenido exitosamente"
	LogDiagnosticControllerUpdateRequest = "Solicitud de actualización de diagnóstico recibida"
	LogDiagnosticControllerUpdateError   = "Error actualizando diagnóstico"
	LogDiagnosticControllerUpdateSuccess = "Diagnóstico actualizado exitosamente"
	LogDiagnosticControllerDeleteRequest = "Solicitud de eliminación de diagnóstico recibida"
	LogDiagnosticControllerDeleteError   = "Error eliminando diagnóstico"
	LogDiagnosticControllerDeleteSuccess = "Diagnóstico eliminado exitosamente"

	// Set Solution (representative)
	LogDiagnosticControllerSetSolutionRequest = "Solicitud de asignación de solución recibida"
	LogDiagnosticControllerSetSolutionError   = "Error asignando solución a diagnóstico"
	LogDiagnosticControllerSetSolutionSuccess = "Solución asignada a diagnóstico exitosamente"

	// Dependency initialization
	LogDepDiagnosticRepoInitOK       = "Repositorio de diagnósticos inicializado"
	LogDepDiagnosticRepoInitErr      = "Error inicializando repositorio de diagnósticos"
	LogDepDiagnosticInteractorInitOK = "Interactor de diagnósticos inicializado"
)

// ============================================
// DIAGNOSTIC PERMISSION REPOSITORY
// ============================================
const (
	LogDiagPermRepoSaveError          = "Error guardando permiso de diagnóstico en BD"
	LogDiagPermRepoDeleteError        = "Error eliminando permiso de diagnóstico de BD"
	LogDiagPermRepoGetError           = "Error obteniendo permiso de diagnóstico"
	LogDiagPermRepoListError          = "Error listando permisos de diagnóstico"
	LogDiagPermRepoScanError          = "Error escaneando fila de permiso de diagnóstico"
	LogDiagPermRepoPrepareSaveError   = "Error preparando statement de inserción de permiso de diagnóstico"
	LogDiagPermRepoPrepareDeleteError = "Error preparando statement de eliminación de permiso de diagnóstico"
	LogDiagPermRepoPrepareGetError    = "Error preparando statement de consulta de permiso de diagnóstico"
	LogDiagPermRepoPrepareListError   = "Error preparando statement de listado de permiso de diagnóstico"
)

// ============================================
// DIAGNOSTIC PERMISSION INTERACTOR
// ============================================
const (
	LogDiagPermInteractorGrantStart    = "Concesión de permiso de diagnóstico iniciada"
	LogDiagPermInteractorMotoError     = "Error obteniendo motocicleta para permiso"
	LogDiagPermInteractorOwnerError    = "El usuario no es propietario de la motocicleta"
	LogDiagPermInteractorBranchError   = "Error validando sede para permiso"
	LogDiagPermInteractorBeginTxError  = "Error iniciando transacción para permiso"
	LogDiagPermInteractorSaveError     = "Error guardando permiso de diagnóstico"
	LogDiagPermInteractorCommitError   = "Error confirmando transacción de permiso"
	LogDiagPermInteractorRollbackError = "ROLLBACK BD FALLÓ - ALERTA CRÍTICA (permiso)"
	LogDiagPermInteractorRollbackOK    = "Rollback BD ejecutado correctamente (permiso)"
	LogDiagPermInteractorGrantSuccess  = "Permiso de diagnóstico concedido exitosamente"
	LogDiagPermInteractorRevokeStart   = "Revocación de permiso de diagnóstico iniciada"
	LogDiagPermInteractorDeleteError   = "Error revocando permiso de diagnóstico"
	LogDiagPermInteractorRevokeSuccess = "Permiso de diagnóstico revocado exitosamente"
	LogDiagPermInteractorListStart     = "Listado de permisos de diagnóstico iniciado"
	LogDiagPermInteractorListError     = LogDiagPermRepoListError
	LogDiagPermInteractorListSuccess   = "Permisos de diagnóstico listados exitosamente"

	// LookupPermissions (representative plate lookup - no ownership check)
	LogDiagPermInteractorLookupStart   = "Consulta de permisos para búsqueda por placa iniciada"
	LogDiagPermInteractorLookupError   = "Error consultando permisos para búsqueda por placa"
	LogDiagPermInteractorLookupSuccess = "Permisos para búsqueda por placa obtenidos exitosamente"
)

// ============================================
// DIAGNOSTIC PERMISSION CONTROLLER
// ============================================
const (
	LogDiagPermControllerGrantRequest  = "Solicitud de concesión de permiso de diagnóstico recibida"
	LogDiagPermControllerGrantError    = "Error concediendo permiso de diagnóstico"
	LogDiagPermControllerGrantSuccess  = "Permiso de diagnóstico concedido exitosamente"
	LogDiagPermControllerRevokeRequest = "Solicitud de revocación de permiso de diagnóstico recibida"
	LogDiagPermControllerRevokeError   = "Error revocando permiso de diagnóstico"
	LogDiagPermControllerRevokeSuccess = "Permiso de diagnóstico revocado exitosamente"
	LogDiagPermControllerListRequest   = "Solicitud de listado de permisos de diagnóstico recibida"
	LogDiagPermControllerListError     = LogDiagPermRepoListError
	LogDiagPermControllerListSuccess   = "Permisos de diagnóstico listados exitosamente"

	// Dependency initialization
	LogDepDiagPermRepoInitOK  = "Repositorio de permisos de diagnóstico inicializado"
	LogDepDiagPermRepoInitErr = "Error inicializando repositorio de permisos de diagnóstico"
)

// ============================================
// COMPLETED SERVICE (HU64)
// ============================================
const (
	// Controller
	LogCSControllerCreateRequest     = "Solicitud de registro de servicio realizado recibida"
	LogCSControllerCreateError       = "Error registrando servicio realizado"
	LogCSControllerCreateSuccess     = "Servicio realizado registrado exitosamente"
	LogCSControllerListRequest       = "Solicitud de listado de servicios realizados recibida"
	LogCSControllerListError         = "Error listando servicios realizados"
	LogCSControllerListSuccess       = "Servicios realizados listados exitosamente"
	LogCSControllerListByMotoReq     = "Solicitud de listado de servicios realizados por moto recibida"
	LogCSControllerListByMotoError   = "Error listando servicios realizados por moto"
	LogCSControllerListByMotoSuccess = "Servicios realizados por moto listados exitosamente"

	// Interactor
	LogCSInteractorRegStart        = "Iniciando registro de servicio realizado"
	LogCSInteractorValidateSvcErr  = "Error al validar servicios de la sede"
	LogCSInteractorSvcValidated    = "Servicios de sede validados"
	LogCSInteractorValidateDiagErr = "Error al validar diagnóstico para motocicleta"
	LogCSInteractorDiagValidated   = "Diagnóstico validado para motocicleta"
	LogCSInteractorActiveCheckErr  = "Error al verificar servicio activo existente"
	LogCSInteractorActiveExists    = "Ya existe un servicio activo para esta motocicleta en esta sede"
	LogCSInteractorTxError         = "Error al iniciar transacción"
	LogCSInteractorRollbackError   = "Error en rollback de transacción"
	LogCSInteractorRollbackOK      = "Rollback de transacción ejecutado correctamente"
	LogCSInteractorSaveError       = "Error al guardar servicio realizado"
	LogCSInteractorSaved           = "Servicio realizado guardado"
	LogCSInteractorSaveItemsErr    = "Error al guardar items de servicio"
	LogCSInteractorItemsSaved      = "Items de servicio guardados"
	LogCSInteractorSaveHistoryErr  = "Error al guardar historial de estado"
	LogCSInteractorCommitError     = "Error al hacer commit de transacción"
	LogCSInteractorRegSuccess      = "Servicio realizado registrado exitosamente"
	LogCSInteractorGetByID         = "Consultando servicio realizado por ID"
	LogCSInteractorGetByIDErr      = "Error al obtener servicio realizado"
	LogCSInteractorGetByIDOK       = "Servicio realizado obtenido"
	LogCSInteractorGetByBranch     = "Consultando servicios realizados por sede"
	LogCSInteractorGetByBranchErr  = "Error al obtener servicios por sede"
	LogCSInteractorGetByBranchOK   = "Servicios por sede obtenidos"
	LogCSInteractorGetByMoto       = "Consultando servicios realizados por motocicleta"
	LogCSInteractorGetByMotoErr    = "Error al obtener servicios por motocicleta"
	LogCSInteractorGetByMotoOK     = "Servicios por motocicleta obtenidos"

	// Service
	LogCSServiceNoServices     = "No se proporcionaron servicios para validar"
	LogCSServiceDiagGetErr     = "Error al obtener diagnóstico"
	LogCSServiceDiagNotForMoto = "Diagnóstico no pertenece a la motocicleta"

	// Repository - Operations
	LogCSRepoPrepareError       = "Error preparando statement de servicio realizado"
	LogCSRepoSaveError          = "Error guardando servicio realizado en BD"
	LogCSRepoSaveItemErr        = "Error al guardar item de servicio"
	LogCSRepoSaveHistoryErr     = "Error al guardar historial de estado"
	LogCSRepoValidateSvcErr     = "Error al validar servicios de la sede"
	LogCSRepoGetError           = "Error obteniendo servicio realizado de BD"
	LogCSRepoGetByMotoErr       = "Error al obtener servicios por motocicleta"
	LogCSRepoGetByBranchErr     = "Error al obtener servicios por sede"
	LogCSRepoGetItemsErr        = "Error al obtener items del servicio"
	LogCSRepoScanItemErr        = "Error al escanear item del servicio"
	LogCSRepoScanError          = "Error escaneando fila de servicio realizado"
	LogCSRepoHasActiveErr       = "Error al verificar servicio activo"
	LogCSRepoPrepareInsert      = "Error preparando stmtInsert de completed_services"
	LogCSRepoPrepareInsertItem  = "Error preparando stmtInsertItem de completed_service_items"
	LogCSRepoPrepareHistory     = "Error preparando stmtInsertStatusHistory"
	LogCSRepoPrepareGetByID     = "Error preparando stmtGetByID de completed_services"
	LogCSRepoPrepareGetByMoto   = "Error preparando stmtGetByMotorcycleID"
	LogCSRepoPrepareGetByBranch = "Error preparando stmtGetByBranchID"
	LogCSRepoPrepareGetItems    = "Error preparando stmtGetItemsByCSID"
	LogCSRepoPrepareHasActive   = "Error preparando stmtHasActiveService"
	LogCSRepoPrepareDelete      = "Error preparando stmtDelete de completed_services"
	LogCSRepoDeleteError        = "Error eliminando servicio realizado de BD"
	LogCSRepoRowsCloseError     = "Error cerrando filas de resultado de BD"

	// Delete (HU65) - Controller
	LogCSControllerDeleteRequest = "Solicitud de eliminación de servicio realizado recibida"
	LogCSControllerDeleteError   = "Error eliminando servicio realizado"
	LogCSControllerDeleteSuccess = "Servicio realizado eliminado exitosamente"

	// Delete (HU65) - Interactor
	LogCSInteractorDelStart       = "Iniciando eliminación de servicio realizado"
	LogCSInteractorDelNotFound    = "Servicio realizado no encontrado para eliminar"
	LogCSInteractorDelTxErr       = "Error al iniciar transacción para eliminación"
	LogCSInteractorDelRollbackErr = "Error en rollback de transacción de eliminación"
	LogCSInteractorDelRollbackOK  = "Rollback de eliminación ejecutado correctamente"
	LogCSServiceDelStrategy       = "Estrategia de eliminación seleccionada"
	LogCSInteractorDelError       = "Error al eliminar servicio realizado"
	LogCSInteractorDelCommitErr   = "Error al hacer commit de eliminación"
	LogCSInteractorDelSuccess     = "Servicio realizado eliminado exitosamente"

	// Dependency initialization
	LogDepCSRepoInitOK       = "Repositorio de servicios realizados inicializado"
	LogDepCSRepoInitErr      = "Error inicializando repositorio de servicios realizados"
	LogDepCSInteractorInitOK = "Interactor de servicios realizados inicializado"

	// Motorcycle Controller - Completed Service Lookup
	LogMotorcycleControllerRepBranchErr = "Error obteniendo sedes del representante"
	LogMotorcycleControllerPermErr      = "Error consultando permisos de la motocicleta"
	LogMotorcycleControllerNoPerm       = "Representante sin permisos para esta motocicleta"

	// Status Transitions (HU73/HU74) - Controller
	LogCSControllerTransRequest  = "Solicitud de consulta de transiciones de estado"
	LogCSControllerTransError    = "Error al consultar transiciones de estado"
	LogCSControllerTransSuccess  = "Transiciones de estado consultadas exitosamente"
	LogCSControllerStatusRequest = "Solicitud de actualización de estado de servicio"
	LogCSControllerStatusError   = "Error al actualizar estado de servicio"
	LogCSControllerStatusSuccess = "Estado de servicio actualizado exitosamente"

	// Status Transitions (HU73/HU74) - Interactor
	LogCSInteractorTransStart     = "Consultando transiciones de estado"
	LogCSInteractorTransNotFound  = "Servicio no encontrado para consultar transiciones"
	LogCSInteractorTransError     = "Error al obtener transiciones de estado"
	LogCSInteractorTransSuccess   = "Transiciones de estado obtenidas"
	LogCSInteractorStatusStart    = "Iniciando transición de estado"
	LogCSInteractorStatusNotFound = "Servicio no encontrado para transición"
	LogCSInteractorStatusInvalid  = "Transición de estado no válida"
	LogCSInteractorStatusTxErr    = "Error al iniciar transacción para transición"
	LogCSInteractorStatusUpdErr   = "Error al actualizar estado en BD"
	LogCSInteractorStatusHistErr  = "Error al guardar historial de transición"
	LogCSInteractorStatusCommErr  = "Error al hacer commit de transición"
	LogCSInteractorStatusRbErr    = "Error en rollback de transición"
	LogCSInteractorStatusRbOK     = "Rollback de transición ejecutado"
	LogCSInteractorStatusSuccess  = "Transición de estado completada exitosamente"

	// Status Transitions (HU73/HU74) - Repository
	LogCSRepoPrepareUpdateStatus = "Error preparando stmtUpdateStatus"
	LogCSRepoPrepareGetHistory   = "Error preparando stmtGetStatusHistory"
	LogCSRepoUpdateStatusErr     = "Error actualizando estado en BD"
	LogCSRepoGetHistoryErr       = "Error obteniendo historial de transiciones"
	LogCSRepoScanHistoryErr      = "Error escaneando transición de estado"
)

// ============================================
// TRANSACTION HELPERS (Shared)
// ============================================
const (
	LogTxRollbackError = "ROLLBACK BD FALLÓ - ALERTA CRÍTICA"
	LogTxRollbackOK    = "Rollback BD ejecutado correctamente"
)
