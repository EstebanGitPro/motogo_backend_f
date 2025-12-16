<#import "template.ftl" as layout>
<@layout.emailLayout>
    <h2>¡Hola<#if user.firstName??> ${user.firstName}</#if>! 👋</h2>
    
    <p>Gracias por registrarte en <strong>MotoGo</strong>, tu plataforma de gestión tributaria para empresas de transporte.</p>
    
    <p>Para activar tu cuenta y comenzar a gestionar tus obligaciones tributarias, necesitamos verificar tu correo electrónico.</p>
    
    <div style="text-align: center;">
        <a href="${link}" class="button">
            ✅ Verificar mi correo
        </a>
    </div>
    
    <div class="info-box">
        <p>⏰ <strong>Este enlace expira en 24 horas.</strong></p>
    </div>
    
    <p>Si el botón no funciona, puedes copiar y pegar este enlace en tu navegador:</p>
    <p style="word-break: break-all; font-size: 14px; color: #6b7280;">${link}</p>
    
    <div class="alert-box">
        <p>🛡️ <strong>¿No creaste esta cuenta?</strong><br>
        Puedes ignorar este mensaje de forma segura. Tu información está protegida.</p>
    </div>
    
    <p style="margin-top: 32px;">¡Bienvenido a bordo! 🏍️</p>
</@layout.emailLayout>
