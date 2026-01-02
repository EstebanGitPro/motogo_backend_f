# 🔧 Guía de Recuperación del Stack de Motogo

## Problema Común: Redes Docker Desconectadas

Cuando ejecutas otros proyectos con Docker (especialmente con `docker-compose`), pueden ocurrir **conflictos de redes** que desconectan tus contenedores de Motogo de sus redes internas.

### 🔍 Síntomas

1. **Métricas no aparecen en Grafana**
2. **Logs no se muestran en los dashboards**
3. **Swagger UI no carga** (Connection refused en localhost:3001)
4. **Errores en logs de Promtail**: `lookup loki on 127.0.0.11:53: no such host`

### 🩺 Diagnóstico Rápido

```bash
# Ver estado de contenedores
docker ps | grep motogo

# Verificar redes
docker inspect motogo-prometheus --format='{{json .NetworkSettings.Networks}}'
docker inspect motogo-loki --format='{{json .NetworkSettings.Networks}}'
docker inspect motogo-swagger-ui --format='{{json .NetworkSettings.Networks}}'

# Si ves {} en lugar de detalles de red, el contenedor está desconectado
```

---

## 🚀 Scripts de Recuperación

### 1. Reparar Stack de Monitoreo Completo

**Cuándo usarlo:**
- No ves métricas en Prometheus/Grafana
- Los logs no aparecen
- Grafana muestra "No data"

**Comando:**
```bash
./docs/scripts/fix-monitoring-stack.sh
```

**Qué hace:**
1. ✅ Detiene todos los servicios de monitoreo (Prometheus, Grafana, Loki, Promtail)
2. ✅ Los reinicia en el orden correcto
3. ✅ Verifica que todos estén en la red `motogo-network`
4. ✅ Valida conectividad entre servicios
5. ✅ Muestra estado final con URLs de acceso

**Tiempo estimado:** 10-15 segundos

---

### 2. Reparar Solo Swagger UI

**Cuándo usarlo:**
- Swagger UI no carga en http://localhost:3001
- Recibes "Connection refused"

**Comando:**
```bash
./docs/scripts/fix-swagger-ui.sh
```

**Qué hace:**
1. ✅ Detiene Swagger UI
2. ✅ Lo reinicia conectado a `motogo-network`
3. ✅ Verifica conectividad

**Tiempo estimado:** 3-5 segundos

---

## 🛠️ Hacer Scripts Ejecutables (Solo la Primera Vez)

```bash
chmod +x docs/scripts/fix-monitoring-stack.sh
chmod +x docs/scripts/fix-swagger-ui.sh
```

---

## 📋 Checklist de Verificación Post-Recuperación

Después de ejecutar los scripts, verifica:

### ✅ Servicios en Estado UP

```bash
docker ps --filter "name=motogo-" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
```

Deberías ver:
- ✅ `motogo-prometheus` → Up → `0.0.0.0:9091->9090/tcp`
- ✅ `motogo-grafana` → Up → `0.0.0.0:3000->3000/tcp`
- ✅ `motogo-loki` → Up → `3100/tcp`
- ✅ `motogo-promtail` → Up
- ✅ `motogo-swagger-ui` → Up → `0.0.0.0:3001->8080/tcp`

### ✅ Accesos Funcionando

| Servicio       | URL                   | Usuario/Password |
| -------------- | --------------------- | ---------------- |
| **Grafana**    | http://localhost:3000 | admin/admin      |
| **Prometheus** | http://localhost:9091 | -                |
| **Swagger UI** | http://localhost:3001 | -                |

### ✅ Sin Errores en Logs

```bash
# Verificar Promtail (no debe tener errores de "no such host")
docker logs motogo-promtail --tail 20

# Verificar Loki
docker logs motogo-loki --tail 20
```

---

## 🔄 Prevención

### Buenas Prácticas

1. **Usa nombres de red únicos** en tus otros proyectos
2. **Detén servicios de otros proyectos** cuando no los uses:

   ```bash
   docker-compose -f otro-proyecto/docker-compose.yml down
   ```
3. **Limpia contenedores huérfanos** periódicamente:

   ```bash
   docker system prune -a --volumes
   ```

### Si Trabajas con Múltiples Proyectos Docker

Considera usar **perfiles de Docker Compose** o **diferentes contextos** para aislar proyectos.

---

## 🆘 Solución de Problemas Avanzados

### Backend No Aparece en Prometheus

**Problema:** Target `motogo-backend` en Prometheus muestra "DOWN"

**Causa:** El backend no está corriendo

**Solución:**
```bash
go run cmd/main.go
# O con logging:
./docs/scripts/run-backend-with-logging.sh
```

### Puerto 9090 Ya en Uso

**Problema:** Otro proyecto usa el puerto 9090

**Solución:** Ya está configurado en puerto 9091. Si necesitas cambiarlo:
1. Edita `docker-compose.grafana.yml` → sección `prometheus.ports`
2. Reinicia: `./docs/scripts/fix-monitoring-stack.sh`

### Logs Antiguos No Aparecen

**Problema:** Solo ves logs nuevos en Grafana

**Causa:** Promtail solo lee desde la posición actual del archivo

**Solución:**
1. Elimina el archivo de posiciones:
   ```bash
   docker exec motogo-promtail rm /tmp/positions.yaml
   ```
2. Reinicia Promtail:
   ```bash
   docker-compose -f docker-compose.grafana.yml restart promtail
   ```

---

## 📞 Contacto

Si los scripts no resuelven tu problema, revisa:
1. Logs de Docker: `docker logs <nombre-contenedor>`
2. Estado de redes: `docker network ls`
3. Inspección de red: `docker network inspect motogo_backend_f_motogo-network`

---

**Última actualización:** 2025-12-13  
**Versión:** 1.0
