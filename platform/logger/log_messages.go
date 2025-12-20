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
	LogPersonServiceSearchByEmail             = "Buscando persona por email"
	LogPersonServiceSearchByID                = "Buscando persona por ID"
	LogPersonServiceFoundByEmail              = "Persona encontrada por email"
	LogPersonServiceFoundByID                 = "Persona encontrada por ID"
	LogPersonServiceErrorByEmail              = "Error buscando persona por email"
	LogPersonServiceErrorByID                 = "Error buscando persona por ID"
	LogPersonServiceValidationStart           = "Iniciando validaciones de registro"
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
