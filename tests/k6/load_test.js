import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { Rate, Trend, Counter } from 'k6/metrics';
import { randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';
import exec from 'k6/execution';

// ============================================
// MotoGo API — Pruebas de Carga Realistas v3.0
// ============================================
// ESCENARIOS:
//   smoke:      Sanity check (1 usuario, 30s)
//   load:       Carga normal progresiva (hasta 100 VUs)
//   stress:     Estrés controlado (hasta 800 VUs)
//   spike:      Pico repentino moderado (hasta 150 VUs)
//   endurance:  Prueba de resistencia (2 horas, carga media)
//   breakpoint: Buscar límite real (incrementos de 50 VUs)
//
// Uso:
//   k6 run tests/k6/load_test.js
//   k6 run --env SCENARIO=smoke tests/k6/load_test.js
//   k6 run --env SCENARIO=stress tests/k6/load_test.js
//
// Con Grafana (Prometheus remote write):
//   k6 run -o experimental-prometheus-rw \
//     --env K6_PROMETHEUS_RW_SERVER_URL=http://localhost:9091/api/v1/write \
//     --env SCENARIO=load tests/k6/load_test.js

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8085';
const API = `${BASE_URL}/motogo/api/v1`;
const USERNAME = __ENV.K6_USERNAME || 'andresortiz@yopmail.com';
const PASSWORD = __ENV.K6_PASSWORD || 'Secret123*';
const SCENARIO = __ENV.SCENARIO || 'load';

// Health check abort: máximo de fallos consecutivos antes de abortar
const MAX_HEALTH_FAILURES = Number.parseInt(__ENV.MAX_HEALTH_FAILURES || '3', 10);

// Token refresh: refrescar 2 minutos antes de expirar (margen de seguridad)
const TOKEN_REFRESH_MARGIN_S = 120;

// ── Métricas de negocio personalizadas ─────────────────
const loginErrors = new Rate('login_errors');
const catalogTime = new Trend('catalog_response_time');
const authTime = new Trend('auth_response_time');
const geoQueryTime = new Trend('geo_query_response_time');
const businessErrors = new Rate('business_errors');
const throughput = new Counter('requests_total');
const tokenRefreshes = new Counter('token_refreshes');

// ── Estado per-VU ──────────────────────────────────────
let consecutiveHealthFailures = 0;
let vuToken = null;
let vuTokenExpiresAt = 0; // timestamp en segundos

// ── ESCENARIOS REALISTAS ───────────────────────────────

const scenarios = {
  smoke: {
    smoke: {
      executor: 'constant-vus',
      vus: 1,
      duration: '30s',
      tags: { test_type: 'smoke' },
    },
  },

  load: {
    morning_rush: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '2m', target: 50 },
        { duration: '5m', target: 50 },
        { duration: '3m', target: 100 },
        { duration: '5m', target: 100 },
        { duration: '2m', target: 80 },
      ],
      tags: { test_type: 'load', period: 'morning' },
    },
    afternoon_normal: {
      executor: 'ramping-vus',
      startVUs: 0,
      startTime: '17m',
      stages: [
        { duration: '2m', target: 60 },
        { duration: '10m', target: 60 },
        { duration: '2m', target: 40 },
      ],
      tags: { test_type: 'load', period: 'afternoon' },
    },
    evening_low: {
      executor: 'ramping-vus',
      startVUs: 0,
      startTime: '31m',
      stages: [
        { duration: '2m', target: 20 },
        { duration: '5m', target: 20 },
        { duration: '2m', target: 0 },
      ],
      tags: { test_type: 'load', period: 'evening' },
    },
  },

  stress: {
    stress_gradual: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '3m', target: 100 },
        { duration: '2m', target: 100 },
        { duration: '3m', target: 200 },
        { duration: '2m', target: 200 },
        { duration: '3m', target: 300 },
        { duration: '2m', target: 300 },
        { duration: '3m', target: 400 },
        { duration: '2m', target: 400 },
        { duration: '3m', target: 500 },
        { duration: '2m', target: 500 },
        { duration: '3m', target: 600 },
        { duration: '2m', target: 600 },
        { duration: '3m', target: 700 },
        { duration: '2m', target: 700 },
        { duration: '3m', target: 800 },
        { duration: '2m', target: 800 },
        { duration: '3m', target: 0 },
      ],
      tags: { test_type: 'stress' },
    },
  },

  spike: {
    spike_moderate: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '5m', target: 50 },
        { duration: '30s', target: 150 },
        { duration: '3m', target: 150 },
        { duration: '1m', target: 80 },
        { duration: '5m', target: 80 },
        { duration: '2m', target: 0 },
      ],
      tags: { test_type: 'spike' },
    },
  },

  endurance: {
    endurance_test: {
      executor: 'constant-vus',
      vus: 75,
      duration: '2h',
      tags: { test_type: 'endurance' },
    },
  },

  breakpoint: {
    breakpoint_search: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '2m', target: 50 },
        { duration: '3m', target: 50 },
        { duration: '2m', target: 100 },
        { duration: '3m', target: 100 },
        { duration: '2m', target: 150 },
        { duration: '3m', target: 150 },
        { duration: '2m', target: 200 },
        { duration: '3m', target: 200 },
        { duration: '2m', target: 250 },
        { duration: '3m', target: 250 },
        { duration: '2m', target: 300 },
        { duration: '3m', target: 300 },
        { duration: '2m', target: 350 },
        { duration: '3m', target: 350 },
        { duration: '2m', target: 400 },
        { duration: '5m', target: 400 },
      ],
      tags: { test_type: 'breakpoint' },
    },
  },
};

