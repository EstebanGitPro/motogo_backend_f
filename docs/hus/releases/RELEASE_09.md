# Release 9: Motocicleta - CRUD

**Tag Git**: `v0.8.0`
**Estado**: 🚧 En progreso
**Total Horas**: 72h

## Historias de Usuario

| HU   | Descripción                             | Horas | Endpoint                                  | Estado |
| ---- | ---------------------------------------- | ----- | ----------------------------------------- | ------ |
| HU43 | Registrar Información de la Motocicleta | 32h   | `POST /motorcycles`                     | ✅     |
| HU46 | Listar Motocicletas del Propietario      | 4h    | `GET /motorcycles`                      | ✅     |
| HU47 | Consultar Motocicleta por Placa          | 8h    | `GET /motorcycles/lookup?plate={placa}` | ✅     |
| HU44 | Modificar Información de la Motocicleta | 16h   | `PUT /motorcycles/{id}`                 | ✅     |
| HU45 | Eliminar Información de la Motocicleta  | 8h    | `DELETE /motorcycles/{id}`              | ✅     |

## Catálogos Incluidos

| HU   | Descripción                     | Horas | Endpoint                         | Estado |
| ---- | -------------------------------- | ----- | -------------------------------- | ------ |
| HU50 | Consultar Referencia Motocicleta | 2h    | `GET /motorcycle-references`     | ✅     |
| HU40 | Consultar Líneas de Marca       | 2h    | `GET /admin/brands/{id}/lines`   | ✅     |

---

## HU46: Listar Motocicletas del Propietario

> [!NOTE]
> **Uso**: Se muestra en el **Home del usuario** después de iniciar sesión.

### Control de Acceso HU46

- Solo usuarios autenticados pueden consultar.
- Solo retorna las motocicletas donde `owner_id` coincide con el usuario logueado.

### Endpoint HU46

```http
GET /motorcycles
Authorization: Bearer {token}
```

### Respuesta HU46

Lista de motocicletas del propietario con información completa (referencia, marca, modelo, etc.).

---

## HU47: Consultar Motocicleta por Placa (Información Pública)

> [!IMPORTANT]
> **Uso**: Para compartir información con el **taller**. El propietario proporciona el número de placa.

### Control de Acceso HU47

- Solo usuarios con rol **Representative** (taller/supervisor) pueden consultar.
- Búsqueda por **placa exacta**.

### Endpoint HU47

```http
GET /motorcycles/lookup?plate={placa}
Authorization: Bearer {token}
```

### Información Retornada HU47

Respuesta simplificada para talleres (sin datos privados del propietario):

- Placa, Año, Kilometraje actual
- Referencia (Marca, Modelo, Categoría, Cilindraje)

> [!IMPORTANT]
> **Datos excluidos**: `owner_notes` (privado del propietario), `brand_id`, `reference.id`

---

## Reglas de Negocio - Ownership (Aplica a HU44, HU45)

> [!IMPORTANT]
> **Solo el propietario puede modificar/eliminar su motocicleta.**

### Validación de Propiedad

| Escenario                             | Respuesta HTTP   | Código Mensaje                 |
| ------------------------------------- | ---------------- | ------------------------------- |
| Usuario no autenticado                | 401 Unauthorized | `GEN_AUTH_ERR_00002`          |
| Usuario autenticado pero NO es dueño | 404 Not Found    | `MOD_MOT_NOT_FOUND_ERR_00001` |
| Usuario autenticado Y es dueño       | 200 OK           | (Éxito)                        |

**Razón de 404 en lugar de 403**: Se retorna 404 para no revelar la existencia de motocicletas de otros usuarios (Security by Obscurity).

### Links HATEOAS

Solo para el propietario:

- `self` → `GET /motorcycles/{id}`
- `update` → `PUT /motorcycles/{id}`
- `delete` → `DELETE /motorcycles/{id}`
- `list` → `GET /motorcycles`

---

## Dependencias

- ✅ Release 11: Usuario (motocicleta pertenece a usuario)

---

## HU40: Consultar Líneas de una Marca (Admin)

> [!IMPORTANT]
> **Acceso**: Solo usuarios con rol **ADMIN**.

### Endpoint HU40

```http
GET /admin/brands/{brandId}/lines
Authorization: Bearer {admin_token}
```

### Respuesta HU40

```json
{
  "success": true,
  "code": "MOD_MOT_BRAND_LINES_EXI_00001",
  "message": "Las líneas de motocicletas se obtuvieron correctamente.",
  "data": {
    "lines": [
      {
        "brand_name": "AKT",
        "model": "CR4 125"
      },
      {
        "brand_name": "AKT",
        "model": "NKD 125"
      }
    ],
    "_links": [...]
  }
}
```

### Archivos Implementados

| Capa | Archivo |
| ---- | ------- |
| Repository | `get_references_by_brand.go` |
| Interactor | `interactor_motorcycle.go` |
| Handler | `motorcycle_controller.go` → `GetBrandLines()` |
| DTO | `motorcycle.go` → `BrandLineItem` |
| HATEOAS | `hateoas.go` → `BuildBrandLinesLinks()` |
| Routing | `server.go` → `/admin/brands/:brandId/lines` |
| SQL | `025_add_brand_lines_messages.sql` |
| Bruno | `docs/bruno_admin/` |
