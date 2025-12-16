<#import "template.ftl" as layout>
<@layout.emailLayout>
    <h2>Recuperación de contraseña 🔐</h2>
    
    <p>Hola<#if user.firstName??> ${user.firstName}</#if>,</p>
    
    <p>Recibimos una solicitud para restablecer la contraseña de tu cuenta en <strong>MotoGo</strong>.</p>
    
    <p>Si fuiste tú quien solicitó esto, haz clic en el siguiente botón para crear una nueva contraseña:</p>
    
    <div style="text-align: center;">
        <a href="${link}" class="button" style="background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);">
            🔑 Restablecer contraseña
        </a>
    </div>
    
    <div class="info-box">
        <p>⏰ <strong>Este enlace expira en 12 horas.</strong></p>
    </div>
    
    <p>Si el botón no funciona, puedes copiar y pegar este enlace en tu navegador:</p>
    <p style="word-break: break-all; font-size: 14px; color: #6b7280;">${link}</p>
    
    <div class="security-tip">
        <p>🛡️ <strong>¿No solicitaste esto?</strong><br>
        No te preocupes, tu contraseña actual permanece <strong>totalmente segura</strong>. Puedes ignorar este mensaje.</p>
    </div>
    
    <div class="alert-box">
        <p>💡 <strong>Consejo de seguridad:</strong><br>
        Nunca compartas tu contraseña con nadie. MotoGo <strong>nunca</strong> te pedirá tu contraseña por correo electrónico.</p>
    </div>
</@layout.emailLayout>