export const options = {
  scenarios: scenarios[SCENARIO] || scenarios.load,

  thresholds: {
    http_req_duration: ['p(95)<1000', 'p(99)<2000'],
    http_req_failed: ['rate<0.10'],
    login_errors: ['rate<0.05'],
    business_errors: ['rate<0.05'],
  },
};

// ── SETUP: Login inicial ───────────────────────────────
export function setup() {
  let token = null;
  let expiresIn = 0;
  let attempts = 0;
  const maxAttempts = 3;

  while (!token && attempts < maxAttempts) {
    attempts++;
    console.log(`🔑 Intentando login (intento ${attempts}/${maxAttempts})...`);

    const loginRes = http.post(
      `${API}/auth/login`,
      JSON.stringify({ email: USERNAME, password: PASSWORD }),
      {
        headers: { 'Content-Type': 'application/json' },
        timeout: '10s',
        tags: { name: 'login' },
      }
    );

    const success = check(loginRes, {
      'login status 200': (r) => r.status === 200,
      'login response time < 2s': (r) => r.timings.duration < 2000,
    });

    if (success) {
      try {
        const body = JSON.parse(loginRes.body);
        token = body.data?.access_token;
        expiresIn = body.data?.expires_in || 900; // default 15 min
        console.log(`✅ Login OK en intento ${attempts} — token expira en ${expiresIn}s`);
        break;
      } catch (e) {
        console.error('❌ Error parseando respuesta de login:', e);
      }
    } else {
      console.warn(`⚠️ Login falló (status ${loginRes.status}): ${loginRes.body?.substring(0, 100)}`);
      sleep(2);
    }
  }

  if (!token) {
    console.error('❌❌ Login falló después de todos los intentos. Abortando tests autenticados.');
  }

  const now = Math.floor(Date.now() / 1000);
  return {
    token,
    loginSuccess: !!token,
    expiresIn,
    tokenObtainedAt: now,
  };
}

// ── FUNCIONES AUXILIARES ───────────────────────────────

function getRandomBogotaCoords() {
  const lat = randomIntBetween(46000, 48300) / 10000;
  const lng = randomIntBetween(-74200, -74000) / 10000;
  return { lat, lng };
}

function categorizeError(status) {
  if (status >= 500) {
    businessErrors.add(1);
    return 'server_error';
  } else if (status >= 400) {
    return 'client_error';
  }
  return 'success';
}

