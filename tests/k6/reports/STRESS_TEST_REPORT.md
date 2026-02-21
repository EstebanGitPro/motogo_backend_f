# Informe de Pruebas de Estrés — MotoGo Backend

**Fecha:** 19 de Febrero de 2026  
**Versión del backend:** v0.13.0  
**Entorno:** MacBook Pro M1 Pro — 16 GB RAM  
**Herramienta:** [K6 v0.55+](https://k6.io/) con Grafana + Prometheus  

---

## 1. Estrategia de Pruebas de Carga

### 1.1 Objetivo

Evaluar la capacidad del backend de MotoGo para manejar carga concurrente creciente, identificando el punto de saturación y validando la estabilidad del sistema bajo estrés.

### 1.2 Escenarios Definidos

Se diseñaron **6 escenarios graduales** para evaluar diferentes aspectos del rendimiento:

| Escenario | VUs Máximos | Duración | Propósito |
|-----------|-------------|----------|-----------|
| **Smoke** | 1 | 30s | Sanity check post-deploy |
| **Load** | 100 | 40 min | Carga normal de producción (3 fases: mañana/tarde/noche) |
| **Stress** | 800 | 43 min | Estrés controlado, buscar límite de capacidad |
| **Spike** | 150 | 16 min | Pico viral moderado (3x tráfico en 30s) |
| **Endurance** | 75 | 2 horas | Detección de memory leaks |
| **Breakpoint** | 400 | 40 min | Encontrar el límite exacto con incrementos de 50 VUs |

### 1.3 Flujo de Usuario Simulado

Cada usuario virtual (VU) ejecuta un flujo realista de motociclista colombiano:

```
1. Health Check (heartbeat)           → 100% de VUs
2. Consultar catálogos                → 70% de VUs
   ├─ GET /brands
   ├─ GET /branch-types
   └─ GET /services
3. Consultar perfil                   → 30% de VUs (autenticado)
4. Consultar motocicletas             → 40% de VUs (autenticado)
5. Buscar talleres cercanos (geo)     → 80% de VUs (autenticado)
   └─ Coordenadas aleatorias en Bogotá (lat 4.6-4.83, lng -74.2 a -74.0)
```

### 1.4 Think-Time Realista

Se implementó distribución de tiempos de espera entre acciones:

| % de usuarios | Tiempo entre acciones | Simula |
|--------------|----------------------|--------|
| 60% | 1–3 segundos | Navegación rápida |
| 30% | 3–5 segundos | Lectura normal |
| 10% | 5–10 segundos | Usuario distraído |

### 1.5 Mecanismos de Protección

| Mecanismo | Descripción |
|-----------|-------------|
| **Health Check Abort** | Si `/health` falla 3 veces consecutivas, el test se aborta automáticamente |
| **Token Refresh Proactivo** | Cada VU refresca su JWT 2 minutos antes de expirar (evita avalancha de 401) |
| **URL Tag Grouping** | Coordenadas GPS agrupadas bajo tag `nearby_search` (evita explosión de métricas) |
| **Login con reintentos** | 3 intentos de autenticación con backoff de 2 segundos |

### 1.6 Optimización del Pool de Conexiones

Antes de las pruebas se optimizó la configuración de conexiones MySQL:

| Parámetro | Valor anterior | Valor optimizado | Razón |
|-----------|---------------|-----------------|-------|
| `max_open_conns` | 50 | 25 | Evitar saturación de MySQL |
| `max_idle_conns` | 20 | 10 | 40% del máximo, balance eficiente |
| `conn_max_lifetime` | 27,000s (7.5h) | 300s (5 min) | Reciclar antes del `wait_timeout` de MySQL |
| `conn_max_idle_time` | 26,000s (7.2h) | 60s (1 min) | Liberar conexiones inactivas rápidamente |

### 1.7 Thresholds (Umbrales de Aceptación)

| Métrica | Umbral |
|---------|--------|
| Latencia p95 | < 1,000 ms |
| Latencia p99 | < 2,000 ms |
| Error rate global | < 10% |
| Errores de login | < 5% |
| Errores de negocio (5xx) | < 5% |

---

## 2. Resultados del Stress Test (800 VUs)

### 2.1 Resumen Ejecutivo

| Métrica | Resultado | Estado |
|---------|-----------|--------|
| **Duración total** | 43 min 07s / 43 min planeados | ✅ Completó al 100% |
| **VUs máximos** | 800 | ✅ Alcanzó el objetivo |
| **Iteraciones completadas** | 214,524 | ✅ Sin interrupciones |
| **Error rate global** | **0.00%** | ✅ Perfecto |
| **Checks fallidos** | 61 de 1,974,732 (0.003%) | ✅ Insignificante |
| **Throughput** | 381.65 req/s | ✅ Excelente |

### 2.2 Rendimiento Global

| Métrica | Valor |
|---------|-------|
| **Throughput** | 381.65 req/s |
| **Latencia p50** | ~3 ms |
| **Latencia p95** | 6.06 ms |
| **Latencia p99** | 24.4 ms |
| **Latencia máxima** | 1,562 ms |
| **Data recibida** | 3.3 GB (1.3 MB/s) |
| **Duración promedio de iteración** | 5.04 s |

### 2.3 Autenticación

| Métrica | Resultado |
|---------|-----------|
| Login status 200 | 100% |
| Login response time < 2s | 100% |
| Login errors | 0.00% |

El sistema de autenticación con Keycloak resultó sólido y **no fue cuello de botella**. El token duró 54,000s (15 horas), eliminando cualquier problema de expiración durante la prueba.

### 2.4 Health Checks

| Métrica | Resultado |
|---------|-----------|
| Health 200 | 100% |
| Health fast < 500ms | 99.998% (3 excedencias de 214,524) |

Las 3 excedencias de 500ms coincidieron con picos de carga donde el servidor alcanzó su capacidad máxima. El sistema **nunca dejó de responder**.

### 2.5 Catálogos (Endpoints Públicos)

| Endpoint | Éxito HTTP | Éxito Latencia (< 1s) | Observación |
|----------|-----------|----------------------|-------------|
| `/brands` | 100% | 99.99% (16 lentos) | Queries ligeras, bien cacheadas |
| `/services` | 100% | 99.99% (13 lentos) | Catálogo de servicios estable |
| `/branch-types` | 100% | 100% | El endpoint más rápido |

### 2.6 Endpoints Autenticados

| Endpoint | Éxito HTTP | Éxito Latencia | Observación |
|----------|-----------|----------------|-------------|
| `persons/me` | 100% | 99.98% (< 800ms) | 14 requests lentos de 64K+ |
| `motorcycles` | 100% | 99.98% (< 800ms) | 15 requests lentos de 85K+ |

### 2.7 Queries Geoespaciales (La métrica más crítica)

| Métrica | Valor |
|---------|-------|
| Éxito HTTP (200/204) | **100%** |
| Éxito Latencia (< 2s) | **100%** |
| Latencia promedio | 4.89 ms |
| Latencia p95 | 10.47 ms |
| Latencia máxima | 1,562 ms (1.5 s) |

- El **promedio de 4.89ms** para queries geoespaciales es impresionante.
- Solo **~20 queries** de más de 70,000 tardaron entre 1.1s y 1.5s, probablemente por concentración de talleres en ciertas zonas de Bogotá.

---

## 3. Progresión por Fase de Carga

| Fase | VUs | Duración | Throughput | p95 Latencia | Error Rate | Estado |
|------|-----|----------|------------|-------------|------------|--------|
| Ramp-up 1 | 0 → 100 | 3 min | ~50 req/s | < 5 ms | 0% | ✅ |
| Estabilización 1 | 100 | 2 min | ~50 req/s | < 5 ms | 0% | ✅ |
| Ramp-up 2 | 100 → 200 | 3 min | ~100 req/s | < 5 ms | 0% | ✅ |
| Estabilización 2 | 200 | 2 min | ~100 req/s | < 5 ms | 0% | ✅ |
| Ramp-up 3 | 200 → 300 | 3 min | ~150 req/s | < 6 ms | 0% | ✅ |
| Estabilización 3 | 300 | 2 min | ~150 req/s | < 6 ms | 0% | ✅ |
| Ramp-up 4 | 300 → 400 | 3 min | ~200 req/s | < 6 ms | 0% | ✅ |
| Estabilización 4 | 400 | 2 min | ~200 req/s | < 6 ms | 0% | ✅ |
| Ramp-up 5 | 400 → 500 | 3 min | ~250 req/s | ~6 ms | 0% | ✅ |
| Estabilización 5 | 500 | 2 min | ~250 req/s | ~6 ms | 0% | ✅ |
| **Zona de riesgo** | 500 → 800 | 6 min | ~380 req/s | 6.06 ms | 0% | ✅ |
| **Pico máximo** | **800** | **2 min** | **381 req/s** | **6.06 ms** | **0%** | **✅** |
| Ramp-down | 800 → 0 | 3 min | Decreciente | < 10 ms | 0% | ✅ |

---

## 4. Comparativa: Expectativa vs Realidad

| Escenario | Latencia p95 Esperada | Latencia p95 Real | Factor de mejora |
|-----------|----------------------|-------------------|-----------------|
| 100 VUs | < 500 ms | ~2 ms | **250x mejor** |
| 400 VUs | < 1,000 ms | ~4 ms | **250x mejor** |
| 800 VUs | < 2,000 ms | 6.06 ms | **330x mejor** |

El sistema resultó **~100-330x más rápido** de lo requerido en todos los escenarios.

---

## 5. Hallazgos y Recomendaciones

### 5.1 Lo que funcionó bien

1. **Pool de conexiones optimizado** — El cambio de lifetime de 7.5h a 5min eliminó los problemas de conexiones stale.
2. **Go + Gin** — La combinación demostró excelente rendimiento bajo carga con goroutines.
3. **Queries geoespaciales** — Los índices espaciales de MySQL respondieron bien incluso a 800 VUs.
4. **Token refresh proactivo** — Eliminó la avalancha de 401s vista en iteraciones anteriores.

### 5.2 Mejoras sugeridas para producción

| Mejora | Prioridad | Impacto |
|--------|-----------|---------|
| Agregar paginación a `/branches/nearby` | Media | Reducir payload en zonas con muchos talleres |
| Optimizar índice geoespacial (SPATIAL INDEX) | Baja | Reducir las ~20 queries que tardaron > 1s |
| Implementar caché Redis para catálogos | Baja | Reducir carga a MySQL en endpoints públicos |

---

## 6. Conclusión

La aplicación MotoGo soportó satisfactoriamente **800 usuarios concurrentes** durante **43 minutos** de prueba de estrés gradual, con un **error rate del 0.00%** y latencias p95 de **6.06ms**, muy por debajo de los umbrales aceptables (1–2s).

El punto de saturación **no fue alcanzado**, indicando que el límite real está por encima de 800 usuarios concurrentes en el hardware de pruebas.

### Límite detectado

| Recurso | Uso observado | Límite estimado |
|---------|--------------|-----------------|
| Usuarios concurrentes | 800 | > 800 (no alcanzado) |
| Requests/segundo | 381 req/s | > 400 req/s |
| Latencia p95 | 6.06 ms | < 10 ms |
| RAM del sistema | ~12–14 GB | 16 GB (MacBook) |

### Comando de ejecución

```bash
# Stress test con métricas en Grafana
k6 run -o experimental-prometheus-rw \
  --env K6_PROMETHEUS_RW_SERVER_URL=http://localhost:9091/api/v1/write \
  --env SCENARIO=stress \
  tests/k6/load_test.js
```

---

## 7. Glosario

| Término | Significado |
|---------|-------------|
| **VU** | Virtual User — un usuario simulado ejecutando el script |
| **p50/p95/p99** | Percentil 50/95/99 de la latencia |
| **Throughput** | Cantidad de requests procesados por segundo (req/s) |
| **Error rate** | Porcentaje de requests que devolvieron un HTTP status de error |
| **Ramp up** | Fase donde se incrementan gradualmente los VUs |
| **Think-time** | Tiempo de espera entre acciones, simula lectura del usuario |
| **Breakpoint** | Punto donde el sistema empieza a degradarse |
| **OOM** | Out Of Memory — cuando un proceso agota la RAM disponible |
