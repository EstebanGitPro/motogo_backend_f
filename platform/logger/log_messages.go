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
	LogDepInit          = "Inicializando dependencias"
	LogDepInitComplete  = "Dependencias inicializadas completamente"
	LogDepInitError     = "Error inicializando dependencias"
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
	LogPersonInteractorCommit_Error         = "COMMIT FALLÓ - ALERTA CRÍTICA"
	LogPersonInteractorCommit_OK            = "Transacción confirmada exitosamente"
	LogPersonInteractorRegComplete          = "Registro completado exitosamente"
	LogPersonInteractorRollbackDB_Error     = "ROLLBACK BD FALLÓ - ALERTA CRÍTICA"
	LogPersonInteractorRollbackDB_OK        = "Rollback BD ejecutado correctamente"
	LogPersonInteractorRollbackKeycloak_Err = "ROLLBACK KEYCLOAK FALLÓ - ALERTA CRÍTICA"
	LogPersonInteractorRollbackKeycloak_OK  = "Rollback Keycloak ejecutado correctamente"
	LogPersonInteractorIncompleteDetected   = "Registro incompleto detectado"
	LogPersonInteractorCleanup_Error        = "Error limpiando estado inconsistente"
	LogPersonInteractorCleanup_OK           = "Estado inconsistente limpiado exitosamente"
	LogPersonInteractorLoginStart           = "Inicio de sesión en Keycloak"
	LogPersonInteractorLoginError           = "Error iniciando sesión en Keycloak"
	LogPersonInteractorLoginOK              = "Sesión iniciada exitosamente"
	LogPersonInteractorLoginComplete        = "Sesión iniciada exitosamente"
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
	LogMessageInteractorCreateCommitErr  = "COMMIT FALLÓ - ALERTA CRÍTICA"
	LogMessageInteractorCreateCommitOK   = "Transacción confirmada exitosamente"
	LogMessageInteractorCreateComplete   = "Mensaje creado exitosamente"

	// UPDATE flow
	LogMessageInteractorUpdateStep1Error = "[PASO 1/4] Mensaje no encontrado"
	LogMessageInteractorUpdateStep1OK    = "[PASO 1/4] Mensaje encontrado"
	LogMessageInteractorUpdateStep2Error = "[PASO 2/4] Validación de mensaje fallida"
	LogMessageInteractorUpdateStep2OK    = "[PASO 2/4] Validación de mensaje completada"
	LogMessageInteractorUpdateStep3Error = "[PASO 3/4] Error iniciando transacción"
	LogMessageInteractorUpdateStep3OK    = "[PASO 3/4] Transacción iniciada"
	LogMessageInteractorUpdateStep4Error = "[PASO 4/4] Error actualizando mensaje"
	LogMessageInteractorUpdateStep4OK    = "[PASO 4/4] Mensaje actualizado en BD"
	LogMessageInteractorUpdateCommitErr  = "COMMIT FALLÓ - ALERTA CRÍTICA"
	LogMessageInteractorUpdateCommitOK   = "Transacción confirmada exitosamente"
	LogMessageInteractorUpdateComplete   = "Mensaje actualizado exitosamente"

	// DELETE flow
	LogMessageInteractorDeleteStep1Error = "[PASO 1/3] Mensaje no encontrado"
	LogMessageInteractorDeleteStep1OK    = "[PASO 1/3] Mensaje encontrado"
	LogMessageInteractorDeleteStep2Error = "[PASO 2/3] Error iniciando transacción"
	LogMessageInteractorDeleteStep2OK    = "[PASO 2/3] Transacción iniciada"
	LogMessageInteractorDeleteStep3Error = "[PASO 3/3] Error eliminando mensaje"
	LogMessageInteractorDeleteStep3OK    = "[PASO 3/3] Mensaje eliminado de BD"
	LogMessageInteractorDeleteCommitErr  = "COMMIT FALLÓ - ALERTA CRÍTICA"
	LogMessageInteractorDeleteCommitOK   = "Transacción confirmada exitosamente"
	LogMessageInteractorDeleteComplete   = "Mensaje eliminado exitosamente"

	// Common rollback
	LogMessageInteractorRollbackError = "ROLLBACK BD FALLÓ - ALERTA CRÍTICA"
	LogMessageInteractorRollbackOK    = "Rollback BD ejecutado correctamente"
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
	LogBranchInteractorRollbackError   = "ROLLBACK BD FALLÓ - ALERTA CRÍTICA"
	LogBranchInteractorRollbackOK      = "Rollback BD ejecutado correctamente"
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
)

