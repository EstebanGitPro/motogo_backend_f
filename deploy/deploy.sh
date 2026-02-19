#!/bin/bash
# ============================================
# Script de Deploy - MotoGo Backend
# ============================================
# Uso: ./deploy.sh [primera-vez|actualizar]
#
# primera-vez: Setup completo (Docker, MySQL, Keycloak, Nginx, SSL)
# actualizar:  Solo actualiza el backend Go (binary + restart)
# ============================================

set -e

# --- Colores ---
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

DEPLOY_DIR="/home/motogo/motogo-backend"
BACKUP_DIR="/home/motogo/backups"

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  MotoGo Deploy Script${NC}"
echo -e "${GREEN}========================================${NC}"

# --- Verificar que se ejecuta como motogo ---
if [ "$(whoami)" != "motogo" ]; then
    echo -e "${RED}Error: Ejecutar como usuario 'motogo'${NC}"
    echo "Uso: su - motogo -c './deploy.sh $1'"
    exit 1
fi

case "$1" in
    primera-vez)
        echo -e "${YELLOW}▶ Modo: PRIMERA VEZ (setup completo)${NC}"
        echo ""

        # 1. Crear directorios
        echo -e "${GREEN}[1/8] Creando directorios...${NC}"
        mkdir -p $DEPLOY_DIR/config
        mkdir -p $DEPLOY_DIR/deploy
        mkdir -p $DEPLOY_DIR/docs/sql
        mkdir -p $DEPLOY_DIR/tools/keycloak-themes
        mkdir -p $BACKUP_DIR

        echo -e "${GREEN}[2/8] Verificando archivos necesarios...${NC}"
        REQUIRED_FILES=(
            "$DEPLOY_DIR/motogo-api"
            "$DEPLOY_DIR/config/prod-config.json"
            "$DEPLOY_DIR/config/serviceAccountKey.json"
            "$DEPLOY_DIR/deploy/docker-compose.production.yml"
            "$DEPLOY_DIR/deploy/.env.production"
        )
        for f in "${REQUIRED_FILES[@]}"; do
            if [ ! -f "$f" ]; then
                echo -e "${RED}  ✗ Falta: $f${NC}"
                echo "Sube los archivos primero con scp"
                exit 1
            fi
            echo -e "  ✓ $f"
        done

        # 3. Levantar Docker containers
        echo -e "${GREEN}[3/8] Levantando MySQL + Keycloak...${NC}"
        cd $DEPLOY_DIR/deploy
        docker compose --env-file .env.production -f docker-compose.production.yml up -d
        echo "Esperando que MySQL inicie (30 seg)..."
        sleep 30

        # 4. Crear usuario MySQL para la app
        echo -e "${GREEN}[4/8] Creando usuario MySQL para la app...${NC}"
        source .env.production
        docker exec motogo-mysql-app mysql -uroot -p"$MYSQL_ROOT_PASSWORD" -e "
            CREATE USER IF NOT EXISTS 'motogo_app'@'%' IDENTIFIED BY '$(grep -oP 'password.*\"(.*)\"' $DEPLOY_DIR/config/prod-config.json | head -1 | cut -d'"' -f2)';
            GRANT ALL PRIVILEGES ON motogoDb.* TO 'motogo_app'@'%';
            FLUSH PRIVILEGES;
        " 2>/dev/null || echo -e "${YELLOW}  (usuario ya existente, ok)${NC}"

        # 5. Ejecutar migraciones SQL
        echo -e "${GREEN}[5/8] Ejecutando migraciones SQL...${NC}"
        for sql_file in $(ls $DEPLOY_DIR/docs/sql/*.sql | sort); do
            echo -e "  Ejecutando: $(basename $sql_file)"
            docker exec -i motogo-mysql-app mysql -uroot -p"$MYSQL_ROOT_PASSWORD" motogoDb < "$sql_file" 2>/dev/null || true
        done

        # 6. Configurar permisos del binario
        echo -e "${GREEN}[6/8] Configurando binario...${NC}"
        chmod +x $DEPLOY_DIR/motogo-api

        # 7. Instalar servicio systemd
        echo -e "${GREEN}[7/8] Instalando servicio systemd...${NC}"
        sudo cp $DEPLOY_DIR/deploy/motogo-api.service /etc/systemd/system/
        sudo systemctl daemon-reload
        sudo systemctl enable motogo-api
        sudo systemctl start motogo-api

        # 8. Verificar
        echo -e "${GREEN}[8/8] Verificando servicios...${NC}"
        sleep 3
        echo ""
        echo "Docker containers:"
        docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" | grep motogo
        echo ""
        echo "Backend API:"
        sudo systemctl status motogo-api --no-pager | head -5
        echo ""
        echo -e "${GREEN}========================================${NC}"
        echo -e "${GREEN}  ✓ Deploy completado!${NC}"
        echo -e "${GREEN}========================================${NC}"
        echo ""
        echo "Siguiente paso: Configurar Nginx + SSL"
        echo "  sudo cp $DEPLOY_DIR/deploy/nginx-motogo.conf /etc/nginx/sites-available/motogo"
        echo "  sudo ln -s /etc/nginx/sites-available/motogo /etc/nginx/sites-enabled/"
        echo "  sudo nginx -t && sudo systemctl reload nginx"
        echo "  sudo certbot --nginx -d api.tudominio.com -d auth.tudominio.com"
        ;;

    actualizar)
        echo -e "${YELLOW}▶ Modo: ACTUALIZAR (solo backend)${NC}"
        echo ""

        # 1. Backup del binario actual
        echo -e "${GREEN}[1/3] Backup del binario actual...${NC}"
        if [ -f "$DEPLOY_DIR/motogo-api" ]; then
            cp $DEPLOY_DIR/motogo-api $BACKUP_DIR/motogo-api.$(date +%Y%m%d_%H%M%S)
        fi

        # 2. Copiar nuevo binario (ya debe estar subido)
        echo -e "${GREEN}[2/3] Verificando nuevo binario...${NC}"
        if [ ! -f "$DEPLOY_DIR/motogo-api-new" ]; then
            echo -e "${RED}Error: No se encuentra motogo-api-new${NC}"
            echo "Primero sube el nuevo binario:"
            echo "  scp motogo-api motogo@tuserver:/home/motogo/motogo-backend/motogo-api-new"
            exit 1
        fi
        mv $DEPLOY_DIR/motogo-api-new $DEPLOY_DIR/motogo-api
        chmod +x $DEPLOY_DIR/motogo-api

        # 3. Restart
        echo -e "${GREEN}[3/3] Reiniciando servicio...${NC}"
        sudo systemctl restart motogo-api
        sleep 2
        sudo systemctl status motogo-api --no-pager | head -5
        echo ""
        echo -e "${GREEN}✓ Backend actualizado exitosamente${NC}"
        ;;

    *)
        echo "Uso: $0 {primera-vez|actualizar}"
        echo ""
        echo "  primera-vez  Setup completo (Docker, MySQL, Keycloak, backend)"
        echo "  actualizar   Solo actualiza el binario del backend"
        exit 1
        ;;
esac
