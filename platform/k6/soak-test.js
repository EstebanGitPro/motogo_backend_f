import http from 'k6/http';
import { check, sleep } from 'k6';

// Soak Test: Detectar problemas de memory leaks y estabilidad
export const options = {
  stages: [
    { duration: '5m', target: 15 },   // Ramp-up: 0 → 15 usuarios en 5 min
    { duration: '1h', target: 15 },   // Mantener: 15 usuarios por 1 hora
    { duration: '5m', target: 0 },    // Ramp-down: 15 → 0 en 5 min
  ],
  thresholds: {
    http_req_duration: ['p(95)<800'],  // 95% < 800ms durante toda la hora
    http_req_failed: ['rate<0.01'],    // < 1% errores durante toda la hora
  },
};

const BASE_URL = 'http://host.docker.internal:8085/motogo/api/v1';

export default function soakTest() {
  // Simular comportamiento real de usuario
  
  // 1. Listar mensajes
  const listRes = http.get(`${BASE_URL}/messages`);
  check(listRes, {
    'list status is 200': (r) => r.status === 200,
  });
  
  sleep(2); // Usuario lee la lista
  
  // 2. Ver detalles de un mensaje
  const detailRes = http.get(`${BASE_URL}/messages/1`);
  check(detailRes, {
    'detail status is 2xx or 404': (r) => (r.status >= 200 && r.status < 300) || r.status === 404,
  });
  
  sleep(3); // Usuario lee el detalle
  
  // 3. Filtrar mensajes
  const filterRes = http.get(`${BASE_URL}/messages?module=users`);
  check(filterRes, {
    'filter status is 200': (r) => r.status === 200,
  });
  
  sleep(5); // Usuario analiza resultados
  
  // 4. Check de métricas (simular monitoreo)
  const metricsRes = http.get('http://host.docker.internal:8085/metrics');
  check(metricsRes, {
    'metrics ok': (r) => r.status === 200,
  });
  
  sleep(10); // Tiempo entre ciclos completos
}

export function setup() {
  console.log('🔄 Starting Soak Test...');
  console.log('⏱️  Duration: 1 hour 10 minutes');
  console.log('👥 Users: 15 concurrent');
  console.log('🎯 Goal: Detect memory leaks and stability issues');
}

export function teardown(data) {
  console.log('✅ Soak Test completed!');
  console.log('📊 Check for degradation in response times');
  console.log('💾 Check memory usage trends');
}
