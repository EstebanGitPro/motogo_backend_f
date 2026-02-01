# Mapeo HUs → Endpoints REST

Este documento presenta el mapeo entre las Historias de Usuario (HUs) en español y los endpoints REST en inglés, siguiendo las mejores prácticas RESTful.

---

## Convenciones de Nomenclatura

| Aspecto | Convención |
|---------|------------|
| **Endpoints (URIs)** | Inglés, kebab-case, plural para colecciones |
| **Verbos HTTP** | GET (consultar), POST (crear), PUT/PATCH (modificar), DELETE (eliminar) |
| **Nombres HUs** | Español, orientados al negocio |
| **Código interno** | Inglés (funciones, variables, servicios) |

---

## Concepto Clave: Person como Recurso Unificado

> [!IMPORTANT]
> Tanto el **Representante de Sede** como el **Usuario de Motocicleta** son `person` en la base de datos.
> Los endpoints usan `/persons` como recurso base, diferenciando por un campo `type` o `role`.

**Estructura:**
```
persons (tabla BD)
├── type: 'branch_representative' → Representante de Sede
└── type: 'motorcycle_user' → Usuario de Motocicleta
```

**Endpoints unificados:**
```
POST   /persons              → Crear persona (con type en body)
GET    /persons/me           → Obtener datos del usuario autenticado
PUT    /persons/me           → Modificar datos del usuario autenticado
DELETE /persons/{id}         → Eliminar persona
GET    /persons              → Listar personas (con filtro opcional ?type=...)
GET    /persons/{id}         → Obtener persona por ID
PUT    /persons/me/password  → Modificar contraseña del usuario autenticado
```

---

## Mapeo por Módulo

### Persona (Person) - Incluye Representante y Usuario

| HU | Nombre HU | Método | Endpoint Sugerido | Endpoint Actual | Estado |
|----|-----------|--------|-------------------|-----------------|--------|
| HU51 | Registrar Información del Representante Sede | POST | `/persons` (con `type: branch_representative`) | `/persons` | ✅ OK |
| HU52 | Modificar Información del Representante Sede | PUT | `/persons/me` | — | ❌ No implementado |
| HU53 | Eliminar Información del Representante Sede | DELETE | `/persons/{id}` | — | ❌ No implementado |
| HU54 | Consultar Información detallada del Representante Sede | GET | `/persons/me` | `/persons/me` | ✅ OK |
| HU55 | Consultar Información General del Representante Sede | GET | `/persons?type=branch_representative` | — | ❌ No implementado |
| HU56 | Restablecer la Clave del Representante Sede | POST | `/auth/password-reset` | `/auth/password-reset` | ✅ OK |
| HU57 | Modificar la Clave del Representante Sede | PUT | `/persons/me/password` | — | ❌ No implementado |
| HU58 | Autenticar Información del Representante Sede | POST | `/auth/login` | `/auth/login` | ✅ OK |
| HU81 | Registrar Información del Usuario | POST | `/persons` (con `type: motorcycle_user`) | — | ❌ No implementado |
| HU82 | Modificar Información del Usuario | PUT | `/persons/me` | — | ❌ No implementado |
| HU83 | Eliminar Información del Usuario | DELETE | `/persons/{id}` | — | ❌ No implementado |
| HU84 | Consultar Información Detallada del Usuario | GET | `/persons/me` | — | ❌ No implementado |
| HU85 | Consultar Información general del Usuario | GET | `/persons?type=motorcycle_user` | — | ❌ No implementado |
| HU86 | Restablecer la Clave del Usuario | POST | `/auth/password-reset` | `/auth/password-reset` | ✅ OK (compartido) |
| HU87 | Modificar la Clave del Usuario | PUT | `/persons/me/password` | — | ❌ No implementado |
| HU88 | Autenticar el Usuario | POST | `/auth/login` | `/auth/login` | ✅ OK (compartido) |

---

### Sede (Branch)