// ⭐ Refresh PROACTIVO de token (por tiempo, no por 401)
// Cada VU refresca su propio token antes de que expire
function ensureFreshToken(data) {
  const now = Math.floor(Date.now() / 1000);

  // Si el VU ya tiene un token vigente, usarlo
  if (vuToken && now < vuTokenExpiresAt) {
    return vuToken;
  }

  // Si el token del setup aún es vigente, usarlo
  const setupTokenExpiresAt = data.tokenObtainedAt + data.expiresIn - TOKEN_REFRESH_MARGIN_S;
  if (!vuToken && now < setupTokenExpiresAt) {
    return data.token;
  }

  // Token expirado o por expirar → refrescar
  const loginRes = http.post(
    `${API}/auth/login`,
    JSON.stringify({ email: USERNAME, password: PASSWORD }),
    {
      headers: { 'Content-Type': 'application/json' },
      timeout: '10s',
      tags: { name: 'token_refresh' },
    }
  );

  if (loginRes.status === 200) {
    try {
      const body = JSON.parse(loginRes.body);
      vuToken = body.data?.access_token;
      const expiresIn = body.data?.expires_in || 900;
      vuTokenExpiresAt = now + expiresIn - TOKEN_REFRESH_MARGIN_S;
      tokenRefreshes.add(1);
      return vuToken;
    } catch (e) {
      console.warn('⚠️ Error parseando refresh response:', e);
    }
  }

  // Si el refresh falla (rate limit, etc), usar el token que tengamos
  return vuToken || data.token;
}

