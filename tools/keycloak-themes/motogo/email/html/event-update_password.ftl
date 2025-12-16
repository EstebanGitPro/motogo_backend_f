<#import "template.ftl" as layout>
<@layout.emailLayout>
    <div class="success-box" style="text-align: center;">
        <div style="font-size: 48px; margin-bottom: 16px;">✅</div>
        <h2 style="color: #065f46; margin: 0;">Contraseña actualizada</h2>
    </div>
    
    <p>Hola<#if user.firstName??> ${user.firstName}</#if>,</p>
    
    <p>Te confirmamos que tu contraseña de <strong>MotoGo</strong> ha sido actualizada exitosamente.</p>
    
    <div class="info-box">
        <p>📅 <strong>Fecha:</strong> ${.now?string('dd/MM/yyyy, HH:mm')}</p>
        <p>🔐 <strong>Cambio:</strong> Contraseña actualizada</p>
    </div>
    
    <p>A partir de ahora, podrás iniciar sesión con tu nueva contraseña.</p>
    
    <div class="security-tip">
        <p>🚨 <strong>¿No fuiste tú?</strong><br>
        Si <strong>NO</strong> realizaste este cambio, tu cuenta podría estar comprometida. 
        <strong>Contacta a soporte inmediatamente:</strong> <a href="mailto:soporte@motogo.com">soporte@motogo.com</a></p>
    </div>
    
    <div class="alert-box">
        <p>💡 <strong>Recomendaciones de seguridad:</strong></p>
        <ul style="margin: 8px 0; padding-left: 20px; color: #92400e;">
            <li>Usa una contraseña única y fuerte</li>
            <li>Nunca compartas tus credenciales</li>
            <li>Habilita la autenticación de dos factores si está disponible</li>
        </ul>
    </div>
</@layout.emailLayout>
