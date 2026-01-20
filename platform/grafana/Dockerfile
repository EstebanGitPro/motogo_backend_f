FROM grafana/grafana:11.3.0

# Metadata
LABEL maintainer="Motogo Backend Team"
LABEL description="Custom Grafana image for Motogo Backend monitoring"
LABEL version="11.3.0"

# Configurar usuario root temporalmente para instalar plugins
USER root

# Instalar plugins útiles para monitoreo
RUN grafana-cli plugins install grafana-clock-panel && \
    grafana-cli plugins install grafana-simple-json-datasource && \
    grafana-cli plugins install grafana-piechart-panel

# Volver al usuario grafana por seguridad
USER grafana

# Exponer puerto por defecto de Grafana
EXPOSE 3000

# Variables de entorno por defecto
ENV GF_SECURITY_ADMIN_USER=admin
ENV GF_SECURITY_ADMIN_PASSWORD=0000
ENV GF_USERS_ALLOW_SIGN_UP=false
ENV GF_SERVER_ROOT_URL=http://localhost:3000

# El comando de inicio ya está definido en la imagen base
