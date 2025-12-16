<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>MotoGo - Recuperar Contraseña</title>
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
                <h1 style="color: #333; text-align: center; margin-bottom: 1.5rem; font-size: 22px;">Recuperación de Contraseña</h1>
                
                <p style="margin: 1rem 0; color: #666;">Hola<#if user.firstName??> ${user.firstName}</#if>,</p>
                
                <p style="margin: 1rem 0; color: #666;">Recibimos una solicitud para restablecer la contraseña de tu cuenta en MotoGo.</p>
                
                <p style="margin: 1rem 0; color: #666;">Si fuiste tú quien solicitó esto, haz clic en el siguiente botón para crear una nueva contraseña:</p>
                
                <div style="text-align: center; margin: 1.5rem 0;">
                    <a href="${link}" style="display: inline-block; padding: 12px 24px; background-color: #007BFF; color: white; text-decoration: none; border-radius: 5px; font-weight: 600;">
                        Restablecer contraseña
                    </a>
                </div>
                
                <div style="background-color: #fff3cd; border-left: 4px solid #ffc107; padding: 1rem; margin: 1rem 0; border-radius: 4px;">
                    <p style="margin: 0; color: #856404; font-size: 14px;"><strong>IMPORTANTE:</strong> Este enlace expira en 12 horas.</p>
                </div>
                
                <p style="color: #666; font-size: 14px; word-break: break-all;">Si no puedes hacer clic en el botón, copia y pega este enlace en tu navegador:</p>
                <p style="color: #666; font-size: 14px; word-break: break-all;">${link}</p>
                
                <div style="background-color: #e7f3ff; border-left: 4px solid #007BFF; padding: 1rem; margin: 1rem 0; border-radius: 4px;">
                    <p style="margin: 0; color: #004085; font-size: 14px;"><strong>¿No solicitaste esto?</strong><br>
                    No te preocupes, tu contraseña actual permanece totalmente segura. Puedes ignorar este mensaje.</p>
                </div>
                
                <p style="color: #999; font-size: 12px; margin-top: 30px;">
                    <strong>Consejo de seguridad:</strong><br>
                    Nunca compartas tu contraseña con nadie. MotoGo nunca te pedirá tu contraseña por correo electrónico.
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
