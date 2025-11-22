# Richardson Maturity Model & HATEOAS - MotoGo Backend

## Implementación de REST Nivel 3

Este documento describe cómo los servicios REST de MotoGo Backend siguen el **Richardson Maturity Model** y las mejores prácticas de diseño RESTful.

---

## Richardson Maturity Model

### Nivel 0: The Swamp of POX

- ❌ No implementado - Superado

### Nivel 1: Resources

- ✅ Recursos con URIs únicas: `/motogo/api/v1/accounts`
- ✅ Cada recurso tiene su propia URI identificable

### Nivel 2: HTTP Verbs

- ✅ Uso correcto de verbos HTTP:
  - `POST /accounts` - Crear nuevo recurso
  - `GET /accounts/:id` - Obtener recurso existente
  - `PUT /accounts/:id` - Actualizar recurso (futuro)
  - `DELETE /accounts/:id` - Eliminar recurso (futuro)
- ✅ Códigos de estado HTTP apropiados:
  - `201 Created` - Recurso creado exitosamente
  - `200 OK` - Recurso obtenido exitosamente
  - `404 Not Found` - Recurso no encontrado
  - `400 Bad Request` - Solicitud inválida

### Nivel 3: Hypermedia Controls (HATEOAS)

- ✅ **HATEOAS** (Hypermedia As The Engine Of Application State)
- ✅ Respuestas incluyen hipervínculos navegables
- ✅ Cliente puede descubrir acciones disponibles a través de los links

---

## HATEOAS Implementation

### Estructura de Links

Cada recurso incluye un array `_links` con hipervínculos relacionados:

```json
{
  "_links": [
    {
      "href": "http://localhost:8080/motogo/api/v1/accounts/123",
      "rel": "self",
      "method": "GET"
    },
    {
      "href": "http://localhost:8080/motogo/api/v1/accounts/123",
      "rel": "update",
      "method": "PUT"
    },
    {
      "href": "http://localhost:8080/motogo/api/v1/accounts/123",
      "rel": "delete",
      "method": "DELETE"
    },
    {
      "href": "http://localhost:8080/motogo/api/v1/accounts",
      "rel": "collection",
      "method": "GET"
    }
  ]
}
```

### Relaciones (rel)

- **self**: El recurso mismo
- **update**: Endpoint para actualizar el recurso
- **delete**: Endpoint para eliminar el recurso
- **collection**: Colección a la que pertenece el recurso

---

## Endpoints Implementados

### POST /motogo/api/v1/accounts

**Crear nueva cuenta de usuario**

#### Request

```json
{
  "identity_number": "123456789",
  "first_name": "Juan",
  "last_name": "Pérez",
  "second_last_name": "García",
  "email": "juan@example.com",
  "phone_number": "+57300123456",
  "password": "SecurePass123!",
  "role": "user"
}
```

#### Response

**Status:** `201 Created`  
**Header:** `Location: http://localhost:8080/motogo/api/v1/accounts/{id}`

```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "identity_number": "123456789",
    "first_name": "Juan",
    "last_name": "Pérez",
    "second_last_name": "García",
    "email": "juan@example.com",
    "phone_number": "+57300123456",
    "role": "user",
    "keycloak_user_id": "keycloak-uuid-here",
    "_links": [
      {
        "href": "http://localhost:8080/motogo/api/v1/accounts/550e8400-e29b-41d4-a716-446655440000",
        "rel": "self",
        "method": "GET"
      },
      {
        "href": "http://localhost:8080/motogo/api/v1/accounts/550e8400-e29b-41d4-a716-446655440000",
        "rel": "update",
        "method": "PUT"
      },
      {
        "href": "http://localhost:8080/motogo/api/v1/accounts/550e8400-e29b-41d4-a716-446655440000",
        "rel": "delete",
        "method": "DELETE"
      },
      {
        "href": "http://localhost:8080/motogo/api/v1/accounts",
        "rel": "collection",
        "method": "GET"
      }
    ]
  },
  "message": "Usuario registrado exitosamente",
  "_links": [...]
}
```

