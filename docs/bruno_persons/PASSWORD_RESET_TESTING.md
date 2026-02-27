# Testing Password Reset Flow

## Prerrequisitos

1. **Keycloak debe tener SMTP configurado**:
   - Ve a Keycloak Admin Console → Realm Settings → Email
   - Configura el servidor SMTP (Gmail, Mailtrap, etc.)
   - Prueba el email con "Test connection"

2. **Theme de MotoGo debe estar aplicado**:
   - El realm debe usar el theme "motogo" para emails
   - Verificar en Realm Settings → Themes → Email theme

3. **Usuario debe existir en el sistema**:
   - Registra un usuario primero (usa `01_register_person.bru`)
   - El email debe ser válido para recibir el correo

## Flujo de Testing End-to-End

### Paso 1: Registrar Usuario (Si no existe)
```bash
POST /accounts
{
  "email": "test@example.com",
  "first_name": "Test",
  "last_name": "User",
  "password": "OldPassword123",
  ...
}
```

### Paso 2: Solicitar Password Reset
```bash
POST /auth/password-reset
{
  "email": "test@example.com"
}
```

**Respuesta esperada**: 200 OK
```json
{
  "success": true,
  "message": {
    "codigo_mensaje": "MOD_P_RESET_INFO_00001",
    "titulo_mensaje": "Email de recuperación enviado",
    "contenido_mensaje": "Si el email existe en nuestro sistema, recibirás un correo..."
  }
}
```

### Paso 3: Revisar Email

1. **Abre tu bandeja de entrada** del email registrado
2. **Busca el email** con asunto: "Recuperación de Contraseña" o "Acción requerida en tu cuenta"
3. **El email debe tener**:
   - Botón "Restablecer contraseña" o "Completar acciones"
   - Link alternativo si el botón no funciona

### Paso 4: Extraer el Token

El link del email apunta a:
```
http://localhost:8085/reset-password.html?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Extrae solo el token** (la parte después de `?token=`)

### Paso 5: Reset Password Con El Token

```bash
POST /auth/password/reset
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "new_password": "NewSecurePassword123"
}
```

**Respuesta esperada**: 200 OK
```json
{
  "success": true,
  "message": {
    "codigo_mensaje": "MOD_P_RESET_EXI_00001",
    "titulo_mensaje": "Contraseña actualizada",
    "contenido_mensaje": "Tu contraseña ha sido actualizada exitosamente."
  }
}
```

### Paso 6: Verificar Nueva Contraseña

Intenta login con la **nueva contraseña**:

```bash
POST /auth/login
{
  "email": "test@example.com",
  "password": "NewSecurePassword123"
}
```

**Debe funcionar** ✅

Intenta login con la **contraseña antigua**:

```bash
POST /auth/login
{
  "email": "test@example.com",
  "password": "OldPassword123"
}
```

**Debe fallar** ❌

## Casos de Prueba

### ✅ Caso Exitoso
- Token válido
- Usuario existe
- Password cumple requisitos (min 8 caracteres)
- **Resultado**: 200 OK, contraseña actualizada

### ❌ Token Inválido
```json
{
  "token": "token-invalido",
  "new_password": "NewPassword123"
}
```
**Resultado**: 400 Bad Request, `MOD_P_RESET_ERR_00001`

### ❌ Token Expirado
- Espera más de 12 horas después de solicitar reset
- **Resultado**: 400 Bad Request, `MOD_P_RESET_ERR_00001`

### ❌ Usuario No Encontrado
- Token válido pero email no existe en Keycloak
- **Resultado**: 404 Not Found, `MOD_P_RESET_ERR_00002`

### ❌ Password Muy Corta
```json
{
  "token": "token-valido",
  "new_password": "123"
}
```
**Resultado**: 400 Bad Request (validación de Gin)

## Troubleshooting

### No recibo el email

**Verificar**:
1. SMTP configurado en Keycloak
2. Email del usuario es válido
3. Revisar spam/correo no deseado
4. Logs de Keycloak para ver errores de envío

### Token inválido

**Verificar**:
1. Token completo (sin espacios o saltos de línea)
2. Token no expirado (menos de 12 horas)
3. Token generado por Keycloak (no editado manualmente)

### Link del email apunta a Keycloak en lugar de nuestra página

**Verificar**:
1. Theme "motogo" aplicado al realm
2. Archivo `executeActions.ftl` modificado correctamente
3. Keycloak reiniciado después de cambiar theme

## Configuración de SMTP para Testing

### Opción 1: Gmail (Para desarrollo)
```
Host: smtp.gmail.com
Port: 587
From: tu-email@gmail.com
Username: tu-email@gmail.com
Password: [App Password]
Enable TLS: Yes
```

### Opción 2: Mailtrap (Recomendado para testing)
```
Host: smtp.mailtrap.io
Port: 2525
Username: [Tu username de Mailtrap]
Password: [Tu password de Mailtrap]
```

Mailtrap captura todos los emails sin enviarlos realmente.

### Opción 3: MailHog (Local)
```
Host: localhost
Port: 1025
```

MailHog es un servidor SMTP local que captura emails.

## Logs para Debug

Backend logs mostrarán:
```
[INFO] Iniciando proceso de recuperación de contraseña
[DEBUG] Email extraído exitosamente del token: test@example.com
[DEBUG] Usuario encontrado para reset de contraseña
[SUCCESS] Contraseña actualizada exitosamente
```

Si hay error:
```
[ERROR] Error extrayendo email del token de reset
[ERROR] Usuario no encontrado para reset de contraseña
[ERROR] Error actualizando contraseña en Keycloak
```
