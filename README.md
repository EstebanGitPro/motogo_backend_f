# MotoGo Backend

[![Go Version](https://img.shields.io/badge/Go-1.25.2-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Version](https://img.shields.io/badge/version-0.14.0-blue)](CHANGELOG.md)
[![SonarCloud](https://img.shields.io/badge/SonarCloud-Analyzed-F3702A?logo=sonarcloud)](https://sonarcloud.io/project/overview?id=EstebanGitPro_motogo_backend_f)

> API RESTful del backend de **MotoGo** — plataforma móvil que ayuda a motociclistas a encontrar talleres y servicios confiables, filtrados por los vehículos registrados en su perfil, con búsqueda por cercanía y la posibilidad de solicitar cotizaciones o diagnósticos antes de desplazarse.

---

## 📖 Tabla de Contenidos

- [Visión](#-visión)
- [Funcionalidades Principales](#-funcionalidades-principales)
- [Stack Tecnológico](#-stack-tecnológico)
- [Arquitectura](#-arquitectura)
- [Estructura del Proyecto](#-estructura-del-proyecto)
- [Inicio Rápido](#-inicio-rápido)
- [Configuración](#️-configuración)
- [Comandos Disponibles](#-comandos-disponibles)
- [Observabilidad](#-observabilidad)
- [Documentación de la API](#-documentación-de-la-api)
- [Testing](#-testing)
- [Despliegue](#-despliegue)
- [Referencia de Puertos](#-referencia-de-puertos)
- [Contribuciones](#-contribuciones)
- [Seguridad](#️-seguridad)
- [Licencia](#-licencia)

---

## 🎯 Visión

Para los motociclistas que no cuentan con un mecanismo que permita encontrar establecimientos confiables que brinden servicios técnicos y/o venta de consumibles que se ajuste a sus necesidades, **MotoGo** es una aplicación móvil que permite a los usuarios ahorrar tanto su tiempo como su dinero en la búsqueda de los servicios para su motocicleta.

A diferencia de Yelp, 4WorldLover, Google Maps y Páginas Amarillas, nuestro producto permite acceder al **catálogo de servicios dependiendo de los vehículos que cada usuario tiene registrados en su perfil**, así como el acceso a la información de los servicios más cercanos al usuario, con la posibilidad de solicitar una cotización o un diagnóstico aproximado antes de dirigirse al lugar.

---

## ✨ Funcionalidades Principales

- **Gestión de Usuarios** — Registro, autenticación, verificación de email y flujos de restablecimiento de contraseña
- **Integración con Identity Provider** — Keycloak para gestión segura de identidad y acceso (OIDC/JWT)
- **Descubrimiento de Sedes y Talleres** — Búsqueda geoespacial con descubrimiento por cercanía en mapa
- **Gestión de Motocicletas** — Registro de vehículos, imágenes de perfil y galerías de evidencia
- **Catálogo de Servicios** — Catálogo global con asociaciones de servicios por sede
- **Diagnósticos y Servicios** — Flujo de diagnóstico de dos lados con integración WhatsApp
- **Sistema de Calificaciones** — Reseñas de servicios y rangos de calificación
- **Mensajería Dinámica** — Sistema de mensajes centralizado y manejado desde base de datos con caché en memoria
- **API HATEOAS** — Respuestas hypermedia (Richardson Maturity Model Level 3)
- **Observabilidad Integral** — Métricas con Prometheus, dashboards en Grafana, logs centralizados con Loki
- **Documentación de API** — Especificaciones OpenAPI/Swagger generadas automáticamente
- **Pruebas de Carga** — Scripts K6 para pruebas de rendimiento y resistencia

---

## 🛠 Stack Tecnológico

| Categoría | Tecnología |
| --- | --- |
| **Lenguaje** | Go 1.25.2 |
| **Web Framework** | [Gin](https://github.com/gin-gonic/gin) |
| **Base de Datos** | MySQL 8.0 |
| **Identity Management** | [Keycloak](https://www.keycloak.org/) (OIDC + JWKS validation) |
| **Almacenamiento de Archivos** | Firebase Storage |
| **Email** | [Resend](https://resend.com/) |
| **Notificaciones Push** | Firebase Cloud Messaging (FCM) |
| **Métricas** | Prometheus |
| **Dashboards** | Grafana |
| **Log Aggregation** | Loki + Promtail |
| **Structured Logging** | `log/slog` |
| **Documentación API** | Swagger (swaggo) |
| **Calidad de Código** | SonarCloud, golangci-lint, staticcheck |
| **Pruebas de Carga** | K6 |
| **Containerización** | Docker + Docker Compose |
| **Proxy HTTP y Gestión DNS** | Cloudflare |
| **Despliegue** | Dokploy en VPS Contabo |

---

## 🏗 Arquitectura

El backend de MotoGo sigue una **Clean / Hexagonal Architecture** (Ports & Adapters) con 4 capas:

```
┌─────────────────────────────────────────────────┐
│                   handlers/                      │  ← Presentation Layer (HTTP controllers)
├─────────────────────────────────────────────────┤
│                  middleware/                      │  ← Cross-cutting concerns (auth, CORS, logging)
├─────────────────────────────────────────────────┤
│                    core/                         │
│   ┌──────────────┐  ┌────────────────────────┐  │
│   │    ports/     │  │     interactor/        │  │  ← Domain Layer (business logic + interfaces)
│   │ (interfaces)  │  │  (use cases / services)│  │
│   └──────────────┘  └────────────────────────┘  │
├─────────────────────────────────────────────────┤
│                  platform/                       │  ← Infrastructure (DB, Firebase, Keycloak, etc.)
└─────────────────────────────────────────────────┘
```

**Principios clave:**

- Las dependencias apuntan **hacia adentro** — `platform` implementa las interfaces de `core/ports`
- La lógica de negocio en `core/interactor` es **agnóstica al framework**
- `handlers` solo orquesta el mapeo de request/response HTTP
- `middleware` provee autenticación, trazabilidad de requests y CORS

---

## 📁 Estructura del Proyecto

```
motogo_backend_f/
├── cmd/
│   ├── main.go                 # Punto de entrada de la aplicación
│   └── dependency/             # Wiring de inyección de dependencias
├── config/                     # Archivos de configuración (JSON) + config loader
├── core/
│   ├── ports/                  # Interfaces (contratos de repositorio y servicio)
│   └── interactor/             # Casos de uso / lógica de negocio
├── handlers/                   # Controladores HTTP (Gin handlers)
├── middleware/                  # Auth, CORS, request ID, métricas
├── server/                     # Bootstrap del servidor y registro de rutas
├── platform/                   # Implementaciones de infraestructura
│   ├── cache/                  # Caché en memoria (mensajes del sistema)
│   ├── constants/              # Constantes compartidas
│   ├── cookie/                 # Utilidades de manejo de cookies HTTP
│   ├── databases/              # Repositorios MySQL
│   ├── firebase/               # Adaptador de Firebase Storage
│   ├── geocoding/              # Servicio de geocodificación (Google / Mapbox)
│   ├── grafana/                # Dashboards y provisioning de Grafana
│   ├── identity_provider/      # Integración con Keycloak
│   ├── jwt/                    # Validación JWT / JWKS
│   ├── k6/                     # Escenarios de pruebas de carga K6
│   ├── log-rotator/            # Servicio de rotación de logs
│   ├── logger/                 # Logging estructurado (slog)
│   ├── loki/                   # Configuración de Loki
│   ├── prometheus/             # Configuración de scrape de Prometheus
│   ├── promtail/               # Configuración de envío de logs con Promtail
│   ├── schema/                 # Validación con JSON Schema
│   └── swaggo/                 # Documentación Swagger generada
├── mocks/                      # Implementaciones mock con testify
├── public/                     # Assets estáticos (templates de email)
├── tools/                      # Utilidades, temas de Keycloak, moderación
├── scripts/                    # Scripts de automatización (Docker, observabilidad)
├── tests/                      # Tests de integración / carga (K6)
├── deploy/                     # Configuraciones de despliegue a producción
│   ├── docker-compose.production.yml
│   ├── deploy.sh
│   ├── nginx-motogo.conf
│   └── .env.production.example
├── docs/                       # Swagger docs + colecciones Bruno API
├── docker-compose.mysql.yml    # MySQL local (app + keycloak)
├── docker-compose.keycloak.yml # Keycloak local
├── docker-compose.grafana.yml  # Stack de observabilidad
├── docker-compose.swagger-ui.yml
├── Containerfile               # Imagen custom de Keycloak (con tema MotoGo)
├── Makefile                    # Automatización de build, test y lint
├── .golangci.yml               # Configuración del linter
└── sonar-project.properties    # Integración con SonarCloud
```

---

## 🚀 Inicio Rápido

**Tiempo estimado:** 15–20 minutos

### Prerrequisitos

| Herramienta | Versión | Propósito |
| --- | --- | --- |
| **Go** | ≥ 1.25 | Runtime del lenguaje |
| **Docker** | ≥ 24.0 | Runtime de contenedores |
| **Docker Compose** | ≥ 2.20 | Orquestación multi-contenedor |
| **Make** | Cualquiera | Automatización de build (preinstalado en macOS/Linux) |
| **Git** | Cualquiera | Control de versiones |

### Paso 1: Clonar el Repositorio

```bash
git clone https://github.com/EstebanGitPro/motogo_backend_f.git
cd motogo_backend_f
```

### Paso 2: Iniciar las Bases de Datos MySQL

El proyecto usa **dos instancias MySQL**: una para la aplicación y otra para Keycloak.

```bash
docker compose -f docker-compose.mysql.yml up -d
```

**Verificación:**

```bash
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
# Esperado: motogo-mysql-app (healthy), motogo-mysql-keycloak (healthy)
```

> **Esperar** a que ambos contenedores reporten `(healthy)` antes de continuar (~30 segundos).

### Paso 3: Iniciar Keycloak

```bash
docker compose -f docker-compose.keycloak.yml up -d
```

**Verificación:** Abrir [http://localhost:8080](http://localhost:8080) — debería verse la consola de administración de Keycloak.

> Credenciales por defecto: `motogo-admin` / `12345678` (solo desarrollo local).

### Paso 4: Configurar la Aplicación

#### 4a. Variables de Entorno

El archivo `.env` en la raíz del proyecto configura todos los servicios externos y secretos. Copiar y editar:

```bash
cp .env.example .env
```

El `.env.example` contiene las siguientes secciones:

| Sección | Variables clave | Propósito |
| --- | --- | --- |
| **Database** | `DB_HOST`, `DB_PORT`, `DB_USERNAME`, `DB_PASSWORD`, `DB_NAME` | Conexión MySQL |
| **Server** | `SERVER_HOST`, `SERVER_PORT` | Binding del servidor API |
| **Resend** | `RESEND_API_KEY`, `RESEND_FROM_EMAIL` | Envío de emails (reset de contraseña, verificación) |
| **Keycloak** | `KEYCLOAK_SERVER_URL`, `KEYCLOAK_REALM`, `KEYCLOAK_CLIENT_SECRET` | Identity provider |
| **Firebase** | `FIREBASE_CREDENTIALS_PATH` | Ruta al service account de Firebase Storage |
| **ID Encoder** | `ID_ENCODER_SECRET`, `ID_ENCODER_MIN_LENGTH` | Codificación de HashID |
| **Geocoding** | `GEOCODING_PROVIDER`, `GEOCODING_API_KEY`, `GEOCODING_COUNTRY_CODE` | Resolución de direcciones (Google / Mapbox) |
| **Verification** | `VERIFICATION_BASE_URL` | URL de callback para verificación de email |

> **Tip:** Las variables de entorno tienen prioridad sobre el archivo JSON. Ver `config/config.go` para detalles.

#### 4b. Configuración por Archivo (JSON)

La aplicación también soporta configuración basada en archivos JSON por entorno:

| Entorno | Archivo | Trigger |
| --- | --- | --- |
| **Local** (por defecto) | `config/local-config.json` | `APP_ENV` sin definir o `local` |
| **Producción** | `config/prod-config.json` | `APP_ENV=production` |

Editar `config/local-config.json` para establecer valores no sobreescritos por `.env`.

#### 4c. Firebase Service Account

Colocar la clave del service account de Firebase en el directorio `config/`:

```bash
# Descargar desde Firebase Console → Project Settings → Service Accounts
cp ~/Downloads/your-firebase-adminsdk.json config/serviceAccountKey.json
```

> ⚠️ Este archivo está en `.gitignore` y **nunca** debe ser commiteado.

### Paso 5: Instalar Dependencias

```bash
go mod tidy
```

### Paso 6: Instalar Herramientas de Desarrollo (Opcional)

```bash
make install-tools
# Instala: staticcheck
# Para golangci-lint: brew install golangci-lint
```

### Paso 7: Ejecutar la Aplicación

```bash
make run
```

O directamente:

```bash
go run ./cmd/api
```

### ✅ Verificación

| Comprobación | Cómo verificar |
| --- | --- |
| **Servidor corriendo** | La consola muestra `Listening on 0.0.0.0:8085` |
| **Health endpoint** | `curl http://localhost:8085/health` retorna `200 OK` |
| **Metrics endpoint** | `curl http://localhost:8085/metrics` retorna datos de Prometheus |
| **Swagger UI** | Abrir [http://localhost:8085/swagger/index.html](http://localhost:8085/swagger/index.html) |

---

## ⚙️ Configuración

La aplicación soporta múltiples entornos mediante archivos de configuración JSON y sobreescrituras por variables de entorno.

**Precedencia:** Variables de entorno (`.env`) → Archivo JSON → Valores por defecto.

Todas las secciones de la configuración pueden sobreescribirse con variables de entorno. Las más comúnmente modificadas son:

```bash
# Base de datos
DB_HOST=localhost
DB_PORT=3309
DB_PASSWORD=your-password

# Servidor
SERVER_PORT=8085

# Keycloak
KEYCLOAK_SERVER_URL=https://auth.example.com
KEYCLOAK_CLIENT_SECRET=your-secret

# Servicios externos
RESEND_API_KEY=re_your_api_key
GEOCODING_API_KEY=your-google-maps-key
```

> Ver `.env.example` para la lista completa de variables de entorno disponibles.

---

## 📋 Comandos Disponibles

Toda la automatización del proyecto está disponible a través del `Makefile`:

| Comando | Descripción |
| --- | --- |
| `make help` | Muestra todos los targets disponibles |
| `make run` | Ejecuta la aplicación (`go run ./cmd/api`) |
| `make build` | Compila el binario en `bin/motogo-api` |
| `make test` | Ejecuta todos los tests con output verbose |
| `make test-short` | Ejecuta tests sin output verbose |
| `make coverage` | Genera reporte HTML de cobertura |
| `make coverage-check` | Verifica que la cobertura cumpla con el mínimo de 65% |
| `make lint` | Ejecuta `go vet` + `staticcheck` |
| `make lint-full` | Ejecuta `golangci-lint` con el conjunto completo de reglas |
| `make fmt` | Formatea todos los archivos fuente Go |
| `make tidy` | Ejecuta `go mod tidy` |
| `make clean` | Elimina artefactos de build y archivos de cobertura |
| `make pre-commit` | Ejecuta la suite completa de pre-commit (fmt + vet + staticcheck + test) |
| `make pre-push` | Ejecuta tests con enforcement del threshold de cobertura promedio |
| `make setup-hooks` | Instala el hook pre-commit de Git |
| `make setup-hooks-all` | Instala los hooks pre-commit y pre-push |
| `make install-tools` | Instala `staticcheck` y sugiere `golangci-lint` |

---

## 📊 Observabilidad

### Iniciar el Stack Completo de Observabilidad

```bash
docker compose -f docker-compose.grafana.yml up -d
```

Esto inicia **5 servicios**:

| Servicio | Puerto | Propósito |
| --- | --- | --- |
| **Grafana** | [localhost:3000](http://localhost:3000) | Dashboards y visualización |
| **Prometheus** | [localhost:9091](http://localhost:9091) | Scraping y almacenamiento de métricas |
| **Loki** | localhost:3100 | Agregación de logs |
| **Promtail** | — | Agente de envío de logs |
| **Log Rotator** | — | Rotación automática de logs (retención de 7 días) |

> Credenciales por defecto de Grafana: `admin` / `admin`

### Script de Inicio Rápido

Para la primera configuración con provisioning automático:

```bash
./scripts/quick-start-observability.sh
```

---

## 📚 Documentación de la API

### Swagger UI (Integrado)

La documentación de la API se sirve directamente desde la aplicación en ejecución:

```
http://localhost:8085/swagger/index.html
```

### Swagger UI (Contenedor Independiente)

Para un Swagger UI independiente con hot-reload:

```bash
docker compose -f docker-compose.swagger-ui.yml up -d
# Abrir http://localhost:3001
```

---

## 🧪 Testing

### Ejecutar Todos los Tests

```bash
make test
```

### Reporte de Cobertura

```bash
make coverage
# Genera: coverage.html (abrir en navegador)
# Imprime: porcentaje total de cobertura
```

### Threshold de Cobertura

El proyecto exige un **mínimo de 65% de cobertura**:

```bash
make coverage-check
# Falla el CI si la cobertura < 65%
```

### Pruebas de Carga (K6)

```bash
# Iniciar el stack K6 + InfluxDB
docker compose -f docker-compose.k6-influxdb.yml up -d

# Ejecutar pruebas de carga (scripts en el directorio tests/)
```

---

## 🚢 Despliegue

Los archivos de despliegue a producción están en el directorio `deploy/`:

| Archivo | Propósito |
| --- | --- |
| `docker-compose.production.yml` | MySQL de producción (app + keycloak) + Keycloak |
| `.env.production.example` | Plantilla para secrets de producción |
| `deploy.sh` | Script de despliegue automatizado |
| `nginx-motogo.conf` | Configuración de reverse proxy con Nginx |
| `motogo-api.service` | Archivo de servicio systemd para el binario Go |
| `ssh.sh` | Helper de conexión SSH |

### Despliegue Rápido

```bash
# 1. Copiar y configurar los secrets de producción
cp deploy/.env.production.example deploy/.env.production

# 2. Generar contraseñas seguras
openssl rand -base64 24  # Usar para cada campo de contraseña

# 3. Ejecutar el script de despliegue
./deploy/deploy.sh
```

> Ver `deploy/.env.production.example` para todas las variables requeridas.

---

## 🔌 Referencia de Puertos

| Puerto | Servicio | Entorno |
| --- | --- | --- |
| `8085` | MotoGo API | Local + Producción |
| `3306` | MySQL (Aplicación) | Local |
| `3309` | MySQL (Keycloak) | Local |
| `8080` | Keycloak (HTTP) | Local + Producción |
| `8443` | Keycloak (HTTPS) | Producción |
| `9000` | Keycloak (Health/Metrics) | Local + Producción |
| `3000` | Grafana | Local |
| `9091` | Prometheus | Local |
| `3100` | Loki | Local |
| `3001` | Swagger UI (standalone) | Local |

---

## 🤝 Contribuciones

¡Las contribuciones son bienvenidas! Por favor leer [CONTRIBUTING.md](CONTRIBUTING.md) para detalles sobre:

- Estrategia de branching con Git Flow
- Formato de Conventional Commits
- Estándares de calidad de código (golangci-lint, pre-commit hooks)
- Proceso de Pull Requests

---

## 🛡️ Seguridad

### Nota sobre Credenciales Incluidas

> **Este repositorio es un proyecto de grado académico.** Las credenciales incluidas en `config/local-config.json` son exclusivamente para el entorno de desarrollo local y se incluyen intencionalmente para facilitar la ejecución y evaluación del proyecto sin configuración adicional.

En un entorno productivo, **todas las credenciales se inyectan mediante variables de entorno** (archivo `.env` o secretos del proveedor de hosting), nunca desde archivos de configuración commiteados. El archivo `config/prod-config.json` sirve como plantilla con placeholders (`CAMBIAR_POR_...`) y no contiene valores reales.

| Archivo | ¿Commiteado? | Propósito |
| --- | :---: | --- |
| `config/local-config.json` | ✅ Sí | Configuración local lista para ejecutar |
| `config/prod-config.json` | ❌ No | Plantilla de producción (solo placeholders) |
| `.env` | ❌ No | Variables de entorno con secrets reales |
| `serviceAccountKey.json` | ❌ No | Credenciales de Firebase (nunca commiteado) |
| `deploy/.env.production` | ❌ No | Secrets de producción |

### Estrategia de Configuración

```
Prioridad de configuración (de mayor a menor):

  Variables de entorno (.env)  →  Archivo JSON (local/prod-config.json)  →  Valores por defecto
```

Para más detalles, ver la sección [Configuración](#️-configuración) y el archivo `config/config.go`.

---

## 📄 Licencia

Este proyecto está licenciado bajo la [Licencia MIT](LICENSE).

---

<p align="center">
  <sub>Hecho con ❤️ por Esteban Agudelo
</sub>
</p>