// ============================================
// BRANCH REPOSITORY (HU59)
// ============================================
const (
	LogBranchRepoSaveError        = "Error guardando sede en BD"
	LogBranchRepoUpdateError      = "Error actualizando sede en BD"
	LogBranchRepoDeleteError      = "Error eliminando sede de BD"
	LogBranchRepoGetByIDError     = "Error obteniendo sede por ID"
	LogBranchRepoGetByNameError   = "Error obteniendo sede por nombre"
	LogBranchRepoGetByRepError    = "Error obteniendo sedes por representante"
	LogBranchRepoScanError        = "Error escaneando fila de sede"
	LogBranchRepoLocationSaveErr  = "Error guardando ubicación"
	LogBranchRepoLocationUpdErr   = "Error actualizando ubicación"
	LogBranchRepoBrandSaveError   = "Error guardando marca de sede"
	LogBranchRepoBrandDelError    = "Error eliminando marcas de sede"
	LogBranchRepoBrandGetError    = "Error obteniendo marcas de sede"
	LogBranchRepoBrandValidateErr = "Error validando marcas"
)

// ============================================
// BRANCH SERVICE (HU59)
// ============================================
const (
	LogBranchServiceInvalidType  = "Tipo de establecimiento inválido"
	LogBranchServiceDupNameCheck = "Error verificando nombre duplicado"
	LogBranchServiceDupName      = "Nombre de sede duplicado en franquicia"
	LogBranchServiceSaveError    = "Error guardando sede"
	LogBranchServiceLocSaveError = "Error guardando ubicación"
	LogBranchServiceBrandSaveErr = "Error guardando marcas"
	LogBranchServiceRegComplete  = "Sede registrada exitosamente"
	LogBranchServiceGetError     = "Error obteniendo sede por ID"
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
	LogLocationRepoGetDepartmentsError     = "Error obteniendo departamentos"
	LogLocationRepoGetDepartmentsScanError = "Error escaneando departamento"
	LogLocationRepoGetDepartmentsIterError = "Error iterando departamentos"
	LogLocationRepoGetCitiesError          = "Error obteniendo ciudades"
	LogLocationRepoGetCitiesScanError      = "Error escaneando ciudad"
	LogLocationRepoGetCitiesIterError      = "Error iterando ciudades"
	LogLocationRepoValidateCityError       = "Error validando ciudad en departamento"
	LogLocationRepoGetDeptByIDError        = "Error obteniendo departamento por ID"
	LogLocationRepoSaveError               = "Error guardando ubicación"
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
	LogFranchiseServiceUpdateError  = "Error actualizando franquicia"
	LogFranchiseServiceDeleteError  = "Error eliminando franquicia"
	LogFranchiseServiceDeleted      = "Franquicia eliminada exitosamente"
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
	LogFranchiseControllerUpdateError      = "Error actualizando franquicia"
	LogFranchiseControllerUpdateSuccess    = "Franquicia actualizada exitosamente"
	LogFranchiseControllerDeleteError      = "Error eliminando franquicia"
	LogFranchiseControllerDeleteSuccess    = "Franquicia eliminada exitosamente"
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
	LogDepBranchRepoInit       = "Inicializando repositorio de sedes"
	LogDepBranchRepoInitOK     = "Repositorio de sedes inicializado"
	LogDepBranchRepoInitErr    = "Error inicializando repositorio de sedes"
	LogDepLocationRepoInitOK   = "Repositorio de ubicaciones inicializado"
	LogDepLocationRepoInitErr  = "Error inicializando repositorio de ubicaciones"
	LogDepBrandRepoInitOK      = "Repositorio de marcas inicializado"
	LogDepBrandRepoInitErr     = "Error inicializando repositorio de marcas"
	LogDepFranchiseRepoInitOK  = "Repositorio de franquicias inicializado"
	LogDepFranchiseRepoInitErr = "Error inicializando repositorio de franquicias"

	// Interactor initialization
	LogDepBranchInteractorInitOK    = "Interactor de sedes inicializado"
	LogDepBrandInteractorInitOK     = "Interactor de marcas inicializado"
	LogDepLocationInteractorInitOK  = "Interactor de ubicaciones inicializado"
	LogDepFranchiseInteractorInitOK = "Interactor de franquicias inicializado"

	// External services
	LogDepGeocodingClientInitOK = "Cliente de geocodificación inicializado"
	LogDepFirebaseClientInitOK  = "Cliente de Firebase inicializado"
	LogDepJWKSValidatorInitOK   = "Validador JWKS inicializado"
)
