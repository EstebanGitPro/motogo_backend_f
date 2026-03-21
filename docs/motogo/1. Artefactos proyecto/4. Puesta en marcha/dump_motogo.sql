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

LOCK TABLES `branch_displacement_ranges` WRITE;
/*!40000 ALTER TABLE `branch_displacement_ranges` DISABLE KEYS */;
/*!40000 ALTER TABLE `branch_displacement_ranges` ENABLE KEYS */;
UNLOCK TABLES;
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

LOCK TABLES `branch_schedules` WRITE;
/*!40000 ALTER TABLE `branch_schedules` DISABLE KEYS */;
/*!40000 ALTER TABLE `branch_schedules` ENABLE KEYS */;
UNLOCK TABLES;
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

LOCK TABLES `branch_services` WRITE;
/*!40000 ALTER TABLE `branch_services` DISABLE KEYS */;
/*!40000 ALTER TABLE `branch_services` ENABLE KEYS */;
UNLOCK TABLES;
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

LOCK TABLES `branches` WRITE;
/*!40000 ALTER TABLE `branches` DISABLE KEYS */;
/*!40000 ALTER TABLE `branches` ENABLE KEYS */;
UNLOCK TABLES;
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

LOCK TABLES `cities` WRITE;
/*!40000 ALTER TABLE `cities` DISABLE KEYS */;
INSERT INTO `cities` VALUES ('b2c3d4e5-2222-4000-8000-000000000001','BogotÃ¡','a1b2c3d4-1111-4000-8000-000000000001','2026-03-13 01:28:22','2026-03-13 01:28:22'),('b2c3d4e5-2222-4000-8000-000000000002','Soacha','a1b2c3d4-1111-4000-8000-000000000001','2026-03-13 01:28:22','2026-03-13 01:28:22'),('b2c3d4e5-2222-4000-8000-000000000003','ZipaquirÃ¡','a1b2c3d4-1111-4000-8000-000000000001','2026-03-13 01:28:22','2026-03-13 01:28:22'),('b2c3d4e5-2222-4000-8000-000000000004','MedellÃ­n','a1b2c3d4-1111-4000-8000-000000000002','2026-03-13 01:28:22','2026-03-13 01:28:22'),('b2c3d4e5-2222-4000-8000-000000000005','Envigado','a1b2c3d4-1111-4000-8000-000000000002','2026-03-13 01:28:22','2026-03-13 01:28:22'),('b2c3d4e5-2222-4000-8000-000000000006','Bello','a1b2c3d4-1111-4000-8000-000000000002','2026-03-13 01:28:22','2026-03-13 01:28:22'),('b2c3d4e5-2222-4000-8000-000000000007','Cali','a1b2c3d4-1111-4000-8000-000000000003','2026-03-13 01:28:22','2026-03-13 01:28:22'),('b2c3d4e5-2222-4000-8000-000000000008','Palmira','a1b2c3d4-1111-4000-8000-000000000003','2026-03-13 01:28:22','2026-03-13 01:28:22'),('b2c3d4e5-2222-4000-8000-000000000009','Barranquilla','a1b2c3d4-1111-4000-8000-000000000004','2026-03-13 01:28:22','2026-03-13 01:28:22'),('b2c3d4e5-2222-4000-8000-000000000010','Soledad','a1b2c3d4-1111-4000-8000-000000000004','2026-03-13 01:28:22','2026-03-13 01:28:22'),('b2c3d4e5-2222-4000-8000-000000000011','Bucaramanga','a1b2c3d4-1111-4000-8000-000000000005','2026-03-13 01:28:22','2026-03-13 01:28:22'),('b2c3d4e5-2222-4000-8000-000000000012','Floridablanca','a1b2c3d4-1111-4000-8000-000000000005','2026-03-13 01:28:22','2026-03-13 01:28:22');
/*!40000 ALTER TABLE `cities` ENABLE KEYS */;
UNLOCK TABLES;
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

LOCK TABLES `ciudades` WRITE;
/*!40000 ALTER TABLE `ciudades` DISABLE KEYS */;
/*!40000 ALTER TABLE `ciudades` ENABLE KEYS */;
UNLOCK TABLES;
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

LOCK TABLES `departamentos` WRITE;
/*!40000 ALTER TABLE `departamentos` DISABLE KEYS */;
/*!40000 ALTER TABLE `departamentos` ENABLE KEYS */;
UNLOCK TABLES;
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

LOCK TABLES `departments` WRITE;
/*!40000 ALTER TABLE `departments` DISABLE KEYS */;
INSERT INTO `departments` VALUES ('a1b2c3d4-1111-4000-8000-000000000001','Cundinamarca','2026-03-13 01:28:22','2026-03-13 01:28:22'),('a1b2c3d4-1111-4000-8000-000000000002','Antioquia','2026-03-13 01:28:22','2026-03-13 01:28:22'),('a1b2c3d4-1111-4000-8000-000000000003','Valle del Cauca','2026-03-13 01:28:22','2026-03-13 01:28:22'),('a1b2c3d4-1111-4000-8000-000000000004','AtlÃ¡ntico','2026-03-13 01:28:22','2026-03-13 01:28:22'),('a1b2c3d4-1111-4000-8000-000000000005','Santander','2026-03-13 01:28:22','2026-03-13 01:28:22');
/*!40000 ALTER TABLE `departments` ENABLE KEYS */;
UNLOCK TABLES;
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

LOCK TABLES `detalles_horario` WRITE;
/*!40000 ALTER TABLE `detalles_horario` DISABLE KEYS */;
/*!40000 ALTER TABLE `detalles_horario` ENABLE KEYS */;
UNLOCK TABLES;
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

LOCK TABLES `diagnosticos` WRITE;
/*!40000 ALTER TABLE `diagnosticos` DISABLE KEYS */;
/*!40000 ALTER TABLE `diagnosticos` ENABLE KEYS */;
UNLOCK TABLES;
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

LOCK TABLES `evidencias_diagnostico` WRITE;
/*!40000 ALTER TABLE `evidencias_diagnostico` DISABLE KEYS */;
/*!40000 ALTER TABLE `evidencias_diagnostico` ENABLE KEYS */;
UNLOCK TABLES;
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

LOCK TABLES `evidencias_moto` WRITE;
/*!40000 ALTER TABLE `evidencias_moto` DISABLE KEYS */;
/*!40000 ALTER TABLE `evidencias_moto` ENABLE KEYS */;
UNLOCK TABLES;
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

LOCK TABLES `franchises` WRITE;
/*!40000 ALTER TABLE `franchises` DISABLE KEYS */;
INSERT INTO `franchises` VALUES ('c3d4e5f6-3333-4000-8000-000000000001','MotoMax','Cadena nacional de talleres y repuestos','2026-03-13 01:28:22','2026-03-13 01:28:22'),('c3d4e5f6-3333-4000-8000-000000000002','Repuestos La Moto','Red de almacenes especializados en repuestos','2026-03-13 01:28:22','2026-03-13 01:28:22'),('c3d4e5f6-3333-4000-8000-000000000003','TallerPro','Talleres profesionales de motocicletas','2026-03-13 01:28:22','2026-03-13 01:28:22');
/*!40000 ALTER TABLE `franchises` ENABLE KEYS */;
UNLOCK TABLES;
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

LOCK TABLES `franquicias` WRITE;
/*!40000 ALTER TABLE `franquicias` DISABLE KEYS */;
/*!40000 ALTER TABLE `franquicias` ENABLE KEYS */;
UNLOCK TABLES;
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

LOCK TABLES `horarios_sede` WRITE;
/*!40000 ALTER TABLE `horarios_sede` DISABLE KEYS */;
/*!40000 ALTER TABLE `horarios_sede` ENABLE KEYS */;
UNLOCK TABLES;
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

LOCK TABLES `items_servicio_realizado` WRITE;
/*!40000 ALTER TABLE `items_servicio_realizado` DISABLE KEYS */;
/*!40000 ALTER TABLE `items_servicio_realizado` ENABLE KEYS */;
UNLOCK TABLES;
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

LOCK TABLES `locations` WRITE;
/*!40000 ALTER TABLE `locations` DISABLE KEYS */;
/*!40000 ALTER TABLE `locations` ENABLE KEYS */;
UNLOCK TABLES;
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

LOCK TABLES `motocicletas` WRITE;
/*!40000 ALTER TABLE `motocicletas` DISABLE KEYS */;
/*!40000 ALTER TABLE `motocicletas` ENABLE KEYS */;
UNLOCK TABLES;
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

LOCK TABLES `permisos_moto_sede` WRITE;
/*!40000 ALTER TABLE `permisos_moto_sede` DISABLE KEYS */;
/*!40000 ALTER TABLE `permisos_moto_sede` ENABLE KEYS */;
UNLOCK TABLES;
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

LOCK TABLES `persons` WRITE;
/*!40000 ALTER TABLE `persons` DISABLE KEYS */;
/*!40000 ALTER TABLE `persons` ENABLE KEYS */;
UNLOCK TABLES;
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

LOCK TABLES `rangos_cilindraje_sede` WRITE;
/*!40000 ALTER TABLE `rangos_cilindraje_sede` DISABLE KEYS */;
/*!40000 ALTER TABLE `rangos_cilindraje_sede` ENABLE KEYS */;
UNLOCK TABLES;
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

LOCK TABLES `schedule_details` WRITE;
/*!40000 ALTER TABLE `schedule_details` DISABLE KEYS */;
/*!40000 ALTER TABLE `schedule_details` ENABLE KEYS */;
UNLOCK TABLES;
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

LOCK TABLES `sedes` WRITE;
/*!40000 ALTER TABLE `sedes` DISABLE KEYS */;
/*!40000 ALTER TABLE `sedes` ENABLE KEYS */;
UNLOCK TABLES;
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

