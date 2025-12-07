# Swagger/OpenAPI Generated Files

Esta carpeta contiene los archivos generados automáticamente por `swag` CLI que son necesarios para que la aplicación ejecute la documentación Swagger/OpenAPI.

## 📁 Archivos

### Archivos Generados
- **`docs.go`** - Paquete Go con la documentación embebida (importado por `server/server.go`)
- **`swagger.json`** - Especificación OpenAPI 3.0 en formato JSON
- **`swagger.yaml`** - Especificación OpenAPI 3.0 en formato YAML

### Archivos de Infraestructura
- **`Dockerfile.swag`** - Imagen Docker personalizada con Go + swag (para generación con Docker)
- **`README.md`** - Este archivo

## 🔄 Generación

Estos archivos se regeneran automáticamente ejecutando:

```bash
./generate-swagger.sh
```

El script intenta 3 métodos en orden:
1. **swag en PATH** - Si tienes swag instalado globalmente
2. **swag en GOPATH** - Si lo instalaste con `go install`
3. **Docker personalizado** - Usa `Dockerfile.swag` con Go incluido

**⚠️ IMPORTANTE**: No edites los archivos generados manualmente. Cualquier cambio se sobrescribirá en la próxima generación.

## 📖 Documentación para Humanos

Para guías y documentación de cómo usar Swagger, consulta:

```
docs/swaggo/
├── README.md                    # Guía completa
├── QUICK_START.md              # Inicio rápido
├── EJEMPLOS_ANOTACIONES.md     # Ejemplos prácticos
└── RESUMEN_IMPLEMENTACION.md   # Resumen técnico
```

## 🏗️ Arquitectura

Esta ubicación (`platform/swaggo/`) fue elegida para mantener la separación de responsabilidades:

- **`platform/`** - Infraestructura y herramientas de la aplicación
- **`docs/swaggo/`** - Documentación para desarrolladores

Los archivos en esta carpeta son parte de la capa de infraestructura/platform y son consumidos por el servidor en tiempo de ejecución.

---

**Generado automáticamente por swag CLI**
