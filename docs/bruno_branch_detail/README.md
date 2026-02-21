# Bruno Collection: Branch Detail View

Esta colección prueba los endpoints necesarios para mostrar el detalle de una sede/taller en la app móvil.

## Prerequisitos

1. **Backend corriendo** en `http://localhost:8085`
2. **Sedes existentes** con coordenadas en la base de datos
3. **Bruno instalado** - [https://www.usebruno.com/](https://www.usebruno.com/)

## Configuración

El archivo `environments/local.bru` contiene:
- `base_url`: URL base del API
- `lat`, `lng`, `radius`: Coordenadas para buscar sedes cercanas (por defecto: Rionegro, Antioquia)
- `branch_id`: ID de sede de prueba (se actualiza automáticamente)

## Orden de Ejecución

```
╔═══════════════════════════════════════════════════════════════╗
║  PASO 1: LOGIN                                                ║
║  00_login.bru → Obtiene token, lo guarda en access_token      ║
╠═══════════════════════════════════════════════════════════════╣
║  PASO 2: BUSCAR SEDES CERCANAS                                ║
║  04_get_nearby_branches.bru → Lista sedes, guarda branch_id   ║
╠═══════════════════════════════════════════════════════════════╣
║  PASO 3-5: CONSULTAR DETALLE (en cualquier orden)             ║
║  01_get_branch.bru → Info básica de la sede                   ║
║  02_get_branch_services.bru → Servicios de la sede            ║
║  03_get_branch_schedules.bru → Horarios de la sede            ║
╚═══════════════════════════════════════════════════════════════╝
```

## Requests

### 00 - Login (Get Token)
- **POST** `/auth/login`
- Credenciales: `negro@yopmail.com` / `Secret123*` (rol USER)
- **Guarda:** `access_token` para los siguientes requests

### 04 - Get Nearby Branches
- **GET** `/branches/nearby?lat={{lat}}&lng={{lng}}&radius={{radius}}`
- Busca sedes cercanas a las coordenadas del environment
- **Guarda:** `branch_id` del primer resultado

### 01 - Get Branch Detail
- **GET** `/branches/{{branch_id}}`
- Retorna: nombre, tipo, estado, teléfono, ubicación, marcas

### 02 - Get Branch Services
- **GET** `/branches/{{branch_id}}/services`
- Retorna: lista de servicios asociados a la sede

### 03 - Get Branch Schedule Details
- **GET** `/branches/{{branch_id}}/schedules/details`
- Retorna: franjas horarias por día de la semana

## Campos importantes (UI del mapa)

| Campo | Descripción |
|-------|-------------|
| `contact_phone` | Teléfono del representante (para botón "Llamar") |
| `establishment_type_label` | Tipo de sede en español (Taller, Tienda, etc.) |
| `distance_km` | Distancia en km desde las coordenadas de búsqueda |
| `_links` | Links HATEOAS para navegar entre recursos |

## Notas

- Las coordenadas por defecto (6.14, -75.37) corresponden a Rionegro, Antioquia
- Si no aparecen sedes en nearby, verificar que las sedes tengan lat/lng configurados
- El endpoint de horarios ahora es accesible por rol USER (no solo REPRESENTANTE)
