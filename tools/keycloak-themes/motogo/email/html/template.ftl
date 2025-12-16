<#macro emailLayout>
<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <title>MotoGo</title>
    <style>
        @import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;600;700&display=swap');
        
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Arial, sans-serif;
            line-height: 1.6;
            background-color: #f3f4f6;
            padding: 20px;
        }
        
        .email-wrapper {
            max-width: 600px;
            margin: 0 auto;
            background: #ffffff;
            border-radius: 16px;
            box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
            overflow: hidden;
        }
        
        .email-header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            padding: 48px 40px;
            text-align: center;
            color: #ffffff;
        }
        
        .email-header .logo {
            font-size: 56px;
            margin-bottom: 12px;
        }
        
        .email-header h1 {
            font-size: 32px;
            font-weight: 700;
            margin: 0;
            letter-spacing: -0.5px;
        }
        
        .email-body {
            padding: 40px;
            color: #1f2937;
        }
        
        .email-body h2 {
            font-size: 24px;
            font-weight: 700;
            color: #111827;
            margin-bottom: 16px;
        }
        
        .email-body p {
            font-size: 16px;
            color: #4b5563;
            margin-bottom: 16px;
            line-height: 1.7;
        }
        
        .button {
            display: inline-block;
            padding: 16px 32px;
            margin: 24px 0;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: #ffffff !important;
            text-decoration: none;
            border-radius: 10px;
            font-weight: 600;
            font-size: 16px;
            text-align: center;
            transition: transform 0.2s;
        }
        
        .button:hover {
            transform: translateY(-2px);
            box-shadow: 0 10px 15px -3px rgba(102, 126, 234, 0.3);
        }
        
        .alert-box {
            background: #fef3c7;
            border-left: 4px solid #f59e0b;
            padding: 16px;
            margin: 24px 0;
            border-radius: 8px;
        }
        
        .alert-box p {
            color: #92400e;
            margin: 0;
            font-size: 14px;
        }
        
        .info-box {
            background: #dbeafe;
            border-left: 4px solid #3b82f6;
            padding: 16px;
            margin: 24px 0;
            border-radius: 8px;
        }
        
        .info-box p {
            color: #1e40af;
            margin: 0;
            font-size: 14px;
        }
        
        .success-box {
            background: #d1fae5;
            border-left: 4px solid #10b981;
            padding: 16px;
            margin: 24px 0;
            border-radius: 8px;
        }
        
        .success-box p {
            color: #065f46;
            margin: 0;
            font-size: 14px;
        }
        
        .email-footer {
            background: #f9fafb;
            padding: 32px 40px;
            text-align: center;
            border-top: 1px solid #e5e7eb;
        }
        
        .email-footer p {
            font-size: 13px;
            color: #6b7280;
            margin: 8px 0;
        }
        
        .email-footer a {
            color: #667eea;
            text-decoration: none;
        }
        
        .security-tip {
            background: #fef2f2;
            border-left: 4px solid #ef4444;
            padding: 16px;
            margin: 24px 0;
            border-radius: 8px;
        }
        
        .security-tip p {
            color: #991b1b;
            margin: 0;
            font-size: 14px;
        }
        
        @media only screen and (max-width: 600px) {
            .email-wrapper {
                border-radius: 0;
            }
            
            .email-header,
            .email-body,
            .email-footer {
                padding: 24px 20px;
            }
            
            .email-header h1 {
                font-size: 24px;
            }
            
            .button {
                display: block;
                width: 100%;
            }
        }
    </style>
</head>
<body>
    <div class="email-wrapper">
        <div class="email-header">
            <div class="logo">🏍️</div>
            <h1>MotoGo</h1>
        </div>
        <div class="email-body">
            <#nested>
        </div>
        <div class="email-footer">
            <p><strong>MotoGo</strong> - Gestión Tributaria Inteligente</p>
            <p>¿Necesitas ayuda? <a href="mailto:soporte@motogo.com">soporte@motogo.com</a></p>
            <p>© ${.now?string('yyyy')} MotoGo. Todos los derechos reservados.</p>
        </div>
    </div>
</body>
</html>
</#macro>
