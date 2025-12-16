<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>MotoGo - Test SMTP</title>
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
                <h1 style="color: #333; text-align: center; margin-bottom: 1.5rem; font-size: 22px;">Test de Configuración SMTP</h1>
                
                <p style="text-align: center; font-size: 18px; color: #666;">
                    Si estás viendo este mensaje, significa que MotoGo está enviando correos electrónicos correctamente.
                </p>
                
                <div style="background-color: #e7f3ff; border-left: 4px solid #007BFF; padding: 1rem; margin: 1rem 0; border-radius: 4px;">
                    <p style="margin: 0; color: #004085; font-size: 14px;"><strong>Configuración SMTP exitosa</strong></p>
                    <p style="margin: 0.5rem 0 0 0; color: #004085; font-size: 14px;">Tu servidor de correo está funcionando perfectamente con Keycloak.</p>
                </div>
                
                <div style="background: #f0fdf4; padding: 24px; border-radius: 12px; text-align: center; margin: 24px 0; border: 1px solid #10b981;">
                    <p style="color: #065f46; font-size: 20px; font-weight: 600; margin: 0;">
                        Todo está listo para enviar notificaciones
                    </p>
                </div>
                
                <p style="color: #666; font-size: 14px; text-align: center;">
                    Este es un correo de prueba generado automáticamente por Keycloak.
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