LOCK TABLES `services` WRITE;
/*!40000 ALTER TABLE `services` DISABLE KEYS */;
INSERT INTO `services` VALUES ('ecf34328-1e7b-11f1-9193-72232b494de6','Cambio de aceite','Cambio de aceite de motor','Mantenimiento','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf346ac-1e7b-11f1-9193-72232b494de6','Cambio de aceite y filtro','Cambio de aceite de motor con reemplazo de filtro','Mantenimiento','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf34802-1e7b-11f1-9193-72232b494de6','RevisiÃ³n general preventiva','InspecciÃ³n completa del estado general de la motocicleta','Mantenimiento','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf34886-1e7b-11f1-9193-72232b494de6','Mantenimiento preventivo completo','Servicio integral de mantenimiento preventivo','Mantenimiento','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf348f4-1e7b-11f1-9193-72232b494de6','Mantenimiento por kilometraje','Servicio segÃºn kilometraje recorrido','Mantenimiento','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf3495c-1e7b-11f1-9193-72232b494de6','Ajuste general de tornillerÃ­a','RevisiÃ³n y ajuste de toda la tornillerÃ­a','Mantenimiento','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf34a08-1e7b-11f1-9193-72232b494de6','Limpieza y lubricaciÃ³n de partes mÃ³viles','LubricaciÃ³n de cadena, cables y partes mÃ³viles','Mantenimiento','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf34a73-1e7b-11f1-9193-72232b494de6','RevisiÃ³n de niveles','VerificaciÃ³n de niveles de aceite, frenos y refrigerante','Mantenimiento','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf34ae1-1e7b-11f1-9193-72232b494de6','SincronizaciÃ³n de motor','Ajuste de sincronizaciÃ³n del motor','Mantenimiento','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf34b79-1e7b-11f1-9193-72232b494de6','SincronizaciÃ³n de carburador','Ajuste y sincronizaciÃ³n del carburador','Mantenimiento','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf34c06-1e7b-11f1-9193-72232b494de6','Limpieza de carburador','Limpieza profunda del carburador','Mantenimiento','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf34c96-1e7b-11f1-9193-72232b494de6','Ajuste de ralentÃ­','CalibraciÃ³n del rÃ©gimen de ralentÃ­','Mantenimiento','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf34e68-1e7b-11f1-9193-72232b494de6','Ajuste de vÃ¡lvulas','CalibraciÃ³n de vÃ¡lvulas del motor','Mantenimiento','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf34ef1-1e7b-11f1-9193-72232b494de6','RevisiÃ³n de sistema de refrigeraciÃ³n','InspecciÃ³n del sistema de refrigeraciÃ³n','Mantenimiento','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf34f5d-1e7b-11f1-9193-72232b494de6','Cambio de refrigerante','Reemplazo del lÃ­quido refrigerante','Mantenimiento','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf34fc2-1e7b-11f1-9193-72232b494de6','Limpieza de radiador','Limpieza y mantenimiento del radiador','Mantenimiento','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf35033-1e7b-11f1-9193-72232b494de6','RevisiÃ³n de bujÃ­as','InspecciÃ³n del estado de las bujÃ­as','Mantenimiento','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf3511e-1e7b-11f1-9193-72232b494de6','Cambio de bujÃ­as','Reemplazo de bujÃ­as','Mantenimiento','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf35195-1e7b-11f1-9193-72232b494de6','RevisiÃ³n de baterÃ­a','InspecciÃ³n del estado y carga de la baterÃ­a','Mantenimiento','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf35201-1e7b-11f1-9193-72232b494de6','Carga de baterÃ­a','Recarga de baterÃ­a con equipo especializado','Mantenimiento','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf35268-1e7b-11f1-9193-72232b494de6','Cambio de baterÃ­a','Reemplazo de baterÃ­a','Mantenimiento','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf352cd-1e7b-11f1-9193-72232b494de6','RevisiÃ³n de sistema elÃ©ctrico','InspecciÃ³n general del sistema elÃ©ctrico','Mantenimiento','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf3bf96-1e7b-11f1-9193-72232b494de6','ReparaciÃ³n de motor','ReparaciÃ³n de componentes internos del motor','ReparaciÃ³n','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf3c208-1e7b-11f1-9193-72232b494de6','ReparaciÃ³n de culata','ReparaciÃ³n y rectificaciÃ³n de culata','ReparaciÃ³n','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf3c323-1e7b-11f1-9193-72232b494de6','RectificaciÃ³n de motor','RectificaciÃ³n completa del motor','ReparaciÃ³n','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf3c39e-1e7b-11f1-9193-72232b494de6','Cambio de pistÃ³n','Reemplazo de pistÃ³n del motor','ReparaciÃ³n','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf3c40a-1e7b-11f1-9193-72232b494de6','Cambio de anillos','Reemplazo de anillos del pistÃ³n','ReparaciÃ³n','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf3c472-1e7b-11f1-9193-72232b494de6','ReparaciÃ³n de caja de cambios','ReparaciÃ³n de transmisiÃ³n','ReparaciÃ³n','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf3c4dd-1e7b-11f1-9193-72232b494de6','ReparaciÃ³n de embrague','ReparaciÃ³n del sistema de embrague','ReparaciÃ³n','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf3c547-1e7b-11f1-9193-72232b494de6','Cambio de discos de embrague','Reemplazo de discos de embrague','ReparaciÃ³n','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf3c5b0-1e7b-11f1-9193-72232b494de6','ReparaciÃ³n de sistema de frenos','ReparaciÃ³n integral del sistema de frenos','ReparaciÃ³n','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf3c621-1e7b-11f1-9193-72232b494de6','Cambio de pastillas de freno','Reemplazo de pastillas de freno','ReparaciÃ³n','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf3c68e-1e7b-11f1-9193-72232b494de6','Cambio de bandas de freno','Reemplazo de bandas de freno trasero','ReparaciÃ³n','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf3c6f8-1e7b-11f1-9193-72232b494de6','ReparaciÃ³n de freno delantero','ReparaciÃ³n del sistema de freno delantero','ReparaciÃ³n','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf3c787-1e7b-11f1-9193-72232b494de6','ReparaciÃ³n de freno trasero','ReparaciÃ³n del sistema de freno trasero','ReparaciÃ³n','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf3c818-1e7b-11f1-9193-72232b494de6','Cambio de cadena','Reemplazo de cadena de transmisiÃ³n','ReparaciÃ³n','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf3c895-1e7b-11f1-9193-72232b494de6','Cambio de piÃ±Ã³n','Reemplazo del piÃ±Ã³n de transmisiÃ³n','ReparaciÃ³n','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf3c905-1e7b-11f1-9193-72232b494de6','Cambio de catalina','Reemplazo de la catalina','ReparaciÃ³n','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf3c96d-1e7b-11f1-9193-72232b494de6','Cambio de kit de arrastre','Reemplazo completo de cadena, piÃ±Ã³n y catalina','ReparaciÃ³n','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf3cc02-1e7b-11f1-9193-72232b494de6','ReparaciÃ³n de suspensiÃ³n delantera','ReparaciÃ³n de barras y amortiguadores delanteros','ReparaciÃ³n','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf3cd6f-1e7b-11f1-9193-72232b494de6','ReparaciÃ³n de suspensiÃ³n trasera','ReparaciÃ³n de amortiguadores traseros','ReparaciÃ³n','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf3cdec-1e7b-11f1-9193-72232b494de6','Cambio de amortiguadores','Reemplazo de amortiguadores','ReparaciÃ³n','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf3cea6-1e7b-11f1-9193-72232b494de6','ReparaciÃ³n de chasis','ReparaciÃ³n y enderezado de chasis','ReparaciÃ³n','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf3cf16-1e7b-11f1-9193-72232b494de6','Soldadura de piezas','Soldadura de partes metÃ¡licas','ReparaciÃ³n','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf3cf83-1e7b-11f1-9193-72232b494de6','ReparaciÃ³n de sistema de escape','ReparaciÃ³n del sistema de escape','ReparaciÃ³n','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf415eb-1e7b-11f1-9193-72232b494de6','Cambio de llantas','Reemplazo de neumÃ¡ticos','Llantas','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf417ca-1e7b-11f1-9193-72232b494de6','Venta de llantas','Venta e instalaciÃ³n de llantas nuevas','Llantas','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf41845-1e7b-11f1-9193-72232b494de6','Montaje de llantas','Montaje de llantas en rines','Llantas','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf418ae-1e7b-11f1-9193-72232b494de6','Desmontaje de llantas','Desmontaje de llantas de rines','Llantas','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf4191c-1e7b-11f1-9193-72232b494de6','ReparaciÃ³n de pinchazos','ReparaciÃ³n de pinchazos y fugas','Llantas','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf41989-1e7b-11f1-9193-72232b494de6','CalibraciÃ³n de llantas','CalibraciÃ³n de presiÃ³n de aire','Llantas','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf419f5-1e7b-11f1-9193-72232b494de6','Balanceo de ruedas','Balanceo de ruedas con pesos','Llantas','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf41aa5-1e7b-11f1-9193-72232b494de6','AlineaciÃ³n de ruedas','AlineaciÃ³n de direcciÃ³n y ruedas','Llantas','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf41b0a-1e7b-11f1-9193-72232b494de6','Enderezado de rines','Enderezado de rines deformados','Llantas','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf41b6d-1e7b-11f1-9193-72232b494de6','Cambio de rines','Reemplazo de rines','Llantas','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf44d42-1e7b-11f1-9193-72232b494de6','DiagnÃ³stico electrÃ³nico','AnÃ¡lisis computarizado de sensores y ECU','DiagnÃ³stico','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf44f0c-1e7b-11f1-9193-72232b494de6','Escaneo computarizado','Lectura de cÃ³digos de error con escÃ¡ner','DiagnÃ³stico','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf44f83-1e7b-11f1-9193-72232b494de6','DiagnÃ³stico de motor','DiagnÃ³stico completo del motor','DiagnÃ³stico','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf44fe7-1e7b-11f1-9193-72232b494de6','DiagnÃ³stico elÃ©ctrico','DiagnÃ³stico del sistema elÃ©ctrico','DiagnÃ³stico','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf45050-1e7b-11f1-9193-72232b494de6','DiagnÃ³stico de inyecciÃ³n','DiagnÃ³stico del sistema de inyecciÃ³n','DiagnÃ³stico','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf450b6-1e7b-11f1-9193-72232b494de6','DiagnÃ³stico de frenos','DiagnÃ³stico del sistema de frenos','DiagnÃ³stico','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf4511b-1e7b-11f1-9193-72232b494de6','DiagnÃ³stico de suspensiÃ³n','DiagnÃ³stico del sistema de suspensiÃ³n','DiagnÃ³stico','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf45185-1e7b-11f1-9193-72232b494de6','DiagnÃ³stico precompra','InspecciÃ³n completa antes de compra de moto usada','DiagnÃ³stico','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf451ed-1e7b-11f1-9193-72232b494de6','RevisiÃ³n tÃ©cnica general','RevisiÃ³n tÃ©cnica integral','DiagnÃ³stico','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf47a60-1e7b-11f1-9193-72232b494de6','Lavado bÃ¡sico','Lavado exterior bÃ¡sico','EstÃ©tica','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf47d00-1e7b-11f1-9193-72232b494de6','Lavado completo','Lavado completo interior y exterior','EstÃ©tica','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf47db8-1e7b-11f1-9193-72232b494de6','Lavado y encerado','Lavado con aplicaciÃ³n de cera protectora','EstÃ©tica','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf47e8a-1e7b-11f1-9193-72232b494de6','Detailing de motocicleta','Detallado profesional completo','EstÃ©tica','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf47f01-1e7b-11f1-9193-72232b494de6','Limpieza profunda de motor','Desengrase y limpieza profunda del motor','EstÃ©tica','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf47f70-1e7b-11f1-9193-72232b494de6','RestauraciÃ³n estÃ©tica','RestauraciÃ³n de partes plÃ¡sticas y pintura','EstÃ©tica','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf47fdf-1e7b-11f1-9193-72232b494de6','Pulido de pintura','Pulido y correcciÃ³n de pintura','EstÃ©tica','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf4804b-1e7b-11f1-9193-72232b494de6','ProtecciÃ³n cerÃ¡mica','AplicaciÃ³n de coating cerÃ¡mico','EstÃ©tica','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf480b4-1e7b-11f1-9193-72232b494de6','Lavado de chasis','Limpieza especializada del chasis','EstÃ©tica','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf4811c-1e7b-11f1-9193-72232b494de6','Limpieza de cadena','Limpieza y lubricaciÃ³n de cadena','EstÃ©tica','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf4b5f6-1e7b-11f1-9193-72232b494de6','InstalaciÃ³n de accesorios','InstalaciÃ³n de accesorios varios','Accesorios','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf4b818-1e7b-11f1-9193-72232b494de6','InstalaciÃ³n de luces auxiliares','InstalaciÃ³n de luces LED auxiliares','Accesorios','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf4b8a1-1e7b-11f1-9193-72232b494de6','InstalaciÃ³n de exploradoras','InstalaciÃ³n de luces exploradoras','Accesorios','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf4b910-1e7b-11f1-9193-72232b494de6','InstalaciÃ³n de sliders','InstalaciÃ³n de sliders protectores','Accesorios','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf4b979-1e7b-11f1-9193-72232b494de6','InstalaciÃ³n de defensas','InstalaciÃ³n de defensas y crashbars','Accesorios','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf4b9e3-1e7b-11f1-9193-72232b494de6','InstalaciÃ³n de maleteros','InstalaciÃ³n de maleteros y baÃºles','Accesorios','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf4ba4e-1e7b-11f1-9193-72232b494de6','InstalaciÃ³n de parrillas','InstalaciÃ³n de parrillas traseras','Accesorios','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf4bab9-1e7b-11f1-9193-72232b494de6','InstalaciÃ³n de alarma','InstalaciÃ³n de sistema de alarma','Accesorios','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf4bb21-1e7b-11f1-9193-72232b494de6','InstalaciÃ³n de GPS','InstalaciÃ³n de rastreador GPS','Accesorios','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf4bb8c-1e7b-11f1-9193-72232b494de6','InstalaciÃ³n de cÃ¡maras','InstalaciÃ³n de cÃ¡mara DVR o dashcam','Accesorios','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf4bbf6-1e7b-11f1-9193-72232b494de6','PersonalizaciÃ³n estÃ©tica','PersonalizaciÃ³n de partes y diseÃ±o','Accesorios','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf4bc60-1e7b-11f1-9193-72232b494de6','ModificaciÃ³n de escape','InstalaciÃ³n de escape deportivo o modificado','Accesorios','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf4f634-1e7b-11f1-9193-72232b494de6','ReparaciÃ³n de sistema elÃ©ctrico','ReparaciÃ³n integral del sistema elÃ©ctrico','ElÃ©ctrico','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf4f875-1e7b-11f1-9193-72232b494de6','Cambio de regulador de voltaje','Reemplazo del regulador rectificador','ElÃ©ctrico','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf4f907-1e7b-11f1-9193-72232b494de6','ReparaciÃ³n de alternador','ReparaciÃ³n del alternador o estator','ElÃ©ctrico','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf4f977-1e7b-11f1-9193-72232b494de6','ReparaciÃ³n de arranque','ReparaciÃ³n del motor de arranque','ElÃ©ctrico','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf4fa1f-1e7b-11f1-9193-72232b494de6','Cambio de luces','Reemplazo de bombillos o luces LED','ElÃ©ctrico','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf4fa8e-1e7b-11f1-9193-72232b494de6','ReparaciÃ³n de tablero','ReparaciÃ³n del tablero de instrumentos','ElÃ©ctrico','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf4fb98-1e7b-11f1-9193-72232b494de6','ReparaciÃ³n de cableado','ReparaciÃ³n de arnÃ©s elÃ©ctrico','ElÃ©ctrico','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf4fc14-1e7b-11f1-9193-72232b494de6','InstalaciÃ³n de toma USB','InstalaciÃ³n de cargador USB en la moto','ElÃ©ctrico','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf52cc1-1e7b-11f1-9193-72232b494de6','RevisiÃ³n pretecnomecÃ¡nica','RevisiÃ³n previa a la tecnomecÃ¡nica oficial','Legal','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf52ee2-1e7b-11f1-9193-72232b494de6','Alistamiento para tecnomecÃ¡nica','PreparaciÃ³n completa para pasar tecnomecÃ¡nica','Legal','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf52f75-1e7b-11f1-9193-72232b494de6','InspecciÃ³n para venta','InspecciÃ³n certificada para venta de moto usada','Legal','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf52fe8-1e7b-11f1-9193-72232b494de6','AvalÃºo de motocicleta','AvalÃºo comercial de la motocicleta','Legal','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf53055-1e7b-11f1-9193-72232b494de6','CertificaciÃ³n de estado mecÃ¡nico','Certificado del estado mecÃ¡nico actual','Legal','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf530bd-1e7b-11f1-9193-72232b494de6','Asistencia en carretera','Servicio de asistencia en carretera','Legal','2026-03-13 01:28:23','2026-03-13 01:28:23',1),('ecf53123-1e7b-11f1-9193-72232b494de6','Servicio a domicilio','Servicio mecÃ¡nico a domicilio','Legal','2026-03-13 01:28:23','2026-03-13 01:28:23',1);
/*!40000 ALTER TABLE `services` ENABLE KEYS */;
UNLOCK TABLES;
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