| HU | Nombre HU | Método | Endpoint Sugerido | Endpoint Actual | Estado |
|----|-----------|--------|-------------------|-----------------|--------|
| HU59 | Registrar Información de una nueva Sede | POST | `/branches` | `/branches` | ✅ OK |
| HU60 | Modificar la Información de la Sede | PUT | `/branches/{id}` | — | ❌ No implementado |
| HU61 | Eliminar Información de la Sede | DELETE | `/branches/{id}` | — | ❌ No implementado |
| HU62 | Consultar Información de la Sede | GET | `/branches/{id}` | — | ❌ No implementado |
| HU76 | Consultar Información del Tipo Sede | GET | `/branch-types` | — | ❌ No implementado |


---

### Servicio Sede (Branch Service)

| HU | Nombre HU | Método | Endpoint Sugerido | Endpoint Actual | Estado |
|----|-----------|--------|-------------------|-----------------|--------|
| HU67 | Crear Información del Servicio Sede | POST | `/branches/{branchId}/services` | — | ❌ No implementado |
| HU68 | Modificar Información del Servicio Sede | PUT | `/branches/{branchId}/services/{id}` | — | ❌ No implementado |
| HU69 | Eliminar Información del Servicio Sede | DELETE | `/branches/{branchId}/services/{id}` | — | ❌ No implementado |
| HU70 | Consultar Información del Servicio Sede | GET | `/branches/{branchId}/services/{id}` | — | ❌ No implementado |
| HU71 | Activar el Servicio Sede | PATCH | `/branches/{branchId}/services/{id}/activate` | — | ❌ No implementado |
| HU72 | Desactivar el Servicio Sede | PATCH | `/branches/{branchId}/services/{id}/deactivate` | — | ❌ No implementado |
| HU75 | Consultar Información del Tipo Servicio | GET | `/service-types` | — | ❌ No implementado |

---

### Servicio Realizado (Completed Service)

| HU | Nombre HU | Método | Endpoint Sugerido | Endpoint Actual | Estado |
|----|-----------|--------|-------------------|-----------------|--------|
| HU64 | Registrar Información del Servicio Realizado | POST | `/completed-services` | — | ❌ No implementado |
| HU65 | Eliminar Información del Servicio Realizado | DELETE | `/completed-services/{id}` | — | ❌ No implementado |
| HU66 | Consultar Información del Servicio Realizado | GET | `/completed-services/{id}` | — | ❌ No implementado |
| HU15 | Consultar Información del Estado del Servicio Realizado | GET | `/service-statuses` | — | ❌ No implementado |
| HU73 | Consultar Información de la Transición Estado Servicio Realizado | GET | `/completed-services/{id}/transitions` | — | ❌ No implementado |
| HU74 | Actualizar la Información de la Transición Estado Servicio Realizado | PUT | `/completed-services/{id}/status` | — | ❌ No implementado |

---

### Franquicia (Franchise)

| HU | Nombre HU | Método | Endpoint Sugerido | Endpoint Actual | Estado |
|----|-----------|--------|-------------------|-----------------|--------|
| HU26 | Registrar Información de la Franquicia | POST | `/franchises` | — | ❌ No implementado |
| HU27 | Modificar Información de la Franquicia | PUT | `/franchises/{id}` | — | ❌ No implementado |
| HU28 | Eliminar Información de la Franquicia | DELETE | `/franchises/{id}` | — | ❌ No implementado |
| HU29 | Consultar Información de la Franquicia | GET | `/franchises/{id}` | — | ❌ No implementado |

---

### Horario Sede (Branch Schedule)

| HU | Nombre HU | Método | Endpoint Sugerido | Endpoint Actual | Estado |
|----|-----------|--------|-------------------|-----------------|--------|
| HU30 | Registrar Información del Horario Sede | POST | `/branches/{branchId}/schedules` | — | ❌ No implementado |
| HU31 | Modificar Información del Horario Sede | PUT | `/branches/{branchId}/schedules/{id}` | — | ❌ No implementado |
| HU32 | Consultar Información del Horario Sede | GET | `/branches/{branchId}/schedules` | — | ❌ No implementado |
| HU33 | Eliminar Información del Horario Sede | DELETE | `/branches/{branchId}/schedules/{id}` | — | ❌ No implementado |
| HU34 | Activar el Horario Sede | PATCH | `/branches/{branchId}/schedules/{id}/activate` | — | ❌ No implementado |
| HU35 | Desactivar el Horario Sede | PATCH | `/branches/{branchId}/schedules/{id}/deactivate` | — | ❌ No implementado |

