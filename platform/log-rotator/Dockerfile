FROM alpine:3.19

# Instalar logrotate y herramientas necesarias
RUN apk add --no-cache \
    bash \
    dcron \
    docker-cli \
    findutils \
    logrotate \
    tzdata

# Configurar timezone
ENV TZ=America/Bogota
# Configurar timezone, crear directorios, hacer ejecutables los scripts
RUN ln -snf /usr/share/zoneinfo/"$TZ" /etc/localtime && echo "$TZ" > /etc/timezone \
    && mkdir -p /var/log /var/lib/logrotate

# Copiar configuraciones
COPY logrotate.conf /etc/logrotate.conf
COPY logrotate.d/ /etc/logrotate.d/
COPY cleanup-volumes.sh /usr/local/bin/cleanup-volumes.sh
COPY entrypoint.sh /entrypoint.sh

# Hacer ejecutables y crear estado
RUN chmod +x /usr/local/bin/cleanup-volumes.sh /entrypoint.sh \
    && touch /var/lib/logrotate/status

ENTRYPOINT ["/entrypoint.sh"]
