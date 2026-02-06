#!/bin/bash

# =============================================================================
# Script: start-swagger-ui.sh
# Descripción: Regenera la documentación Swagger y reinicia el contenedor
# Uso: ./scripts/start-swagger-ui.sh
# =============================================================================

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_ROOT"

echo "🔨 Regenerando documentación Swagger..."
echo ""

# Verificar si existe la imagen de swag, si no, construirla
if ! docker images | grep -q "motogo-swag"; then
    echo "📦 Construyendo imagen motogo-swag..."
    docker build -t motogo-swag -f platform/swaggo/Dockerfile.swag .
    echo ""
fi

# Generar documentación usando el contenedor de swag
docker run --rm \
  --entrypoint sh \
  -v "$(pwd):/app" \
  -w /app \
  motogo-swag -c "/go/bin/swag init -d ./ --generalInfo cmd/main.go --output platform/swaggo --parseInternal --parseDependency"

echo ""
echo "✅ Documentación regenerada"
echo ""

# Detener contenedor existente si está corriendo
echo "🔄 Reiniciando contenedor Swagger UI..."
docker compose -f docker-compose.swagger-ui.yml down 2>/dev/null || true

# Levantar contenedor de Swagger UI
docker compose -f docker-compose.swagger-ui.yml up -d

echo ""
echo "✨ ¡Listo! Swagger UI disponible en:"
echo "   👉 http://localhost:3001"
echo ""
echo "Comandos útiles:"
echo "   Detener:    docker compose -f docker-compose.swagger-ui.yml down"
echo "   Ver logs:   docker logs -f motogo-swagger-ui"
echo ""