---

### Detalle Horario (Schedule Detail)

| HU | Nombre HU | Método | Endpoint Sugerido | Endpoint Actual | Estado |
|----|-----------|--------|-------------------|-----------------|--------|
| HU6 | Registrar Información del Detalle Horario | POST | `/schedules/{scheduleId}/details` | — | ❌ No implementado |
| HU7 | Modificar Información del Detalle Horario | PUT | `/schedules/{scheduleId}/details/{id}` | — | ❌ No implementado |
| HU8 | Eliminar Información del Detalle Horario | DELETE | `/schedules/{scheduleId}/details/{id}` | — | ❌ No implementado |
| HU9 | Consultar Información del Detalle Horario | GET | `/schedules/{scheduleId}/details` | — | ❌ No implementado |
| HU10 | Consultar Información del Día | GET | `/days` | — | ❌ No implementado |

---

### Excepción Horario (Schedule Exception)

| HU | Nombre HU | Método | Endpoint Sugerido | Endpoint Actual | Estado |
|----|-----------|--------|-------------------|-----------------|--------|
| HU20 | Registrar Información de la Excepción de Horario | POST | `/schedules/{scheduleId}/exceptions` | — | ❌ No implementado |
| HU21 | Modificar Información de la Excepción Horario | PUT | `/schedules/{scheduleId}/exceptions/{id}` | — | ❌ No implementado |
| HU22 | Eliminar Información de la Excepción Horario | DELETE | `/schedules/{scheduleId}/exceptions/{id}` | — | ❌ No implementado |
| HU23 | Consultar Información de la Excepción Horario | GET | `/schedules/{scheduleId}/exceptions` | — | ❌ No implementado |
| HU24 | Activar la Excepción Horario | PATCH | `/schedules/{scheduleId}/exceptions/{id}/activate` | — | ❌ No implementado |
| HU25 | Desactivar la Excepción Horario | PATCH | `/schedules/{scheduleId}/exceptions/{id}/deactivate` | — | ❌ No implementado |

---

### Ubicación (Location)

| HU | Nombre HU | Método | Endpoint Sugerido | Endpoint Actual | Estado |
|----|-----------|--------|-------------------|-----------------|--------|
| HU77 | Crear Información de la Ubicación | POST | `/locations` | — | ❌ No implementado |
| HU78 | Modificar Información de la Ubicación | PUT | `/locations/{id}` | — | ❌ No implementado |
| HU79 | Eliminar Información de la Ubicación | DELETE | `/locations/{id}` | — | ❌ No implementado |
| HU80 | Consultar Información de la Ubicación | GET | `/locations/{id}` | — | ❌ No implementado |
| HU4 | Consultar Información de la Ciudad | GET | `/cities` | — | ❌ No implementado |
| HU5 | Consultar Información del Departamento | GET | `/departments` | — | ❌ No implementado |

---

### Motocicleta (Motorcycle)

| HU | Nombre HU | Método | Endpoint Sugerido | Endpoint Actual | Estado |
|----|-----------|--------|-------------------|-----------------|--------|
| HU43 | Registrar Información de la Motocicleta | POST | `/motorcycles` | — | ❌ No implementado |
| HU44 | Modificar Información de la Motocicleta | PUT | `/motorcycles/{id}` | — | ❌ No implementado |
| HU45 | Eliminar Información de la Motocicleta | DELETE | `/motorcycles/{id}` | — | ❌ No implementado |
| HU46 | Consultar Información de la Motocicleta | GET | `/motorcycles/{id}` | — | ❌ No implementado |
| HU47 | Consultar Información General de la Motocicleta | GET | `/motorcycles` | — | ❌ No implementado |