### GET /motogo/api/v1/accounts/:id

**Locate: Obtener cuenta por ID**

Este es el endpoint referenciado en el `Location` header del POST.

#### Response

**Status:** `200 OK`

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "identity_number": "123456789",
  "first_name": "Juan",
  "last_name": "Pérez",
  "second_last_name": "García",
  "email": "juan@example.com",
  "phone_number": "+57300123456",
  "role": "user",
  "keycloak_user_id": "keycloak-uuid-here",
  "_links": [
    {
      "href": "http://localhost:8080/motogo/api/v1/accounts/550e8400-e29b-41d4-a716-446655440000",
      "rel": "self",
      "method": "GET"
    },
    {
      "href": "http://localhost:8080/motogo/api/v1/accounts/550e8400-e29b-41d4-a716-446655440000",
      "rel": "update",
      "method": "PUT"
    },
    {
      "href": "http://localhost:8080/motogo/api/v1/accounts/550e8400-e29b-41d4-a716-446655440000",
      "rel": "delete",
      "method": "DELETE"
    },
    {
      "href": "http://localhost:8080/motogo/api/v1/accounts",
      "rel": "collection",
      "method": "GET"
    }
  ]
}
```

---

## Best Practices Implementadas

### 1. Location Header (RFC 7231)

✅ **POST** devuelve header `Location` con la URI del recurso creado

### 2. Status Codes Apropiados

✅ **201 Created** para recursos creados  
✅ **200 OK** para recursos obtenidos  
✅ **404 Not Found** para recursos no encontrados  
✅ **400 Bad Request** para solicitudes inválidas

### 3. HATEOAS

✅ Todas las respuestas incluyen `_links` con hipervínculos navegables  
✅ Cliente puede descubrir funcionalidad sin documentación externa

### 4. ID Local + ID Keycloak

✅ **ID Local** (`id`): Identificador único en la base de datos del negocio  
✅ **ID Keycloak** (`keycloak_user_id`): Identificador en el sistema de autenticación

Según las recomendaciones del profesor:
- El ID local sirve para las operaciones del negocio
- Con el ID local se puede traer el ID de Keycloak para autenticación
- La autenticación se maneja con Keycloak
- La lógica de negocio usa el ID local

### 5. Servicios de Negocio

✅ Los endpoints representan operaciones de negocio:
- **POST /accounts** - Registrar nuevo usuario en el sistema
- **GET /accounts/:id** - Obtener información de cuenta (Locate)

---

## Arquitectura

### Clean Architecture + Hexagonal

```
handlers/ (Adapters)
  ├── hateoas.go          - Estructuras HATEOAS
  ├── person.go           - DTOs de request/response
  └── person_controller.go - Controladores REST

core/
  ├── interactor/         - Orquestación de casos de uso
  ├── ports/              - Interfaces (input/output)
  └── services/           - Lógica de negocio

platform/
  ├── databases/          - Implementación BD
  └── identity_provider/  - Implementación Keycloak
```

---

## Próximos Pasos

### Endpoints Futuros

- [ ] `PUT /accounts/:id` - Actualizar cuenta
- [ ] `DELETE /accounts/:id` - Eliminar cuenta
- [ ] `GET /accounts` - Listar cuentas (con paginación)
- [ ] `POST /auth/login` - Autenticación (devuelve tokens Keycloak)
- [ ] `POST /auth/logout` - Cerrar sesión
- [ ] `POST /auth/refresh` - Refrescar token

### Mejoras HATEOAS

- [ ] Links condicionales según permisos del usuario
- [ ] Más relaciones (prev, next para paginación)
- [ ] Templates de URI (RFC 6570)

---

## Referencias

- [Richardson Maturity Model](https://martinfowler.com/articles/richardsonMaturityModel.html)
- [RFC 7231 - HTTP/1.1 Semantics and Content](https://tools.ietf.org/html/rfc7231)
- [HATEOAS](https://restfulapi.net/hateoas/)
- [REST API Best Practices](https://restfulapi.net/)