import http from 'k6/http';
import { check, sleep } from 'k6';

// Load Test: Simular carga normal de producción
export const options = {
  stages: [
    { duration: '2m', target: 10 },  // Ramp-up: 0 → 10 usuarios en 2 min
    { duration: '5m', target: 10 },  // Mantener: 10 usuarios por 5 min
    { duration: '2m', target: 20 },  // Escalar: 10 → 20 usuarios en 2 min
    { duration: '5m', target: 20 },  // Mantener: 20 usuarios por 5 min
    { duration: '2m', target: 0 },   // Ramp-down: 20 → 0 en 2 min
  ],
  thresholds: {
    http_req_duration: ['p(95)<1000'], // 95% de requests < 1s
    http_req_failed: ['rate<0.05'],    // < 5% de errores
  },
};

const BASE_URL = 'http://host.docker.internal:8085/motogo/api/v1';

export default function () {
  // Test 1: Listar mensajes
  const messagesRes = http.get(`${BASE_URL}/messages`);
  
  check(messagesRes, {
    'messages status is 200': (r) => r.status === 200,
    'messages response < 1s': (r) => r.timings.duration < 1000,
  });

  sleep(1); // Pausa entre requests

  // Test 2: Obtener métricas
  const metricsRes = http.get('http://host.docker.internal:8085/metrics');
  
  check(metricsRes, {
    'metrics endpoint works': (r) => r.status === 200,
  });

  sleep(2); // Simular tiempo de "pensamiento" del usuario
}

export function setup() {
  console.log('🔄 Starting Load Test...');
  console.log('📊 Target: 10-20 concurrent users');
  console.log('⏱️  Duration: ~18 minutes');
}

export function teardown(data) {
  console.log('✅ Load Test completed!');
}