---

### Imagen/Evidencia Motocicleta (Motorcycle Images)

| HU | Nombre HU | Método | Endpoint Sugerido | Endpoint Actual | Estado |
|----|-----------|--------|-------------------|-----------------|--------|
| HU36 | Agregar la imagen de la Motocicleta | POST | `/motorcycles/{id}/image` | — | ❌ No implementado |
| HU37 | Actualizar la Imagen de la Motocicleta | PUT | `/motorcycles/{id}/image` | — | ❌ No implementado |
| HU38 | Consultar la Imagen de la Motocicleta | GET | `/motorcycles/{id}/image` | — | ❌ No implementado |
| HU39 | Eliminar la Imagen de la Motocicleta | DELETE | `/motorcycles/{id}/image` | — | ❌ No implementado |
| HU16 | Cargar Evidencia Fotográfica de la Motocicleta | POST | `/motorcycles/{id}/evidence` | — | ❌ No implementado |
| HU17 | Actualizar Evidencia Fotográfica de la Motocicleta | PUT | `/motorcycles/{id}/evidence/{evidenceId}` | — | ❌ No implementado |
| HU18 | Consultar Evidencia Fotográfica de la Motocicleta | GET | `/motorcycles/{id}/evidence` | — | ❌ No implementado |
| HU19 | Eliminar Evidencia Fotográfica de la Motocicleta | DELETE | `/motorcycles/{id}/evidence/{evidenceId}` | — | ❌ No implementado |

---

### Catálogos Motocicleta (Motorcycle Catalogs)

| HU | Nombre HU | Método | Endpoint Sugerido | Endpoint Actual | Estado |
|----|-----------|--------|-------------------|-----------------|--------|
| HU1 | Consultar la Información de una Categoría Motocicleta | GET | `/motorcycle-categories` | — | ❌ No implementado |
| HU40 | Consultar Información de las líneas de una marca de motocicletas | GET | `/admin/brands/{brandId}/lines` | `/admin/brands/{brandId}/lines` | ✅ OK |
| HU41 | Consultar Información de la Línea Categoría | GET | `/motorcycle-line-categories` | — | ❌ No implementado |
| HU42 | Consultar Información de las Marcas de Motocicletas | GET | `/motorcycle-brands` | — | ❌ No implementado |
| HU49 | Consultar Información del Rango Cilindraje | GET | `/engine-displacement-ranges` | — | ❌ No implementado |
| HU50 | Consultar Información de la Referencia Motocicleta | GET | `/motorcycle-references` | — | ❌ No implementado |

---

### Diagnóstico (Diagnosis)

| HU | Nombre HU | Método | Endpoint Sugerido | Endpoint Actual | Estado |
|----|-----------|--------|-------------------|-----------------|--------|
| HU11 | Registrar Información del Diagnóstico | POST | `/diagnoses` | — | ❌ No implementado |
| HU12 | Modificar Información del Diagnóstico | PUT | `/diagnoses/{id}` | — | ❌ No implementado |
| HU13 | Eliminar Información del Diagnóstico | DELETE | `/diagnoses/{id}` | — | ❌ No implementado |
| HU14 | Consultar Información del Diagnóstico | GET | `/diagnoses/{id}` | — | ❌ No implementado |

---

### Calificación (Rating)

| HU | Nombre HU | Método | Endpoint Sugerido | Endpoint Actual | Estado |
|----|-----------|--------|-------------------|-----------------|--------|
| HU2 | Registrar Información de la Calificación | POST | `/ratings` | — | ❌ No implementado |
| HU3 | Consultar Información de la Calificación | GET | `/ratings/{id}` | — | ❌ No implementado |
| HU48 | Consultar Información del Rango Calificación | GET | `/rating-ranges` | — | ❌ No implementado |

---

## Resumen de Estado

