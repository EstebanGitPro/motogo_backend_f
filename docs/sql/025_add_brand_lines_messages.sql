-- ============================================
-- HU40: Consultar Líneas de una Marca
-- Fecha: 2026-01-30
-- Módulo: Motorcycles (MOD_MOT_BRAND_LINES_*)
-- ============================================
-- Este script agrega el mensaje de éxito para listar
-- las líneas (referencias) de una marca específica.
-- Endpoint: GET /admin/brands/:brandId/lines
-- ============================================

USE motogoDb;

-- ============================================
-- INSERCIÓN DE MENSAJES - HU40
-- ============================================

INSERT INTO system_messages (
        id,
        message_code,
        type,
        category,
        module,
        message_title,
        message_content,
        is_active
    )
VALUES 
    -- Éxito: Líneas de marca obtenidas (HU40)
    (
        UUID(),
        'MOD_MOT_BRAND_LINES_EXI_00001',
        'EXITO',
        'usuario_final',
        'motorcycles',
        'Líneas de Marca Obtenidas',
        'Las líneas de motocicletas para la marca seleccionada se obtuvieron correctamente.',
        TRUE
    );

-- ============================================
-- Verificar datos insertados
-- ============================================
SELECT 'Mensaje HU40 Insertado' AS estado,
    message_code,
    message_title
FROM system_messages
WHERE message_code = 'MOD_MOT_BRAND_LINES_EXI_00001';
