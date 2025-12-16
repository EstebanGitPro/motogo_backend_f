<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>MotoGo - Contraseña Actualizada</title>
</head>
<body style="font-family: Arial, sans-serif; background-color: #f2f2f2; margin: 0; padding: 2rem;">
    <div style="max-width: 600px; margin: 0 auto;">
        <div style="background: white; padding: 2rem; border-radius: 10px; box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);">
            
            <!-- Header -->
            <div style="text-align: center; margin-bottom: 2rem; padding-bottom: 1rem; border-bottom: 2px solid #007BFF;">
                <div style="color: #007BFF; font-size: 18px; font-weight: bold;">MotoGo</div>
            </div>
            
            <!-- Content -->
            <div style="color: #333; line-height: 1.6;">
                <h1 style="color: #333; text-align: center; margin-bottom: 1.5rem; font-size: 22px;">Contraseña Actualizada</h1>
                
                <p style="margin: 1rem 0; color: #666;">Hola<#if user.firstName??> ${user.firstName}</#if>,</p>
                
                <p style="margin: 1rem 0; color: #666;">Te confirmamos que tu contraseña de MotoGo ha sido actualizada exitosamente.</p>
                
                <div style="background-color: #e7f3ff; border-left: 4px solid #007BFF; padding: 1rem; margin: 1rem 0; border-radius: 4px;">
                    <p style="margin: 0; color: #004085; font-size: 14px;"><strong>Fecha:</strong> ${.now?string('dd/MM/yyyy, HH:mm')}</p>
                    <p style="margin: 0.5rem 0 0 0; color: #004085; font-size: 14px;"><strong>Cambio:</strong> Contraseña actualizada</p>
                </div>
                
                <p style="margin: 1rem 0; color: #666;">A partir de ahora, podrás iniciar sesión con tu nueva contraseña.</p>
                
                <div style="background-color: #fef2f2; border-left: 4px solid #ef4444; padding: 1rem; margin: 1rem 0; border-radius: 4px;">
                    <p style="margin: 0; color: #991b1b; font-size: 14px;"><strong>¿NO FUISTE TÚ?</strong><br>
                    Si NO realizaste este cambio, tu cuenta podría estar comprometida.<br>
                    Contacta a soporte INMEDIATAMENTE: soporte@motogo.com</p>
                </div>
                
                <p style="color: #999; font-size: 12px; margin-top: 30px;">
                    <strong>Recomendaciones de seguridad:</strong><br>
                    - Usa una contraseña única y fuerte<br>
                    - Nunca compartas tus credenciales<br>
                    - Habilita la autenticación de dos factores si está disponible
                </p>
            </div>
            
            <!-- Footer -->
            <div style="margin-top: 2rem; padding-top: 1.5rem; border-top: 1px solid #e5e7eb; text-align: center; font-size: 12px; color: #999;">
                <p><a href="mailto:soporte@motogo.com" style="color: #007BFF; text-decoration: none;">soporte@motogo.com</a></p>
                <p>&copy; ${.now?string('yyyy')} MotoGo. Todos los derechos reservados.</p>
            </div>
            
        </div>
    </div>
</body>
</html>