// ── MAIN TEST FUNCTION ─────────────────────────────────
export default function mainTest(data) {
  // ⭐ Obtener token fresco ANTES de hacer requests
  const activeToken = data.loginSuccess ? ensureFreshToken(data) : null;
  const authHeaders = activeToken
    ? { headers: { Authorization: `Bearer ${activeToken}` } }
    : {};

  throughput.add(1);

  // ==========================================
  // 1. HEALTH CHECK CON ABORT AUTOMÁTICO
  // ==========================================
  group('01_Health', () => {
    const res = http.get(`${BASE_URL}/health`, {
      timeout: '5s',
      tags: { name: 'health_check' },
    });
    const ok = check(res, {
      'health 200': (r) => r.status === 200,
      'health fast < 500ms': (r) => r.timings.duration < 500,
    });

    if (ok) {
      consecutiveHealthFailures = 0;
    } else {
      consecutiveHealthFailures++;
      console.warn(
        `💔 Health check falló (${consecutiveHealthFailures}/${MAX_HEALTH_FAILURES}): ` +
        `status=${res.status} (${res.timings.duration}ms)`
      );
      categorizeError(res.status);

      if (consecutiveHealthFailures >= MAX_HEALTH_FAILURES) {
        console.error(
          `🚨🚨🚨 ABORT: El servidor no responde después de ${MAX_HEALTH_FAILURES} intentos consecutivos. ` +
          `Abortando test para evitar miles de errores innecesarios.`
        );
        exec.test.abort(
          `Server health check failed ${MAX_HEALTH_FAILURES} consecutive times — server likely crashed`
        );
      }
    }
  });

  sleep(Math.random() * 0.5 + 0.2);

  // ==========================================
  // 2. CATÁLOGOS (70% de los usuarios)
  // ⭐ Tags agrupan URLs para evitar explosión de métricas
  // ==========================================
  if (Math.random() < 0.7) {
    group('02_Catalogs', () => {
      const endpoints = [
        { path: '/brands', tag: 'catalog_brands' },
        { path: '/branch-types', tag: 'catalog_branch_types' },
        { path: '/services', tag: 'catalog_services' },
      ];

      for (const ep of endpoints) {
        const res = http.get(`${API}${ep.path}`, {
          timeout: '10s',
          tags: { name: ep.tag },
        });
        catalogTime.add(res.timings.duration);

        const ok = check(res, {
          [`${ep.path} status 200`]: (r) => r.status === 200,
          [`${ep.path} response < 1s`]: (r) => r.timings.duration < 1000,
        });

        if (!ok) {
          console.warn(`📚 Catalog ${ep.path} falló: ${res.status} (${res.timings.duration}ms)`);
          categorizeError(res.status);
        }

        sleep(Math.random() * 0.3 + 0.1);
      }
    });
  }

  sleep(Math.random() * 0.5 + 0.3);

  // ==========================================
  // 3. ENDPOINTS AUTENTICADOS (solo si hay token)
  // ==========================================
  if (activeToken && data.loginSuccess) {

    // 3.1 Perfil (30% de los usuarios)
    if (Math.random() < 0.3) {
      group('03_Profile', () => {
        const res = http.get(`${API}/persons/me`, {
          ...authHeaders,
          timeout: '5s',
          tags: { name: 'profile_get' },
        });
        authTime.add(res.timings.duration);

        const ok = check(res, {
          'persons/me 200': (r) => r.status === 200,
          'profile < 800ms': (r) => r.timings.duration < 800,
        });

        if (!ok) {
          console.warn(`👤 Profile falló: ${res.status}`);
          categorizeError(res.status);
        }
      });
      sleep(Math.random() * 0.5 + 0.2);
    }

    // 3.2 Motocicletas del usuario (40% de los usuarios)
    if (Math.random() < 0.4) {
      group('04_Motorcycles', () => {
        const res = http.get(`${API}/motorcycles`, {
          ...authHeaders,
          timeout: '5s',
          tags: { name: 'motorcycles_list' },
        });
        authTime.add(res.timings.duration);

        const ok = check(res, {
          'motorcycles 200/204': (r) => r.status === 200 || r.status === 204,
          'motorcycles < 800ms': (r) => r.timings.duration < 800,
        });

        if (!ok) {
          console.warn(`🏍️ Motorcycles falló: ${res.status}`);
          categorizeError(res.status);
        }
      });
      sleep(Math.random() * 0.5 + 0.3);
    }

    // 3.3 BÚSQUEDA GEO (la más importante - 80% de usuarios)
    // ⭐ Tag 'nearby_search' agrupa TODAS las coordenadas en UNA sola métrica
    if (Math.random() < 0.8) {
      group('05_Nearby_Search', () => {
        const coords = getRandomBogotaCoords();
        const radius = randomIntBetween(5, 20);

        const res = http.get(
          `${API}/branches/nearby?lat=${coords.lat}&lng=${coords.lng}&radius=${radius}`,
          {
            ...authHeaders,
            timeout: '15s',
            tags: { name: 'nearby_search' }, // ⭐ Agrupa todas las variantes geo
          }
        );
        geoQueryTime.add(res.timings.duration);

        const ok = check(res, {
          'nearby 200/204': (r) => r.status === 200 || r.status === 204,
          'nearby < 2s': (r) => r.timings.duration < 2000,
        });

        if (!ok) {
          console.warn(
            `📍 Nearby falló: ${res.status} (${res.timings.duration}ms) ` +
            `- coords: ${coords.lat},${coords.lng}`
          );
          categorizeError(res.status);
        } else if (res.timings.duration > 1000) {
          console.log(`🐌 Query geo lenta: ${res.timings.duration}ms`);
        }
      });
    }

  } else {
    loginErrors.add(1);
    sleep(1);
  }

  // ⭐ THINK TIME REALISTA (tiempo entre acciones del usuario)
  const rand = Math.random();
  let thinkTime;
  if (rand < 0.6) {
    thinkTime = Math.random() * 2 + 1;
  } else if (rand < 0.9) {
    thinkTime = Math.random() * 2 + 3;
  } else {
    thinkTime = Math.random() * 5 + 5;
  }
  sleep(thinkTime);
}

// ── TEARDOWN ───────────────────────────────────────────
export function teardown(data) {
  console.log('\n🏁 RESUMEN DE PRUEBA');
  console.log('===================');

  if (data.loginSuccess) {
    console.log('✅ Autenticación: OK');
  } else {
    console.log('❌ Autenticación: FALLÓ - Solo se probaron endpoints públicos');
  }

  console.log(`📊 Escenario ejecutado: ${SCENARIO}`);
  console.log('📈 Revisa Grafana/dashboard para métricas detalladas');
  console.log('💡 Tip: Si ves muchos errores 5xx, revisa logs del servidor Go');
}
