Acción requerida en tu cuenta - MotoGo

Hola<#if user.firstName??> ${user.firstName}</#if>,

Necesitamos que completes algunas acciones en tu cuenta de MotoGo para continuar usando nuestros servicios.

<#if requiredActions??>
Acciones pendientes:
<#list requiredActions as action>
- ${action}
</#list>
</#if>

Para completar estas acciones, visita el siguiente enlace:
${link}

⏰ Este enlace expira en <#if linkExpiration??>${linkExpiration}<#else>24 horas</#if>.

🛡️ Si no reconoces esta solicitud, contacta a nuestro equipo de soporte de inmediato.

--
MotoGo - Gestión Tributaria Inteligente
¿Necesitas ayuda? soporte@motogo.com
© ${.now?string('yyyy')} MotoGo
