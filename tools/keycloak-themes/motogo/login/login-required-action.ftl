<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>MotoGo - Redirigiendo...</title>
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
        .spinner {
            width: 50px;
            height: 50px;
            border: 4px solid #f3f3f3;
            border-top: 4px solid #007BFF;
            border-radius: 50%;
            animation: spin 1s linear infinite;
            margin: 0 auto 1.5rem;
        }
        @keyframes spin {
            0% { transform: rotate(0deg); }
            100% { transform: rotate(360deg); }
        }
        h1 { color: #333; font-size: 20px; margin-bottom: 1rem; }
        .message { color: #666; font-size: 14px; line-height: 1.6; margin-bottom: 1.5rem; }
        .button {
            display: inline-block;
            padding: 12px 24px;
            background-color: #007BFF;
            color: white;
            text-decoration: none;
            border: none;
            border-radius: 8px;
            font-weight: 600;
            cursor: pointer;
            transition: background-color 0.3s;
        }
        .button:hover { background-color: #0056b3; }
        .footer { margin-top: 1.5rem; color: #999; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="card">
            <div class="logo">MotoGo</div>
            <div class="spinner"></div>
            
            <h1>Procesando tu verificación...</h1>
            
            <p class="message">
                Estamos verificando tu correo electrónico. Serás redirigido automáticamente.
            </p>
            
            <!-- Formulario que se auto-submits para ejecutar la acción -->
            <form id="kc-action-form" action="${url.loginAction}" method="post">
                <noscript>
                    <p class="message">
                        <button type="submit" class="button">Haz clic aquí para continuar</button>
                    </p>
                </noscript>
            </form>
            
            <div class="footer">
                <p>© ${.now?string('yyyy')} MotoGo</p>
            </div>
        </div>
    </div>
    
    <script>
        // Auto-submit el formulario inmediatamente cuando la página carga
        document.addEventListener('DOMContentLoaded', function() {
            document.getElementById('kc-action-form').submit();
        });
    </script>
</body>
</html>
