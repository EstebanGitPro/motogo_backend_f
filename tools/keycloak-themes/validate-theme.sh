#!/bin/bash

# Verificador de estructura del theme de Keycloak
# Valida que todos los archivos necesarios existan

set -e

THEME_PATH="./tools/keycloak-themes/motogo"

echo "🔍 Verificando estructura del theme de Keycloak..."
echo ""

# Función para verificar archivos
check_file() {
    if [ -f "$1" ]; then
        echo "✅ $1"
        return 0
    else
        echo "❌ FALTA: $1"
        return 1
    fi
}

# Función para verificar directorios
check_dir() {
    if [ -d "$1" ]; then
        echo "✅ $1/"
        return 0
    else
        echo "❌ FALTA: $1/"
        return 1
    fi
}

ERRORS=0

# Verificar estructura base
echo "📁 Estructura base:"
check_dir "$THEME_PATH" || ((ERRORS++))
check_file "$THEME_PATH/theme.properties" || ((ERRORS++))
check_dir "$THEME_PATH/email" || ((ERRORS++))
check_file "$THEME_PATH/email/theme.properties" || ((ERRORS++))
echo ""

# Verificar directorios principales
echo "📂 Directorios principales:"
check_dir "$THEME_PATH/email/html" || ((ERRORS++))
check_dir "$THEME_PATH/email/text" || ((ERRORS++))
check_dir "$THEME_PATH/email/messages" || ((ERRORS++))
check_dir "$THEME_PATH/email/resources" || ((ERRORS++))
echo ""

# Verificar messages
echo "🌍 Messages (i18n):"
check_file "$THEME_PATH/email/messages/messages_es.properties" || ((ERRORS++))
echo ""

# Verificar templates HTML
echo "📧 Templates HTML:"
check_file "$THEME_PATH/email/html/template.ftl" || ((ERRORS++))
check_file "$THEME_PATH/email/html/email-verification.ftl" || ((ERRORS++))
check_file "$THEME_PATH/email/html/password-reset.ftl" || ((ERRORS++))
check_file "$THEME_PATH/email/html/executeActions.ftl" || ((ERRORS++))
check_file "$THEME_PATH/email/html/event-update_password.ftl" || ((ERRORS++))
check_file "$THEME_PATH/email/html/email-test.ftl" || ((ERRORS++))
echo ""

# Verificar templates texto plano
echo "📝 Templates Texto Plano:"
check_file "$THEME_PATH/email/text/email-verification.ftl" || ((ERRORS++))
check_file "$THEME_PATH/email/text/password-reset.ftl" || ((ERRORS++))
check_file "$THEME_PATH/email/text/executeActions.ftl" || ((ERRORS++))
check_file "$THEME_PATH/email/text/event-update_password.ftl" || ((ERRORS++))
check_file "$THEME_PATH/email/text/email-test.ftl" || ((ERRORS++))
echo ""

# Resumen
echo "══════════════════════════════════════════════"
if [ $ERRORS -eq 0 ]; then
    echo "✅ VALIDACIÓN EXITOSA"
    echo "   Todos los archivos necesarios están presentes"
    echo "   El theme está listo para instalarse"
    echo ""
    echo "💡 Para instalar, ejecuta:"
    echo "   ./tools/keycloak-themes/install-theme.sh"
    exit 0
else
    echo "❌ ERRORES ENCONTRADOS: $ERRORS"
    echo "   Revisa los archivos faltantes arriba"
    exit 1
fi
