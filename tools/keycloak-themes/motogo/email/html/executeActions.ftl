<#import "template.ftl" as layout>
<@layout.emailLayout>
    <h2>Acción requerida en tu cuenta 📋</h2>
    
    <p>Hola<#if user.firstName??> ${user.firstName}</#if>,</p>
    
    <p>Necesitamos que completes algunas acciones en tu cuenta de <strong>MotoGo</strong> para continuar usando nuestros servicios.</p>
    
    <#if requiredActions??>
    <div class="info-box">
        <p><strong>Acciones pendientes:</strong></p>
        <ul style="margin: 8px 0; padding-left: 20px;">
        <#list requiredActions as action>
            <li>${action}</li>
        </#list>
        </ul>
    </div>
    </#if>
    
    <div style="text-align: center;">
        <a href="${link}" class="button">
            📝 Completar acciones
        </a>
    </div>
    
    <div class="info-box">
        <p>⏰ <strong>Este enlace expira en <#if linkExpiration??>${linkExpiration}<#else>24 horas</#if>.</strong></p>
    </div>
    
    <p>Si el botón no funciona, puedes copiar y pegar este enlace en tu navegador:</p>
    <p style="word-break: break-all; font-size: 14px; color: #6b7280;">${link}</p>
    
    <div class="alert-box">
        <p>🛡️ Si no reconoces esta solicitud, contacta a nuestro equipo de soporte de inmediato.</p>
    </div>
</@layout.emailLayout>
