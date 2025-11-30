#!/bin/bash

# Script para levantar Swagger UI en contenedor Docker
# Primero genera la documentación, luego levanta el contenedor con Swagger UI

set -e

echo "🔨 Generando documentación Swagger..."

# Verificar si existe la imagen, si no, construirla
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
  motogo-swag -c "/go/bin/swag init --generalInfo cmd/main.go --output platform/swaggo --parseInternal --parseDependency"

echo ""
echo "✅ Documentación generada"
echo ""
echo "🚀 Levantando Swagger UI..."

# Levantar contenedor de Swagger UI
docker-compose -f docker-compose.swagger-ui.yml up -d

echo ""
echo "✨ ¡Listo! Swagger UI disponible en:"
echo "   👉 http://localhost:3001"
echo ""
echo "Para detener Swagger UI:"
echo "   docker-compose -f docker-compose.swagger-ui.yml down"
echo ""
