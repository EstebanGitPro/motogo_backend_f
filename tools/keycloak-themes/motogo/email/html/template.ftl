<#macro emailLayout>
<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>MotoGo</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f2f2f2;
            margin: 0;
            padding: 2rem;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
        }
        .card {
            background: white;
            padding: 2rem;
            border-radius: 10px;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
        }
        .header {
            text-align: center;
            margin-bottom: 2rem;
            padding-bottom: 1rem;
            border-bottom: 2px solid #007BFF;
        }
        .header h1 {
            color: #333;
            font-size: 24px;
            margin: 0.5rem 0;
        }
        .logo {
            color: #007BFF;
            font-size: 18px;
            font-weight: bold;
            margin-bottom: 0.5rem;
        }
        .content {
            color: #333;
            line-height: 1.6;
        }
        .content p {
            margin: 1rem 0;
            color: #666;
        }
        .button {
            display: inline-block;
            margin: 1.5rem 0;
            padding: 12px 24px;
            background-color: #007BFF;
            color: white !important;
            text-decoration: none;
            border-radius: 5px;
            font-weight: 600;
        }
        .button:hover {
            background-color: #0056b3;
        }
        .link-text {
            color: #666;
            font-size: 14px;
            word-break: break-all;
        }
        .footer {
            margin-top: 2rem;
            padding-top: 1.5rem;
            border-top: 1px solid #e5e7eb;
            text-align: center;
            font-size: 12px;
            color: #999;
        }
        .footer a {
            color: #007BFF;
            text-decoration: none;
        }
        .info-box {
            background-color: #e7f3ff;
            border-left: 4px solid #007BFF;
            padding: 1rem;
            margin: 1rem 0;
            border-radius: 4px;
        }
        .info-box p {
            margin: 0;
            color: #004085;
            font-size: 14px;
        }
        .warning-box {
            background-color: #fff3cd;
            border-left: 4px solid #ffc107;
            padding: 1rem;
            margin: 1rem 0;
            border-radius: 4px;
        }
        .warning-box p {
            margin: 0;
            color: #856404;
            font-size: 14px;
        }
        @media only screen and (max-width: 600px) {
            body {
                padding: 1rem;
            }
            .card {
                padding: 1.5rem;
            }
            .button {
                display: block;
                text-align: center;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="card">
            <div class="header">
                <div class="logo">MotoGo</div>
            </div>
            <div class="content">
                <#nested>
            </div>
            <div class="footer">
                <p><a href="mailto:soporte@motogo.com">soporte@motogo.com</a></p>
                <p>&copy; ${.now?string('yyyy')} MotoGo. Todos los derechos reservados.</p>
            </div>
        </div>
    </div>
</body>
</html>
</#macro>