| Estado | Cantidad | Porcentaje |
|--------|----------|------------|
| ✅ Implementado correctamente | 6 | ~7% |
| ⚠️ Implementado, renombrar | 0 | 0% |
| ❌ No implementado | ~82 | ~93% |

---

## Endpoints Actuales vs Sugeridos

### ~~Endpoints a Renombrar~~ (Completado ✅)

| Anterior | Actual | Estado |
|----------|--------|--------|
| `POST /accounts` | `POST /persons` | ✅ Completado |
| `GET /person/me` | `GET /persons/me` | ✅ Completado |

### Endpoints OK (mantener)

| Endpoint | HUs relacionadas |
|----------|------------------|
| `POST /auth/login` | HU58, HU88 |
| `POST /auth/password-reset` | HU56, HU86 |
| `POST /auth/resend-verification` | (Funcionalidad auxiliar) |
| `POST /auth/verify-email` | (Funcionalidad auxiliar) |
| `POST /auth/password/reset` | (Completar reset con token) |

---

## Estructura Sugerida para `/persons`

### Request Body para Crear Persona

```json
{
  "type": "branch_representative",  // o "motorcycle_user"
  "first_name": "Juan",
  "last_name": "Pérez",
  "second_last_name": "García",
  "email": "juan@example.com",
  "phone_number": "3001234567",
  "password": "secure123",
  // campos específicos según type...
}
```

### Response de `/persons/me`

```json
{
  "id": "abc123",
  "type": "branch_representative",
  "first_name": "Juan",
  "last_name": "Pérez",
  "email": "juan@example.com",
  // ...otros campos
  "_links": {
    "self": "/persons/me",
    "update": "/persons/me",
    "change_password": "/persons/me/password"
  }
}
```

---

## Glosario Español → Inglés

| Español | Inglés | Tabla BD | Campo Type |
|---------|--------|----------|------------|
| Representante Sede | Branch Representative | `persons` | `type='branch_representative'` |
| Usuario Motocicleta | Motorcycle User | `persons` | `type='motorcycle_user'` |
| Persona | Person | `persons` | — |
| Sede | Branch | `branches` | — |
| Tipo Sede | Branch Type | `branch_types` | — |
| Servicio Sede | Branch Service | `branch_services` | — |
| Tipo Servicio | Service Type | `service_types` | — |
| Servicio Realizado | Completed Service | `completed_services` | — |
| Estado Servicio | Service Status | `service_statuses` | — |
| Franquicia | Franchise | `franchises` | — |
| Horario Sede | Branch Schedule | `branch_schedules` | — |
| Detalle Horario | Schedule Detail | `schedule_details` | — |
| Excepción Horario | Schedule Exception | `schedule_exceptions` | — |
| Día | Day | `days` | — |
| Ubicación | Location | `locations` | — |
| Ciudad | City | `cities` | — |
| Departamento | Department | `departments` | — |
| Motocicleta | Motorcycle | `motorcycles` | — |
| Diagnóstico | Diagnosis | `diagnoses` | — |
| Calificación | Rating | `ratings` | — |
| Rango Calificación | Rating Range | `rating_ranges` | — |
| Rango Cilindraje | Engine Displacement Range | `engine_displacement_ranges` | — |
| Marca Motocicleta | Motorcycle Brand | `motorcycle_brands` | — |
| Línea Motocicleta | Motorcycle Line | `motorcycle_lines` | — |
| Categoría Motocicleta | Motorcycle Category | `motorcycle_categories` | — |

---

## Comandos para Release

```bash
# 1. Desde develop, crear rama release
git checkout develop && git pull origin develop
git checkout -b release/0.2.0

# 2. Push y crear PR a main (o merge directo)
git push origin release/0.2.0
# → PR en GitHub: release/0.2.0 → main

# 3. Después del merge, crear tag en main
git checkout main && git pull origin main
git tag -a v0.2.0 -m "Release v0.2.0"
git push origin v0.2.0

# 4. Sincronizar develop con main
git checkout develop && git pull origin develop
git merge main -m "Merge main into develop after release 0.2.0"
git push origin develop
```
