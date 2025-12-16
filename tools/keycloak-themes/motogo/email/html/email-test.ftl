<#import "template.ftl" as layout>
<@layout.emailLayout>
    <div class="success-box" style="text-align: center;">
        <div style="font-size: 64px; margin-bottom: 16px;">✅</div>
        <h2 style="color: #065f46; margin: 0;">¡SMTP Configurado Correctamente!</h2>
    </div>
    
    <p style="text-align: center; font-size: 18px;">
        Si estás viendo este mensaje, significa que <strong>MotoGo</strong> está enviando correos electrónicos correctamente.
    </p>
    
    <div class="info-box">
        <p><strong>✓ Configuración SMTP exitosa</strong></p>
        <p>Tu servidor de correo está funcionando perfectamente con Keycloak.</p>
    </div>
    
    <div style="background: #f0fdf4; padding: 24px; border-radius: 12px; text-align: center; margin: 24px 0;">
        <p style="color: #065f46; font-size: 20px; font-weight: 600; margin: 0;">
            🎉 ¡Todo está listo para enviar notificaciones!
        </p>
    </div>
    
    <p style="color: #6b7280; font-size: 14px; text-align: center;">
        Este es un correo de prueba generado automáticamente por Keycloak.
    </p>
</@layout.emailLayout>

