# 📧 Cómo Probar la Verificación de Email - Guía Paso a Paso

Esta guía te explica **cómo obtener el token** necesario para probar el endpoint de verificación de email.

---

## 🎯 Flujo Completo de Prueba

### Paso 1: Registrar un Usuario

1. Abre Bruno y ve a la carpeta `bruno_persons`
2. Ejecuta el request **"Create Person"**
3. Usa un email que puedas revisar (por ejemplo: `test@yopmail.com`)
4. Guarda el email que usaste

**Ejemplo de body**:

```json
{
  "identity_number": "1234567890",
  "first_name": "Juan",
  "last_name": "Pérez",
  "second_last_name": "García",
  "email": "test@yopmail.com",
  "phone_number": "3001112233",
  "password": "Secret123",
  "role": "user"
}
```

### Paso 2: Obtener el Token de Verificación

Tienes **3 opciones** para obtener el token:

---

## 📝 **Opción 1: Desde el Correo (RECOMENDADA PARA PRODUCCIÓN)**

1. **Ve al correo** que usaste en el registro (por ejemplo: https://yopmail.com)
2. **Busca el email** de verificación de MotoGo
3. **Encuentra el link** que dice algo como:
   ```
   http://localhost:8000/verify-email?token=eyJhbGciOiJSUzI1NiIsInR5cC...
   ```
4. **Copia TODO** lo que viene después de `token=`
5. Ese es tu token JWT ✅

**Ejemplo de token**:

```
eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJYc3Iw...
```

---

## 🖥️ **Opción 2: Desde los Logs del Servidor (DESARROLLO)**

1. **Ve al terminal** donde está corriendo el backend
2. **Busca el log** después de crear el usuario, verás algo como:
   ```json
   {"level":"INFO","msg":"Email de verificación enviado exitosamente"}
   ```
3. En los logs anteriores o posteriores debería aparecer el **token completo**
4. **Copia el token** desde los logs

> **Nota**: Esta opción solo funciona si el backend está configurado para loggear los tokens (generalmente en modo desarrollo)

---

## 🧪 **Opción 3: Usar Token de Login (TESTING RÁPIDO)**

Si solo quieres probar que el endpoint funciona:

1. **Ejecuta el request "Login"** en Bruno con las credenciales del usuario que creaste:

   ```json
   {
     "email": "test@yopmail.com",
     "password": "Secret123"
   }
   ```
2. **Copia el `access_token`** de la respuesta:

   ```json
   {
     "data": {
       "access_token": "eyJhbGciOiJSUzI1NiIsInR5cC...",
       "refresh_token": "...",
       "expires_in": 300,
       "token_type": "Bearer"
     }
   }
   ```
3. **Usa ese `access_token`** como el token de verificación

> **Nota**: Esta opción funciona porque el sistema extrae el email del token JWT, y tanto el token de login como el de verificación contienen el email.

---

### Paso 3: Verificar el Email en Bruno

1. Abre el request **"Verify Email"** en Bruno
2. En el body JSON, **pega el token** que obtuviste:

   ```json
   {
     "token": "PEGA_AQUI_EL_TOKEN_QUE_COPIASTE"
   }
   ```
3. **Ejecuta el request** (Send)

---

## ✅ Respuestas Esperadas

### Primera vez que verificas el email:

```json
{
  "success": true,
  "message": {
    "codigo": "...",
    "titulo": "Email Verificado",
    "contenido": "Su email ha sido verificado exitosamente."
  },
  "data": {
    "verified": true,
    "email": "test@yopmail.com"
  }
}
```

### Si intentas verificar un email que ya estaba verificado:

```json
{
  "success": false,
  "message": {
    "codigo": "MOD_KC_EMAIL_ALREADY_VERIFIED_WARN_00001",
    "titulo": "Email Ya Verificado",
    "contenido": "Su correo electrónico ya ha sido verificado anteriormente."
  }
}
```

---

## 🔍 Troubleshooting

### ❌ "Token inválido o expirado"

- **Causa**: El token expiró o está mal formado
- **Solución**: Genera un nuevo token haciendo login nuevamente o registrando un usuario nuevo

### ❌ "Usuario no encontrado"

- **Causa**: El email del token no existe en la base de datos
- **Solución**: Asegúrate de haber registrado el usuario primero

### ❌ "Error de formato JSON"

- **Causa**: El body del request no está bien formado
- **Solución**: Verifica que el JSON sea válido y que el campo `token` esté presente

---

## 🎓 Ejemplo Completo

```bash
# 1. Registrar usuario
POST http://localhost:8000/api/v1/motogo/accounts
{
  "email": "test@yopmail.com",
  "password": "Secret123",
  ...
}

# 2. Hacer login para obtener token
POST http://localhost:8000/api/v1/motogo/auth/login
{
  "email": "test@yopmail.com",
  "password": "Secret123"
}

# Respuesta:
{
  "data": {
    "access_token": "eyJhbGc..."  <-- COPIA ESTO
  }
}

# 3. Verificar email con el token
POST http://localhost:8000/api/v1/motogo/auth/verify-email
{
  "token": "eyJhbGc..."  <-- PEGA AQUÍ
}
```

---

## 📋 Checklist de Prueba

- [ ] Usuario registrado correctamente
- [ ] Token obtenido (de correo, logs o login)
- [ ] Token pegado en el body del request "Verify Email"
- [ ] Request ejecutado exitosamente
- [ ] Respuesta con `"verified": true` recibida
- [ ] Segunda ejecución retorna "Email ya verificado"

---

**¡Listo!** Ahora ya sabes cómo probar la verificación de email 🎉
