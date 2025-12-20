<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>MotoGo - Error</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: Arial, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            padding: 20px;
        }
        .container { max-width: 450px; width: 100%; }
        .card {
            background: white;
            padding: 2.5rem 2rem;
            border-radius: 16px;
            box-shadow: 0 20px 40px rgba(0, 0, 0, 0.2);
            text-align: center;
        }
        .logo { color: #007BFF; font-size: 28px; font-weight: bold; margin-bottom: 1.5rem; }
        .icon {
            width: 80px; height: 80px;
            background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
            border-radius: 50%;
            display: flex; align-items: center; justify-content: center;
            margin: 0 auto 1.5rem;
            font-size: 40px; color: white;
        }
        h1 { color: #333; font-size: 22px; margin-bottom: 1rem; }
        .message { color: #666; font-size: 16px; line-height: 1.6; margin-bottom: 1.5rem; }
        .error-box {
            background: #fef2f2;
            border-left: 4px solid #ef4444;
            border-radius: 8px;
            padding: 1rem;
            margin-bottom: 1.5rem;
            text-align: left;
        }
        .error-box p { color: #991b1b; font-size: 14px; margin: 0; }
        .app-box {
            background: linear-gradient(135deg, #f0f9ff 0%, #e0f2fe 100%);
            border: 2px solid #007BFF;
            border-radius: 12px;
            padding: 1.5rem;
            margin-bottom: 1.5rem;
        }
        .app-box h2 { color: #007BFF; font-size: 18px; margin-bottom: 0.5rem; }
        .app-box p { color: #0369a1; font-size: 14px; margin: 0; }
        .app-icon { font-size: 48px; margin-bottom: 0.5rem; }
        .footer { margin-top: 1.5rem; color: #999; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="card">
            <div class="logo">MotoGo</div>
            <div class="icon">✕</div>
            
            <h1>Enlace inválido o expirado</h1>
            
            <p class="message">
                El enlace que utilizaste ya no es válido o ha expirado.
            </p>
            
            <div class="app-box">
                <div class="app-icon">📱</div>
                <h2>Solicita un nuevo enlace</h2>
                <p>Abre la aplicación MotoGo y solicita un nuevo correo de verificación.</p>
            </div>
            
            <div class="error-box">
                <p><strong>¿Por qué pasó esto?</strong></p>
                <p>• El enlace puede haber expirado (15 minutos)<br>
                   • Ya fue utilizado anteriormente<br>
                   • El enlace está incompleto</p>
            </div>
            
            <p class="message" style="margin-bottom: 0;">
                Puedes cerrar esta ventana.
            </p>
            
            <div class="footer">
                <p>© ${.now?string('yyyy')} MotoGo</p>
            </div>
        </div>
    </div>
</body>
</html>