LOCK TABLES `servicios` WRITE;
/*!40000 ALTER TABLE `servicios` DISABLE KEYS */;
/*!40000 ALTER TABLE `servicios` ENABLE KEYS */;
UNLOCK TABLES;
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

LOCK TABLES `servicios_realizados` WRITE;
/*!40000 ALTER TABLE `servicios_realizados` DISABLE KEYS */;
/*!40000 ALTER TABLE `servicios_realizados` ENABLE KEYS */;
UNLOCK TABLES;
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

LOCK TABLES `servicios_sede` WRITE;
/*!40000 ALTER TABLE `servicios_sede` DISABLE KEYS */;
/*!40000 ALTER TABLE `servicios_sede` ENABLE KEYS */;
UNLOCK TABLES;
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

LOCK TABLES `system_messages` WRITE;
/*!40000 ALTER TABLE `system_messages` DISABLE KEYS */;
INSERT INTO `system_messages` VALUES ('67d2269a-df71-11f0-a78d-1286885819d0','MOD_AUTH_LOGIN_SUCCESS_EXI_00001','EXITO','usuario_final','authentication','Inicio de SesiÃ³n Exitoso','Â¡Bienvenido/a! Has iniciado sesiÃ³n correctamente.',1,'2025-12-22 20:04:21','2025-12-22 20:04:21'),('6ddc165f-df76-11f0-a78d-1286885819d0','MOD_U_EMAIL_NV_ERR_00006','ERROR','usuario_final','authentication','Email No Verificado','Tu correo electrÃ³nico aÃºn no ha sido verificado. Hemos enviado un nuevo correo de verificaciÃ³n a tu bandeja de entrada. Por favor, revisa tu correo (incluyendo la carpeta de spam) y haz clic en el enlace de verificaciÃ³n.',1,'2025-12-22 20:40:19','2025-12-22 20:40:19'),('ebf0d3a8-1e7b-11f1-9193-72232b494de6','GEN_LOG_INIT_EXI_00001','EXITO','sistema','logging','Sistema de logs iniciado','El sistema de logging estructurado ha sido inicializado correctamente con JSON.',1,'2026-03-13 01:28:21','2026-03-13 01:28:21'),('ebf0e4bc-1e7b-11f1-9193-72232b494de6','GEN_LOG_LOKI_CONN_EXI_00002','EXITO','sistema','logging','ConexiÃ³n a Loki exitosa','La conexiÃ³n con Loki para agregaciÃ³n de logs se estableciÃ³ correctamente.',1,'2026-03-13 01:28:21','2026-03-13 01:28:21'),('ebf0e713-1e7b-11f1-9193-72232b494de6','GEN_LOG_LOKI_CONN_ERR_00001','ERROR','sistema','logging','Error de conexiÃ³n a Loki','No se pudo establecer conexiÃ³n con Loki: ${0}',1,'2026-03-13 01:28:21','2026-03-13 01:28:21'),('ebf0e7b7-1e7b-11f1-9193-72232b494de6','GEN_LOG_QUERY_ERR_00002','ERROR','sistema','logging','Error al consultar logs','No se pudieron recuperar los logs desde Loki: ${0}',1,'2026-03-13 01:28:21','2026-03-13 01:28:21'),('ebf0e8c6-1e7b-11f1-9193-72232b494de6','GEN_LOG_RETENTION_WARN_00001','WARNING','sistema','logging','RetenciÃ³n de logs prÃ³xima a lÃ­mite','Los logs almacenados estÃ¡n prÃ³ximos al lÃ­mite de retenciÃ³n de ${0} dÃ­as.',1,'2026-03-13 01:28:21','2026-03-13 01:28:21'),('ebf0e95c-1e7b-11f1-9193-72232b494de6','GEN_LOG_PROMTAIL_ERR_00003','ERROR','sistema','logging','Error en agente Promtail','El agente Promtail no pudo enviar logs a Loki: ${0}',1,'2026-03-13 01:28:21','2026-03-13 01:28:21'),('ebf0e9e3-1e7b-11f1-9193-72232b494de6','GEN_LOG_FORMAT_ERR_00004','ERROR','sistema','logging','Error en formato de log','El formato del log no es vÃ¡lido: ${0}. Se esperaba JSON estructurado.',1,'2026-03-13 01:28:21','2026-03-13 01:28:21'),('ebf0ea63-1e7b-11f1-9193-72232b494de6','GEN_LOG_WRITE_ERR_00005','ERROR','sistema','logging','Error al escribir log','No se pudo escribir el log en el sistema: ${0}',1,'2026-03-13 01:28:21','2026-03-13 01:28:21'),('ebf0eade-1e7b-11f1-9193-72232b494de6','GEN_LOG_GRAFANA_CONN_EXI_00003','EXITO','sistema','logging','ConexiÃ³n a Grafana exitosa','La conexiÃ³n con Grafana para visualizaciÃ³n de logs se estableciÃ³ correctamente.',1,'2026-03-13 01:28:21','2026-03-13 01:28:21'),('ebf0eb5d-1e7b-11f1-9193-72232b494de6','GEN_LOG_GRAFANA_CONN_ERR_00006','ERROR','sistema','logging','Error de conexiÃ³n a Grafana','No se pudo establecer conexiÃ³n con Grafana: ${0}',1,'2026-03-13 01:28:21','2026-03-13 01:28:21'),('ebf0ebd4-1e7b-11f1-9193-72232b494de6','GEN_LOG_DASHBOARD_ERR_00007','ERROR','sistema','logging','Error al cargar dashboard','Error al cargar el dashboard de logs en Grafana: ${0}',1,'2026-03-13 01:28:21','2026-03-13 01:28:21'),('ebf0ec4b-1e7b-11f1-9193-72232b494de6','GEN_LOG_STORAGE_WARN_00002','WARNING','sistema','logging','Almacenamiento de logs elevado','El almacenamiento de logs estÃ¡ al ${0}% de capacidad.',1,'2026-03-13 01:28:21','2026-03-13 01:28:21'),('ebf0ecc0-1e7b-11f1-9193-72232b494de6','GEN_LOG_CONFIG_ERR_00008','ERROR','sistema','logging','Error de configuraciÃ³n','Error en la configuraciÃ³n del sistema de logging: ${0}',1,'2026-03-13 01:28:21','2026-03-13 01:28:21'),('ebf0ed37-1e7b-11f1-9193-72232b494de6','GEN_LOG_LEVEL_INFO_00001','INFO','sistema','logging','Nivel de log actualizado','El nivel de logging ha sido actualizado a: ${0}',1,'2026-03-13 01:28:21','2026-03-13 01:28:21'),('ec271db7-1e7b-11f1-9193-72232b494de6','MOD_KC_VERIF_EMAIL_SENT_EXI_00001','EXITO','usuario_final','keycloak','Â¡Revisa tu correo! ðŸ“§','Te hemos enviado un email de verificaciÃ³n a ${0}. Por favor revisa tu bandeja de entrada (y spam si es necesario) para activar tu cuenta.',1,'2026-03-13 01:28:21','2026-03-13 01:28:21'),('ec2721ed-1e7b-11f1-9193-72232b494de6','MOD_KC_VERIF_EMAIL_ERROR_ERR_00001','ERROR','usuario_final','keycloak','No pudimos enviar el email','Tuvimos un problema al enviar el email de verificaciÃ³n. Por favor intenta nuevamente en unos momentos.',1,'2026-03-13 01:28:21','2026-03-13 01:28:21'),('ec272297-1e7b-11f1-9193-72232b494de6','MOD_KC_EMAIL_ALREADY_VERIFIED_WARN_00001','WARNING','usuario_final','keycloak','Â¡Tu email ya estÃ¡ verificado! âœ“','No te preocupes, tu email ya fue verificado anteriormente. Puedes iniciar sesiÃ³n normalmente.',1,'2026-03-13 01:28:21','2026-03-13 01:28:21'),('ec272321-1e7b-11f1-9193-72232b494de6','MOD_KC_PWD_RESET_SENT_EXI_00001','EXITO','usuario_final','keycloak','Revisa tu correo ðŸ“¬','Si tu email estÃ¡ registrado en nuestro sistema, recibirÃ¡s instrucciones para restablecer tu contraseÃ±a en los prÃ³ximos minutos.',1,'2026-03-13 01:28:21','2026-03-13 01:28:21'),('ec27240b-1e7b-11f1-9193-72232b494de6','MOD_KC_PWD_RESET_ERROR_ERR_00001','ERROR','usuario_final','keycloak','Algo saliÃ³ mal','No pudimos procesar tu solicitud en este momento. Por favor intenta nuevamente mÃ¡s tarde.',1,'2026-03-13 01:28:21','2026-03-13 01:28:21'),('ec272494-1e7b-11f1-9193-72232b494de6','MOD_KC_USER_NOT_FOUND_ERR_00001','ERROR','usuario_final','keycloak','Usuario no encontrado','No encontramos una cuenta con ese email. Verifica que estÃ© escrito correctamente.',1,'2026-03-13 01:28:21','2026-03-13 01:28:21'),('ec27250b-1e7b-11f1-9193-72232b494de6','MOD_KC_USER_SEARCH_ERROR_ERR_00002','ERROR','sistema','keycloak','Error de bÃºsqueda','OcurriÃ³ un error al buscar el usuario en el sistema.',1,'2026-03-13 01:28:21','2026-03-13 01:28:21'),('ec347bdf-1e7b-11f1-9193-72232b494de6','MOD_KC_EMAIL_VERIFIED_EXI_00001','EXITO','usuario_final','keycloak','Â¡Email verificado! âœ“','Tu correo electrÃ³nico ha sido verificado exitosamente. Ya puedes iniciar sesiÃ³n en la aplicaciÃ³n.',1,'2026-03-13 01:28:21','2026-03-13 01:28:21'),('ec348002-1e7b-11f1-9193-72232b494de6','MOD_KC_INVALID_TOKEN_ERR_00001','ERROR','usuario_final','keycloak','Token invÃ¡lido','El enlace de verificaciÃ³n no es vÃ¡lido o estÃ¡ mal formado. Por favor solicita un nuevo email de verificaciÃ³n.',1,'2026-03-13 01:28:21','2026-03-13 01:28:21'),('ec3480b9-1e7b-11f1-9193-72232b494de6','MOD_KC_EMAIL_VERIFY_ERROR_ERR_00001','ERROR','usuario_final','keycloak','Error de verificaciÃ³n','No se pudo verificar tu correo electrÃ³nico. Por favor intenta nuevamente o solicita un nuevo enlace de verificaciÃ³n.',1,'2026-03-13 01:28:21','2026-03-13 01:28:21'),('ec3eb397-1e7b-11f1-9193-72232b494de6','MOD_EMAIL_SEND_INFO_00001','INFO','sistema','email','Enviando Correo de VerificaciÃ³n','Se estÃ¡ enviando el correo electrÃ³nico de verificaciÃ³n a ${0}.',1,'2026-03-13 01:28:21','2026-03-13 01:28:21'),('ec3eb7ed-1e7b-11f1-9193-72232b494de6','MOD_EMAIL_SEND_EXI_00001','EXITO','usuario_final','email','Correo de VerificaciÃ³n Enviado','Se ha enviado un correo de verificaciÃ³n a ${0}. Por favor, revise su bandeja de entrada y carpeta de spam.',1,'2026-03-13 01:28:21','2026-03-13 01:28:21'),('ec3eb8ac-1e7b-11f1-9193-72232b494de6','MOD_EMAIL_SEND_ERR_00001','ERROR','usuario_final','email','Error al Enviar Correo','No pudimos enviar el correo de verificaciÃ³n a ${0}. Por favor, solicite un nuevo envÃ­o o contacte con soporte tÃ©cnico.',1,'2026-03-13 01:28:21','2026-03-13 01:28:21'),('ec3eb933-1e7b-11f1-9193-72232b494de6','MOD_EMAIL_CONFIG_ERR_00002','ERROR','sistema','email','Error de ConfiguraciÃ³n de Email','Error en la configuraciÃ³n del servicio de correo electrÃ³nico: ${0}',1,'2026-03-13 01:28:21','2026-03-13 01:28:21'),('ec3f0861-1e7b-11f1-9193-72232b494de6','MOD_REG_START_INFO_00001','INFO','sistema','registration','Iniciando Registro','Iniciando proceso de registro para el usuario ${0}.',1,'2026-03-13 01:28:21','2026-03-13 01:28:21'),('ec3f0a8f-1e7b-11f1-9193-72232b494de6','MOD_REG_VALID_EXI_00002','EXITO','sistema','registration','ValidaciÃ³n Exitosa','Los datos del usuario han sido validados correctamente.',1,'2026-03-13 01:28:21','2026-03-13 01:28:21'),('ec3f0b37-1e7b-11f1-9193-72232b494de6','MOD_REG_COMPLETE_EXI_00001','EXITO','usuario_final','registration','Registro Completado','Â¡Bienvenido/a ${0}! Su registro ha sido completado exitosamente. Por favor, verifique su correo electrÃ³nico para activar su cuenta.',1,'2026-03-13 01:28:21','2026-03-13 01:28:21'),('ec3f0bd3-1e7b-11f1-9193-72232b494de6','MOD_REG_KC_UPD_INFO_00001','INFO','sistema','registration','Actualizando Keycloak ID','Actualizando el identificador de Keycloak en la base de datos.',1,'2026-03-13 01:28:21','2026-03-13 01:28:21'),('ec3f0c59-1e7b-11f1-9193-72232b494de6','MOD_REG_KC_UPD_EXI_00001','EXITO','sistema','registration','ID de Keycloak Actualizado','El identificador de Keycloak ha sido actualizado correctamente en la base de datos.',1,'2026-03-13 01:28:21','2026-03-13 01:28:21'),('ec3f0ce3-1e7b-11f1-9193-72232b494de6','MOD_REG_KC_UPD_ERR_00001','ERROR','sistema','registration','Error al Actualizar Keycloak ID','No se pudo actualizar el identificador de Keycloak en la base de datos: ${0}',1,'2026-03-13 01:28:21','2026-03-13 01:28:21'),('ec695470-1e7b-11f1-9193-72232b494de6','MOD_P_RESET_EXI_00001','EXITO','usuario_final','person','ContraseÃ±a actualizada','Tu contraseÃ±a ha sido actualizada exitosamente.',1,'2026-03-13 01:28:22','2026-03-13 01:28:22'),('ec6958db-1e7b-11f1-9193-72232b494de6','MOD_P_RESET_ERR_00001','ERROR','usuario_final','person','Token invÃ¡lido','El token de recuperaciÃ³n no es vÃ¡lido o ha expirado. Por favor, solicita un nuevo enlace de recuperaciÃ³n.',1,'2026-03-13 01:28:22','2026-03-13 01:28:22'),('ec69598e-1e7b-11f1-9193-72232b494de6','MOD_P_RESET_ERR_00002','ERROR','usuario_final','person','Usuario no encontrado','No se encontrÃ³ un usuario asociado al token proporcionado.',1,'2026-03-13 01:28:22','2026-03-13 01:28:22'),('ec695a16-1e7b-11f1-9193-72232b494de6','MOD_P_RESET_ERR_00003','ERROR','usuario_final','person','Error al actualizar contraseÃ±a','OcurriÃ³ un error al actualizar tu contraseÃ±a. Por favor, intenta nuevamente.',1,'2026-03-13 01:28:22','2026-03-13 01:28:22'),('ec755039-1e7b-11f1-9193-72232b494de6','MOD_P_CHANGE_EXI_00001','EXITO','usuario_final','person','ContraseÃ±a actualizada','Tu contraseÃ±a ha sido actualizada exitosamente.',1,'2026-03-13 01:28:22','2026-03-13 01:28:22'),('ec755609-1e7b-11f1-9193-72232b494de6','MOD_P_CHANGE_ERR_00001','ERROR','usuario_final','person','ContraseÃ±a actual incorrecta','La contraseÃ±a actual que ingresaste es incorrecta. Por favor, verifica e intenta nuevamente.',1,'2026-03-13 01:28:22','2026-03-13 01:28:22'),('ec755702-1e7b-11f1-9193-72232b494de6','MOD_P_CHANGE_ERR_00002','ERROR','usuario_final','person','Error al actualizar contraseÃ±a','OcurriÃ³ un error al actualizar tu contraseÃ±a. Por favor, intenta nuevamente.',1,'2026-03-13 01:28:22','2026-03-13 01:28:22'),('ec7557a8-1e7b-11f1-9193-72232b494de6','MOD_P_CHANGE_ERR_00003','ERROR','usuario_final','person','ContraseÃ±a no cumple requisitos','La contraseÃ±a debe tener mÃ­nimo 8 caracteres, incluir al menos 1 letra mayÃºscula y 1 carÃ¡cter especial.',1,'2026-03-13 01:28:22','2026-03-13 01:28:22'),('ec84feb0-1e7b-11f1-9193-72232b494de6','MOD_P_UPDATE_ERR_00001','ERROR','usuario_final','person','Error al actualizar perfil','No se pudo actualizar tu perfil. Por favor, intenta nuevamente.',1,'2026-03-13 01:28:22','2026-03-13 01:28:22'),('ec887ee8-1e7b-11f1-9193-72232b494de6','MOD_P_DUP_PHONE_ERR_00001','ERROR','usuario_final','person','NÃºmero de telÃ©fono duplicado','El nÃºmero de telÃ©fono ya estÃ¡ registrado por otro usuario.',1,'2026-03-13 01:28:22','2026-03-13 01:28:22'),('ec88beb3-1e7b-11f1-9193-72232b494de6','MOD_P_DUP_ID_ERR_00001','ERROR','usuario_final','person','NÃºmero de identidad duplicado','El nÃºmero de identidad ya estÃ¡ registrado por otro usuario.',1,'2026-03-13 01:28:22','2026-03-13 01:28:22'),('ecd444bb-1e7b-11f1-9193-72232b494de6','MOD_F_REG_EXI_00001','EXITO','end_user','franchises','Franquicia Registrada','La franquicia ha sido registrada exitosamente.',1,'2026-03-13 01:28:22','2026-03-13 01:28:22'),('ecd44883-1e7b-11f1-9193-72232b494de6','MOD_F_GET_EXI_00001','EXITO','end_user','franchises','Franquicia Encontrada','InformaciÃ³n de la franquicia obtenida correctamente.',1,'2026-03-13 01:28:22','2026-03-13 01:28:22'),('ecd44940-1e7b-11f1-9193-72232b494de6','MOD_F_LIST_EXI_00001','EXITO','end_user','franchises','Franquicias Listadas','Lista de franquicias obtenida correctamente.',1,'2026-03-13 01:28:22','2026-03-13 01:28:22'),('ecd449d8-1e7b-11f1-9193-72232b494de6','MOD_F_UPD_EXI_00001','EXITO','end_user','franchises','Franquicia Actualizada','La informaciÃ³n de la franquicia ha sido actualizada.',1,'2026-03-13 01:28:22','2026-03-13 01:28:22'),('ecd44a65-1e7b-11f1-9193-72232b494de6','MOD_F_DEL_EXI_00001','EXITO','end_user','franchises','Franquicia Eliminada','La franquicia ha sido eliminada correctamente.',1,'2026-03-13 01:28:22','2026-03-13 01:28:22'),('ecd44af1-1e7b-11f1-9193-72232b494de6','MOD_F_NOT_FOUND_ERR_00001','ERROR','end_user','franchises','Franquicia No Encontrada','La franquicia solicitada no existe o no tienes acceso a ella.',1,'2026-03-13 01:28:22','2026-03-13 01:28:22'),('ecd44b76-1e7b-11f1-9193-72232b494de6','MOD_F_DUP_NAME_ERR_00001','ERROR','end_user','franchises','Nombre Duplicado','Ya existe una franquicia con ese nombre.',1,'2026-03-13 01:28:22','2026-03-13 01:28:22'),('ecd44bfd-1e7b-11f1-9193-72232b494de6','MOD_F_NO_BRANCHES_ERR_00001','ERROR','end_user','franchises','Sedes Requeridas','Debes asociar al menos una sede para crear una franquicia.',1,'2026-03-13 01:28:22','2026-03-13 01:28:22'),('ecd44c77-1e7b-11f1-9193-72232b494de6','MOD_F_BRANCH_NOT_OWNED_ERR_00001','ERROR','end_user','franchises','Sede No Autorizada','Una o mÃ¡s sedes no te pertenecen o ya estÃ¡n asociadas a otra franquicia.',1,'2026-03-13 01:28:22','2026-03-13 01:28:22'),('ecd44cf3-1e7b-11f1-9193-72232b494de6','MOD_F_HAS_BRANCHES_ERR_00001','ERROR','end_user','franchises','Franquicia Con Sedes','No se puede eliminar una franquicia que tiene sedes asociadas.',1,'2026-03-13 01:28:22','2026-03-13 01:28:22'),('ecd44d6e-1e7b-11f1-9193-72232b494de6','MOD_F_BRANCH_ADD_EXI_00001','EXITO','end_user','franchises','Sede Vinculada','La sede ha sido vinculada a la franquicia exitosamente.',1,'2026-03-13 01:28:22','2026-03-13 01:28:22'),('ecd44de5-1e7b-11f1-9193-72232b494de6','MOD_F_BRANCH_REM_EXI_00001','EXITO','end_user','franchises','Sede Desvinculada','La sede ha sido desvinculada de la franquicia exitosamente.',1,'2026-03-13 01:28:22','2026-03-13 01:28:22'),('ecd44e64-1e7b-11f1-9193-72232b494de6','MOD_F_MIN_BRANCHES_ERR_00001','ERROR','end_user','franchises','MÃ­nimo Una Sede','No se puede desvincular la Ãºltima sede. La franquicia debe tener al menos una sede.',1,'2026-03-13 01:28:22','2026-03-13 01:28:22'),('ed0a46a3-1e7b-11f1-9193-72232b494de6','MOD_H_CREATE_EXI_00001','EXITO','end_user','schedules','Horario Configurado','El horario de atenciÃ³n de la sede ha sido configurado exitosamente.',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed0a4a95-1e7b-11f1-9193-72232b494de6','MOD_H_GET_EXI_00001','EXITO','end_user','schedules','Horario Consultado','InformaciÃ³n del horario de atenciÃ³n obtenida correctamente.',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed0a4b47-1e7b-11f1-9193-72232b494de6','MOD_H_UPDATE_EXI_00001','EXITO','end_user','schedules','Horario Actualizado','El horario de atenciÃ³n ha sido actualizado exitosamente.',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed0a4bce-1e7b-11f1-9193-72232b494de6','MOD_H_DELETE_EXI_00001','EXITO','end_user','schedules','Horario Eliminado','El horario de atenciÃ³n ha sido eliminado correctamente.',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed0a4c52-1e7b-11f1-9193-72232b494de6','MOD_H_ACTIV_EXI_00001','EXITO','end_user','schedules','Horario Activado','Â¡Tu horario de atenciÃ³n ahora es visible para los clientes!',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed0a4cd3-1e7b-11f1-9193-72232b494de6','MOD_H_DEACT_EXI_00001','EXITO','end_user','schedules','Horario Desactivado','El horario de atenciÃ³n ha sido ocultado temporalmente.',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed0a4d50-1e7b-11f1-9193-72232b494de6','MOD_H_DAYS_EXI_00001','EXITO','end_user','schedules','DÃ­as Disponibles','CatÃ¡logo de dÃ­as de la semana obtenido correctamente.',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed0a4dce-1e7b-11f1-9193-72232b494de6','MOD_H_NOT_FOUND_ERR_00001','ERROR','end_user','schedules','Horario No Encontrado','Esta sede aÃºn no tiene un horario de atenciÃ³n configurado.',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed0a4e4b-1e7b-11f1-9193-72232b494de6','MOD_H_EXISTS_ERR_00001','ERROR','end_user','schedules','Horario Ya Existe','Esta sede ya tiene un horario de atenciÃ³n configurado.',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed0a4ec4-1e7b-11f1-9193-72232b494de6','MOD_H_DAY_ERR_00001','ERROR','end_user','schedules','DÃ­a InvÃ¡lido','El dÃ­a seleccionado no es vÃ¡lido. Por favor selecciona un dÃ­a de la semana.',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed0a4f3d-1e7b-11f1-9193-72232b494de6','MOD_H_TIME_ERR_00001','ERROR','end_user','schedules','Hora InvÃ¡lida','El formato de hora no es vÃ¡lido. Usa el formato HH:MM (ej: 08:00).',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed0a4fb2-1e7b-11f1-9193-72232b494de6','MOD_H_TIME_ORDER_ERR_00001','ERROR','end_user','schedules','Hora de Cierre InvÃ¡lida','La hora de cierre debe ser posterior a la hora de apertura.',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed0a86aa-1e7b-11f1-9193-72232b494de6','MOD_H_INACTIVE_ERR_00001','ERROR','end_user','schedules','Horario Desactivado','No puedes modificar un horario desactivado. ActÃ­valo primero.',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed151907-1e7b-11f1-9193-72232b494de6','MOD_HD_CREATE_EXI_00001','EXITO','end_user','schedule_detail','Horario Registrado','La franja horaria ha sido registrada exitosamente',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed157fe1-1e7b-11f1-9193-72232b494de6','MOD_HD_GET_EXI_00001','EXITO','end_user','schedule_detail','Horario Encontrado','Se encontrÃ³ la informaciÃ³n de la franja horaria',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed15ab54-1e7b-11f1-9193-72232b494de6','MOD_HD_UPDATE_EXI_00001','EXITO','end_user','schedule_detail','Horario Actualizado','La franja horaria ha sido actualizada exitosamente',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed15eea4-1e7b-11f1-9193-72232b494de6','MOD_HD_DELETE_EXI_00001','EXITO','end_user','schedule_detail','Horario Eliminado','La franja horaria ha sido eliminada exitosamente',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed162a90-1e7b-11f1-9193-72232b494de6','MOD_HD_LIST_EXI_00001','EXITO','end_user','schedule_detail','Horarios Listados','Se listaron las franjas horarias exitosamente',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed166440-1e7b-11f1-9193-72232b494de6','MOD_HD_NOT_FOUND_ERR_00001','ERROR','end_user','schedule_detail','Horario No Encontrado','No se encontrÃ³ la franja horaria solicitada',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed169e28-1e7b-11f1-9193-72232b494de6','MOD_HD_CONFLICT_ERR_00001','ERROR','end_user','schedule_detail','Conflicto de Horario','Ya existe una franja horaria que se superpone con la solicitada para este dÃ­a',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed16cc93-1e7b-11f1-9193-72232b494de6','MOD_HD_TIME_ERR_00001','ERROR','end_user','schedule_detail','Hora InvÃ¡lida','El formato de hora debe ser HH:MM (ej: 08:00) y la hora de cierre debe ser posterior a la de apertura',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed16f9b0-1e7b-11f1-9193-72232b494de6','MOD_HD_DAY_ERR_00001','ERROR','end_user','schedule_detail','DÃ­a InvÃ¡lido','El dÃ­a de la semana debe estar entre 1 (Lunes) y 7 (Domingo)',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed20a8a1-1e7b-11f1-9193-72232b494de6','MOD_EH_CREATE_EXI_00001','EXITO','end_user','schedule_exception','ExcepciÃ³n Registrada','La excepciÃ³n de horario ha sido registrada exitosamente',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed210091-1e7b-11f1-9193-72232b494de6','MOD_EH_GET_EXI_00001','EXITO','end_user','schedule_exception','ExcepciÃ³n Encontrada','Se encontrÃ³ la informaciÃ³n de la excepciÃ³n de horario',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed2157f4-1e7b-11f1-9193-72232b494de6','MOD_EH_LIST_EXI_00001','EXITO','end_user','schedule_exception','Excepciones Listadas','Se listaron las excepciones de horario exitosamente',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed221358-1e7b-11f1-9193-72232b494de6','MOD_EH_UPDATE_EXI_00001','EXITO','end_user','schedule_exception','ExcepciÃ³n Actualizada','La excepciÃ³n de horario ha sido actualizada exitosamente',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed224d21-1e7b-11f1-9193-72232b494de6','MOD_EH_DELETE_EXI_00001','EXITO','end_user','schedule_exception','ExcepciÃ³n Eliminada','La excepciÃ³n de horario ha sido eliminada exitosamente',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed227b9e-1e7b-11f1-9193-72232b494de6','MOD_EH_ACTIVATE_EXI_00001','EXITO','end_user','schedule_exception','ExcepciÃ³n Activada','La excepciÃ³n de horario ha sido activada exitosamente',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed22a737-1e7b-11f1-9193-72232b494de6','MOD_EH_DEACTIVATE_EXI_00001','EXITO','end_user','schedule_exception','ExcepciÃ³n Desactivada','La excepciÃ³n de horario ha sido desactivada exitosamente',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed22daef-1e7b-11f1-9193-72232b494de6','MOD_EH_NOT_FOUND_ERR_00001','ERROR','end_user','schedule_exception','ExcepciÃ³n No Encontrada','No se encontrÃ³ la excepciÃ³n de horario solicitada',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed2300b3-1e7b-11f1-9193-72232b494de6','MOD_EH_DATE_CONFLICT_ERR_00001','ERROR','end_user','schedule_exception','Fecha Duplicada','Ya existe una excepciÃ³n de horario para esta fecha',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed2324d7-1e7b-11f1-9193-72232b494de6','MOD_EH_DATE_PAST_ERR_00001','ERROR','end_user','schedule_exception','Fecha Pasada','No se puede crear una excepciÃ³n para una fecha pasada',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed234f14-1e7b-11f1-9193-72232b494de6','MOD_EH_TIME_ERR_00001','ERROR','end_user','schedule_exception','Hora InvÃ¡lida','El formato de hora debe ser HH:MM y la hora de cierre debe ser posterior a la de apertura',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed2cb99e-1e7b-11f1-9193-72232b494de6','MOD_AUTH_REFRESH_SUCCESS_EXI_00001','EXITO','usuario_final','authentication','SesiÃ³n Renovada','Tu sesiÃ³n ha sido renovada exitosamente.',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed49324b-1e7b-11f1-9193-72232b494de6','MOD_V_GEO_LAT_REQ_ERR_00001','ERROR','usuario_final','validation','Latitud Requerida','La latitud es requerida para realizar la bÃºsqueda.',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed49381b-1e7b-11f1-9193-72232b494de6','MOD_V_GEO_LAT_INV_ERR_00001','ERROR','usuario_final','validation','Latitud InvÃ¡lida','La latitud debe estar entre -90 y 90 grados.',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed495ec6-1e7b-11f1-9193-72232b494de6','MOD_V_GEO_LNG_REQ_ERR_00001','ERROR','usuario_final','validation','Longitud Requerida','La longitud es requerida para realizar la bÃºsqueda.',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed495fa5-1e7b-11f1-9193-72232b494de6','MOD_V_GEO_LNG_INV_ERR_00001','ERROR','usuario_final','validation','Longitud InvÃ¡lida','La longitud debe estar entre -180 y 180 grados.',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed496043-1e7b-11f1-9193-72232b494de6','MOD_V_GEO_RAD_INV_ERR_00001','ERROR','usuario_final','validation','Radio de BÃºsqueda InvÃ¡lido','El radio de bÃºsqueda debe ser mayor a 0 y menor o igual a 50 kilÃ³metros.',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed499eb0-1e7b-11f1-9193-72232b494de6','MOD_B_NEARBY_EXI_00001','EXITO','usuario_final','branches','Sedes Cercanas','Se encontraron sedes cercanas a tu ubicaciÃ³n.',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed540dbd-1e7b-11f1-9193-72232b494de6','MOD_MOT_REF_LIST_EXI_00001','EXITO','usuario_final','motorcycles','Referencias Obtenidas','El catÃ¡logo de referencias de motocicletas se obtuvo correctamente.',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed54128b-1e7b-11f1-9193-72232b494de6','MOD_MOT_REF_REQ_ERR_00001','ERROR','usuario_final','motorcycles','Referencia Requerida','Debe seleccionar una marca y modelo de motocicleta del catÃ¡logo.',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed5ca96c-1e7b-11f1-9193-72232b494de6','MOD_MOT_CAT_LIST_EXI_00001','EXITO','usuario_final','motorcycles','CategorÃ­as Obtenidas','El catÃ¡logo de categorÃ­as de motocicletas se obtuvo correctamente.',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed5cadf4-1e7b-11f1-9193-72232b494de6','MOD_MOT_CAT_LINES_EXI_00001','EXITO','usuario_final','motorcycles','LÃ­neas de CategorÃ­a Obtenidas','Las lÃ­neas de motocicletas para la categorÃ­a seleccionada se obtuvieron correctamente.',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed65b267-1e7b-11f1-9193-72232b494de6','MOD_MOT_DISP_LIST_EXI_00001','EXITO','usuario_final','motorcycles','Rangos de Cilindraje Obtenidos','Los rangos de cilindraje han sido obtenidos exitosamente.',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed7099fd-1e7b-11f1-9193-72232b494de6','MOD_MOT_IMG_UPDATE_EXI_00001','EXITO','usuario_final','motorcycles','Imagen Actualizada','La imagen de perfil de su motocicleta se actualizÃ³ correctamente.',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed709e71-1e7b-11f1-9193-72232b494de6','MOD_MOT_IMG_GET_EXI_00001','EXITO','usuario_final','motorcycles','Imagen Obtenida','La imagen de perfil de la motocicleta se obtuvo correctamente.',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed709f3d-1e7b-11f1-9193-72232b494de6','MOD_MOT_IMG_DELETE_EXI_00001','EXITO','usuario_final','motorcycles','Imagen Eliminada','La imagen de perfil de su motocicleta se eliminÃ³ correctamente.',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed709fd4-1e7b-11f1-9193-72232b494de6','MOD_MOT_IMG_UPDATE_ERR_00001','ERROR','usuario_final','motorcycles','Error al Actualizar Imagen','No se pudo actualizar la imagen de perfil de su motocicleta. Por favor, intente nuevamente.',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed70a0c3-1e7b-11f1-9193-72232b494de6','MOD_MOT_IMG_NOT_FOUND_ERR_00001','ERROR','usuario_final','motorcycles','Imagen No Encontrada','La motocicleta no tiene una imagen de perfil configurada.',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed70a15d-1e7b-11f1-9193-72232b494de6','MOD_MOT_IMG_URL_REQ_ERR_00001','ERROR','usuario_final','motorcycles','URL de Imagen Requerida','Debe proporcionar una URL vÃ¡lida para la imagen de perfil.',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed81fdc6-1e7b-11f1-9193-72232b494de6','MOD_EVD_CREATE_EXI_00001','EXITO','usuario_final','evidence','Evidencia Cargada','La evidencia fotogrÃ¡fica ha sido cargada exitosamente.',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed82126e-1e7b-11f1-9193-72232b494de6','MOD_EVD_GET_EXI_00001','EXITO','usuario_final','evidence','Evidencia Encontrada','La informaciÃ³n de la evidencia se obtuvo correctamente.',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed821371-1e7b-11f1-9193-72232b494de6','MOD_EVD_UPDATE_EXI_00001','EXITO','usuario_final','evidence','Evidencia Actualizada','La evidencia fotogrÃ¡fica ha sido actualizada exitosamente.',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed82140e-1e7b-11f1-9193-72232b494de6','MOD_EVD_DELETE_EXI_00001','EXITO','usuario_final','evidence','Evidencia Eliminada','La evidencia fotogrÃ¡fica ha sido eliminada exitosamente.',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed8214be-1e7b-11f1-9193-72232b494de6','MOD_EVD_LIST_EXI_00001','EXITO','usuario_final','evidence','Evidencias Listadas','Las evidencias de la motocicleta se obtuvieron correctamente.',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed82157b-1e7b-11f1-9193-72232b494de6','MOD_EVD_NOT_FOUND_ERR_00001','ERROR','usuario_final','evidence','Evidencia No Encontrada','La evidencia solicitada no fue encontrada.',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed821607-1e7b-11f1-9193-72232b494de6','MOD_EVD_CREATE_ERR_00001','ERROR','usuario_final','evidence','Error al Cargar Evidencia','OcurriÃ³ un error al cargar la evidencia fotogrÃ¡fica. Por favor, intente nuevamente.',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed8216b6-1e7b-11f1-9193-72232b494de6','MOD_EVD_UPDATE_ERR_00001','ERROR','usuario_final','evidence','Error al Actualizar','OcurriÃ³ un error al actualizar la evidencia fotogrÃ¡fica. Por favor, intente nuevamente.',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed82173f-1e7b-11f1-9193-72232b494de6','MOD_EVD_DELETE_ERR_00001','ERROR','usuario_final','evidence','Error al Eliminar','OcurriÃ³ un error al eliminar la evidencia fotogrÃ¡fica. Por favor, intente nuevamente.',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed8217e5-1e7b-11f1-9193-72232b494de6','MOD_EVD_LIMIT_ERR_00001','ERROR','usuario_final','evidence','LÃ­mite de ImÃ¡genes Alcanzado','Ha alcanzado el lÃ­mite mÃ¡ximo de 5 evidencias fotogrÃ¡ficas por motocicleta.',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed82186c-1e7b-11f1-9193-72232b494de6','MOD_EVD_URL_ERR_00001','ERROR','usuario_final','evidence','URL de Imagen InvÃ¡lida','La URL de la imagen proporcionada no es vÃ¡lida. AsegÃºrese de subir la imagen primero.',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('ed8218f4-1e7b-11f1-9193-72232b494de6','MOD_EVD_ANGLE_ERR_00001','ERROR','usuario_final','evidence','Ãngulo de Foto InvÃ¡lido','El Ã¡ngulo de la foto debe ser uno de: frontal, lateral, trasera.',1,'2026-03-13 01:28:23','2026-03-13 01:28:23'),('eda0d884-1e7b-11f1-9193-72232b494de6','MOD_DGP_GRANT_EXI_00001','EXITO','usuario_final','diagnostic_permission','Permiso Concedido','Â¡Listo! La sede ahora puede ver los detalles de tu moto para atenderte mejor.',1,'2026-03-13 01:28:24','2026-03-13 01:28:24'),('eda0dd55-1e7b-11f1-9193-72232b494de6','MOD_DGP_REVOKE_EXI_00001','EXITO','usuario_final','diagnostic_permission','Permiso Revocado','El permiso ha sido revocado. La sede ya no podrÃ¡ ver los detalles de tu moto.',1,'2026-03-13 01:28:24','2026-03-13 01:28:24'),('eda0de1a-1e7b-11f1-9193-72232b494de6','MOD_DGP_LIST_EXI_00001','EXITO','usuario_final','diagnostic_permission','Permisos Obtenidos','Los permisos de diagnÃ³stico se obtuvieron correctamente.',1,'2026-03-13 01:28:24','2026-03-13 01:28:24'),('eda0deb3-1e7b-11f1-9193-72232b494de6','MOD_DGP_NOT_FOUND_ERR_00001','ERROR','usuario_final','diagnostic_permission','Permiso No Encontrado','No se encontrÃ³ un permiso activo para esta sede.',1,'2026-03-13 01:28:24','2026-03-13 01:28:24'),('eda0df76-1e7b-11f1-9193-72232b494de6','MOD_DGP_SAVE_ERR_00001','ERROR','usuario_final','diagnostic_permission','Error al Conceder Permiso','No fue posible conceder el permiso. Por favor, intenta nuevamente.',1,'2026-03-13 01:28:24','2026-03-13 01:28:24'),('eda0e004-1e7b-11f1-9193-72232b494de6','MOD_DGP_DELETE_ERR_00001','ERROR','usuario_final','diagnostic_permission','Error al Revocar Permiso','No fue posible revocar el permiso. Por favor, intenta nuevamente.',1,'2026-03-13 01:28:24','2026-03-13 01:28:24'),('edbc283e-1e7b-11f1-9193-72232b494de6','MOD_MOT_RATE_LIST_EXI_00001','EXITO','usuario_final','motorcycles','Rangos de CalificaciÃ³n Obtenidos','Los rangos de calificaciÃ³n han sido obtenidos exitosamente.',1,'2026-03-13 01:28:24','2026-03-13 01:28:24'),('edc707e8-1e7b-11f1-9193-72232b494de6','MOD_CS_CREATE_EXI_00001','EXITO','usuario_final','completed_service','Servicio Registrado','El servicio realizado ha sido registrado exitosamente.',1,'2026-03-13 01:28:24','2026-03-13 01:28:24'),('edc70c23-1e7b-11f1-9193-72232b494de6','MOD_CS_GET_EXI_00001','EXITO','usuario_final','completed_service','Servicio Encontrado','La informaciÃ³n del servicio realizado se obtuvo correctamente.',1,'2026-03-13 01:28:24','2026-03-13 01:28:24'),('edc70d1c-1e7b-11f1-9193-72232b494de6','MOD_CS_LIST_EXI_00001','EXITO','usuario_final','completed_service','Servicios Obtenidos','Los servicios realizados se obtuvieron correctamente.',1,'2026-03-13 01:28:24','2026-03-13 01:28:24'),('edc70db2-1e7b-11f1-9193-72232b494de6','MOD_CS_NOT_FOUND_ERR_00001','ERROR','usuario_final','completed_service','Servicio No Encontrado','No se encontrÃ³ el servicio realizado solicitado.',1,'2026-03-13 01:28:24','2026-03-13 01:28:24'),('edc70e60-1e7b-11f1-9193-72232b494de6','MOD_CS_CREATE_ERR_00001','ERROR','usuario_final','completed_service','Error al Registrar Servicio','No fue posible registrar el servicio realizado. Por favor, intenta nuevamente.',1,'2026-03-13 01:28:24','2026-03-13 01:28:24'),('edc70eeb-1e7b-11f1-9193-72232b494de6','MOD_CS_BRANCH_SVC_ERR_00001','ERROR','usuario_final','completed_service','Servicios No VÃ¡lidos','Los servicios seleccionados no pertenecen a la sede indicada.',1,'2026-03-13 01:28:24','2026-03-13 01:28:24'),('edc70f76-1e7b-11f1-9193-72232b494de6','MOD_CS_DGN_MOTO_ERR_00001','ERROR','usuario_final','completed_service','DiagnÃ³stico No Corresponde','El diagnÃ³stico seleccionado no corresponde a la motocicleta indicada.',1,'2026-03-13 01:28:24','2026-03-13 01:28:24'),('edc70ff3-1e7b-11f1-9193-72232b494de6','MOD_CS_ACTIVE_ERR_00001','ERROR','usuario_final','completed_service','Servicio Activo Existente','Ya existe un servicio pendiente o en proceso para esta motocicleta en esta sede. Finalice o cancele el servicio actual antes de registrar uno nuevo.',1,'2026-03-13 01:28:24','2026-03-13 01:28:24'),('edd2f472-1e7b-11f1-9193-72232b494de6','GEN_RATE_LIMIT_ERR_00004','ERROR','usuario_final','general','Demasiadas solicitudes','Ha superado el lÃ­mite de solicitudes permitidas. Por favor, espere unos momentos antes de intentar nuevamente.',1,'2026-03-13 01:28:24','2026-03-13 01:28:24'),('eddce4e7-1e7b-11f1-9193-72232b494de6','MOD_CS_DEL_EXI_00001','EXITO','usuario_final','completed_service','Servicio Eliminado','El servicio realizado ha sido eliminado exitosamente.',1,'2026-03-13 01:28:24','2026-03-13 01:28:24'),('eddce977-1e7b-11f1-9193-72232b494de6','MOD_CS_DEL_ERR_00001','ERROR','usuario_final','completed_service','Error al Eliminar Servicio','No fue posible eliminar el servicio realizado. Por favor, intenta nuevamente.',1,'2026-03-13 01:28:24','2026-03-13 01:28:24'),('eddceb22-1e7b-11f1-9193-72232b494de6','MOD_CS_TRANS_EXI_00001','EXITO','usuario_final','completed_service','Transiciones Obtenidas','Las transiciones de estado del servicio se obtuvieron correctamente.',1,'2026-03-13 01:28:24','2026-03-13 01:28:24'),('eddcebc2-1e7b-11f1-9193-72232b494de6','MOD_CS_TRANS_ERR_00001','ERROR','usuario_final','completed_service','Error al Consultar Transiciones','No fue posible obtener las transiciones de estado del servicio.',1,'2026-03-13 01:28:24','2026-03-13 01:28:24'),('eddcec7c-1e7b-11f1-9193-72232b494de6','MOD_CS_STATUS_EXI_00001','EXITO','usuario_final','completed_service','Estado Actualizado','El estado del servicio realizado ha sido actualizado exitosamente.',1,'2026-03-13 01:28:24','2026-03-13 01:28:24'),('eddced08-1e7b-11f1-9193-72232b494de6','MOD_CS_STATUS_ERR_00001','ERROR','usuario_final','completed_service','TransiciÃ³n de Estado No VÃ¡lida','La transiciÃ³n de estado solicitada no es vÃ¡lida para el estado actual del servicio.',1,'2026-03-13 01:28:24','2026-03-13 01:28:24'),('edee5e32-1e7b-11f1-9193-72232b494de6','MOD_CS_DETAILS_UPD_EXI_00001','EXITO','usuario_final','completed_service','Detalles Actualizados','Los detalles del servicio realizado han sido actualizados exitosamente.',1,'2026-03-13 01:28:24','2026-03-13 01:28:24'),('edee62b0-1e7b-11f1-9193-72232b494de6','MOD_CS_UPDATE_ERR_00001','ERROR','usuario_final','completed_service','No es Posible Editar el Servicio','No es posible editar los detalles del servicio en su estado actual. Solo se pueden editar servicios en estado Pendiente o En Proceso.',1,'2026-03-13 01:28:24','2026-03-13 01:28:24'),('edf78638-1e7b-11f1-9193-72232b494de6','MOD_CS_RATE_EXI_00001','EXITO','usuario_final','rating','CalificaciÃ³n Registrada','La calificaciÃ³n del servicio ha sido registrada exitosamente.',1,'2026-03-13 01:28:24','2026-03-13 01:28:24'),('edf78b0a-1e7b-11f1-9193-72232b494de6','MOD_CS_RATE_OFFENSIVE_EXI_00001','EXITO','usuario_final','rating','CalificaciÃ³n Registrada','Tu calificaciÃ³n fue registrada. Sin embargo, el comentario contiene lenguaje inapropiado y no serÃ¡ visible para otros usuarios.',1,'2026-03-13 01:28:24','2026-03-13 01:28:24'),('edf78bfe-1e7b-11f1-9193-72232b494de6','MOD_CS_RATE_LIST_EXI_00001','EXITO','usuario_final','rating','Calificaciones Obtenidas','Las calificaciones del servicio han sido obtenidas exitosamente.',1,'2026-03-13 01:28:24','2026-03-13 01:28:24'),('edf78c90-1e7b-11f1-9193-72232b494de6','MOD_CS_ITEM_NF_ERR_00001','ERROR','usuario_final','rating','Ãtem de Servicio No Encontrado','El Ã­tem de servicio que desea calificar no fue encontrado.',1,'2026-03-13 01:28:24','2026-03-13 01:28:24'),('edf78d1b-1e7b-11f1-9193-72232b494de6','MOD_CS_RATE_DUP_ERR_00001','ERROR','usuario_final','rating','Servicio Ya Calificado','Este servicio ya fue calificado. Solo se permite una calificaciÃ³n por servicio.',1,'2026-03-13 01:28:24','2026-03-13 01:28:24'),('edf78da6-1e7b-11f1-9193-72232b494de6','MOD_CS_RATE_RANGE_ERR_00001','ERROR','usuario_final','rating','CalificaciÃ³n Fuera de Rango','La calificaciÃ³n debe ser un valor entre 1 y 5.',1,'2026-03-13 01:28:24','2026-03-13 01:28:24'),('edf78e2d-1e7b-11f1-9193-72232b494de6','MOD_CS_RATE_STATUS_ERR_00001','ERROR','usuario_final','rating','Servicio No Finalizado','Solo se pueden calificar servicios que hayan sido finalizados.',1,'2026-03-13 01:28:24','2026-03-13 01:28:24'),('edf78eb0-1e7b-11f1-9193-72232b494de6','MOD_CS_RATE_SAVE_ERR_00001','ERROR','usuario_final','rating','Error al Guardar CalificaciÃ³n','No fue posible guardar la calificaciÃ³n. Por favor, intente de nuevo mÃ¡s tarde.',1,'2026-03-13 01:28:24','2026-03-13 01:28:24'),('f8f92faa-1e7b-11f1-9193-72232b494de6','MOD_INFRA_KC_UNAVAIL_ERR_00004','ERROR','usuario_final','infrastructure','Servicio de AutenticaciÃ³n No Disponible','El servicio de autenticaciÃ³n no estÃ¡ disponible en este momento. Por favor intente mÃ¡s tarde.',1,'2026-03-13 01:28:43','2026-03-13 01:28:43'),('f8f9471a-1e7b-11f1-9193-72232b494de6','MOD_INFRA_DB_UNAVAIL_ERR_00005','ERROR','usuario_final','infrastructure','Servicio Temporalmente No Disponible','El sistema no puede procesar su solicitud en este momento. Por favor intente mÃ¡s tarde.',1,'2026-03-13 01:28:43','2026-03-13 01:28:43'),('f8f94925-1e7b-11f1-9193-72232b494de6','MOD_INFRA_DEP_FAIL_ERR_00006','ERROR','usuario_final','infrastructure','Servicio Temporalmente No Disponible','No se pudo completar la operaciÃ³n debido a un problema tÃ©cnico temporal.',1,'2026-03-13 01:28:43','2026-03-13 01:28:43'),('f8f949da-1e7b-11f1-9193-72232b494de6','MOD_INFRA_KC_USER_EXISTS_ERR_00007','ERROR','usuario_final','infrastructure','Usuario Ya Registrado','El correo electrÃ³nico ya estÃ¡ registrado en el sistema. Por favor use otro correo electrÃ³nico.',1,'2026-03-13 01:28:43','2026-03-13 01:28:43'),('f8f94af1-1e7b-11f1-9193-72232b494de6','MOD_INFRA_DB_USER_EXISTS_ERR_00008','ERROR','usuario_final','infrastructure','Usuario Ya Registrado','El correo electrÃ³nico ya estÃ¡ registrado en el sistema. Por favor use otro correo electrÃ³nico.',1,'2026-03-13 01:28:43','2026-03-13 01:28:43'),('f8f94b89-1e7b-11f1-9193-72232b494de6','MOD_INFRA_INCOMPLETE_REG_ERR_00009','ERROR','usuario_final','infrastructure','Registro Incompleto en CorrecciÃ³n','Detectamos un registro incompleto que estamos corrigiendo automÃ¡ticamente. Por favor intente registrarse nuevamente en unos momentos.',1,'2026-03-13 01:28:43','2026-03-13 01:28:43'),('f9057187-1e7b-11f1-9193-72232b494de6','MOD_MOT_NO_PERMISSION_ERR_00001','ERROR','usuario_final','motorcycle','Sin Permiso de DiagnÃ³stico','No tiene permisos activos para ver los diagnÃ³sticos de esta motocicleta. El propietario debe otorgar permiso a su sede.',1,'2026-03-13 01:28:43','2026-03-13 01:28:43');
/*!40000 ALTER TABLE `system_messages` ENABLE KEYS */;
UNLOCK TABLES;
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

LOCK TABLES `transiciones_estado_servicio` WRITE;
/*!40000 ALTER TABLE `transiciones_estado_servicio` DISABLE KEYS */;
/*!40000 ALTER TABLE `transiciones_estado_servicio` ENABLE KEYS */;
UNLOCK TABLES;
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

LOCK TABLES `ubicaciones` WRITE;
/*!40000 ALTER TABLE `ubicaciones` DISABLE KEYS */;
/*!40000 ALTER TABLE `ubicaciones` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

