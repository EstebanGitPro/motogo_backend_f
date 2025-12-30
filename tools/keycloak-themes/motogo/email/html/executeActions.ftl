<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>MotoGo - Acción Requerida</title>
</head>
<body style="font-family: Arial, sans-serif; background-color: #f2f2f2; margin: 0; padding: 2rem;">
    <div style="max-width: 600px; margin: 0 auto;">
        <div style="background: white; padding: 2rem; border-radius: 10px; box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);">
            
            <!-- Header -->
            <div style="text-align: center; margin-bottom: 2rem; padding-bottom: 1rem; border-bottom: 2px solid #007BFF;">
                <div style="color: #007BFF; font-size: 18px; font-weight: bold;">MotoGo</div>
            </div>
            
            <!-- Content -->
            <div style="color: #333; line-height: 1.6;">
                <h1 style="color: #333; text-align: center; margin-bottom: 1.5rem; font-size: 22px;">Acción requerida en tu cuenta</h1>
                
                <p style="margin: 1rem 0; color: #666;">Hola<#if user.firstName??> ${user.firstName}</#if>,</p>
                
                <p style="margin: 1rem 0; color: #666;">Necesitamos que completes algunas acciones en tu cuenta de MotoGo para continuar.</p>
                
                <#if requiredActions??>
                <div style="background-color: #e7f3ff; border-left: 4px solid #007BFF; padding: 1rem; margin: 1rem 0; border-radius: 4px;">
                    <p style="margin: 0; color: #004085; font-size: 14px;"><strong>Acciones pendientes:</strong></p>
                    <ul style="margin: 8px 0; padding-left: 20px; color: #004085;">
                    <#list requiredActions as action>
                        <li>
                            <#if action == "VERIFY_EMAIL">
                                Verificar correo electrónico
                            <#elseif action == "UPDATE_PASSWORD">
                                Actualizar contraseña
                            <#elseif action == "CONFIGURE_TOTP">
                                Configurar autenticación de dos factores
                            <#elseif action == "UPDATE_PROFILE">
                                Actualizar perfil
                            <#elseif action == "TERMS_AND_CONDITIONS">
                                Aceptar términos y condiciones
                            <#else>
                                ${action}
                            </#if>
                        </li>
                    </#list>
                    </ul>
                </div>
                </#if>
                
                <#-- 
                    CUSTOM LOGIC FOR PASSWORD RESET:
                    If action is UPDATE_PASSWORD, redirect to our custom page instead of Keycloak
                    Extract the token from the Keycloak link and pass it to our page
                -->
                <#assign customLink = link>
                <#if requiredActions?? && requiredActions?seq_contains("UPDATE_PASSWORD")>
                    <#-- Extract token from Keycloak link: http://keycloak/realms/realm/login-actions/action-token?key=TOKEN -->
                    <#assign customLink = "http://localhost:8085/reset-password.html?token=" + link?replace(".*key=", "", "r")>
                </#if>
                
                <div style="text-align: center; margin: 1.5rem 0;">
                    <a href="${customLink}" style="display: inline-block; padding: 12px 24px; background-color: #007BFF; color: white; text-decoration: none; border-radius: 5px; font-weight: 600;">
                        Completar acciones
                    </a>
                </div>
                
                <div style="background-color: #fff3cd; border-left: 4px solid #ffc107; padding: 1rem; margin: 1rem 0; border-radius: 4px;">
                    <p style="margin: 0; color: #856404; font-size: 14px;"><strong>IMPORTANTE:</strong> Este enlace expira en 15 minutos.</p>
                </div>
                
                <p style="color: #666; font-size: 14px; word-break: break-all;">Si el botón no funciona, copia y pega este enlace en tu navegador:</p>
                <p style="color: #666; font-size: 14px; word-break: break-all;">${link}</p>
                
                <div style="background-color: #e7f3ff; border-left: 4px solid #007BFF; padding: 1rem; margin: 1rem 0; border-radius: 4px;">
                    <p style="margin: 0; color: #004085; font-size: 14px;"><strong>Nota:</strong> Si no reconoces esta solicitud, contacta a nuestro equipo de soporte de inmediato.</p>
                </div>
            </div>
            
            <!-- Footer -->
            <div style="margin-top: 2rem; padding-top: 1.5rem; border-top: 1px solid #e5e7eb; text-align: center; font-size: 12px; color: #999;">
                <p><a href="mailto:soporte@motogo.com" style="color: #007BFF; text-decoration: none;">soporte@motogo.com</a></p>
                <p>&copy; ${.now?string('yyyy')} MotoGo. Todos los derechos reservados.</p>
            </div>
            
        </div>
    </div>
</body>
</html>
