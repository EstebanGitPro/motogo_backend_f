import http from 'k6/http';
import { check, sleep } from 'k6';

// Realistic Test: Simular flujo real de usuario
export const options = {
  stages: [
    { duration: '1m', target: 5 },   // Ramp-up: 5 usuarios
    { duration: '3m', target: 5 },   // Mantener: 5 usuarios
    { duration: '1m', target: 0 },   // Ramp-down
  ],
  thresholds: {
    http_req_duration: ['p(95)<1000'],
    http_req_failed: ['rate<0.05'],
  },
};

const BASE_URL = 'http://host.docker.internal:8085/motogo/api/v1';

export default function () {
  // Flujo realista de usuario
  
  // 1. Usuario visita homepage/mensajes
  let res1 = http.get(`${BASE_URL}/messages`);
  check(res1, {
    'homepage loaded': (r) => r.status === 200,
  });
  sleep(2); // Usuario lee la página
  
  // 2. Usuario busca mensajes específicos
  let res2 = http.get(`${BASE_URL}/messages?module=users`);
  check(res2, {
    'search worked': (r) => r.status === 200,
  });
  sleep(3); // Usuario revisa resultados
  
  // 3. Usuario ve detalle de un mensaje
  let res3 = http.get(`${BASE_URL}/messages/1`);
  check(res3, {
    'detail loaded': (r) => r.status === 200 || r.status === 404,
  });
  sleep(2); // Usuario lee el detalle
  
  // TODO: Cuando tengas más endpoints, añadir:
  // 4. Usuario intenta registrarse
  // const registerData = {
  //   identity_number: `TEST${Date.now()}`,
  //   first_name: 'Test',
  //   last_name: 'User',
  //   email: `test${Date.now()}@example.com`,
  //   password: 'Test123!',
  //   phone_number: '3001234567',
  // };
  // let res4 = http.post(`${BASE_URL}/accounts`, JSON.stringify(registerData), {
  //   headers: { 'Content-Type': 'application/json' },
  // });
  // check(res4, {
  //   'registration successful': (r) => r.status === 201,
  // });
  
  sleep(5); // Tiempo entre ciclos
}

export function setup() {
  console.log('🔄 Starting Realistic Test...');
  console.log('👤 Simulating real user behavior');
  console.log('⏱️  Duration: ~5 minutes');
}

export function teardown(data) {
  console.log('✅ Realistic Test completed!');
  console.log('💡 Tip: Add more endpoints (register, login, etc.) for better tests');
}
