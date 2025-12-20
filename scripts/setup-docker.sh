#!/bin/bash

# ============================================
# 🐳 MotoGo Backend - Docker Recovery Script
# ============================================
# Script para recrear todos los contenedores del backend MotoGo
# Útil cuando Docker Desktop se resetea o pierde los contenedores
#
# Autor: MotoGo Backend Team
# Fecha: 2025-12-17
# ============================================

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# ============================================
# Helper Functions
# ============================================

print_header() {
    echo -e "${BLUE}============================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}============================================${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_step() {
    echo -e "${BLUE}🔹 $1${NC}"
}

# ============================================
# Pre-flight Checks
# ============================================

print_header "VERIFICACIONES PREVIAS"

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    print_error "Docker no está corriendo. Por favor inicia Docker Desktop."
    exit 1
fi
print_success "Docker está corriendo"

# Check if we're in the correct directory
if [ ! -f "docker-compose.keycloak.yml" ]; then
    print_error "No se encontraron los archivos docker-compose"
    print_error "Por favor ejecuta este script desde la raíz del proyecto: /Users/estebandev/Documents/go/motogo_backend_f"
    exit 1
fi
print_success "Directorio correcto detectado"

# Check if .env file exists
if [ ! -f ".env" ]; then
    print_warning "Archivo .env no encontrado. Algunos servicios podrían no funcionar correctamente."
fi

echo ""

# ============================================
# Configuration
# ============================================

print_header "CONFIGURACIÓN"

# Ask user which services to start
echo "¿Qué servicios quieres iniciar?"
echo ""
echo "1) 💾 Solo MySQL (Base de datos)"
echo "2) 🔐 MySQL + Keycloak (Autenticación)"
echo "3) 📊 MySQL + Grafana Stack (Base de datos + Monitoreo)"
echo "4) 📖 Swagger UI (Documentación API)"
echo "5) 📈 InfluxDB (Opcional - solo si usas k6 con InfluxDB)"
echo "6) ✨ TODOS los servicios esenciales (MySQL + Keycloak + Grafana + Swagger)"
echo "7) 🚀 TODO (incluye InfluxDB)"
echo ""
read -p "Selecciona una opción [1-7] (default: 6): " choice
choice=${choice:-6}

START_MYSQL=false
START_KEYCLOAK=false
START_GRAFANA=false
START_SWAGGER=false
START_INFLUXDB=false

case $choice in
    1)
        START_MYSQL=true
        ;;
    2)
        START_MYSQL=true
        START_KEYCLOAK=true
        ;;
    3)
        START_MYSQL=true
        START_GRAFANA=true
        ;;
    4)
        START_SWAGGER=true
        ;;
    5)
        START_INFLUXDB=true
        ;;
    6)
        START_MYSQL=true
        START_KEYCLOAK=true
        START_GRAFANA=true
        START_SWAGGER=true
        ;;
    7)
        START_MYSQL=true
        START_KEYCLOAK=true
        START_GRAFANA=true
        START_SWAGGER=true
        START_INFLUXDB=true
        ;;
    *)
        print_error "Opción inválida. Saliendo..."
        exit 1
        ;;
esac

echo ""

# ============================================
# Stop and Remove Old Containers
# ============================================

print_header "LIMPIEZA DE CONTENEDORES ANTIGUOS"

print_step "Deteniendo contenedores existentes..."
docker-compose -f docker-compose.mysql.yml down 2>/dev/null || true
docker-compose -f docker-compose.keycloak.yml down 2>/dev/null || true
docker-compose -f docker-compose.grafana.yml down 2>/dev/null || true
docker-compose -f docker-compose.swagger-ui.yml down 2>/dev/null || true
docker-compose -f docker-compose.k6-influxdb.yml down 2>/dev/null || true

print_success "Contenedores antiguos detenidos"
echo ""

# ============================================
# Create Docker Network
# ============================================

print_header "CONFIGURANDO RED DOCKER"

# Create network if it doesn't exist
if ! docker network inspect motogo-network >/dev/null 2>&1; then
    print_step "Creando red motogo-network..."
    docker network create motogo-network
    print_success "Red motogo-network creada"
else
    print_success "Red motogo-network ya existe"
fi

echo ""

# ============================================
# Build and Start Services
# ============================================

print_header "INICIANDO SERVICIOS"

# Start MySQL FIRST (Keycloak depends on it)
if [ "$START_MYSQL" = true ]; then
    print_step "💾 Iniciando MySQL (App + Keycloak)..."
    docker-compose -f docker-compose.mysql.yml up -d
    print_success "MySQL App iniciado en puerto 3306"
    print_success "MySQL Keycloak iniciado en puerto 3309"
    echo "   App DB:"
    echo "     Host: localhost:3306"
    echo "     Database: motogoDb"
    echo "   Keycloak DB:"
    echo "     Host: localhost:3309"
    echo "     Database: keydb"
    echo "   User: root | Pass: 0000"
    echo ""
    print_warning "Esperando 10 segundos a que MySQL esté listo..."
    sleep 10
fi

