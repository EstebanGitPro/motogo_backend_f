# Estrategia de Despliegue — MotoGo Backend

> **Fecha:** Febrero 2026 | **Dokploy:** v0.27.1 | **Servidor:** Contabo VPS

---

## 1. Estrategias Evaluadas

### Opción A: VPS Manual (Nginx + systemd + binario)

El approach original documentado en `DEPLOY_GUIDE.md`:

```
Mac (go build) → scp binario → VPS → systemd → Nginx → Certbot SSL
```

| Ventaja | Desventaja |
|---------|------------|
| Control total del servidor | Mucha configuración manual |
| Sin overhead de Docker | Deploy manual con `scp` cada vez |
| | Nginx + Certbot que mantener |
| | Rollback manual (guardar binario anterior) |
| | No hay CI/CD nativo |

### Opción B: Docker Compose en VPS (manual)

Containerizar todo y correr con `docker compose up`:

```
GitHub → git pull en VPS → docker compose build → containers
```

| Ventaja | Desventaja |
|---------|------------|
| Reproducible | Requiere SSH al VPS para cada deploy |
| Rollback con image tags | Nginx/Certbot sigue siendo manual |
| | No hay interfaz visual |

### Opción C: Dokploy ✅ (elegida)

Plataforma de deploy self-hosted (panel de administración web para gestión de deploys):

```
GitHub Push → Dokploy detecta → Docker Build → Deploy automático → Traefik SSL
```

| Ventaja | Desventaja |
|---------|------------|
| Deploy automático desde GitHub | Consume ~500MB RAM extra en el VPS |
| SSL automático con Traefik | Dependencia de Dokploy como plataforma |
| UI web para logs, env vars, dominios | |
| Rollback con un click | |
| Múltiples apps en un solo servidor | |
| Healthchecks y restart automático | |

### ¿Por qué Dokploy?

1. **Elimina Nginx manual** — Traefik (integrado) maneja routing + SSL automáticamente
2. **Elimina systemd** — Dokploy maneja el lifecycle del container (restart, health)
3. **Elimina `scp`** — Dokploy conecta directo al repo de GitHub y buildea
4. **UI visual** — Logs, métricas, env vars, dominios, todo desde el navegador
5. **Multi-proyecto** — Un solo Dokploy puede deployar N aplicaciones sin conflicto

---

## 2. Archivos Necesarios (ya creados)

| Archivo | Propósito | Estado |
|---------|-----------|--------|
| `Dockerfile` | Multi-stage build de la API Go | ✅ Creado |
| `.dockerignore` | Reduce build context (~100MB menos) | ✅ Actualizado |
| `.env.example` | Referencia de variables de entorno | ✅ Ya existía |

---

## 3. Configurar GitHub para Dokploy

### 3.1 Crear GitHub App (una sola vez)

1. Ir a **Dokploy UI** → **Settings** → **Git Providers**
2. Click en **GitHub** → **Configure**
3. Dokploy te redirige a GitHub para instalar una **GitHub App**
4. Seleccionar la organización/cuenta (`EstebanGitPro`)
5. Dar acceso al repositorio `motogo_backend_f`
6. Confirmar → Dokploy queda conectado a GitHub

> [!NOTE]
> La GitHub App permite a Dokploy leer el código y recibir webhooks de push.
> Solo necesitas hacer esto **una vez** por cuenta de GitHub.

### 3.2 Dar acceso a repos adicionales

Si después quieres deployar otro repo (ej: `rb-go-backend`):

1. Ir a **GitHub** → **Settings** → **Applications** → **Dokploy**
2. Click en **Configure** → **Repository access**
3. Agregar el nuevo repo
4. Listo — aparecerá disponible en Dokploy

---

## 4. Configurar la App en Dokploy

### 4.1 Crear la aplicación

1. En Dokploy → **Projects** → **Create Project** (ej: "MotoGo")
2. Dentro del proyecto → **Create Service** → **Application**
3. Configurar:

| Campo | Valor |
|-------|-------|
| **Provider** | GitHub |
| **Repository** | `EstebanGitPro/motogo_backend_f` |
| **Branch** | `main` |
| **Build Type** | Dockerfile |
| **Dockerfile Path** | `Dockerfile` |
| **Docker Context** | `.` |

### 4.2 Configurar variables de entorno

En **Environment** tab, agregar todas las variables del `.env.example`:

