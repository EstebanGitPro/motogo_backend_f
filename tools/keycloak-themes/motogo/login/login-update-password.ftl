<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>MotoGo - Cambiar Contraseña</title>
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
            background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
            border-radius: 50%;
            display: flex; align-items: center; justify-content: center;
            margin: 0 auto 1.5rem;
            font-size: 40px; color: white;
        }
        h1 { color: #333; font-size: 22px; margin-bottom: 1rem; }
        .message { color: #666; font-size: 16px; line-height: 1.6; margin-bottom: 1.5rem; }
        .app-box {
            background: linear-gradient(135deg, #fef3c7 0%, #fde68a 100%);
            border: 2px solid #f59e0b;
            border-radius: 12px;
            padding: 1.5rem;
            margin-bottom: 1.5rem;
        }
        .app-box h2 { color: #92400e; font-size: 18px; margin-bottom: 0.5rem; }
        .app-box p { color: #a16207; font-size: 14px; margin: 0; }
        .app-icon { font-size: 48px; margin-bottom: 0.5rem; }
        .footer { margin-top: 1.5rem; color: #999; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="card">
            <div class="logo">MotoGo</div>
            <div class="icon">🔐</div>
            
            <h1>Cambio de contraseña</h1>
            
            <p class="message">
                Para cambiar tu contraseña, utiliza la aplicación móvil de MotoGo.
            </p>
            
            <div class="app-box">
                <div class="app-icon">📱</div>
                <h2>Abre la aplicación MotoGo</h2>
                <p>Ve a "Configuración" → "Cambiar contraseña" para actualizar tu contraseña de forma segura.</p>
            </div>
            
            <p class="message" style="color: #f59e0b; font-weight: 600; margin-bottom: 0;">
                Puedes cerrar esta ventana.
            </p>
            
            <div class="footer">
                <p>© ${.now?string('yyyy')} MotoGo</p>
            </div>
        </div>
    </div>
</body>
</html>
