# 🔄 Log Rotator Service

Servicio automatizado para **rotación de logs** y **limpieza de volúmenes** de Docker en el proyecto Motogo Backend.

## 📋 Características

### ✅ Rotación de Logs

- **Frecuencia:** Cada hora (verifica si hay logs para rotar)
- **Retención:** 7 días
- **Compresión:** Logs antiguos se comprimen automáticamente con gzip
- **Logs afectados:** `/tmp/motogo-logs/backend-*.log`

### ✅ Limpieza de Volúmenes

- **Frecuencia:** Diariamente a las 2:00 AM
- **Retención:** 7 días
- **Volúmenes limpiados:**
  - `prometheus-data` (bloques TSDB y WAL)
  - `loki-data` (chunks, índices, WAL)
  - `grafana-data` (sesiones y logs temporales)

---

## 🚀 Uso

### Iniciar el Servicio

El log-rotator se inicia automáticamente con el stack de monitoreo:

```bash
docker-compose -f docker-compose.grafana.yml up -d log-rotator
```

O iniciar todo el stack:

```bash
./docs/scripts/fix-monitoring-stack.sh
```

### Verificar Estado

```bash
# Ver logs del rotator
docker logs motogo-log-rotator

# Ver logs de rotación
docker exec motogo-log-rotator cat /var/log/logrotate.log

# Ver logs de limpieza de volúmenes
docker exec motogo-log-rotator cat /var/log/volume-cleanup.log
```

### Ejecutar Manualmente

```bash
# Rotar logs manualmente
docker exec motogo-log-rotator /usr/sbin/logrotate -f /etc/logrotate.conf

# Limpiar volúmenes manualmente
docker exec motogo-log-rotator /usr/local/bin/cleanup-volumes.sh
```

---

## 📁 Estructura de Archivos

```
platform/log-rotator/
├── Dockerfile                    # Imagen del contenedor
├── logrotate.conf               # Configuración principal de logrotate
├── logrotate.d/
│   └── motogo-backend          # Reglas específicas para logs del backend
├── cleanup-volumes.sh           # Script de limpieza de volúmenes Docker
├── entrypoint.sh               # Script de inicio del contenedor
└── README.md                   # Esta documentación
```

---

## ⚙️ Configuración

### Cambiar Retención (días)

Edita el archivo `docker-compose.grafana.yml`:

```yaml
log-rotator:
  environment:
    - RETENTION_DAYS=7  # Cambiar a 14, 30, etc.
```

Luego reinicia:

```bash
docker-compose -f docker-compose.grafana.yml restart log-rotator
```

### Cambiar Horario de Limpieza

Edita `platform/log-rotator/entrypoint.sh` y modifica el cron:

```bash
# Cambiar "0 2 * * *" a tu horario preferido (formato: min hora día mes día_semana)
0 2 * * * /usr/local/bin/cleanup-volumes.sh
```

Luego reconstruye la imagen:

```bash
docker-compose -f docker-compose.grafana.yml up -d --build log-rotator
```

---

## 📊 Monitoreo

### Logs de Actividad

El servicio genera logs en dos archivos dentro del contenedor:

1. **`/var/log/logrotate.log`** - Actividad de rotación de logs

   ```bash
   docker exec motogo-log-rotator tail -f /var/log/logrotate.log
   ```
2. **`/var/log/volume-cleanup.log`** - Actividad de limpieza de volúmenes

   ```bash
   docker exec motogo-log-rotator tail -f /var/log/volume-cleanup.log
   ```

### Heartbeat

El contenedor registra un "heartbeat" cada 6 horas para confirmar que está activo:

```bash
docker exec motogo-log-rotator grep "is alive" /var/log/logrotate.log
```

---

## 🔍 Troubleshooting

### El servicio no inicia

```bash
# Ver logs de error
docker logs motogo-log-rotator

# Verificar que los volúmenes estén montados
docker exec motogo-log-rotator ls -la /var/log/motogo /prometheus /loki /grafana
```

### Los logs no se rotan

```bash
# Verificar configuración
docker exec motogo-log-rotator cat /etc/logrotate.d/motogo-backend

# Ejecutar en modo debug
docker exec motogo-log-rotator /usr/sbin/logrotate -d /etc/logrotate.conf
```

### Los volúmenes no se limpian

```bash
# Ver último log de limpieza
docker exec motogo-log-rotator tail -50 /var/log/volume-cleanup.log

# Ejecutar limpieza manual con verbose
docker exec motogo-log-rotator bash -x /usr/local/bin/cleanup-volumes.sh
```

### Verificar espacio en disco

```bash
# Dentro del contenedor
docker exec motogo-log-rotator df -h

# Tamaño de volúmenes
docker exec motogo-log-rotator du -sh /prometheus /loki /grafana
```

---

## 🛡️ Seguridad

- El contenedor tiene acceso **read-only** al socket de Docker
- Tiene acceso **read-write** a volúmenes específicos para limpieza
- Corre con timezone configurado (America/Bogota)
- No expone puertos externos

---

## 📈 Impacto en Rendimiento

- **CPU:** Mínimo (solo durante ejecuciones programadas)
- **RAM:** ~20MB en idle
- **Disco:** Libera espacio eliminando datos antiguos
- **Red:** No usa red (solo operaciones locales)

---

## 🔄 Actualización

Para actualizar el servicio con nuevos scripts o configuraciones:

```bash
# Reconstruir imagen
docker-compose -f docker-compose.grafana.yml build log-rotator

# Reiniciar con nueva imagen
docker-compose -f docker-compose.grafana.yml up -d log-rotator
```

---

## 📝 Logs de Ejemplo

### Rotación Exitosa

```
[2025-12-13 02:00:01] Ejecutando logrotate...
rotating pattern: /var/log/motogo/backend-*.log  after 1 days (7 rotations)
rotating file /var/log/motogo/backend-20251213.log
compressing log with: /bin/gzip
2025-12-13 02:00:01 - Logs rotados exitosamente
```

### Limpieza de Volúmenes

```
[2025-12-13 02:00:05] === Iniciando limpieza de volúmenes Docker ===
[2025-12-13 02:00:05] Limpiando datos de Prometheus (>7 días)...
[2025-12-13 02:00:05]   Eliminando bloque Prometheus: 01KBAVQVYN2VEQ2QHKRKVPZQV1
[2025-12-13 02:00:06] Prometheus: 3 bloques eliminados
[2025-12-13 02:00:06] Loki: 156 chunks eliminados
[2025-12-13 02:00:06] === Limpieza de volúmenes completada ===
```

---

## 📞 Contacto

Si tienes problemas con el log-rotator, revisa:

1. Logs del contenedor: `docker logs motogo-log-rotator`
2. Estado del cron: `docker exec motogo-log-rotator crontab -l`
3. Espacio en disco: `df -h`

---

**Última actualización:** 2025-12-13  
**Versión:** 1.0