```env
# ── Entorno ──────────────────────────────────
APP_ENV=production

# ── Base de Datos ────────────────────────────
DB_DRIVER=mysql
DB_HOST=<IP_O_HOST_MYSQL>
DB_PORT=3306
DB_USERNAME=motogo_app
DB_PASSWORD=<password_seguro>
DB_NAME=motogoDb
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=10
DB_CONN_MAX_LIFETIME=300
DB_CONN_MAX_IDLE_TIME=60

# ── Servidor ─────────────────────────────────
SERVER_PORT=8085
SERVER_HOST=0.0.0.0

# ── Resend (Email) ───────────────────────────
RESEND_API_KEY=re_xxxxx
RESEND_FROM_EMAIL=motogo@rbsuport.com

# ── Keycloak ─────────────────────────────────
KEYCLOAK_SERVER_URL=https://auth.tudominio.com
KEYCLOAK_REALM=motogo
KEYCLOAK_CLIENT_ID=stifler
KEYCLOAK_CLIENT_SECRET=<secret>
KEYCLOAK_ADMIN=motogo-admin
KEYCLOAK_ADMIN_PASSWORD=<password_seguro>

# ── Verificación ─────────────────────────────
VERIFICATION_BASE_URL=https://api.tudominio.com

# ── ID Encoder ───────────────────────────────
ID_ENCODER_SECRET=<generar_con_openssl_rand_-base64_32>
ID_ENCODER_MIN_LENGTH=36

# ── Firebase ─────────────────────────────────
FIREBASE_CREDENTIALS_PATH=config/serviceAccountKey.json

# ── Geocoding ────────────────────────────────
GEOCODING_PROVIDER=google
GEOCODING_API_KEY=<tu_google_key>
GEOCODING_BASE_URL=https://maps.googleapis.com/maps/api/geocode
GEOCODING_TIMEOUT=5
GEOCODING_COUNTRY_CODE=co
```

> [!WARNING]
> El `serviceAccountKey.json` de Firebase necesita estar dentro de la imagen o montado como volumen.
> Opción 1: Dejarlo en `config/` del repo (⚠️ asegurarse que `.gitignore` lo excluye).
> Opción 2: Usar **Advanced → Volumes** en Dokploy para montar el archivo desde el host.

### 4.3 Configurar dominio

1. En **Domains** tab → **Add Domain**
2. Poner: `api.tudominio.com`
3. Activar **HTTPS** → Dokploy genera el certificado SSL con Let's Encrypt automáticamente
4. **Puerto del container:** `8085`

### 4.4 Configurar Health Check (opcional pero recomendado)

En **Advanced** → **Health Check**:

| Campo | Valor |
|-------|-------|
| **Path** | `/health` |
| **Port** | `8085` |
| **Interval** | `30s` |

### 4.5 Deploy

1. Click en **Deploy** → Dokploy clona el repo, ejecuta `docker build`, y levanta el container
2. Ver logs en **Logs** tab para confirmar que arranca
3. Verificar: `https://api.tudominio.com/health`

### 4.6 Auto-deploy (opcional)

En **General** → habilitar **Auto Deploy** → cada push a `main` trigerea deploy automático.

---

## 5. ¿Necesito otro Dokploy para otra API?

**No.** Un solo Dokploy maneja múltiples aplicaciones.

### Cómo deployar una segunda API (ej: Rb-Go-Backend)

1. En el **mismo Dokploy** → crear un nuevo **Project** (ej: "Alertax")
2. Crear un **Service** → **Application** → conectar el otro repo
3. Configurar Dockerfile, env vars y dominio (ej: `api-alertax.tudominio.com`)
4. Deploy

```
Dokploy (1 instancia)
├── Project: MotoGo
│   └── App: motogo-api → api.tudominio.com
├── Project: Alertax
│   └── App: alertax-api → api-alertax.tudominio.com
└── Traefik (maneja routing + SSL de todas las apps)
```

> [!TIP]
> Cada app tiene su propio container, dominio, env vars y logs.
> Traefik rutea por dominio automáticamente — ambas apps pueden usar el mismo puerto interno (8085)
> porque Traefik distingue por hostname, no por puerto.

### Límites prácticos del VPS (Contabo 6 vCPU / 12 GB RAM)

| Servicio | RAM estimada |
|----------|-------------|
| Dokploy + Traefik | ~500MB |
| MySQL App | ~500MB |
| MySQL Keycloak | ~300MB |
| Keycloak | ~1.5GB |
| MotoGo API | ~50MB |
| Segunda API | ~50MB |
| **Total** | **~2.9GB de 12GB** |

Sobra capacidad de sobra para 5-10 aplicaciones Go.

---

## 6. Flujo de Deploy Resumido

```
Developer                    GitHub                      Dokploy (VPS)
    │                           │                            │
    ├── git push main ────────► │                            │
    │                           ├── webhook ───────────────► │
    │                           │                            ├── git clone
    │                           │                            ├── docker build (Dockerfile)
    │                           │                            ├── docker run (container)
    │                           │                            ├── Traefik SSL + routing
    │                           │                            └── ✅ Live en api.tudominio.com
    │                           │                            │
    └── Verificar en navegador ◄─────────────────────────────┘
```
