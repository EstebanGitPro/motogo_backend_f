# MotoGo — Paquete de Instalación para Jurados

> Todo lo necesario para probar la aplicación MotoGo.

**ENTORNO CONTROLADO DE PRESENTACION**

Este paquete esta disenado exclusivamente para la evaluacion academica del proyecto de grado. El servidor al que se conecta la app es un entorno de presentacion con datos de prueba, no un entorno de produccion real. Los archivos de configuracion incluidos son templates de referencia y **no contienen contrasenas ni secrets reales**.

---

## Opcion Rapida: Instalar el APK

El archivo `MotoGo.apk` se encuentra en la carpeta `frontend/`. A continuacion se describen las formas de transferirlo e instalarlo.

### Como transferir el APK al celular

**Opcion A: Google Drive (sin cable)**

1. Desde el computador, subir el archivo `MotoGo.apk` a Google Drive.
2. En el celular Android, abrir la app de Google Drive.
3. Buscar el archivo y descargarlo.
4. Tocar el archivo descargado para iniciar la instalacion.

**Opcion B: Cable USB**

1. Conectar el celular al computador con un cable USB.
2. En el celular, seleccionar el modo "Transferencia de archivos" cuando aparezca el dialogo USB.
3. En Mac: descargar e instalar [Android File Transfer](https://www.android.com/filetransfer/) (macOS no monta dispositivos Android nativamente). En Windows: el dispositivo aparece directamente en el Explorador de archivos.
4. Copiar `MotoGo.apk` a la carpeta `Download` del celular.
5. En el celular, abrir el administrador de archivos y tocar el APK para instalarlo.

**Opcion C: WhatsApp**

1. Enviarse el archivo `MotoGo.apk` a uno mismo (chat personal o grupo individual).
2. En el celular, descargar el archivo desde el chat.
3. Tocar el archivo descargado para instalarlo.

**Opcion D: Emulador de Android Studio (sin celular fisico)**

1. Abrir Android Studio y acceder a Device Manager.
2. Crear un dispositivo virtual (ej: Pixel 6, API 33+) o usar uno existente.
3. Iniciar el emulador con el boton de play.
4. Arrastrar el archivo `MotoGo.apk` sobre la ventana del emulador. La app se instala automaticamente.

### Instalacion en el celular

1. Si el celular muestra un aviso de seguridad, ir a **Configuracion > Seguridad > Instalar apps desconocidas** y activar el permiso para la app con la que se abrio el APK (Archivos, Drive, Chrome, etc.).
2. Tocar "Instalar" cuando aparezca el dialogo.
3. Al abrir la app, esta ya viene configurada para conectarse al servidor de presentacion.

## Opción Completa: Ejecutar desde código fuente

El código fuente se encuentra en los siguientes repositorios:

- **Backend (Go):** <https://github.com/EstebanGitPro/motogo_backend_f>
- **Frontend (Flutter):** <https://github.com/EstebanGitPro/motogo_frontend>

### Prerrequisitos

| Herramienta    | Versión | Verificar con         |
| -------------- | -------- | --------------------- |
| Docker Desktop | Última  | `docker --version`  |
| Go             | 1.25+    | `go version`        |
| Flutter SDK    | 3.8.1+   | `flutter --version` |
| Android SDK    | API 33+  | Android Studio > SDK  |

### Paso 1: Levantar la base de datos y Keycloak

```bash
cd backend/
docker compose -f docker-compose.mysql.yml up -d
docker compose -f docker-compose.keycloak.yml up -d
```

### Paso 2: Configurar el backend

```bash
cp backend/.env.example .env
# Editar .env con las credenciales de desarrollo
go run ./cmd
```

### Paso 3: Ejecutar la app Flutter

```bash
cd <ruta-del-frontend>
flutter pub get
flutter run --dart-define=BASE_URL=http://10.0.2.2:8085/motogo/api/v1
```

> `10.0.2.2` es cómo el emulador de Android accede al localhost de la máquina host.
> Para un dispositivo físico, usar la IP local de la máquina (ej: `192.168.1.X`).

---

## Estructura de la carpeta

```
Anexos_para_jurados/
├── README.md                           ← Este archivo
├── Guia_Instalacion_APK_MotoGo.xlsx    ← Guía detallada paso a paso
│
├── frontend/
│   └── MotoGo.apk                     ← App Android (servidor en la nube)
│
└── backend/
    ├── Dockerfile                      ← Imagen del backend (Go, distroless)
    ├── Containerfile                   ← Imagen de Keycloak
    ├── docker-compose.mysql.yml        ← MySQL local
    ├── docker-compose.keycloak.yml     ← Keycloak local
    ├── docker-compose.production.yml   ← Stack producción (referencia)
    ├── docker-compose.grafana.yml      ← Observabilidad (Grafana + Loki)
    ├── .env.example                    ← Variables de entorno backend
    └── .env.production.example         ← Variables producción (referencia)
```

---

## Notas importantes

- El APK viene firmado con clave de debug — es solo para evaluación
- Los archivos `.env.production.example` son **solo de referencia**, no se usan directamente
- El servidor de **presentación** está en: `https://api.rbsuport.com`
