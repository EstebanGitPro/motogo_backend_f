<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>MotoGo - Verificar Email</title>
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
                <h1 style="color: #333; text-align: center; margin-bottom: 1.5rem; font-size: 22px;">Bienvenido a MotoGo</h1>
                
                <p style="margin: 1rem 0; color: #666;">Gracias por registrarte en nuestra plataforma.</p>
                
                <p style="margin: 1rem 0; color: #666;">Para completar tu registro, por favor haz clic en el siguiente enlace:</p>
                
                <div style="text-align: center; margin: 1.5rem 0;">
                    <a href="${link}" style="display: inline-block; padding: 12px 24px; background-color: #007BFF; color: white; text-decoration: none; border-radius: 5px; font-weight: 600;">
                        Verificar mi email
                    </a>
                </div>
                
                <p style="color: #666; font-size: 14px; word-break: break-all;">Si no puedes hacer clic en el botón, copia y pega este enlace en tu navegador:</p>
                <p style="color: #666; font-size: 14px; word-break: break-all;">${link}</p>
                
                <div style="background-color: #fff3cd; border-left: 4px solid #ffc107; padding: 1rem; margin: 1rem 0; border-radius: 4px;">
                    <p style="margin: 0; color: #856404; font-size: 14px;"><strong>IMPORTANTE:</strong> Este enlace expirará en 15 minutos por seguridad.</p>
                </div>
                
                <p style="color: #999; font-size: 12px; margin-top: 30px;">
                    Si no solicitaste este registro, puedes ignorar este correo de forma segura.
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