# Start Keycloak (depends on MySQL)
if [ "$START_KEYCLOAK" = true ]; then
    if [ "$START_MYSQL" = false ]; then
        print_warning "⚠️  Keycloak necesita MySQL. Iniciando MySQL primero..."
        docker-compose -f docker-compose.mysql.yml up -d
        print_success "MySQL iniciado"
        sleep 10
    fi
    print_step "🔐 Iniciando Keycloak..."
    docker-compose -f docker-compose.keycloak.yml up -d --build
    print_success "Keycloak iniciado en puerto 8080"
    echo "   Admin Console: http://localhost:8080"
    echo ""
    
    # Configurar Keycloak (deshabilitar SSL para desarrollo)
    print_step "⏳ Esperando a que Keycloak esté listo (60 segundos)..."
    sleep 60
    
    print_step "🔧 Configurando Keycloak (deshabilitando SSL para desarrollo)..."
    
    # Autenticar con Keycloak
    if docker exec keycloak-prod /opt/keycloak/bin/kcadm.sh config credentials \
        --server http://localhost:8080 \
        --realm master \
        --user motogo-admin \
        --password 12345678 2>/dev/null; then
        
        # Deshabilitar SSL en realm master
        docker exec keycloak-prod /opt/keycloak/bin/kcadm.sh update realms/master -s sslRequired=NONE 2>/dev/null
        print_success "SSL deshabilitado en realm master"
        
        # Deshabilitar SSL en realm motogo
        docker exec keycloak-prod /opt/keycloak/bin/kcadm.sh update realms/motogo -s sslRequired=NONE 2>/dev/null
        print_success "SSL deshabilitado en realm motogo"
    else
        print_warning "No se pudo configurar Keycloak automáticamente"
        print_warning "Ejecuta manualmente después de que Keycloak esté listo:"
        echo "   docker exec -it keycloak-prod /opt/keycloak/bin/kcadm.sh config credentials --server http://localhost:8080 --realm master --user motogo-admin --password 12345678"
        echo "   docker exec -it keycloak-prod /opt/keycloak/bin/kcadm.sh update realms/master -s sslRequired=NONE"
        echo "   docker exec -it keycloak-prod /opt/keycloak/bin/kcadm.sh update realms/motogo -s sslRequired=NONE"
    fi
    echo ""
fi

# Start Grafana Stack
if [ "$START_GRAFANA" = true ]; then
    print_step "📊 Iniciando Grafana + Prometheus + Loki + Promtail + Log-Rotator..."
    docker-compose -f docker-compose.grafana.yml up -d --build
    print_success "Stack de Monitoreo iniciado"
    echo "   Grafana:    http://localhost:3000 (admin/admin)"
    echo "   Prometheus: http://localhost:9091"
    echo "   Loki:       http://localhost:3100"
    echo ""
fi

# Start Swagger UI
if [ "$START_SWAGGER" = true ]; then
    print_step "📖 Iniciando Swagger UI..."
    docker-compose -f docker-compose.swagger-ui.yml up -d
    print_success "Swagger UI iniciado en puerto 3001"
    echo "   Swagger UI: http://localhost:3001"
    echo ""
fi

# Start InfluxDB (optional)
if [ "$START_INFLUXDB" = true ]; then
    print_step "📈 Iniciando InfluxDB..."
    docker-compose -f docker-compose.k6-influxdb.yml up -d
    print_success "InfluxDB iniciado en puerto 8086"
    echo "   InfluxDB: http://localhost:8086"
    echo "   Username: admin"
    echo "   Password: adminpassword"
    echo ""
fi

# ============================================
# Health Checks
# ============================================

print_header "VERIFICANDO ESTADO DE SERVICIOS"

sleep 5  # Wait a bit for containers to start

print_step "Verificando contenedores..."
echo ""

# List running containers
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" | grep "motogo\|keycloak"

echo ""

# ============================================
# Summary
# ============================================

print_header "RESUMEN"

echo "Servicios iniciados correctamente:"
echo ""

if [ "$START_MYSQL" = true ]; then
    echo "💾 MySQL App:"
    echo "   URL: localhost:3306"
    echo "   DB: motogoDb"
    echo ""
    echo "💾 MySQL Keycloak:"
    echo "   URL: localhost:3309"
    echo "   DB: keydb"
    echo "   (root/0000 para ambos)"
    echo ""
fi

if [ "$START_KEYCLOAK" = true ]; then
    echo "🔐 Keycloak:"
    echo "   URL: http://localhost:8080"
    echo "   Realm: motogo"
    echo ""
fi

if [ "$START_GRAFANA" = true ]; then
    echo "📊 Grafana:"
    echo "   URL: http://localhost:3000"
    echo "   User: admin"
    echo "   Pass: admin"
    echo ""
    echo "📈 Prometheus:"
    echo "   URL: http://localhost:9091"
    echo ""
    echo "📝 Loki:"
    echo "   URL: http://localhost:3100"
    echo ""
fi

if [ "$START_SWAGGER" = true ]; then
    echo "📖 Swagger UI:"
    echo "   URL: http://localhost:3001"
    echo ""
fi

if [ "$START_INFLUXDB" = true ]; then
    echo "📈 InfluxDB:"
    echo "   URL: http://localhost:8086"
    echo "   User: admin"
    echo "   Pass: adminpassword"
    echo ""
fi

print_success "¡Todos los servicios están corriendo!"
echo ""
print_warning "Nota: Keycloak puede tardar 30-60 segundos en estar completamente listo"
print_warning "      Espera a que el healthcheck esté OK antes de usarlo"
echo ""

# ============================================
# Useful Commands
# ============================================

print_header "COMANDOS ÚTILES"

echo "Ver logs de un servicio:"
echo "  docker logs -f [nombre-contenedor]"
echo ""
echo "Reiniciar un servicio:"
echo "  docker restart [nombre-contenedor]"
echo ""
echo "Detener todos los servicios:"
echo "  docker-compose -f docker-compose.keycloak.yml down"
echo "  docker-compose -f docker-compose.grafana.yml down"
echo "  docker-compose -f docker-compose.swagger-ui.yml down"
echo ""
echo "Ver estado de todos los contenedores:"
echo "  docker ps -a"
echo ""

print_success "Script completado exitosamente 🎉"
