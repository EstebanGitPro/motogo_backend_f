mysqldump: [Warning] Using a password on the command line interface can be insecure.

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
DROP TABLE IF EXISTS `branch_displacement_ranges`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `branch_displacement_ranges` (
  `id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'UUID primary key',
  `branch_id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'FK to branch',
  `displacement_range` enum('BAJO','MEDIO','ALTO') COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Displacement range: BAJO(50-200cc), MEDIO(201-400cc), ALTO(401-3000cc)',
  `active` tinyint(1) DEFAULT '1' COMMENT 'Indicates if range is active at branch',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_branch_displacement_ranges` (`branch_id`,`displacement_range`),
  CONSTRAINT `fk_branch_displacement_ranges_branch` FOREIGN KEY (`branch_id`) REFERENCES `branches` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Displacement ranges handled by each branch';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `branch_schedules`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `branch_schedules` (
  `id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'UUID primary key',
  `branch_id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'FK to branch (unique per branch)',
  `active` tinyint(1) DEFAULT '1' COMMENT 'Master switch to enable/disable entire schedule',
  `start_date` date NOT NULL DEFAULT (curdate()) COMMENT 'Schedule validity start date',
  `end_date` date DEFAULT NULL COMMENT 'Schedule validity end date (NULL = indefinite)',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_branch_schedules_branch` (`branch_id`),
  CONSTRAINT `fk_branch_schedules_branch` FOREIGN KEY (`branch_id`) REFERENCES `branches` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `chk_branch_schedules_dates` CHECK (((`end_date` is null) or (`end_date` >= `start_date`)))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Branch schedule configuration (HU30-35)';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `branch_services`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `branch_services` (
  `id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'UUID primary key',
  `branch_id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'FK to branch',
  `service_id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'FK to service',
  `active` tinyint(1) DEFAULT '1' COMMENT 'Indicates if service is active at branch',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_branch_services` (`branch_id`,`service_id`),
  KEY `fk_branch_services_service` (`service_id`),
  CONSTRAINT `fk_branch_services_branch` FOREIGN KEY (`branch_id`) REFERENCES `branches` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_branch_services_service` FOREIGN KEY (`service_id`) REFERENCES `services` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Services offered by each branch';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `branches`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `branches` (
  `id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'UUID primary key',
  `franchise_id` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'FK to franchise (optional)',
  `representative_id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'FK to person with REPRESENTATIVE role',
  `name` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Branch name',
  `establishment_type` varchar(30) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'WORKSHOP, STORE, or WORKSHOP_STORE',
  `profile_image_url` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'Profile image URL',
  `status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'ACTIVE' COMMENT 'ACTIVE or INACTIVE',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `fk_branches_franchise` (`franchise_id`),
  KEY `fk_branches_representative` (`representative_id`),
  CONSTRAINT `fk_branches_franchise` FOREIGN KEY (`franchise_id`) REFERENCES `franchises` (`id`) ON DELETE SET NULL ON UPDATE CASCADE,
  CONSTRAINT `fk_branches_representative` FOREIGN KEY (`representative_id`) REFERENCES `persons` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE,
  CONSTRAINT `chk_branches_status` CHECK ((`status` in (_utf8mb4'ACTIVE',_utf8mb4'INACTIVE'))),
  CONSTRAINT `chk_branches_type` CHECK ((`establishment_type` in (_latin1'WORKSHOP',_latin1'STORE',_latin1'WORKSHOP_STORE')))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Workshop and store branches';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `cities`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `cities` (
  `id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'UUID primary key',
  `name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'City name',
  `department_id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'FK to department',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `fk_cities_department` (`department_id`),
  KEY `idx_cities_name` (`name`),
  CONSTRAINT `fk_cities_department` FOREIGN KEY (`department_id`) REFERENCES `departments` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Catalog of cities by department';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `ciudades`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ciudades` (
  `id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'UUID primary key',
  `nombre` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Nombre de la ciudad',
  `departamento_id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'FK al departamento',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `fk_ciudades_departamento` (`departamento_id`),
  KEY `idx_ciudades_nombre` (`nombre`),
  CONSTRAINT `fk_ciudades_departamento` FOREIGN KEY (`departamento_id`) REFERENCES `departamentos` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='CatÃ¡logo de ciudades por departamento';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `departamentos`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `departamentos` (
  `id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'UUID primary key',
  `nombre` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Nombre del departamento',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_departamentos_nombre` (`nombre`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='CatÃ¡logo de departamentos de Colombia';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `departments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `departments` (
  `id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'UUID primary key',
  `name` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Department name',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_departments_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Catalog of Colombian departments';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `detalles_horario`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `detalles_horario` (
  `id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'UUID primary key',
  `horario_id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'FK a horarios_sede',
  `tipo_entrada` enum('REGULAR','EXCEPCION') COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'REGULAR=por dÃ­a, EXCEPCION=por fecha',
  `dia_semana` tinyint DEFAULT NULL COMMENT '1=Lunes...7=Domingo (NULL para EXCEPCION)',
  `fecha_excepcion` date DEFAULT NULL COMMENT 'Fecha especÃ­fica para excepciones (NULL para REGULAR)',
  `hora_apertura` time DEFAULT NULL COMMENT 'Hora apertura (NULL si es_cerrado=true)',
  `hora_cierre` time DEFAULT NULL COMMENT 'Hora cierre (NULL si es_cerrado=true)',
  `es_cerrado` tinyint(1) DEFAULT '0' COMMENT 'True = cerrado este dÃ­a/fecha',
  `activo` tinyint(1) DEFAULT '1' COMMENT 'Activar/desactivar esta entrada especÃ­fica',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `fecha_inicio_excepcion` date DEFAULT NULL COMMENT 'Fecha inicio para excepciones (NULL para REGULAR)',
  `fecha_fin_excepcion` date DEFAULT NULL COMMENT 'Fecha fin para excepciones (igual a inicio si es un solo dÃ­a)',
  PRIMARY KEY (`id`),
  KEY `idx_detalles_horario_tipo` (`horario_id`,`tipo_entrada`),
  KEY `idx_detalles_horario_dia` (`horario_id`,`dia_semana`),
  KEY `idx_detalles_horario_fecha` (`horario_id`,`fecha_excepcion`),
  KEY `idx_detalles_horario_fechas_excepcion` (`horario_id`,`fecha_inicio_excepcion`,`fecha_fin_excepcion`),
  CONSTRAINT `fk_entradas_horario_horario` FOREIGN KEY (`horario_id`) REFERENCES `horarios_sede` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `chk_entradas_dia` CHECK (((`dia_semana` is null) or ((`dia_semana` >= 1) and (`dia_semana` <= 7))))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Franjas horarias y excepciones (HU6-9, HU20-25)';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `diagnosticos`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `diagnosticos` (
  `id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'UUID primary key',
  `motocicleta_id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'FK a motocicleta',
  `sede_id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'FK a sede que emite el diagnÃ³stico',
  `fecha` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'Fecha del diagnÃ³stico',
  `descripcion_problema` text COLLATE utf8mb4_unicode_ci COMMENT 'DescripciÃ³n del problema por el usuario',
  `posible_solucion` text COLLATE utf8mb4_unicode_ci COMMENT 'Respuesta de la sede',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `fk_diagnosticos_motocicleta` (`motocicleta_id`),
  KEY `fk_diagnosticos_sede` (`sede_id`),
  KEY `idx_diagnosticos_fecha` (`fecha`),
  CONSTRAINT `fk_diagnosticos_motocicleta` FOREIGN KEY (`motocicleta_id`) REFERENCES `motocicletas` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_diagnosticos_sede` FOREIGN KEY (`sede_id`) REFERENCES `sedes` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='DiagnÃ³sticos y cotizaciones para las motocicletas';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `evidencias_diagnostico`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `evidencias_diagnostico` (
  `id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'UUID primary key',
  `diagnostico_id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'FK a solicitud de diagnÃ³stico',
  `url_imagen` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'URL de Firebase Storage',
  `descripcion` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'DescripciÃ³n opcional de la foto',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_evidencias_diagnostico_diagnostico` (`diagnostico_id`),
  CONSTRAINT `fk_evidencias_diagnostico_diagnostico` FOREIGN KEY (`diagnostico_id`) REFERENCES `diagnosticos` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Fotos adjuntas a solicitudes de diagnÃ³stico';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `evidencias_moto`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `evidencias_moto` (
  `id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'UUID primary key',
  `motocicleta_id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'FK a motocicleta',
  `titulo` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'Perfil, Motor, DaÃ±o, etc.',
  `url_imagen` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'URL de la imagen',
  `fecha_carga` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'Fecha de carga',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `fk_evidencias_motocicleta` (`motocicleta_id`),
  CONSTRAINT `fk_evidencias_motocicleta` FOREIGN KEY (`motocicleta_id`) REFERENCES `motocicletas` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Fotos y evidencias de las motocicletas';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `franchises`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `franchises` (
  `id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'UUID primary key',
  `name` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Franchise name',
  `description` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'Franchise description',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_franchises_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Franchises for workshops and stores';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `franquicias`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `franquicias` (
  `id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'UUID primary key',
  `nombre` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Nombre de la franquicia',
  `descripcion` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'DescripciÃ³n de la franquicia',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_franquicias_nombre` (`nombre`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Franquicias de talleres y almacenes';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `horarios_sede`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `horarios_sede` (
  `id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'UUID primary key',
  `sede_id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'FK a sede (Ãºnico por sede)',
  `activo` tinyint(1) DEFAULT '1' COMMENT 'Interruptor maestro para activar/desactivar todo el horario',
  `fecha_inicio` date NOT NULL DEFAULT (curdate()) COMMENT 'Fecha inicio vigencia del horario',
  `fecha_fin` date DEFAULT NULL COMMENT 'Fecha fin vigencia del horario (NULL = indefinido)',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_horarios_sede_sede` (`sede_id`),
  CONSTRAINT `fk_horarios_sede_sede` FOREIGN KEY (`sede_id`) REFERENCES `sedes` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `chk_horarios_sede_fechas` CHECK (((`fecha_fin` is null) or (`fecha_fin` >= `fecha_inicio`)))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='ConfiguraciÃ³n de horario de sede (HU30-35)';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `items_servicio_realizado`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `items_servicio_realizado` (
  `id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'UUID primary key',
  `servicio_realizado_id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'FK a servicios_realizados',
  `servicio_id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'FK al catÃ¡logo de servicios',
  `calificacion` tinyint DEFAULT NULL COMMENT 'CalificaciÃ³n del 1 al 5 (NULL hasta que se califique)',
  `comentario` text COLLATE utf8mb4_unicode_ci COMMENT 'ReseÃ±a opcional del motociclista',
  `fecha_calificacion` datetime DEFAULT NULL COMMENT 'CuÃ¡ndo se enviÃ³ la calificaciÃ³n',
  `es_comentario_ofensivo` tinyint(1) NOT NULL DEFAULT '0' COMMENT 'Indicador: 0=normal, 1=ofensivo',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `fk_items_sr_servicio` (`servicio_id`),
  KEY `idx_items_sr_servicio_realizado` (`servicio_realizado_id`),
  CONSTRAINT `fk_items_sr_servicio` FOREIGN KEY (`servicio_id`) REFERENCES `servicios` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE,
  CONSTRAINT `fk_items_sr_servicio_realizado` FOREIGN KEY (`servicio_realizado_id`) REFERENCES `servicios_realizados` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `chk_items_sr_calificacion` CHECK (((`calificacion` is null) or ((`calificacion` >= 1) and (`calificacion` <= 5))))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Servicios realizados con calificaciÃ³n opcional (HU48)';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `locations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `locations` (
  `id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'UUID primary key',
  `branch_id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'FK to branch (unique per branch)',
  `city_id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'FK to city',
  `address` varchar(150) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Physical address',
  `latitude` decimal(10,8) DEFAULT NULL COMMENT 'GPS latitude',
  `longitude` decimal(11,8) DEFAULT NULL COMMENT 'GPS longitude',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_locations_branch` (`branch_id`),
  KEY `fk_locations_city` (`city_id`),
  CONSTRAINT `fk_locations_branch` FOREIGN KEY (`branch_id`) REFERENCES `branches` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_locations_city` FOREIGN KEY (`city_id`) REFERENCES `cities` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Physical location of each branch';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `motocicletas`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `motocicletas` (
  `id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'UUID primary key',
  `matricula` varchar(10) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Placa de la motocicleta',
  `referencia_id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'FK a referencia_moto',
  `dueno_id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'FK a persona con rol USUARIO',
  `anio` int DEFAULT NULL COMMENT 'AÃ±o de la motocicleta',
  `color` varchar(30) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'Color de la motocicleta',
  `kilometraje_actual` int DEFAULT NULL COMMENT 'Kilometraje actual',
  `observaciones_dueno` text COLLATE utf8mb4_unicode_ci COMMENT 'Observaciones del dueÃ±o',
  `url_imagen_perfil` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'Foto de perfil de la motocicleta (URL de Firebase Storage)',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT 'Timestamp de eliminaciÃ³n lÃ³gica (HU45)',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_motocicletas_matricula` (`matricula`),
  KEY `fk_motocicletas_referencia` (`referencia_id`),
  KEY `fk_motocicletas_dueno` (`dueno_id`),
  KEY `idx_motocicletas_eliminado` (`deleted_at`),
  CONSTRAINT `fk_motocicletas_dueno` FOREIGN KEY (`dueno_id`) REFERENCES `persons` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE,
  CONSTRAINT `fk_motocicletas_referencia` FOREIGN KEY (`referencia_id`) REFERENCES `referencias_moto` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Motocicletas registradas por los usuarios';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `permisos_moto_sede`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `permisos_moto_sede` (
  `id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'UUID primary key',
  `motocicleta_id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'FK a motocicleta',
  `sede_id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'FK a sede autorizada',
  `activo` tinyint(1) DEFAULT '1' COMMENT 'Indica si el permiso estÃ¡ activo',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_permisos_moto_sede` (`motocicleta_id`,`sede_id`),
  KEY `fk_permisos_moto_sede_sede` (`sede_id`),
  CONSTRAINT `fk_permisos_moto_sede_moto` FOREIGN KEY (`motocicleta_id`) REFERENCES `motocicletas` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_permisos_moto_sede_sede` FOREIGN KEY (`sede_id`) REFERENCES `sedes` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Permisos de diagnÃ³stico por moto y sede';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `persons`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `persons` (
  `id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Unique UUID identifier for each person',
  `identity_number` varchar(10) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Identity document number (ID card, passport, DNI, etc.)',
  `first_name` varchar(120) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'First name of the person',
  `last_name` varchar(120) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Last name of the person',
  `second_last_name` varchar(120) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'Second last name of the person',
  `email` varchar(250) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Email address',
  `phone_number` varchar(10) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Cell phone number',
  `role` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Role of the person in the system',
  `keycloak_user_id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Reference to the Keycloak user UUID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_persons_identity_number` (`identity_number`),
  UNIQUE KEY `uk_persons_email` (`email`),
  UNIQUE KEY `uk_persons_phone_number` (`phone_number`),
  UNIQUE KEY `uk_persons_keycloak_user_id` (`keycloak_user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `rangos_cilindraje_sede`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `rangos_cilindraje_sede` (
  `id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'UUID primary key',
  `sede_id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'FK a sede',
  `rango_cilindraje` enum('BAJO','MEDIO','ALTO') COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Rango de cilindraje: BAJO(50-200cc), MEDIO(201-400cc), ALTO(401-3000cc)',
  `activo` tinyint(1) DEFAULT '1' COMMENT 'Indica si el rango estÃ¡ activo en la sede',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_rangos_cilindraje_sede` (`sede_id`,`rango_cilindraje`),
  CONSTRAINT `fk_rangos_cilindraje_sede_sede` FOREIGN KEY (`sede_id`) REFERENCES `sedes` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Rangos de cilindraje manejados por cada sede';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `schedule_details`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `schedule_details` (
  `id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'UUID primary key',
  `schedule_id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'FK to branch_schedules',
  `entry_type` enum('REGULAR','EXCEPTION') COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'REGULAR=by day, EXCEPTION=by date',
  `day_of_week` tinyint DEFAULT NULL COMMENT '1=Lunes...7=Domingo (NULL for EXCEPTION)',
  `exception_date` date DEFAULT NULL COMMENT 'Specific date for exceptions (NULL for REGULAR)',
  `opening_time` time DEFAULT NULL COMMENT 'Opening time (NULL if is_closed=true)',
  `closing_time` time DEFAULT NULL COMMENT 'Closing time (NULL if is_closed=true)',
  `is_closed` tinyint(1) DEFAULT '0' COMMENT 'True = closed on this day/date',
  `active` tinyint(1) DEFAULT '1' COMMENT 'Enable/disable this specific entry',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `exception_start_date` date DEFAULT NULL COMMENT 'Start date for exceptions (NULL for REGULAR)',
  `exception_end_date` date DEFAULT NULL COMMENT 'End date for exceptions (same as start if single day)',
  PRIMARY KEY (`id`),
  KEY `idx_schedule_details_type` (`schedule_id`,`entry_type`),
  KEY `idx_schedule_details_day` (`schedule_id`,`day_of_week`),
  KEY `idx_schedule_details_date` (`schedule_id`,`exception_date`),
  KEY `idx_schedule_details_exception_dates` (`schedule_id`,`exception_start_date`,`exception_end_date`),
  CONSTRAINT `fk_schedule_entries_schedule` FOREIGN KEY (`schedule_id`) REFERENCES `branch_schedules` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `chk_entries_day` CHECK (((`day_of_week` is null) or ((`day_of_week` >= 1) and (`day_of_week` <= 7))))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Schedule time slots and exceptions (HU6-9, HU20-25)';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `sedes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sedes` (
  `id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'UUID primary key',
  `franquicia_id` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'FK a franquicia (opcional)',
  `representante_id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'FK a persona con rol REPRESENTANTE',
  `nombre` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Nombre de la sede',
  `tipo_establecimiento` varchar(30) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'TALLER, ALMACEN o TALLER_ALMACEN',
  `url_imagen_perfil` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'URL de imagen de perfil',
  `estado` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'ACTIVO' COMMENT 'ACTIVO o INACTIVO',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `fk_sedes_franquicia` (`franquicia_id`),
  KEY `fk_sedes_representante` (`representante_id`),
  CONSTRAINT `fk_sedes_franquicia` FOREIGN KEY (`franquicia_id`) REFERENCES `franquicias` (`id`) ON DELETE SET NULL ON UPDATE CASCADE,
  CONSTRAINT `fk_sedes_representante` FOREIGN KEY (`representante_id`) REFERENCES `persons` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE,
  CONSTRAINT `chk_sedes_estado` CHECK ((`estado` in (_latin1'ACTIVO',_latin1'INACTIVO'))),
  CONSTRAINT `chk_sedes_tipo` CHECK ((`tipo_establecimiento` in (_latin1'TALLER',_latin1'ALMACEN',_latin1'TALLER_ALMACEN')))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Sedes de talleres y almacenes';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `services`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `services` (
  `id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'UUID primary key',
  `name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Service name',
  `description` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'Service description',
  `service_type` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'Maintenance, Repair, etc.',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `is_active` tinyint(1) DEFAULT '1' COMMENT 'Indicates if service is active in catalog',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_services_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Catalog of available services';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `servicios`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `servicios` (
  `id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'UUID primary key',
  `nombre` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Nombre del servicio',
  `descripcion` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'DescripciÃ³n del servicio',
  `tipo_servicio` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'Mantenimiento, Reparacion, etc.',
  `is_active` tinyint(1) DEFAULT '1' COMMENT 'Indica si el servicio estÃ¡ activo en el catÃ¡logo',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_servicios_nombre` (`nombre`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='CatÃ¡logo de servicios disponibles';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `servicios_realizados`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `servicios_realizados` (
  `id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'UUID primary key',
  `sede_id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'FK a sede que realiza el servicio',
  `motocicleta_id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'FK a motocicleta',
  `diagnostico_id` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'FK a diagnÃ³stico que originÃ³ el servicio (opcional)',
  `fecha_solicitud` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'Fecha de solicitud del servicio',
  `fecha_finalizacion` datetime DEFAULT NULL COMMENT 'Fecha de finalizaciÃ³n del servicio',
  `estado` enum('PENDIENTE','EN_PROCESO','FINALIZADO','CANCELADO') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'PENDIENTE' COMMENT 'Estado del servicio',
  `motivo_ingreso` text COLLATE utf8mb4_unicode_ci COMMENT 'Motivo de ingreso de la moto',
  `precio_cotizado` decimal(12,2) DEFAULT NULL COMMENT 'Precio cotizado antes del servicio',
  `precio_final` decimal(12,2) DEFAULT NULL COMMENT 'Precio final cobrado al terminar',
  `notas_representante` text COLLATE utf8mb4_unicode_ci COMMENT 'Notas del representante para el cliente',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `fk_servicios_realizados_sede` (`sede_id`),
  KEY `fk_servicios_realizados_motocicleta` (`motocicleta_id`),
  KEY `fk_servicios_realizados_diagnostico` (`diagnostico_id`),
  KEY `idx_servicios_realizados_fecha` (`fecha_solicitud`),
  CONSTRAINT `fk_servicios_realizados_diagnostico` FOREIGN KEY (`diagnostico_id`) REFERENCES `diagnosticos` (`id`) ON DELETE SET NULL ON UPDATE CASCADE,
  CONSTRAINT `fk_servicios_realizados_motocicleta` FOREIGN KEY (`motocicleta_id`) REFERENCES `motocicletas` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_servicios_realizados_sede` FOREIGN KEY (`sede_id`) REFERENCES `sedes` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Servicios realizados a las motocicletas';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `servicios_sede`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `servicios_sede` (
  `id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'UUID primary key',
  `sede_id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'FK a sede',
  `servicio_id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'FK a servicio',
  `activo` tinyint(1) DEFAULT '1' COMMENT 'Indica si el servicio estÃ¡ activo en la sede',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_servicios_sede` (`sede_id`,`servicio_id`),
  KEY `fk_servicios_sede_servicio` (`servicio_id`),
  CONSTRAINT `fk_servicios_sede_sede` FOREIGN KEY (`sede_id`) REFERENCES `sedes` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_servicios_sede_servicio` FOREIGN KEY (`servicio_id`) REFERENCES `servicios` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='RelaciÃ³n de servicios ofrecidos por cada sede';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `system_messages`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `system_messages` (
  `id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'UUID primary key',
  `message_code` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Unique message code (e.g.: MOD_U_DUP_ERR_00001)',
  `type` enum('ERROR','EXITO','WARNING','INFO','DEBUG') COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Message type',
  `category` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Category: end_user, system, validation',
  `module` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Module it belongs to: general, users, validation, authorization',
  `message_title` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Short message title (in Spanish)',
  `message_content` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Message content with placeholders NULL, NULL, etc. (in Spanish)',
  `is_active` tinyint(1) DEFAULT '1' COMMENT 'Indicates if the message is active',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `message_code` (`message_code`),
  KEY `idx_code` (`message_code`),
  KEY `idx_type` (`type`),
  KEY `idx_module` (`module`),
  KEY `idx_active` (`is_active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Table for dynamic system messages';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `transiciones_estado_servicio`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `transiciones_estado_servicio` (
  `id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'UUID primary key',
  `servicio_realizado_id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'FK a servicios_realizados',
  `estado_anterior` enum('PENDIENTE','EN_PROCESO','FINALIZADO','CANCELADO') COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'Estado anterior (NULL si es creaciÃ³n)',
  `estado_nuevo` enum('PENDIENTE','EN_PROCESO','FINALIZADO','CANCELADO') COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Nuevo estado del servicio',
  `creado_por` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'FK a persona (representante) que realizÃ³ el cambio',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `fk_transiciones_estado_persona` (`creado_por`),
  KEY `idx_transiciones_estado_servicio` (`servicio_realizado_id`),
  KEY `idx_transiciones_estado_fecha` (`created_at`),
  CONSTRAINT `fk_transiciones_estado_persona` FOREIGN KEY (`creado_por`) REFERENCES `persons` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE,
  CONSTRAINT `fk_transiciones_estado_servicio` FOREIGN KEY (`servicio_realizado_id`) REFERENCES `servicios_realizados` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Transiciones de estado de servicios realizados';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `ubicaciones`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ubicaciones` (
  `id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'UUID primary key',
  `sede_id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'FK a sede (Ãºnica por sede)',
  `ciudad_id` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'FK a ciudad',
  `direccion` varchar(150) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'DirecciÃ³n fÃ­sica',
  `latitud` decimal(10,8) DEFAULT NULL COMMENT 'Latitud GPS',
  `longitud` decimal(11,8) DEFAULT NULL COMMENT 'Longitud GPS',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_ubicaciones_sede` (`sede_id`),
  KEY `fk_ubicaciones_ciudad` (`ciudad_id`),
  CONSTRAINT `fk_ubicaciones_ciudad` FOREIGN KEY (`ciudad_id`) REFERENCES `ciudades` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE,
  CONSTRAINT `fk_ubicaciones_sede` FOREIGN KEY (`sede_id`) REFERENCES `sedes` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='UbicaciÃ³n fÃ­sica de cada sede';
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

