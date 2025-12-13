# Docker Setup para Swagger

Configuración de Docker para generar documentación Swagger sin necesidad de tener `swag` instalado localmente.

## 🐳 Problema Original

La imagen oficial `ghcr.io/swaggo/swag:latest` no funciona porque:
- No incluye Go en el contenedor
- `swag` necesita Go para analizar el código
- El contenedor se apaga con error: `"go": executable file not found in $PATH`

## ✅ Solución Implementada

Creamos un **Dockerfile personalizado** (`Dockerfile.swag`) que incluye:
- ✅ Go 1.23
- ✅ swag CLI instalado
- ✅ Todas las dependencias necesarias

## 🚀 Uso

### Opción 1: Script Solo Docker (Más Simple)

Si solo quieres usar Docker sin verificaciones:

```bash
./generate-swagger-docker.sh
```

Este script:
- ✅ Solo usa Docker (no verifica instalaciones locales)
- ✅ Construye la imagen si no existe
- ✅ Genera la documentación
- ✅ Más simple y directo

### Opción 2: Script Automático

El script original intenta múltiples métodos:

```bash
./generate-swagger.sh
```

Este script intenta:
1. swag en PATH
2. swag en GOPATH
3. Docker (como último recurso)

### Opción 3: Docker Manual

Si quieres más control:

```bash
# Construir la imagen
docker build -t motogo-swag -f platform/swaggo/Dockerfile.swag .

# Generar documentación
docker run --rm \
  -v "$(pwd):/app" \
  -w /app \
  motogo-swag init \
  --generalInfo cmd/main.go \
  --output platform/swaggo \
  --parseInternal \
  --parseDependency
```

## 📦 Imagen Docker

### Construcción

```bash
docker build -t motogo-swag -f platform/swaggo/Dockerfile.swag .
```

### Detalles de la Imagen

- **Imagen base:** `golang:1.23-alpine`
- **Tamaño:** ~500MB (Go + Alpine + swag)
- **Herramientas incluidas:**
  - Go 1.23
  - swag CLI (latest)
  - git (para dependencias)

### Verificar la Imagen

```bash
# Listar imágenes
docker images | grep motogo-swag

# Verificar versión de swag (nota: usa --version, no version)
docker run --rm motogo-swag --version
```

## 🔄 Workflow con Docker

### Primera Vez

```bash
./generate-swagger.sh
# → Detecta que no tienes swag local
# → Construye imagen Docker (toma ~2-3 min)
# → Genera documentación
```

### Siguientes Veces

```bash
./generate-swagger.sh
# → Usa imagen ya construida
# → Genera documentación (toma ~10 seg)
```

## 🆚 Comparación: Local vs Docker

| Aspecto                    | swag Local           | swag Docker                 |
| -------------------------- | -------------------- | --------------------------- |
| **Primera ejecución**      | Rápida (<5s)         | Lenta (~3min build)         |
| **Ejecuciones siguientes** | Muy rápida (<5s)     | Rápida (~10s)               |
| **Requisitos**             | Go + swag instalados | Solo Docker                 |
| **Portabilidad**           | Depende del sistema  | Funciona en cualquier lugar |
| **Tamaño**                 | ~20MB (binario swag) | ~500MB (imagen completa)    |

## 💡 Recomendaciones

### Para Desarrollo Local

✅ Instala `swag` localmente:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

**Ventajas:**
- Más rápido
- Menos espacio en disco
- Mejor para iteraciones frecuentes

### Para CI/CD o Entornos Sin Go

✅ Usa Docker:
```bash
docker build -t motogo-swag -f platform/swaggo/Dockerfile.swag .
```

**Ventajas:**
- No requiere Go instalado
- Entorno consistente
- Fácil de integrar en pipelines

## 🧹 Limpieza

### Eliminar Imagen Docker

```bash
docker rmi motogo-swag
```

### Eliminar Archivos Generados

```bash
rm -rf platform/swaggo/docs.go
rm -rf platform/swaggo/swagger.json
rm -rf platform/swaggo/swagger.yaml
```

## 🐛 Troubleshooting

### "Cannot connect to Docker daemon"

**Problema:** Docker no está corriendo

**Solución:**
```bash
# macOS
open -a Docker

# Linux

sudo systemctl start docker
```

### "Permission denied"

**Problema:** El usuario no tiene permisos para Docker

**Solución:**
```bash
# Agregar usuario al grupo docker (Linux)
sudo usermod -aG docker $USER
newgrp docker
```

### "Image build failed"

**Problema:** Error al construir la imagen

**Solución:**
```bash
# Limpiar caché de Docker
docker builder prune

# Reconstruir sin caché

docker build --no-cache -t motogo-swag -f platform/swaggo/Dockerfile.swag .
```

## 📚 Referencias

- [Documentación oficial de swaggo](https://github.com/swaggo/swag)
- [Dockerfile best practices](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/)
- [Go Docker images](https://hub.docker.com/_/golang)

---

**Nota:** Esta configuración es parte de la infraestructura de documentación de la API Motogo.
