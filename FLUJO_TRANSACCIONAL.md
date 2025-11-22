# Flujo Transaccional del Registro de Personas

## 📋 Tabla de Contenidos
1. [Introducción](#introducción)
2. [Flujo Completo del Interactor](#flujo-completo-del-interactor)
3. [Cómo Funcionan las Transacciones](#cómo-funcionan-las-transacciones)
4. [Casos Borde y Manejo de Errores](#casos-borde-y-manejo-de-errores)
5. [Recomendaciones](#recomendaciones)

---

## Introducción

El registro de personas es un caso de uso complejo que requiere **coordinación entre dos sistemas independientes**:
- **Base de Datos (MySQL)**: Sistema transaccional con soporte ACID
- **Keycloak**: Sistema externo de autenticación sin transacciones

Esta coordinación implementa un **patrón Saga compensatorio** para garantizar consistencia eventual entre ambos sistemas.

---

## Flujo Completo del Interactor

### Diagrama de Flujo

```
┌─────────────────────────────────────────────────────────────┐
│                   RegisterPerson (Interactor)                │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│ 1. Validaciones Iniciales (RegisterPerson)                   │
│    - Verificar email duplicado en BD                         │
│    - Si existe → Error: ErrDuplicateUser                     │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│ 2. Generar ID de Persona (person.SetID)                      │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│ 3. INICIAR TRANSACCIÓN DE BD (BeginTx)                       │
│    - Crea una TX que vivirá hasta el paso 8                  │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│ 4. Guardar Persona en BD (SavePersonToDB con TX)             │
│    ✅ Dentro de la transacción                               │
│    ❌ Si falla → Rollback TX + retornar error                │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│ 5. Crear Usuario en Keycloak (CreateUserInKeycloak)          │
│    ⚠️ Operación EXTERNA (sin TX)                             │
│    ❌ Si falla → Rollback TX + retornar error                │
│       (La BD aún no ha hecho commit)                         │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│ 6. Configurar Password en Keycloak (SetUserPassword)         │
│    ⚠️ Operación EXTERNA (sin TX)                             │
│    ❌ Si falla → Rollback Keycloak + Rollback TX             │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│ 7. Asignar Rol en Keycloak (AssignUserRole)                  │
│    ⚠️ Operación EXTERNA (sin TX)                             │
│    ❌ Si falla → Rollback Keycloak + Rollback TX             │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│ 8. Actualizar KeycloakID en BD (UpdatePersonKeycloakID)      │
│    ✅ Dentro de la MISMA transacción del paso 4              │
│    ❌ Si falla → Rollback Keycloak + Rollback TX             │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│ 9. COMMIT DE LA TRANSACCIÓN (tx.Commit)                      │
│    ✅ Si OK → Todo persistido en BD + Keycloak               │
│    ❌ Si falla → Rollback Keycloak + retornar error          │
│       (La TX ya no puede hacer rollback)                     │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│ 10. Retornar Resultado Exitoso                               │
└─────────────────────────────────────────────────────────────┘
```

### Código del Flujo

```go
func (i *Interactor) RegisterPerson(ctx context.Context, person domain.Person) (*dto.RegistrationResult, error) {
	// 1. Validaciones iniciales (email duplicado)
	result, err := i.service.RegisterPerson(ctx, person)
	if err != nil {
		return nil, err
	}

	// 2. Generar ID
	person.SetID()

	// 3. Iniciar transacción de BD
	tx, err := i.service.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	// 4. Guardar persona en BD dentro de la transacción
	if err = i.service.SavePersonToDB(ctx, tx, person); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	// 5. Crear usuario en Keycloak
	keycloakUserID, err := i.service.CreateUserInKeycloak(ctx, &person)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	// 6. Configurar password en Keycloak
	err = i.service.SetUserPassword(ctx, keycloakUserID, person.Password)
	if err != nil {
		_ = i.service.RollbackKeycloakUser(ctx, keycloakUserID)
		_ = tx.Rollback()
		return nil, err
	}

	// 7. Asignar rol en Keycloak
	err = i.service.AssignUserRole(ctx, keycloakUserID, person.Role)
	if err != nil {
		_ = i.service.RollbackKeycloakUser(ctx, keycloakUserID)
		_ = tx.Rollback()
		return nil, err
	}

	// 8. Actualizar keycloakID en BD dentro de la transacción
	err = i.service.UpdatePersonKeycloakID(ctx, tx, person.ID, keycloakUserID)
	if err != nil {
		_ = i.service.RollbackKeycloakUser(ctx, keycloakUserID)
		_ = tx.Rollback()
		return nil, err
	}

	// 9. Confirmar transacción de BD (commit final)
	if err = tx.Commit(); err != nil {
		_ = i.service.RollbackKeycloakUser(ctx, keycloakUserID)
		return nil, err
	}

	// 10. Retornar resultado exitoso
	person.KeycloakUserID = keycloakUserID
	result.Person = person
	result.Message = "Usuario registrado exitosamente"

	return result, nil
}
```

---

## Cómo Funcionan las Transacciones

### 1. Arquitectura de Capas

```
┌──────────────────────────────────────────────────────┐
│                    INTERACTOR                         │
│  - Coordina el flujo completo                         │
│  - Maneja la TX de principio a fin                    │
│  - Implementa Saga compensatorio                      │
└──────────────────────────────────────────────────────┘
                       │
                       ▼
┌──────────────────────────────────────────────────────┐
│                     SERVICE                           │
│  - BeginTx(): crea la transacción                     │
│  - SavePersonToDB(tx): pasa TX al repositorio         │
│  - UpdatePersonKeycloakID(tx): pasa TX al repositorio │
└──────────────────────────────────────────────────────┘
                       │
                       ▼
┌──────────────────────────────────────────────────────┐
│                   REPOSITORY                          │
│  - BeginTx(): crea sql.Tx y lo envuelve en sqlTx      │
│  - SavePerson(tx): usa TX recibida o crea nueva       │
│  - PatchPerson(tx): usa TX recibida o crea nueva      │
└──────────────────────────────────────────────────────┘
                       │
                       ▼
┌──────────────────────────────────────────────────────┐
│                   DATABASE (MySQL)                    │
│  - Maneja transacciones ACID                          │
└──────────────────────────────────────────────────────┘
```

### 2. Wrapper de Transacciones

Para mantener la independencia de capas, creamos un wrapper `sqlTx` que implementa la interfaz `output.Tx`:

```go
// En repository.go
type sqlTx struct {
	*sql.Tx
}

func (t *sqlTx) Commit() error {
	return t.Tx.Commit()
}

func (t *sqlTx) Rollback() error {
	return t.Tx.Rollback()
}

func (r *repository) BeginTx(ctx context.Context) (output.Tx, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &sqlTx{Tx: tx}, nil
}
```

### 3. Repositorio con TX Opcional

Los métodos del repositorio aceptan una TX opcional:
- **Si `tx != nil`**: usan esa transacción (modo controlado por interactor)
- **Si `tx == nil`**: crean su propia TX y hacen commit/rollback (modo autónomo)

```go
func (r *repository) SavePerson(ctx context.Context, tx output.Tx, person domain.Person) error {
	var dbTx *sqlTx
	var shouldCommit bool

	if tx != nil {
		// Usar la transacción existente
		dbTx = tx.(*sqlTx)
		shouldCommit = false
	} else {
		// Crear nueva transacción
		newTx, err := r.db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		dbTx = &sqlTx{Tx: newTx}
		shouldCommit = true
	}

	// Ejecutar operación
	_, err := dbTx.ExecContext(ctx, querySave, ...)
	if err != nil {
		if shouldCommit {
			dbTx.Rollback()
		}
		return err
	}

	// Commit solo si creamos la TX
	if shouldCommit {
		return dbTx.Commit()
	}
	return nil
}
```

### 4. Atomicidad Entre Operaciones

La clave es que `SavePerson` y `UpdatePersonKeycloakID` **comparten la misma transacción**:

```go
tx, _ := i.service.BeginTx(ctx)

// Ambas operaciones usan la MISMA tx
i.service.SavePersonToDB(ctx, tx, person)      // Dentro de TX
i.service.UpdatePersonKeycloakID(ctx, tx, ...) // Dentro de LA MISMA TX

tx.Commit() // Ambas se persisten juntas o ninguna
```

Esto garantiza que:
- ✅ Si ambas tienen éxito → ambas se persisten
- ✅ Si alguna falla → ninguna se persiste
- ✅ No hay estado inconsistente en BD

---

## Casos Borde y Manejo de Errores

### Caso 1: ❌ Falla durante SavePerson (paso 4)

**Escenario:**
```go
tx, _ := i.service.BeginTx(ctx)
err := i.service.SavePersonToDB(ctx, tx, person) // ❌ FALLA
```

**¿Qué pasa?**
- ❌ TX hace rollback automático
- ✅ Nada se persiste en BD
- ✅ Keycloak no se ha tocado aún
- ✅ Estado consistente

**Código:**
```go
if err = i.service.SavePersonToDB(ctx, tx, person); err != nil {
	_ = tx.Rollback() // Limpia todo
	return nil, err
}
```

### Caso 2: ❌ Falla durante CreateUserInKeycloak (paso 5)

**Escenario:**
```go
tx, _ := i.service.BeginTx(ctx)
i.service.SavePersonToDB(ctx, tx, person)             // ✅ OK (en TX pendiente)
_, err := i.service.CreateUserInKeycloak(ctx, &person) // ❌ FALLA
```

**¿Qué pasa?**
- ✅ Persona guardada en BD pero **TX aún no committeada**
- ❌ Usuario no creado en Keycloak
- ✅ TX hace rollback
- ✅ Nada se persiste
- ✅ Estado consistente

**Código:**
```go
keycloakUserID, err := i.service.CreateUserInKeycloak(ctx, &person)
if err != nil {
	_ = tx.Rollback() // La BD vuelve al estado anterior
	return nil, err
}
```

### Caso 3: ❌ Falla durante SetUserPassword (paso 6)

**Escenario:**
```go
tx, _ := i.service.BeginTx(ctx)
i.service.SavePersonToDB(ctx, tx, person)        // ✅ OK (en TX pendiente)
keycloakID, _ := i.service.CreateUserInKeycloak(...)  // ✅ OK (usuario creado)
err := i.service.SetUserPassword(...)            // ❌ FALLA
```

**¿Qué pasa?**
- ✅ Persona guardada en BD pero **TX aún no committeada**
- ✅ Usuario creado en Keycloak (SIN password asignado)
- ❌ Error al configurar password
- ✅ Se ejecuta rollback compensatorio de Keycloak
- ✅ TX de BD hace rollback
- ✅ Estado consistente

**Código:**
```go
err = i.service.SetUserPassword(ctx, keycloakUserID, person.Password)
if err != nil {
	_ = i.service.RollbackKeycloakUser(ctx, keycloakUserID) // Elimina usuario de Keycloak
	_ = tx.Rollback()                                       // Elimina persona de BD
	return nil, err
}
```

### Caso 4: ⚠️ Falla el Rollback de Keycloak

**Escenario:**
```go
err = i.service.SetUserPassword(...)
if err != nil {
	deleteErr := i.service.RollbackKeycloakUser(ctx, keycloakUserID) // ❌ FALLA
	_ = tx.Rollback() // ✅ OK
	return nil, err
}
```

**¿Qué pasa?**
- ❌ Usuario queda huérfano en Keycloak (sin registro en BD)
- ✅ BD rollback exitoso (no hay registro)
- ⚠️ **Inconsistencia temporal**: Keycloak tiene usuario que la BD no conoce

**Solución:**
1. **Logging crítico**: registrar el error de rollback de Keycloak
2. **Job de limpieza**: proceso batch que detecta usuarios huérfanos en Keycloak
3. **Retry automático**: intentar eliminar el usuario antes de retornar

**Código mejorado:**
```go
err = i.service.SetUserPassword(...)
if err != nil {
	// Intentar rollback de Keycloak con retry
	if rbErr := i.service.RollbackKeycloakUser(ctx, keycloakUserID); rbErr != nil {
		// ⚠️ CRÍTICO: Registrar para limpieza posterior
		log.Error().
			Err(rbErr).
			Str("keycloakUserID", keycloakUserID).
			Str("personID", person.ID).
			Msg("CRÍTICO: Rollback de Keycloak falló - usuario huérfano")
	}
	_ = tx.Rollback()
	return nil, err
}
```

### Caso 5: ❌ Falla el Commit de la TX (paso 9)

**Escenario:**
```go
tx, _ := i.service.BeginTx(ctx)
i.service.SavePersonToDB(ctx, tx, person)              // ✅ OK
keycloakID, _ := i.service.CreateUserInKeycloak(...)   // ✅ OK
i.service.SetUserPassword(...)                         // ✅ OK
i.service.AssignUserRole(...)                          // ✅ OK
i.service.UpdatePersonKeycloakID(ctx, tx, ...)         // ✅ OK
err := tx.Commit()                                     // ❌ FALLA
```

**¿Qué pasa?**
- ✅ Usuario creado en Keycloak con password y rol
- ❌ TX de BD falla al commitear (pérdida de conexión, constraint violation, etc.)
- ⚠️ **Inconsistencia**: Keycloak tiene usuario que la BD no tiene

**⚠️ CRÍTICO: No podemos hacer rollback de la TX después del commit**

La transacción de BD se descarta automáticamente cuando falla el commit, pero Keycloak ya está modificado.

**Solución:**
```go
if err = tx.Commit(); err != nil {
	// Intentar rollback compensatorio de Keycloak
	if rbErr := i.service.RollbackKeycloakUser(ctx, keycloakUserID); rbErr != nil {
		// ⚠️ DOBLE ERROR: Commit falló Y rollback de Keycloak falló
		log.Error().
			Err(err).
			Err(rbErr).
			Str("keycloakUserID", keycloakUserID).
			Msg("CRÍTICO: Commit falló y rollback de Keycloak también falló")
	}
	return nil, err // Retornar error original del commit
}
```

### Caso 6: ⚠️ Falla el Rollback de la TX

**Escenario:**
```go
tx, _ := i.service.BeginTx(ctx)
i.service.SavePersonToDB(ctx, tx, person)
err := i.service.CreateUserInKeycloak(...)  // ❌ FALLA
rbErr := tx.Rollback()                      // ❌ TAMBIÉN FALLA
```

**¿Qué pasa?**
- ❌ Operación de Keycloak falló
- ❌ Rollback de TX falló (pérdida de conexión, timeout, etc.)
- ⚠️ Estado de la TX es **indefinido**

**Comportamiento de MySQL:**
- Si se pierde la conexión → la TX se rollbackea automáticamente por el servidor
- Si el servidor está caído → la TX eventualmente se descarta
- En la mayoría de casos, el sistema se auto-recupera

**Recomendación:**
```go
if err = i.service.SavePersonToDB(ctx, tx, person); err != nil {
	if rbErr := tx.Rollback(); rbErr != nil {
		// Registrar pero NO intentar compensar
		// La BD probablemente hará rollback automático
		log.Error().
			Err(rbErr).
			Err(err).
			Msg("Rollback de TX falló - la BD debería auto-recuperarse")
	}
	return nil, err
}
```

---

## Recomendaciones

### ✅ Buenas Prácticas Implementadas

1. **Interactor como orquestador**
   - ✅ El interactor controla el ciclo de vida de la TX
   - ✅ El service solo expone operaciones atómicas
   - ✅ El repository ejecuta sin conocer el contexto completo

2. **TX explícitas en la firma**
   - ✅ `SavePersonToDB(ctx, tx, person)` deja claro que acepta una TX
   - ✅ El que llama decide si pasa TX o `nil`
   - ✅ Backward compatible

3. **Separación de responsabilidades**
   - ✅ Repository: maneja la infraestructura de BD
   - ✅ Service: agrupa operaciones relacionadas
   - ✅ Interactor: coordina el caso de uso completo

4. **Rollback compensatorio**
   - ✅ Si falla algo después de Keycloak, se elimina el usuario
   - ✅ Implementa patrón Saga de forma simple

### 🔧 Mejoras Recomendadas

#### 1. Logging Estructurado para Casos Críticos

```go
// En el interactor
import "github.com/rs/zerolog/log"

err = i.service.SetUserPassword(ctx, keycloakUserID, person.Password)
if err != nil {
	if rbErr := i.service.RollbackKeycloakUser(ctx, keycloakUserID); rbErr != nil {
		log.Error().
			Err(rbErr).
			Str("keycloakUserID", keycloakUserID).
			Str("personID", person.ID).
			Str("email", person.Email).
			Msg("CRITICAL: Keycloak rollback failed - orphaned user")
	}
	_ = tx.Rollback()
	return nil, err
}
```

#### 2. Job de Limpieza de Usuarios Huérfanos

Crea un proceso batch que:
- Busca usuarios en Keycloak sin registro en BD
- Los elimina automáticamente
- Se ejecuta periódicamente (ej: cada hora)

```go
// En un worker separado
func CleanOrphanedKeycloakUsers(ctx context.Context) error {
	// 1. Obtener todos los usuarios de Keycloak
	keycloakUsers, _ := keycloak.GetAllUsers(ctx)
	
	// 2. Para cada usuario, verificar si existe en BD
	for _, kcUser := range keycloakUsers {
		person, err := repo.GetPersonByKeycloakID(ctx, nil, kcUser.ID)
		if err != nil || person == nil {
			// Usuario huérfano - eliminar
			log.Warn().
				Str("keycloakUserID", kcUser.ID).
				Msg("Limpiando usuario huérfano de Keycloak")
			_ = keycloak.DeleteUser(ctx, kcUser.ID)
		}
	}
	return nil
}
```

#### 3. Retry con Backoff para Rollback

```go
import "github.com/cenkalti/backoff/v4"

func (s *service) RollbackKeycloakUserWithRetry(ctx context.Context, userID string) error {
	operation := func() error {
		return s.keycloak.DeleteUser(ctx, userID)
	}
	
	exponentialBackoff := backoff.NewExponentialBackOff()
	exponentialBackoff.MaxElapsedTime = 30 * time.Second
	
	return backoff.Retry(operation, exponentialBackoff)
}
```

#### 4. Métricas y Alertas

```go
import "github.com/prometheus/client_golang/prometheus"

var (
	registrationAttempts = prometheus.NewCounter(...)
	registrationSuccesses = prometheus.NewCounter(...)
	registrationFailures = prometheus.NewCounterVec(...)
	keycloakRollbackFailures = prometheus.NewCounter(...)
)

// En el interactor
registrationAttempts.Inc()

if err := tx.Commit(); err != nil {
	registrationFailures.WithLabelValues("commit_failed").Inc()
	
	if rbErr := i.service.RollbackKeycloakUser(ctx, keycloakUserID); rbErr != nil {
		keycloakRollbackFailures.Inc() // ⚠️ Alerta crítica
	}
	return nil, err
}

registrationSuccesses.Inc()
```

#### 5. Context con Timeout

```go
func (i *Interactor) RegisterPerson(ctx context.Context, person domain.Person) (*dto.RegistrationResult, error) {
	// Timeout de 30 segundos para todo el flujo
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	
	// ... resto del código
}
```

#### 6. Idempotencia para Reintentos

Si el cliente reintenta el registro, debemos evitar duplicados:

```go
func (i *Interactor) RegisterPerson(ctx context.Context, person domain.Person) (*dto.RegistrationResult, error) {
	// Verificar si ya existe
	existing, _ := i.service.GetPersonByEmail(ctx, person.Email)
	if existing != nil {
		// Ya existe - verificar si está completo
		if existing.KeycloakUserID != "" {
			// Registro completo - retornar éxito (idempotente)
			return &dto.RegistrationResult{
				Person:  *existing,
				Message: "Usuario ya registrado",
			}, nil
		}
		// Registro incompleto - eliminar y reintentar
		// (Saga anterior falló a mitad de camino)
		_ = i.service.RollbackPerson(ctx, existing.ID)
	}
	
	// ... continuar con el flujo normal
}
```

### 🚫 Anti-Patrones a Evitar

1. **❌ NO hacer commit/rollback en el service**
   ```go
   // MAL - el service no debe controlar la TX
   func (s *service) SavePerson(ctx context.Context, person domain.Person) error {
       tx, _ := s.repo.BeginTx(ctx)
       s.repo.SavePerson(ctx, tx, person)
       tx.Commit() // ❌ No
   }
   ```

2. **❌ NO crear múltiples transacciones para una misma operación lógica**
   ```go
   // MAL - dos TXs separadas pueden quedar inconsistentes
   tx1, _ := i.service.BeginTx(ctx)
   i.service.SavePersonToDB(ctx, tx1, person)
   tx1.Commit()
   
   tx2, _ := i.service.BeginTx(ctx) // ❌ Nueva TX
   i.service.UpdatePersonKeycloakID(ctx, tx2, ...)
   tx2.Commit()
   ```

3. **❌ NO ignorar errores de rollback silenciosamente**
   ```go
   // MAL - si el rollback falla, hay que registrarlo
   if err != nil {
       tx.Rollback() // ❌ Error ignorado
       return err
   }
   
   // BIEN
   if err != nil {
       if rbErr := tx.Rollback(); rbErr != nil {
           log.Error().Err(rbErr).Msg("Rollback failed")
       }
       return err
   }
   ```

4. **❌ NO hacer operaciones externas después del commit**
   ```go
   // MAL - si falla el email, la BD ya commiteo
   tx.Commit()
   sendWelcomeEmail(person.Email) // ❌ Después del commit
   
   // BIEN - hacer antes del commit o en proceso asíncrono
   tx.Commit()
   go sendWelcomeEmail(person.Email) // Async, no crítico
   ```

---

## Resumen

### ✅ Lo que funciona bien

- ✅ Atomicidad entre operaciones de BD (SavePerson + UpdatePersonKeycloakID)
- ✅ Rollback compensatorio de Keycloak cuando falla algo
- ✅ Separación clara de responsabilidades
- ✅ Código simple sin patrones innecesarios

### ⚠️ Casos borde que requieren atención

- ⚠️ Falla el rollback de Keycloak → usuario huérfano
- ⚠️ Falla el commit de TX → usuario en Keycloak sin registro en BD
- ⚠️ Falla el rollback de TX → estado indefinido (usualmente se auto-recupera)

### 🔧 Mejoras sugeridas

1. Logging estructurado para casos críticos
2. Job de limpieza de usuarios huérfanos
3. Retry con backoff para rollback de Keycloak
4. Métricas y alertas
5. Timeout en el contexto
6. Idempotencia para reintentos del cliente

---

**Última actualización:** Noviembre 2025  
**Autor:** Esteban (con asistencia de Cascade)
