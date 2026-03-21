# Arquitectura de Base de Datos

## Visión General

El proyecto utiliza un único servidor MySQL ejecutado en Docker:

```
Container: motogo-mysql-keycloak
```

Dentro del servidor existen dos bases separadas:

| Base de Datos | Propósito                                 |
| ------------- | ------------------------------------------ |
| `motogoDb`  | Base de datos de negocio                   |
| `keydb`     | Base de datos de autenticación (Keycloak) |

Esta separación permite:

- **Aislamiento lógico** de responsabilidades
- **Organización** clara entre negocio y autenticación
- **Independencia de backups** por dominio
- **Escalabilidad futura** (migrar Keycloak a otro servidor sin afectar negocio)

---

## Diagrama de Arquitectura

```
┌─────────────────────────────────────────┐
│         Docker Container                │
│     motogo-mysql-keycloak               │
│                                         │
│   ┌─────────────┐  ┌─────────────┐     │
│   │  motogoDb   │  │   keydb     │     │
│   │  (negocio)  │  │ (Keycloak)  │     │
│   │             │  │             │     │
│   │ • persons   │  │ • realm     │     │
│   │ • branches  │  │ • users     │     │
│   │ • motors    │  │ • clients   │     │
│   │ • services  │  │ • sessions  │     │
│   │ • ratings   │  │ • roles     │     │
│   │ • ...       │  │ • ...       │     │
│   └─────────────┘  └─────────────┘     │
│                                         │
└─────────────────────────────────────────┘
```

---

## Exportación de Estructura (sin datos)

Para exportar la estructura DDL de ambas bases:

```bash
chmod +x database-export.sh
./database-export.sh
```

Esto genera:

| Archivo                 | Contenido                              |
| ----------------------- | -------------------------------------- |
| `schema-motogo.sql`   | Estructura de tablas de negocio (DDL)  |
| `schema-keycloak.sql` | Estructura de tablas de Keycloak (DDL) |

> [!IMPORTANT]
> Ambos archivos contienen **únicamente la definición de tablas** (CREATE TABLE), sin datos. Son ideales para versionar en Git.

---

## Tablas de Negocio (`motogoDb`)

### Base del Sistema

| Tabla               | Descripción                                       |
| ------------------- | -------------------------------------------------- |
| `persons`         | Usuarios del sistema (USER, REPRESENTATIVE, ADMIN) |
| `system_messages` | Mensajes dinámicos del sistema                    |

### Catálogos Geográficos

| Tabla           | Descripción              |
| --------------- | ------------------------- |
| `departments` | Departamentos de Colombia |
| `cities`      | Ciudades por departamento |

### Establecimientos

| Tabla                | Descripción                    |
| -------------------- | ------------------------------- |
| `franchises`       | Franquicias                     |
| `branches`         | Sedes (talleres y tiendas)      |
| `branch_brands`    | Marcas atendidas por sede       |
| `locations`        | Ubicación geográfica de sedes |
| `branch_schedules` | Configuración de horarios      |
| `schedule_details` | Franjas horarias y excepciones  |

### Servicios

| Tabla               | Descripción                 |
| ------------------- | ---------------------------- |
| `services`        | Catálogo de servicios       |
| `branch_services` | Servicios ofrecidos por sede |

### Motocicletas

| Tabla                     | Descripción             |
| ------------------------- | ------------------------ |
| `motorcycle_references` | Catálogo de referencias |
| `motorcycles`           | Motocicletas registradas |
| `motorcycle_evidence`   | Fotos/evidencias         |

### Diagnósticos y Servicios Realizados

| Tabla                                 | Descripción                          |
| ------------------------------------- | ------------------------------------- |
| `diagnostics`                       | Diagnósticos y cotizaciones          |
| `motorcycle_diagnostic_permissions` | Permisos de diagnóstico por sede     |
| `completed_services`                | Servicios realizados                  |
| `completed_service_items`           | Ítems de servicio con calificaciones |
| `service_status_transitions`        | Trazabilidad de cambios de estado     |

## Estrategia de Integridad Referencial

| Relación                                 | ON DELETE    | Razón                               |
| ----------------------------------------- | ------------ | ------------------------------------ |
| `locations` → `branches`             | `CASCADE`  | La ubicación muere con la sede      |
| `locations` → `cities`               | `RESTRICT` | Protege catálogos geográficos      |
| `motorcycles` → `persons`            | `RESTRICT` | No se puede borrar persona con motos |
| `completed_services` → `branches`    | `RESTRICT` | Protege historial de servicios       |
| `completed_services` → `motorcycles` | `CASCADE`  | Cascada con moto (usa soft-delete)   |
| `branch_brands` → `brands`           | `RESTRICT` | Protege catálogo de marcas          |

> [!NOTE]
> Las motocicletas usan **soft-delete** (`deleted_at`), por lo que el `CASCADE` nunca se dispara en producción. Las calificaciones se preservan siempre.
