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
DROP TABLE IF EXISTS `ADMIN_EVENT_ENTITY`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ADMIN_EVENT_ENTITY` (
  `ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ADMIN_EVENT_TIME` bigint DEFAULT NULL,
  `REALM_ID` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `OPERATION_TYPE` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `AUTH_REALM_ID` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `AUTH_CLIENT_ID` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `AUTH_USER_ID` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `IP_ADDRESS` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `RESOURCE_PATH` text COLLATE utf8mb4_unicode_ci,
  `REPRESENTATION` text COLLATE utf8mb4_unicode_ci,
  `ERROR` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `RESOURCE_TYPE` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `DETAILS_JSON` longtext CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  PRIMARY KEY (`ID`),
  KEY `IDX_ADMIN_EVENT_TIME` (`REALM_ID`,`ADMIN_EVENT_TIME`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `ADMIN_EVENT_ENTITY` WRITE;
/*!40000 ALTER TABLE `ADMIN_EVENT_ENTITY` DISABLE KEYS */;
/*!40000 ALTER TABLE `ADMIN_EVENT_ENTITY` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `ASSOCIATED_POLICY`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ASSOCIATED_POLICY` (
  `POLICY_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ASSOCIATED_POLICY_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`POLICY_ID`,`ASSOCIATED_POLICY_ID`),
  KEY `IDX_ASSOC_POL_ASSOC_POL_ID` (`ASSOCIATED_POLICY_ID`),
  CONSTRAINT `FK_FRSR5S213XCX4WNKOG82SSRFY` FOREIGN KEY (`ASSOCIATED_POLICY_ID`) REFERENCES `RESOURCE_SERVER_POLICY` (`ID`),
  CONSTRAINT `FK_FRSRPAS14XCX4WNKOG82SSRFY` FOREIGN KEY (`POLICY_ID`) REFERENCES `RESOURCE_SERVER_POLICY` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `ASSOCIATED_POLICY` WRITE;
/*!40000 ALTER TABLE `ASSOCIATED_POLICY` DISABLE KEYS */;
/*!40000 ALTER TABLE `ASSOCIATED_POLICY` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `AUTHENTICATION_EXECUTION`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `AUTHENTICATION_EXECUTION` (
  `ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ALIAS` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `AUTHENTICATOR` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `REALM_ID` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `FLOW_ID` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `REQUIREMENT` int DEFAULT NULL,
  `PRIORITY` int DEFAULT NULL,
  `AUTHENTICATOR_FLOW` tinyint NOT NULL DEFAULT '0',
  `AUTH_FLOW_ID` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `AUTH_CONFIG` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`ID`),
  KEY `IDX_AUTH_EXEC_REALM_FLOW` (`REALM_ID`,`FLOW_ID`),
  KEY `IDX_AUTH_EXEC_FLOW` (`FLOW_ID`),
  CONSTRAINT `FK_AUTH_EXEC_FLOW` FOREIGN KEY (`FLOW_ID`) REFERENCES `AUTHENTICATION_FLOW` (`ID`),
  CONSTRAINT `FK_AUTH_EXEC_REALM` FOREIGN KEY (`REALM_ID`) REFERENCES `REALM` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `AUTHENTICATION_EXECUTION` WRITE;
/*!40000 ALTER TABLE `AUTHENTICATION_EXECUTION` DISABLE KEYS */;
INSERT INTO `AUTHENTICATION_EXECUTION` VALUES ('07837254-8779-4987-9703-f2f21b44a038',NULL,'registration-page-form','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','26b753d6-6520-4727-b63d-dac95814059f',0,10,1,'50b77434-6401-443e-9e81-1edbf8122259',NULL),('0818d688-e107-40a7-8d89-d476bf18604e',NULL,'reset-password','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','e9deb89c-1118-49eb-83ba-9a9c1dfb9a2a',0,30,0,NULL,NULL),('0a117767-f1ea-4900-8ea4-8d18548f4242',NULL,'client-x509','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','a0edc4f0-693e-42be-8fd9-f6f6da0515a0',2,40,0,NULL,NULL),('0d7484e2-5e5e-4a22-b6ee-b93b0ab27dea',NULL,'reset-credentials-choose-user','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','d0f21bd5-8d4d-4dee-b58e-77f8cd44ad8c',0,10,0,NULL,NULL),('154f9309-f166-46b0-b556-a15e27270970',NULL,'client-jwt','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','de9091e2-e7a5-42f3-aefb-56c5360be436',2,20,0,NULL,NULL),('15603824-04b7-4318-9179-3a64f24e83b5',NULL,'http-basic-authenticator','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','7ba91da4-c6bd-47ce-9f9d-da97f7f72703',0,10,0,NULL,NULL),('17d96f5e-72a0-4c8a-ac7a-88220a9165e4',NULL,'auth-spnego','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','676c5a06-ae2d-44a2-89a8-9017809fa2b7',3,20,0,NULL,NULL),('1a91f016-6888-4ca2-a52b-95f9d5e4fd2f',NULL,'direct-grant-validate-otp','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','37b360cb-ce8c-4e6a-a7d6-d3380139c06d',0,20,0,NULL,NULL),('1d53f533-ca58-4c47-85fa-dbc6028304f2',NULL,'reset-password','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','d0f21bd5-8d4d-4dee-b58e-77f8cd44ad8c',0,30,0,NULL,NULL),('1db2b857-9166-4bce-a93f-768c23330450',NULL,'reset-credentials-choose-user','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','e9deb89c-1118-49eb-83ba-9a9c1dfb9a2a',0,10,0,NULL,NULL),('20dcdf12-26a6-4945-8c29-8dcc2c258691',NULL,'idp-username-password-form','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','c5735bc4-e5c1-4037-ba8e-85d2640c770d',0,10,0,NULL,NULL),('2486b16c-c6f9-4799-a85a-242e678f77ec',NULL,'webauthn-authenticator','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','6348a2e8-c07e-4834-9f6f-431b58891ec0',3,40,0,NULL,NULL),('26400daa-1c83-419d-9fb5-ce9859ffe272',NULL,'auth-otp-form','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','6348a2e8-c07e-4834-9f6f-431b58891ec0',2,30,0,NULL,NULL),('2edc2385-ab46-4e47-8d4a-467c18c1e8b6',NULL,'auth-cookie','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','676c5a06-ae2d-44a2-89a8-9017809fa2b7',2,10,0,NULL,NULL),('2f1daf6b-8dab-4c65-8f32-71696113f397',NULL,'registration-terms-and-conditions','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','50b77434-6401-443e-9e81-1edbf8122259',3,70,0,NULL,NULL),('311e52f5-9373-4bc1-b9c3-62534fc5e6b7',NULL,NULL,'4c2886d8-e4c4-4b40-9064-445c7ab90a1c','79d34d67-cb53-4f81-bbd8-78febfecd274',1,10,1,'c08954fb-ccf3-4181-8864-c0a2ec7d0d84',NULL),('328358a0-a267-4192-a087-17b1af2e83d4',NULL,'docker-http-basic-authenticator','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','216f221b-9be0-40f5-9c68-f5bee839e1bc',0,10,0,NULL,NULL),('32a1b49d-55a9-4e20-a2d4-37350772d871',NULL,'client-secret','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','de9091e2-e7a5-42f3-aefb-56c5360be436',2,10,0,NULL,NULL),('32cc9a64-764a-4419-85a3-ce881ac6d465',NULL,NULL,'4c2886d8-e4c4-4b40-9064-445c7ab90a1c','c5735bc4-e5c1-4037-ba8e-85d2640c770d',1,20,1,'8862e11d-2025-4230-8693-4e0cfff5d962',NULL),('34794a2f-83d1-4380-9772-f0c258c695bb',NULL,'client-secret-jwt','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','a0edc4f0-693e-42be-8fd9-f6f6da0515a0',2,30,0,NULL,NULL),('35913774-4c62-473f-b69c-abcdcefacb84',NULL,'auth-username-password-form','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','b920157a-a297-4663-afda-7bcfe73e349f',0,10,0,NULL,NULL),('36c8dac2-0e0c-4dd2-8e2a-4d263fd2a07b',NULL,'client-jwt','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','a0edc4f0-693e-42be-8fd9-f6f6da0515a0',2,20,0,NULL,NULL),('3879e22d-da7f-40de-95ba-a7b9e8636c73',NULL,'webauthn-authenticator','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','5f3a045d-a67d-43ce-b2e0-9625663e8dd4',3,40,0,NULL,NULL),('3c66d74d-fffb-43dc-bbf5-9b8fbc01ab31',NULL,NULL,'4c2886d8-e4c4-4b40-9064-445c7ab90a1c','06cd7bfa-72e7-48fa-9701-d8e7c0a49633',0,20,1,'2d8f5f5a-cbd8-46b4-a07f-10c9523759c5',NULL),('3ff2f719-9bf9-467c-bbce-4880564f57b0',NULL,'auth-spnego','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','e23ba648-6f5b-447f-8d45-33da2f4f7e9e',3,20,0,NULL,NULL),('49fd8337-64e3-454f-a09a-5ed9b20113f3',NULL,'conditional-user-configured','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','6348a2e8-c07e-4834-9f6f-431b58891ec0',0,10,0,NULL,NULL),('4f688044-818f-4f4b-b5d9-1e81543e672d',NULL,'conditional-credential','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','5f3a045d-a67d-43ce-b2e0-9625663e8dd4',0,20,0,NULL,'f7f468af-378f-496f-9ca5-69dc30098555'),('5133dcf0-4f45-4925-bd6f-f5987eb62dab',NULL,'conditional-credential','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','6348a2e8-c07e-4834-9f6f-431b58891ec0',0,20,0,NULL,'64e20831-27ab-4fd0-ab17-aeec0be646a1'),('51adea09-3247-4fe9-9517-aa1c7e83d6c3',NULL,'conditional-credential','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','8862e11d-2025-4230-8693-4e0cfff5d962',0,20,0,NULL,'2be225dc-227b-4f4a-8e19-1956e8d32f94'),('531254f0-6ad7-44b4-b6bd-c83f36746a93',NULL,'conditional-user-configured','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','c08954fb-ccf3-4181-8864-c0a2ec7d0d84',0,10,0,NULL,NULL),('543b43b7-0288-4658-9954-d8d0136804e5',NULL,'registration-password-action','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','50b77434-6401-443e-9e81-1edbf8122259',0,50,0,NULL,NULL),('5afa7b79-d6cf-4b7e-859f-11ab29ea6b0c',NULL,NULL,'4c2886d8-e4c4-4b40-9064-445c7ab90a1c','e23ba648-6f5b-447f-8d45-33da2f4f7e9e',2,26,1,'79d34d67-cb53-4f81-bbd8-78febfecd274',NULL),('5df1e640-4a9c-452a-843f-99a2e658e9c7',NULL,NULL,'5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','676c5a06-ae2d-44a2-89a8-9017809fa2b7',2,30,1,'7ec6b431-8ec8-46bb-8611-e0ca04703203',NULL),('652c1ccd-c660-4f4f-9cd8-b5673ffbdda3',NULL,NULL,'4c2886d8-e4c4-4b40-9064-445c7ab90a1c','e9deb89c-1118-49eb-83ba-9a9c1dfb9a2a',1,40,1,'7bca3933-de65-4855-8568-7388b2b12092',NULL),('655e33a6-12d0-422a-b620-ae2ada492e12',NULL,'idp-add-organization-member','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','9b4399ab-2f01-4ccf-a46c-6aa764d61f1a',0,20,0,NULL,NULL),('678315a2-66ac-472a-85f9-1d2fab5ba621',NULL,'conditional-user-configured','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','9b4399ab-2f01-4ccf-a46c-6aa764d61f1a',0,10,0,NULL,NULL),('6966a136-8ab9-4cf4-a07d-506d7c208ad3',NULL,'idp-confirm-link','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','8a307303-1570-4786-9d86-061a1e51d66f',0,10,0,NULL,NULL),('6bf8b895-6684-43f0-a09f-c95dc3246bdc',NULL,'registration-password-action','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','7549030d-3672-4d2b-bbf0-ffd7368e4524',0,50,0,NULL,NULL),('72d37736-58a6-480c-8826-7bdc3e6b72b4',NULL,'conditional-user-configured','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','b606b87b-37af-4482-8b2d-1c7a18160d44',0,10,0,NULL,NULL),('732c9369-52d9-4acd-ac00-2a4e4b8f4660',NULL,'registration-user-creation','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','50b77434-6401-443e-9e81-1edbf8122259',0,20,0,NULL,NULL),('74d17584-1866-452f-9ceb-a633d6b67a58',NULL,'registration-recaptcha-action','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','7549030d-3672-4d2b-bbf0-ffd7368e4524',3,60,0,NULL,NULL),('7688eedd-8ebb-4478-8d91-8acb8cfaed64',NULL,NULL,'5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','0f2a187b-6a65-4599-acd3-8404ee73bf62',2,20,1,'75ba53df-d698-4559-88ca-8d9143e67d75',NULL),('77a4ee53-5354-4f99-bb3b-2042cbef952c',NULL,'organization','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','c08954fb-ccf3-4181-8864-c0a2ec7d0d84',2,20,0,NULL,NULL),('790b9e04-ad79-465b-bbf5-7fbb9ea692bd',NULL,'reset-credential-email','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','d0f21bd5-8d4d-4dee-b58e-77f8cd44ad8c',0,20,0,NULL,NULL),('7c4cc7f5-47ed-4e35-98e5-a8941e63cb61',NULL,'registration-page-form','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','49326f20-3ceb-4c4b-8171-8e7d27d44093',0,10,1,'7549030d-3672-4d2b-bbf0-ffd7368e4524',NULL),('7ecc0643-2e19-4910-aa17-ea345e34571f',NULL,NULL,'5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','8503d694-8a57-4077-8624-5d9711b52163',2,20,1,'8a307303-1570-4786-9d86-061a1e51d66f',NULL),('7f102cff-72f6-42d0-ae67-00ed45d5f206',NULL,'identity-provider-redirector','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','e23ba648-6f5b-447f-8d45-33da2f4f7e9e',2,25,0,NULL,NULL),('80c43fd2-b090-4888-bd4d-17337ec46627',NULL,'webauthn-authenticator','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','8862e11d-2025-4230-8693-4e0cfff5d962',3,40,0,NULL,NULL),('80fc9197-d48c-43b4-8556-8f8294ad2de5',NULL,NULL,'5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','d0f21bd5-8d4d-4dee-b58e-77f8cd44ad8c',1,40,1,'0526685c-2944-462b-8f75-855e0c7761a0',NULL),('81ed3c94-7039-4c1e-914c-720115bd4c12',NULL,'auth-otp-form','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','5f3a045d-a67d-43ce-b2e0-9625663e8dd4',2,30,0,NULL,NULL),('840a602b-505d-41b6-ace7-6b5b127fb615',NULL,'conditional-user-configured','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','7bca3933-de65-4855-8568-7388b2b12092',0,10,0,NULL,NULL),('84419b8e-e23f-4974-b095-aa26bb5d0a43',NULL,'auth-recovery-authn-code-form','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','e167e17a-2d12-4dcf-9a7f-c9b3c830be01',3,50,0,NULL,NULL),('84accf05-fb30-4e73-81d9-d7680100a9a2',NULL,'auth-recovery-authn-code-form','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','5f3a045d-a67d-43ce-b2e0-9625663e8dd4',3,50,0,NULL,NULL),('86626fe4-1430-467c-ada3-9c33b093b773',NULL,'docker-http-basic-authenticator','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','906eb7da-9867-4708-9777-30c21db80523',0,10,0,NULL,NULL),('87cfd48a-6cb5-4aff-8fd7-312fd6d0cac4',NULL,NULL,'4c2886d8-e4c4-4b40-9064-445c7ab90a1c','f3d39771-85ea-427d-b233-fa26308a3952',1,60,1,'9b4399ab-2f01-4ccf-a46c-6aa764d61f1a',NULL),('885872a2-2f0d-449e-a0f7-21621a61e1fd',NULL,'auth-recovery-authn-code-form','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','6348a2e8-c07e-4834-9f6f-431b58891ec0',3,50,0,NULL,NULL),('89ab2b99-cde9-40fd-926c-1a192b86f53a',NULL,'conditional-user-configured','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','0526685c-2944-462b-8f75-855e0c7761a0',0,10,0,NULL,NULL),('8ae20e35-dc94-48de-8f13-a73e11fb35e1',NULL,'auth-recovery-authn-code-form','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','8862e11d-2025-4230-8693-4e0cfff5d962',3,50,0,NULL,NULL),('8b11fe4f-dc1c-467d-ab5c-c15ca204147b',NULL,'direct-grant-validate-otp','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','b606b87b-37af-4482-8b2d-1c7a18160d44',0,20,0,NULL,NULL),('8b5ef2ec-aaca-4e72-8aae-f077fe3f1b26',NULL,'client-secret-jwt','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','de9091e2-e7a5-42f3-aefb-56c5360be436',2,30,0,NULL,NULL),('8c1aa08c-5fed-4fd5-80ca-3e0b3587493d',NULL,NULL,'5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','8a307303-1570-4786-9d86-061a1e51d66f',0,20,1,'0f2a187b-6a65-4599-acd3-8404ee73bf62',NULL),('8f47bc47-e147-47fe-a133-1261a820bb08',NULL,'direct-grant-validate-password','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','a3c9dcd5-2130-44c7-bab3-55f37e7c2da6',0,20,0,NULL,NULL),('927964d8-f06c-4658-9a5b-00053a07537b',NULL,NULL,'4c2886d8-e4c4-4b40-9064-445c7ab90a1c','f662b0fb-ef63-4db6-934e-e4d37ca37296',1,30,1,'37b360cb-ce8c-4e6a-a7d6-d3380139c06d',NULL),('95229b9a-42e3-4fbf-80a5-b4eeb5c4d46c',NULL,'registration-user-creation','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','7549030d-3672-4d2b-bbf0-ffd7368e4524',0,20,0,NULL,NULL),('98889fbd-5a88-4918-bd10-5ebf8e7d6ce5',NULL,'conditional-user-configured','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','37b360cb-ce8c-4e6a-a7d6-d3380139c06d',0,10,0,NULL,NULL),('9c55e5c7-d2dd-4aeb-9ce8-9e5b081ffe39',NULL,'idp-email-verification','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','0f2a187b-6a65-4599-acd3-8404ee73bf62',2,10,0,NULL,NULL),('9c6eb4cc-aad9-4ec1-864e-280bb1de77a1',NULL,'registration-recaptcha-action','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','50b77434-6401-443e-9e81-1edbf8122259',3,60,0,NULL,NULL),('a13f0cf1-ddf3-4610-a4ca-8f555a8035bb',NULL,NULL,'4c2886d8-e4c4-4b40-9064-445c7ab90a1c','23d3b2c0-2acd-42f9-b25e-3088d6c1caf3',2,20,1,'06cd7bfa-72e7-48fa-9701-d8e7c0a49633',NULL),('a2092408-ab3c-431c-9909-29c2e2ea2eed',NULL,'conditional-credential','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','e167e17a-2d12-4dcf-9a7f-c9b3c830be01',0,20,0,NULL,'5b99313d-012b-42fe-a42e-2733a3d43904'),('a6c815eb-1e90-4032-8648-7047a08d2b5c',NULL,NULL,'5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','403e0820-f1af-4eff-8cbc-6a5423453d76',0,20,1,'8503d694-8a57-4077-8624-5d9711b52163',NULL),('b0accc07-e9d2-4e2d-a3d1-b71f808f39e5',NULL,NULL,'5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','a3c9dcd5-2130-44c7-bab3-55f37e7c2da6',1,30,1,'b606b87b-37af-4482-8b2d-1c7a18160d44',NULL),('b6e7b083-1d28-45b0-9c6e-75d62911e968',NULL,'registration-terms-and-conditions','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','7549030d-3672-4d2b-bbf0-ffd7368e4524',3,70,0,NULL,NULL),('bb053ec1-0d1b-4192-82a1-7767986bec4b',NULL,'reset-credential-email','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','e9deb89c-1118-49eb-83ba-9a9c1dfb9a2a',0,20,0,NULL,NULL),('bd81181c-0a0c-47ab-9b22-b07db241942b',NULL,'auth-otp-form','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','e167e17a-2d12-4dcf-9a7f-c9b3c830be01',2,30,0,NULL,NULL),('c151d1dd-9e72-43a6-b071-88c1abc38f09',NULL,'identity-provider-redirector','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','676c5a06-ae2d-44a2-89a8-9017809fa2b7',2,25,0,NULL,NULL),('c42376bb-8b2f-4927-b8cb-cae335169005',NULL,'idp-review-profile','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','f3d39771-85ea-427d-b233-fa26308a3952',0,10,0,NULL,'ffc1376a-dddb-4503-aaa6-6008371a716b'),('c75f3810-3ed3-49b5-932c-64f7ae79de49',NULL,'conditional-user-configured','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','5f3a045d-a67d-43ce-b2e0-9625663e8dd4',0,10,0,NULL,NULL),('cca65279-19c1-416b-a66b-97fa172d6eaa',NULL,'conditional-user-configured','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','e167e17a-2d12-4dcf-9a7f-c9b3c830be01',0,10,0,NULL,NULL),('cfa83450-7fb9-4a91-9fde-777d8d681f64',NULL,'direct-grant-validate-password','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','f662b0fb-ef63-4db6-934e-e4d37ca37296',0,20,0,NULL,NULL),('d19c3789-3911-4e6e-a22b-42f67d05a4c6',NULL,'idp-confirm-link','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','06cd7bfa-72e7-48fa-9701-d8e7c0a49633',0,10,0,NULL,NULL),('d1e23cb3-3d77-4b76-b6bd-9e7502ca20e4',NULL,'auth-otp-form','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','8862e11d-2025-4230-8693-4e0cfff5d962',2,30,0,NULL,NULL),('d1ead887-10e0-4c5c-86e7-bc5432425c18',NULL,'client-secret','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','a0edc4f0-693e-42be-8fd9-f6f6da0515a0',2,10,0,NULL,NULL),('d6189ed4-93ac-44f6-834d-52b6546928e8',NULL,NULL,'4c2886d8-e4c4-4b40-9064-445c7ab90a1c','e23ba648-6f5b-447f-8d45-33da2f4f7e9e',2,30,1,'b920157a-a297-4663-afda-7bcfe73e349f',NULL),('d7d3fd69-3124-4afa-9b05-d37fa93ed4a9',NULL,NULL,'4c2886d8-e4c4-4b40-9064-445c7ab90a1c','2d8f5f5a-cbd8-46b4-a07f-10c9523759c5',2,20,1,'c5735bc4-e5c1-4037-ba8e-85d2640c770d',NULL),('d879b7c1-4e57-4c8c-9d49-0d6f233919f3',NULL,NULL,'5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','7ec6b431-8ec8-46bb-8611-e0ca04703203',1,20,1,'e167e17a-2d12-4dcf-9a7f-c9b3c830be01',NULL),('d96442ea-d0af-4eb9-8539-61bbb93578dd',NULL,'direct-grant-validate-username','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','f662b0fb-ef63-4db6-934e-e4d37ca37296',0,10,0,NULL,NULL),('da045cb2-f3bd-4813-98d5-729483ad7171',NULL,NULL,'4c2886d8-e4c4-4b40-9064-445c7ab90a1c','f3d39771-85ea-427d-b233-fa26308a3952',0,20,1,'23d3b2c0-2acd-42f9-b25e-3088d6c1caf3',NULL),('db6726e2-49bb-4b7a-be5d-220be0c08d24',NULL,NULL,'5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','75ba53df-d698-4559-88ca-8d9143e67d75',1,20,1,'6348a2e8-c07e-4834-9f6f-431b58891ec0',NULL),('dca1c608-d5b2-44b0-a8ac-491f3cc7ed1b',NULL,'idp-create-user-if-unique','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','8503d694-8a57-4077-8624-5d9711b52163',2,10,0,NULL,'f231ed1d-aaab-4f62-acaf-5e9de2a3b399'),('e384bf7f-a802-410e-aa0d-a2d5b8199c87',NULL,NULL,'4c2886d8-e4c4-4b40-9064-445c7ab90a1c','b920157a-a297-4663-afda-7bcfe73e349f',1,20,1,'5f3a045d-a67d-43ce-b2e0-9625663e8dd4',NULL),('e470521a-7b4c-41e6-a13a-91cefd13ae30',NULL,'idp-create-user-if-unique','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','23d3b2c0-2acd-42f9-b25e-3088d6c1caf3',2,10,0,NULL,'7117ace8-ea89-4a9b-9b13-d0016374af12'),('e5b56b75-07a0-4f7d-b334-0cb9ecefa65c',NULL,'webauthn-authenticator','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','e167e17a-2d12-4dcf-9a7f-c9b3c830be01',3,40,0,NULL,NULL),('e5c389d3-3359-4e1c-bd34-45d8eacf1b10',NULL,'conditional-user-configured','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','8862e11d-2025-4230-8693-4e0cfff5d962',0,10,0,NULL,NULL),('e95a513f-31a7-4ec3-afb5-51928d8b39c4',NULL,'idp-email-verification','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','2d8f5f5a-cbd8-46b4-a07f-10c9523759c5',2,10,0,NULL,NULL),('efeb0168-75d3-4148-ae87-2ba67a381353',NULL,'direct-grant-validate-username','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','a3c9dcd5-2130-44c7-bab3-55f37e7c2da6',0,10,0,NULL,NULL),('f0077528-7816-4b92-99f7-5aacd434a0b1',NULL,'client-x509','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','de9091e2-e7a5-42f3-aefb-56c5360be436',2,40,0,NULL,NULL),('f19039c0-77f8-4fc5-8b85-60209e730471',NULL,'idp-review-profile','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','403e0820-f1af-4eff-8cbc-6a5423453d76',0,10,0,NULL,'b6ebb592-9ed6-431a-9343-e37b1f0a0b91'),('f6dafb2b-0e74-4c64-8fb2-113e97b2576f',NULL,'http-basic-authenticator','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','2833329d-525f-48db-bb58-53003ca3bd40',0,10,0,NULL,NULL),('f8c74fd5-3677-46fb-9b7b-46ca93587a7e',NULL,'reset-otp','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','0526685c-2944-462b-8f75-855e0c7761a0',0,20,0,NULL,NULL),('fc09002b-d368-4845-ad60-7aea0d3765cc',NULL,'auth-username-password-form','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','7ec6b431-8ec8-46bb-8611-e0ca04703203',0,10,0,NULL,NULL),('fd3e9650-01b0-49ac-9a0c-2a6ac2910936',NULL,'idp-username-password-form','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','75ba53df-d698-4559-88ca-8d9143e67d75',0,10,0,NULL,NULL),('fe898b70-931f-4ecd-8993-2f29515a1918',NULL,'reset-otp','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','7bca3933-de65-4855-8568-7388b2b12092',0,20,0,NULL,NULL),('ff4f6d56-cc63-40ee-b761-9ca37061268e',NULL,'auth-cookie','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','e23ba648-6f5b-447f-8d45-33da2f4f7e9e',2,10,0,NULL,NULL);
/*!40000 ALTER TABLE `AUTHENTICATION_EXECUTION` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `AUTHENTICATION_FLOW`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `AUTHENTICATION_FLOW` (
  `ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ALIAS` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `DESCRIPTION` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `REALM_ID` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `PROVIDER_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'basic-flow',
  `TOP_LEVEL` tinyint NOT NULL DEFAULT '0',
  `BUILT_IN` tinyint NOT NULL DEFAULT '0',
  PRIMARY KEY (`ID`),
  KEY `IDX_AUTH_FLOW_REALM` (`REALM_ID`),
  CONSTRAINT `FK_AUTH_FLOW_REALM` FOREIGN KEY (`REALM_ID`) REFERENCES `REALM` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `AUTHENTICATION_FLOW` WRITE;
/*!40000 ALTER TABLE `AUTHENTICATION_FLOW` DISABLE KEYS */;
INSERT INTO `AUTHENTICATION_FLOW` VALUES ('0526685c-2944-462b-8f75-855e0c7761a0','Reset - Conditional OTP','Flow to determine if the OTP should be reset or not. Set to REQUIRED to force.','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','basic-flow',0,1),('06cd7bfa-72e7-48fa-9701-d8e7c0a49633','Handle Existing Account','Handle what to do if there is existing account with same email/username like authenticated identity provider','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','basic-flow',0,1),('0f2a187b-6a65-4599-acd3-8404ee73bf62','Account verification options','Method with which to verity the existing account','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','basic-flow',0,1),('216f221b-9be0-40f5-9c68-f5bee839e1bc','docker auth','Used by Docker clients to authenticate against the IDP','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','basic-flow',1,1),('23d3b2c0-2acd-42f9-b25e-3088d6c1caf3','User creation or linking','Flow for the existing/non-existing user alternatives','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','basic-flow',0,1),('26b753d6-6520-4727-b63d-dac95814059f','registration','Registration flow','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','basic-flow',1,1),('2833329d-525f-48db-bb58-53003ca3bd40','saml ecp','SAML ECP Profile Authentication Flow','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','basic-flow',1,1),('2d8f5f5a-cbd8-46b4-a07f-10c9523759c5','Account verification options','Method with which to verity the existing account','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','basic-flow',0,1),('37b360cb-ce8c-4e6a-a7d6-d3380139c06d','Direct Grant - Conditional OTP','Flow to determine if the OTP is required for the authentication','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','basic-flow',0,1),('403e0820-f1af-4eff-8cbc-6a5423453d76','first broker login','Actions taken after first broker login with identity provider account, which is not yet linked to any Keycloak account','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','basic-flow',1,1),('49326f20-3ceb-4c4b-8171-8e7d27d44093','registration','Registration flow','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','basic-flow',1,1),('50b77434-6401-443e-9e81-1edbf8122259','registration form','Registration form','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','form-flow',0,1),('5f3a045d-a67d-43ce-b2e0-9625663e8dd4','Browser - Conditional 2FA','Flow to determine if any 2FA is required for the authentication','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','basic-flow',0,1),('6348a2e8-c07e-4834-9f6f-431b58891ec0','First broker login - Conditional 2FA','Flow to determine if any 2FA is required for the authentication','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','basic-flow',0,1),('676c5a06-ae2d-44a2-89a8-9017809fa2b7','browser','Browser based authentication','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','basic-flow',1,1),('7549030d-3672-4d2b-bbf0-ffd7368e4524','registration form','Registration form','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','form-flow',0,1),('75ba53df-d698-4559-88ca-8d9143e67d75','Verify Existing Account by Re-authentication','Reauthentication of existing account','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','basic-flow',0,1),('79d34d67-cb53-4f81-bbd8-78febfecd274','Organization',NULL,'4c2886d8-e4c4-4b40-9064-445c7ab90a1c','basic-flow',0,1),('7ba91da4-c6bd-47ce-9f9d-da97f7f72703','saml ecp','SAML ECP Profile Authentication Flow','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','basic-flow',1,1),('7bca3933-de65-4855-8568-7388b2b12092','Reset - Conditional OTP','Flow to determine if the OTP should be reset or not. Set to REQUIRED to force.','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','basic-flow',0,1),('7ec6b431-8ec8-46bb-8611-e0ca04703203','forms','Username, password, otp and other auth forms.','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','basic-flow',0,1),('8503d694-8a57-4077-8624-5d9711b52163','User creation or linking','Flow for the existing/non-existing user alternatives','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','basic-flow',0,1),('8862e11d-2025-4230-8693-4e0cfff5d962','First broker login - Conditional 2FA','Flow to determine if any 2FA is required for the authentication','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','basic-flow',0,1),('8a307303-1570-4786-9d86-061a1e51d66f','Handle Existing Account','Handle what to do if there is existing account with same email/username like authenticated identity provider','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','basic-flow',0,1),('906eb7da-9867-4708-9777-30c21db80523','docker auth','Used by Docker clients to authenticate against the IDP','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','basic-flow',1,1),('9b4399ab-2f01-4ccf-a46c-6aa764d61f1a','First Broker Login - Conditional Organization','Flow to determine if the authenticator that adds organization members is to be used','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','basic-flow',0,1),('a0edc4f0-693e-42be-8fd9-f6f6da0515a0','clients','Base authentication for clients','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','client-flow',1,1),('a3c9dcd5-2130-44c7-bab3-55f37e7c2da6','direct grant','OpenID Connect Resource Owner Grant','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','basic-flow',1,1),('b606b87b-37af-4482-8b2d-1c7a18160d44','Direct Grant - Conditional OTP','Flow to determine if the OTP is required for the authentication','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','basic-flow',0,1),('b920157a-a297-4663-afda-7bcfe73e349f','forms','Username, password, otp and other auth forms.','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','basic-flow',0,1),('c08954fb-ccf3-4181-8864-c0a2ec7d0d84','Browser - Conditional Organization','Flow to determine if the organization identity-first login is to be used','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','basic-flow',0,1),('c5735bc4-e5c1-4037-ba8e-85d2640c770d','Verify Existing Account by Re-authentication','Reauthentication of existing account','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','basic-flow',0,1),('d0f21bd5-8d4d-4dee-b58e-77f8cd44ad8c','reset credentials','Reset credentials for a user if they forgot their password or something','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','basic-flow',1,1),('de9091e2-e7a5-42f3-aefb-56c5360be436','clients','Base authentication for clients','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','client-flow',1,1),('e167e17a-2d12-4dcf-9a7f-c9b3c830be01','Browser - Conditional 2FA','Flow to determine if any 2FA is required for the authentication','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','basic-flow',0,1),('e23ba648-6f5b-447f-8d45-33da2f4f7e9e','browser','Browser based authentication','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','basic-flow',1,1),('e9deb89c-1118-49eb-83ba-9a9c1dfb9a2a','reset credentials','Reset credentials for a user if they forgot their password or something','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','basic-flow',1,1),('f3d39771-85ea-427d-b233-fa26308a3952','first broker login','Actions taken after first broker login with identity provider account, which is not yet linked to any Keycloak account','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','basic-flow',1,1),('f662b0fb-ef63-4db6-934e-e4d37ca37296','direct grant','OpenID Connect Resource Owner Grant','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','basic-flow',1,1);
/*!40000 ALTER TABLE `AUTHENTICATION_FLOW` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `AUTHENTICATOR_CONFIG`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `AUTHENTICATOR_CONFIG` (
  `ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ALIAS` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `REALM_ID` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`ID`),
  KEY `IDX_AUTH_CONFIG_REALM` (`REALM_ID`),
  CONSTRAINT `FK_AUTH_REALM` FOREIGN KEY (`REALM_ID`) REFERENCES `REALM` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `AUTHENTICATOR_CONFIG` WRITE;
/*!40000 ALTER TABLE `AUTHENTICATOR_CONFIG` DISABLE KEYS */;
INSERT INTO `AUTHENTICATOR_CONFIG` VALUES ('2be225dc-227b-4f4a-8e19-1956e8d32f94','first-broker-login-conditional-credential','4c2886d8-e4c4-4b40-9064-445c7ab90a1c'),('5b99313d-012b-42fe-a42e-2733a3d43904','browser-conditional-credential','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581'),('64e20831-27ab-4fd0-ab17-aeec0be646a1','first-broker-login-conditional-credential','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581'),('7117ace8-ea89-4a9b-9b13-d0016374af12','create unique user config','4c2886d8-e4c4-4b40-9064-445c7ab90a1c'),('b6ebb592-9ed6-431a-9343-e37b1f0a0b91','review profile config','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581'),('f231ed1d-aaab-4f62-acaf-5e9de2a3b399','create unique user config','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581'),('f7f468af-378f-496f-9ca5-69dc30098555','browser-conditional-credential','4c2886d8-e4c4-4b40-9064-445c7ab90a1c'),('ffc1376a-dddb-4503-aaa6-6008371a716b','review profile config','4c2886d8-e4c4-4b40-9064-445c7ab90a1c');
/*!40000 ALTER TABLE `AUTHENTICATOR_CONFIG` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `AUTHENTICATOR_CONFIG_ENTRY`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `AUTHENTICATOR_CONFIG_ENTRY` (
  `AUTHENTICATOR_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `VALUE` longtext COLLATE utf8mb4_unicode_ci,
  `NAME` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`AUTHENTICATOR_ID`,`NAME`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `AUTHENTICATOR_CONFIG_ENTRY` WRITE;
/*!40000 ALTER TABLE `AUTHENTICATOR_CONFIG_ENTRY` DISABLE KEYS */;
INSERT INTO `AUTHENTICATOR_CONFIG_ENTRY` VALUES ('2be225dc-227b-4f4a-8e19-1956e8d32f94','webauthn-passwordless','credentials'),('5b99313d-012b-42fe-a42e-2733a3d43904','webauthn-passwordless','credentials'),('64e20831-27ab-4fd0-ab17-aeec0be646a1','webauthn-passwordless','credentials'),('7117ace8-ea89-4a9b-9b13-d0016374af12','false','require.password.update.after.registration'),('b6ebb592-9ed6-431a-9343-e37b1f0a0b91','missing','update.profile.on.first.login'),('f231ed1d-aaab-4f62-acaf-5e9de2a3b399','false','require.password.update.after.registration'),('f7f468af-378f-496f-9ca5-69dc30098555','webauthn-passwordless','credentials'),('ffc1376a-dddb-4503-aaa6-6008371a716b','missing','update.profile.on.first.login');
/*!40000 ALTER TABLE `AUTHENTICATOR_CONFIG_ENTRY` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `BROKER_LINK`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `BROKER_LINK` (
  `IDENTITY_PROVIDER` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `STORAGE_PROVIDER_ID` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `REALM_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `BROKER_USER_ID` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `BROKER_USERNAME` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `TOKEN` text COLLATE utf8mb4_unicode_ci,
  `USER_ID` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`IDENTITY_PROVIDER`,`USER_ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `BROKER_LINK` WRITE;
/*!40000 ALTER TABLE `BROKER_LINK` DISABLE KEYS */;
/*!40000 ALTER TABLE `BROKER_LINK` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `CLIENT`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `CLIENT` (
  `ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ENABLED` tinyint NOT NULL DEFAULT '0',
  `FULL_SCOPE_ALLOWED` tinyint NOT NULL DEFAULT '0',
  `CLIENT_ID` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `NOT_BEFORE` int DEFAULT NULL,
  `PUBLIC_CLIENT` tinyint NOT NULL DEFAULT '0',
  `SECRET` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `BASE_URL` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `BEARER_ONLY` tinyint NOT NULL DEFAULT '0',
  `MANAGEMENT_URL` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `SURROGATE_AUTH_REQUIRED` tinyint NOT NULL DEFAULT '0',
  `REALM_ID` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `PROTOCOL` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `NODE_REREG_TIMEOUT` int DEFAULT '0',
  `FRONTCHANNEL_LOGOUT` tinyint NOT NULL DEFAULT '0',
  `CONSENT_REQUIRED` tinyint NOT NULL DEFAULT '0',
  `NAME` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `SERVICE_ACCOUNTS_ENABLED` tinyint NOT NULL DEFAULT '0',
  `CLIENT_AUTHENTICATOR_TYPE` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `ROOT_URL` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `DESCRIPTION` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `REGISTRATION_TOKEN` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `STANDARD_FLOW_ENABLED` tinyint NOT NULL DEFAULT '1',
  `IMPLICIT_FLOW_ENABLED` tinyint NOT NULL DEFAULT '0',
  `DIRECT_ACCESS_GRANTS_ENABLED` tinyint NOT NULL DEFAULT '0',
  `ALWAYS_DISPLAY_IN_CONSOLE` tinyint NOT NULL DEFAULT '0',
  PRIMARY KEY (`ID`),
  UNIQUE KEY `UK_B71CJLBENV945RB6GCON438AT` (`REALM_ID`,`CLIENT_ID`),
  KEY `IDX_CLIENT_ID` (`CLIENT_ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `CLIENT` WRITE;
/*!40000 ALTER TABLE `CLIENT` DISABLE KEYS */;
INSERT INTO `CLIENT` VALUES ('1085b6f8-5701-4591-8dd5-7b68238b51b8',1,0,'broker',0,0,NULL,NULL,1,NULL,0,'5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','openid-connect',0,0,0,'${client_broker}',0,'client-secret',NULL,NULL,NULL,1,0,0,0),('13fd5765-ac20-4558-ae2c-cf8e2df4c22c',1,0,'account-console',0,1,NULL,'/realms/master/account/',0,NULL,0,'5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','openid-connect',0,0,0,'${client_account-console}',0,'client-secret','${authBaseUrl}',NULL,NULL,1,0,0,0),('232e71d1-663b-4fd9-981f-0f08b56cd567',1,0,'motogo-realm',0,0,NULL,NULL,1,NULL,0,'5d56fd15-57bc-496f-8f3f-bf0d7a5c5581',NULL,0,0,0,'motogo Realm',0,'client-secret',NULL,NULL,NULL,1,0,0,0),('286f0039-9d0d-4f17-8b87-9cbeddaa0a89',1,0,'account',0,1,NULL,'/realms/master/account/',0,NULL,0,'5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','openid-connect',0,0,0,'${client_account}',0,'client-secret','${authBaseUrl}',NULL,NULL,1,0,0,0),('35ba46a4-b1e1-4dd0-8cff-fc98df4c2659',1,1,'stifler',0,0,'sO7Jk9qgViBGFaO9GGpAR6ggz1Oqtuaf','',0,'',0,'4c2886d8-e4c4-4b40-9064-445c7ab90a1c','openid-connect',-1,1,0,'Esteban',0,'client-secret','','admin',NULL,1,0,1,1),('3c6f327b-8018-44e9-8540-511170baeb2f',1,1,'security-admin-console',0,1,NULL,'/admin/motogo/console/',0,NULL,0,'4c2886d8-e4c4-4b40-9064-445c7ab90a1c','openid-connect',0,0,0,'${client_security-admin-console}',0,'client-secret','${authAdminUrl}',NULL,NULL,1,0,0,0),('4c9063cc-607c-4f14-b55d-eb3d700e742c',1,1,'security-admin-console',0,1,NULL,'/admin/master/console/',0,NULL,0,'5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','openid-connect',0,0,0,'${client_security-admin-console}',0,'client-secret','${authAdminUrl}',NULL,NULL,1,0,0,0),('5fc55bdb-f575-4b2e-a17c-e39bc6fb9412',1,0,'broker',0,0,NULL,NULL,1,NULL,0,'4c2886d8-e4c4-4b40-9064-445c7ab90a1c','openid-connect',0,0,0,'${client_broker}',0,'client-secret',NULL,NULL,NULL,1,0,0,0),('6b5c5dea-9656-4dfb-8a26-f699429eecd3',1,0,'master-realm',0,0,NULL,NULL,1,NULL,0,'5d56fd15-57bc-496f-8f3f-bf0d7a5c5581',NULL,0,0,0,'master Realm',0,'client-secret',NULL,NULL,NULL,1,0,0,0),('8b0a29f2-46dd-4f3c-8446-c0f63ad56ebc',1,0,'account',0,1,NULL,'/realms/motogo/account/',0,NULL,0,'4c2886d8-e4c4-4b40-9064-445c7ab90a1c','openid-connect',0,0,0,'${client_account}',0,'client-secret','${authBaseUrl}',NULL,NULL,1,0,0,0),('8b786b3a-c968-49cf-b5eb-5e7b36889df7',1,0,'realm-management',0,0,NULL,NULL,1,NULL,0,'4c2886d8-e4c4-4b40-9064-445c7ab90a1c','openid-connect',0,0,0,'${client_realm-management}',0,'client-secret',NULL,NULL,NULL,1,0,0,0),('ba9af61d-e6d3-41c9-88d0-ab69b58e061b',1,1,'admin-cli',0,1,NULL,NULL,0,NULL,0,'5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','openid-connect',0,0,0,'${client_admin-cli}',0,'client-secret',NULL,NULL,NULL,0,0,1,0),('d2eb4296-2046-420d-beb2-25cf9e05b715',1,0,'account-console',0,1,NULL,'/realms/motogo/account/',0,NULL,0,'4c2886d8-e4c4-4b40-9064-445c7ab90a1c','openid-connect',0,0,0,'${client_account-console}',0,'client-secret','${authBaseUrl}',NULL,NULL,1,0,0,0),('fe3360eb-3c93-4e56-99df-b8fd05dbf60a',1,1,'admin-cli',0,1,NULL,NULL,0,NULL,0,'4c2886d8-e4c4-4b40-9064-445c7ab90a1c','openid-connect',0,0,0,'${client_admin-cli}',0,'client-secret',NULL,NULL,NULL,0,0,1,0);
/*!40000 ALTER TABLE `CLIENT` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `CLIENT_ATTRIBUTES`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `CLIENT_ATTRIBUTES` (
  `CLIENT_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `NAME` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `VALUE` longtext CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  PRIMARY KEY (`CLIENT_ID`,`NAME`),
  KEY `IDX_CLIENT_ATT_BY_NAME_VALUE` (`NAME`,`VALUE`(255)),
  CONSTRAINT `FK3C47C64BEACCA966` FOREIGN KEY (`CLIENT_ID`) REFERENCES `CLIENT` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `CLIENT_ATTRIBUTES` WRITE;
/*!40000 ALTER TABLE `CLIENT_ATTRIBUTES` DISABLE KEYS */;
INSERT INTO `CLIENT_ATTRIBUTES` VALUES ('13fd5765-ac20-4558-ae2c-cf8e2df4c22c','pkce.code.challenge.method','S256'),('13fd5765-ac20-4558-ae2c-cf8e2df4c22c','post.logout.redirect.uris','+'),('286f0039-9d0d-4f17-8b87-9cbeddaa0a89','post.logout.redirect.uris','+'),('35ba46a4-b1e1-4dd0-8cff-fc98df4c2659','access.token.header.type.rfc9068','false'),('35ba46a4-b1e1-4dd0-8cff-fc98df4c2659','acr.loa.map','{}'),('35ba46a4-b1e1-4dd0-8cff-fc98df4c2659','backchannel.logout.revoke.offline.tokens','false'),('35ba46a4-b1e1-4dd0-8cff-fc98df4c2659','backchannel.logout.session.required','true'),('35ba46a4-b1e1-4dd0-8cff-fc98df4c2659','client_credentials.use_refresh_token','false'),('35ba46a4-b1e1-4dd0-8cff-fc98df4c2659','client.introspection.response.allow.jwt.claim.enabled','false'),('35ba46a4-b1e1-4dd0-8cff-fc98df4c2659','client.secret.creation.time','1766007023'),('35ba46a4-b1e1-4dd0-8cff-fc98df4c2659','client.use.lightweight.access.token.enabled','false'),('35ba46a4-b1e1-4dd0-8cff-fc98df4c2659','display.on.consent.screen','false'),('35ba46a4-b1e1-4dd0-8cff-fc98df4c2659','dpop.bound.access.tokens','false'),('35ba46a4-b1e1-4dd0-8cff-fc98df4c2659','frontchannel.logout.session.required','true'),('35ba46a4-b1e1-4dd0-8cff-fc98df4c2659','id.token.as.detached.signature','false'),('35ba46a4-b1e1-4dd0-8cff-fc98df4c2659','oauth2.device.authorization.grant.enabled','false'),('35ba46a4-b1e1-4dd0-8cff-fc98df4c2659','oidc.ciba.grant.enabled','false'),('35ba46a4-b1e1-4dd0-8cff-fc98df4c2659','realm_client','false'),('35ba46a4-b1e1-4dd0-8cff-fc98df4c2659','request.object.encryption.alg','any'),('35ba46a4-b1e1-4dd0-8cff-fc98df4c2659','request.object.encryption.enc','any'),('35ba46a4-b1e1-4dd0-8cff-fc98df4c2659','request.object.required','not required'),('35ba46a4-b1e1-4dd0-8cff-fc98df4c2659','request.object.signature.alg','any'),('35ba46a4-b1e1-4dd0-8cff-fc98df4c2659','require.pushed.authorization.requests','false'),('35ba46a4-b1e1-4dd0-8cff-fc98df4c2659','standard.token.exchange.enabled','false'),('35ba46a4-b1e1-4dd0-8cff-fc98df4c2659','tls.client.certificate.bound.access.tokens','false'),('35ba46a4-b1e1-4dd0-8cff-fc98df4c2659','token.response.type.bearer.lower-case','false'),('35ba46a4-b1e1-4dd0-8cff-fc98df4c2659','use.jwks.url','false'),('35ba46a4-b1e1-4dd0-8cff-fc98df4c2659','use.refresh.tokens','true'),('3c6f327b-8018-44e9-8540-511170baeb2f','client.use.lightweight.access.token.enabled','true'),('3c6f327b-8018-44e9-8540-511170baeb2f','pkce.code.challenge.method','S256'),('3c6f327b-8018-44e9-8540-511170baeb2f','post.logout.redirect.uris','+'),('4c9063cc-607c-4f14-b55d-eb3d700e742c','client.use.lightweight.access.token.enabled','true'),('4c9063cc-607c-4f14-b55d-eb3d700e742c','pkce.code.challenge.method','S256'),('4c9063cc-607c-4f14-b55d-eb3d700e742c','post.logout.redirect.uris','+'),('8b0a29f2-46dd-4f3c-8446-c0f63ad56ebc','post.logout.redirect.uris','+'),('ba9af61d-e6d3-41c9-88d0-ab69b58e061b','client.use.lightweight.access.token.enabled','true'),('d2eb4296-2046-420d-beb2-25cf9e05b715','pkce.code.challenge.method','S256'),('d2eb4296-2046-420d-beb2-25cf9e05b715','post.logout.redirect.uris','+'),('fe3360eb-3c93-4e56-99df-b8fd05dbf60a','client.use.lightweight.access.token.enabled','true');
/*!40000 ALTER TABLE `CLIENT_ATTRIBUTES` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `CLIENT_AUTH_FLOW_BINDINGS`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `CLIENT_AUTH_FLOW_BINDINGS` (
  `CLIENT_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `FLOW_ID` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `BINDING_NAME` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`CLIENT_ID`,`BINDING_NAME`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `CLIENT_AUTH_FLOW_BINDINGS` WRITE;
/*!40000 ALTER TABLE `CLIENT_AUTH_FLOW_BINDINGS` DISABLE KEYS */;
/*!40000 ALTER TABLE `CLIENT_AUTH_FLOW_BINDINGS` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `CLIENT_INITIAL_ACCESS`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `CLIENT_INITIAL_ACCESS` (
  `ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `REALM_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `TIMESTAMP` int DEFAULT NULL,
  `EXPIRATION` int DEFAULT NULL,
  `COUNT` int DEFAULT NULL,
  `REMAINING_COUNT` int DEFAULT NULL,
  PRIMARY KEY (`ID`),
  KEY `IDX_CLIENT_INIT_ACC_REALM` (`REALM_ID`),
  CONSTRAINT `FK_CLIENT_INIT_ACC_REALM` FOREIGN KEY (`REALM_ID`) REFERENCES `REALM` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `CLIENT_INITIAL_ACCESS` WRITE;
/*!40000 ALTER TABLE `CLIENT_INITIAL_ACCESS` DISABLE KEYS */;
/*!40000 ALTER TABLE `CLIENT_INITIAL_ACCESS` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `CLIENT_NODE_REGISTRATIONS`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `CLIENT_NODE_REGISTRATIONS` (
  `CLIENT_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `VALUE` int DEFAULT NULL,
  `NAME` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`CLIENT_ID`,`NAME`),
  CONSTRAINT `FK4129723BA992F594` FOREIGN KEY (`CLIENT_ID`) REFERENCES `CLIENT` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `CLIENT_NODE_REGISTRATIONS` WRITE;
/*!40000 ALTER TABLE `CLIENT_NODE_REGISTRATIONS` DISABLE KEYS */;
/*!40000 ALTER TABLE `CLIENT_NODE_REGISTRATIONS` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `CLIENT_SCOPE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `CLIENT_SCOPE` (
  `ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `NAME` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `REALM_ID` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `DESCRIPTION` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `PROTOCOL` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `UK_CLI_SCOPE` (`REALM_ID`,`NAME`),
  KEY `IDX_REALM_CLSCOPE` (`REALM_ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `CLIENT_SCOPE` WRITE;
/*!40000 ALTER TABLE `CLIENT_SCOPE` DISABLE KEYS */;
INSERT INTO `CLIENT_SCOPE` VALUES ('09dab565-4627-474a-ba4a-2c1dd91593a6','organization','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','Additional claims about the organization a subject belongs to','openid-connect'),('1daedf09-7b9a-4ed1-b0a7-db5dec46c9b5','service_account','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','Specific scope for a client enabled for service accounts','openid-connect'),('23ac037f-c944-45b5-93f5-bd8395bf49e0','email','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','OpenID Connect built-in scope: email','openid-connect'),('35746b28-fed9-4ae4-92a2-9b42f120f272','email','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','OpenID Connect built-in scope: email','openid-connect'),('3ab1f10c-ca69-4372-8fb3-dfb2bc434ca4','profile','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','OpenID Connect built-in scope: profile','openid-connect'),('4178816c-d39e-4c12-a66e-099a02427fa7','basic','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','OpenID Connect scope for add all basic claims to the token','openid-connect'),('44b8cc1f-4adc-4465-86ad-c7ba53036921','roles','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','OpenID Connect scope for add user roles to the access token','openid-connect'),('4b7b21fc-bddf-42ce-97f0-36242b13dcfd','organization','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','Additional claims about the organization a subject belongs to','openid-connect'),('5dcfd9e0-f7f3-42be-8cc7-fd9079e16baf','role_list','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','SAML role list','saml'),('70413ad2-e95a-413c-af6e-d1741fe9dfe6','basic','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','OpenID Connect scope for add all basic claims to the token','openid-connect'),('7b89f755-1355-4936-a90e-04777d87b75b','microprofile-jwt','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','Microprofile - JWT built-in scope','openid-connect'),('8a6cb980-3684-4ef9-9d40-88f23541f569','phone','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','OpenID Connect built-in scope: phone','openid-connect'),('9f17cc3f-a321-4d62-aca0-b8a889d8629d','offline_access','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','OpenID Connect built-in scope: offline_access','openid-connect'),('9f9d3e99-a0c3-4e59-96b4-d9177299709b','phone','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','OpenID Connect built-in scope: phone','openid-connect'),('a687fb3e-0aca-481d-a697-9c23538f74c2','roles','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','OpenID Connect scope for add user roles to the access token','openid-connect'),('a84876a0-9e4b-4a9f-9369-a49e29c2cdac','web-origins','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','OpenID Connect scope for add allowed web origins to the access token','openid-connect'),('b679321c-7425-4c23-8834-e0209a8ae8c5','profile','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','OpenID Connect built-in scope: profile','openid-connect'),('b82ffaf7-8ece-41c1-a3f9-00bfbe9e9739','saml_organization','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','Organization Membership','saml'),('be17e97a-a116-447b-bfb3-aacefeb3fc4f','web-origins','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','OpenID Connect scope for add allowed web origins to the access token','openid-connect'),('bea2ab00-c561-421a-8d40-5a0dca74e80b','acr','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','OpenID Connect scope for add acr (authentication context class reference) to the token','openid-connect'),('c80eb4e6-b9ce-4164-82e3-6e00ae132571','saml_organization','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','Organization Membership','saml'),('c83dfc09-88fa-4c04-bf69-4d9afd7d6d4e','acr','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','OpenID Connect scope for add acr (authentication context class reference) to the token','openid-connect'),('cb48b63f-5f46-4cef-a1bf-a1ddc25d2b24','role_list','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','SAML role list','saml'),('db7d6d32-c0a5-4b47-abf8-31ba61143504','address','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','OpenID Connect built-in scope: address','openid-connect'),('dbf26c31-238a-4e9d-8336-3daf3d2dd7bd','offline_access','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','OpenID Connect built-in scope: offline_access','openid-connect'),('df0e3ca7-158b-4a1c-b58d-9b423067475e','microprofile-jwt','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','Microprofile - JWT built-in scope','openid-connect'),('f06d6bc3-3c00-4cf3-a474-e5becd828d7a','service_account','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','Specific scope for a client enabled for service accounts','openid-connect'),('f8102d42-5d53-4610-bf3e-aacec43635e4','address','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','OpenID Connect built-in scope: address','openid-connect');
/*!40000 ALTER TABLE `CLIENT_SCOPE` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `CLIENT_SCOPE_ATTRIBUTES`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `CLIENT_SCOPE_ATTRIBUTES` (
  `SCOPE_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `VALUE` text COLLATE utf8mb4_unicode_ci,
  `NAME` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`SCOPE_ID`,`NAME`),
  KEY `IDX_CLSCOPE_ATTRS` (`SCOPE_ID`),
  CONSTRAINT `FK_CL_SCOPE_ATTR_SCOPE` FOREIGN KEY (`SCOPE_ID`) REFERENCES `CLIENT_SCOPE` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `CLIENT_SCOPE_ATTRIBUTES` WRITE;
/*!40000 ALTER TABLE `CLIENT_SCOPE_ATTRIBUTES` DISABLE KEYS */;
INSERT INTO `CLIENT_SCOPE_ATTRIBUTES` VALUES ('09dab565-4627-474a-ba4a-2c1dd91593a6','${organizationScopeConsentText}','consent.screen.text'),('09dab565-4627-474a-ba4a-2c1dd91593a6','true','display.on.consent.screen'),('09dab565-4627-474a-ba4a-2c1dd91593a6','true','include.in.token.scope'),('1daedf09-7b9a-4ed1-b0a7-db5dec46c9b5','false','display.on.consent.screen'),('1daedf09-7b9a-4ed1-b0a7-db5dec46c9b5','false','include.in.token.scope'),('23ac037f-c944-45b5-93f5-bd8395bf49e0','${emailScopeConsentText}','consent.screen.text'),('23ac037f-c944-45b5-93f5-bd8395bf49e0','true','display.on.consent.screen'),('23ac037f-c944-45b5-93f5-bd8395bf49e0','true','include.in.token.scope'),('35746b28-fed9-4ae4-92a2-9b42f120f272','${emailScopeConsentText}','consent.screen.text'),('35746b28-fed9-4ae4-92a2-9b42f120f272','true','display.on.consent.screen'),('35746b28-fed9-4ae4-92a2-9b42f120f272','true','include.in.token.scope'),('3ab1f10c-ca69-4372-8fb3-dfb2bc434ca4','${profileScopeConsentText}','consent.screen.text'),('3ab1f10c-ca69-4372-8fb3-dfb2bc434ca4','true','display.on.consent.screen'),('3ab1f10c-ca69-4372-8fb3-dfb2bc434ca4','true','include.in.token.scope'),('4178816c-d39e-4c12-a66e-099a02427fa7','false','display.on.consent.screen'),('4178816c-d39e-4c12-a66e-099a02427fa7','false','include.in.token.scope'),('44b8cc1f-4adc-4465-86ad-c7ba53036921','${rolesScopeConsentText}','consent.screen.text'),('44b8cc1f-4adc-4465-86ad-c7ba53036921','true','display.on.consent.screen'),('44b8cc1f-4adc-4465-86ad-c7ba53036921','false','include.in.token.scope'),('4b7b21fc-bddf-42ce-97f0-36242b13dcfd','${organizationScopeConsentText}','consent.screen.text'),('4b7b21fc-bddf-42ce-97f0-36242b13dcfd','true','display.on.consent.screen'),('4b7b21fc-bddf-42ce-97f0-36242b13dcfd','true','include.in.token.scope'),('5dcfd9e0-f7f3-42be-8cc7-fd9079e16baf','${samlRoleListScopeConsentText}','consent.screen.text'),('5dcfd9e0-f7f3-42be-8cc7-fd9079e16baf','true','display.on.consent.screen'),('70413ad2-e95a-413c-af6e-d1741fe9dfe6','false','display.on.consent.screen'),('70413ad2-e95a-413c-af6e-d1741fe9dfe6','false','include.in.token.scope'),('7b89f755-1355-4936-a90e-04777d87b75b','false','display.on.consent.screen'),('7b89f755-1355-4936-a90e-04777d87b75b','true','include.in.token.scope'),('8a6cb980-3684-4ef9-9d40-88f23541f569','${phoneScopeConsentText}','consent.screen.text'),('8a6cb980-3684-4ef9-9d40-88f23541f569','true','display.on.consent.screen'),('8a6cb980-3684-4ef9-9d40-88f23541f569','true','include.in.token.scope'),('9f17cc3f-a321-4d62-aca0-b8a889d8629d','${offlineAccessScopeConsentText}','consent.screen.text'),('9f17cc3f-a321-4d62-aca0-b8a889d8629d','true','display.on.consent.screen'),('9f9d3e99-a0c3-4e59-96b4-d9177299709b','${phoneScopeConsentText}','consent.screen.text'),('9f9d3e99-a0c3-4e59-96b4-d9177299709b','true','display.on.consent.screen'),('9f9d3e99-a0c3-4e59-96b4-d9177299709b','true','include.in.token.scope'),('a687fb3e-0aca-481d-a697-9c23538f74c2','${rolesScopeConsentText}','consent.screen.text'),('a687fb3e-0aca-481d-a697-9c23538f74c2','true','display.on.consent.screen'),('a687fb3e-0aca-481d-a697-9c23538f74c2','false','include.in.token.scope'),('a84876a0-9e4b-4a9f-9369-a49e29c2cdac','','consent.screen.text'),('a84876a0-9e4b-4a9f-9369-a49e29c2cdac','false','display.on.consent.screen'),('a84876a0-9e4b-4a9f-9369-a49e29c2cdac','false','include.in.token.scope'),('b679321c-7425-4c23-8834-e0209a8ae8c5','${profileScopeConsentText}','consent.screen.text'),('b679321c-7425-4c23-8834-e0209a8ae8c5','true','display.on.consent.screen'),('b679321c-7425-4c23-8834-e0209a8ae8c5','true','include.in.token.scope'),('b82ffaf7-8ece-41c1-a3f9-00bfbe9e9739','false','display.on.consent.screen'),('be17e97a-a116-447b-bfb3-aacefeb3fc4f','','consent.screen.text'),('be17e97a-a116-447b-bfb3-aacefeb3fc4f','false','display.on.consent.screen'),('be17e97a-a116-447b-bfb3-aacefeb3fc4f','false','include.in.token.scope'),('bea2ab00-c561-421a-8d40-5a0dca74e80b','false','display.on.consent.screen'),('bea2ab00-c561-421a-8d40-5a0dca74e80b','false','include.in.token.scope'),('c80eb4e6-b9ce-4164-82e3-6e00ae132571','false','display.on.consent.screen'),('c83dfc09-88fa-4c04-bf69-4d9afd7d6d4e','false','display.on.consent.screen'),('c83dfc09-88fa-4c04-bf69-4d9afd7d6d4e','false','include.in.token.scope'),('cb48b63f-5f46-4cef-a1bf-a1ddc25d2b24','${samlRoleListScopeConsentText}','consent.screen.text'),('cb48b63f-5f46-4cef-a1bf-a1ddc25d2b24','true','display.on.consent.screen'),('db7d6d32-c0a5-4b47-abf8-31ba61143504','${addressScopeConsentText}','consent.screen.text'),('db7d6d32-c0a5-4b47-abf8-31ba61143504','true','display.on.consent.screen'),('db7d6d32-c0a5-4b47-abf8-31ba61143504','true','include.in.token.scope'),('dbf26c31-238a-4e9d-8336-3daf3d2dd7bd','${offlineAccessScopeConsentText}','consent.screen.text'),('dbf26c31-238a-4e9d-8336-3daf3d2dd7bd','true','display.on.consent.screen'),('df0e3ca7-158b-4a1c-b58d-9b423067475e','false','display.on.consent.screen'),('df0e3ca7-158b-4a1c-b58d-9b423067475e','true','include.in.token.scope'),('f06d6bc3-3c00-4cf3-a474-e5becd828d7a','false','display.on.consent.screen'),('f06d6bc3-3c00-4cf3-a474-e5becd828d7a','false','include.in.token.scope'),('f8102d42-5d53-4610-bf3e-aacec43635e4','${addressScopeConsentText}','consent.screen.text'),('f8102d42-5d53-4610-bf3e-aacec43635e4','true','display.on.consent.screen'),('f8102d42-5d53-4610-bf3e-aacec43635e4','true','include.in.token.scope');
/*!40000 ALTER TABLE `CLIENT_SCOPE_ATTRIBUTES` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `CLIENT_SCOPE_CLIENT`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `CLIENT_SCOPE_CLIENT` (
  `CLIENT_ID` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `SCOPE_ID` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `DEFAULT_SCOPE` tinyint NOT NULL DEFAULT '0',
  PRIMARY KEY (`CLIENT_ID`,`SCOPE_ID`),
  KEY `IDX_CLSCOPE_CL` (`CLIENT_ID`),
  KEY `IDX_CL_CLSCOPE` (`SCOPE_ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `CLIENT_SCOPE_CLIENT` WRITE;
/*!40000 ALTER TABLE `CLIENT_SCOPE_CLIENT` DISABLE KEYS */;
INSERT INTO `CLIENT_SCOPE_CLIENT` VALUES ('1085b6f8-5701-4591-8dd5-7b68238b51b8','09dab565-4627-474a-ba4a-2c1dd91593a6',0),('1085b6f8-5701-4591-8dd5-7b68238b51b8','23ac037f-c944-45b5-93f5-bd8395bf49e0',1),('1085b6f8-5701-4591-8dd5-7b68238b51b8','4178816c-d39e-4c12-a66e-099a02427fa7',1),('1085b6f8-5701-4591-8dd5-7b68238b51b8','44b8cc1f-4adc-4465-86ad-c7ba53036921',1),('1085b6f8-5701-4591-8dd5-7b68238b51b8','7b89f755-1355-4936-a90e-04777d87b75b',0),('1085b6f8-5701-4591-8dd5-7b68238b51b8','9f17cc3f-a321-4d62-aca0-b8a889d8629d',0),('1085b6f8-5701-4591-8dd5-7b68238b51b8','9f9d3e99-a0c3-4e59-96b4-d9177299709b',0),('1085b6f8-5701-4591-8dd5-7b68238b51b8','b679321c-7425-4c23-8834-e0209a8ae8c5',1),('1085b6f8-5701-4591-8dd5-7b68238b51b8','be17e97a-a116-447b-bfb3-aacefeb3fc4f',1),('1085b6f8-5701-4591-8dd5-7b68238b51b8','c83dfc09-88fa-4c04-bf69-4d9afd7d6d4e',1),('1085b6f8-5701-4591-8dd5-7b68238b51b8','db7d6d32-c0a5-4b47-abf8-31ba61143504',0),('13fd5765-ac20-4558-ae2c-cf8e2df4c22c','09dab565-4627-474a-ba4a-2c1dd91593a6',0),('13fd5765-ac20-4558-ae2c-cf8e2df4c22c','23ac037f-c944-45b5-93f5-bd8395bf49e0',1),('13fd5765-ac20-4558-ae2c-cf8e2df4c22c','4178816c-d39e-4c12-a66e-099a02427fa7',1),('13fd5765-ac20-4558-ae2c-cf8e2df4c22c','44b8cc1f-4adc-4465-86ad-c7ba53036921',1),('13fd5765-ac20-4558-ae2c-cf8e2df4c22c','7b89f755-1355-4936-a90e-04777d87b75b',0),('13fd5765-ac20-4558-ae2c-cf8e2df4c22c','9f17cc3f-a321-4d62-aca0-b8a889d8629d',0),('13fd5765-ac20-4558-ae2c-cf8e2df4c22c','9f9d3e99-a0c3-4e59-96b4-d9177299709b',0),('13fd5765-ac20-4558-ae2c-cf8e2df4c22c','b679321c-7425-4c23-8834-e0209a8ae8c5',1),('13fd5765-ac20-4558-ae2c-cf8e2df4c22c','be17e97a-a116-447b-bfb3-aacefeb3fc4f',1),('13fd5765-ac20-4558-ae2c-cf8e2df4c22c','c83dfc09-88fa-4c04-bf69-4d9afd7d6d4e',1),('13fd5765-ac20-4558-ae2c-cf8e2df4c22c','db7d6d32-c0a5-4b47-abf8-31ba61143504',0),('286f0039-9d0d-4f17-8b87-9cbeddaa0a89','09dab565-4627-474a-ba4a-2c1dd91593a6',0),('286f0039-9d0d-4f17-8b87-9cbeddaa0a89','23ac037f-c944-45b5-93f5-bd8395bf49e0',1),('286f0039-9d0d-4f17-8b87-9cbeddaa0a89','4178816c-d39e-4c12-a66e-099a02427fa7',1),('286f0039-9d0d-4f17-8b87-9cbeddaa0a89','44b8cc1f-4adc-4465-86ad-c7ba53036921',1),('286f0039-9d0d-4f17-8b87-9cbeddaa0a89','7b89f755-1355-4936-a90e-04777d87b75b',0),('286f0039-9d0d-4f17-8b87-9cbeddaa0a89','9f17cc3f-a321-4d62-aca0-b8a889d8629d',0),('286f0039-9d0d-4f17-8b87-9cbeddaa0a89','9f9d3e99-a0c3-4e59-96b4-d9177299709b',0),('286f0039-9d0d-4f17-8b87-9cbeddaa0a89','b679321c-7425-4c23-8834-e0209a8ae8c5',1),('286f0039-9d0d-4f17-8b87-9cbeddaa0a89','be17e97a-a116-447b-bfb3-aacefeb3fc4f',1),('286f0039-9d0d-4f17-8b87-9cbeddaa0a89','c83dfc09-88fa-4c04-bf69-4d9afd7d6d4e',1),('286f0039-9d0d-4f17-8b87-9cbeddaa0a89','db7d6d32-c0a5-4b47-abf8-31ba61143504',0),('35ba46a4-b1e1-4dd0-8cff-fc98df4c2659','35746b28-fed9-4ae4-92a2-9b42f120f272',1),('35ba46a4-b1e1-4dd0-8cff-fc98df4c2659','3ab1f10c-ca69-4372-8fb3-dfb2bc434ca4',1),('35ba46a4-b1e1-4dd0-8cff-fc98df4c2659','4b7b21fc-bddf-42ce-97f0-36242b13dcfd',0),('35ba46a4-b1e1-4dd0-8cff-fc98df4c2659','70413ad2-e95a-413c-af6e-d1741fe9dfe6',1),('35ba46a4-b1e1-4dd0-8cff-fc98df4c2659','8a6cb980-3684-4ef9-9d40-88f23541f569',0),('35ba46a4-b1e1-4dd0-8cff-fc98df4c2659','a687fb3e-0aca-481d-a697-9c23538f74c2',1),('35ba46a4-b1e1-4dd0-8cff-fc98df4c2659','a84876a0-9e4b-4a9f-9369-a49e29c2cdac',1),('35ba46a4-b1e1-4dd0-8cff-fc98df4c2659','bea2ab00-c561-421a-8d40-5a0dca74e80b',1),('35ba46a4-b1e1-4dd0-8cff-fc98df4c2659','dbf26c31-238a-4e9d-8336-3daf3d2dd7bd',0),('35ba46a4-b1e1-4dd0-8cff-fc98df4c2659','df0e3ca7-158b-4a1c-b58d-9b423067475e',0),('35ba46a4-b1e1-4dd0-8cff-fc98df4c2659','f8102d42-5d53-4610-bf3e-aacec43635e4',0),('3c6f327b-8018-44e9-8540-511170baeb2f','35746b28-fed9-4ae4-92a2-9b42f120f272',1),('3c6f327b-8018-44e9-8540-511170baeb2f','3ab1f10c-ca69-4372-8fb3-dfb2bc434ca4',1),('3c6f327b-8018-44e9-8540-511170baeb2f','4b7b21fc-bddf-42ce-97f0-36242b13dcfd',0),('3c6f327b-8018-44e9-8540-511170baeb2f','70413ad2-e95a-413c-af6e-d1741fe9dfe6',1),('3c6f327b-8018-44e9-8540-511170baeb2f','8a6cb980-3684-4ef9-9d40-88f23541f569',0),('3c6f327b-8018-44e9-8540-511170baeb2f','a687fb3e-0aca-481d-a697-9c23538f74c2',1),('3c6f327b-8018-44e9-8540-511170baeb2f','a84876a0-9e4b-4a9f-9369-a49e29c2cdac',1),('3c6f327b-8018-44e9-8540-511170baeb2f','bea2ab00-c561-421a-8d40-5a0dca74e80b',1),('3c6f327b-8018-44e9-8540-511170baeb2f','dbf26c31-238a-4e9d-8336-3daf3d2dd7bd',0),('3c6f327b-8018-44e9-8540-511170baeb2f','df0e3ca7-158b-4a1c-b58d-9b423067475e',0),('3c6f327b-8018-44e9-8540-511170baeb2f','f8102d42-5d53-4610-bf3e-aacec43635e4',0),('4c9063cc-607c-4f14-b55d-eb3d700e742c','09dab565-4627-474a-ba4a-2c1dd91593a6',0),('4c9063cc-607c-4f14-b55d-eb3d700e742c','23ac037f-c944-45b5-93f5-bd8395bf49e0',1),('4c9063cc-607c-4f14-b55d-eb3d700e742c','4178816c-d39e-4c12-a66e-099a02427fa7',1),('4c9063cc-607c-4f14-b55d-eb3d700e742c','44b8cc1f-4adc-4465-86ad-c7ba53036921',1),('4c9063cc-607c-4f14-b55d-eb3d700e742c','7b89f755-1355-4936-a90e-04777d87b75b',0),('4c9063cc-607c-4f14-b55d-eb3d700e742c','9f17cc3f-a321-4d62-aca0-b8a889d8629d',0),('4c9063cc-607c-4f14-b55d-eb3d700e742c','9f9d3e99-a0c3-4e59-96b4-d9177299709b',0),('4c9063cc-607c-4f14-b55d-eb3d700e742c','b679321c-7425-4c23-8834-e0209a8ae8c5',1),('4c9063cc-607c-4f14-b55d-eb3d700e742c','be17e97a-a116-447b-bfb3-aacefeb3fc4f',1),('4c9063cc-607c-4f14-b55d-eb3d700e742c','c83dfc09-88fa-4c04-bf69-4d9afd7d6d4e',1),('4c9063cc-607c-4f14-b55d-eb3d700e742c','db7d6d32-c0a5-4b47-abf8-31ba61143504',0),('5fc55bdb-f575-4b2e-a17c-e39bc6fb9412','35746b28-fed9-4ae4-92a2-9b42f120f272',1),('5fc55bdb-f575-4b2e-a17c-e39bc6fb9412','3ab1f10c-ca69-4372-8fb3-dfb2bc434ca4',1),('5fc55bdb-f575-4b2e-a17c-e39bc6fb9412','4b7b21fc-bddf-42ce-97f0-36242b13dcfd',0),('5fc55bdb-f575-4b2e-a17c-e39bc6fb9412','70413ad2-e95a-413c-af6e-d1741fe9dfe6',1),('5fc55bdb-f575-4b2e-a17c-e39bc6fb9412','8a6cb980-3684-4ef9-9d40-88f23541f569',0),('5fc55bdb-f575-4b2e-a17c-e39bc6fb9412','a687fb3e-0aca-481d-a697-9c23538f74c2',1),('5fc55bdb-f575-4b2e-a17c-e39bc6fb9412','a84876a0-9e4b-4a9f-9369-a49e29c2cdac',1),('5fc55bdb-f575-4b2e-a17c-e39bc6fb9412','bea2ab00-c561-421a-8d40-5a0dca74e80b',1),('5fc55bdb-f575-4b2e-a17c-e39bc6fb9412','dbf26c31-238a-4e9d-8336-3daf3d2dd7bd',0),('5fc55bdb-f575-4b2e-a17c-e39bc6fb9412','df0e3ca7-158b-4a1c-b58d-9b423067475e',0),('5fc55bdb-f575-4b2e-a17c-e39bc6fb9412','f8102d42-5d53-4610-bf3e-aacec43635e4',0),('6b5c5dea-9656-4dfb-8a26-f699429eecd3','09dab565-4627-474a-ba4a-2c1dd91593a6',0),('6b5c5dea-9656-4dfb-8a26-f699429eecd3','23ac037f-c944-45b5-93f5-bd8395bf49e0',1),('6b5c5dea-9656-4dfb-8a26-f699429eecd3','4178816c-d39e-4c12-a66e-099a02427fa7',1),('6b5c5dea-9656-4dfb-8a26-f699429eecd3','44b8cc1f-4adc-4465-86ad-c7ba53036921',1),('6b5c5dea-9656-4dfb-8a26-f699429eecd3','7b89f755-1355-4936-a90e-04777d87b75b',0),('6b5c5dea-9656-4dfb-8a26-f699429eecd3','9f17cc3f-a321-4d62-aca0-b8a889d8629d',0),('6b5c5dea-9656-4dfb-8a26-f699429eecd3','9f9d3e99-a0c3-4e59-96b4-d9177299709b',0),('6b5c5dea-9656-4dfb-8a26-f699429eecd3','b679321c-7425-4c23-8834-e0209a8ae8c5',1),('6b5c5dea-9656-4dfb-8a26-f699429eecd3','be17e97a-a116-447b-bfb3-aacefeb3fc4f',1),('6b5c5dea-9656-4dfb-8a26-f699429eecd3','c83dfc09-88fa-4c04-bf69-4d9afd7d6d4e',1),('6b5c5dea-9656-4dfb-8a26-f699429eecd3','db7d6d32-c0a5-4b47-abf8-31ba61143504',0),('8b0a29f2-46dd-4f3c-8446-c0f63ad56ebc','35746b28-fed9-4ae4-92a2-9b42f120f272',1),('8b0a29f2-46dd-4f3c-8446-c0f63ad56ebc','3ab1f10c-ca69-4372-8fb3-dfb2bc434ca4',1),('8b0a29f2-46dd-4f3c-8446-c0f63ad56ebc','4b7b21fc-bddf-42ce-97f0-36242b13dcfd',0),('8b0a29f2-46dd-4f3c-8446-c0f63ad56ebc','70413ad2-e95a-413c-af6e-d1741fe9dfe6',1),('8b0a29f2-46dd-4f3c-8446-c0f63ad56ebc','8a6cb980-3684-4ef9-9d40-88f23541f569',0),('8b0a29f2-46dd-4f3c-8446-c0f63ad56ebc','a687fb3e-0aca-481d-a697-9c23538f74c2',1),('8b0a29f2-46dd-4f3c-8446-c0f63ad56ebc','a84876a0-9e4b-4a9f-9369-a49e29c2cdac',1),('8b0a29f2-46dd-4f3c-8446-c0f63ad56ebc','bea2ab00-c561-421a-8d40-5a0dca74e80b',1),('8b0a29f2-46dd-4f3c-8446-c0f63ad56ebc','dbf26c31-238a-4e9d-8336-3daf3d2dd7bd',0),('8b0a29f2-46dd-4f3c-8446-c0f63ad56ebc','df0e3ca7-158b-4a1c-b58d-9b423067475e',0),('8b0a29f2-46dd-4f3c-8446-c0f63ad56ebc','f8102d42-5d53-4610-bf3e-aacec43635e4',0),('8b786b3a-c968-49cf-b5eb-5e7b36889df7','35746b28-fed9-4ae4-92a2-9b42f120f272',1),('8b786b3a-c968-49cf-b5eb-5e7b36889df7','3ab1f10c-ca69-4372-8fb3-dfb2bc434ca4',1),('8b786b3a-c968-49cf-b5eb-5e7b36889df7','4b7b21fc-bddf-42ce-97f0-36242b13dcfd',0),('8b786b3a-c968-49cf-b5eb-5e7b36889df7','70413ad2-e95a-413c-af6e-d1741fe9dfe6',1),('8b786b3a-c968-49cf-b5eb-5e7b36889df7','8a6cb980-3684-4ef9-9d40-88f23541f569',0),('8b786b3a-c968-49cf-b5eb-5e7b36889df7','a687fb3e-0aca-481d-a697-9c23538f74c2',1),('8b786b3a-c968-49cf-b5eb-5e7b36889df7','a84876a0-9e4b-4a9f-9369-a49e29c2cdac',1),('8b786b3a-c968-49cf-b5eb-5e7b36889df7','bea2ab00-c561-421a-8d40-5a0dca74e80b',1),('8b786b3a-c968-49cf-b5eb-5e7b36889df7','dbf26c31-238a-4e9d-8336-3daf3d2dd7bd',0),('8b786b3a-c968-49cf-b5eb-5e7b36889df7','df0e3ca7-158b-4a1c-b58d-9b423067475e',0),('8b786b3a-c968-49cf-b5eb-5e7b36889df7','f8102d42-5d53-4610-bf3e-aacec43635e4',0),('ba9af61d-e6d3-41c9-88d0-ab69b58e061b','09dab565-4627-474a-ba4a-2c1dd91593a6',0),('ba9af61d-e6d3-41c9-88d0-ab69b58e061b','23ac037f-c944-45b5-93f5-bd8395bf49e0',1),('ba9af61d-e6d3-41c9-88d0-ab69b58e061b','4178816c-d39e-4c12-a66e-099a02427fa7',1),('ba9af61d-e6d3-41c9-88d0-ab69b58e061b','44b8cc1f-4adc-4465-86ad-c7ba53036921',1),('ba9af61d-e6d3-41c9-88d0-ab69b58e061b','7b89f755-1355-4936-a90e-04777d87b75b',0),('ba9af61d-e6d3-41c9-88d0-ab69b58e061b','9f17cc3f-a321-4d62-aca0-b8a889d8629d',0),('ba9af61d-e6d3-41c9-88d0-ab69b58e061b','9f9d3e99-a0c3-4e59-96b4-d9177299709b',0),('ba9af61d-e6d3-41c9-88d0-ab69b58e061b','b679321c-7425-4c23-8834-e0209a8ae8c5',1),('ba9af61d-e6d3-41c9-88d0-ab69b58e061b','be17e97a-a116-447b-bfb3-aacefeb3fc4f',1),('ba9af61d-e6d3-41c9-88d0-ab69b58e061b','c83dfc09-88fa-4c04-bf69-4d9afd7d6d4e',1),('ba9af61d-e6d3-41c9-88d0-ab69b58e061b','db7d6d32-c0a5-4b47-abf8-31ba61143504',0),('d2eb4296-2046-420d-beb2-25cf9e05b715','35746b28-fed9-4ae4-92a2-9b42f120f272',1),('d2eb4296-2046-420d-beb2-25cf9e05b715','3ab1f10c-ca69-4372-8fb3-dfb2bc434ca4',1),('d2eb4296-2046-420d-beb2-25cf9e05b715','4b7b21fc-bddf-42ce-97f0-36242b13dcfd',0),('d2eb4296-2046-420d-beb2-25cf9e05b715','70413ad2-e95a-413c-af6e-d1741fe9dfe6',1),('d2eb4296-2046-420d-beb2-25cf9e05b715','8a6cb980-3684-4ef9-9d40-88f23541f569',0),('d2eb4296-2046-420d-beb2-25cf9e05b715','a687fb3e-0aca-481d-a697-9c23538f74c2',1),('d2eb4296-2046-420d-beb2-25cf9e05b715','a84876a0-9e4b-4a9f-9369-a49e29c2cdac',1),('d2eb4296-2046-420d-beb2-25cf9e05b715','bea2ab00-c561-421a-8d40-5a0dca74e80b',1),('d2eb4296-2046-420d-beb2-25cf9e05b715','dbf26c31-238a-4e9d-8336-3daf3d2dd7bd',0),('d2eb4296-2046-420d-beb2-25cf9e05b715','df0e3ca7-158b-4a1c-b58d-9b423067475e',0),('d2eb4296-2046-420d-beb2-25cf9e05b715','f8102d42-5d53-4610-bf3e-aacec43635e4',0),('fe3360eb-3c93-4e56-99df-b8fd05dbf60a','35746b28-fed9-4ae4-92a2-9b42f120f272',1),('fe3360eb-3c93-4e56-99df-b8fd05dbf60a','3ab1f10c-ca69-4372-8fb3-dfb2bc434ca4',1),('fe3360eb-3c93-4e56-99df-b8fd05dbf60a','4b7b21fc-bddf-42ce-97f0-36242b13dcfd',0),('fe3360eb-3c93-4e56-99df-b8fd05dbf60a','70413ad2-e95a-413c-af6e-d1741fe9dfe6',1),('fe3360eb-3c93-4e56-99df-b8fd05dbf60a','8a6cb980-3684-4ef9-9d40-88f23541f569',0),('fe3360eb-3c93-4e56-99df-b8fd05dbf60a','a687fb3e-0aca-481d-a697-9c23538f74c2',1),('fe3360eb-3c93-4e56-99df-b8fd05dbf60a','a84876a0-9e4b-4a9f-9369-a49e29c2cdac',1),('fe3360eb-3c93-4e56-99df-b8fd05dbf60a','bea2ab00-c561-421a-8d40-5a0dca74e80b',1),('fe3360eb-3c93-4e56-99df-b8fd05dbf60a','dbf26c31-238a-4e9d-8336-3daf3d2dd7bd',0),('fe3360eb-3c93-4e56-99df-b8fd05dbf60a','df0e3ca7-158b-4a1c-b58d-9b423067475e',0),('fe3360eb-3c93-4e56-99df-b8fd05dbf60a','f8102d42-5d53-4610-bf3e-aacec43635e4',0);
/*!40000 ALTER TABLE `CLIENT_SCOPE_CLIENT` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `CLIENT_SCOPE_ROLE_MAPPING`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `CLIENT_SCOPE_ROLE_MAPPING` (
  `SCOPE_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ROLE_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`SCOPE_ID`,`ROLE_ID`),
  KEY `IDX_CLSCOPE_ROLE` (`SCOPE_ID`),
  KEY `IDX_ROLE_CLSCOPE` (`ROLE_ID`),
  CONSTRAINT `FK_CL_SCOPE_RM_SCOPE` FOREIGN KEY (`SCOPE_ID`) REFERENCES `CLIENT_SCOPE` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `CLIENT_SCOPE_ROLE_MAPPING` WRITE;
/*!40000 ALTER TABLE `CLIENT_SCOPE_ROLE_MAPPING` DISABLE KEYS */;
INSERT INTO `CLIENT_SCOPE_ROLE_MAPPING` VALUES ('9f17cc3f-a321-4d62-aca0-b8a889d8629d','afdaca80-10d0-4be8-b18c-4933151fb78c'),('dbf26c31-238a-4e9d-8336-3daf3d2dd7bd','bbbc95f8-5d80-4426-9b22-11e1b5c6db03');
/*!40000 ALTER TABLE `CLIENT_SCOPE_ROLE_MAPPING` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `COMPONENT`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `COMPONENT` (
  `ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `NAME` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `PARENT_ID` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `PROVIDER_ID` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `PROVIDER_TYPE` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `REALM_ID` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `SUB_TYPE` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`ID`),
  KEY `IDX_COMPONENT_REALM` (`REALM_ID`),
  KEY `IDX_COMPONENT_PROVIDER_TYPE` (`PROVIDER_TYPE`),
  CONSTRAINT `FK_COMPONENT_REALM` FOREIGN KEY (`REALM_ID`) REFERENCES `REALM` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `COMPONENT` WRITE;
/*!40000 ALTER TABLE `COMPONENT` DISABLE KEYS */;
INSERT INTO `COMPONENT` VALUES ('2055c140-8e1c-4d4f-8c12-5af9d7492f9a','Full Scope Disabled','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','scope','org.keycloak.services.clientregistration.policy.ClientRegistrationPolicy','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','anonymous'),('272acfab-2609-4470-aa38-03608ae64547','hmac-generated-hs512','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','hmac-generated','org.keycloak.keys.KeyProvider','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581',NULL),('29cddc7d-c913-48c6-90d6-73f99abb4972','rsa-generated','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','rsa-generated','org.keycloak.keys.KeyProvider','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',NULL),('319d59ac-77ea-4712-ac2c-21ff6f237fee','rsa-enc-generated','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','rsa-enc-generated','org.keycloak.keys.KeyProvider','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',NULL),('4c39c7c3-ef4a-4d77-94b6-e29eae7d79f2','aes-generated','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','aes-generated','org.keycloak.keys.KeyProvider','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581',NULL),('4db06bd5-28bc-4564-a464-005b66a32b50','Consent Required','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','consent-required','org.keycloak.services.clientregistration.policy.ClientRegistrationPolicy','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','anonymous'),('76f0e063-b83f-4308-bd67-af4985f83991','Allowed Protocol Mapper Types','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','allowed-protocol-mappers','org.keycloak.services.clientregistration.policy.ClientRegistrationPolicy','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','authenticated'),('84f3ab09-9205-48ab-bda1-02ee21730f96','rsa-generated','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','rsa-generated','org.keycloak.keys.KeyProvider','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581',NULL),('889c6ce2-796a-4817-9f2e-fd7e229e0c79','Trusted Hosts','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','trusted-hosts','org.keycloak.services.clientregistration.policy.ClientRegistrationPolicy','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','anonymous'),('a4ad6213-fef8-432a-bfed-89b8f27bde8f','Max Clients Limit','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','max-clients','org.keycloak.services.clientregistration.policy.ClientRegistrationPolicy','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','anonymous'),('a82eeef9-f37e-4cb1-a010-749ccd86a52b','Allowed Protocol Mapper Types','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','allowed-protocol-mappers','org.keycloak.services.clientregistration.policy.ClientRegistrationPolicy','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','authenticated'),('aaa5dd06-de5b-4f3d-9461-6bd3ae654622','Consent Required','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','consent-required','org.keycloak.services.clientregistration.policy.ClientRegistrationPolicy','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','anonymous'),('ae3ee582-a6c4-4729-bb0d-641df186f7fe','Full Scope Disabled','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','scope','org.keycloak.services.clientregistration.policy.ClientRegistrationPolicy','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','anonymous'),('bbf18037-cd03-4297-b42f-aa6fc206c80e','hmac-generated-hs512','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','hmac-generated','org.keycloak.keys.KeyProvider','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',NULL),('c2940d79-c0fc-456e-9047-6d3eef82c715',NULL,'5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','declarative-user-profile','org.keycloak.userprofile.UserProfileProvider','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581',NULL),('ce4a5c77-367f-4988-a090-bada4848688e','Allowed Protocol Mapper Types','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','allowed-protocol-mappers','org.keycloak.services.clientregistration.policy.ClientRegistrationPolicy','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','anonymous'),('d4f962bd-4d23-48b1-83a2-b12ab249c152','Allowed Client Scopes','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','allowed-client-templates','org.keycloak.services.clientregistration.policy.ClientRegistrationPolicy','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','authenticated'),('d977b6b2-0bad-4ac5-94ca-82c5644be1de','Allowed Protocol Mapper Types','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','allowed-protocol-mappers','org.keycloak.services.clientregistration.policy.ClientRegistrationPolicy','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','anonymous'),('d9dec3a4-2881-439c-85c8-3c0bd7f2e0c6','Max Clients Limit','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','max-clients','org.keycloak.services.clientregistration.policy.ClientRegistrationPolicy','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','anonymous'),('da8ab1c7-fb41-4121-bff6-9b43ff563220','Allowed Client Scopes','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','allowed-client-templates','org.keycloak.services.clientregistration.policy.ClientRegistrationPolicy','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','anonymous'),('ddfa79f3-dd3b-4753-9659-3266cb92b599','Allowed Client Scopes','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','allowed-client-templates','org.keycloak.services.clientregistration.policy.ClientRegistrationPolicy','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','anonymous'),('e223336a-ff9f-4616-bca2-2519d6e197b5','aes-generated','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','aes-generated','org.keycloak.keys.KeyProvider','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',NULL),('e77544a7-70a7-4a7a-b724-3b8f39b19c47','rsa-enc-generated','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','rsa-enc-generated','org.keycloak.keys.KeyProvider','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581',NULL),('ebc211ab-256f-471a-bfac-798284a67310','Trusted Hosts','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','trusted-hosts','org.keycloak.services.clientregistration.policy.ClientRegistrationPolicy','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','anonymous'),('f3839f85-b917-4a4e-bab6-1424a039329f','Allowed Client Scopes','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','allowed-client-templates','org.keycloak.services.clientregistration.policy.ClientRegistrationPolicy','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','authenticated');
/*!40000 ALTER TABLE `COMPONENT` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `COMPONENT_CONFIG`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `COMPONENT_CONFIG` (
  `ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `COMPONENT_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `NAME` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `VALUE` longtext CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  PRIMARY KEY (`ID`),
  KEY `IDX_COMPO_CONFIG_COMPO` (`COMPONENT_ID`),
  CONSTRAINT `FK_COMPONENT_CONFIG` FOREIGN KEY (`COMPONENT_ID`) REFERENCES `COMPONENT` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `COMPONENT_CONFIG` WRITE;
/*!40000 ALTER TABLE `COMPONENT_CONFIG` DISABLE KEYS */;
INSERT INTO `COMPONENT_CONFIG` VALUES ('07e9970b-3333-4f43-8f9d-ba506bd5b98a','76f0e063-b83f-4308-bd67-af4985f83991','allowed-protocol-mapper-types','saml-role-list-mapper'),('08df0bc2-01cd-408a-940c-37d0b312c579','e223336a-ff9f-4616-bca2-2519d6e197b5','kid','924360e9-dd92-4545-a8d3-f64023668f4d'),('10282eac-b214-45b9-b994-845d850e7118','d977b6b2-0bad-4ac5-94ca-82c5644be1de','allowed-protocol-mapper-types','oidc-usermodel-attribute-mapper'),('1093ed4b-90a3-4344-81d9-2c74c91a2cbd','a82eeef9-f37e-4cb1-a010-749ccd86a52b','allowed-protocol-mapper-types','saml-user-attribute-mapper'),('117c9683-5902-4585-afb2-4ced23f04a0e','ce4a5c77-367f-4988-a090-bada4848688e','allowed-protocol-mapper-types','oidc-usermodel-attribute-mapper'),('11d13391-8434-41c2-909f-2f139ca8e110','a82eeef9-f37e-4cb1-a010-749ccd86a52b','allowed-protocol-mapper-types','oidc-usermodel-attribute-mapper'),('16f96b10-a980-47dc-918f-34d235170316','4c39c7c3-ef4a-4d77-94b6-e29eae7d79f2','secret','_Ge_f7a5-mkn1ejcxkNxkw'),('1cfccb9a-c42f-4a05-8dd3-f2853e225a57','76f0e063-b83f-4308-bd67-af4985f83991','allowed-protocol-mapper-types','oidc-address-mapper'),('2037a92e-0e94-479b-a601-f1ced04acb33','e77544a7-70a7-4a7a-b724-3b8f39b19c47','keyUse','ENC'),('2171c3a2-9f1e-4272-ba76-a535f146b188','d977b6b2-0bad-4ac5-94ca-82c5644be1de','allowed-protocol-mapper-types','saml-user-property-mapper'),('23c4d9e5-470e-4700-ac50-6d3ab2aa975e','d977b6b2-0bad-4ac5-94ca-82c5644be1de','allowed-protocol-mapper-types','oidc-usermodel-property-mapper'),('26214863-67e2-4a35-bb2f-50a13a700bf8','a82eeef9-f37e-4cb1-a010-749ccd86a52b','allowed-protocol-mapper-types','oidc-usermodel-property-mapper'),('2e2727a5-c5a4-4669-9a15-48a648f0841b','29cddc7d-c913-48c6-90d6-73f99abb4972','certificate','MIICmzCCAYMCBgGbLi/oqTANBgkqhkiG9w0BAQsFADARMQ8wDQYDVQQDDAZtb3RvZ28wHhcNMjUxMjE3MjExOTEwWhcNMzUxMjE3MjEyMDUwWjARMQ8wDQYDVQQDDAZtb3RvZ28wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC2TPlI3nnU32HOvi+CJ2uFHi+hwdQc5RKHYKdddKanSfGgltwkhZB2uc4l95mGSkYt8FzGx5RRnRP7P6ctS+p1IcV74OWn5lb8rCcD3CzMVGN7jsGazLqytaLo8Bu8JTdk5zOx2vN/5ALhaYPscwSuOlEJnx8a52VWSaxbvaXJQEVC29hb5jr/sDZUdlZ3dou56kYIovgO2+jwhxZ4f0Yx+yjunzAuQPwhf6cJincm9ZS6+dBe07H73/auWtOnTCrppvIqIUFvgYicwXGaMlAiSWm9ojEJUi0KrdLOE3o1I1KhEk/AMjsEKhe7G686PZAvciTSY5zXnr5OnADqEIPfAgMBAAEwDQYJKoZIhvcNAQELBQADggEBAKvMpstew4iItE8Paf2YfYrXBr+3PfdUnwgZ+OohbQ6nCFdqxKmKVyQgrD33mLaVuKWxMFxuEF+N+QEx4QtjtzDu5DYvGkKz65ekkwpLI954Up4eA4TNS09k4ZjKdqGsb78zP5VnRsvX6sP9K/yr/RGmf5PepXhSBMhV4ix73dhE/KAmgV9ttNJQ4qNaJxog5Cyujy0RlhSGsdxL5qKJ+2QIzu0PyYKJ7q1jNrEgepg9cFGDK9dAwT9MeL6qtYZvm4a+TyEwZP8AFbwwjpEANCJ6kuweWr2U/hE/SztuIjY//8qssPyhcjw66z1XgyucffmGwOg91U4HDI2Dyb9Tk1Y='),('2fb13193-f960-4b9d-a8db-6d2c2ac67d4c','ce4a5c77-367f-4988-a090-bada4848688e','allowed-protocol-mapper-types','saml-user-attribute-mapper'),('3a425e1e-053a-4e72-85bf-eea913a61871','ce4a5c77-367f-4988-a090-bada4848688e','allowed-protocol-mapper-types','saml-role-list-mapper'),('3cce17fe-919c-4160-bea0-4735970d6657','e77544a7-70a7-4a7a-b724-3b8f39b19c47','algorithm','RSA-OAEP'),('3cfd19c9-b2b6-4e44-aa06-bd892ff8422f','272acfab-2609-4470-aa38-03608ae64547','secret','NyGOiXvtC5NgdpzRUOLtA5K0m_CMefLA4_psTJEy1ptSDUk3eHaagmKU_VzbfQeENqtL3ohrWJocZ0bnwj9BIhoWYyw0AO5o-SXgxH5b7RLU242lF2XJVfaR6glKaae9Rwr8lPpfe1sRpdaA0RS21nv37cqvV23bzFk9aMbUDqs'),('40b9c614-ac82-4379-a3a9-6b28a2029c23','c2940d79-c0fc-456e-9047-6d3eef82c715','kc.user.profile.config','{\"attributes\":[{\"name\":\"username\",\"displayName\":\"${username}\",\"validations\":{\"length\":{\"min\":3,\"max\":255},\"username-prohibited-characters\":{},\"up-username-not-idn-homograph\":{}},\"permissions\":{\"view\":[\"admin\",\"user\"],\"edit\":[\"admin\",\"user\"]},\"multivalued\":false},{\"name\":\"email\",\"displayName\":\"${email}\",\"validations\":{\"email\":{},\"length\":{\"max\":255}},\"permissions\":{\"view\":[\"admin\",\"user\"],\"edit\":[\"admin\",\"user\"]},\"multivalued\":false},{\"name\":\"firstName\",\"displayName\":\"${firstName}\",\"validations\":{\"length\":{\"max\":255},\"person-name-prohibited-characters\":{}},\"permissions\":{\"view\":[\"admin\",\"user\"],\"edit\":[\"admin\",\"user\"]},\"multivalued\":false},{\"name\":\"lastName\",\"displayName\":\"${lastName}\",\"validations\":{\"length\":{\"max\":255},\"person-name-prohibited-characters\":{}},\"permissions\":{\"view\":[\"admin\",\"user\"],\"edit\":[\"admin\",\"user\"]},\"multivalued\":false}],\"groups\":[{\"name\":\"user-metadata\",\"displayHeader\":\"User metadata\",\"displayDescription\":\"Attributes, which refer to user metadata\"}]}'),('412f5cb4-4e27-4f11-920e-4fc2d1ae257d','272acfab-2609-4470-aa38-03608ae64547','priority','100'),('4a8b37cf-2ea0-4e8c-935c-5942a2f88921','a82eeef9-f37e-4cb1-a010-749ccd86a52b','allowed-protocol-mapper-types','oidc-full-name-mapper'),('4aecc4f9-fe23-4b3f-b901-2b92f04ab25a','d4f962bd-4d23-48b1-83a2-b12ab249c152','allow-default-scopes','true'),('4af2070a-af57-4ff8-b179-84b508479fc5','84f3ab09-9205-48ab-bda1-02ee21730f96','privateKey','MIIEowIBAAKCAQEAxC1gCIRt0JQiRPGb44KZ7Xd9QiLDmNzZ75SMcISHhnFkjJ2vPFyuC4kgx3eJMQWPWXnAePMphD8orspQ48ADIx+7JJzqBCkgGxYszgasq5eRTt1cHDKPKwOfYZtIa2wyT7D9mQK60U9fJS/HhunvLDVii3noYSWPXq2t46aNHkRO9bp0ElfQbggDAatAVM7OZZQRan6ROu/8TfI6OAzEVELDy02Yl/NZBqBWuCcKDFyhR563YxsY2Z7AyJKirBh2c9xO5n/ToDwyMJgwnS1Xw98ivjzjkL6Z8yk4rqHQqJAfqQD+BXcS1LAXUvnBOmIUwSt1xr/U0/xrAP4N6abW7wIDAQABAoIBAGExj+o4XkXgFohs4cyauJOHx0W8fobeKJWD1f6rEpAf3u6x8w9nZY1uj5Bzn0Nr20U7xonGA3QzbH8I2oGz8Hn27xupKGIOVAtHdcRPVz4edk4J6uhHfloTzeOocP1T2itUZ5hm21L3TY7ZeLPVc/tFqcCrkd+LodgJGWV372OmGfDRKq8QjDfTumqYgQRK0YjvPssselMzMS5z95yB/Ci/o3aIomwbpY0CdDU6p7V/5jbDOLif8aCYknvgYMgIG8A3SuPNx+D4LXZqp3pUINqS8kIH0lXYIdASWCtUMpKJyelFE1u8dERJQcVAtm05nKQoLmJtQokVVbyXSQl4fIECgYEA+HI4QJFloJvEc1Nh99YZNa5k3f/AvkwsO/XwlFz16J68Fd4siCWeEz/Ym4DkrMA0ic2LIC1YNWioeX6TTCCZJqVskqX8c1hrSocEKx1QLely0Vx/rC3SRZkjtG+JOvsz1bQfzgId3GCJpR+daDDsl3r23qf46G3yF5qVp+6FNoECgYEAyiRSARrz3y94BorTwC46PWOz9Y6lAhBR+D0QhWl6uUO7ty9gXCRfU+hsLUXguru+GN3+u5nvKjDfa9NN0T7Rt2lmBmU5pfrbUhaT+pyBv02KLFdn9dW5HJiXHU2u3ldX+P3zXidTfRlqzoGZ1xSFNTKVL0xyIr8YQ0LrmK7btW8CgYEA0LecHFpJp1jU/BGzvvKFcTyWDo56TSRPEPx3ZRyHliosVq60eJgpqYA8p77r3VkHijE8TQbrEMtgpdi5Q/4qy3WztatuXzSLxiuRS+svXbK+O7Zap7pC6SdTAJWnEv7QOBs/kODkVyxe1nKt2Ucoq8Woubshd5pOBJhgW5kdHoECgYB0QDMRA6kgv7HtQUXKDkvcJkzGecSXzT723ld0Aldb2IJKNWrbUuen39DmsHsvqO0IbBMlXL/Yqcazwg/MuK3vJyRPpUNkG0KqkOWFiSBtHnDKU3IsRT9yHKc7aQOaRu0G+zVM/pzbIxTzskEBGoIJOA/n5bT2GqQfeor3gghE7QKBgGucNCkHpW4zCVEe/ohPzHwHzKeFdN1CppSBms2Oj1q5chu55h+okQk5yn9s1pQYQBESXKGhKqUltjzRGiUk2ZWNv4nlq8e3v01DetZw7pX5zqcZpIlqaHmZgMFBTIw6/g8Kxu7ocDRlzUINtXWZdKGnU/iM14NY0b2lRZCe4QwT'),('4efdb731-1114-4397-9b9b-4a06bd05c468','ce4a5c77-367f-4988-a090-bada4848688e','allowed-protocol-mapper-types','saml-user-property-mapper'),('51daae90-a281-48a9-9c06-99e00f250c0b','bbf18037-cd03-4297-b42f-aa6fc206c80e','kid','bc922b2c-fdef-4f41-b688-68580ef183fc'),('5280b6c0-d9c3-415c-92a6-7903021e0679','a82eeef9-f37e-4cb1-a010-749ccd86a52b','allowed-protocol-mapper-types','saml-user-property-mapper'),('532b8d7d-da62-4a0d-ab28-fd3570713e6a','319d59ac-77ea-4712-ac2c-21ff6f237fee','privateKey','MIIEowIBAAKCAQEA2BABa64w3RHfeY9RLV7GTn8i+CvCgZcU9nwfoySbpPa9hpAL+9RUds5gPdN3c1Hc6DJHQH30qmzKfm11fvgRfiin4ycXTUivmUsoI83dbKhTA0avmwLWHNqF2gcBTndzvRi4NMAiRF0gzIFeOECiGxrdpjon3KurUb/ZUKxYrju6Iqzop7k+7+j3tj2uTyLV+O+0EeH4x3rs+3y+RqaTY5vTd8h5UpQ86qwxcirUYLdyWjwZRiR3U0sK3WPfZeBrf4+lA4wWXkHQAfpa/uALgxoNBLfb10kKPACGighLUGSNHDj1CaUl8c08LhbhmIZ/lMLoe/NB+I3ESJdPunMDnwIDAQABAoIBAApXScTI7ixxLIxO4F3BiXrm+YTTOK38eH5g68q3zSiUvqFwHNgNY4rZecsREfMtBIZV0bSJKBVNZlZYECWRjeochrHdLxpTKO8MpSaBizV7dkL+i8Jph6yTMOYsgP2BGMMA1IrXH0y9AFLeu48JrjiJ4qM6vqmrIGSl6qrCdO7HH9KmX6+C8TsntZDQNmhjBahK+a/pIoKszcTU+Eee5+9nt5lySk/EUtBZqedR2QNgtZpkmxwWPPrZGhQX22PtTe2fHb9DBJX4oKKhwP5B776UtNtCs3Gpcdt6aKGl0S0HsZeFuKGmmB4w8e9iCRYgu7xymZ5d7CY4dy/tbyszTvECgYEA+CHztQi5f5aAdht4yqJaYawNbz2phXxAS77ag4+gy+CMXMS9PgKpKKy4gZumvab5fapLDdeKNPbdPaU2mrugZ9s62wGgO/GMmZkh/ZTPEmWDNceBF3/ZFlgfQf/ytI2urr9P2IfS/MuGX5tbn20luVPbTBO3cjcIoDNlcF92rLUCgYEA3um/DPHecTVP2qGklC/nWw40waDKSizeo+7ETuhGnj5eYqLIwx2JeWPWOvBhkMbEyS2OxBkfeibUBGcVbJ6ojEyXVHNbyZXn+29XtyLEeOIHFJ3UAmTbx06qW+BzfPSWiC/OYiYmBpbfopl7Eg638188+aRJGZQzWzkAG0N+94MCgYB1g67l6BHiKScByJT2ctmXNGLQckWc4HGa7fKTKCyhIeGX9d0iRyRSDSr8xvi7DR9Psp0FGzVrTBcPoPUxZvba+Orw2d5HRWivPeTpMda4AbeYernapUPs3yM+oLcTZ555SEXKaNk1vDYRhYRpBso88UgEybmxr/q31F+QTsXnmQKBgF+LHKI1wTwQt6tkASCq76Ttj9s2Sb5FQSBOpMGVeT82SEpKr3ZB4Lyw+A/HHGprZ8k6x6hxhDZdEh7Y+CvYTTHlWi+OSlc683WnoZcTRdbyOkqx4z8rb7AbnZVHvPhf9W1l4DXdWAsW+7+7vS5eNKYRuhbePZ7DEMgWxF7rLS51AoGBAONYRcdepOnQHKrKW8p5WVwOCUoHo2XkXnwlW+IX46BFb/m2PAkCYopLbZ/HvNZYaS65LtKJrjlj8Mbecmus0wAUCot9mVL5oZz6tie5jtawTzPa1puK4LMRw2IPnEX3TmAak598a3jcIkTh3zqCCvUSq74OkhyoOOr7FFnMV5LB'),('53ade693-16b9-4cdb-94cd-ed7cfb4b27ac','ce4a5c77-367f-4988-a090-bada4848688e','allowed-protocol-mapper-types','oidc-address-mapper'),('56de67c0-4ecf-4b09-8915-6713b6f0310b','889c6ce2-796a-4817-9f2e-fd7e229e0c79','host-sending-registration-request-must-match','true'),('5b418d8c-9030-4466-9678-5e13b0b25a1c','29cddc7d-c913-48c6-90d6-73f99abb4972','privateKey','MIIEowIBAAKCAQEAtkz5SN551N9hzr4vgidrhR4vocHUHOUSh2CnXXSmp0nxoJbcJIWQdrnOJfeZhkpGLfBcxseUUZ0T+z+nLUvqdSHFe+Dlp+ZW/KwnA9wszFRje47Bmsy6srWi6PAbvCU3ZOczsdrzf+QC4WmD7HMErjpRCZ8fGudlVkmsW72lyUBFQtvYW+Y6/7A2VHZWd3aLuepGCKL4Dtvo8IcWeH9GMfso7p8wLkD8IX+nCYp3JvWUuvnQXtOx+9/2rlrTp0wq6abyKiFBb4GInMFxmjJQIklpvaIxCVItCq3SzhN6NSNSoRJPwDI7BCoXuxuvOj2QL3Ik0mOc156+TpwA6hCD3wIDAQABAoIBAFYACBX1RlpDK3wBwp3PBO0eqW/Q5cTf2BbfZPU0o64FtySAJshf7h/kVbkZjmaIwYTIgKwR2QmLZbLAe5hl2Ox7CxP+q7jCdETs3db3VnbqVKgkV3c0n9rQz1K+ewcOM4y5eWytEvRBV46JXX/9XrbmsfYtB9qEUQ0c4l/BHuXKTBvXh7QpwUz2Qldwor+wf+bhbK0BAY4rjZaC4ofgTF2+PckhZ9Z18R1KLqRxEbGC3IlrnfjNGDdQG1yIeFyCA2XY8Ug+NkYmaXQCN4QS/HzoSA6f+uqmKVmLUEIV1+Tcif929YOCsNUgZ57p39roA5YtlX6w6WWuCVzvmDKLm3kCgYEA8MLsZMrwdDNT8WJhSDzmpuwHv7U2Foyt4TkMPQ0hGBzVknHqoTQfuUs0RSc/Nk0n0w8nDNd5fd1VebqPskaENJl6aHbajGGzLzLaeXApzxtQaxfqEIGTur8TJbXsqhvZqg67QCDz49/0TerKdymnkFrg14WSCFYcWUHJiSsoIN0CgYEAwdbOaKWsZK4KtvvHUbXA27Fdcs5KQbC0nnuRh0xLKIa2pzXKszqigKOK3hLzdPY3hyCqFojb2viOjm9T4zoPQ5gAj3yUxvxlM61bEgu5MPduu04SmcfhyVkN2PntjAVpXVqdTWv8t3mm4cHuLlh52+i6gdG/BftKnEVtGZtJresCgYAhFNeAffYOMWNOXpM5ZszWkMX7zH1zGb9mm/0Xgm9SrduWtLpljhDx1+iFu//OiUl4kbNu60iRR9FP2ZeXCmNuCsyr//DFq0MAsD8ewo+sGnexT2bk/7j5xvltjBLJyOvelwmZpS192S7raeCKARo6FpKgaoUHSjW4+fYkN97arQKBgQC+THiP3DaHESaLxltEKL+R72rizwuvuh9Yh8+zmH8g5kBjoDfEBYiL6gLnSouhBCg3jSXuwaW2L9cklr8wkUty/MNwGsqlIBzqwwJpkQDaxc/8KwWv28AawEcMATHpCVcLw24zGBUeu1vOv2lqLopUDwMreBOld7MkADS/suh+yQKBgBbHps5LqnMmEYV4Jl48c0z31iJWz620pvxxJcgrPZKxAzm26PiKnD7IEaLlsXFNCz8D34SyiIO6/1gdO+AW1SWGqlUW9U0LjnJXvKMjR3DOCME+7oiXbehtcg2B03Y6bTWoNnf19f5vk5VuatZflcNnG6tWuMkRXjxUwURzfaJ7'),('5df80b9b-2fdf-4068-9997-e262cfe297c9','ebc211ab-256f-471a-bfac-798284a67310','client-uris-must-match','true'),('653dbbd2-b6ec-4139-b7c9-acd354692c8f','bbf18037-cd03-4297-b42f-aa6fc206c80e','secret','SVSRs0cy3hfjSA3UCwEqZref0L5kynr6rTunkudvoIjXKsQ45gw7ZLdCZ1OMxUQNc1QpVlj44xB_qgew3LryVy5VUXmG8OaXTUoqajETOMsgzSYZixg8yimI0YGZqniv-Jslt3-gniS7IGnDysOduzb2nDRwKQdIyASQij3pzaw'),('6ae49e15-ab78-483b-ae17-410c12a0b949','d977b6b2-0bad-4ac5-94ca-82c5644be1de','allowed-protocol-mapper-types','oidc-sha256-pairwise-sub-mapper'),('6cd516af-5c27-4b3d-9d50-4f6abce13e89','d9dec3a4-2881-439c-85c8-3c0bd7f2e0c6','max-clients','200'),('71d3a71c-9358-4786-87fc-7f0efa5f9f0a','76f0e063-b83f-4308-bd67-af4985f83991','allowed-protocol-mapper-types','saml-user-attribute-mapper'),('74864388-75db-40da-9a5b-b6d5338e6e07','319d59ac-77ea-4712-ac2c-21ff6f237fee','keyUse','ENC'),('749cbb19-8dd0-4011-be36-e474bf265fc0','29cddc7d-c913-48c6-90d6-73f99abb4972','keyUse','SIG'),('755dd126-8757-4add-9533-ccd130b620f3','d977b6b2-0bad-4ac5-94ca-82c5644be1de','allowed-protocol-mapper-types','oidc-address-mapper'),('7a89fafb-81b5-48d9-b111-3b232739b4e2','ce4a5c77-367f-4988-a090-bada4848688e','allowed-protocol-mapper-types','oidc-sha256-pairwise-sub-mapper'),('84f798b3-83cb-47a3-948e-8e55538b92a4','272acfab-2609-4470-aa38-03608ae64547','kid','e0ebff3e-44b0-4280-996b-29de22786649'),('8612bcb0-c7f8-4583-a843-13ec0f2209dc','e77544a7-70a7-4a7a-b724-3b8f39b19c47','privateKey','MIIEowIBAAKCAQEAzf1y4sRH88pgpnBaNrmryTR3emR8O+n0WEjeIbcsA/YiKjMJY3wEYIJ4R11cdha1rMc+qKjIb3G+A8Y5FtUPI3ZXrZggWnr91K78ur85/bGFrtAj+gEM2nIT4/SizMhvCAHfLuobw/2BamwgPaT584bm/CvXSUJP75yTl4mGSRLoYNm6muFzCUaCT12clExAxLOkXdKa5C6pqzarIkeqwuTmzoVbgz9faybUWMAth8c7oJkirUIAvuWVI2MeVG8HgMm4K6kPjXEmyk+xy/UQ93FKgwXr1sCOn4exIE/bRdkrrhpnANpTi5gRD9RxWCn2UVNfXgC6l5wozyXuf3LvLQIDAQABAoIBACLBqwrcN4/JC/nJZ09okLXIZoh+QAAT/6iN0v05XPHmK+mW5vSkTlDCcu0HsbwzJMBcdhQbNoanWHpmjeduV7mSFEU0L/FUY1Ppmrcf4PcPZxESCZ/7YWHSccy60GVhRGrDdhcCZWDye02rpN0B4zLWBk9yTkcMK1SwxK6DA+d8PbLECBHUhtMA4ao9ZiigH8ap1Q62wEFv6zJazEPXhYEACkD4ad2d2uX1xmNHqpU4hLLK5Skux0dQTrrVsRS12+kFsk9gjrhXLFdiR46eivDsidu8NCd1bFfTaZ5HvdDZ9RRTroqsHQRPO06fRCOYCi/VtrVvNkfvYtRRfANyQLkCgYEA7c1NU7MSnNK1UdwP5HiaOdeuKWmQa781fUSPhbkOEBby/o1Lh8hpofSyfz4bZd9ZT/SOlAfsNRrVk40rZ6RCTHZGNHWcO/qD0DMIJ2w9zzGtLIdxacRCsrkfAMyfjJ6w6czQOxrv7weDI7jhMFcGrxc/Fp8AT0ue17QjGZakAFcCgYEA3cDuFBNdGwcdd385hccwX8KGnwKdysJBP/ro7QwqD+T/rQ7tXT553ckKVzWi1qOWqumVyU7Mo+mGKol5mjQcmniI5qHXt3LACi5MzA394GnHBBgN0fMnyhdmwcokbcDzzK9QH84l2cmTA7Th35EJxXJ5/t5xK1/cSUmVOpDIihsCgYEA4jku2PNAP6RMDiX1kHoS7ff6nTbpxzbNoe/5QVN8x12tqO+t7OyBYSsyK85frhg+FMwoXzbyRMrUZ1wi8KPY1nil7P39tSlrAEOr/CUC3/r4LaZkBCRK908Y6V/AYLwOvTy006Wan6Bb5z3YLG5900+gGtsEJmtgHwAlbBBc0F0CgYBhuDMyBZZx4qDJl/bxTOW0TTSlvbFrOE8ospkSXNnsRy2kRX+0J2PsqDs2Uxx/GXe7uVOYcnEbijzlYc+EdJJmP9eUSp2U6axW0DBne4L08agR1cfBTTWriMRvgXSsKOtGt7rHSUDEyF9QstUWgAOeTuQTvFKCmquoYBSyyjqVIQKBgA68K0LSvT0jrIhJ81qwlj7pqtMTrxHqpGxalMLGs0uZ1zCLwbWL9MAi6csbDtEdh99hQ/Rpbq5E2zE2UUU3wOyU7SPfJ1jObyRoDDlbsdY5qwWPiuguUML27vr8iRIWILyOWR85cwbrpvya8PFlMSRjeB3I9b9iXXlhPdyDLrWV'),('882b0ceb-6146-49fd-8565-34ae484975cf','84f3ab09-9205-48ab-bda1-02ee21730f96','keyUse','SIG'),('8bc2cdc9-9621-409a-8b2f-a5694e6a57e2','a82eeef9-f37e-4cb1-a010-749ccd86a52b','allowed-protocol-mapper-types','saml-role-list-mapper'),('8e4ef8d5-897d-4cf2-95b0-dc3be258ff64','319d59ac-77ea-4712-ac2c-21ff6f237fee','certificate','MIICmzCCAYMCBgGbLi/pHzANBgkqhkiG9w0BAQsFADARMQ8wDQYDVQQDDAZtb3RvZ28wHhcNMjUxMjE3MjExOTEwWhcNMzUxMjE3MjEyMDUwWjARMQ8wDQYDVQQDDAZtb3RvZ28wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDYEAFrrjDdEd95j1EtXsZOfyL4K8KBlxT2fB+jJJuk9r2GkAv71FR2zmA903dzUdzoMkdAffSqbMp+bXV++BF+KKfjJxdNSK+ZSygjzd1sqFMDRq+bAtYc2oXaBwFOd3O9GLg0wCJEXSDMgV44QKIbGt2mOifcq6tRv9lQrFiuO7oirOinuT7v6Pe2Pa5PItX477QR4fjHeuz7fL5GppNjm9N3yHlSlDzqrDFyKtRgt3JaPBlGJHdTSwrdY99l4Gt/j6UDjBZeQdAB+lr+4AuDGg0Et9vXSQo8AIaKCEtQZI0cOPUJpSXxzTwuFuGYhn+Uwuh780H4jcRIl0+6cwOfAgMBAAEwDQYJKoZIhvcNAQELBQADggEBACjhE+M6ow5O0fYiVvx9Ha1ilqracaX5QMq6q/Ts8ChEacgw0jANjrKO5bzjwA96Al28moQuDtHhW3/bYfG3pZotN3bnnNlmTOmeDJ0HXYOxgDlSIeNKWxt453QJGqn3Y1RIZyFDKmtH+8aupQw+PJW43KHL2JhhMCZv5ZvL5ym/3/r4+MBChaFe3MWd/8BwsjgFobZ5Pv+WwokhPD9s+K0nMtH7m5euSroWK/lsjEx6nsODaBejztwaFR4ddXK02WGvkycSl/EkKimGk1TLcTjFe0KHUPNhnY6dQAXWTcq39i0ykga75j6ILFc9cYYU0Tqx171YWmR+7cMtVqKEko8='),('971b32ed-dc93-4005-ae18-ca20009d9295','e77544a7-70a7-4a7a-b724-3b8f39b19c47','certificate','MIICmzCCAYMCBgGbLWwECTANBgkqhkiG9w0BAQsFADARMQ8wDQYDVQQDDAZtYXN0ZXIwHhcNMjUxMjE3MTc0NTEyWhcNMzUxMjE3MTc0NjUyWjARMQ8wDQYDVQQDDAZtYXN0ZXIwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDN/XLixEfzymCmcFo2uavJNHd6ZHw76fRYSN4htywD9iIqMwljfARggnhHXVx2FrWsxz6oqMhvcb4DxjkW1Q8jdletmCBaev3Urvy6vzn9sYWu0CP6AQzachPj9KLMyG8IAd8u6hvD/YFqbCA9pPnzhub8K9dJQk/vnJOXiYZJEuhg2bqa4XMJRoJPXZyUTEDEs6Rd0prkLqmrNqsiR6rC5ObOhVuDP19rJtRYwC2HxzugmSKtQgC+5ZUjYx5UbweAybgrqQ+NcSbKT7HL9RD3cUqDBevWwI6fh7EgT9tF2SuuGmcA2lOLmBEP1HFYKfZRU19eALqXnCjPJe5/cu8tAgMBAAEwDQYJKoZIhvcNAQELBQADggEBAJey0oVbJaWnxylERrsl0Hy/qOOLRmSQHVMwGSz1fpkrXV+yHS7m34lVsICbZy3+XREkJ8CNfaD/sqOAEShMMTxJVUfKQ8f6aH2X2k3lJVUMAM6yOeqUk9GYUJvSXpJMvZ5AwJ+/2gnnRighdvCyJeQTrzvXCAG8+SfjNkNy7i+eamYmiTp3wHYsiv7bB8r0pJ5Owu2OQvAT+PobPIifvhWoXHgINwQe0rAG9phfs1ysA5nVwKH4dCRZTknnU7Pt99ZCq+81qLOJiXr7GEwZ/KjMY4UGOyJyi3HVR/+bZe8EbfupqeWXrBD5VJ6Ox+sSEefmmXbj7qj27vPFyru0IgI='),('9b275b08-a5c6-46c0-8241-7bb845a7d3bb','d977b6b2-0bad-4ac5-94ca-82c5644be1de','allowed-protocol-mapper-types','saml-user-attribute-mapper'),('9db90e2b-0f7b-40bd-a730-70a2ce2b6080','84f3ab09-9205-48ab-bda1-02ee21730f96','priority','100'),('a61e3de9-618a-4d9a-a431-397a92223a1d','76f0e063-b83f-4308-bd67-af4985f83991','allowed-protocol-mapper-types','saml-user-property-mapper'),('ac6d71d1-e13f-4b2e-abd4-4d3c66444142','84f3ab09-9205-48ab-bda1-02ee21730f96','certificate','MIICmzCCAYMCBgGbLWwDgDANBgkqhkiG9w0BAQsFADARMQ8wDQYDVQQDDAZtYXN0ZXIwHhcNMjUxMjE3MTc0NTEyWhcNMzUxMjE3MTc0NjUyWjARMQ8wDQYDVQQDDAZtYXN0ZXIwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDELWAIhG3QlCJE8Zvjgpntd31CIsOY3NnvlIxwhIeGcWSMna88XK4LiSDHd4kxBY9ZecB48ymEPyiuylDjwAMjH7sknOoEKSAbFizOBqyrl5FO3VwcMo8rA59hm0hrbDJPsP2ZArrRT18lL8eG6e8sNWKLeehhJY9era3jpo0eRE71unQSV9BuCAMBq0BUzs5llBFqfpE67/xN8jo4DMRUQsPLTZiX81kGoFa4JwoMXKFHnrdjGxjZnsDIkqKsGHZz3E7mf9OgPDIwmDCdLVfD3yK+POOQvpnzKTiuodCokB+pAP4FdxLUsBdS+cE6YhTBK3XGv9TT/GsA/g3pptbvAgMBAAEwDQYJKoZIhvcNAQELBQADggEBAJmgLJIbDPsEHnbisgbDttFtqG6pulf+vyp859IBh3UGeeEdsmacQbUaWxoUbVoa+c2QmkokPsZDgUxr1/zkGmLK5HAOsR2/JBTjPiqfYoNv/NgqdhXW1raN24nfC3ynNWEaK7Z4YTKPl/5X6hs5w2d2vma3mPtv/BXnGfHEJ2X/N2nu3b7yyQ+XCFKsyYyVnn0U+JhNZPrLBxnlc6PxhaFWv1XuiDxebLrpthTBCc4O9N8o0qLrHMuvldez+J/KQ2ZVdk8lKgQJQQVYG6QJ6j2pND3tLNdcAaSPRzYZwdVQG2/HODjbIqavapolZ1SFr2xHFwBq+pxc4vqhq9auya4='),('ae4d8efe-cb64-4872-b82e-064706f944f9','e223336a-ff9f-4616-bca2-2519d6e197b5','secret','dqYy5tnC6yee_DQMFeYruw'),('b3a388a6-f2eb-46b8-9682-5fca43e34546','29cddc7d-c913-48c6-90d6-73f99abb4972','priority','100'),('b3ff2a04-0d6c-47e7-8f44-2e08995899da','a82eeef9-f37e-4cb1-a010-749ccd86a52b','allowed-protocol-mapper-types','oidc-address-mapper'),('b8e94d46-9487-4c79-b7a0-efa7fe1a0909','a4ad6213-fef8-432a-bfed-89b8f27bde8f','max-clients','200'),('b9c9b202-2858-4bfb-8f8f-c29288c7ba35','a82eeef9-f37e-4cb1-a010-749ccd86a52b','allowed-protocol-mapper-types','oidc-sha256-pairwise-sub-mapper'),('c67f774f-13bc-49f0-b4f8-c38e53563f31','ce4a5c77-367f-4988-a090-bada4848688e','allowed-protocol-mapper-types','oidc-usermodel-property-mapper'),('ca37d240-ce21-4f22-8930-f1dc098a1636','319d59ac-77ea-4712-ac2c-21ff6f237fee','algorithm','RSA-OAEP'),('cb1ba56f-6380-4b1b-b2fb-2bf3bdb1b04f','76f0e063-b83f-4308-bd67-af4985f83991','allowed-protocol-mapper-types','oidc-full-name-mapper'),('cb3459f6-feb8-448c-92f5-b8ce6404e2c2','ddfa79f3-dd3b-4753-9659-3266cb92b599','allow-default-scopes','true'),('cc01c5d3-b2c7-4417-aa52-6f42e49d6a1c','d977b6b2-0bad-4ac5-94ca-82c5644be1de','allowed-protocol-mapper-types','oidc-full-name-mapper'),('ce6206da-1066-4462-a884-88723dd954f7','76f0e063-b83f-4308-bd67-af4985f83991','allowed-protocol-mapper-types','oidc-sha256-pairwise-sub-mapper'),('d2b95055-517a-4501-9a9a-8b963030bc2c','bbf18037-cd03-4297-b42f-aa6fc206c80e','algorithm','HS512'),('da6f6814-816a-44b2-918f-4e35d04e22fc','889c6ce2-796a-4817-9f2e-fd7e229e0c79','client-uris-must-match','true'),('dc2628ec-02a1-447e-9601-748d52f3543a','319d59ac-77ea-4712-ac2c-21ff6f237fee','priority','100'),('e2a5fd09-49bc-46f3-93f6-3b0d35569bf2','d977b6b2-0bad-4ac5-94ca-82c5644be1de','allowed-protocol-mapper-types','saml-role-list-mapper'),('e376eb1a-7c03-4e9a-a50c-c6fed3847bd8','e77544a7-70a7-4a7a-b724-3b8f39b19c47','priority','100'),('e4649341-66bf-426f-900a-2ca046b27882','ebc211ab-256f-471a-bfac-798284a67310','host-sending-registration-request-must-match','true'),('e4bf6517-fdca-47b3-90c2-c1eed35fee29','ce4a5c77-367f-4988-a090-bada4848688e','allowed-protocol-mapper-types','oidc-full-name-mapper'),('e81c27f0-a4de-41ad-96a4-a3b0d738f2d1','da8ab1c7-fb41-4121-bff6-9b43ff563220','allow-default-scopes','true'),('eb7db513-c383-4540-b5f2-e42b52925047','76f0e063-b83f-4308-bd67-af4985f83991','allowed-protocol-mapper-types','oidc-usermodel-property-mapper'),('ee95a89b-6ea3-4e3b-b291-a62a657cdf8b','272acfab-2609-4470-aa38-03608ae64547','algorithm','HS512'),('f06f5a61-711a-45d4-bb99-7445fca38bc7','4c39c7c3-ef4a-4d77-94b6-e29eae7d79f2','priority','100'),('f39d6524-4e28-477a-8a73-c39adb57f3d6','e223336a-ff9f-4616-bca2-2519d6e197b5','priority','100'),('f7829a4e-ac58-4b4a-82a9-fa6b5475af01','bbf18037-cd03-4297-b42f-aa6fc206c80e','priority','100'),('f90599d3-19d8-4ee2-9592-ddecac26cc6b','4c39c7c3-ef4a-4d77-94b6-e29eae7d79f2','kid','8d29ec4c-7f62-4ca0-9c47-12b6df98a42e'),('fa7b890e-d11f-4baf-a2f1-7beecd5cc5ba','76f0e063-b83f-4308-bd67-af4985f83991','allowed-protocol-mapper-types','oidc-usermodel-attribute-mapper'),('fbde35a5-9c37-4e57-a900-c0c91516d893','f3839f85-b917-4a4e-bab6-1424a039329f','allow-default-scopes','true');
/*!40000 ALTER TABLE `COMPONENT_CONFIG` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `COMPOSITE_ROLE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `COMPOSITE_ROLE` (
  `COMPOSITE` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `CHILD_ROLE` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`COMPOSITE`,`CHILD_ROLE`),
  KEY `IDX_COMPOSITE` (`COMPOSITE`),
  KEY `IDX_COMPOSITE_CHILD` (`CHILD_ROLE`),
  CONSTRAINT `FK_A63WVEKFTU8JO1PNJ81E7MCE2` FOREIGN KEY (`COMPOSITE`) REFERENCES `KEYCLOAK_ROLE` (`ID`),
  CONSTRAINT `FK_GR7THLLB9LU8Q4VQA4524JJY8` FOREIGN KEY (`CHILD_ROLE`) REFERENCES `KEYCLOAK_ROLE` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `COMPOSITE_ROLE` WRITE;
/*!40000 ALTER TABLE `COMPOSITE_ROLE` DISABLE KEYS */;
INSERT INTO `COMPOSITE_ROLE` VALUES ('0c50da1d-eeb7-4074-b31e-f638245df45e','737be905-164a-4495-980f-d6f1efeeb69a'),('0c50da1d-eeb7-4074-b31e-f638245df45e','a8ded98f-bc0d-490e-99d6-320855772b20'),('0c50da1d-eeb7-4074-b31e-f638245df45e','afdaca80-10d0-4be8-b18c-4933151fb78c'),('0c50da1d-eeb7-4074-b31e-f638245df45e','ba23c820-628a-4fd3-aecc-07491f156fe2'),('251d9d2a-f9fc-4ca1-ac83-52a391885e02','a25dee31-81ff-4364-b45d-2a9799a7d6b9'),('5125458c-9476-41c9-8fb6-0543f7083773','49d17406-d5b7-4b46-801c-7aad314461cd'),('5125458c-9476-41c9-8fb6-0543f7083773','7a94500c-c30a-42ed-ae1a-19bc77834854'),('5125458c-9476-41c9-8fb6-0543f7083773','bbbc95f8-5d80-4426-9b22-11e1b5c6db03'),('5125458c-9476-41c9-8fb6-0543f7083773','bc732c7b-1dda-4fcb-bbfd-d2c11138eebf'),('59f35481-657d-4cc2-9fb2-6ff586c5bcb1','418c2f98-e334-4a20-aeb7-3cb62f58f764'),('78eb646a-3134-4823-b2ce-daf30bbfcf9b','326e598f-eaae-4da1-acbf-daafb4789421'),('807232af-7777-454a-b1cf-82a74b4ddf4b','08e2bffc-779d-4dc5-b4fc-2b9c27ddf898'),('807232af-7777-454a-b1cf-82a74b4ddf4b','1302a0c9-cff5-4f38-a524-387fb61ecae5'),('807232af-7777-454a-b1cf-82a74b4ddf4b','1438fd35-067c-4657-a9ae-99394fcdddf4'),('807232af-7777-454a-b1cf-82a74b4ddf4b','14ba4893-70c6-4be5-9beb-225ca62ed5ec'),('807232af-7777-454a-b1cf-82a74b4ddf4b','31f48475-6b51-42b4-8706-db2d938ed265'),('807232af-7777-454a-b1cf-82a74b4ddf4b','326e598f-eaae-4da1-acbf-daafb4789421'),('807232af-7777-454a-b1cf-82a74b4ddf4b','3dd13d47-55dd-4ac7-8e2a-c428fafc3739'),('807232af-7777-454a-b1cf-82a74b4ddf4b','5867419e-37e5-4166-ac97-92b7c53e6eb9'),('807232af-7777-454a-b1cf-82a74b4ddf4b','78eb646a-3134-4823-b2ce-daf30bbfcf9b'),('807232af-7777-454a-b1cf-82a74b4ddf4b','8868fe7b-aadb-4011-a14b-dea7ef1df1e1'),('807232af-7777-454a-b1cf-82a74b4ddf4b','9a839442-0a1f-42d6-8e91-2b8c6996e199'),('807232af-7777-454a-b1cf-82a74b4ddf4b','b2a6016f-8d43-4ab7-be9d-7119c2b71b80'),('807232af-7777-454a-b1cf-82a74b4ddf4b','b723d2ac-460d-49a1-abbb-077d65b47e98'),('807232af-7777-454a-b1cf-82a74b4ddf4b','e7b39c0b-fef7-44a1-b209-df633fcbc8a5'),('807232af-7777-454a-b1cf-82a74b4ddf4b','eacb6d31-9ae2-4618-867c-2165b879f81d'),('807232af-7777-454a-b1cf-82a74b4ddf4b','ec508af1-4bcb-48f3-95af-d8e3f7efe47c'),('807232af-7777-454a-b1cf-82a74b4ddf4b','f1b75db6-0157-406f-ade0-bc90cf859fda'),('807232af-7777-454a-b1cf-82a74b4ddf4b','ff7b90a1-7585-44e1-89c5-0a39963f236d'),('87792728-d289-4c80-9c2f-c4568b5d8801','4bac0b9c-8ffe-491a-b070-a7bad02d88e4'),('87792728-d289-4c80-9c2f-c4568b5d8801','75e66f19-ba72-4c66-a699-21e5c35720b5'),('95fe56c8-1533-4dcc-9f8f-c0bd67254df3','38c40b55-9be6-435d-9e4a-595377c1fa77'),('95fe56c8-1533-4dcc-9f8f-c0bd67254df3','a8bd6cb5-d0ff-4a0a-8cf8-eac4f15a541e'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','00d5c5a8-f8df-41fd-bc8d-b9ef9a93d3e6'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','04291111-9b7e-473f-b432-9a3cd9a06c33'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','118607ab-1f22-4fb8-8b89-243907ef7f5d'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','1245baa1-e64c-4cba-966f-0a8069172814'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','1fb46915-dcb5-4305-b4d6-67798ae7498f'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','21082b5c-dd66-4536-b23b-c8678d7d60d5'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','251d9d2a-f9fc-4ca1-ac83-52a391885e02'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','26f9df1c-293a-4cf2-932a-c3a709cd1185'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','2b35e642-e0c2-4f86-a716-4d5f546d800d'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','2e7701b8-55a9-4fd2-8bc3-3d77c8249ce0'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','2fbe9598-8e10-4bba-ba61-899a02a01877'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','301e1fbd-e2dc-4aaa-a133-6ac79d60317d'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','35f07854-7a7a-4408-aca4-fad5c10cb1c6'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','38c40b55-9be6-435d-9e4a-595377c1fa77'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','418c2f98-e334-4a20-aeb7-3cb62f58f764'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','4729d7e4-7b7f-4927-8af1-e74e487a1357'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','4bac0b9c-8ffe-491a-b070-a7bad02d88e4'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','4c83be7d-a5c9-4b31-98e7-95e0320b16a5'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','4fd8b333-076a-4b26-91df-5fdc3516511b'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','542a60ff-97e0-4127-bfcf-1fd2b102251d'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','5634fdb9-5681-4755-9efe-ed28a3559530'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','59f35481-657d-4cc2-9fb2-6ff586c5bcb1'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','5e86f90e-ac50-4faa-88d7-d54cb32eb4eb'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','64516292-6475-48b1-a156-399de0901f86'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','6dda702d-8bae-46a3-8308-97b6c5a82fe1'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','71cbf586-9211-43be-af4e-3f49778b8300'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','726fa1e9-a766-47ec-8d2e-0ff9c38fcd51'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','737be905-164a-4495-980f-d6f1efeeb69a'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','75e66f19-ba72-4c66-a699-21e5c35720b5'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','7a5e73a4-f0c7-40dc-bb8f-01055b209de5'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','7be7e323-385b-4bce-a3cc-95eb5a1b22b2'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','87792728-d289-4c80-9c2f-c4568b5d8801'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','87cd1ca9-fabd-463a-9f38-9c7aab1d8694'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','95fe56c8-1533-4dcc-9f8f-c0bd67254df3'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','a25dee31-81ff-4364-b45d-2a9799a7d6b9'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','a8bd6cb5-d0ff-4a0a-8cf8-eac4f15a541e'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','a8ded98f-bc0d-490e-99d6-320855772b20'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','c0412fa5-7c3f-4928-80bd-74c26fb16dea'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','c08aac22-9f0a-43d1-a2f9-8ae0e672828e'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','caf491f4-fba2-4246-99e4-2d3f85248e64'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','cb13b85a-805b-4c06-bbc9-5548b3737bef'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','d2fc3a79-57f8-4447-8a23-297a4770240f'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','daae15c7-3895-4aff-b1c9-3e753b44a99a'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','f8772141-c42e-4f39-b42c-6652dd8a5b70'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','f89122ff-d079-4f7c-b403-4d20643f6649'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','f9a83f1b-cc89-4a2c-a873-0055fe745268'),('a3fa1093-7cb2-444c-82f9-dc04cf28ca42','098e5980-f81f-4f7a-8eae-685dec71faa4'),('a8ded98f-bc0d-490e-99d6-320855772b20','7a5e73a4-f0c7-40dc-bb8f-01055b209de5'),('b723d2ac-460d-49a1-abbb-077d65b47e98','1438fd35-067c-4657-a9ae-99394fcdddf4'),('b723d2ac-460d-49a1-abbb-077d65b47e98','5867419e-37e5-4166-ac97-92b7c53e6eb9'),('bc732c7b-1dda-4fcb-bbfd-d2c11138eebf','4fba4722-371b-47c9-a0ca-3a2b05b0f197'),('caf491f4-fba2-4246-99e4-2d3f85248e64','d2fc3a79-57f8-4447-8a23-297a4770240f');
/*!40000 ALTER TABLE `COMPOSITE_ROLE` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `CREDENTIAL`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `CREDENTIAL` (
  `ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `SALT` tinyblob,
  `TYPE` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `USER_ID` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `CREATED_DATE` bigint DEFAULT NULL,
  `USER_LABEL` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `SECRET_DATA` longtext COLLATE utf8mb4_unicode_ci,
  `CREDENTIAL_DATA` longtext COLLATE utf8mb4_unicode_ci,
  `PRIORITY` int DEFAULT NULL,
  `VERSION` int DEFAULT '0',
  PRIMARY KEY (`ID`),
  KEY `IDX_USER_CREDENTIAL` (`USER_ID`),
  CONSTRAINT `FK_PFYR0GLASQYL0DEI3KL69R6V0` FOREIGN KEY (`USER_ID`) REFERENCES `USER_ENTITY` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `CREDENTIAL` WRITE;
/*!40000 ALTER TABLE `CREDENTIAL` DISABLE KEYS */;
INSERT INTO `CREDENTIAL` VALUES ('82886297-2430-43fb-8cb0-bc5ea5f72419',NULL,'password','19b6d072-21e8-40f9-9d22-711e5814459f',1765993612454,NULL,'{\"value\":\"8fQgPbDATkm6Cl4iIP8+swWDmCB6ToZ4XkD87N3ZuHE=\",\"salt\":\"2lezhqia29SGT015lT9ECA==\",\"additionalParameters\":{}}','{\"hashIterations\":5,\"algorithm\":\"argon2\",\"additionalParameters\":{\"hashLength\":[\"32\"],\"memory\":[\"7168\"],\"type\":[\"id\"],\"version\":[\"1.3\"],\"parallelism\":[\"1\"]}}',10,0),('96514475-73f2-46c1-a453-16fede0987ba',NULL,'password','d62612e2-e122-441d-8e71-7fee81f4c53c',1768842442111,'My password','{\"value\":\"UND9CTANI1w8nR/0gfDYk0Uz7DInoHQihBx506690fk=\",\"salt\":\"OonxLO7zjsiB8biqSn/bDA==\",\"additionalParameters\":{}}','{\"hashIterations\":5,\"algorithm\":\"argon2\",\"additionalParameters\":{\"hashLength\":[\"32\"],\"memory\":[\"7168\"],\"type\":[\"id\"],\"version\":[\"1.3\"],\"parallelism\":[\"1\"]}}',10,1),('967f4592-f55d-4932-ae6e-b00d7ca1e6ba',NULL,'password','79c7eea4-7f93-4b1e-ad06-b244ba9ca466',1773366116962,NULL,'{\"value\":\"0FN8f/dA0UKw95p0ySvmtAKCAZGr4beg5Ig2u07zGTs=\",\"salt\":\"ADySFRZQ6VEB1Zj5L6IPpA==\",\"additionalParameters\":{}}','{\"hashIterations\":5,\"algorithm\":\"argon2\",\"additionalParameters\":{\"hashLength\":[\"32\"],\"memory\":[\"7168\"],\"type\":[\"id\"],\"version\":[\"1.3\"],\"parallelism\":[\"1\"]}}',10,0),('b6e41ac4-bf0a-43c3-beab-bcfb75b839d4',NULL,'password','fb214e14-d12c-451e-aca8-c62b09f77b11',1773366451963,NULL,'{\"value\":\"reMdiaDOVgTFXcTiCwCNcW1Fum61fImi+P77/gFLd4g=\",\"salt\":\"nSyhZ5bz/oDMatjFup9Xjw==\",\"additionalParameters\":{}}','{\"hashIterations\":5,\"algorithm\":\"argon2\",\"additionalParameters\":{\"hashLength\":[\"32\"],\"memory\":[\"7168\"],\"type\":[\"id\"],\"version\":[\"1.3\"],\"parallelism\":[\"1\"]}}',10,0);
/*!40000 ALTER TABLE `CREDENTIAL` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `DATABASECHANGELOG`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `DATABASECHANGELOG` (
  `ID` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `AUTHOR` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `FILENAME` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `DATEEXECUTED` datetime NOT NULL,
  `ORDEREXECUTED` int NOT NULL,
  `EXECTYPE` varchar(10) COLLATE utf8mb4_unicode_ci NOT NULL,
  `MD5SUM` varchar(35) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `DESCRIPTION` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `COMMENTS` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `TAG` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `LIQUIBASE` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `CONTEXTS` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `LABELS` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `DEPLOYMENT_ID` varchar(10) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`ID`,`AUTHOR`,`FILENAME`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `DATABASECHANGELOG` WRITE;
/*!40000 ALTER TABLE `DATABASECHANGELOG` DISABLE KEYS */;
INSERT INTO `DATABASECHANGELOG` VALUES ('1.0.0.Final-KEYCLOAK-5461','sthorger@redhat.com','META-INF/db2-jpa-changelog-1.0.0.Final.xml','2025-12-17 12:46:34',2,'MARK_RAN','9:828775b1596a07d1200ba1d49e5e3941','createTable tableName=APPLICATION_DEFAULT_ROLES; createTable tableName=CLIENT; createTable tableName=CLIENT_SESSION; createTable tableName=CLIENT_SESSION_ROLE; createTable tableName=COMPOSITE_ROLE; createTable tableName=CREDENTIAL; createTable tab...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('1.0.0.Final-KEYCLOAK-5461','sthorger@redhat.com','META-INF/jpa-changelog-1.0.0.Final.xml','2025-12-17 12:46:34',1,'EXECUTED','9:6f1016664e21e16d26517a4418f5e3df','createTable tableName=APPLICATION_DEFAULT_ROLES; createTable tableName=CLIENT; createTable tableName=CLIENT_SESSION; createTable tableName=CLIENT_SESSION_ROLE; createTable tableName=COMPOSITE_ROLE; createTable tableName=CREDENTIAL; createTable tab...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('1.1.0.Beta1','sthorger@redhat.com','META-INF/jpa-changelog-1.1.0.Beta1.xml','2025-12-17 12:46:35',3,'EXECUTED','9:5f090e44a7d595883c1fb61f4b41fd38','delete tableName=CLIENT_SESSION_ROLE; delete tableName=CLIENT_SESSION; delete tableName=USER_SESSION; createTable tableName=CLIENT_ATTRIBUTES; createTable tableName=CLIENT_SESSION_NOTE; createTable tableName=APP_NODE_REGISTRATIONS; addColumn table...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('1.1.0.Final','sthorger@redhat.com','META-INF/jpa-changelog-1.1.0.Final.xml','2025-12-17 12:46:35',4,'EXECUTED','9:c07e577387a3d2c04d1adc9aaad8730e','renameColumn newColumnName=EVENT_TIME, oldColumnName=TIME, tableName=EVENT_ENTITY','',NULL,'4.33.0',NULL,NULL,'5993589629'),('1.2.0.Beta1','psilva@redhat.com','META-INF/db2-jpa-changelog-1.2.0.Beta1.xml','2025-12-17 12:46:35',6,'MARK_RAN','9:543b5c9989f024fe35c6f6c5a97de88e','delete tableName=CLIENT_SESSION_ROLE; delete tableName=CLIENT_SESSION_NOTE; delete tableName=CLIENT_SESSION; delete tableName=USER_SESSION; createTable tableName=PROTOCOL_MAPPER; createTable tableName=PROTOCOL_MAPPER_CONFIG; createTable tableName=...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('1.2.0.Beta1','psilva@redhat.com','META-INF/jpa-changelog-1.2.0.Beta1.xml','2025-12-17 12:46:35',5,'EXECUTED','9:b68ce996c655922dbcd2fe6b6ae72686','delete tableName=CLIENT_SESSION_ROLE; delete tableName=CLIENT_SESSION_NOTE; delete tableName=CLIENT_SESSION; delete tableName=USER_SESSION; createTable tableName=PROTOCOL_MAPPER; createTable tableName=PROTOCOL_MAPPER_CONFIG; createTable tableName=...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('1.2.0.Final','keycloak','META-INF/jpa-changelog-1.2.0.Final.xml','2025-12-17 12:46:36',9,'EXECUTED','9:9d05c7be10cdb873f8bcb41bc3a8ab23','update tableName=CLIENT; update tableName=CLIENT; update tableName=CLIENT','',NULL,'4.33.0',NULL,NULL,'5993589629'),('1.2.0.RC1','bburke@redhat.com','META-INF/db2-jpa-changelog-1.2.0.CR1.xml','2025-12-17 12:46:36',8,'MARK_RAN','9:db4a145ba11a6fdaefb397f6dbf829a1','delete tableName=CLIENT_SESSION_ROLE; delete tableName=CLIENT_SESSION_NOTE; delete tableName=CLIENT_SESSION; delete tableName=USER_SESSION_NOTE; delete tableName=USER_SESSION; createTable tableName=MIGRATION_MODEL; createTable tableName=IDENTITY_P...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('1.2.0.RC1','bburke@redhat.com','META-INF/jpa-changelog-1.2.0.CR1.xml','2025-12-17 12:46:36',7,'EXECUTED','9:765afebbe21cf5bbca048e632df38336','delete tableName=CLIENT_SESSION_ROLE; delete tableName=CLIENT_SESSION_NOTE; delete tableName=CLIENT_SESSION; delete tableName=USER_SESSION_NOTE; delete tableName=USER_SESSION; createTable tableName=MIGRATION_MODEL; createTable tableName=IDENTITY_P...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('1.3.0','bburke@redhat.com','META-INF/jpa-changelog-1.3.0.xml','2025-12-17 12:46:37',10,'EXECUTED','9:18593702353128d53111f9b1ff0b82b8','delete tableName=CLIENT_SESSION_ROLE; delete tableName=CLIENT_SESSION_PROT_MAPPER; delete tableName=CLIENT_SESSION_NOTE; delete tableName=CLIENT_SESSION; delete tableName=USER_SESSION_NOTE; delete tableName=USER_SESSION; createTable tableName=ADMI...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('1.4.0','bburke@redhat.com','META-INF/db2-jpa-changelog-1.4.0.xml','2025-12-17 12:46:37',12,'MARK_RAN','9:e1ff28bf7568451453f844c5d54bb0b5','delete tableName=CLIENT_SESSION_AUTH_STATUS; delete tableName=CLIENT_SESSION_ROLE; delete tableName=CLIENT_SESSION_PROT_MAPPER; delete tableName=CLIENT_SESSION_NOTE; delete tableName=CLIENT_SESSION; delete tableName=USER_SESSION_NOTE; delete table...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('1.4.0','bburke@redhat.com','META-INF/jpa-changelog-1.4.0.xml','2025-12-17 12:46:37',11,'EXECUTED','9:6122efe5f090e41a85c0f1c9e52cbb62','delete tableName=CLIENT_SESSION_AUTH_STATUS; delete tableName=CLIENT_SESSION_ROLE; delete tableName=CLIENT_SESSION_PROT_MAPPER; delete tableName=CLIENT_SESSION_NOTE; delete tableName=CLIENT_SESSION; delete tableName=USER_SESSION_NOTE; delete table...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('1.5.0','bburke@redhat.com','META-INF/jpa-changelog-1.5.0.xml','2025-12-17 12:46:37',13,'EXECUTED','9:7af32cd8957fbc069f796b61217483fd','delete tableName=CLIENT_SESSION_AUTH_STATUS; delete tableName=CLIENT_SESSION_ROLE; delete tableName=CLIENT_SESSION_PROT_MAPPER; delete tableName=CLIENT_SESSION_NOTE; delete tableName=CLIENT_SESSION; delete tableName=USER_SESSION_NOTE; delete table...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('1.6.1','mposolda@redhat.com','META-INF/jpa-changelog-1.6.1.xml','2025-12-17 12:46:37',17,'EXECUTED','9:d41d8cd98f00b204e9800998ecf8427e','empty','',NULL,'4.33.0',NULL,NULL,'5993589629'),('1.6.1_from15','mposolda@redhat.com','META-INF/jpa-changelog-1.6.1.xml','2025-12-17 12:46:37',14,'EXECUTED','9:6005e15e84714cd83226bf7879f54190','addColumn tableName=REALM; addColumn tableName=KEYCLOAK_ROLE; addColumn tableName=CLIENT; createTable tableName=OFFLINE_USER_SESSION; createTable tableName=OFFLINE_CLIENT_SESSION; addPrimaryKey constraintName=CONSTRAINT_OFFL_US_SES_PK2, tableName=...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('1.6.1_from16','mposolda@redhat.com','META-INF/jpa-changelog-1.6.1.xml','2025-12-17 12:46:37',16,'MARK_RAN','9:f8dadc9284440469dcf71e25ca6ab99b','dropPrimaryKey constraintName=CONSTRAINT_OFFLINE_US_SES_PK, tableName=OFFLINE_USER_SESSION; dropPrimaryKey constraintName=CONSTRAINT_OFFLINE_CL_SES_PK, tableName=OFFLINE_CLIENT_SESSION; addColumn tableName=OFFLINE_USER_SESSION; update tableName=OF...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('1.6.1_from16-pre','mposolda@redhat.com','META-INF/jpa-changelog-1.6.1.xml','2025-12-17 12:46:37',15,'MARK_RAN','9:bf656f5a2b055d07f314431cae76f06c','delete tableName=OFFLINE_CLIENT_SESSION; delete tableName=OFFLINE_USER_SESSION','',NULL,'4.33.0',NULL,NULL,'5993589629'),('1.7.0','bburke@redhat.com','META-INF/jpa-changelog-1.7.0.xml','2025-12-17 12:46:38',18,'EXECUTED','9:3368ff0be4c2855ee2dd9ca813b38d8e','createTable tableName=KEYCLOAK_GROUP; createTable tableName=GROUP_ROLE_MAPPING; createTable tableName=GROUP_ATTRIBUTE; createTable tableName=USER_GROUP_MEMBERSHIP; createTable tableName=REALM_DEFAULT_GROUPS; addColumn tableName=IDENTITY_PROVIDER; ...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('1.8.0','mposolda@redhat.com','META-INF/db2-jpa-changelog-1.8.0.xml','2025-12-17 12:46:38',21,'MARK_RAN','9:831e82914316dc8a57dc09d755f23c51','addColumn tableName=IDENTITY_PROVIDER; createTable tableName=CLIENT_TEMPLATE; createTable tableName=CLIENT_TEMPLATE_ATTRIBUTES; createTable tableName=TEMPLATE_SCOPE_MAPPING; dropNotNullConstraint columnName=CLIENT_ID, tableName=PROTOCOL_MAPPER; ad...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('1.8.0','mposolda@redhat.com','META-INF/jpa-changelog-1.8.0.xml','2025-12-17 12:46:38',19,'EXECUTED','9:8ac2fb5dd030b24c0570a763ed75ed20','addColumn tableName=IDENTITY_PROVIDER; createTable tableName=CLIENT_TEMPLATE; createTable tableName=CLIENT_TEMPLATE_ATTRIBUTES; createTable tableName=TEMPLATE_SCOPE_MAPPING; dropNotNullConstraint columnName=CLIENT_ID, tableName=PROTOCOL_MAPPER; ad...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('1.8.0-2','keycloak','META-INF/db2-jpa-changelog-1.8.0.xml','2025-12-17 12:46:38',22,'MARK_RAN','9:f91ddca9b19743db60e3057679810e6c','dropDefaultValue columnName=ALGORITHM, tableName=CREDENTIAL; update tableName=CREDENTIAL','',NULL,'4.33.0',NULL,NULL,'5993589629'),('1.8.0-2','keycloak','META-INF/jpa-changelog-1.8.0.xml','2025-12-17 12:46:38',20,'EXECUTED','9:f91ddca9b19743db60e3057679810e6c','dropDefaultValue columnName=ALGORITHM, tableName=CREDENTIAL; update tableName=CREDENTIAL','',NULL,'4.33.0',NULL,NULL,'5993589629'),('1.9.0','mposolda@redhat.com','META-INF/jpa-changelog-1.9.0.xml','2025-12-17 12:46:38',23,'EXECUTED','9:bc3d0f9e823a69dc21e23e94c7a94bb1','update tableName=REALM; update tableName=REALM; update tableName=REALM; update tableName=REALM; update tableName=CREDENTIAL; update tableName=CREDENTIAL; update tableName=CREDENTIAL; update tableName=REALM; update tableName=REALM; customChange; dr...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('1.9.1','keycloak','META-INF/db2-jpa-changelog-1.9.1.xml','2025-12-17 12:46:38',25,'MARK_RAN','9:0d6c65c6f58732d81569e77b10ba301d','modifyDataType columnName=PRIVATE_KEY, tableName=REALM; modifyDataType columnName=CERTIFICATE, tableName=REALM','',NULL,'4.33.0',NULL,NULL,'5993589629'),('1.9.1','keycloak','META-INF/jpa-changelog-1.9.1.xml','2025-12-17 12:46:38',24,'EXECUTED','9:c9999da42f543575ab790e76439a2679','modifyDataType columnName=PRIVATE_KEY, tableName=REALM; modifyDataType columnName=PUBLIC_KEY, tableName=REALM; modifyDataType columnName=CERTIFICATE, tableName=REALM','',NULL,'4.33.0',NULL,NULL,'5993589629'),('1.9.2','keycloak','META-INF/jpa-changelog-1.9.2.xml','2025-12-17 12:46:38',26,'EXECUTED','9:fc576660fc016ae53d2d4778d84d86d0','createIndex indexName=IDX_USER_EMAIL, tableName=USER_ENTITY; createIndex indexName=IDX_USER_ROLE_MAPPING, tableName=USER_ROLE_MAPPING; createIndex indexName=IDX_USER_GROUP_MAPPING, tableName=USER_GROUP_MEMBERSHIP; createIndex indexName=IDX_USER_CO...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('12.1.0-add-realm-localization-table','keycloak','META-INF/jpa-changelog-12.0.0.xml','2025-12-17 12:46:44',88,'EXECUTED','9:fffabce2bc01e1a8f5110d5278500065','createTable tableName=REALM_LOCALIZATIONS; addPrimaryKey tableName=REALM_LOCALIZATIONS','',NULL,'4.33.0',NULL,NULL,'5993589629'),('13.0.0-increase-column-size-federated','keycloak','META-INF/jpa-changelog-13.0.0.xml','2025-12-17 12:46:45',94,'EXECUTED','9:43c0c1055b6761b4b3e89de76d612ccf','modifyDataType columnName=CLIENT_ID, tableName=CLIENT_SCOPE_CLIENT; modifyDataType columnName=SCOPE_ID, tableName=CLIENT_SCOPE_CLIENT','',NULL,'4.33.0',NULL,NULL,'5993589629'),('13.0.0-KEYCLOAK-16844','keycloak','META-INF/jpa-changelog-13.0.0.xml','2025-12-17 12:46:44',91,'EXECUTED','9:ad1194d66c937e3ffc82386c050ba089','createIndex indexName=IDX_OFFLINE_USS_PRELOAD, tableName=OFFLINE_USER_SESSION','',NULL,'4.33.0',NULL,NULL,'5993589629'),('13.0.0-KEYCLOAK-17992-drop-constraints','keycloak','META-INF/jpa-changelog-13.0.0.xml','2025-12-17 12:46:45',93,'MARK_RAN','9:544d201116a0fcc5a5da0925fbbc3bde','dropPrimaryKey constraintName=C_CLI_SCOPE_BIND, tableName=CLIENT_SCOPE_CLIENT; dropIndex indexName=IDX_CLSCOPE_CL, tableName=CLIENT_SCOPE_CLIENT; dropIndex indexName=IDX_CL_CLSCOPE, tableName=CLIENT_SCOPE_CLIENT','',NULL,'4.33.0',NULL,NULL,'5993589629'),('13.0.0-KEYCLOAK-17992-recreate-constraints','keycloak','META-INF/jpa-changelog-13.0.0.xml','2025-12-17 12:46:45',95,'MARK_RAN','9:8bd711fd0330f4fe980494ca43ab1139','addNotNullConstraint columnName=CLIENT_ID, tableName=CLIENT_SCOPE_CLIENT; addNotNullConstraint columnName=SCOPE_ID, tableName=CLIENT_SCOPE_CLIENT; addPrimaryKey constraintName=C_CLI_SCOPE_BIND, tableName=CLIENT_SCOPE_CLIENT; createIndex indexName=...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('14.0.0-KEYCLOAK-11019','keycloak','META-INF/jpa-changelog-14.0.0.xml','2025-12-17 12:46:45',97,'EXECUTED','9:24fb8611e97f29989bea412aa38d12b7','createIndex indexName=IDX_OFFLINE_CSS_PRELOAD, tableName=OFFLINE_CLIENT_SESSION; createIndex indexName=IDX_OFFLINE_USS_BY_USER, tableName=OFFLINE_USER_SESSION; createIndex indexName=IDX_OFFLINE_USS_BY_USERSESS, tableName=OFFLINE_USER_SESSION','',NULL,'4.33.0',NULL,NULL,'5993589629'),('14.0.0-KEYCLOAK-18286','keycloak','META-INF/jpa-changelog-14.0.0.xml','2025-12-17 12:46:45',98,'MARK_RAN','9:259f89014ce2506ee84740cbf7163aa7','createIndex indexName=IDX_CLIENT_ATT_BY_NAME_VALUE, tableName=CLIENT_ATTRIBUTES','',NULL,'4.33.0',NULL,NULL,'5993589629'),('14.0.0-KEYCLOAK-18286-revert','keycloak','META-INF/jpa-changelog-14.0.0.xml','2025-12-17 12:46:45',99,'MARK_RAN','9:04baaf56c116ed19951cbc2cca584022','dropIndex indexName=IDX_CLIENT_ATT_BY_NAME_VALUE, tableName=CLIENT_ATTRIBUTES','',NULL,'4.33.0',NULL,NULL,'5993589629'),('14.0.0-KEYCLOAK-18286-supported-dbs','keycloak','META-INF/jpa-changelog-14.0.0.xml','2025-12-17 12:46:45',100,'EXECUTED','9:bd2bd0fc7768cf0845ac96a8786fa735','createIndex indexName=IDX_CLIENT_ATT_BY_NAME_VALUE, tableName=CLIENT_ATTRIBUTES','',NULL,'4.33.0',NULL,NULL,'5993589629'),('14.0.0-KEYCLOAK-18286-unsupported-dbs','keycloak','META-INF/jpa-changelog-14.0.0.xml','2025-12-17 12:46:45',101,'MARK_RAN','9:d3d977031d431db16e2c181ce49d73e9','createIndex indexName=IDX_CLIENT_ATT_BY_NAME_VALUE, tableName=CLIENT_ATTRIBUTES','',NULL,'4.33.0',NULL,NULL,'5993589629'),('15.0.0-KEYCLOAK-18467','keycloak','META-INF/jpa-changelog-15.0.0.xml','2025-12-17 12:46:45',104,'EXECUTED','9:47a760639ac597360a8219f5b768b4de','addColumn tableName=REALM_LOCALIZATIONS; update tableName=REALM_LOCALIZATIONS; dropColumn columnName=TEXTS, tableName=REALM_LOCALIZATIONS; renameColumn newColumnName=TEXTS, oldColumnName=TEXTS_NEW, tableName=REALM_LOCALIZATIONS; addNotNullConstrai...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('17.0.0-9562','keycloak','META-INF/jpa-changelog-17.0.0.xml','2025-12-17 12:46:45',105,'EXECUTED','9:a6272f0576727dd8cad2522335f5d99e','createIndex indexName=IDX_USER_SERVICE_ACCOUNT, tableName=USER_ENTITY','',NULL,'4.33.0',NULL,NULL,'5993589629'),('18.0.0-10625-IDX_ADMIN_EVENT_TIME','keycloak','META-INF/jpa-changelog-18.0.0.xml','2025-12-17 12:46:45',106,'EXECUTED','9:015479dbd691d9cc8669282f4828c41d','createIndex indexName=IDX_ADMIN_EVENT_TIME, tableName=ADMIN_EVENT_ENTITY','',NULL,'4.33.0',NULL,NULL,'5993589629'),('18.0.15-30992-index-consent','keycloak','META-INF/jpa-changelog-18.0.15.xml','2025-12-17 12:46:45',107,'EXECUTED','9:80071ede7a05604b1f4906f3bf3b00f0','createIndex indexName=IDX_USCONSENT_SCOPE_ID, tableName=USER_CONSENT_CLIENT_SCOPE','',NULL,'4.33.0',NULL,NULL,'5993589629'),('19.0.0-10135','keycloak','META-INF/jpa-changelog-19.0.0.xml','2025-12-17 12:46:45',108,'EXECUTED','9:9518e495fdd22f78ad6425cc30630221','customChange','',NULL,'4.33.0',NULL,NULL,'5993589629'),('2.1.0-KEYCLOAK-5461','bburke@redhat.com','META-INF/jpa-changelog-2.1.0.xml','2025-12-17 12:46:39',29,'EXECUTED','9:bd88e1f833df0420b01e114533aee5e8','createTable tableName=BROKER_LINK; createTable tableName=FED_USER_ATTRIBUTE; createTable tableName=FED_USER_CONSENT; createTable tableName=FED_USER_CONSENT_ROLE; createTable tableName=FED_USER_CONSENT_PROT_MAPPER; createTable tableName=FED_USER_CR...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('2.2.0','bburke@redhat.com','META-INF/jpa-changelog-2.2.0.xml','2025-12-17 12:46:39',30,'EXECUTED','9:a7022af5267f019d020edfe316ef4371','addColumn tableName=ADMIN_EVENT_ENTITY; createTable tableName=CREDENTIAL_ATTRIBUTE; createTable tableName=FED_CREDENTIAL_ATTRIBUTE; modifyDataType columnName=VALUE, tableName=CREDENTIAL; addForeignKeyConstraint baseTableName=FED_CREDENTIAL_ATTRIBU...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('2.3.0','bburke@redhat.com','META-INF/jpa-changelog-2.3.0.xml','2025-12-17 12:46:39',31,'EXECUTED','9:fc155c394040654d6a79227e56f5e25a','createTable tableName=FEDERATED_USER; addPrimaryKey constraintName=CONSTR_FEDERATED_USER, tableName=FEDERATED_USER; dropDefaultValue columnName=TOTP, tableName=USER_ENTITY; dropColumn columnName=TOTP, tableName=USER_ENTITY; addColumn tableName=IDE...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('2.4.0','bburke@redhat.com','META-INF/jpa-changelog-2.4.0.xml','2025-12-17 12:46:39',32,'EXECUTED','9:eac4ffb2a14795e5dc7b426063e54d88','customChange','',NULL,'4.33.0',NULL,NULL,'5993589629'),('2.5.0','bburke@redhat.com','META-INF/jpa-changelog-2.5.0.xml','2025-12-17 12:46:39',33,'EXECUTED','9:54937c05672568c4c64fc9524c1e9462','customChange; modifyDataType columnName=USER_ID, tableName=OFFLINE_USER_SESSION','',NULL,'4.33.0',NULL,NULL,'5993589629'),('2.5.0-duplicate-email-support','slawomir@dabek.name','META-INF/jpa-changelog-2.5.0.xml','2025-12-17 12:46:40',36,'EXECUTED','9:61b6d3d7a4c0e0024b0c839da283da0c','addColumn tableName=REALM','',NULL,'4.33.0',NULL,NULL,'5993589629'),('2.5.0-unicode-oracle','hmlnarik@redhat.com','META-INF/jpa-changelog-2.5.0.xml','2025-12-17 12:46:39',34,'MARK_RAN','9:79a309a7f4ddc2f8db9b23e46010152d','modifyDataType columnName=DESCRIPTION, tableName=AUTHENTICATION_FLOW; modifyDataType columnName=DESCRIPTION, tableName=CLIENT_TEMPLATE; modifyDataType columnName=DESCRIPTION, tableName=RESOURCE_SERVER_POLICY; modifyDataType columnName=DESCRIPTION,...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('2.5.0-unicode-other-dbs','hmlnarik@redhat.com','META-INF/jpa-changelog-2.5.0.xml','2025-12-17 12:46:40',35,'EXECUTED','9:33d72168746f81f98ae3a1e8e0ca3554','modifyDataType columnName=DESCRIPTION, tableName=AUTHENTICATION_FLOW; modifyDataType columnName=DESCRIPTION, tableName=CLIENT_TEMPLATE; modifyDataType columnName=DESCRIPTION, tableName=RESOURCE_SERVER_POLICY; modifyDataType columnName=DESCRIPTION,...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('2.5.0-unique-group-names','hmlnarik@redhat.com','META-INF/jpa-changelog-2.5.0.xml','2025-12-17 12:46:40',37,'EXECUTED','9:8dcac7bdf7378e7d823cdfddebf72fda','addUniqueConstraint constraintName=SIBLING_NAMES, tableName=KEYCLOAK_GROUP','',NULL,'4.33.0',NULL,NULL,'5993589629'),('2.5.1','bburke@redhat.com','META-INF/jpa-changelog-2.5.1.xml','2025-12-17 12:46:40',38,'EXECUTED','9:a2b870802540cb3faa72098db5388af3','addColumn tableName=FED_USER_CONSENT','',NULL,'4.33.0',NULL,NULL,'5993589629'),('20.0.0-12964-supported-dbs','keycloak','META-INF/jpa-changelog-20.0.0.xml','2025-12-17 12:46:45',109,'EXECUTED','9:f2e1331a71e0aa85e5608fe42f7f681c','createIndex indexName=IDX_GROUP_ATT_BY_NAME_VALUE, tableName=GROUP_ATTRIBUTE','',NULL,'4.33.0',NULL,NULL,'5993589629'),('20.0.0-12964-supported-dbs-edb-migration','keycloak','META-INF/jpa-changelog-20.0.0.xml','2025-12-17 12:46:45',110,'MARK_RAN','9:a6b18a8e38062df5793edbe064f4aecd','dropIndex indexName=IDX_GROUP_ATT_BY_NAME_VALUE, tableName=GROUP_ATTRIBUTE; createIndex indexName=IDX_GROUP_ATT_BY_NAME_VALUE, tableName=GROUP_ATTRIBUTE','',NULL,'4.33.0',NULL,NULL,'5993589629'),('20.0.0-12964-unsupported-dbs','keycloak','META-INF/jpa-changelog-20.0.0.xml','2025-12-17 12:46:45',111,'MARK_RAN','9:1a6fcaa85e20bdeae0a9ce49b41946a5','createIndex indexName=IDX_GROUP_ATT_BY_NAME_VALUE, tableName=GROUP_ATTRIBUTE','',NULL,'4.33.0',NULL,NULL,'5993589629'),('21.0.2-17277','keycloak','META-INF/jpa-changelog-21.0.2.xml','2025-12-17 12:46:45',115,'EXECUTED','9:7ee1f7a3fb8f5588f171fb9a6ab623c0','customChange','',NULL,'4.33.0',NULL,NULL,'5993589629'),('21.1.0-19404','keycloak','META-INF/jpa-changelog-21.1.0.xml','2025-12-17 12:46:45',116,'EXECUTED','9:3d7e830b52f33676b9d64f7f2b2ea634','modifyDataType columnName=DECISION_STRATEGY, tableName=RESOURCE_SERVER_POLICY; modifyDataType columnName=LOGIC, tableName=RESOURCE_SERVER_POLICY; modifyDataType columnName=POLICY_ENFORCE_MODE, tableName=RESOURCE_SERVER','',NULL,'4.33.0',NULL,NULL,'5993589629'),('21.1.0-19404-2','keycloak','META-INF/jpa-changelog-21.1.0.xml','2025-12-17 12:46:45',117,'MARK_RAN','9:627d032e3ef2c06c0e1f73d2ae25c26c','addColumn tableName=RESOURCE_SERVER_POLICY; update tableName=RESOURCE_SERVER_POLICY; dropColumn columnName=DECISION_STRATEGY, tableName=RESOURCE_SERVER_POLICY; renameColumn newColumnName=DECISION_STRATEGY, oldColumnName=DECISION_STRATEGY_NEW, tabl...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('22.0.0-17484-updated','keycloak','META-INF/jpa-changelog-22.0.0.xml','2025-12-17 12:46:45',118,'EXECUTED','9:90af0bfd30cafc17b9f4d6eccd92b8b3','customChange','',NULL,'4.33.0',NULL,NULL,'5993589629'),('22.0.5-24031','keycloak','META-INF/jpa-changelog-22.0.0.xml','2025-12-17 12:46:45',119,'MARK_RAN','9:a60d2d7b315ec2d3eba9e2f145f9df28','customChange','',NULL,'4.33.0',NULL,NULL,'5993589629'),('23.0.0-12062','keycloak','META-INF/jpa-changelog-23.0.0.xml','2025-12-17 12:46:45',120,'EXECUTED','9:2168fbe728fec46ae9baf15bf80927b8','addColumn tableName=COMPONENT_CONFIG; update tableName=COMPONENT_CONFIG; dropColumn columnName=VALUE, tableName=COMPONENT_CONFIG; renameColumn newColumnName=VALUE, oldColumnName=VALUE_NEW, tableName=COMPONENT_CONFIG','',NULL,'4.33.0',NULL,NULL,'5993589629'),('23.0.0-17258','keycloak','META-INF/jpa-changelog-23.0.0.xml','2025-12-17 12:46:45',121,'EXECUTED','9:36506d679a83bbfda85a27ea1864dca8','addColumn tableName=EVENT_ENTITY','',NULL,'4.33.0',NULL,NULL,'5993589629'),('24.0.0-26618-drop-index-if-present','keycloak','META-INF/jpa-changelog-24.0.0.xml','2025-12-17 12:46:45',124,'MARK_RAN','9:04baaf56c116ed19951cbc2cca584022','dropIndex indexName=IDX_CLIENT_ATT_BY_NAME_VALUE, tableName=CLIENT_ATTRIBUTES','',NULL,'4.33.0',NULL,NULL,'5993589629'),('24.0.0-26618-edb-migration','keycloak','META-INF/jpa-changelog-24.0.0.xml','2025-12-17 12:46:45',126,'MARK_RAN','9:2f684b29d414cd47efe3a3599f390741','dropIndex indexName=IDX_CLIENT_ATT_BY_NAME_VALUE, tableName=CLIENT_ATTRIBUTES; createIndex indexName=IDX_CLIENT_ATT_BY_NAME_VALUE, tableName=CLIENT_ATTRIBUTES','',NULL,'4.33.0',NULL,NULL,'5993589629'),('24.0.0-26618-reindex','keycloak','META-INF/jpa-changelog-24.0.0.xml','2025-12-17 12:46:45',125,'EXECUTED','9:bd2bd0fc7768cf0845ac96a8786fa735','createIndex indexName=IDX_CLIENT_ATT_BY_NAME_VALUE, tableName=CLIENT_ATTRIBUTES','',NULL,'4.33.0',NULL,NULL,'5993589629'),('24.0.0-9758','keycloak','META-INF/jpa-changelog-24.0.0.xml','2025-12-17 12:46:45',122,'EXECUTED','9:502c557a5189f600f0f445a9b49ebbce','addColumn tableName=USER_ATTRIBUTE; addColumn tableName=FED_USER_ATTRIBUTE; createIndex indexName=USER_ATTR_LONG_VALUES, tableName=USER_ATTRIBUTE; createIndex indexName=FED_USER_ATTR_LONG_VALUES, tableName=FED_USER_ATTRIBUTE; createIndex indexName...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('24.0.0-9758-2','keycloak','META-INF/jpa-changelog-24.0.0.xml','2025-12-17 12:46:45',123,'EXECUTED','9:bf0fdee10afdf597a987adbf291db7b2','customChange','',NULL,'4.33.0',NULL,NULL,'5993589629'),('24.0.2-27228','keycloak','META-INF/jpa-changelog-24.0.2.xml','2025-12-17 12:46:45',127,'EXECUTED','9:eaee11f6b8aa25d2cc6a84fb86fc6238','customChange','',NULL,'4.33.0',NULL,NULL,'5993589629'),('24.0.2-27967-drop-index-if-present','keycloak','META-INF/jpa-changelog-24.0.2.xml','2025-12-17 12:46:45',128,'MARK_RAN','9:04baaf56c116ed19951cbc2cca584022','dropIndex indexName=IDX_CLIENT_ATT_BY_NAME_VALUE, tableName=CLIENT_ATTRIBUTES','',NULL,'4.33.0',NULL,NULL,'5993589629'),('24.0.2-27967-reindex','keycloak','META-INF/jpa-changelog-24.0.2.xml','2025-12-17 12:46:45',129,'MARK_RAN','9:d3d977031d431db16e2c181ce49d73e9','createIndex indexName=IDX_CLIENT_ATT_BY_NAME_VALUE, tableName=CLIENT_ATTRIBUTES','',NULL,'4.33.0',NULL,NULL,'5993589629'),('25.0.0-28265-index-2-mysql','keycloak','META-INF/jpa-changelog-25.0.0.xml','2025-12-17 12:46:45',136,'EXECUTED','9:b7ef76036d3126bb83c2423bf4d449d6','createIndex indexName=IDX_OFFLINE_USS_BY_BROKER_SESSION_ID, tableName=OFFLINE_USER_SESSION','',NULL,'4.33.0',NULL,NULL,'5993589629'),('25.0.0-28265-index-2-not-mysql','keycloak','META-INF/jpa-changelog-25.0.0.xml','2025-12-17 12:46:45',137,'MARK_RAN','9:23396cf51ab8bc1ae6f0cac7f9f6fcf7','createIndex indexName=IDX_OFFLINE_USS_BY_BROKER_SESSION_ID, tableName=OFFLINE_USER_SESSION','',NULL,'4.33.0',NULL,NULL,'5993589629'),('25.0.0-28265-index-cleanup-css-preload','keycloak','META-INF/jpa-changelog-25.0.0.xml','2025-12-17 12:46:45',135,'EXECUTED','9:5411d2fb2891d3e8d63ddb55dfa3c0c9','dropIndex indexName=IDX_OFFLINE_CSS_PRELOAD, tableName=OFFLINE_CLIENT_SESSION','',NULL,'4.33.0',NULL,NULL,'5993589629'),('25.0.0-28265-index-cleanup-uss-by-usersess','keycloak','META-INF/jpa-changelog-25.0.0.xml','2025-12-17 12:46:45',134,'EXECUTED','9:6eee220d024e38e89c799417ec33667f','dropIndex indexName=IDX_OFFLINE_USS_BY_USERSESS, tableName=OFFLINE_USER_SESSION','',NULL,'4.33.0',NULL,NULL,'5993589629'),('25.0.0-28265-index-cleanup-uss-createdon','keycloak','META-INF/jpa-changelog-25.0.0.xml','2025-12-17 12:46:45',132,'EXECUTED','9:78ab4fc129ed5e8265dbcc3485fba92f','dropIndex indexName=IDX_OFFLINE_USS_CREATEDON, tableName=OFFLINE_USER_SESSION','',NULL,'4.33.0',NULL,NULL,'5993589629'),('25.0.0-28265-index-cleanup-uss-preload','keycloak','META-INF/jpa-changelog-25.0.0.xml','2025-12-17 12:46:45',133,'EXECUTED','9:de5f7c1f7e10994ed8b62e621d20eaab','dropIndex indexName=IDX_OFFLINE_USS_PRELOAD, tableName=OFFLINE_USER_SESSION','',NULL,'4.33.0',NULL,NULL,'5993589629'),('25.0.0-28265-index-creation','keycloak','META-INF/jpa-changelog-25.0.0.xml','2025-12-17 12:46:45',131,'EXECUTED','9:3e96709818458ae49f3c679ae58d263a','createIndex indexName=IDX_OFFLINE_USS_BY_LAST_SESSION_REFRESH, tableName=OFFLINE_USER_SESSION','',NULL,'4.33.0',NULL,NULL,'5993589629'),('25.0.0-28265-tables','keycloak','META-INF/jpa-changelog-25.0.0.xml','2025-12-17 12:46:45',130,'EXECUTED','9:deda2df035df23388af95bbd36c17cef','addColumn tableName=OFFLINE_USER_SESSION; addColumn tableName=OFFLINE_CLIENT_SESSION','',NULL,'4.33.0',NULL,NULL,'5993589629'),('25.0.0-28861-index-creation','keycloak','META-INF/jpa-changelog-25.0.0.xml','2025-12-17 12:46:46',142,'EXECUTED','9:b9acb58ac958d9ada0fe12a5d4794ab1','createIndex indexName=IDX_PERM_TICKET_REQUESTER, tableName=RESOURCE_SERVER_PERM_TICKET; createIndex indexName=IDX_PERM_TICKET_OWNER, tableName=RESOURCE_SERVER_PERM_TICKET','',NULL,'4.33.0',NULL,NULL,'5993589629'),('25.0.0-org','keycloak','META-INF/jpa-changelog-25.0.0.xml','2025-12-17 12:46:46',138,'EXECUTED','9:5c859965c2c9b9c72136c360649af157','createTable tableName=ORG; addUniqueConstraint constraintName=UK_ORG_NAME, tableName=ORG; addUniqueConstraint constraintName=UK_ORG_GROUP, tableName=ORG; createTable tableName=ORG_DOMAIN','',NULL,'4.33.0',NULL,NULL,'5993589629'),('26.0.0-32583-drop-redundant-index-on-client-session','keycloak','META-INF/jpa-changelog-26.0.0.xml','2025-12-17 12:46:46',150,'EXECUTED','9:24972d83bf27317a055d234187bb4af9','dropIndex indexName=IDX_US_SESS_ID_ON_CL_SESS, tableName=OFFLINE_CLIENT_SESSION','',NULL,'4.33.0',NULL,NULL,'5993589629'),('26.0.0-33201-org-redirect-url','keycloak','META-INF/jpa-changelog-26.0.0.xml','2025-12-17 12:46:46',152,'EXECUTED','9:4d0e22b0ac68ebe9794fa9cb752ea660','addColumn tableName=ORG','',NULL,'4.33.0',NULL,NULL,'5993589629'),('26.0.0-idps-for-login','keycloak','META-INF/jpa-changelog-26.0.0.xml','2025-12-17 12:46:46',149,'EXECUTED','9:51f5fffadf986983d4bd59582c6c1604','addColumn tableName=IDENTITY_PROVIDER; createIndex indexName=IDX_IDP_REALM_ORG, tableName=IDENTITY_PROVIDER; createIndex indexName=IDX_IDP_FOR_LOGIN, tableName=IDENTITY_PROVIDER; customChange','',NULL,'4.33.0',NULL,NULL,'5993589629'),('26.0.0-org-alias','keycloak','META-INF/jpa-changelog-26.0.0.xml','2025-12-17 12:46:46',143,'EXECUTED','9:6ef7d63e4412b3c2d66ed179159886a4','addColumn tableName=ORG; update tableName=ORG; addNotNullConstraint columnName=ALIAS, tableName=ORG; addUniqueConstraint constraintName=UK_ORG_ALIAS, tableName=ORG','',NULL,'4.33.0',NULL,NULL,'5993589629'),('26.0.0-org-group','keycloak','META-INF/jpa-changelog-26.0.0.xml','2025-12-17 12:46:46',144,'EXECUTED','9:da8e8087d80ef2ace4f89d8c5b9ca223','addColumn tableName=KEYCLOAK_GROUP; update tableName=KEYCLOAK_GROUP; addNotNullConstraint columnName=TYPE, tableName=KEYCLOAK_GROUP; customChange','',NULL,'4.33.0',NULL,NULL,'5993589629'),('26.0.0-org-group-membership','keycloak','META-INF/jpa-changelog-26.0.0.xml','2025-12-17 12:46:46',146,'EXECUTED','9:a6ace2ce583a421d89b01ba2a28dc2d4','addColumn tableName=USER_GROUP_MEMBERSHIP; update tableName=USER_GROUP_MEMBERSHIP; addNotNullConstraint columnName=MEMBERSHIP_TYPE, tableName=USER_GROUP_MEMBERSHIP','',NULL,'4.33.0',NULL,NULL,'5993589629'),('26.0.0-org-indexes','keycloak','META-INF/jpa-changelog-26.0.0.xml','2025-12-17 12:46:46',145,'EXECUTED','9:79b05dcd610a8c7f25ec05135eec0857','createIndex indexName=IDX_ORG_DOMAIN_ORG_ID, tableName=ORG_DOMAIN','',NULL,'4.33.0',NULL,NULL,'5993589629'),('26.0.0.32582-remove-tables-user-session-user-session-note-and-client-session','keycloak','META-INF/jpa-changelog-26.0.0.xml','2025-12-17 12:46:46',151,'EXECUTED','9:febdc0f47f2ed241c59e60f58c3ceea5','dropTable tableName=CLIENT_SESSION_ROLE; dropTable tableName=CLIENT_SESSION_NOTE; dropTable tableName=CLIENT_SESSION_PROT_MAPPER; dropTable tableName=CLIENT_SESSION_AUTH_STATUS; dropTable tableName=CLIENT_USER_SESSION_NOTE; dropTable tableName=CLI...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('26.1.0-34013','keycloak','META-INF/jpa-changelog-26.1.0.xml','2025-12-17 12:46:46',154,'EXECUTED','9:e6b686a15759aef99a6d758a5c4c6a26','addColumn tableName=ADMIN_EVENT_ENTITY','',NULL,'4.33.0',NULL,NULL,'5993589629'),('26.1.0-34380','keycloak','META-INF/jpa-changelog-26.1.0.xml','2025-12-17 12:46:46',155,'EXECUTED','9:ac8b9edb7c2b6c17a1c7a11fcf5ccf01','dropTable tableName=USERNAME_LOGIN_FAILURE','',NULL,'4.33.0',NULL,NULL,'5993589629'),('26.2.0-26106','keycloak','META-INF/jpa-changelog-26.2.0.xml','2025-12-17 12:46:46',157,'EXECUTED','9:b5877d5dab7d10ff3a9d209d7beb6680','addColumn tableName=CREDENTIAL','',NULL,'4.33.0',NULL,NULL,'5993589629'),('26.2.0-36750','keycloak','META-INF/jpa-changelog-26.2.0.xml','2025-12-17 12:46:46',156,'EXECUTED','9:b49ce951c22f7eb16480ff085640a33a','createTable tableName=SERVER_CONFIG','',NULL,'4.33.0',NULL,NULL,'5993589629'),('26.2.6-39866-duplicate','keycloak','META-INF/jpa-changelog-26.2.6.xml','2025-12-17 12:46:46',158,'EXECUTED','9:1dc67ccee24f30331db2cba4f372e40e','customChange','',NULL,'4.33.0',NULL,NULL,'5993589629'),('26.2.6-39866-uk','keycloak','META-INF/jpa-changelog-26.2.6.xml','2025-12-17 12:46:46',159,'EXECUTED','9:b70b76f47210cf0a5f4ef0e219eac7cd','addUniqueConstraint constraintName=UK_MIGRATION_VERSION, tableName=MIGRATION_MODEL','',NULL,'4.33.0',NULL,NULL,'5993589629'),('26.2.6-40088-duplicate','keycloak','META-INF/jpa-changelog-26.2.6.xml','2025-12-17 12:46:46',160,'EXECUTED','9:cc7e02ed69ab31979afb1982f9670e8f','customChange','',NULL,'4.33.0',NULL,NULL,'5993589629'),('26.2.6-40088-uk','keycloak','META-INF/jpa-changelog-26.2.6.xml','2025-12-17 12:46:46',161,'EXECUTED','9:5bb848128da7bc4595cc507383325241','addUniqueConstraint constraintName=UK_MIGRATION_UPDATE_TIME, tableName=MIGRATION_MODEL','',NULL,'4.33.0',NULL,NULL,'5993589629'),('26.3.0-groups-description','keycloak','META-INF/jpa-changelog-26.3.0.xml','2025-12-17 12:46:46',162,'EXECUTED','9:e1a3c05574326fb5b246b73b9a4c4d49','addColumn tableName=KEYCLOAK_GROUP','',NULL,'4.33.0',NULL,NULL,'5993589629'),('26.4.0-40933-saml-encryption-attributes','keycloak','META-INF/jpa-changelog-26.4.0.xml','2025-12-17 12:46:46',163,'EXECUTED','9:7e9eaba362ca105efdda202303a4fe49','customChange','',NULL,'4.33.0',NULL,NULL,'5993589629'),('26.4.0-51321','keycloak','META-INF/jpa-changelog-26.4.0.xml','2025-12-17 12:46:46',164,'EXECUTED','9:34bab2bc56f75ffd7e347c580874e306','createIndex indexName=IDX_EVENT_ENTITY_USER_ID_TYPE, tableName=EVENT_ENTITY','',NULL,'4.33.0',NULL,NULL,'5993589629'),('26.5.0-index-offline-css-by-client','keycloak','META-INF/jpa-changelog-26.5.0.xml','2025-12-17 12:46:46',166,'EXECUTED','9:680b59ca7854fa5b77a303301bb2a941','createIndex indexName=IDX_OFFLINE_CSS_BY_CLIENT, tableName=OFFLINE_CLIENT_SESSION','',NULL,'4.33.0',NULL,NULL,'5993589629'),('26.5.0-index-offline-css-by-client-storage-provider','keycloak','META-INF/jpa-changelog-26.5.0.xml','2025-12-17 12:46:46',167,'EXECUTED','9:809bc160e2bc92f9c28eea39db323ae2','createIndex indexName=IDX_OFFLINE_CSS_BY_CLIENT_STORAGE_PROVIDER, tableName=OFFLINE_CLIENT_SESSION','',NULL,'4.33.0',NULL,NULL,'5993589629'),('29399-jdbc-ping-default','keycloak','META-INF/jpa-changelog-26.1.0.xml','2025-12-17 12:46:46',153,'EXECUTED','9:007dbe99d7203fca403b89d4edfdf21e','createTable tableName=JGROUPS_PING; addPrimaryKey constraintName=CONSTRAINT_JGROUPS_PING, tableName=JGROUPS_PING','',NULL,'4.33.0',NULL,NULL,'5993589629'),('3.0.0','bburke@redhat.com','META-INF/jpa-changelog-3.0.0.xml','2025-12-17 12:46:40',39,'EXECUTED','9:132a67499ba24bcc54fb5cbdcfe7e4c0','addColumn tableName=IDENTITY_PROVIDER','',NULL,'4.33.0',NULL,NULL,'5993589629'),('3.2.0-fix','keycloak','META-INF/jpa-changelog-3.2.0.xml','2025-12-17 12:46:40',40,'MARK_RAN','9:938f894c032f5430f2b0fafb1a243462','addNotNullConstraint columnName=REALM_ID, tableName=CLIENT_INITIAL_ACCESS','',NULL,'4.33.0',NULL,NULL,'5993589629'),('3.2.0-fix-offline-sessions','hmlnarik','META-INF/jpa-changelog-3.2.0.xml','2025-12-17 12:46:40',42,'EXECUTED','9:fc86359c079781adc577c5a217e4d04c','customChange','',NULL,'4.33.0',NULL,NULL,'5993589629'),('3.2.0-fix-with-keycloak-5416','keycloak','META-INF/jpa-changelog-3.2.0.xml','2025-12-17 12:46:40',41,'MARK_RAN','9:845c332ff1874dc5d35974b0babf3006','dropIndex indexName=IDX_CLIENT_INIT_ACC_REALM, tableName=CLIENT_INITIAL_ACCESS; addNotNullConstraint columnName=REALM_ID, tableName=CLIENT_INITIAL_ACCESS; createIndex indexName=IDX_CLIENT_INIT_ACC_REALM, tableName=CLIENT_INITIAL_ACCESS','',NULL,'4.33.0',NULL,NULL,'5993589629'),('3.2.0-fixed','keycloak','META-INF/jpa-changelog-3.2.0.xml','2025-12-17 12:46:41',43,'EXECUTED','9:59a64800e3c0d09b825f8a3b444fa8f4','addColumn tableName=REALM; dropPrimaryKey constraintName=CONSTRAINT_OFFL_CL_SES_PK2, tableName=OFFLINE_CLIENT_SESSION; dropColumn columnName=CLIENT_SESSION_ID, tableName=OFFLINE_CLIENT_SESSION; addPrimaryKey constraintName=CONSTRAINT_OFFL_CL_SES_P...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('3.3.0','keycloak','META-INF/jpa-changelog-3.3.0.xml','2025-12-17 12:46:41',44,'EXECUTED','9:d48d6da5c6ccf667807f633fe489ce88','addColumn tableName=USER_ENTITY','',NULL,'4.33.0',NULL,NULL,'5993589629'),('3.4.0','keycloak','META-INF/jpa-changelog-3.4.0.xml','2025-12-17 12:46:42',50,'EXECUTED','9:cfdd8736332ccdd72c5256ccb42335db','addPrimaryKey constraintName=CONSTRAINT_REALM_DEFAULT_ROLES, tableName=REALM_DEFAULT_ROLES; addPrimaryKey constraintName=CONSTRAINT_COMPOSITE_ROLE, tableName=COMPOSITE_ROLE; addPrimaryKey constraintName=CONSTR_REALM_DEFAULT_GROUPS, tableName=REALM...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('3.4.0-KEYCLOAK-5230','hmlnarik@redhat.com','META-INF/jpa-changelog-3.4.0.xml','2025-12-17 12:46:42',51,'EXECUTED','9:7c84de3d9bd84d7f077607c1a4dcb714','createIndex indexName=IDX_FU_ATTRIBUTE, tableName=FED_USER_ATTRIBUTE; createIndex indexName=IDX_FU_CONSENT, tableName=FED_USER_CONSENT; createIndex indexName=IDX_FU_CONSENT_RU, tableName=FED_USER_CONSENT; createIndex indexName=IDX_FU_CREDENTIAL, t...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('3.4.1','psilva@redhat.com','META-INF/jpa-changelog-3.4.1.xml','2025-12-17 12:46:42',52,'EXECUTED','9:5a6bb36cbefb6a9d6928452c0852af2d','modifyDataType columnName=VALUE, tableName=CLIENT_ATTRIBUTES','',NULL,'4.33.0',NULL,NULL,'5993589629'),('3.4.2','keycloak','META-INF/jpa-changelog-3.4.2.xml','2025-12-17 12:46:42',53,'EXECUTED','9:8f23e334dbc59f82e0a328373ca6ced0','update tableName=REALM','',NULL,'4.33.0',NULL,NULL,'5993589629'),('3.4.2-KEYCLOAK-5172','mkanis@redhat.com','META-INF/jpa-changelog-3.4.2.xml','2025-12-17 12:46:42',54,'EXECUTED','9:9156214268f09d970cdf0e1564d866af','update tableName=CLIENT','',NULL,'4.33.0',NULL,NULL,'5993589629'),('31296-persist-revoked-access-tokens','keycloak','META-INF/jpa-changelog-26.0.0.xml','2025-12-17 12:46:46',147,'EXECUTED','9:64ef94489d42a358e8304b0e245f0ed4','createTable tableName=REVOKED_TOKEN; addPrimaryKey constraintName=CONSTRAINT_RT, tableName=REVOKED_TOKEN','',NULL,'4.33.0',NULL,NULL,'5993589629'),('31725-index-persist-revoked-access-tokens','keycloak','META-INF/jpa-changelog-26.0.0.xml','2025-12-17 12:46:46',148,'EXECUTED','9:b994246ec2bf7c94da881e1d28782c7b','createIndex indexName=IDX_REV_TOKEN_ON_EXPIRE, tableName=REVOKED_TOKEN','',NULL,'4.33.0',NULL,NULL,'5993589629'),('4.0.0-CLEANUP-UNUSED-TABLE','bburke@redhat.com','META-INF/jpa-changelog-4.0.0.xml','2025-12-17 12:46:42',56,'EXECUTED','9:229a041fb72d5beac76bb94a5fa709de','dropTable tableName=CLIENT_IDENTITY_PROV_MAPPING','',NULL,'4.33.0',NULL,NULL,'5993589629'),('4.0.0-KEYCLOAK-5579-fixed','mposolda@redhat.com','META-INF/jpa-changelog-4.0.0.xml','2025-12-17 12:46:43',58,'EXECUTED','9:139b79bcbbfe903bb1c2d2a4dbf001d9','dropForeignKeyConstraint baseTableName=CLIENT_TEMPLATE_ATTRIBUTES, constraintName=FK_CL_TEMPL_ATTR_TEMPL; renameTable newTableName=CLIENT_SCOPE_ATTRIBUTES, oldTableName=CLIENT_TEMPLATE_ATTRIBUTES; renameColumn newColumnName=SCOPE_ID, oldColumnName...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('4.0.0-KEYCLOAK-6228','bburke@redhat.com','META-INF/jpa-changelog-4.0.0.xml','2025-12-17 12:46:42',57,'EXECUTED','9:079899dade9c1e683f26b2aa9ca6ff04','dropUniqueConstraint constraintName=UK_JKUWUVD56ONTGSUHOGM8UEWRT, tableName=USER_CONSENT; dropNotNullConstraint columnName=CLIENT_ID, tableName=USER_CONSENT; addColumn tableName=USER_CONSENT; addUniqueConstraint constraintName=UK_JKUWUVD56ONTGSUHO...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('4.0.0-KEYCLOAK-6335','bburke@redhat.com','META-INF/jpa-changelog-4.0.0.xml','2025-12-17 12:46:42',55,'EXECUTED','9:db806613b1ed154826c02610b7dbdf74','createTable tableName=CLIENT_AUTH_FLOW_BINDINGS; addPrimaryKey constraintName=C_CLI_FLOW_BIND, tableName=CLIENT_AUTH_FLOW_BINDINGS','',NULL,'4.33.0',NULL,NULL,'5993589629'),('4.2.0-KEYCLOAK-6313','wadahiro@gmail.com','META-INF/jpa-changelog-4.2.0.xml','2025-12-17 12:46:43',63,'EXECUTED','9:92143a6daea0a3f3b8f598c97ce55c3d','addColumn tableName=REQUIRED_ACTION_PROVIDER','',NULL,'4.33.0',NULL,NULL,'5993589629'),('4.3.0-KEYCLOAK-7984','wadahiro@gmail.com','META-INF/jpa-changelog-4.3.0.xml','2025-12-17 12:46:43',64,'EXECUTED','9:82bab26a27195d889fb0429003b18f40','update tableName=REQUIRED_ACTION_PROVIDER','',NULL,'4.33.0',NULL,NULL,'5993589629'),('4.6.0-KEYCLOAK-7950','psilva@redhat.com','META-INF/jpa-changelog-4.6.0.xml','2025-12-17 12:46:43',65,'EXECUTED','9:e590c88ddc0b38b0ae4249bbfcb5abc3','update tableName=RESOURCE_SERVER_RESOURCE','',NULL,'4.33.0',NULL,NULL,'5993589629'),('4.6.0-KEYCLOAK-8377','keycloak','META-INF/jpa-changelog-4.6.0.xml','2025-12-17 12:46:43',66,'EXECUTED','9:5c1f475536118dbdc38d5d7977950cc0','createTable tableName=ROLE_ATTRIBUTE; addPrimaryKey constraintName=CONSTRAINT_ROLE_ATTRIBUTE_PK, tableName=ROLE_ATTRIBUTE; addForeignKeyConstraint baseTableName=ROLE_ATTRIBUTE, constraintName=FK_ROLE_ATTRIBUTE_ID, referencedTableName=KEYCLOAK_ROLE...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('4.6.0-KEYCLOAK-8555','gideonray@gmail.com','META-INF/jpa-changelog-4.6.0.xml','2025-12-17 12:46:43',67,'EXECUTED','9:e7c9f5f9c4d67ccbbcc215440c718a17','createIndex indexName=IDX_COMPONENT_PROVIDER_TYPE, tableName=COMPONENT','',NULL,'4.33.0',NULL,NULL,'5993589629'),('4.7.0-KEYCLOAK-1267','sguilhen@redhat.com','META-INF/jpa-changelog-4.7.0.xml','2025-12-17 12:46:43',68,'EXECUTED','9:88e0bfdda924690d6f4e430c53447dd5','addColumn tableName=REALM','',NULL,'4.33.0',NULL,NULL,'5993589629'),('4.7.0-KEYCLOAK-7275','keycloak','META-INF/jpa-changelog-4.7.0.xml','2025-12-17 12:46:43',69,'EXECUTED','9:f53177f137e1c46b6a88c59ec1cb5218','renameColumn newColumnName=CREATED_ON, oldColumnName=LAST_SESSION_REFRESH, tableName=OFFLINE_USER_SESSION; addNotNullConstraint columnName=CREATED_ON, tableName=OFFLINE_USER_SESSION; addColumn tableName=OFFLINE_USER_SESSION; customChange; createIn...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('4.8.0-KEYCLOAK-8835','sguilhen@redhat.com','META-INF/jpa-changelog-4.8.0.xml','2025-12-17 12:46:43',70,'EXECUTED','9:a74d33da4dc42a37ec27121580d1459f','addNotNullConstraint columnName=SSO_MAX_LIFESPAN_REMEMBER_ME, tableName=REALM; addNotNullConstraint columnName=SSO_IDLE_TIMEOUT_REMEMBER_ME, tableName=REALM','',NULL,'4.33.0',NULL,NULL,'5993589629'),('40343-workflow-state-table','keycloak','META-INF/jpa-changelog-26.4.0.xml','2025-12-17 12:46:46',165,'EXECUTED','9:ed3ab4723ceed210e5b5e60ac4562106','createTable tableName=WORKFLOW_STATE; addPrimaryKey constraintName=PK_WORKFLOW_STATE, tableName=WORKFLOW_STATE; addUniqueConstraint constraintName=UQ_WORKFLOW_RESOURCE, tableName=WORKFLOW_STATE; createIndex indexName=IDX_WORKFLOW_STATE_STEP, table...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('8.0.0-adding-credential-columns','keycloak','META-INF/jpa-changelog-8.0.0.xml','2025-12-17 12:46:44',72,'EXECUTED','9:aa072ad090bbba210d8f18781b8cebf4','addColumn tableName=CREDENTIAL; addColumn tableName=FED_USER_CREDENTIAL','',NULL,'4.33.0',NULL,NULL,'5993589629'),('8.0.0-credential-cleanup-fixed','keycloak','META-INF/jpa-changelog-8.0.0.xml','2025-12-17 12:46:44',75,'EXECUTED','9:2b9cc12779be32c5b40e2e67711a218b','dropDefaultValue columnName=COUNTER, tableName=CREDENTIAL; dropDefaultValue columnName=DIGITS, tableName=CREDENTIAL; dropDefaultValue columnName=PERIOD, tableName=CREDENTIAL; dropDefaultValue columnName=ALGORITHM, tableName=CREDENTIAL; dropColumn ...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('8.0.0-resource-tag-support','keycloak','META-INF/jpa-changelog-8.0.0.xml','2025-12-17 12:46:44',76,'EXECUTED','9:91fa186ce7a5af127a2d7a91ee083cc5','addColumn tableName=MIGRATION_MODEL; createIndex indexName=IDX_UPDATE_TIME, tableName=MIGRATION_MODEL','',NULL,'4.33.0',NULL,NULL,'5993589629'),('8.0.0-updating-credential-data-not-oracle-fixed','keycloak','META-INF/jpa-changelog-8.0.0.xml','2025-12-17 12:46:44',73,'EXECUTED','9:1ae6be29bab7c2aa376f6983b932be37','update tableName=CREDENTIAL; update tableName=CREDENTIAL; update tableName=CREDENTIAL; update tableName=FED_USER_CREDENTIAL; update tableName=FED_USER_CREDENTIAL; update tableName=FED_USER_CREDENTIAL','',NULL,'4.33.0',NULL,NULL,'5993589629'),('8.0.0-updating-credential-data-oracle-fixed','keycloak','META-INF/jpa-changelog-8.0.0.xml','2025-12-17 12:46:44',74,'MARK_RAN','9:14706f286953fc9a25286dbd8fb30d97','update tableName=CREDENTIAL; update tableName=CREDENTIAL; update tableName=CREDENTIAL; update tableName=FED_USER_CREDENTIAL; update tableName=FED_USER_CREDENTIAL; update tableName=FED_USER_CREDENTIAL','',NULL,'4.33.0',NULL,NULL,'5993589629'),('9.0.0-always-display-client','keycloak','META-INF/jpa-changelog-9.0.0.xml','2025-12-17 12:46:44',77,'EXECUTED','9:6335e5c94e83a2639ccd68dd24e2e5ad','addColumn tableName=CLIENT','',NULL,'4.33.0',NULL,NULL,'5993589629'),('9.0.0-drop-constraints-for-column-increase','keycloak','META-INF/jpa-changelog-9.0.0.xml','2025-12-17 12:46:44',78,'MARK_RAN','9:6bdb5658951e028bfe16fa0a8228b530','dropUniqueConstraint constraintName=UK_FRSR6T700S9V50BU18WS5PMT, tableName=RESOURCE_SERVER_PERM_TICKET; dropUniqueConstraint constraintName=UK_FRSR6T700S9V50BU18WS5HA6, tableName=RESOURCE_SERVER_RESOURCE; dropPrimaryKey constraintName=CONSTRAINT_O...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('9.0.0-increase-column-size-federated-fk','keycloak','META-INF/jpa-changelog-9.0.0.xml','2025-12-17 12:46:44',79,'EXECUTED','9:d5bc15a64117ccad481ce8792d4c608f','modifyDataType columnName=CLIENT_ID, tableName=FED_USER_CONSENT; modifyDataType columnName=CLIENT_REALM_CONSTRAINT, tableName=KEYCLOAK_ROLE; modifyDataType columnName=OWNER, tableName=RESOURCE_SERVER_POLICY; modifyDataType columnName=CLIENT_ID, ta...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('9.0.0-recreate-constraints-after-column-increase','keycloak','META-INF/jpa-changelog-9.0.0.xml','2025-12-17 12:46:44',80,'MARK_RAN','9:077cba51999515f4d3e7ad5619ab592c','addNotNullConstraint columnName=CLIENT_ID, tableName=OFFLINE_CLIENT_SESSION; addNotNullConstraint columnName=OWNER, tableName=RESOURCE_SERVER_PERM_TICKET; addNotNullConstraint columnName=REQUESTER, tableName=RESOURCE_SERVER_PERM_TICKET; addNotNull...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('9.0.1-add-index-to-client.client_id','keycloak','META-INF/jpa-changelog-9.0.1.xml','2025-12-17 12:46:44',81,'EXECUTED','9:be969f08a163bf47c6b9e9ead8ac2afb','createIndex indexName=IDX_CLIENT_ID, tableName=CLIENT','',NULL,'4.33.0',NULL,NULL,'5993589629'),('9.0.1-add-index-to-events','keycloak','META-INF/jpa-changelog-9.0.1.xml','2025-12-17 12:46:44',85,'EXECUTED','9:7d93d602352a30c0c317e6a609b56599','createIndex indexName=IDX_EVENT_TIME, tableName=EVENT_ENTITY','',NULL,'4.33.0',NULL,NULL,'5993589629'),('9.0.1-KEYCLOAK-12579-add-not-null-constraint','keycloak','META-INF/jpa-changelog-9.0.1.xml','2025-12-17 12:46:44',83,'EXECUTED','9:966bda61e46bebf3cc39518fbed52fa7','addNotNullConstraint columnName=PARENT_GROUP, tableName=KEYCLOAK_GROUP','',NULL,'4.33.0',NULL,NULL,'5993589629'),('9.0.1-KEYCLOAK-12579-drop-constraints','keycloak','META-INF/jpa-changelog-9.0.1.xml','2025-12-17 12:46:44',82,'MARK_RAN','9:6d3bb4408ba5a72f39bd8a0b301ec6e3','dropUniqueConstraint constraintName=SIBLING_NAMES, tableName=KEYCLOAK_GROUP','',NULL,'4.33.0',NULL,NULL,'5993589629'),('9.0.1-KEYCLOAK-12579-recreate-constraints','keycloak','META-INF/jpa-changelog-9.0.1.xml','2025-12-17 12:46:44',84,'MARK_RAN','9:8dcac7bdf7378e7d823cdfddebf72fda','addUniqueConstraint constraintName=SIBLING_NAMES, tableName=KEYCLOAK_GROUP','',NULL,'4.33.0',NULL,NULL,'5993589629'),('authn-3.4.0.CR1-refresh-token-max-reuse','glavoie@gmail.com','META-INF/jpa-changelog-authz-3.4.0.CR1.xml','2025-12-17 12:46:41',49,'EXECUTED','9:d198654156881c46bfba39abd7769e69','addColumn tableName=REALM','',NULL,'4.33.0',NULL,NULL,'5993589629'),('authz-2.0.0','psilva@redhat.com','META-INF/jpa-changelog-authz-2.0.0.xml','2025-12-17 12:46:39',27,'EXECUTED','9:43ed6b0da89ff77206289e87eaa9c024','createTable tableName=RESOURCE_SERVER; addPrimaryKey constraintName=CONSTRAINT_FARS, tableName=RESOURCE_SERVER; addUniqueConstraint constraintName=UK_AU8TT6T700S9V50BU18WS5HA6, tableName=RESOURCE_SERVER; createTable tableName=RESOURCE_SERVER_RESOU...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('authz-2.5.1','psilva@redhat.com','META-INF/jpa-changelog-authz-2.5.1.xml','2025-12-17 12:46:39',28,'EXECUTED','9:44bae577f551b3738740281eceb4ea70','update tableName=RESOURCE_SERVER_POLICY','',NULL,'4.33.0',NULL,NULL,'5993589629'),('authz-3.4.0.CR1-resource-server-pk-change-part1','glavoie@gmail.com','META-INF/jpa-changelog-authz-3.4.0.CR1.xml','2025-12-17 12:46:41',45,'EXECUTED','9:dde36f7973e80d71fceee683bc5d2951','addColumn tableName=RESOURCE_SERVER_POLICY; addColumn tableName=RESOURCE_SERVER_RESOURCE; addColumn tableName=RESOURCE_SERVER_SCOPE','',NULL,'4.33.0',NULL,NULL,'5993589629'),('authz-3.4.0.CR1-resource-server-pk-change-part2-KEYCLOAK-6095','hmlnarik@redhat.com','META-INF/jpa-changelog-authz-3.4.0.CR1.xml','2025-12-17 12:46:41',46,'EXECUTED','9:b855e9b0a406b34fa323235a0cf4f640','customChange','',NULL,'4.33.0',NULL,NULL,'5993589629'),('authz-3.4.0.CR1-resource-server-pk-change-part3-fixed','glavoie@gmail.com','META-INF/jpa-changelog-authz-3.4.0.CR1.xml','2025-12-17 12:46:41',47,'MARK_RAN','9:51abbacd7b416c50c4421a8cabf7927e','dropIndex indexName=IDX_RES_SERV_POL_RES_SERV, tableName=RESOURCE_SERVER_POLICY; dropIndex indexName=IDX_RES_SRV_RES_RES_SRV, tableName=RESOURCE_SERVER_RESOURCE; dropIndex indexName=IDX_RES_SRV_SCOPE_RES_SRV, tableName=RESOURCE_SERVER_SCOPE','',NULL,'4.33.0',NULL,NULL,'5993589629'),('authz-3.4.0.CR1-resource-server-pk-change-part3-fixed-nodropindex','glavoie@gmail.com','META-INF/jpa-changelog-authz-3.4.0.CR1.xml','2025-12-17 12:46:41',48,'EXECUTED','9:bdc99e567b3398bac83263d375aad143','addNotNullConstraint columnName=RESOURCE_SERVER_CLIENT_ID, tableName=RESOURCE_SERVER_POLICY; addNotNullConstraint columnName=RESOURCE_SERVER_CLIENT_ID, tableName=RESOURCE_SERVER_RESOURCE; addNotNullConstraint columnName=RESOURCE_SERVER_CLIENT_ID, ...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('authz-4.0.0.Beta3','psilva@redhat.com','META-INF/jpa-changelog-authz-4.0.0.Beta3.xml','2025-12-17 12:46:43',60,'EXECUTED','9:e0057eac39aa8fc8e09ac6cfa4ae15fe','addColumn tableName=RESOURCE_SERVER_POLICY; addColumn tableName=RESOURCE_SERVER_PERM_TICKET; addForeignKeyConstraint baseTableName=RESOURCE_SERVER_PERM_TICKET, constraintName=FK_FRSRPO2128CX4WNKOG82SSRFY, referencedTableName=RESOURCE_SERVER_POLICY','',NULL,'4.33.0',NULL,NULL,'5993589629'),('authz-4.0.0.CR1','psilva@redhat.com','META-INF/jpa-changelog-authz-4.0.0.CR1.xml','2025-12-17 12:46:43',59,'EXECUTED','9:b55738ad889860c625ba2bf483495a04','createTable tableName=RESOURCE_SERVER_PERM_TICKET; addPrimaryKey constraintName=CONSTRAINT_FAPMT, tableName=RESOURCE_SERVER_PERM_TICKET; addForeignKeyConstraint baseTableName=RESOURCE_SERVER_PERM_TICKET, constraintName=FK_FRSRHO213XCX4WNKOG82SSPMT...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('authz-4.2.0.Final','mhajas@redhat.com','META-INF/jpa-changelog-authz-4.2.0.Final.xml','2025-12-17 12:46:43',61,'EXECUTED','9:42a33806f3a0443fe0e7feeec821326c','createTable tableName=RESOURCE_URIS; addForeignKeyConstraint baseTableName=RESOURCE_URIS, constraintName=FK_RESOURCE_SERVER_URIS, referencedTableName=RESOURCE_SERVER_RESOURCE; customChange; dropColumn columnName=URI, tableName=RESOURCE_SERVER_RESO...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('authz-4.2.0.Final-KEYCLOAK-9944','hmlnarik@redhat.com','META-INF/jpa-changelog-authz-4.2.0.Final.xml','2025-12-17 12:46:43',62,'EXECUTED','9:9968206fca46eecc1f51db9c024bfe56','addPrimaryKey constraintName=CONSTRAINT_RESOUR_URIS_PK, tableName=RESOURCE_URIS','',NULL,'4.33.0',NULL,NULL,'5993589629'),('authz-7.0.0-KEYCLOAK-10443','psilva@redhat.com','META-INF/jpa-changelog-authz-7.0.0.xml','2025-12-17 12:46:44',71,'EXECUTED','9:fd4ade7b90c3b67fae0bfcfcb42dfb5f','addColumn tableName=RESOURCE_SERVER','',NULL,'4.33.0',NULL,NULL,'5993589629'),('client-attributes-string-accomodation-fixed','keycloak','META-INF/jpa-changelog-20.0.0.xml','2025-12-17 12:46:45',113,'EXECUTED','9:3f332e13e90739ed0c35b0b25b7822ca','addColumn tableName=CLIENT_ATTRIBUTES; update tableName=CLIENT_ATTRIBUTES; dropColumn columnName=VALUE, tableName=CLIENT_ATTRIBUTES; renameColumn newColumnName=VALUE, oldColumnName=VALUE_NEW, tableName=CLIENT_ATTRIBUTES','',NULL,'4.33.0',NULL,NULL,'5993589629'),('client-attributes-string-accomodation-fixed-post-create-index','keycloak','META-INF/jpa-changelog-20.0.0.xml','2025-12-17 12:46:45',114,'MARK_RAN','9:bd2bd0fc7768cf0845ac96a8786fa735','createIndex indexName=IDX_CLIENT_ATT_BY_NAME_VALUE, tableName=CLIENT_ATTRIBUTES','',NULL,'4.33.0',NULL,NULL,'5993589629'),('client-attributes-string-accomodation-fixed-pre-drop-index','keycloak','META-INF/jpa-changelog-20.0.0.xml','2025-12-17 12:46:45',112,'EXECUTED','9:04baaf56c116ed19951cbc2cca584022','dropIndex indexName=IDX_CLIENT_ATT_BY_NAME_VALUE, tableName=CLIENT_ATTRIBUTES','',NULL,'4.33.0',NULL,NULL,'5993589629'),('default-roles','keycloak','META-INF/jpa-changelog-13.0.0.xml','2025-12-17 12:46:44',89,'EXECUTED','9:fa8a5b5445e3857f4b010bafb5009957','addColumn tableName=REALM; customChange','',NULL,'4.33.0',NULL,NULL,'5993589629'),('default-roles-cleanup','keycloak','META-INF/jpa-changelog-13.0.0.xml','2025-12-17 12:46:44',90,'EXECUTED','9:67ac3241df9a8582d591c5ed87125f39','dropTable tableName=REALM_DEFAULT_ROLES; dropTable tableName=CLIENT_DEFAULT_ROLES','',NULL,'4.33.0',NULL,NULL,'5993589629'),('json-string-accomodation-fixed','keycloak','META-INF/jpa-changelog-13.0.0.xml','2025-12-17 12:46:45',96,'EXECUTED','9:e07d2bc0970c348bb06fb63b1f82ddbf','addColumn tableName=REALM_ATTRIBUTE; update tableName=REALM_ATTRIBUTE; dropColumn columnName=VALUE, tableName=REALM_ATTRIBUTE; renameColumn newColumnName=VALUE, oldColumnName=VALUE_NEW, tableName=REALM_ATTRIBUTE','',NULL,'4.33.0',NULL,NULL,'5993589629'),('KEYCLOAK-17267-add-index-to-user-attributes','keycloak','META-INF/jpa-changelog-14.0.0.xml','2025-12-17 12:46:45',102,'EXECUTED','9:0b305d8d1277f3a89a0a53a659ad274c','createIndex indexName=IDX_USER_ATTRIBUTE_NAME, tableName=USER_ATTRIBUTE','',NULL,'4.33.0',NULL,NULL,'5993589629'),('KEYCLOAK-18146-add-saml-art-binding-identifier','keycloak','META-INF/jpa-changelog-14.0.0.xml','2025-12-17 12:46:45',103,'EXECUTED','9:2c374ad2cdfe20e2905a84c8fac48460','customChange','',NULL,'4.33.0',NULL,NULL,'5993589629'),('map-remove-ri','keycloak','META-INF/jpa-changelog-11.0.0.xml','2025-12-17 12:46:44',86,'EXECUTED','9:71c5969e6cdd8d7b6f47cebc86d37627','dropForeignKeyConstraint baseTableName=REALM, constraintName=FK_TRAF444KK6QRKMS7N56AIWQ5Y; dropForeignKeyConstraint baseTableName=KEYCLOAK_ROLE, constraintName=FK_KJHO5LE2C0RAL09FL8CM9WFW9','',NULL,'4.33.0',NULL,NULL,'5993589629'),('map-remove-ri','keycloak','META-INF/jpa-changelog-12.0.0.xml','2025-12-17 12:46:44',87,'EXECUTED','9:a9ba7d47f065f041b7da856a81762021','dropForeignKeyConstraint baseTableName=REALM_DEFAULT_GROUPS, constraintName=FK_DEF_GROUPS_GROUP; dropForeignKeyConstraint baseTableName=REALM_DEFAULT_ROLES, constraintName=FK_H4WPD7W4HSOOLNI3H0SW7BTJE; dropForeignKeyConstraint baseTableName=CLIENT...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('map-remove-ri-13.0.0','keycloak','META-INF/jpa-changelog-13.0.0.xml','2025-12-17 12:46:45',92,'EXECUTED','9:d9be619d94af5a2f5d07b9f003543b91','dropForeignKeyConstraint baseTableName=DEFAULT_CLIENT_SCOPE, constraintName=FK_R_DEF_CLI_SCOPE_SCOPE; dropForeignKeyConstraint baseTableName=CLIENT_SCOPE_CLIENT, constraintName=FK_C_CLI_SCOPE_SCOPE; dropForeignKeyConstraint baseTableName=CLIENT_SC...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('unique-consentuser','keycloak','META-INF/jpa-changelog-25.0.0.xml','2025-12-17 12:46:46',139,'MARK_RAN','9:5857626a2ea8767e9a6c66bf3a2cb32f','customChange; dropUniqueConstraint constraintName=UK_JKUWUVD56ONTGSUHOGM8UEWRT, tableName=USER_CONSENT; addUniqueConstraint constraintName=UK_LOCAL_CONSENT, tableName=USER_CONSENT; addUniqueConstraint constraintName=UK_EXTERNAL_CONSENT, tableName=...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('unique-consentuser-edb-migration','keycloak','META-INF/jpa-changelog-25.0.0.xml','2025-12-17 12:46:46',140,'MARK_RAN','9:5857626a2ea8767e9a6c66bf3a2cb32f','customChange; dropUniqueConstraint constraintName=UK_JKUWUVD56ONTGSUHOGM8UEWRT, tableName=USER_CONSENT; addUniqueConstraint constraintName=UK_LOCAL_CONSENT, tableName=USER_CONSENT; addUniqueConstraint constraintName=UK_EXTERNAL_CONSENT, tableName=...','',NULL,'4.33.0',NULL,NULL,'5993589629'),('unique-consentuser-mysql','keycloak','META-INF/jpa-changelog-25.0.0.xml','2025-12-17 12:46:46',141,'EXECUTED','9:b79478aad5adaa1bc428e31563f55e8e','customChange; dropUniqueConstraint constraintName=UK_JKUWUVD56ONTGSUHOGM8UEWRT, tableName=USER_CONSENT; addUniqueConstraint constraintName=UK_LOCAL_CONSENT, tableName=USER_CONSENT; addUniqueConstraint constraintName=UK_EXTERNAL_CONSENT, tableName=...','',NULL,'4.33.0',NULL,NULL,'5993589629');
/*!40000 ALTER TABLE `DATABASECHANGELOG` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `DATABASECHANGELOGLOCK`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `DATABASECHANGELOGLOCK` (
  `ID` int NOT NULL,
  `LOCKED` tinyint NOT NULL,
  `LOCKGRANTED` datetime DEFAULT NULL,
  `LOCKEDBY` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `DATABASECHANGELOGLOCK` WRITE;
/*!40000 ALTER TABLE `DATABASECHANGELOGLOCK` DISABLE KEYS */;
INSERT INTO `DATABASECHANGELOGLOCK` VALUES (1,0,NULL,NULL),(1000,0,NULL,NULL);
/*!40000 ALTER TABLE `DATABASECHANGELOGLOCK` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `DEFAULT_CLIENT_SCOPE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `DEFAULT_CLIENT_SCOPE` (
  `REALM_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `SCOPE_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `DEFAULT_SCOPE` tinyint NOT NULL DEFAULT '0',
  PRIMARY KEY (`REALM_ID`,`SCOPE_ID`),
  KEY `IDX_DEFCLS_REALM` (`REALM_ID`),
  KEY `IDX_DEFCLS_SCOPE` (`SCOPE_ID`),
  CONSTRAINT `FK_R_DEF_CLI_SCOPE_REALM` FOREIGN KEY (`REALM_ID`) REFERENCES `REALM` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `DEFAULT_CLIENT_SCOPE` WRITE;
/*!40000 ALTER TABLE `DEFAULT_CLIENT_SCOPE` DISABLE KEYS */;
INSERT INTO `DEFAULT_CLIENT_SCOPE` VALUES ('4c2886d8-e4c4-4b40-9064-445c7ab90a1c','35746b28-fed9-4ae4-92a2-9b42f120f272',1),('4c2886d8-e4c4-4b40-9064-445c7ab90a1c','3ab1f10c-ca69-4372-8fb3-dfb2bc434ca4',1),('4c2886d8-e4c4-4b40-9064-445c7ab90a1c','4b7b21fc-bddf-42ce-97f0-36242b13dcfd',0),('4c2886d8-e4c4-4b40-9064-445c7ab90a1c','70413ad2-e95a-413c-af6e-d1741fe9dfe6',1),('4c2886d8-e4c4-4b40-9064-445c7ab90a1c','8a6cb980-3684-4ef9-9d40-88f23541f569',0),('4c2886d8-e4c4-4b40-9064-445c7ab90a1c','a687fb3e-0aca-481d-a697-9c23538f74c2',1),('4c2886d8-e4c4-4b40-9064-445c7ab90a1c','a84876a0-9e4b-4a9f-9369-a49e29c2cdac',1),('4c2886d8-e4c4-4b40-9064-445c7ab90a1c','bea2ab00-c561-421a-8d40-5a0dca74e80b',1),('4c2886d8-e4c4-4b40-9064-445c7ab90a1c','c80eb4e6-b9ce-4164-82e3-6e00ae132571',1),('4c2886d8-e4c4-4b40-9064-445c7ab90a1c','cb48b63f-5f46-4cef-a1bf-a1ddc25d2b24',1),('4c2886d8-e4c4-4b40-9064-445c7ab90a1c','dbf26c31-238a-4e9d-8336-3daf3d2dd7bd',0),('4c2886d8-e4c4-4b40-9064-445c7ab90a1c','df0e3ca7-158b-4a1c-b58d-9b423067475e',0),('4c2886d8-e4c4-4b40-9064-445c7ab90a1c','f8102d42-5d53-4610-bf3e-aacec43635e4',0),('5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','09dab565-4627-474a-ba4a-2c1dd91593a6',0),('5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','23ac037f-c944-45b5-93f5-bd8395bf49e0',1),('5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','4178816c-d39e-4c12-a66e-099a02427fa7',1),('5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','44b8cc1f-4adc-4465-86ad-c7ba53036921',1),('5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','5dcfd9e0-f7f3-42be-8cc7-fd9079e16baf',1),('5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','7b89f755-1355-4936-a90e-04777d87b75b',0),('5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','9f17cc3f-a321-4d62-aca0-b8a889d8629d',0),('5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','9f9d3e99-a0c3-4e59-96b4-d9177299709b',0),('5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','b679321c-7425-4c23-8834-e0209a8ae8c5',1),('5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','b82ffaf7-8ece-41c1-a3f9-00bfbe9e9739',1),('5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','be17e97a-a116-447b-bfb3-aacefeb3fc4f',1),('5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','c83dfc09-88fa-4c04-bf69-4d9afd7d6d4e',1),('5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','db7d6d32-c0a5-4b47-abf8-31ba61143504',0);
/*!40000 ALTER TABLE `DEFAULT_CLIENT_SCOPE` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `EVENT_ENTITY`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `EVENT_ENTITY` (
  `ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `CLIENT_ID` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `DETAILS_JSON` text COLLATE utf8mb4_unicode_ci,
  `ERROR` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `IP_ADDRESS` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `REALM_ID` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `SESSION_ID` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `EVENT_TIME` bigint DEFAULT NULL,
  `TYPE` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `USER_ID` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `DETAILS_JSON_LONG_VALUE` longtext CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  PRIMARY KEY (`ID`),
  KEY `IDX_EVENT_TIME` (`REALM_ID`,`EVENT_TIME`),
  KEY `IDX_EVENT_ENTITY_USER_ID_TYPE` (`USER_ID`,`TYPE`,`EVENT_TIME`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `EVENT_ENTITY` WRITE;
/*!40000 ALTER TABLE `EVENT_ENTITY` DISABLE KEYS */;
/*!40000 ALTER TABLE `EVENT_ENTITY` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `FEDERATED_IDENTITY`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `FEDERATED_IDENTITY` (
  `IDENTITY_PROVIDER` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `REALM_ID` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `FEDERATED_USER_ID` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `FEDERATED_USERNAME` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `TOKEN` text COLLATE utf8mb4_unicode_ci,
  `USER_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`IDENTITY_PROVIDER`,`USER_ID`),
  KEY `IDX_FEDIDENTITY_USER` (`USER_ID`),
  KEY `IDX_FEDIDENTITY_FEDUSER` (`FEDERATED_USER_ID`),
  CONSTRAINT `FK404288B92EF007A6` FOREIGN KEY (`USER_ID`) REFERENCES `USER_ENTITY` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `FEDERATED_IDENTITY` WRITE;
/*!40000 ALTER TABLE `FEDERATED_IDENTITY` DISABLE KEYS */;
/*!40000 ALTER TABLE `FEDERATED_IDENTITY` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `FEDERATED_USER`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `FEDERATED_USER` (
  `ID` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `STORAGE_PROVIDER_ID` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `REALM_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `FEDERATED_USER` WRITE;
/*!40000 ALTER TABLE `FEDERATED_USER` DISABLE KEYS */;
/*!40000 ALTER TABLE `FEDERATED_USER` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `FED_USER_ATTRIBUTE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `FED_USER_ATTRIBUTE` (
  `ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `NAME` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `USER_ID` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `REALM_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `STORAGE_PROVIDER_ID` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `VALUE` text COLLATE utf8mb4_unicode_ci,
  `LONG_VALUE_HASH` binary(64) DEFAULT NULL,
  `LONG_VALUE_HASH_LOWER_CASE` binary(64) DEFAULT NULL,
  `LONG_VALUE` longtext CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  PRIMARY KEY (`ID`),
  KEY `IDX_FU_ATTRIBUTE` (`USER_ID`,`REALM_ID`,`NAME`),
  KEY `FED_USER_ATTR_LONG_VALUES` (`LONG_VALUE_HASH`,`NAME`),
  KEY `FED_USER_ATTR_LONG_VALUES_LOWER_CASE` (`LONG_VALUE_HASH_LOWER_CASE`,`NAME`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `FED_USER_ATTRIBUTE` WRITE;
/*!40000 ALTER TABLE `FED_USER_ATTRIBUTE` DISABLE KEYS */;
/*!40000 ALTER TABLE `FED_USER_ATTRIBUTE` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `FED_USER_CONSENT`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `FED_USER_CONSENT` (
  `ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `CLIENT_ID` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `USER_ID` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `REALM_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `STORAGE_PROVIDER_ID` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `CREATED_DATE` bigint DEFAULT NULL,
  `LAST_UPDATED_DATE` bigint DEFAULT NULL,
  `CLIENT_STORAGE_PROVIDER` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `EXTERNAL_CLIENT_ID` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`ID`),
  KEY `IDX_FU_CONSENT` (`USER_ID`,`CLIENT_ID`),
  KEY `IDX_FU_CONSENT_RU` (`REALM_ID`,`USER_ID`),
  KEY `IDX_FU_CNSNT_EXT` (`USER_ID`,`CLIENT_STORAGE_PROVIDER`,`EXTERNAL_CLIENT_ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `FED_USER_CONSENT` WRITE;
/*!40000 ALTER TABLE `FED_USER_CONSENT` DISABLE KEYS */;
/*!40000 ALTER TABLE `FED_USER_CONSENT` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `FED_USER_CONSENT_CL_SCOPE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `FED_USER_CONSENT_CL_SCOPE` (
  `USER_CONSENT_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `SCOPE_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`USER_CONSENT_ID`,`SCOPE_ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `FED_USER_CONSENT_CL_SCOPE` WRITE;
/*!40000 ALTER TABLE `FED_USER_CONSENT_CL_SCOPE` DISABLE KEYS */;
/*!40000 ALTER TABLE `FED_USER_CONSENT_CL_SCOPE` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `FED_USER_CREDENTIAL`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `FED_USER_CREDENTIAL` (
  `ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `SALT` tinyblob,
  `TYPE` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `CREATED_DATE` bigint DEFAULT NULL,
  `USER_ID` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `REALM_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `STORAGE_PROVIDER_ID` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `USER_LABEL` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `SECRET_DATA` longtext COLLATE utf8mb4_unicode_ci,
  `CREDENTIAL_DATA` longtext COLLATE utf8mb4_unicode_ci,
  `PRIORITY` int DEFAULT NULL,
  PRIMARY KEY (`ID`),
  KEY `IDX_FU_CREDENTIAL` (`USER_ID`,`TYPE`),
  KEY `IDX_FU_CREDENTIAL_RU` (`REALM_ID`,`USER_ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `FED_USER_CREDENTIAL` WRITE;
/*!40000 ALTER TABLE `FED_USER_CREDENTIAL` DISABLE KEYS */;
/*!40000 ALTER TABLE `FED_USER_CREDENTIAL` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `FED_USER_GROUP_MEMBERSHIP`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `FED_USER_GROUP_MEMBERSHIP` (
  `GROUP_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `USER_ID` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `REALM_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `STORAGE_PROVIDER_ID` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`GROUP_ID`,`USER_ID`),
  KEY `IDX_FU_GROUP_MEMBERSHIP` (`USER_ID`,`GROUP_ID`),
  KEY `IDX_FU_GROUP_MEMBERSHIP_RU` (`REALM_ID`,`USER_ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `FED_USER_GROUP_MEMBERSHIP` WRITE;
/*!40000 ALTER TABLE `FED_USER_GROUP_MEMBERSHIP` DISABLE KEYS */;
/*!40000 ALTER TABLE `FED_USER_GROUP_MEMBERSHIP` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `FED_USER_REQUIRED_ACTION`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `FED_USER_REQUIRED_ACTION` (
  `REQUIRED_ACTION` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT ' ',
  `USER_ID` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `REALM_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `STORAGE_PROVIDER_ID` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`REQUIRED_ACTION`,`USER_ID`),
  KEY `IDX_FU_REQUIRED_ACTION` (`USER_ID`,`REQUIRED_ACTION`),
  KEY `IDX_FU_REQUIRED_ACTION_RU` (`REALM_ID`,`USER_ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `FED_USER_REQUIRED_ACTION` WRITE;
/*!40000 ALTER TABLE `FED_USER_REQUIRED_ACTION` DISABLE KEYS */;
/*!40000 ALTER TABLE `FED_USER_REQUIRED_ACTION` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `FED_USER_ROLE_MAPPING`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `FED_USER_ROLE_MAPPING` (
  `ROLE_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `USER_ID` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `REALM_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `STORAGE_PROVIDER_ID` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`ROLE_ID`,`USER_ID`),
  KEY `IDX_FU_ROLE_MAPPING` (`USER_ID`,`ROLE_ID`),
  KEY `IDX_FU_ROLE_MAPPING_RU` (`REALM_ID`,`USER_ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `FED_USER_ROLE_MAPPING` WRITE;
/*!40000 ALTER TABLE `FED_USER_ROLE_MAPPING` DISABLE KEYS */;
/*!40000 ALTER TABLE `FED_USER_ROLE_MAPPING` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `GROUP_ATTRIBUTE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `GROUP_ATTRIBUTE` (
  `ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'sybase-needs-something-here',
  `NAME` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `VALUE` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `GROUP_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`ID`),
  KEY `IDX_GROUP_ATTR_GROUP` (`GROUP_ID`),
  KEY `IDX_GROUP_ATT_BY_NAME_VALUE` (`NAME`,`VALUE`),
  CONSTRAINT `FK_GROUP_ATTRIBUTE_GROUP` FOREIGN KEY (`GROUP_ID`) REFERENCES `KEYCLOAK_GROUP` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `GROUP_ATTRIBUTE` WRITE;
/*!40000 ALTER TABLE `GROUP_ATTRIBUTE` DISABLE KEYS */;
/*!40000 ALTER TABLE `GROUP_ATTRIBUTE` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `GROUP_ROLE_MAPPING`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `GROUP_ROLE_MAPPING` (
  `ROLE_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `GROUP_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`ROLE_ID`,`GROUP_ID`),
  KEY `IDX_GROUP_ROLE_MAPP_GROUP` (`GROUP_ID`),
  CONSTRAINT `FK_GROUP_ROLE_GROUP` FOREIGN KEY (`GROUP_ID`) REFERENCES `KEYCLOAK_GROUP` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `GROUP_ROLE_MAPPING` WRITE;
/*!40000 ALTER TABLE `GROUP_ROLE_MAPPING` DISABLE KEYS */;
/*!40000 ALTER TABLE `GROUP_ROLE_MAPPING` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `IDENTITY_PROVIDER`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `IDENTITY_PROVIDER` (
  `INTERNAL_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ENABLED` tinyint NOT NULL DEFAULT '0',
  `PROVIDER_ALIAS` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `PROVIDER_ID` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `STORE_TOKEN` tinyint NOT NULL DEFAULT '0',
  `AUTHENTICATE_BY_DEFAULT` tinyint NOT NULL DEFAULT '0',
  `REALM_ID` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `ADD_TOKEN_ROLE` tinyint NOT NULL DEFAULT '1',
  `TRUST_EMAIL` tinyint NOT NULL DEFAULT '0',
  `FIRST_BROKER_LOGIN_FLOW_ID` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `POST_BROKER_LOGIN_FLOW_ID` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `PROVIDER_DISPLAY_NAME` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `LINK_ONLY` tinyint NOT NULL DEFAULT '0',
  `ORGANIZATION_ID` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `HIDE_ON_LOGIN` tinyint DEFAULT '0',
  PRIMARY KEY (`INTERNAL_ID`),
  UNIQUE KEY `UK_2DAELWNIBJI49AVXSRTUF6XJ33` (`PROVIDER_ALIAS`,`REALM_ID`),
  KEY `IDX_IDENT_PROV_REALM` (`REALM_ID`),
  KEY `IDX_IDP_REALM_ORG` (`REALM_ID`,`ORGANIZATION_ID`),
  KEY `IDX_IDP_FOR_LOGIN` (`REALM_ID`,`ENABLED`,`LINK_ONLY`,`HIDE_ON_LOGIN`,`ORGANIZATION_ID`),
  CONSTRAINT `FK2B4EBC52AE5C3B34` FOREIGN KEY (`REALM_ID`) REFERENCES `REALM` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `IDENTITY_PROVIDER` WRITE;
/*!40000 ALTER TABLE `IDENTITY_PROVIDER` DISABLE KEYS */;
/*!40000 ALTER TABLE `IDENTITY_PROVIDER` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `IDENTITY_PROVIDER_CONFIG`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `IDENTITY_PROVIDER_CONFIG` (
  `IDENTITY_PROVIDER_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `VALUE` longtext COLLATE utf8mb4_unicode_ci,
  `NAME` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`IDENTITY_PROVIDER_ID`,`NAME`),
  CONSTRAINT `FKDC4897CF864C4E43` FOREIGN KEY (`IDENTITY_PROVIDER_ID`) REFERENCES `IDENTITY_PROVIDER` (`INTERNAL_ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `IDENTITY_PROVIDER_CONFIG` WRITE;
/*!40000 ALTER TABLE `IDENTITY_PROVIDER_CONFIG` DISABLE KEYS */;
/*!40000 ALTER TABLE `IDENTITY_PROVIDER_CONFIG` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `IDENTITY_PROVIDER_MAPPER`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `IDENTITY_PROVIDER_MAPPER` (
  `ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `NAME` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `IDP_ALIAS` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `IDP_MAPPER_NAME` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `REALM_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`ID`),
  KEY `IDX_ID_PROV_MAPP_REALM` (`REALM_ID`),
  CONSTRAINT `FK_IDPM_REALM` FOREIGN KEY (`REALM_ID`) REFERENCES `REALM` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `IDENTITY_PROVIDER_MAPPER` WRITE;
/*!40000 ALTER TABLE `IDENTITY_PROVIDER_MAPPER` DISABLE KEYS */;
/*!40000 ALTER TABLE `IDENTITY_PROVIDER_MAPPER` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `IDP_MAPPER_CONFIG`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `IDP_MAPPER_CONFIG` (
  `IDP_MAPPER_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `VALUE` longtext COLLATE utf8mb4_unicode_ci,
  `NAME` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`IDP_MAPPER_ID`,`NAME`),
  CONSTRAINT `FK_IDPMCONFIG` FOREIGN KEY (`IDP_MAPPER_ID`) REFERENCES `IDENTITY_PROVIDER_MAPPER` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `IDP_MAPPER_CONFIG` WRITE;
/*!40000 ALTER TABLE `IDP_MAPPER_CONFIG` DISABLE KEYS */;
/*!40000 ALTER TABLE `IDP_MAPPER_CONFIG` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `JGROUPS_PING`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `JGROUPS_PING` (
  `address` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `cluster_name` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ip` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
  `coord` tinyint DEFAULT NULL,
  PRIMARY KEY (`address`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `JGROUPS_PING` WRITE;
/*!40000 ALTER TABLE `JGROUPS_PING` DISABLE KEYS */;
INSERT INTO `JGROUPS_PING` VALUES ('uuid://00000000-0000-0000-0000-000000000063','db5d0edff897-38714','ISPN','172.18.0.5:7800',1);
/*!40000 ALTER TABLE `JGROUPS_PING` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `KEYCLOAK_GROUP`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `KEYCLOAK_GROUP` (
  `ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `NAME` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `PARENT_GROUP` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `REALM_ID` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `TYPE` int NOT NULL,
  `DESCRIPTION` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `SIBLING_NAMES` (`REALM_ID`,`PARENT_GROUP`,`NAME`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `KEYCLOAK_GROUP` WRITE;
/*!40000 ALTER TABLE `KEYCLOAK_GROUP` DISABLE KEYS */;
/*!40000 ALTER TABLE `KEYCLOAK_GROUP` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `KEYCLOAK_ROLE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `KEYCLOAK_ROLE` (
  `ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `CLIENT_REALM_CONSTRAINT` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `CLIENT_ROLE` tinyint DEFAULT NULL,
  `DESCRIPTION` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `NAME` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `REALM_ID` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `CLIENT` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `REALM` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `UK_J3RWUVD56ONTGSUHOGM184WW2-2` (`NAME`,`CLIENT_REALM_CONSTRAINT`),
  KEY `IDX_KEYCLOAK_ROLE_CLIENT` (`CLIENT`),
  KEY `IDX_KEYCLOAK_ROLE_REALM` (`REALM`),
  CONSTRAINT `FK_6VYQFE4CN4WLQ8R6KT5VDSJ5C` FOREIGN KEY (`REALM`) REFERENCES `REALM` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `KEYCLOAK_ROLE` WRITE;
/*!40000 ALTER TABLE `KEYCLOAK_ROLE` DISABLE KEYS */;
INSERT INTO `KEYCLOAK_ROLE` VALUES ('00d5c5a8-f8df-41fd-bc8d-b9ef9a93d3e6','232e71d1-663b-4fd9-981f-0f08b56cd567',1,'${role_manage-identity-providers}','manage-identity-providers','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','232e71d1-663b-4fd9-981f-0f08b56cd567',NULL),('04291111-9b7e-473f-b432-9a3cd9a06c33','6b5c5dea-9656-4dfb-8a26-f699429eecd3',1,'${role_view-identity-providers}','view-identity-providers','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','6b5c5dea-9656-4dfb-8a26-f699429eecd3',NULL),('08e2bffc-779d-4dc5-b4fc-2b9c27ddf898','8b786b3a-c968-49cf-b5eb-5e7b36889df7',1,'${role_manage-identity-providers}','manage-identity-providers','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','8b786b3a-c968-49cf-b5eb-5e7b36889df7',NULL),('098e5980-f81f-4f7a-8eae-685dec71faa4','8b0a29f2-46dd-4f3c-8446-c0f63ad56ebc',1,'${role_view-consent}','view-consent','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','8b0a29f2-46dd-4f3c-8446-c0f63ad56ebc',NULL),('0c50da1d-eeb7-4074-b31e-f638245df45e','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581',0,'${role_default-roles}','default-roles-master','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581',NULL,NULL),('118607ab-1f22-4fb8-8b89-243907ef7f5d','232e71d1-663b-4fd9-981f-0f08b56cd567',1,'${role_query-realms}','query-realms','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','232e71d1-663b-4fd9-981f-0f08b56cd567',NULL),('1245baa1-e64c-4cba-966f-0a8069172814','232e71d1-663b-4fd9-981f-0f08b56cd567',1,'${role_view-events}','view-events','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','232e71d1-663b-4fd9-981f-0f08b56cd567',NULL),('1302a0c9-cff5-4f38-a524-387fb61ecae5','8b786b3a-c968-49cf-b5eb-5e7b36889df7',1,'${role_manage-authorization}','manage-authorization','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','8b786b3a-c968-49cf-b5eb-5e7b36889df7',NULL),('1438fd35-067c-4657-a9ae-99394fcdddf4','8b786b3a-c968-49cf-b5eb-5e7b36889df7',1,'${role_query-groups}','query-groups','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','8b786b3a-c968-49cf-b5eb-5e7b36889df7',NULL),('14ba4893-70c6-4be5-9beb-225ca62ed5ec','8b786b3a-c968-49cf-b5eb-5e7b36889df7',1,'${role_view-realm}','view-realm','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','8b786b3a-c968-49cf-b5eb-5e7b36889df7',NULL),('15af1b7b-a8c8-46e7-b327-566cbfe05346','8b0a29f2-46dd-4f3c-8446-c0f63ad56ebc',1,'${role_view-applications}','view-applications','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','8b0a29f2-46dd-4f3c-8446-c0f63ad56ebc',NULL),('1fb46915-dcb5-4305-b4d6-67798ae7498f','232e71d1-663b-4fd9-981f-0f08b56cd567',1,'${role_view-realm}','view-realm','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','232e71d1-663b-4fd9-981f-0f08b56cd567',NULL),('21082b5c-dd66-4536-b23b-c8678d7d60d5','6b5c5dea-9656-4dfb-8a26-f699429eecd3',1,'${role_view-authorization}','view-authorization','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','6b5c5dea-9656-4dfb-8a26-f699429eecd3',NULL),('251d9d2a-f9fc-4ca1-ac83-52a391885e02','6b5c5dea-9656-4dfb-8a26-f699429eecd3',1,'${role_view-clients}','view-clients','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','6b5c5dea-9656-4dfb-8a26-f699429eecd3',NULL),('26f9df1c-293a-4cf2-932a-c3a709cd1185','6b5c5dea-9656-4dfb-8a26-f699429eecd3',1,'${role_view-events}','view-events','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','6b5c5dea-9656-4dfb-8a26-f699429eecd3',NULL),('2b35e642-e0c2-4f86-a716-4d5f546d800d','1085b6f8-5701-4591-8dd5-7b68238b51b8',1,'${role_read-token}','read-token','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','1085b6f8-5701-4591-8dd5-7b68238b51b8',NULL),('2e7701b8-55a9-4fd2-8bc3-3d77c8249ce0','232e71d1-663b-4fd9-981f-0f08b56cd567',1,'${role_manage-users}','manage-users','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','232e71d1-663b-4fd9-981f-0f08b56cd567',NULL),('2fbe9598-8e10-4bba-ba61-899a02a01877','232e71d1-663b-4fd9-981f-0f08b56cd567',1,'${role_view-identity-providers}','view-identity-providers','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','232e71d1-663b-4fd9-981f-0f08b56cd567',NULL),('301e1fbd-e2dc-4aaa-a133-6ac79d60317d','232e71d1-663b-4fd9-981f-0f08b56cd567',1,'${role_manage-authorization}','manage-authorization','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','232e71d1-663b-4fd9-981f-0f08b56cd567',NULL),('31f48475-6b51-42b4-8706-db2d938ed265','8b786b3a-c968-49cf-b5eb-5e7b36889df7',1,'${role_view-identity-providers}','view-identity-providers','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','8b786b3a-c968-49cf-b5eb-5e7b36889df7',NULL),('326e598f-eaae-4da1-acbf-daafb4789421','8b786b3a-c968-49cf-b5eb-5e7b36889df7',1,'${role_query-clients}','query-clients','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','8b786b3a-c968-49cf-b5eb-5e7b36889df7',NULL),('35f07854-7a7a-4408-aca4-fad5c10cb1c6','6b5c5dea-9656-4dfb-8a26-f699429eecd3',1,'${role_create-client}','create-client','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','6b5c5dea-9656-4dfb-8a26-f699429eecd3',NULL),('38c40b55-9be6-435d-9e4a-595377c1fa77','6b5c5dea-9656-4dfb-8a26-f699429eecd3',1,'${role_query-groups}','query-groups','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','6b5c5dea-9656-4dfb-8a26-f699429eecd3',NULL),('3dd13d47-55dd-4ac7-8e2a-c428fafc3739','8b786b3a-c968-49cf-b5eb-5e7b36889df7',1,'${role_create-client}','create-client','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','8b786b3a-c968-49cf-b5eb-5e7b36889df7',NULL),('418c2f98-e334-4a20-aeb7-3cb62f58f764','232e71d1-663b-4fd9-981f-0f08b56cd567',1,'${role_query-clients}','query-clients','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','232e71d1-663b-4fd9-981f-0f08b56cd567',NULL),('4729d7e4-7b7f-4927-8af1-e74e487a1357','6b5c5dea-9656-4dfb-8a26-f699429eecd3',1,'${role_manage-users}','manage-users','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','6b5c5dea-9656-4dfb-8a26-f699429eecd3',NULL),('49d17406-d5b7-4b46-801c-7aad314461cd','8b0a29f2-46dd-4f3c-8446-c0f63ad56ebc',1,'${role_view-profile}','view-profile','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','8b0a29f2-46dd-4f3c-8446-c0f63ad56ebc',NULL),('4bac0b9c-8ffe-491a-b070-a7bad02d88e4','232e71d1-663b-4fd9-981f-0f08b56cd567',1,'${role_query-users}','query-users','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','232e71d1-663b-4fd9-981f-0f08b56cd567',NULL),('4c83be7d-a5c9-4b31-98e7-95e0320b16a5','6b5c5dea-9656-4dfb-8a26-f699429eecd3',1,'${role_manage-identity-providers}','manage-identity-providers','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','6b5c5dea-9656-4dfb-8a26-f699429eecd3',NULL),('4fba4722-371b-47c9-a0ca-3a2b05b0f197','8b0a29f2-46dd-4f3c-8446-c0f63ad56ebc',1,'${role_manage-account-links}','manage-account-links','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','8b0a29f2-46dd-4f3c-8446-c0f63ad56ebc',NULL),('4fd8b333-076a-4b26-91df-5fdc3516511b','232e71d1-663b-4fd9-981f-0f08b56cd567',1,'${role_impersonation}','impersonation','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','232e71d1-663b-4fd9-981f-0f08b56cd567',NULL),('5125458c-9476-41c9-8fb6-0543f7083773','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',0,'${role_default-roles}','default-roles-motogo','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',NULL,NULL),('51fe34be-b8a7-4b86-8ff5-d036dcfae6a1','8b0a29f2-46dd-4f3c-8446-c0f63ad56ebc',1,'${role_delete-account}','delete-account','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','8b0a29f2-46dd-4f3c-8446-c0f63ad56ebc',NULL),('542a60ff-97e0-4127-bfcf-1fd2b102251d','232e71d1-663b-4fd9-981f-0f08b56cd567',1,'${role_manage-realm}','manage-realm','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','232e71d1-663b-4fd9-981f-0f08b56cd567',NULL),('5634fdb9-5681-4755-9efe-ed28a3559530','6b5c5dea-9656-4dfb-8a26-f699429eecd3',1,'${role_query-realms}','query-realms','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','6b5c5dea-9656-4dfb-8a26-f699429eecd3',NULL),('56de3a81-f2d8-4b89-ae11-f933ea5efcd8','8b0a29f2-46dd-4f3c-8446-c0f63ad56ebc',1,'${role_view-groups}','view-groups','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','8b0a29f2-46dd-4f3c-8446-c0f63ad56ebc',NULL),('5867419e-37e5-4166-ac97-92b7c53e6eb9','8b786b3a-c968-49cf-b5eb-5e7b36889df7',1,'${role_query-users}','query-users','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','8b786b3a-c968-49cf-b5eb-5e7b36889df7',NULL),('59f35481-657d-4cc2-9fb2-6ff586c5bcb1','232e71d1-663b-4fd9-981f-0f08b56cd567',1,'${role_view-clients}','view-clients','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','232e71d1-663b-4fd9-981f-0f08b56cd567',NULL),('5ce52df1-959c-4528-9456-cec03ee9fb34','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',0,'representative motogo','representative','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',NULL,NULL),('5e86f90e-ac50-4faa-88d7-d54cb32eb4eb','6b5c5dea-9656-4dfb-8a26-f699429eecd3',1,'${role_manage-clients}','manage-clients','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','6b5c5dea-9656-4dfb-8a26-f699429eecd3',NULL),('64516292-6475-48b1-a156-399de0901f86','6b5c5dea-9656-4dfb-8a26-f699429eecd3',1,'${role_manage-realm}','manage-realm','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','6b5c5dea-9656-4dfb-8a26-f699429eecd3',NULL),('6dda702d-8bae-46a3-8308-97b6c5a82fe1','286f0039-9d0d-4f17-8b87-9cbeddaa0a89',1,'${role_view-applications}','view-applications','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','286f0039-9d0d-4f17-8b87-9cbeddaa0a89',NULL),('71cbf586-9211-43be-af4e-3f49778b8300','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581',0,'${role_create-realm}','create-realm','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581',NULL,NULL),('726fa1e9-a766-47ec-8d2e-0ff9c38fcd51','286f0039-9d0d-4f17-8b87-9cbeddaa0a89',1,'${role_delete-account}','delete-account','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','286f0039-9d0d-4f17-8b87-9cbeddaa0a89',NULL),('737be905-164a-4495-980f-d6f1efeeb69a','286f0039-9d0d-4f17-8b87-9cbeddaa0a89',1,'${role_view-profile}','view-profile','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','286f0039-9d0d-4f17-8b87-9cbeddaa0a89',NULL),('73ea483d-68da-4825-a78c-1180b07ec417','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',0,'administrador  motogo','admin','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',NULL,NULL),('75e66f19-ba72-4c66-a699-21e5c35720b5','232e71d1-663b-4fd9-981f-0f08b56cd567',1,'${role_query-groups}','query-groups','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','232e71d1-663b-4fd9-981f-0f08b56cd567',NULL),('78eb646a-3134-4823-b2ce-daf30bbfcf9b','8b786b3a-c968-49cf-b5eb-5e7b36889df7',1,'${role_view-clients}','view-clients','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','8b786b3a-c968-49cf-b5eb-5e7b36889df7',NULL),('7a5e73a4-f0c7-40dc-bb8f-01055b209de5','286f0039-9d0d-4f17-8b87-9cbeddaa0a89',1,'${role_manage-account-links}','manage-account-links','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','286f0039-9d0d-4f17-8b87-9cbeddaa0a89',NULL),('7a94500c-c30a-42ed-ae1a-19bc77834854','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',0,'${role_uma_authorization}','uma_authorization','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',NULL,NULL),('7be7e323-385b-4bce-a3cc-95eb5a1b22b2','232e71d1-663b-4fd9-981f-0f08b56cd567',1,'${role_view-authorization}','view-authorization','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','232e71d1-663b-4fd9-981f-0f08b56cd567',NULL),('807232af-7777-454a-b1cf-82a74b4ddf4b','8b786b3a-c968-49cf-b5eb-5e7b36889df7',1,'${role_realm-admin}','realm-admin','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','8b786b3a-c968-49cf-b5eb-5e7b36889df7',NULL),('87792728-d289-4c80-9c2f-c4568b5d8801','232e71d1-663b-4fd9-981f-0f08b56cd567',1,'${role_view-users}','view-users','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','232e71d1-663b-4fd9-981f-0f08b56cd567',NULL),('87cd1ca9-fabd-463a-9f38-9c7aab1d8694','6b5c5dea-9656-4dfb-8a26-f699429eecd3',1,'${role_manage-authorization}','manage-authorization','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','6b5c5dea-9656-4dfb-8a26-f699429eecd3',NULL),('8868fe7b-aadb-4011-a14b-dea7ef1df1e1','8b786b3a-c968-49cf-b5eb-5e7b36889df7',1,'${role_query-realms}','query-realms','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','8b786b3a-c968-49cf-b5eb-5e7b36889df7',NULL),('95fe56c8-1533-4dcc-9f8f-c0bd67254df3','6b5c5dea-9656-4dfb-8a26-f699429eecd3',1,'${role_view-users}','view-users','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','6b5c5dea-9656-4dfb-8a26-f699429eecd3',NULL),('9a839442-0a1f-42d6-8e91-2b8c6996e199','8b786b3a-c968-49cf-b5eb-5e7b36889df7',1,'${role_manage-users}','manage-users','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','8b786b3a-c968-49cf-b5eb-5e7b36889df7',NULL),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581',0,'${role_admin}','admin','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581',NULL,NULL),('a25dee31-81ff-4364-b45d-2a9799a7d6b9','6b5c5dea-9656-4dfb-8a26-f699429eecd3',1,'${role_query-clients}','query-clients','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','6b5c5dea-9656-4dfb-8a26-f699429eecd3',NULL),('a3fa1093-7cb2-444c-82f9-dc04cf28ca42','8b0a29f2-46dd-4f3c-8446-c0f63ad56ebc',1,'${role_manage-consent}','manage-consent','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','8b0a29f2-46dd-4f3c-8446-c0f63ad56ebc',NULL),('a8bd6cb5-d0ff-4a0a-8cf8-eac4f15a541e','6b5c5dea-9656-4dfb-8a26-f699429eecd3',1,'${role_query-users}','query-users','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','6b5c5dea-9656-4dfb-8a26-f699429eecd3',NULL),('a8ded98f-bc0d-490e-99d6-320855772b20','286f0039-9d0d-4f17-8b87-9cbeddaa0a89',1,'${role_manage-account}','manage-account','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','286f0039-9d0d-4f17-8b87-9cbeddaa0a89',NULL),('afdaca80-10d0-4be8-b18c-4933151fb78c','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581',0,'${role_offline-access}','offline_access','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581',NULL,NULL),('b2a6016f-8d43-4ab7-be9d-7119c2b71b80','8b786b3a-c968-49cf-b5eb-5e7b36889df7',1,'${role_manage-clients}','manage-clients','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','8b786b3a-c968-49cf-b5eb-5e7b36889df7',NULL),('b723d2ac-460d-49a1-abbb-077d65b47e98','8b786b3a-c968-49cf-b5eb-5e7b36889df7',1,'${role_view-users}','view-users','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','8b786b3a-c968-49cf-b5eb-5e7b36889df7',NULL),('ba23c820-628a-4fd3-aecc-07491f156fe2','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581',0,'${role_uma_authorization}','uma_authorization','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581',NULL,NULL),('bbbc95f8-5d80-4426-9b22-11e1b5c6db03','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',0,'${role_offline-access}','offline_access','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',NULL,NULL),('bc732c7b-1dda-4fcb-bbfd-d2c11138eebf','8b0a29f2-46dd-4f3c-8446-c0f63ad56ebc',1,'${role_manage-account}','manage-account','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','8b0a29f2-46dd-4f3c-8446-c0f63ad56ebc',NULL),('c0412fa5-7c3f-4928-80bd-74c26fb16dea','6b5c5dea-9656-4dfb-8a26-f699429eecd3',1,'${role_manage-events}','manage-events','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','6b5c5dea-9656-4dfb-8a26-f699429eecd3',NULL),('c08aac22-9f0a-43d1-a2f9-8ae0e672828e','232e71d1-663b-4fd9-981f-0f08b56cd567',1,'${role_manage-events}','manage-events','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','232e71d1-663b-4fd9-981f-0f08b56cd567',NULL),('caf491f4-fba2-4246-99e4-2d3f85248e64','286f0039-9d0d-4f17-8b87-9cbeddaa0a89',1,'${role_manage-consent}','manage-consent','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','286f0039-9d0d-4f17-8b87-9cbeddaa0a89',NULL),('cb13b85a-805b-4c06-bbc9-5548b3737bef','232e71d1-663b-4fd9-981f-0f08b56cd567',1,'${role_create-client}','create-client','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','232e71d1-663b-4fd9-981f-0f08b56cd567',NULL),('d2fc3a79-57f8-4447-8a23-297a4770240f','286f0039-9d0d-4f17-8b87-9cbeddaa0a89',1,'${role_view-consent}','view-consent','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','286f0039-9d0d-4f17-8b87-9cbeddaa0a89',NULL),('da6b76cf-a450-4e8a-81f6-4286c9fa2ef7','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',0,'user motogo','user','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',NULL,NULL),('daae15c7-3895-4aff-b1c9-3e753b44a99a','232e71d1-663b-4fd9-981f-0f08b56cd567',1,'${role_manage-clients}','manage-clients','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','232e71d1-663b-4fd9-981f-0f08b56cd567',NULL),('e7b39c0b-fef7-44a1-b209-df633fcbc8a5','8b786b3a-c968-49cf-b5eb-5e7b36889df7',1,'${role_manage-realm}','manage-realm','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','8b786b3a-c968-49cf-b5eb-5e7b36889df7',NULL),('eacb6d31-9ae2-4618-867c-2165b879f81d','8b786b3a-c968-49cf-b5eb-5e7b36889df7',1,'${role_view-authorization}','view-authorization','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','8b786b3a-c968-49cf-b5eb-5e7b36889df7',NULL),('ec508af1-4bcb-48f3-95af-d8e3f7efe47c','8b786b3a-c968-49cf-b5eb-5e7b36889df7',1,'${role_view-events}','view-events','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','8b786b3a-c968-49cf-b5eb-5e7b36889df7',NULL),('f13034bc-6106-42fb-950c-1467847ca692','5fc55bdb-f575-4b2e-a17c-e39bc6fb9412',1,'${role_read-token}','read-token','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','5fc55bdb-f575-4b2e-a17c-e39bc6fb9412',NULL),('f1b75db6-0157-406f-ade0-bc90cf859fda','8b786b3a-c968-49cf-b5eb-5e7b36889df7',1,'${role_manage-events}','manage-events','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','8b786b3a-c968-49cf-b5eb-5e7b36889df7',NULL),('f8772141-c42e-4f39-b42c-6652dd8a5b70','286f0039-9d0d-4f17-8b87-9cbeddaa0a89',1,'${role_view-groups}','view-groups','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','286f0039-9d0d-4f17-8b87-9cbeddaa0a89',NULL),('f89122ff-d079-4f7c-b403-4d20643f6649','6b5c5dea-9656-4dfb-8a26-f699429eecd3',1,'${role_view-realm}','view-realm','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','6b5c5dea-9656-4dfb-8a26-f699429eecd3',NULL),('f9a83f1b-cc89-4a2c-a873-0055fe745268','6b5c5dea-9656-4dfb-8a26-f699429eecd3',1,'${role_impersonation}','impersonation','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','6b5c5dea-9656-4dfb-8a26-f699429eecd3',NULL),('ff7b90a1-7585-44e1-89c5-0a39963f236d','8b786b3a-c968-49cf-b5eb-5e7b36889df7',1,'${role_impersonation}','impersonation','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','8b786b3a-c968-49cf-b5eb-5e7b36889df7',NULL);
/*!40000 ALTER TABLE `KEYCLOAK_ROLE` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `MIGRATION_MODEL`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `MIGRATION_MODEL` (
  `ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `VERSION` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `UPDATE_TIME` bigint NOT NULL DEFAULT '0',
  PRIMARY KEY (`ID`),
  UNIQUE KEY `UK_MIGRATION_UPDATE_TIME` (`UPDATE_TIME`),
  UNIQUE KEY `UK_MIGRATION_VERSION` (`VERSION`),
  KEY `IDX_UPDATE_TIME` (`UPDATE_TIME`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `MIGRATION_MODEL` WRITE;
/*!40000 ALTER TABLE `MIGRATION_MODEL` DISABLE KEYS */;
INSERT INTO `MIGRATION_MODEL` VALUES ('sqvj9','26.4.7',1765993610);
/*!40000 ALTER TABLE `MIGRATION_MODEL` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `OFFLINE_CLIENT_SESSION`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `OFFLINE_CLIENT_SESSION` (
  `USER_SESSION_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `CLIENT_ID` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `OFFLINE_FLAG` varchar(4) COLLATE utf8mb4_unicode_ci NOT NULL,
  `TIMESTAMP` int DEFAULT NULL,
  `DATA` longtext COLLATE utf8mb4_unicode_ci,
  `CLIENT_STORAGE_PROVIDER` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'local',
  `EXTERNAL_CLIENT_ID` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'local',
  `VERSION` int DEFAULT '0',
  PRIMARY KEY (`USER_SESSION_ID`,`CLIENT_ID`,`CLIENT_STORAGE_PROVIDER`,`EXTERNAL_CLIENT_ID`,`OFFLINE_FLAG`),
  KEY `IDX_OFFLINE_CSS_BY_CLIENT` (`CLIENT_ID`,`OFFLINE_FLAG`),
  KEY `IDX_OFFLINE_CSS_BY_CLIENT_STORAGE_PROVIDER` (`CLIENT_STORAGE_PROVIDER`,`EXTERNAL_CLIENT_ID`,`OFFLINE_FLAG`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `OFFLINE_CLIENT_SESSION` WRITE;
/*!40000 ALTER TABLE `OFFLINE_CLIENT_SESSION` DISABLE KEYS */;
INSERT INTO `OFFLINE_CLIENT_SESSION` VALUES ('124f803e-4f92-ed61-dc4a-19690c8d60cf','35ba46a4-b1e1-4dd0-8cff-fc98df4c2659','0',1773366192,'{\"authMethod\":\"openid-connect\",\"notes\":{\"clientId\":\"35ba46a4-b1e1-4dd0-8cff-fc98df4c2659\",\"scope\":\"openid\",\"userSessionStartedAt\":\"1773366192\",\"iss\":\"http://localhost:8080/realms/motogo\",\"startedAt\":\"1773366192\",\"level-of-authentication\":\"-1\"}}','local','local',0),('345e3838-a324-0be9-c2f7-4b4342a0f4b0','fe3360eb-3c93-4e56-99df-b8fd05dbf60a','0',1773365732,'{\"authMethod\":\"openid-connect\",\"notes\":{\"clientId\":\"fe3360eb-3c93-4e56-99df-b8fd05dbf60a\",\"userSessionStartedAt\":\"1773365732\",\"iss\":\"http://localhost:8080/realms/motogo\",\"startedAt\":\"1773365732\",\"level-of-authentication\":\"-1\"}}','local','local',0),('3581e48f-f263-8a64-4dc6-36ecbc673b69','35ba46a4-b1e1-4dd0-8cff-fc98df4c2659','0',1773365953,'{\"authMethod\":\"openid-connect\",\"notes\":{\"clientId\":\"35ba46a4-b1e1-4dd0-8cff-fc98df4c2659\",\"scope\":\"openid\",\"userSessionStartedAt\":\"1773365953\",\"iss\":\"http://localhost:8080/realms/motogo\",\"startedAt\":\"1773365953\",\"level-of-authentication\":\"-1\"}}','local','local',0),('4a3fcd4c-a88f-579a-3382-0e206449b125','35ba46a4-b1e1-4dd0-8cff-fc98df4c2659','0',1773175872,'{\"authMethod\":\"openid-connect\",\"notes\":{\"clientId\":\"35ba46a4-b1e1-4dd0-8cff-fc98df4c2659\",\"scope\":\"openid\",\"userSessionStartedAt\":\"1773175872\",\"iss\":\"http://localhost:8080/realms/motogo\",\"startedAt\":\"1773175872\",\"level-of-authentication\":\"-1\"}}','local','local',0),('51ec6ed8-3239-f1b4-44d0-ec94a6e6c733','fe3360eb-3c93-4e56-99df-b8fd05dbf60a','0',1773175644,'{\"authMethod\":\"openid-connect\",\"notes\":{\"clientId\":\"fe3360eb-3c93-4e56-99df-b8fd05dbf60a\",\"userSessionStartedAt\":\"1773175644\",\"iss\":\"http://localhost:8080/realms/motogo\",\"startedAt\":\"1773175644\",\"level-of-authentication\":\"-1\"}}','local','local',0),('5f955bf8-ca18-d35d-fa30-c22a7ca6a3c1','fe3360eb-3c93-4e56-99df-b8fd05dbf60a','0',1773175821,'{\"authMethod\":\"openid-connect\",\"notes\":{\"clientId\":\"fe3360eb-3c93-4e56-99df-b8fd05dbf60a\",\"userSessionStartedAt\":\"1773175821\",\"iss\":\"http://localhost:8080/realms/motogo\",\"startedAt\":\"1773175821\",\"level-of-authentication\":\"-1\"}}','local','local',0),('84dcf58d-5493-8167-b254-c3a211fa5eee','fe3360eb-3c93-4e56-99df-b8fd05dbf60a','0',1773172466,'{\"authMethod\":\"openid-connect\",\"notes\":{\"clientId\":\"fe3360eb-3c93-4e56-99df-b8fd05dbf60a\",\"userSessionStartedAt\":\"1773172464\",\"iss\":\"http://localhost:8080/realms/motogo\",\"startedAt\":\"1773172464\",\"level-of-authentication\":\"-1\"}}','local','local',0),('b302d552-3c2d-418a-a6ac-44839fa3638a','4c9063cc-607c-4f14-b55d-eb3d700e742c','0',1773365646,'{\"authMethod\":\"openid-connect\",\"redirectUri\":\"http://localhost:8080/admin/master/console/\",\"notes\":{\"clientId\":\"4c9063cc-607c-4f14-b55d-eb3d700e742c\",\"iss\":\"http://localhost:8080/realms/master\",\"startedAt\":\"1773365588\",\"response_type\":\"code\",\"level-of-authentication\":\"-1\",\"code_challenge_method\":\"S256\",\"nonce\":\"58b5dc4f-c389-45c6-8671-85a796982f17\",\"response_mode\":\"query\",\"scope\":\"openid\",\"userSessionStartedAt\":\"1773365588\",\"redirect_uri\":\"http://localhost:8080/admin/master/console/\",\"state\":\"24d19ff5-10fa-47bf-84bb-b1627cc72c53\",\"code_challenge\":\"UYUoEs29pBaya-2CrFfTVXZmI3T3CcYkyWXpv9HHj78\"}}','local','local',1),('d885afb4-c519-90b5-bd1b-89c4ae770545','35ba46a4-b1e1-4dd0-8cff-fc98df4c2659','0',1773366530,'{\"authMethod\":\"openid-connect\",\"notes\":{\"clientId\":\"35ba46a4-b1e1-4dd0-8cff-fc98df4c2659\",\"scope\":\"openid\",\"userSessionStartedAt\":\"1773366530\",\"iss\":\"http://localhost:8080/realms/motogo\",\"startedAt\":\"1773366530\",\"level-of-authentication\":\"-1\"}}','local','local',0),('f068a363-0216-fc6b-9181-c5528e30fa42','fe3360eb-3c93-4e56-99df-b8fd05dbf60a','0',1772993759,'{\"authMethod\":\"openid-connect\",\"notes\":{\"clientId\":\"fe3360eb-3c93-4e56-99df-b8fd05dbf60a\",\"userSessionStartedAt\":\"1772993759\",\"iss\":\"http://localhost:8080/realms/motogo\",\"startedAt\":\"1772993759\",\"level-of-authentication\":\"-1\"}}','local','local',0),('fd2751b4-55d8-83a6-e1e9-45f23a83e581','fe3360eb-3c93-4e56-99df-b8fd05dbf60a','0',1773175659,'{\"authMethod\":\"openid-connect\",\"notes\":{\"clientId\":\"fe3360eb-3c93-4e56-99df-b8fd05dbf60a\",\"userSessionStartedAt\":\"1773175659\",\"iss\":\"http://localhost:8080/realms/motogo\",\"startedAt\":\"1773175659\",\"level-of-authentication\":\"-1\"}}','local','local',0);
/*!40000 ALTER TABLE `OFFLINE_CLIENT_SESSION` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `OFFLINE_USER_SESSION`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `OFFLINE_USER_SESSION` (
  `USER_SESSION_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `USER_ID` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `REALM_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `CREATED_ON` int NOT NULL,
  `OFFLINE_FLAG` varchar(4) COLLATE utf8mb4_unicode_ci NOT NULL,
  `DATA` longtext COLLATE utf8mb4_unicode_ci,
  `LAST_SESSION_REFRESH` int NOT NULL DEFAULT '0',
  `BROKER_SESSION_ID` text COLLATE utf8mb4_unicode_ci,
  `VERSION` int DEFAULT '0',
  PRIMARY KEY (`USER_SESSION_ID`,`OFFLINE_FLAG`),
  KEY `IDX_OFFLINE_USS_BY_USER` (`USER_ID`,`REALM_ID`,`OFFLINE_FLAG`),
  KEY `IDX_OFFLINE_USS_BY_LAST_SESSION_REFRESH` (`REALM_ID`,`OFFLINE_FLAG`,`LAST_SESSION_REFRESH`),
  KEY `IDX_OFFLINE_USS_BY_BROKER_SESSION_ID` (`BROKER_SESSION_ID`(255),`REALM_ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `OFFLINE_USER_SESSION` WRITE;
/*!40000 ALTER TABLE `OFFLINE_USER_SESSION` DISABLE KEYS */;
INSERT INTO `OFFLINE_USER_SESSION` VALUES ('124f803e-4f92-ed61-dc4a-19690c8d60cf','79c7eea4-7f93-4b1e-ad06-b244ba9ca466','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',1773366192,'0','{\"ipAddress\":\"192.168.65.1\",\"authMethod\":\"openid-connect\",\"rememberMe\":false,\"started\":0,\"notes\":{\"KC_DEVICE_NOTE\":\"eyJpcEFkZHJlc3MiOiIxOTIuMTY4LjY1LjEiLCJvcyI6Ik90aGVyIiwib3NWZXJzaW9uIjoiVW5rbm93biIsImJyb3dzZXIiOiJnby1yZXN0eS8yLjcuMCIsImRldmljZSI6Ik90aGVyIiwibGFzdEFjY2VzcyI6MCwibW9iaWxlIjpmYWxzZX0=\",\"authenticators-completed\":\"{\\\"d96442ea-d0af-4eb9-8539-61bbb93578dd\\\":1773366192,\\\"cfa83450-7fb9-4a91-9fde-777d8d681f64\\\":1773366192}\"},\"state\":\"LOGGED_IN\"}',1773366192,NULL,0),('345e3838-a324-0be9-c2f7-4b4342a0f4b0','d62612e2-e122-441d-8e71-7fee81f4c53c','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',1773365732,'0','{\"ipAddress\":\"192.168.65.1\",\"authMethod\":\"openid-connect\",\"rememberMe\":false,\"started\":0,\"notes\":{\"KC_DEVICE_NOTE\":\"eyJpcEFkZHJlc3MiOiIxOTIuMTY4LjY1LjEiLCJvcyI6Ik90aGVyIiwib3NWZXJzaW9uIjoiVW5rbm93biIsImJyb3dzZXIiOiJnby1yZXN0eS8yLjcuMCIsImRldmljZSI6Ik90aGVyIiwibGFzdEFjY2VzcyI6MCwibW9iaWxlIjpmYWxzZX0=\",\"authenticators-completed\":\"{\\\"d96442ea-d0af-4eb9-8539-61bbb93578dd\\\":1773365732,\\\"cfa83450-7fb9-4a91-9fde-777d8d681f64\\\":1773365732}\"},\"state\":\"LOGGED_IN\"}',1773365732,NULL,0),('3581e48f-f263-8a64-4dc6-36ecbc673b69','d62612e2-e122-441d-8e71-7fee81f4c53c','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',1773365953,'0','{\"ipAddress\":\"192.168.65.1\",\"authMethod\":\"openid-connect\",\"rememberMe\":false,\"started\":0,\"notes\":{\"KC_DEVICE_NOTE\":\"eyJpcEFkZHJlc3MiOiIxOTIuMTY4LjY1LjEiLCJvcyI6Ik90aGVyIiwib3NWZXJzaW9uIjoiVW5rbm93biIsImJyb3dzZXIiOiJnby1yZXN0eS8yLjcuMCIsImRldmljZSI6Ik90aGVyIiwibGFzdEFjY2VzcyI6MCwibW9iaWxlIjpmYWxzZX0=\",\"authenticators-completed\":\"{\\\"d96442ea-d0af-4eb9-8539-61bbb93578dd\\\":1773365953,\\\"cfa83450-7fb9-4a91-9fde-777d8d681f64\\\":1773365953}\"},\"state\":\"LOGGED_IN\"}',1773365953,NULL,0),('4a3fcd4c-a88f-579a-3382-0e206449b125','d62612e2-e122-441d-8e71-7fee81f4c53c','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',1773175872,'0','{\"ipAddress\":\"172.18.0.1\",\"authMethod\":\"openid-connect\",\"rememberMe\":false,\"started\":0,\"notes\":{\"KC_DEVICE_NOTE\":\"eyJpcEFkZHJlc3MiOiIxNzIuMTguMC4xIiwib3MiOiJPdGhlciIsIm9zVmVyc2lvbiI6IlVua25vd24iLCJicm93c2VyIjoiZ28tcmVzdHkvMi43LjAiLCJkZXZpY2UiOiJPdGhlciIsImxhc3RBY2Nlc3MiOjAsIm1vYmlsZSI6ZmFsc2V9\",\"authenticators-completed\":\"{\\\"d96442ea-d0af-4eb9-8539-61bbb93578dd\\\":1773175872,\\\"cfa83450-7fb9-4a91-9fde-777d8d681f64\\\":1773175872}\"},\"state\":\"LOGGED_IN\"}',1773175872,NULL,0),('51ec6ed8-3239-f1b4-44d0-ec94a6e6c733','d62612e2-e122-441d-8e71-7fee81f4c53c','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',1773175644,'0','{\"ipAddress\":\"172.18.0.1\",\"authMethod\":\"openid-connect\",\"rememberMe\":false,\"started\":0,\"notes\":{\"KC_DEVICE_NOTE\":\"eyJpcEFkZHJlc3MiOiIxNzIuMTguMC4xIiwib3MiOiJPdGhlciIsIm9zVmVyc2lvbiI6IlVua25vd24iLCJicm93c2VyIjoiZ28tcmVzdHkvMi43LjAiLCJkZXZpY2UiOiJPdGhlciIsImxhc3RBY2Nlc3MiOjAsIm1vYmlsZSI6ZmFsc2V9\",\"authenticators-completed\":\"{\\\"d96442ea-d0af-4eb9-8539-61bbb93578dd\\\":1773175643,\\\"cfa83450-7fb9-4a91-9fde-777d8d681f64\\\":1773175643}\"},\"state\":\"LOGGED_IN\"}',1773175644,NULL,0),('5f955bf8-ca18-d35d-fa30-c22a7ca6a3c1','d62612e2-e122-441d-8e71-7fee81f4c53c','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',1773175821,'0','{\"ipAddress\":\"172.18.0.1\",\"authMethod\":\"openid-connect\",\"rememberMe\":false,\"started\":0,\"notes\":{\"KC_DEVICE_NOTE\":\"eyJpcEFkZHJlc3MiOiIxNzIuMTguMC4xIiwib3MiOiJPdGhlciIsIm9zVmVyc2lvbiI6IlVua25vd24iLCJicm93c2VyIjoiZ28tcmVzdHkvMi43LjAiLCJkZXZpY2UiOiJPdGhlciIsImxhc3RBY2Nlc3MiOjAsIm1vYmlsZSI6ZmFsc2V9\",\"authenticators-completed\":\"{\\\"d96442ea-d0af-4eb9-8539-61bbb93578dd\\\":1773175821,\\\"cfa83450-7fb9-4a91-9fde-777d8d681f64\\\":1773175821}\"},\"state\":\"LOGGED_IN\"}',1773175821,NULL,0),('84dcf58d-5493-8167-b254-c3a211fa5eee','d62612e2-e122-441d-8e71-7fee81f4c53c','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',1773172464,'0','{\"ipAddress\":\"172.18.0.1\",\"authMethod\":\"openid-connect\",\"rememberMe\":false,\"started\":0,\"notes\":{\"KC_DEVICE_NOTE\":\"eyJpcEFkZHJlc3MiOiIxNzIuMTguMC4xIiwib3MiOiJPdGhlciIsIm9zVmVyc2lvbiI6IlVua25vd24iLCJicm93c2VyIjoiZ28tcmVzdHkvMi43LjAiLCJkZXZpY2UiOiJPdGhlciIsImxhc3RBY2Nlc3MiOjAsIm1vYmlsZSI6ZmFsc2V9\",\"authenticators-completed\":\"{\\\"d96442ea-d0af-4eb9-8539-61bbb93578dd\\\":1773172462,\\\"cfa83450-7fb9-4a91-9fde-777d8d681f64\\\":1773172463}\"},\"state\":\"LOGGED_IN\"}',1773172466,NULL,0),('b302d552-3c2d-418a-a6ac-44839fa3638a','19b6d072-21e8-40f9-9d22-711e5814459f','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581',1773365588,'0','{\"ipAddress\":\"192.168.65.1\",\"authMethod\":\"openid-connect\",\"rememberMe\":false,\"started\":0,\"notes\":{\"KC_DEVICE_NOTE\":\"eyJpcEFkZHJlc3MiOiIxOTIuMTY4LjY1LjEiLCJvcyI6Ik1hYyBPUyBYIiwib3NWZXJzaW9uIjoiMTAuMTUuNyIsImJyb3dzZXIiOiJDaHJvbWUvMTQ2LjAuMCIsImRldmljZSI6Ik1hYyIsImxhc3RBY2Nlc3MiOjAsIm1vYmlsZSI6ZmFsc2V9\",\"AUTH_TIME\":\"1773365588\",\"authenticators-completed\":\"{\\\"fc09002b-d368-4845-ad60-7aea0d3765cc\\\":1773365588}\"},\"state\":\"LOGGED_IN\"}',1773365646,NULL,1),('d885afb4-c519-90b5-bd1b-89c4ae770545','fb214e14-d12c-451e-aca8-c62b09f77b11','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',1773366530,'0','{\"ipAddress\":\"192.168.65.1\",\"authMethod\":\"openid-connect\",\"rememberMe\":false,\"started\":0,\"notes\":{\"KC_DEVICE_NOTE\":\"eyJpcEFkZHJlc3MiOiIxOTIuMTY4LjY1LjEiLCJvcyI6Ik90aGVyIiwib3NWZXJzaW9uIjoiVW5rbm93biIsImJyb3dzZXIiOiJnby1yZXN0eS8yLjcuMCIsImRldmljZSI6Ik90aGVyIiwibGFzdEFjY2VzcyI6MCwibW9iaWxlIjpmYWxzZX0=\",\"authenticators-completed\":\"{\\\"d96442ea-d0af-4eb9-8539-61bbb93578dd\\\":1773366530,\\\"cfa83450-7fb9-4a91-9fde-777d8d681f64\\\":1773366530}\"},\"state\":\"LOGGED_IN\"}',1773366530,NULL,0),('f068a363-0216-fc6b-9181-c5528e30fa42','d62612e2-e122-441d-8e71-7fee81f4c53c','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',1772993759,'0','{\"ipAddress\":\"172.18.0.1\",\"authMethod\":\"openid-connect\",\"rememberMe\":false,\"started\":0,\"notes\":{\"KC_DEVICE_NOTE\":\"eyJpcEFkZHJlc3MiOiIxNzIuMTguMC4xIiwib3MiOiJPdGhlciIsIm9zVmVyc2lvbiI6IlVua25vd24iLCJicm93c2VyIjoiZ28tcmVzdHkvMi43LjAiLCJkZXZpY2UiOiJPdGhlciIsImxhc3RBY2Nlc3MiOjAsIm1vYmlsZSI6ZmFsc2V9\",\"authenticators-completed\":\"{\\\"d96442ea-d0af-4eb9-8539-61bbb93578dd\\\":1772993759,\\\"cfa83450-7fb9-4a91-9fde-777d8d681f64\\\":1772993759}\"},\"state\":\"LOGGED_IN\"}',1772993759,NULL,0),('fd2751b4-55d8-83a6-e1e9-45f23a83e581','d62612e2-e122-441d-8e71-7fee81f4c53c','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',1773175659,'0','{\"ipAddress\":\"172.18.0.1\",\"authMethod\":\"openid-connect\",\"rememberMe\":false,\"started\":0,\"notes\":{\"KC_DEVICE_NOTE\":\"eyJpcEFkZHJlc3MiOiIxNzIuMTguMC4xIiwib3MiOiJPdGhlciIsIm9zVmVyc2lvbiI6IlVua25vd24iLCJicm93c2VyIjoiZ28tcmVzdHkvMi43LjAiLCJkZXZpY2UiOiJPdGhlciIsImxhc3RBY2Nlc3MiOjAsIm1vYmlsZSI6ZmFsc2V9\",\"authenticators-completed\":\"{\\\"d96442ea-d0af-4eb9-8539-61bbb93578dd\\\":1773175658,\\\"cfa83450-7fb9-4a91-9fde-777d8d681f64\\\":1773175658}\"},\"state\":\"LOGGED_IN\"}',1773175659,NULL,0);
/*!40000 ALTER TABLE `OFFLINE_USER_SESSION` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `ORG`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ORG` (
  `ID` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ENABLED` tinyint NOT NULL,
  `REALM_ID` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `GROUP_ID` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `NAME` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `DESCRIPTION` text COLLATE utf8mb4_unicode_ci,
  `ALIAS` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `REDIRECT_URL` text COLLATE utf8mb4_unicode_ci,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `UK_ORG_NAME` (`REALM_ID`,`NAME`),
  UNIQUE KEY `UK_ORG_GROUP` (`GROUP_ID`),
  UNIQUE KEY `UK_ORG_ALIAS` (`REALM_ID`,`ALIAS`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `ORG` WRITE;
/*!40000 ALTER TABLE `ORG` DISABLE KEYS */;
/*!40000 ALTER TABLE `ORG` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `ORG_DOMAIN`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ORG_DOMAIN` (
  `ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `NAME` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `VERIFIED` tinyint NOT NULL,
  `ORG_ID` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`ID`,`NAME`),
  KEY `IDX_ORG_DOMAIN_ORG_ID` (`ORG_ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `ORG_DOMAIN` WRITE;
/*!40000 ALTER TABLE `ORG_DOMAIN` DISABLE KEYS */;
/*!40000 ALTER TABLE `ORG_DOMAIN` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `POLICY_CONFIG`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `POLICY_CONFIG` (
  `POLICY_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `NAME` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `VALUE` longtext COLLATE utf8mb4_unicode_ci,
  PRIMARY KEY (`POLICY_ID`,`NAME`),
  CONSTRAINT `FKDC34197CF864C4E43` FOREIGN KEY (`POLICY_ID`) REFERENCES `RESOURCE_SERVER_POLICY` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `POLICY_CONFIG` WRITE;
/*!40000 ALTER TABLE `POLICY_CONFIG` DISABLE KEYS */;
/*!40000 ALTER TABLE `POLICY_CONFIG` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `PROTOCOL_MAPPER`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `PROTOCOL_MAPPER` (
  `ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `NAME` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `PROTOCOL` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `PROTOCOL_MAPPER_NAME` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `CLIENT_ID` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `CLIENT_SCOPE_ID` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`ID`),
  KEY `IDX_PROTOCOL_MAPPER_CLIENT` (`CLIENT_ID`),
  KEY `IDX_CLSCOPE_PROTMAP` (`CLIENT_SCOPE_ID`),
  CONSTRAINT `FK_CLI_SCOPE_MAPPER` FOREIGN KEY (`CLIENT_SCOPE_ID`) REFERENCES `CLIENT_SCOPE` (`ID`),
  CONSTRAINT `FK_PCM_REALM` FOREIGN KEY (`CLIENT_ID`) REFERENCES `CLIENT` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `PROTOCOL_MAPPER` WRITE;
/*!40000 ALTER TABLE `PROTOCOL_MAPPER` DISABLE KEYS */;
INSERT INTO `PROTOCOL_MAPPER` VALUES ('0543c8c9-67f9-47e1-bc8f-c76e3bf21e64','family name','openid-connect','oidc-usermodel-attribute-mapper',NULL,'3ab1f10c-ca69-4372-8fb3-dfb2bc434ca4'),('06dac6f6-8866-4b4e-9ad3-e102f86ff548','middle name','openid-connect','oidc-usermodel-attribute-mapper',NULL,'b679321c-7425-4c23-8834-e0209a8ae8c5'),('07d77d71-061e-4eb7-92e6-c468f90657dc','Client ID','openid-connect','oidc-usersessionmodel-note-mapper',NULL,'1daedf09-7b9a-4ed1-b0a7-db5dec46c9b5'),('0d54ff44-d2b0-4217-b491-a75e411e68eb','nickname','openid-connect','oidc-usermodel-attribute-mapper',NULL,'b679321c-7425-4c23-8834-e0209a8ae8c5'),('0ef20646-a5a9-4c30-81e7-5bcc5836abe2','nickname','openid-connect','oidc-usermodel-attribute-mapper',NULL,'3ab1f10c-ca69-4372-8fb3-dfb2bc434ca4'),('1730c42a-1d72-411a-984a-9a819d4060c2','updated at','openid-connect','oidc-usermodel-attribute-mapper',NULL,'b679321c-7425-4c23-8834-e0209a8ae8c5'),('17387ac5-5633-4984-a987-f193e600f9b0','upn','openid-connect','oidc-usermodel-attribute-mapper',NULL,'df0e3ca7-158b-4a1c-b58d-9b423067475e'),('185e5efd-282f-4a83-b52c-40e9c72a0134','gender','openid-connect','oidc-usermodel-attribute-mapper',NULL,'3ab1f10c-ca69-4372-8fb3-dfb2bc434ca4'),('1aab8d5c-7edb-45c7-b4e1-87b9364bfe87','Client Host','openid-connect','oidc-usersessionmodel-note-mapper',NULL,'1daedf09-7b9a-4ed1-b0a7-db5dec46c9b5'),('1fbc6ca0-99ae-44d2-83fd-f650feef91d4','username','openid-connect','oidc-usermodel-attribute-mapper',NULL,'3ab1f10c-ca69-4372-8fb3-dfb2bc434ca4'),('242421cc-38ad-40d2-b6f9-9ed40e42390d','profile','openid-connect','oidc-usermodel-attribute-mapper',NULL,'3ab1f10c-ca69-4372-8fb3-dfb2bc434ca4'),('25206647-5675-4ffe-b195-893e3c62bd55','sub','openid-connect','oidc-sub-mapper',NULL,'4178816c-d39e-4c12-a66e-099a02427fa7'),('277f7ed0-4ffb-4efc-8c00-2a4354362da2','client roles','openid-connect','oidc-usermodel-client-role-mapper',NULL,'a687fb3e-0aca-481d-a697-9c23538f74c2'),('280c6f3a-cbe8-44d4-a630-40166862a4d3','Client IP Address','openid-connect','oidc-usersessionmodel-note-mapper',NULL,'1daedf09-7b9a-4ed1-b0a7-db5dec46c9b5'),('2a6e41cd-993b-424b-8858-564276c6ef8c','locale','openid-connect','oidc-usermodel-attribute-mapper',NULL,'3ab1f10c-ca69-4372-8fb3-dfb2bc434ca4'),('335a8d35-b400-4b7a-a306-01e53e9ff6ff','organization','saml','saml-organization-membership-mapper',NULL,'c80eb4e6-b9ce-4164-82e3-6e00ae132571'),('388800ab-5b9e-435e-ad59-d44d200ac5e1','realm roles','openid-connect','oidc-usermodel-realm-role-mapper',NULL,'44b8cc1f-4adc-4465-86ad-c7ba53036921'),('38a846a4-579c-436c-8800-02a1a6272c07','email verified','openid-connect','oidc-usermodel-property-mapper',NULL,'35746b28-fed9-4ae4-92a2-9b42f120f272'),('38cd59f8-9311-4e55-a636-cc85ea77eb3a','audience resolve','openid-connect','oidc-audience-resolve-mapper',NULL,'a687fb3e-0aca-481d-a697-9c23538f74c2'),('39808fda-1527-4cbd-8220-d698f76d9556','profile','openid-connect','oidc-usermodel-attribute-mapper',NULL,'b679321c-7425-4c23-8834-e0209a8ae8c5'),('3eb3fcc9-f7a9-4929-92d0-e2f69209771d','allowed web origins','openid-connect','oidc-allowed-origins-mapper',NULL,'be17e97a-a116-447b-bfb3-aacefeb3fc4f'),('43814bec-d0ae-426d-b3ef-95464ea1be67','email','openid-connect','oidc-usermodel-attribute-mapper',NULL,'35746b28-fed9-4ae4-92a2-9b42f120f272'),('4fbfb957-8785-4380-8149-5be64cba1c01','locale','openid-connect','oidc-usermodel-attribute-mapper','4c9063cc-607c-4f14-b55d-eb3d700e742c',NULL),('5566cade-a4a4-4a22-b524-258986841372','Client IP Address','openid-connect','oidc-usersessionmodel-note-mapper',NULL,'f06d6bc3-3c00-4cf3-a474-e5becd828d7a'),('5ee47111-68f0-4fe7-9686-247d1f92b5c7','picture','openid-connect','oidc-usermodel-attribute-mapper',NULL,'3ab1f10c-ca69-4372-8fb3-dfb2bc434ca4'),('5f495607-3a9b-4850-acac-eefbe989f02e','zoneinfo','openid-connect','oidc-usermodel-attribute-mapper',NULL,'3ab1f10c-ca69-4372-8fb3-dfb2bc434ca4'),('66a447c6-cc42-4b19-853f-6b66681d6082','middle name','openid-connect','oidc-usermodel-attribute-mapper',NULL,'3ab1f10c-ca69-4372-8fb3-dfb2bc434ca4'),('6ce94053-6071-44cc-ab4d-0be1d37144f3','organization','saml','saml-organization-membership-mapper',NULL,'b82ffaf7-8ece-41c1-a3f9-00bfbe9e9739'),('6f2e805e-d911-440e-92d0-490bdef3cf2b','website','openid-connect','oidc-usermodel-attribute-mapper',NULL,'3ab1f10c-ca69-4372-8fb3-dfb2bc434ca4'),('71010a33-ce8f-41aa-a090-f042de837139','audience resolve','openid-connect','oidc-audience-resolve-mapper',NULL,'44b8cc1f-4adc-4465-86ad-c7ba53036921'),('71e65459-0242-489b-a3b5-07e3faa7304c','phone number verified','openid-connect','oidc-usermodel-attribute-mapper',NULL,'9f9d3e99-a0c3-4e59-96b4-d9177299709b'),('73c2a0a6-3325-43f5-86dd-1da5b6b84510','acr loa level','openid-connect','oidc-acr-mapper',NULL,'c83dfc09-88fa-4c04-bf69-4d9afd7d6d4e'),('75f8d3f9-1d00-48c2-9425-c82b140ef2af','phone number verified','openid-connect','oidc-usermodel-attribute-mapper',NULL,'8a6cb980-3684-4ef9-9d40-88f23541f569'),('76be9620-7989-440c-8c68-16da57b40d19','updated at','openid-connect','oidc-usermodel-attribute-mapper',NULL,'3ab1f10c-ca69-4372-8fb3-dfb2bc434ca4'),('7a99cbf4-af03-4ceb-80cd-b540d2876037','address','openid-connect','oidc-address-mapper',NULL,'db7d6d32-c0a5-4b47-abf8-31ba61143504'),('7e9aa70f-9359-4fd8-a047-b52720f65c50','family name','openid-connect','oidc-usermodel-attribute-mapper',NULL,'b679321c-7425-4c23-8834-e0209a8ae8c5'),('81ded2d3-54f8-431c-a864-cd8ae9485e88','full name','openid-connect','oidc-full-name-mapper',NULL,'3ab1f10c-ca69-4372-8fb3-dfb2bc434ca4'),('83d5038c-9b1a-432c-be40-1cfd8ec8c8e8','groups','openid-connect','oidc-usermodel-realm-role-mapper',NULL,'7b89f755-1355-4936-a90e-04777d87b75b'),('881ccc55-c62b-4e4a-b95d-7aeb8f73fe80','full name','openid-connect','oidc-full-name-mapper',NULL,'b679321c-7425-4c23-8834-e0209a8ae8c5'),('88556cd3-317c-4ea8-ab61-09618a18da2b','organization','openid-connect','oidc-organization-membership-mapper',NULL,'09dab565-4627-474a-ba4a-2c1dd91593a6'),('90bc110b-9196-46a1-92f7-8f80d3ff465c','email','openid-connect','oidc-usermodel-attribute-mapper',NULL,'23ac037f-c944-45b5-93f5-bd8395bf49e0'),('91182dc0-9ad6-43a6-b15b-b062821d6d97','role list','saml','saml-role-list-mapper',NULL,'cb48b63f-5f46-4cef-a1bf-a1ddc25d2b24'),('9b115eee-e55c-489d-b203-85013572f223','email verified','openid-connect','oidc-usermodel-property-mapper',NULL,'23ac037f-c944-45b5-93f5-bd8395bf49e0'),('9f6277a2-296e-43ac-98df-c152cec85e9f','allowed web origins','openid-connect','oidc-allowed-origins-mapper',NULL,'a84876a0-9e4b-4a9f-9369-a49e29c2cdac'),('a15f4aa2-6d01-4636-8f96-fc33234b1956','phone number','openid-connect','oidc-usermodel-attribute-mapper',NULL,'9f9d3e99-a0c3-4e59-96b4-d9177299709b'),('a1f1a727-50e7-4f79-b239-87ac4fed7ad2','Client ID','openid-connect','oidc-usersessionmodel-note-mapper',NULL,'f06d6bc3-3c00-4cf3-a474-e5becd828d7a'),('a29fbd5d-e05d-4cee-bff5-81b1359085b3','birthdate','openid-connect','oidc-usermodel-attribute-mapper',NULL,'3ab1f10c-ca69-4372-8fb3-dfb2bc434ca4'),('a454fc4c-9349-4c5c-b7bd-77d9fefcfa8c','sub','openid-connect','oidc-sub-mapper',NULL,'70413ad2-e95a-413c-af6e-d1741fe9dfe6'),('a7fa25ad-c21f-433c-bc89-3b7489b2f745','username','openid-connect','oidc-usermodel-attribute-mapper',NULL,'b679321c-7425-4c23-8834-e0209a8ae8c5'),('aa852eb2-6df6-49ff-943a-f8b86e67e992','organization','openid-connect','oidc-organization-membership-mapper',NULL,'4b7b21fc-bddf-42ce-97f0-36242b13dcfd'),('aba0fb9f-bf9d-4413-8d0f-1d85bbf0383f','client roles','openid-connect','oidc-usermodel-client-role-mapper',NULL,'44b8cc1f-4adc-4465-86ad-c7ba53036921'),('b6401f2e-62e8-45cc-bbb4-3af887fdb464','audience resolve','openid-connect','oidc-audience-resolve-mapper','d2eb4296-2046-420d-beb2-25cf9e05b715',NULL),('b671eb28-0bd1-4f17-9b73-9e20fe1cc45a','locale','openid-connect','oidc-usermodel-attribute-mapper',NULL,'b679321c-7425-4c23-8834-e0209a8ae8c5'),('bb0728a9-4774-4d72-86cc-4edea3d8f789','realm roles','openid-connect','oidc-usermodel-realm-role-mapper',NULL,'a687fb3e-0aca-481d-a697-9c23538f74c2'),('bbfd59ca-4897-454c-b6c7-4f33dd935efe','auth_time','openid-connect','oidc-usersessionmodel-note-mapper',NULL,'4178816c-d39e-4c12-a66e-099a02427fa7'),('c0ddb572-44f1-4264-bf8b-a4ab876aeaf4','picture','openid-connect','oidc-usermodel-attribute-mapper',NULL,'b679321c-7425-4c23-8834-e0209a8ae8c5'),('c862c840-8335-41cf-aef3-318a023b586c','birthdate','openid-connect','oidc-usermodel-attribute-mapper',NULL,'b679321c-7425-4c23-8834-e0209a8ae8c5'),('cc66e6d1-871c-48b6-9bd7-42d400903043','audience resolve','openid-connect','oidc-audience-resolve-mapper','13fd5765-ac20-4558-ae2c-cf8e2df4c22c',NULL),('ce536f16-d5dc-4688-8913-2b0f34d7671f','address','openid-connect','oidc-address-mapper',NULL,'f8102d42-5d53-4610-bf3e-aacec43635e4'),('d2f52858-16d3-425b-b7ce-7270a0ee3e78','locale','openid-connect','oidc-usermodel-attribute-mapper','3c6f327b-8018-44e9-8540-511170baeb2f',NULL),('da0f3ba9-f24d-43c6-8e21-42f77f2f846e','given name','openid-connect','oidc-usermodel-attribute-mapper',NULL,'b679321c-7425-4c23-8834-e0209a8ae8c5'),('dde2d6fa-647f-434f-847a-0c7c915c9e12','Client Host','openid-connect','oidc-usersessionmodel-note-mapper',NULL,'f06d6bc3-3c00-4cf3-a474-e5becd828d7a'),('e073d946-8b03-4f79-b4ed-de5fcad5a543','given name','openid-connect','oidc-usermodel-attribute-mapper',NULL,'3ab1f10c-ca69-4372-8fb3-dfb2bc434ca4'),('e2cffa4c-6764-49b1-bf78-9d3a7ef656e2','acr loa level','openid-connect','oidc-acr-mapper',NULL,'bea2ab00-c561-421a-8d40-5a0dca74e80b'),('e47981cc-6126-4f12-983d-e5af93ad8261','website','openid-connect','oidc-usermodel-attribute-mapper',NULL,'b679321c-7425-4c23-8834-e0209a8ae8c5'),('ec976618-f1b8-4285-8f69-313d975ff126','gender','openid-connect','oidc-usermodel-attribute-mapper',NULL,'b679321c-7425-4c23-8834-e0209a8ae8c5'),('ed02d7d5-3f23-4c6c-aa11-744bb3a03002','phone number','openid-connect','oidc-usermodel-attribute-mapper',NULL,'8a6cb980-3684-4ef9-9d40-88f23541f569'),('ed4234fb-cf69-414b-8fde-bfb47e0d5575','groups','openid-connect','oidc-usermodel-realm-role-mapper',NULL,'df0e3ca7-158b-4a1c-b58d-9b423067475e'),('f01049bb-0c22-4527-a7f5-1c93bd4235ff','zoneinfo','openid-connect','oidc-usermodel-attribute-mapper',NULL,'b679321c-7425-4c23-8834-e0209a8ae8c5'),('f11bc60a-0f09-4588-8415-7e44aa362618','role list','saml','saml-role-list-mapper',NULL,'5dcfd9e0-f7f3-42be-8cc7-fd9079e16baf'),('f76a4e44-b836-4d58-80f8-1012c94278f4','auth_time','openid-connect','oidc-usersessionmodel-note-mapper',NULL,'70413ad2-e95a-413c-af6e-d1741fe9dfe6'),('f883dadf-2b2d-4feb-ae39-746e6007a128','upn','openid-connect','oidc-usermodel-attribute-mapper',NULL,'7b89f755-1355-4936-a90e-04777d87b75b');
/*!40000 ALTER TABLE `PROTOCOL_MAPPER` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `PROTOCOL_MAPPER_CONFIG`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `PROTOCOL_MAPPER_CONFIG` (
  `PROTOCOL_MAPPER_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `VALUE` longtext COLLATE utf8mb4_unicode_ci,
  `NAME` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`PROTOCOL_MAPPER_ID`,`NAME`),
  CONSTRAINT `FK_PMCONFIG` FOREIGN KEY (`PROTOCOL_MAPPER_ID`) REFERENCES `PROTOCOL_MAPPER` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `PROTOCOL_MAPPER_CONFIG` WRITE;
/*!40000 ALTER TABLE `PROTOCOL_MAPPER_CONFIG` DISABLE KEYS */;
INSERT INTO `PROTOCOL_MAPPER_CONFIG` VALUES ('0543c8c9-67f9-47e1-bc8f-c76e3bf21e64','true','access.token.claim'),('0543c8c9-67f9-47e1-bc8f-c76e3bf21e64','family_name','claim.name'),('0543c8c9-67f9-47e1-bc8f-c76e3bf21e64','true','id.token.claim'),('0543c8c9-67f9-47e1-bc8f-c76e3bf21e64','true','introspection.token.claim'),('0543c8c9-67f9-47e1-bc8f-c76e3bf21e64','String','jsonType.label'),('0543c8c9-67f9-47e1-bc8f-c76e3bf21e64','lastName','user.attribute'),('0543c8c9-67f9-47e1-bc8f-c76e3bf21e64','true','userinfo.token.claim'),('06dac6f6-8866-4b4e-9ad3-e102f86ff548','true','access.token.claim'),('06dac6f6-8866-4b4e-9ad3-e102f86ff548','middle_name','claim.name'),('06dac6f6-8866-4b4e-9ad3-e102f86ff548','true','id.token.claim'),('06dac6f6-8866-4b4e-9ad3-e102f86ff548','true','introspection.token.claim'),('06dac6f6-8866-4b4e-9ad3-e102f86ff548','String','jsonType.label'),('06dac6f6-8866-4b4e-9ad3-e102f86ff548','middleName','user.attribute'),('06dac6f6-8866-4b4e-9ad3-e102f86ff548','true','userinfo.token.claim'),('07d77d71-061e-4eb7-92e6-c468f90657dc','true','access.token.claim'),('07d77d71-061e-4eb7-92e6-c468f90657dc','client_id','claim.name'),('07d77d71-061e-4eb7-92e6-c468f90657dc','true','id.token.claim'),('07d77d71-061e-4eb7-92e6-c468f90657dc','true','introspection.token.claim'),('07d77d71-061e-4eb7-92e6-c468f90657dc','String','jsonType.label'),('07d77d71-061e-4eb7-92e6-c468f90657dc','client_id','user.session.note'),('0d54ff44-d2b0-4217-b491-a75e411e68eb','true','access.token.claim'),('0d54ff44-d2b0-4217-b491-a75e411e68eb','nickname','claim.name'),('0d54ff44-d2b0-4217-b491-a75e411e68eb','true','id.token.claim'),('0d54ff44-d2b0-4217-b491-a75e411e68eb','true','introspection.token.claim'),('0d54ff44-d2b0-4217-b491-a75e411e68eb','String','jsonType.label'),('0d54ff44-d2b0-4217-b491-a75e411e68eb','nickname','user.attribute'),('0d54ff44-d2b0-4217-b491-a75e411e68eb','true','userinfo.token.claim'),('0ef20646-a5a9-4c30-81e7-5bcc5836abe2','true','access.token.claim'),('0ef20646-a5a9-4c30-81e7-5bcc5836abe2','nickname','claim.name'),('0ef20646-a5a9-4c30-81e7-5bcc5836abe2','true','id.token.claim'),('0ef20646-a5a9-4c30-81e7-5bcc5836abe2','true','introspection.token.claim'),('0ef20646-a5a9-4c30-81e7-5bcc5836abe2','String','jsonType.label'),('0ef20646-a5a9-4c30-81e7-5bcc5836abe2','nickname','user.attribute'),('0ef20646-a5a9-4c30-81e7-5bcc5836abe2','true','userinfo.token.claim'),('1730c42a-1d72-411a-984a-9a819d4060c2','true','access.token.claim'),('1730c42a-1d72-411a-984a-9a819d4060c2','updated_at','claim.name'),('1730c42a-1d72-411a-984a-9a819d4060c2','true','id.token.claim'),('1730c42a-1d72-411a-984a-9a819d4060c2','true','introspection.token.claim'),('1730c42a-1d72-411a-984a-9a819d4060c2','long','jsonType.label'),('1730c42a-1d72-411a-984a-9a819d4060c2','updatedAt','user.attribute'),('1730c42a-1d72-411a-984a-9a819d4060c2','true','userinfo.token.claim'),('17387ac5-5633-4984-a987-f193e600f9b0','true','access.token.claim'),('17387ac5-5633-4984-a987-f193e600f9b0','upn','claim.name'),('17387ac5-5633-4984-a987-f193e600f9b0','true','id.token.claim'),('17387ac5-5633-4984-a987-f193e600f9b0','true','introspection.token.claim'),('17387ac5-5633-4984-a987-f193e600f9b0','String','jsonType.label'),('17387ac5-5633-4984-a987-f193e600f9b0','username','user.attribute'),('17387ac5-5633-4984-a987-f193e600f9b0','true','userinfo.token.claim'),('185e5efd-282f-4a83-b52c-40e9c72a0134','true','access.token.claim'),('185e5efd-282f-4a83-b52c-40e9c72a0134','gender','claim.name'),('185e5efd-282f-4a83-b52c-40e9c72a0134','true','id.token.claim'),('185e5efd-282f-4a83-b52c-40e9c72a0134','true','introspection.token.claim'),('185e5efd-282f-4a83-b52c-40e9c72a0134','String','jsonType.label'),('185e5efd-282f-4a83-b52c-40e9c72a0134','gender','user.attribute'),('185e5efd-282f-4a83-b52c-40e9c72a0134','true','userinfo.token.claim'),('1aab8d5c-7edb-45c7-b4e1-87b9364bfe87','true','access.token.claim'),('1aab8d5c-7edb-45c7-b4e1-87b9364bfe87','clientHost','claim.name'),('1aab8d5c-7edb-45c7-b4e1-87b9364bfe87','true','id.token.claim'),('1aab8d5c-7edb-45c7-b4e1-87b9364bfe87','true','introspection.token.claim'),('1aab8d5c-7edb-45c7-b4e1-87b9364bfe87','String','jsonType.label'),('1aab8d5c-7edb-45c7-b4e1-87b9364bfe87','clientHost','user.session.note'),('1fbc6ca0-99ae-44d2-83fd-f650feef91d4','true','access.token.claim'),('1fbc6ca0-99ae-44d2-83fd-f650feef91d4','preferred_username','claim.name'),('1fbc6ca0-99ae-44d2-83fd-f650feef91d4','true','id.token.claim'),('1fbc6ca0-99ae-44d2-83fd-f650feef91d4','true','introspection.token.claim'),('1fbc6ca0-99ae-44d2-83fd-f650feef91d4','String','jsonType.label'),('1fbc6ca0-99ae-44d2-83fd-f650feef91d4','username','user.attribute'),('1fbc6ca0-99ae-44d2-83fd-f650feef91d4','true','userinfo.token.claim'),('242421cc-38ad-40d2-b6f9-9ed40e42390d','true','access.token.claim'),('242421cc-38ad-40d2-b6f9-9ed40e42390d','profile','claim.name'),('242421cc-38ad-40d2-b6f9-9ed40e42390d','true','id.token.claim'),('242421cc-38ad-40d2-b6f9-9ed40e42390d','true','introspection.token.claim'),('242421cc-38ad-40d2-b6f9-9ed40e42390d','String','jsonType.label'),('242421cc-38ad-40d2-b6f9-9ed40e42390d','profile','user.attribute'),('242421cc-38ad-40d2-b6f9-9ed40e42390d','true','userinfo.token.claim'),('25206647-5675-4ffe-b195-893e3c62bd55','true','access.token.claim'),('25206647-5675-4ffe-b195-893e3c62bd55','true','introspection.token.claim'),('277f7ed0-4ffb-4efc-8c00-2a4354362da2','true','access.token.claim'),('277f7ed0-4ffb-4efc-8c00-2a4354362da2','resource_access.${client_id}.roles','claim.name'),('277f7ed0-4ffb-4efc-8c00-2a4354362da2','true','introspection.token.claim'),('277f7ed0-4ffb-4efc-8c00-2a4354362da2','String','jsonType.label'),('277f7ed0-4ffb-4efc-8c00-2a4354362da2','true','multivalued'),('277f7ed0-4ffb-4efc-8c00-2a4354362da2','foo','user.attribute'),('280c6f3a-cbe8-44d4-a630-40166862a4d3','true','access.token.claim'),('280c6f3a-cbe8-44d4-a630-40166862a4d3','clientAddress','claim.name'),('280c6f3a-cbe8-44d4-a630-40166862a4d3','true','id.token.claim'),('280c6f3a-cbe8-44d4-a630-40166862a4d3','true','introspection.token.claim'),('280c6f3a-cbe8-44d4-a630-40166862a4d3','String','jsonType.label'),('280c6f3a-cbe8-44d4-a630-40166862a4d3','clientAddress','user.session.note'),('2a6e41cd-993b-424b-8858-564276c6ef8c','true','access.token.claim'),('2a6e41cd-993b-424b-8858-564276c6ef8c','locale','claim.name'),('2a6e41cd-993b-424b-8858-564276c6ef8c','true','id.token.claim'),('2a6e41cd-993b-424b-8858-564276c6ef8c','true','introspection.token.claim'),('2a6e41cd-993b-424b-8858-564276c6ef8c','String','jsonType.label'),('2a6e41cd-993b-424b-8858-564276c6ef8c','locale','user.attribute'),('2a6e41cd-993b-424b-8858-564276c6ef8c','true','userinfo.token.claim'),('388800ab-5b9e-435e-ad59-d44d200ac5e1','true','access.token.claim'),('388800ab-5b9e-435e-ad59-d44d200ac5e1','realm_access.roles','claim.name'),('388800ab-5b9e-435e-ad59-d44d200ac5e1','true','introspection.token.claim'),('388800ab-5b9e-435e-ad59-d44d200ac5e1','String','jsonType.label'),('388800ab-5b9e-435e-ad59-d44d200ac5e1','true','multivalued'),('388800ab-5b9e-435e-ad59-d44d200ac5e1','foo','user.attribute'),('38a846a4-579c-436c-8800-02a1a6272c07','true','access.token.claim'),('38a846a4-579c-436c-8800-02a1a6272c07','email_verified','claim.name'),('38a846a4-579c-436c-8800-02a1a6272c07','true','id.token.claim'),('38a846a4-579c-436c-8800-02a1a6272c07','true','introspection.token.claim'),('38a846a4-579c-436c-8800-02a1a6272c07','boolean','jsonType.label'),('38a846a4-579c-436c-8800-02a1a6272c07','emailVerified','user.attribute'),('38a846a4-579c-436c-8800-02a1a6272c07','true','userinfo.token.claim'),('38cd59f8-9311-4e55-a636-cc85ea77eb3a','true','access.token.claim'),('38cd59f8-9311-4e55-a636-cc85ea77eb3a','true','introspection.token.claim'),('39808fda-1527-4cbd-8220-d698f76d9556','true','access.token.claim'),('39808fda-1527-4cbd-8220-d698f76d9556','profile','claim.name'),('39808fda-1527-4cbd-8220-d698f76d9556','true','id.token.claim'),('39808fda-1527-4cbd-8220-d698f76d9556','true','introspection.token.claim'),('39808fda-1527-4cbd-8220-d698f76d9556','String','jsonType.label'),('39808fda-1527-4cbd-8220-d698f76d9556','profile','user.attribute'),('39808fda-1527-4cbd-8220-d698f76d9556','true','userinfo.token.claim'),('3eb3fcc9-f7a9-4929-92d0-e2f69209771d','true','access.token.claim'),('3eb3fcc9-f7a9-4929-92d0-e2f69209771d','true','introspection.token.claim'),('43814bec-d0ae-426d-b3ef-95464ea1be67','true','access.token.claim'),('43814bec-d0ae-426d-b3ef-95464ea1be67','email','claim.name'),('43814bec-d0ae-426d-b3ef-95464ea1be67','true','id.token.claim'),('43814bec-d0ae-426d-b3ef-95464ea1be67','true','introspection.token.claim'),('43814bec-d0ae-426d-b3ef-95464ea1be67','String','jsonType.label'),('43814bec-d0ae-426d-b3ef-95464ea1be67','email','user.attribute'),('43814bec-d0ae-426d-b3ef-95464ea1be67','true','userinfo.token.claim'),('4fbfb957-8785-4380-8149-5be64cba1c01','true','access.token.claim'),('4fbfb957-8785-4380-8149-5be64cba1c01','locale','claim.name'),('4fbfb957-8785-4380-8149-5be64cba1c01','true','id.token.claim'),('4fbfb957-8785-4380-8149-5be64cba1c01','true','introspection.token.claim'),('4fbfb957-8785-4380-8149-5be64cba1c01','String','jsonType.label'),('4fbfb957-8785-4380-8149-5be64cba1c01','locale','user.attribute'),('4fbfb957-8785-4380-8149-5be64cba1c01','true','userinfo.token.claim'),('5566cade-a4a4-4a22-b524-258986841372','true','access.token.claim'),('5566cade-a4a4-4a22-b524-258986841372','clientAddress','claim.name'),('5566cade-a4a4-4a22-b524-258986841372','true','id.token.claim'),('5566cade-a4a4-4a22-b524-258986841372','true','introspection.token.claim'),('5566cade-a4a4-4a22-b524-258986841372','String','jsonType.label'),('5566cade-a4a4-4a22-b524-258986841372','clientAddress','user.session.note'),('5ee47111-68f0-4fe7-9686-247d1f92b5c7','true','access.token.claim'),('5ee47111-68f0-4fe7-9686-247d1f92b5c7','picture','claim.name'),('5ee47111-68f0-4fe7-9686-247d1f92b5c7','true','id.token.claim'),('5ee47111-68f0-4fe7-9686-247d1f92b5c7','true','introspection.token.claim'),('5ee47111-68f0-4fe7-9686-247d1f92b5c7','String','jsonType.label'),('5ee47111-68f0-4fe7-9686-247d1f92b5c7','picture','user.attribute'),('5ee47111-68f0-4fe7-9686-247d1f92b5c7','true','userinfo.token.claim'),('5f495607-3a9b-4850-acac-eefbe989f02e','true','access.token.claim'),('5f495607-3a9b-4850-acac-eefbe989f02e','zoneinfo','claim.name'),('5f495607-3a9b-4850-acac-eefbe989f02e','true','id.token.claim'),('5f495607-3a9b-4850-acac-eefbe989f02e','true','introspection.token.claim'),('5f495607-3a9b-4850-acac-eefbe989f02e','String','jsonType.label'),('5f495607-3a9b-4850-acac-eefbe989f02e','zoneinfo','user.attribute'),('5f495607-3a9b-4850-acac-eefbe989f02e','true','userinfo.token.claim'),('66a447c6-cc42-4b19-853f-6b66681d6082','true','access.token.claim'),('66a447c6-cc42-4b19-853f-6b66681d6082','middle_name','claim.name'),('66a447c6-cc42-4b19-853f-6b66681d6082','true','id.token.claim'),('66a447c6-cc42-4b19-853f-6b66681d6082','true','introspection.token.claim'),('66a447c6-cc42-4b19-853f-6b66681d6082','String','jsonType.label'),('66a447c6-cc42-4b19-853f-6b66681d6082','middleName','user.attribute'),('66a447c6-cc42-4b19-853f-6b66681d6082','true','userinfo.token.claim'),('6f2e805e-d911-440e-92d0-490bdef3cf2b','true','access.token.claim'),('6f2e805e-d911-440e-92d0-490bdef3cf2b','website','claim.name'),('6f2e805e-d911-440e-92d0-490bdef3cf2b','true','id.token.claim'),('6f2e805e-d911-440e-92d0-490bdef3cf2b','true','introspection.token.claim'),('6f2e805e-d911-440e-92d0-490bdef3cf2b','String','jsonType.label'),('6f2e805e-d911-440e-92d0-490bdef3cf2b','website','user.attribute'),('6f2e805e-d911-440e-92d0-490bdef3cf2b','true','userinfo.token.claim'),('71010a33-ce8f-41aa-a090-f042de837139','true','access.token.claim'),('71010a33-ce8f-41aa-a090-f042de837139','true','introspection.token.claim'),('71e65459-0242-489b-a3b5-07e3faa7304c','true','access.token.claim'),('71e65459-0242-489b-a3b5-07e3faa7304c','phone_number_verified','claim.name'),('71e65459-0242-489b-a3b5-07e3faa7304c','true','id.token.claim'),('71e65459-0242-489b-a3b5-07e3faa7304c','true','introspection.token.claim'),('71e65459-0242-489b-a3b5-07e3faa7304c','boolean','jsonType.label'),('71e65459-0242-489b-a3b5-07e3faa7304c','phoneNumberVerified','user.attribute'),('71e65459-0242-489b-a3b5-07e3faa7304c','true','userinfo.token.claim'),('73c2a0a6-3325-43f5-86dd-1da5b6b84510','true','access.token.claim'),('73c2a0a6-3325-43f5-86dd-1da5b6b84510','true','id.token.claim'),('73c2a0a6-3325-43f5-86dd-1da5b6b84510','true','introspection.token.claim'),('75f8d3f9-1d00-48c2-9425-c82b140ef2af','true','access.token.claim'),('75f8d3f9-1d00-48c2-9425-c82b140ef2af','phone_number_verified','claim.name'),('75f8d3f9-1d00-48c2-9425-c82b140ef2af','true','id.token.claim'),('75f8d3f9-1d00-48c2-9425-c82b140ef2af','true','introspection.token.claim'),('75f8d3f9-1d00-48c2-9425-c82b140ef2af','boolean','jsonType.label'),('75f8d3f9-1d00-48c2-9425-c82b140ef2af','phoneNumberVerified','user.attribute'),('75f8d3f9-1d00-48c2-9425-c82b140ef2af','true','userinfo.token.claim'),('76be9620-7989-440c-8c68-16da57b40d19','true','access.token.claim'),('76be9620-7989-440c-8c68-16da57b40d19','updated_at','claim.name'),('76be9620-7989-440c-8c68-16da57b40d19','true','id.token.claim'),('76be9620-7989-440c-8c68-16da57b40d19','true','introspection.token.claim'),('76be9620-7989-440c-8c68-16da57b40d19','long','jsonType.label'),('76be9620-7989-440c-8c68-16da57b40d19','updatedAt','user.attribute'),('76be9620-7989-440c-8c68-16da57b40d19','true','userinfo.token.claim'),('7a99cbf4-af03-4ceb-80cd-b540d2876037','true','access.token.claim'),('7a99cbf4-af03-4ceb-80cd-b540d2876037','true','id.token.claim'),('7a99cbf4-af03-4ceb-80cd-b540d2876037','true','introspection.token.claim'),('7a99cbf4-af03-4ceb-80cd-b540d2876037','country','user.attribute.country'),('7a99cbf4-af03-4ceb-80cd-b540d2876037','formatted','user.attribute.formatted'),('7a99cbf4-af03-4ceb-80cd-b540d2876037','locality','user.attribute.locality'),('7a99cbf4-af03-4ceb-80cd-b540d2876037','postal_code','user.attribute.postal_code'),('7a99cbf4-af03-4ceb-80cd-b540d2876037','region','user.attribute.region'),('7a99cbf4-af03-4ceb-80cd-b540d2876037','street','user.attribute.street'),('7a99cbf4-af03-4ceb-80cd-b540d2876037','true','userinfo.token.claim'),('7e9aa70f-9359-4fd8-a047-b52720f65c50','true','access.token.claim'),('7e9aa70f-9359-4fd8-a047-b52720f65c50','family_name','claim.name'),('7e9aa70f-9359-4fd8-a047-b52720f65c50','true','id.token.claim'),('7e9aa70f-9359-4fd8-a047-b52720f65c50','true','introspection.token.claim'),('7e9aa70f-9359-4fd8-a047-b52720f65c50','String','jsonType.label'),('7e9aa70f-9359-4fd8-a047-b52720f65c50','lastName','user.attribute'),('7e9aa70f-9359-4fd8-a047-b52720f65c50','true','userinfo.token.claim'),('81ded2d3-54f8-431c-a864-cd8ae9485e88','true','access.token.claim'),('81ded2d3-54f8-431c-a864-cd8ae9485e88','true','id.token.claim'),('81ded2d3-54f8-431c-a864-cd8ae9485e88','true','introspection.token.claim'),('81ded2d3-54f8-431c-a864-cd8ae9485e88','true','userinfo.token.claim'),('83d5038c-9b1a-432c-be40-1cfd8ec8c8e8','true','access.token.claim'),('83d5038c-9b1a-432c-be40-1cfd8ec8c8e8','groups','claim.name'),('83d5038c-9b1a-432c-be40-1cfd8ec8c8e8','true','id.token.claim'),('83d5038c-9b1a-432c-be40-1cfd8ec8c8e8','true','introspection.token.claim'),('83d5038c-9b1a-432c-be40-1cfd8ec8c8e8','String','jsonType.label'),('83d5038c-9b1a-432c-be40-1cfd8ec8c8e8','true','multivalued'),('83d5038c-9b1a-432c-be40-1cfd8ec8c8e8','foo','user.attribute'),('881ccc55-c62b-4e4a-b95d-7aeb8f73fe80','true','access.token.claim'),('881ccc55-c62b-4e4a-b95d-7aeb8f73fe80','true','id.token.claim'),('881ccc55-c62b-4e4a-b95d-7aeb8f73fe80','true','introspection.token.claim'),('881ccc55-c62b-4e4a-b95d-7aeb8f73fe80','true','userinfo.token.claim'),('88556cd3-317c-4ea8-ab61-09618a18da2b','true','access.token.claim'),('88556cd3-317c-4ea8-ab61-09618a18da2b','organization','claim.name'),('88556cd3-317c-4ea8-ab61-09618a18da2b','true','id.token.claim'),('88556cd3-317c-4ea8-ab61-09618a18da2b','true','introspection.token.claim'),('88556cd3-317c-4ea8-ab61-09618a18da2b','String','jsonType.label'),('88556cd3-317c-4ea8-ab61-09618a18da2b','true','multivalued'),('90bc110b-9196-46a1-92f7-8f80d3ff465c','true','access.token.claim'),('90bc110b-9196-46a1-92f7-8f80d3ff465c','email','claim.name'),('90bc110b-9196-46a1-92f7-8f80d3ff465c','true','id.token.claim'),('90bc110b-9196-46a1-92f7-8f80d3ff465c','true','introspection.token.claim'),('90bc110b-9196-46a1-92f7-8f80d3ff465c','String','jsonType.label'),('90bc110b-9196-46a1-92f7-8f80d3ff465c','email','user.attribute'),('90bc110b-9196-46a1-92f7-8f80d3ff465c','true','userinfo.token.claim'),('91182dc0-9ad6-43a6-b15b-b062821d6d97','Role','attribute.name'),('91182dc0-9ad6-43a6-b15b-b062821d6d97','Basic','attribute.nameformat'),('91182dc0-9ad6-43a6-b15b-b062821d6d97','false','single'),('9b115eee-e55c-489d-b203-85013572f223','true','access.token.claim'),('9b115eee-e55c-489d-b203-85013572f223','email_verified','claim.name'),('9b115eee-e55c-489d-b203-85013572f223','true','id.token.claim'),('9b115eee-e55c-489d-b203-85013572f223','true','introspection.token.claim'),('9b115eee-e55c-489d-b203-85013572f223','boolean','jsonType.label'),('9b115eee-e55c-489d-b203-85013572f223','emailVerified','user.attribute'),('9b115eee-e55c-489d-b203-85013572f223','true','userinfo.token.claim'),('9f6277a2-296e-43ac-98df-c152cec85e9f','true','access.token.claim'),('9f6277a2-296e-43ac-98df-c152cec85e9f','true','introspection.token.claim'),('a15f4aa2-6d01-4636-8f96-fc33234b1956','true','access.token.claim'),('a15f4aa2-6d01-4636-8f96-fc33234b1956','phone_number','claim.name'),('a15f4aa2-6d01-4636-8f96-fc33234b1956','true','id.token.claim'),('a15f4aa2-6d01-4636-8f96-fc33234b1956','true','introspection.token.claim'),('a15f4aa2-6d01-4636-8f96-fc33234b1956','String','jsonType.label'),('a15f4aa2-6d01-4636-8f96-fc33234b1956','phoneNumber','user.attribute'),('a15f4aa2-6d01-4636-8f96-fc33234b1956','true','userinfo.token.claim'),('a1f1a727-50e7-4f79-b239-87ac4fed7ad2','true','access.token.claim'),('a1f1a727-50e7-4f79-b239-87ac4fed7ad2','client_id','claim.name'),('a1f1a727-50e7-4f79-b239-87ac4fed7ad2','true','id.token.claim'),('a1f1a727-50e7-4f79-b239-87ac4fed7ad2','true','introspection.token.claim'),('a1f1a727-50e7-4f79-b239-87ac4fed7ad2','String','jsonType.label'),('a1f1a727-50e7-4f79-b239-87ac4fed7ad2','client_id','user.session.note'),('a29fbd5d-e05d-4cee-bff5-81b1359085b3','true','access.token.claim'),('a29fbd5d-e05d-4cee-bff5-81b1359085b3','birthdate','claim.name'),('a29fbd5d-e05d-4cee-bff5-81b1359085b3','true','id.token.claim'),('a29fbd5d-e05d-4cee-bff5-81b1359085b3','true','introspection.token.claim'),('a29fbd5d-e05d-4cee-bff5-81b1359085b3','String','jsonType.label'),('a29fbd5d-e05d-4cee-bff5-81b1359085b3','birthdate','user.attribute'),('a29fbd5d-e05d-4cee-bff5-81b1359085b3','true','userinfo.token.claim'),('a454fc4c-9349-4c5c-b7bd-77d9fefcfa8c','true','access.token.claim'),('a454fc4c-9349-4c5c-b7bd-77d9fefcfa8c','true','introspection.token.claim'),('a7fa25ad-c21f-433c-bc89-3b7489b2f745','true','access.token.claim'),('a7fa25ad-c21f-433c-bc89-3b7489b2f745','preferred_username','claim.name'),('a7fa25ad-c21f-433c-bc89-3b7489b2f745','true','id.token.claim'),('a7fa25ad-c21f-433c-bc89-3b7489b2f745','true','introspection.token.claim'),('a7fa25ad-c21f-433c-bc89-3b7489b2f745','String','jsonType.label'),('a7fa25ad-c21f-433c-bc89-3b7489b2f745','username','user.attribute'),('a7fa25ad-c21f-433c-bc89-3b7489b2f745','true','userinfo.token.claim'),('aa852eb2-6df6-49ff-943a-f8b86e67e992','true','access.token.claim'),('aa852eb2-6df6-49ff-943a-f8b86e67e992','organization','claim.name'),('aa852eb2-6df6-49ff-943a-f8b86e67e992','true','id.token.claim'),('aa852eb2-6df6-49ff-943a-f8b86e67e992','true','introspection.token.claim'),('aa852eb2-6df6-49ff-943a-f8b86e67e992','String','jsonType.label'),('aa852eb2-6df6-49ff-943a-f8b86e67e992','true','multivalued'),('aba0fb9f-bf9d-4413-8d0f-1d85bbf0383f','true','access.token.claim'),('aba0fb9f-bf9d-4413-8d0f-1d85bbf0383f','resource_access.${client_id}.roles','claim.name'),('aba0fb9f-bf9d-4413-8d0f-1d85bbf0383f','true','introspection.token.claim'),('aba0fb9f-bf9d-4413-8d0f-1d85bbf0383f','String','jsonType.label'),('aba0fb9f-bf9d-4413-8d0f-1d85bbf0383f','true','multivalued'),('aba0fb9f-bf9d-4413-8d0f-1d85bbf0383f','foo','user.attribute'),('b671eb28-0bd1-4f17-9b73-9e20fe1cc45a','true','access.token.claim'),('b671eb28-0bd1-4f17-9b73-9e20fe1cc45a','locale','claim.name'),('b671eb28-0bd1-4f17-9b73-9e20fe1cc45a','true','id.token.claim'),('b671eb28-0bd1-4f17-9b73-9e20fe1cc45a','true','introspection.token.claim'),('b671eb28-0bd1-4f17-9b73-9e20fe1cc45a','String','jsonType.label'),('b671eb28-0bd1-4f17-9b73-9e20fe1cc45a','locale','user.attribute'),('b671eb28-0bd1-4f17-9b73-9e20fe1cc45a','true','userinfo.token.claim'),('bb0728a9-4774-4d72-86cc-4edea3d8f789','true','access.token.claim'),('bb0728a9-4774-4d72-86cc-4edea3d8f789','realm_access.roles','claim.name'),('bb0728a9-4774-4d72-86cc-4edea3d8f789','true','introspection.token.claim'),('bb0728a9-4774-4d72-86cc-4edea3d8f789','String','jsonType.label'),('bb0728a9-4774-4d72-86cc-4edea3d8f789','true','multivalued'),('bb0728a9-4774-4d72-86cc-4edea3d8f789','foo','user.attribute'),('bbfd59ca-4897-454c-b6c7-4f33dd935efe','true','access.token.claim'),('bbfd59ca-4897-454c-b6c7-4f33dd935efe','auth_time','claim.name'),('bbfd59ca-4897-454c-b6c7-4f33dd935efe','true','id.token.claim'),('bbfd59ca-4897-454c-b6c7-4f33dd935efe','true','introspection.token.claim'),('bbfd59ca-4897-454c-b6c7-4f33dd935efe','long','jsonType.label'),('bbfd59ca-4897-454c-b6c7-4f33dd935efe','AUTH_TIME','user.session.note'),('c0ddb572-44f1-4264-bf8b-a4ab876aeaf4','true','access.token.claim'),('c0ddb572-44f1-4264-bf8b-a4ab876aeaf4','picture','claim.name'),('c0ddb572-44f1-4264-bf8b-a4ab876aeaf4','true','id.token.claim'),('c0ddb572-44f1-4264-bf8b-a4ab876aeaf4','true','introspection.token.claim'),('c0ddb572-44f1-4264-bf8b-a4ab876aeaf4','String','jsonType.label'),('c0ddb572-44f1-4264-bf8b-a4ab876aeaf4','picture','user.attribute'),('c0ddb572-44f1-4264-bf8b-a4ab876aeaf4','true','userinfo.token.claim'),('c862c840-8335-41cf-aef3-318a023b586c','true','access.token.claim'),('c862c840-8335-41cf-aef3-318a023b586c','birthdate','claim.name'),('c862c840-8335-41cf-aef3-318a023b586c','true','id.token.claim'),('c862c840-8335-41cf-aef3-318a023b586c','true','introspection.token.claim'),('c862c840-8335-41cf-aef3-318a023b586c','String','jsonType.label'),('c862c840-8335-41cf-aef3-318a023b586c','birthdate','user.attribute'),('c862c840-8335-41cf-aef3-318a023b586c','true','userinfo.token.claim'),('ce536f16-d5dc-4688-8913-2b0f34d7671f','true','access.token.claim'),('ce536f16-d5dc-4688-8913-2b0f34d7671f','true','id.token.claim'),('ce536f16-d5dc-4688-8913-2b0f34d7671f','true','introspection.token.claim'),('ce536f16-d5dc-4688-8913-2b0f34d7671f','country','user.attribute.country'),('ce536f16-d5dc-4688-8913-2b0f34d7671f','formatted','user.attribute.formatted'),('ce536f16-d5dc-4688-8913-2b0f34d7671f','locality','user.attribute.locality'),('ce536f16-d5dc-4688-8913-2b0f34d7671f','postal_code','user.attribute.postal_code'),('ce536f16-d5dc-4688-8913-2b0f34d7671f','region','user.attribute.region'),('ce536f16-d5dc-4688-8913-2b0f34d7671f','street','user.attribute.street'),('ce536f16-d5dc-4688-8913-2b0f34d7671f','true','userinfo.token.claim'),('d2f52858-16d3-425b-b7ce-7270a0ee3e78','true','access.token.claim'),('d2f52858-16d3-425b-b7ce-7270a0ee3e78','locale','claim.name'),('d2f52858-16d3-425b-b7ce-7270a0ee3e78','true','id.token.claim'),('d2f52858-16d3-425b-b7ce-7270a0ee3e78','true','introspection.token.claim'),('d2f52858-16d3-425b-b7ce-7270a0ee3e78','String','jsonType.label'),('d2f52858-16d3-425b-b7ce-7270a0ee3e78','locale','user.attribute'),('d2f52858-16d3-425b-b7ce-7270a0ee3e78','true','userinfo.token.claim'),('da0f3ba9-f24d-43c6-8e21-42f77f2f846e','true','access.token.claim'),('da0f3ba9-f24d-43c6-8e21-42f77f2f846e','given_name','claim.name'),('da0f3ba9-f24d-43c6-8e21-42f77f2f846e','true','id.token.claim'),('da0f3ba9-f24d-43c6-8e21-42f77f2f846e','true','introspection.token.claim'),('da0f3ba9-f24d-43c6-8e21-42f77f2f846e','String','jsonType.label'),('da0f3ba9-f24d-43c6-8e21-42f77f2f846e','firstName','user.attribute'),('da0f3ba9-f24d-43c6-8e21-42f77f2f846e','true','userinfo.token.claim'),('dde2d6fa-647f-434f-847a-0c7c915c9e12','true','access.token.claim'),('dde2d6fa-647f-434f-847a-0c7c915c9e12','clientHost','claim.name'),('dde2d6fa-647f-434f-847a-0c7c915c9e12','true','id.token.claim'),('dde2d6fa-647f-434f-847a-0c7c915c9e12','true','introspection.token.claim'),('dde2d6fa-647f-434f-847a-0c7c915c9e12','String','jsonType.label'),('dde2d6fa-647f-434f-847a-0c7c915c9e12','clientHost','user.session.note'),('e073d946-8b03-4f79-b4ed-de5fcad5a543','true','access.token.claim'),('e073d946-8b03-4f79-b4ed-de5fcad5a543','given_name','claim.name'),('e073d946-8b03-4f79-b4ed-de5fcad5a543','true','id.token.claim'),('e073d946-8b03-4f79-b4ed-de5fcad5a543','true','introspection.token.claim'),('e073d946-8b03-4f79-b4ed-de5fcad5a543','String','jsonType.label'),('e073d946-8b03-4f79-b4ed-de5fcad5a543','firstName','user.attribute'),('e073d946-8b03-4f79-b4ed-de5fcad5a543','true','userinfo.token.claim'),('e2cffa4c-6764-49b1-bf78-9d3a7ef656e2','true','access.token.claim'),('e2cffa4c-6764-49b1-bf78-9d3a7ef656e2','true','id.token.claim'),('e2cffa4c-6764-49b1-bf78-9d3a7ef656e2','true','introspection.token.claim'),('e47981cc-6126-4f12-983d-e5af93ad8261','true','access.token.claim'),('e47981cc-6126-4f12-983d-e5af93ad8261','website','claim.name'),('e47981cc-6126-4f12-983d-e5af93ad8261','true','id.token.claim'),('e47981cc-6126-4f12-983d-e5af93ad8261','true','introspection.token.claim'),('e47981cc-6126-4f12-983d-e5af93ad8261','String','jsonType.label'),('e47981cc-6126-4f12-983d-e5af93ad8261','website','user.attribute'),('e47981cc-6126-4f12-983d-e5af93ad8261','true','userinfo.token.claim'),('ec976618-f1b8-4285-8f69-313d975ff126','true','access.token.claim'),('ec976618-f1b8-4285-8f69-313d975ff126','gender','claim.name'),('ec976618-f1b8-4285-8f69-313d975ff126','true','id.token.claim'),('ec976618-f1b8-4285-8f69-313d975ff126','true','introspection.token.claim'),('ec976618-f1b8-4285-8f69-313d975ff126','String','jsonType.label'),('ec976618-f1b8-4285-8f69-313d975ff126','gender','user.attribute'),('ec976618-f1b8-4285-8f69-313d975ff126','true','userinfo.token.claim'),('ed02d7d5-3f23-4c6c-aa11-744bb3a03002','true','access.token.claim'),('ed02d7d5-3f23-4c6c-aa11-744bb3a03002','phone_number','claim.name'),('ed02d7d5-3f23-4c6c-aa11-744bb3a03002','true','id.token.claim'),('ed02d7d5-3f23-4c6c-aa11-744bb3a03002','true','introspection.token.claim'),('ed02d7d5-3f23-4c6c-aa11-744bb3a03002','String','jsonType.label'),('ed02d7d5-3f23-4c6c-aa11-744bb3a03002','phoneNumber','user.attribute'),('ed02d7d5-3f23-4c6c-aa11-744bb3a03002','true','userinfo.token.claim'),('ed4234fb-cf69-414b-8fde-bfb47e0d5575','true','access.token.claim'),('ed4234fb-cf69-414b-8fde-bfb47e0d5575','groups','claim.name'),('ed4234fb-cf69-414b-8fde-bfb47e0d5575','true','id.token.claim'),('ed4234fb-cf69-414b-8fde-bfb47e0d5575','true','introspection.token.claim'),('ed4234fb-cf69-414b-8fde-bfb47e0d5575','String','jsonType.label'),('ed4234fb-cf69-414b-8fde-bfb47e0d5575','true','multivalued'),('ed4234fb-cf69-414b-8fde-bfb47e0d5575','foo','user.attribute'),('f01049bb-0c22-4527-a7f5-1c93bd4235ff','true','access.token.claim'),('f01049bb-0c22-4527-a7f5-1c93bd4235ff','zoneinfo','claim.name'),('f01049bb-0c22-4527-a7f5-1c93bd4235ff','true','id.token.claim'),('f01049bb-0c22-4527-a7f5-1c93bd4235ff','true','introspection.token.claim'),('f01049bb-0c22-4527-a7f5-1c93bd4235ff','String','jsonType.label'),('f01049bb-0c22-4527-a7f5-1c93bd4235ff','zoneinfo','user.attribute'),('f01049bb-0c22-4527-a7f5-1c93bd4235ff','true','userinfo.token.claim'),('f11bc60a-0f09-4588-8415-7e44aa362618','Role','attribute.name'),('f11bc60a-0f09-4588-8415-7e44aa362618','Basic','attribute.nameformat'),('f11bc60a-0f09-4588-8415-7e44aa362618','false','single'),('f76a4e44-b836-4d58-80f8-1012c94278f4','true','access.token.claim'),('f76a4e44-b836-4d58-80f8-1012c94278f4','auth_time','claim.name'),('f76a4e44-b836-4d58-80f8-1012c94278f4','true','id.token.claim'),('f76a4e44-b836-4d58-80f8-1012c94278f4','true','introspection.token.claim'),('f76a4e44-b836-4d58-80f8-1012c94278f4','long','jsonType.label'),('f76a4e44-b836-4d58-80f8-1012c94278f4','AUTH_TIME','user.session.note'),('f883dadf-2b2d-4feb-ae39-746e6007a128','true','access.token.claim'),('f883dadf-2b2d-4feb-ae39-746e6007a128','upn','claim.name'),('f883dadf-2b2d-4feb-ae39-746e6007a128','true','id.token.claim'),('f883dadf-2b2d-4feb-ae39-746e6007a128','true','introspection.token.claim'),('f883dadf-2b2d-4feb-ae39-746e6007a128','String','jsonType.label'),('f883dadf-2b2d-4feb-ae39-746e6007a128','username','user.attribute'),('f883dadf-2b2d-4feb-ae39-746e6007a128','true','userinfo.token.claim');
/*!40000 ALTER TABLE `PROTOCOL_MAPPER_CONFIG` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `REALM`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `REALM` (
  `ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ACCESS_CODE_LIFESPAN` int DEFAULT NULL,
  `USER_ACTION_LIFESPAN` int DEFAULT NULL,
  `ACCESS_TOKEN_LIFESPAN` int DEFAULT NULL,
  `ACCOUNT_THEME` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `ADMIN_THEME` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `EMAIL_THEME` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `ENABLED` tinyint NOT NULL DEFAULT '0',
  `EVENTS_ENABLED` tinyint NOT NULL DEFAULT '0',
  `EVENTS_EXPIRATION` bigint DEFAULT NULL,
  `LOGIN_THEME` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `NAME` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `NOT_BEFORE` int DEFAULT NULL,
  `PASSWORD_POLICY` text COLLATE utf8mb4_unicode_ci,
  `REGISTRATION_ALLOWED` tinyint NOT NULL DEFAULT '0',
  `REMEMBER_ME` tinyint NOT NULL DEFAULT '0',
  `RESET_PASSWORD_ALLOWED` tinyint NOT NULL DEFAULT '0',
  `SOCIAL` tinyint NOT NULL DEFAULT '0',
  `SSL_REQUIRED` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `SSO_IDLE_TIMEOUT` int DEFAULT NULL,
  `SSO_MAX_LIFESPAN` int DEFAULT NULL,
  `UPDATE_PROFILE_ON_SOC_LOGIN` tinyint NOT NULL DEFAULT '0',
  `VERIFY_EMAIL` tinyint NOT NULL DEFAULT '0',
  `MASTER_ADMIN_CLIENT` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `LOGIN_LIFESPAN` int DEFAULT NULL,
  `INTERNATIONALIZATION_ENABLED` tinyint NOT NULL DEFAULT '0',
  `DEFAULT_LOCALE` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `REG_EMAIL_AS_USERNAME` tinyint NOT NULL DEFAULT '0',
  `ADMIN_EVENTS_ENABLED` tinyint NOT NULL DEFAULT '0',
  `ADMIN_EVENTS_DETAILS_ENABLED` tinyint NOT NULL DEFAULT '0',
  `EDIT_USERNAME_ALLOWED` tinyint NOT NULL DEFAULT '0',
  `OTP_POLICY_COUNTER` int DEFAULT '0',
  `OTP_POLICY_WINDOW` int DEFAULT '1',
  `OTP_POLICY_PERIOD` int DEFAULT '30',
  `OTP_POLICY_DIGITS` int DEFAULT '6',
  `OTP_POLICY_ALG` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT 'HmacSHA1',
  `OTP_POLICY_TYPE` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT 'totp',
  `BROWSER_FLOW` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `REGISTRATION_FLOW` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `DIRECT_GRANT_FLOW` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `RESET_CREDENTIALS_FLOW` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `CLIENT_AUTH_FLOW` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `OFFLINE_SESSION_IDLE_TIMEOUT` int DEFAULT '0',
  `REVOKE_REFRESH_TOKEN` tinyint NOT NULL DEFAULT '0',
  `ACCESS_TOKEN_LIFE_IMPLICIT` int DEFAULT '0',
  `LOGIN_WITH_EMAIL_ALLOWED` tinyint NOT NULL DEFAULT '1',
  `DUPLICATE_EMAILS_ALLOWED` tinyint NOT NULL DEFAULT '0',
  `DOCKER_AUTH_FLOW` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `REFRESH_TOKEN_MAX_REUSE` int DEFAULT '0',
  `ALLOW_USER_MANAGED_ACCESS` tinyint NOT NULL DEFAULT '0',
  `SSO_MAX_LIFESPAN_REMEMBER_ME` int NOT NULL,
  `SSO_IDLE_TIMEOUT_REMEMBER_ME` int NOT NULL,
  `DEFAULT_ROLE` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `UK_ORVSDMLA56612EAEFIQ6WL5OI` (`NAME`),
  KEY `IDX_REALM_MASTER_ADM_CLI` (`MASTER_ADMIN_CLIENT`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `REALM` WRITE;
/*!40000 ALTER TABLE `REALM` DISABLE KEYS */;
INSERT INTO `REALM` VALUES ('4c2886d8-e4c4-4b40-9064-445c7ab90a1c',60,300,54000,'','','motogo',1,0,0,'motogo','motogo',0,'upperCase(1) and length(8) and specialChars(1) and maxLength(64)',0,0,0,0,'NONE',604800,691200,0,1,'232e71d1-663b-4fd9-981f-0f08b56cd567',1800,0,NULL,0,0,0,0,0,1,30,6,'HmacSHA1','totp','e23ba648-6f5b-447f-8d45-33da2f4f7e9e','49326f20-3ceb-4c4b-8171-8e7d27d44093','f662b0fb-ef63-4db6-934e-e4d37ca37296','e9deb89c-1118-49eb-83ba-9a9c1dfb9a2a','a0edc4f0-693e-42be-8fd9-f6f6da0515a0',2592000,1,900,1,0,'906eb7da-9867-4708-9777-30c21db80523',20,0,0,0,'5125458c-9476-41c9-8fb6-0543f7083773'),('5d56fd15-57bc-496f-8f3f-bf0d7a5c5581',60,300,60,NULL,NULL,NULL,1,0,0,NULL,'master',0,NULL,0,0,0,0,'NONE',1800,36000,0,0,'6b5c5dea-9656-4dfb-8a26-f699429eecd3',1800,0,NULL,0,0,0,0,0,1,30,6,'HmacSHA1','totp','676c5a06-ae2d-44a2-89a8-9017809fa2b7','26b753d6-6520-4727-b63d-dac95814059f','a3c9dcd5-2130-44c7-bab3-55f37e7c2da6','d0f21bd5-8d4d-4dee-b58e-77f8cd44ad8c','de9091e2-e7a5-42f3-aefb-56c5360be436',2592000,0,86400,1,0,'216f221b-9be0-40f5-9c68-f5bee839e1bc',0,0,0,0,'0c50da1d-eeb7-4074-b31e-f638245df45e');
/*!40000 ALTER TABLE `REALM` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `REALM_ATTRIBUTE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `REALM_ATTRIBUTE` (
  `NAME` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `REALM_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `VALUE` longtext CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  PRIMARY KEY (`NAME`,`REALM_ID`),
  KEY `IDX_REALM_ATTR_REALM` (`REALM_ID`),
  CONSTRAINT `FK_8SHXD6L3E9ATQUKACXGPFFPTW` FOREIGN KEY (`REALM_ID`) REFERENCES `REALM` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `REALM_ATTRIBUTE` WRITE;
/*!40000 ALTER TABLE `REALM_ATTRIBUTE` DISABLE KEYS */;
INSERT INTO `REALM_ATTRIBUTE` VALUES ('_browser_header.contentSecurityPolicy','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','frame-src \'self\'; frame-ancestors \'self\'; object-src \'none\';'),('_browser_header.contentSecurityPolicy','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','frame-src \'self\'; frame-ancestors \'self\'; object-src \'none\';'),('_browser_header.contentSecurityPolicyReportOnly','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',''),('_browser_header.contentSecurityPolicyReportOnly','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581',''),('_browser_header.referrerPolicy','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','no-referrer'),('_browser_header.referrerPolicy','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','no-referrer'),('_browser_header.strictTransportSecurity','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','max-age=31536000; includeSubDomains'),('_browser_header.strictTransportSecurity','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','max-age=31536000; includeSubDomains'),('_browser_header.xContentTypeOptions','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','nosniff'),('_browser_header.xContentTypeOptions','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','nosniff'),('_browser_header.xFrameOptions','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','SAMEORIGIN'),('_browser_header.xFrameOptions','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','SAMEORIGIN'),('_browser_header.xRobotsTag','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','none'),('_browser_header.xRobotsTag','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','none'),('actionTokenGeneratedByAdminLifespan','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','43200'),('actionTokenGeneratedByAdminLifespan','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','43200'),('actionTokenGeneratedByUserLifespan','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','300'),('actionTokenGeneratedByUserLifespan','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','300'),('actionTokenGeneratedByUserLifespan.execute-actions','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',''),('actionTokenGeneratedByUserLifespan.execute-actions','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581',''),('actionTokenGeneratedByUserLifespan.idp-verify-account-via-email','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',''),('actionTokenGeneratedByUserLifespan.idp-verify-account-via-email','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581',''),('actionTokenGeneratedByUserLifespan.reset-credentials','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',''),('actionTokenGeneratedByUserLifespan.reset-credentials','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581',''),('actionTokenGeneratedByUserLifespan.verify-email','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',''),('actionTokenGeneratedByUserLifespan.verify-email','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581',''),('adminPermissionsEnabled','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','false'),('adminPermissionsEnabled','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','false'),('bruteForceProtected','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','true'),('bruteForceProtected','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','false'),('bruteForceStrategy','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','MULTIPLE'),('bruteForceStrategy','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','MULTIPLE'),('cibaAuthRequestedUserHint','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','login_hint'),('cibaAuthRequestedUserHint','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','login_hint'),('cibaBackchannelTokenDeliveryMode','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','poll'),('cibaBackchannelTokenDeliveryMode','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','poll'),('cibaExpiresIn','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','120'),('cibaExpiresIn','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','120'),('cibaInterval','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','5'),('cibaInterval','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','5'),('client-policies.policies','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','{\"policies\":[]}'),('client-policies.policies','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','{\"policies\":[]}'),('client-policies.profiles','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','{\"profiles\":[]}'),('client-policies.profiles','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','{\"profiles\":[]}'),('clientOfflineSessionIdleTimeout','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','0'),('clientOfflineSessionIdleTimeout','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','0'),('clientOfflineSessionMaxLifespan','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','0'),('clientOfflineSessionMaxLifespan','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','0'),('clientSessionIdleTimeout','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','0'),('clientSessionIdleTimeout','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','0'),('clientSessionMaxLifespan','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','0'),('clientSessionMaxLifespan','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','0'),('darkMode','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','true'),('defaultSignatureAlgorithm','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','RS256'),('defaultSignatureAlgorithm','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','RS256'),('displayName','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','Keycloak'),('displayNameHtml','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','<div class=\"kc-logo-text\"><span>Keycloak</span></div>'),('failureFactor','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','5'),('failureFactor','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','30'),('firstBrokerLoginFlowId','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','f3d39771-85ea-427d-b233-fa26308a3952'),('firstBrokerLoginFlowId','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','403e0820-f1af-4eff-8cbc-6a5423453d76'),('maxDeltaTimeSeconds','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','3600'),('maxDeltaTimeSeconds','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','43200'),('maxFailureWaitSeconds','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','900'),('maxFailureWaitSeconds','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','900'),('maxTemporaryLockouts','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','0'),('maxTemporaryLockouts','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','0'),('minimumQuickLoginWaitSeconds','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','3600'),('minimumQuickLoginWaitSeconds','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','60'),('oauth2DeviceCodeLifespan','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','600'),('oauth2DeviceCodeLifespan','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','600'),('oauth2DevicePollingInterval','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','7'),('oauth2DevicePollingInterval','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','5'),('offlineSessionMaxLifespan','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','5184000'),('offlineSessionMaxLifespan','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','5184000'),('offlineSessionMaxLifespanEnabled','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','false'),('offlineSessionMaxLifespanEnabled','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','false'),('organizationsEnabled','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','false'),('organizationsEnabled','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','false'),('parRequestUriLifespan','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','60'),('parRequestUriLifespan','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','60'),('permanentLockout','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','false'),('permanentLockout','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','false'),('quickLoginCheckMilliSeconds','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','1000'),('quickLoginCheckMilliSeconds','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','1000'),('realmReusableOtpCode','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','false'),('realmReusableOtpCode','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','false'),('shortVerificationUri','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',''),('shortVerificationUri','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581',''),('verifiableCredentialsEnabled','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','false'),('verifiableCredentialsEnabled','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','false'),('waitIncrementSeconds','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','60'),('waitIncrementSeconds','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','60'),('webAuthnPolicyAttestationConveyancePreference','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','not specified'),('webAuthnPolicyAttestationConveyancePreference','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','not specified'),('webAuthnPolicyAttestationConveyancePreferencePasswordless','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','not specified'),('webAuthnPolicyAttestationConveyancePreferencePasswordless','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','not specified'),('webAuthnPolicyAuthenticatorAttachment','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','not specified'),('webAuthnPolicyAuthenticatorAttachment','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','not specified'),('webAuthnPolicyAuthenticatorAttachmentPasswordless','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','not specified'),('webAuthnPolicyAuthenticatorAttachmentPasswordless','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','not specified'),('webAuthnPolicyAvoidSameAuthenticatorRegister','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','false'),('webAuthnPolicyAvoidSameAuthenticatorRegister','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','false'),('webAuthnPolicyAvoidSameAuthenticatorRegisterPasswordless','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','false'),('webAuthnPolicyAvoidSameAuthenticatorRegisterPasswordless','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','false'),('webAuthnPolicyCreateTimeout','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','0'),('webAuthnPolicyCreateTimeout','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','0'),('webAuthnPolicyCreateTimeoutPasswordless','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','0'),('webAuthnPolicyCreateTimeoutPasswordless','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','0'),('webAuthnPolicyRequireResidentKey','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','not specified'),('webAuthnPolicyRequireResidentKey','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','not specified'),('webAuthnPolicyRequireResidentKeyPasswordless','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','Yes'),('webAuthnPolicyRequireResidentKeyPasswordless','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','Yes'),('webAuthnPolicyRpEntityName','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','keycloak'),('webAuthnPolicyRpEntityName','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','keycloak'),('webAuthnPolicyRpEntityNamePasswordless','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','keycloak'),('webAuthnPolicyRpEntityNamePasswordless','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','keycloak'),('webAuthnPolicyRpId','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',''),('webAuthnPolicyRpId','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581',''),('webAuthnPolicyRpIdPasswordless','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',''),('webAuthnPolicyRpIdPasswordless','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581',''),('webAuthnPolicySignatureAlgorithms','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','ES256,RS256'),('webAuthnPolicySignatureAlgorithms','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','ES256,RS256'),('webAuthnPolicySignatureAlgorithmsPasswordless','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','ES256,RS256'),('webAuthnPolicySignatureAlgorithmsPasswordless','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','ES256,RS256'),('webAuthnPolicyUserVerificationRequirement','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','not specified'),('webAuthnPolicyUserVerificationRequirement','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','not specified'),('webAuthnPolicyUserVerificationRequirementPasswordless','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','required'),('webAuthnPolicyUserVerificationRequirementPasswordless','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','required');
/*!40000 ALTER TABLE `REALM_ATTRIBUTE` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `REALM_DEFAULT_GROUPS`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `REALM_DEFAULT_GROUPS` (
  `REALM_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `GROUP_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`REALM_ID`,`GROUP_ID`),
  UNIQUE KEY `CON_GROUP_ID_DEF_GROUPS` (`GROUP_ID`),
  KEY `IDX_REALM_DEF_GRP_REALM` (`REALM_ID`),
  CONSTRAINT `FK_DEF_GROUPS_REALM` FOREIGN KEY (`REALM_ID`) REFERENCES `REALM` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `REALM_DEFAULT_GROUPS` WRITE;
/*!40000 ALTER TABLE `REALM_DEFAULT_GROUPS` DISABLE KEYS */;
/*!40000 ALTER TABLE `REALM_DEFAULT_GROUPS` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `REALM_ENABLED_EVENT_TYPES`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `REALM_ENABLED_EVENT_TYPES` (
  `REALM_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `VALUE` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`REALM_ID`,`VALUE`),
  KEY `IDX_REALM_EVT_TYPES_REALM` (`REALM_ID`),
  CONSTRAINT `FK_H846O4H0W8EPX5NWEDRF5Y69J` FOREIGN KEY (`REALM_ID`) REFERENCES `REALM` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `REALM_ENABLED_EVENT_TYPES` WRITE;
/*!40000 ALTER TABLE `REALM_ENABLED_EVENT_TYPES` DISABLE KEYS */;
/*!40000 ALTER TABLE `REALM_ENABLED_EVENT_TYPES` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `REALM_EVENTS_LISTENERS`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `REALM_EVENTS_LISTENERS` (
  `REALM_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `VALUE` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`REALM_ID`,`VALUE`),
  KEY `IDX_REALM_EVT_LIST_REALM` (`REALM_ID`),
  CONSTRAINT `FK_H846O4H0W8EPX5NXEV9F5Y69J` FOREIGN KEY (`REALM_ID`) REFERENCES `REALM` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `REALM_EVENTS_LISTENERS` WRITE;
/*!40000 ALTER TABLE `REALM_EVENTS_LISTENERS` DISABLE KEYS */;
INSERT INTO `REALM_EVENTS_LISTENERS` VALUES ('4c2886d8-e4c4-4b40-9064-445c7ab90a1c','jboss-logging'),('5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','jboss-logging');
/*!40000 ALTER TABLE `REALM_EVENTS_LISTENERS` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `REALM_LOCALIZATIONS`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `REALM_LOCALIZATIONS` (
  `REALM_ID` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `LOCALE` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `TEXTS` longtext CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL,
  PRIMARY KEY (`REALM_ID`,`LOCALE`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `REALM_LOCALIZATIONS` WRITE;
/*!40000 ALTER TABLE `REALM_LOCALIZATIONS` DISABLE KEYS */;
/*!40000 ALTER TABLE `REALM_LOCALIZATIONS` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `REALM_REQUIRED_CREDENTIAL`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `REALM_REQUIRED_CREDENTIAL` (
  `TYPE` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `FORM_LABEL` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `INPUT` tinyint NOT NULL DEFAULT '0',
  `SECRET` tinyint NOT NULL DEFAULT '0',
  `REALM_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`REALM_ID`,`TYPE`),
  CONSTRAINT `FK_5HG65LYBEVAVKQFKI3KPONH9V` FOREIGN KEY (`REALM_ID`) REFERENCES `REALM` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `REALM_REQUIRED_CREDENTIAL` WRITE;
/*!40000 ALTER TABLE `REALM_REQUIRED_CREDENTIAL` DISABLE KEYS */;
INSERT INTO `REALM_REQUIRED_CREDENTIAL` VALUES ('password','password',1,1,'4c2886d8-e4c4-4b40-9064-445c7ab90a1c'),('password','password',1,1,'5d56fd15-57bc-496f-8f3f-bf0d7a5c5581');
/*!40000 ALTER TABLE `REALM_REQUIRED_CREDENTIAL` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `REALM_SMTP_CONFIG`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `REALM_SMTP_CONFIG` (
  `REALM_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `VALUE` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `NAME` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`REALM_ID`,`NAME`),
  CONSTRAINT `FK_70EJ8XDXGXD0B9HH6180IRR0O` FOREIGN KEY (`REALM_ID`) REFERENCES `REALM` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `REALM_SMTP_CONFIG` WRITE;
/*!40000 ALTER TABLE `REALM_SMTP_CONFIG` DISABLE KEYS */;
INSERT INTO `REALM_SMTP_CONFIG` VALUES ('4c2886d8-e4c4-4b40-9064-445c7ab90a1c','','allowutf8'),('4c2886d8-e4c4-4b40-9064-445c7ab90a1c','true','auth'),('4c2886d8-e4c4-4b40-9064-445c7ab90a1c','basic','authType'),('4c2886d8-e4c4-4b40-9064-445c7ab90a1c','false','debug'),('4c2886d8-e4c4-4b40-9064-445c7ab90a1c','','envelopeFrom'),('4c2886d8-e4c4-4b40-9064-445c7ab90a1c','noreply_motogo@rbsuport.com','from'),('4c2886d8-e4c4-4b40-9064-445c7ab90a1c','MotoGo','fromDisplayName'),('4c2886d8-e4c4-4b40-9064-445c7ab90a1c','smtp.resend.com','host'),('4c2886d8-e4c4-4b40-9064-445c7ab90a1c','re_eefMc9Ca_BVDyzRzNKo5YCRHp1nRsAtCd','password'),('4c2886d8-e4c4-4b40-9064-445c7ab90a1c','587','port'),('4c2886d8-e4c4-4b40-9064-445c7ab90a1c','','replyTo'),('4c2886d8-e4c4-4b40-9064-445c7ab90a1c','','replyToDisplayName'),('4c2886d8-e4c4-4b40-9064-445c7ab90a1c','false','ssl'),('4c2886d8-e4c4-4b40-9064-445c7ab90a1c','true','starttls'),('4c2886d8-e4c4-4b40-9064-445c7ab90a1c','resend','user');
/*!40000 ALTER TABLE `REALM_SMTP_CONFIG` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `REALM_SUPPORTED_LOCALES`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `REALM_SUPPORTED_LOCALES` (
  `REALM_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `VALUE` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`REALM_ID`,`VALUE`),
  KEY `IDX_REALM_SUPP_LOCAL_REALM` (`REALM_ID`),
  CONSTRAINT `FK_SUPPORTED_LOCALES_REALM` FOREIGN KEY (`REALM_ID`) REFERENCES `REALM` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `REALM_SUPPORTED_LOCALES` WRITE;
/*!40000 ALTER TABLE `REALM_SUPPORTED_LOCALES` DISABLE KEYS */;
/*!40000 ALTER TABLE `REALM_SUPPORTED_LOCALES` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `REDIRECT_URIS`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `REDIRECT_URIS` (
  `CLIENT_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `VALUE` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`CLIENT_ID`,`VALUE`),
  KEY `IDX_REDIR_URI_CLIENT` (`CLIENT_ID`),
  CONSTRAINT `FK_1BURS8PB4OUJ97H5WUPPAHV9F` FOREIGN KEY (`CLIENT_ID`) REFERENCES `CLIENT` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `REDIRECT_URIS` WRITE;
/*!40000 ALTER TABLE `REDIRECT_URIS` DISABLE KEYS */;
INSERT INTO `REDIRECT_URIS` VALUES ('13fd5765-ac20-4558-ae2c-cf8e2df4c22c','/realms/master/account/*'),('286f0039-9d0d-4f17-8b87-9cbeddaa0a89','/realms/master/account/*'),('35ba46a4-b1e1-4dd0-8cff-fc98df4c2659','http://localhost:8085/*'),('3c6f327b-8018-44e9-8540-511170baeb2f','/admin/motogo/console/*'),('4c9063cc-607c-4f14-b55d-eb3d700e742c','/admin/master/console/*'),('8b0a29f2-46dd-4f3c-8446-c0f63ad56ebc','/realms/motogo/account/*'),('d2eb4296-2046-420d-beb2-25cf9e05b715','/realms/motogo/account/*');
/*!40000 ALTER TABLE `REDIRECT_URIS` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `REQUIRED_ACTION_CONFIG`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `REQUIRED_ACTION_CONFIG` (
  `REQUIRED_ACTION_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `VALUE` longtext COLLATE utf8mb4_unicode_ci,
  `NAME` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`REQUIRED_ACTION_ID`,`NAME`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `REQUIRED_ACTION_CONFIG` WRITE;
/*!40000 ALTER TABLE `REQUIRED_ACTION_CONFIG` DISABLE KEYS */;
/*!40000 ALTER TABLE `REQUIRED_ACTION_CONFIG` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `REQUIRED_ACTION_PROVIDER`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `REQUIRED_ACTION_PROVIDER` (
  `ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ALIAS` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `NAME` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `REALM_ID` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `ENABLED` tinyint NOT NULL DEFAULT '0',
  `DEFAULT_ACTION` tinyint NOT NULL DEFAULT '0',
  `PROVIDER_ID` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `PRIORITY` int DEFAULT NULL,
  PRIMARY KEY (`ID`),
  KEY `IDX_REQ_ACT_PROV_REALM` (`REALM_ID`),
  CONSTRAINT `FK_REQ_ACT_REALM` FOREIGN KEY (`REALM_ID`) REFERENCES `REALM` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `REQUIRED_ACTION_PROVIDER` WRITE;
/*!40000 ALTER TABLE `REQUIRED_ACTION_PROVIDER` DISABLE KEYS */;
INSERT INTO `REQUIRED_ACTION_PROVIDER` VALUES ('1d53645f-efda-451d-8600-b9f9deeee8a0','TERMS_AND_CONDITIONS','Terms and Conditions','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',0,0,'TERMS_AND_CONDITIONS',20),('22a659c3-3f0a-4f9d-ae85-516bbdacb739','VERIFY_EMAIL','Verify Email','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',1,0,'VERIFY_EMAIL',50),('249d6b65-da48-4c1f-8f4c-724f0db7df97','delete_credential','Delete Credential','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',1,0,'delete_credential',110),('2640d858-05a9-4009-b01a-363070a4f732','CONFIGURE_RECOVERY_AUTHN_CODES','Recovery Authentication Codes','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',1,0,'CONFIGURE_RECOVERY_AUTHN_CODES',130),('273d655f-b498-4983-bd95-3229f3ea9c4f','VERIFY_PROFILE','Verify Profile','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',1,0,'VERIFY_PROFILE',100),('2aafd5b4-81d8-48ec-8e62-1cfd160d5eab','delete_account','Delete Account','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',0,0,'delete_account',60),('2fa26728-973f-4f80-a7ef-66fd4cb161a4','UPDATE_PROFILE','Update Profile','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',0,0,'UPDATE_PROFILE',40),('4491c30b-a297-40af-8d62-ae5812e8027d','UPDATE_EMAIL','Update Email','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581',0,0,'UPDATE_EMAIL',80),('579179ac-774f-4f26-b246-426db0aed711','UPDATE_PASSWORD','Update Password','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581',0,0,'UPDATE_PASSWORD',40),('57f36c23-bf95-41da-9d55-bd50f564c190','CONFIGURE_TOTP','Configure OTP','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',0,0,'CONFIGURE_TOTP',10),('58d1a518-8ab9-4e7a-aad5-9ab7446cde11','UPDATE_EMAIL','Update Email','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',0,0,'UPDATE_EMAIL',70),('65da426c-7357-42da-80bf-7ffcfb9d1324','idp_link','Linking Identity Provider','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581',1,0,'idp_link',130),('6f385838-2867-4b67-b758-9534a45a89f1','webauthn-register-passwordless','Webauthn Register Passwordless','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581',1,0,'webauthn-register-passwordless',100),('71f6a18a-bcd0-4a6f-9e5f-f5abb0307eb5','delete_account','Delete Account','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581',0,0,'delete_account',70),('7e373aba-3333-41e4-a92d-0db1a72cdfe0','webauthn-register','Webauthn Register','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',1,0,'webauthn-register',80),('7ed33aca-749f-4191-b3cd-6e88bc918d2f','VERIFY_PROFILE','Verify Profile','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581',1,0,'VERIFY_PROFILE',110),('7ed4b8df-a421-4573-b39b-b27e270adae6','CONFIGURE_RECOVERY_AUTHN_CODES','Recovery Authentication Codes','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581',1,0,'CONFIGURE_RECOVERY_AUTHN_CODES',1000),('861979a0-3a2f-4627-bbfe-4a13379163ac','idp_link','Linking Identity Provider','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',1,0,'idp_link',120),('8c0c6e9d-806c-4c88-a58f-0371dcbd598e','webauthn-register','Webauthn Register','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581',1,0,'webauthn-register',90),('95f237f5-baba-4500-882d-57480b8df7eb','TERMS_AND_CONDITIONS','Terms and Conditions','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581',0,0,'TERMS_AND_CONDITIONS',30),('9b8142ee-8c5b-4113-95c0-6086da41b435','webauthn-register-passwordless','Webauthn Register Passwordless','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',1,0,'webauthn-register-passwordless',90),('b723a48a-bcc4-4f90-9a28-7423b1f7904c','VERIFY_EMAIL','Verify Email','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581',1,0,'VERIFY_EMAIL',60),('bd81f76c-9141-4615-91d7-788ea3c19c77','UPDATE_PROFILE','Update Profile','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581',0,0,'UPDATE_PROFILE',50),('bf864e01-cc04-4b9d-9300-0f33b02c5232','update_user_locale','Update User Locale','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581',1,0,'update_user_locale',10),('c3dfd33b-af4b-4f4f-8473-38b7f8262665','update_user_locale','Update User Locale','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',1,0,'update_user_locale',1000),('c5ac80b9-313c-45a7-8b10-5c9337cb1d1b','CONFIGURE_TOTP','Configure OTP','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581',0,0,'CONFIGURE_TOTP',20),('c6d71c20-f946-4f07-828c-421e288d9b2c','delete_credential','Delete Credential','5d56fd15-57bc-496f-8f3f-bf0d7a5c5581',1,0,'delete_credential',120),('d2080e43-8b0c-4d6e-a932-0be90239109f','UPDATE_PASSWORD','Update Password','4c2886d8-e4c4-4b40-9064-445c7ab90a1c',0,0,'UPDATE_PASSWORD',30);
/*!40000 ALTER TABLE `REQUIRED_ACTION_PROVIDER` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `RESOURCE_ATTRIBUTE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `RESOURCE_ATTRIBUTE` (
  `ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'sybase-needs-something-here',
  `NAME` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `VALUE` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `RESOURCE_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`ID`),
  KEY `FK_5HRM2VLF9QL5FU022KQEPOVBR` (`RESOURCE_ID`),
  CONSTRAINT `FK_5HRM2VLF9QL5FU022KQEPOVBR` FOREIGN KEY (`RESOURCE_ID`) REFERENCES `RESOURCE_SERVER_RESOURCE` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `RESOURCE_ATTRIBUTE` WRITE;
/*!40000 ALTER TABLE `RESOURCE_ATTRIBUTE` DISABLE KEYS */;
/*!40000 ALTER TABLE `RESOURCE_ATTRIBUTE` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `RESOURCE_POLICY`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `RESOURCE_POLICY` (
  `RESOURCE_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `POLICY_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`RESOURCE_ID`,`POLICY_ID`),
  KEY `IDX_RES_POLICY_POLICY` (`POLICY_ID`),
  CONSTRAINT `FK_FRSRPOS53XCX4WNKOG82SSRFY` FOREIGN KEY (`RESOURCE_ID`) REFERENCES `RESOURCE_SERVER_RESOURCE` (`ID`),
  CONSTRAINT `FK_FRSRPP213XCX4WNKOG82SSRFY` FOREIGN KEY (`POLICY_ID`) REFERENCES `RESOURCE_SERVER_POLICY` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `RESOURCE_POLICY` WRITE;
/*!40000 ALTER TABLE `RESOURCE_POLICY` DISABLE KEYS */;
/*!40000 ALTER TABLE `RESOURCE_POLICY` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `RESOURCE_SCOPE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `RESOURCE_SCOPE` (
  `RESOURCE_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `SCOPE_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`RESOURCE_ID`,`SCOPE_ID`),
  KEY `IDX_RES_SCOPE_SCOPE` (`SCOPE_ID`),
  CONSTRAINT `FK_FRSRPOS13XCX4WNKOG82SSRFY` FOREIGN KEY (`RESOURCE_ID`) REFERENCES `RESOURCE_SERVER_RESOURCE` (`ID`),
  CONSTRAINT `FK_FRSRPS213XCX4WNKOG82SSRFY` FOREIGN KEY (`SCOPE_ID`) REFERENCES `RESOURCE_SERVER_SCOPE` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `RESOURCE_SCOPE` WRITE;
/*!40000 ALTER TABLE `RESOURCE_SCOPE` DISABLE KEYS */;
/*!40000 ALTER TABLE `RESOURCE_SCOPE` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `RESOURCE_SERVER`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `RESOURCE_SERVER` (
  `ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ALLOW_RS_REMOTE_MGMT` tinyint NOT NULL DEFAULT '0',
  `POLICY_ENFORCE_MODE` tinyint DEFAULT NULL,
  `DECISION_STRATEGY` tinyint NOT NULL DEFAULT '1',
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `RESOURCE_SERVER` WRITE;
/*!40000 ALTER TABLE `RESOURCE_SERVER` DISABLE KEYS */;
/*!40000 ALTER TABLE `RESOURCE_SERVER` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `RESOURCE_SERVER_PERM_TICKET`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `RESOURCE_SERVER_PERM_TICKET` (
  `ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `OWNER` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `REQUESTER` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `CREATED_TIMESTAMP` bigint NOT NULL,
  `GRANTED_TIMESTAMP` bigint DEFAULT NULL,
  `RESOURCE_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `SCOPE_ID` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `RESOURCE_SERVER_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `POLICY_ID` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `UK_FRSR6T700S9V50BU18WS5PMT` (`OWNER`,`REQUESTER`,`RESOURCE_SERVER_ID`,`RESOURCE_ID`,`SCOPE_ID`),
  KEY `FK_FRSRHO213XCX4WNKOG82SSPMT` (`RESOURCE_SERVER_ID`),
  KEY `FK_FRSRHO213XCX4WNKOG83SSPMT` (`RESOURCE_ID`),
  KEY `FK_FRSRHO213XCX4WNKOG84SSPMT` (`SCOPE_ID`),
  KEY `FK_FRSRPO2128CX4WNKOG82SSRFY` (`POLICY_ID`),
  KEY `IDX_PERM_TICKET_REQUESTER` (`REQUESTER`),
  KEY `IDX_PERM_TICKET_OWNER` (`OWNER`),
  CONSTRAINT `FK_FRSRHO213XCX4WNKOG82SSPMT` FOREIGN KEY (`RESOURCE_SERVER_ID`) REFERENCES `RESOURCE_SERVER` (`ID`),
  CONSTRAINT `FK_FRSRHO213XCX4WNKOG83SSPMT` FOREIGN KEY (`RESOURCE_ID`) REFERENCES `RESOURCE_SERVER_RESOURCE` (`ID`),
  CONSTRAINT `FK_FRSRHO213XCX4WNKOG84SSPMT` FOREIGN KEY (`SCOPE_ID`) REFERENCES `RESOURCE_SERVER_SCOPE` (`ID`),
  CONSTRAINT `FK_FRSRPO2128CX4WNKOG82SSRFY` FOREIGN KEY (`POLICY_ID`) REFERENCES `RESOURCE_SERVER_POLICY` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `RESOURCE_SERVER_PERM_TICKET` WRITE;
/*!40000 ALTER TABLE `RESOURCE_SERVER_PERM_TICKET` DISABLE KEYS */;
/*!40000 ALTER TABLE `RESOURCE_SERVER_PERM_TICKET` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `RESOURCE_SERVER_POLICY`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `RESOURCE_SERVER_POLICY` (
  `ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `NAME` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `DESCRIPTION` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `TYPE` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `DECISION_STRATEGY` tinyint DEFAULT NULL,
  `LOGIC` tinyint DEFAULT NULL,
  `RESOURCE_SERVER_ID` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `OWNER` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `UK_FRSRPT700S9V50BU18WS5HA6` (`NAME`,`RESOURCE_SERVER_ID`),
  KEY `IDX_RES_SERV_POL_RES_SERV` (`RESOURCE_SERVER_ID`),
  CONSTRAINT `FK_FRSRPO213XCX4WNKOG82SSRFY` FOREIGN KEY (`RESOURCE_SERVER_ID`) REFERENCES `RESOURCE_SERVER` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `RESOURCE_SERVER_POLICY` WRITE;
/*!40000 ALTER TABLE `RESOURCE_SERVER_POLICY` DISABLE KEYS */;
/*!40000 ALTER TABLE `RESOURCE_SERVER_POLICY` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `RESOURCE_SERVER_RESOURCE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `RESOURCE_SERVER_RESOURCE` (
  `ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `NAME` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `TYPE` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `ICON_URI` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `OWNER` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `RESOURCE_SERVER_ID` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `OWNER_MANAGED_ACCESS` tinyint NOT NULL DEFAULT '0',
  `DISPLAY_NAME` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `UK_FRSR6T700S9V50BU18WS5HA6` (`NAME`,`OWNER`,`RESOURCE_SERVER_ID`),
  KEY `IDX_RES_SRV_RES_RES_SRV` (`RESOURCE_SERVER_ID`),
  CONSTRAINT `FK_FRSRHO213XCX4WNKOG82SSRFY` FOREIGN KEY (`RESOURCE_SERVER_ID`) REFERENCES `RESOURCE_SERVER` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `RESOURCE_SERVER_RESOURCE` WRITE;
/*!40000 ALTER TABLE `RESOURCE_SERVER_RESOURCE` DISABLE KEYS */;
/*!40000 ALTER TABLE `RESOURCE_SERVER_RESOURCE` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `RESOURCE_SERVER_SCOPE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `RESOURCE_SERVER_SCOPE` (
  `ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `NAME` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ICON_URI` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `RESOURCE_SERVER_ID` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `DISPLAY_NAME` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `UK_FRSRST700S9V50BU18WS5HA6` (`NAME`,`RESOURCE_SERVER_ID`),
  KEY `IDX_RES_SRV_SCOPE_RES_SRV` (`RESOURCE_SERVER_ID`),
  CONSTRAINT `FK_FRSRSO213XCX4WNKOG82SSRFY` FOREIGN KEY (`RESOURCE_SERVER_ID`) REFERENCES `RESOURCE_SERVER` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `RESOURCE_SERVER_SCOPE` WRITE;
/*!40000 ALTER TABLE `RESOURCE_SERVER_SCOPE` DISABLE KEYS */;
/*!40000 ALTER TABLE `RESOURCE_SERVER_SCOPE` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `RESOURCE_URIS`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `RESOURCE_URIS` (
  `RESOURCE_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `VALUE` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`RESOURCE_ID`,`VALUE`),
  CONSTRAINT `FK_RESOURCE_SERVER_URIS` FOREIGN KEY (`RESOURCE_ID`) REFERENCES `RESOURCE_SERVER_RESOURCE` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `RESOURCE_URIS` WRITE;
/*!40000 ALTER TABLE `RESOURCE_URIS` DISABLE KEYS */;
/*!40000 ALTER TABLE `RESOURCE_URIS` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `REVOKED_TOKEN`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `REVOKED_TOKEN` (
  `ID` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `EXPIRE` bigint NOT NULL,
  PRIMARY KEY (`ID`),
  KEY `IDX_REV_TOKEN_ON_EXPIRE` (`EXPIRE`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `REVOKED_TOKEN` WRITE;
/*!40000 ALTER TABLE `REVOKED_TOKEN` DISABLE KEYS */;
/*!40000 ALTER TABLE `REVOKED_TOKEN` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `ROLE_ATTRIBUTE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ROLE_ATTRIBUTE` (
  `ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ROLE_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `NAME` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `VALUE` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  PRIMARY KEY (`ID`),
  KEY `IDX_ROLE_ATTRIBUTE` (`ROLE_ID`),
  CONSTRAINT `FK_ROLE_ATTRIBUTE_ID` FOREIGN KEY (`ROLE_ID`) REFERENCES `KEYCLOAK_ROLE` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `ROLE_ATTRIBUTE` WRITE;
/*!40000 ALTER TABLE `ROLE_ATTRIBUTE` DISABLE KEYS */;
/*!40000 ALTER TABLE `ROLE_ATTRIBUTE` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `SCOPE_MAPPING`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `SCOPE_MAPPING` (
  `CLIENT_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ROLE_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`CLIENT_ID`,`ROLE_ID`),
  KEY `IDX_SCOPE_MAPPING_ROLE` (`ROLE_ID`),
  CONSTRAINT `FK_OUSE064PLMLR732LXJCN1Q5F1` FOREIGN KEY (`CLIENT_ID`) REFERENCES `CLIENT` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `SCOPE_MAPPING` WRITE;
/*!40000 ALTER TABLE `SCOPE_MAPPING` DISABLE KEYS */;
INSERT INTO `SCOPE_MAPPING` VALUES ('d2eb4296-2046-420d-beb2-25cf9e05b715','56de3a81-f2d8-4b89-ae11-f933ea5efcd8'),('13fd5765-ac20-4558-ae2c-cf8e2df4c22c','a8ded98f-bc0d-490e-99d6-320855772b20'),('d2eb4296-2046-420d-beb2-25cf9e05b715','bc732c7b-1dda-4fcb-bbfd-d2c11138eebf'),('13fd5765-ac20-4558-ae2c-cf8e2df4c22c','f8772141-c42e-4f39-b42c-6652dd8a5b70');
/*!40000 ALTER TABLE `SCOPE_MAPPING` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `SCOPE_POLICY`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `SCOPE_POLICY` (
  `SCOPE_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `POLICY_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`SCOPE_ID`,`POLICY_ID`),
  KEY `IDX_SCOPE_POLICY_POLICY` (`POLICY_ID`),
  CONSTRAINT `FK_FRSRASP13XCX4WNKOG82SSRFY` FOREIGN KEY (`POLICY_ID`) REFERENCES `RESOURCE_SERVER_POLICY` (`ID`),
  CONSTRAINT `FK_FRSRPASS3XCX4WNKOG82SSRFY` FOREIGN KEY (`SCOPE_ID`) REFERENCES `RESOURCE_SERVER_SCOPE` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `SCOPE_POLICY` WRITE;
/*!40000 ALTER TABLE `SCOPE_POLICY` DISABLE KEYS */;
/*!40000 ALTER TABLE `SCOPE_POLICY` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `SERVER_CONFIG`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `SERVER_CONFIG` (
  `SERVER_CONFIG_KEY` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `VALUE` longtext COLLATE utf8mb4_unicode_ci NOT NULL,
  `VERSION` int DEFAULT '0',
  PRIMARY KEY (`SERVER_CONFIG_KEY`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `SERVER_CONFIG` WRITE;
/*!40000 ALTER TABLE `SERVER_CONFIG` DISABLE KEYS */;
INSERT INTO `SERVER_CONFIG` VALUES ('crt_jgroups','{\"prvKey\":\"MIIEowIBAAKCAQEAt9XOs/7YZ5wF9iWhEcXLmHEgc06m0IgZ0a2i6e1bh+hz5Nckz01v1fI0Dk5NfyIggeVUVEdGauRivN3uq53A+2XhWHxWJXtAxBmVSrgp43K7nuOL1tQf2pSXjmFgQNXcurqzH0npensVU0MDkktJvpyQt2HiEfvc8SLBKAe3yu4NI2JQD93sv3zDlGfOo3M7evC+16m38EBL4ab7lmz8d0gGdU0GhYGdqk/OudCVseD+b0bBDdiPgLa/bmpSSZtmDdjfEeuuJhNEH42p4kURSABTZ5xNpja5U0gZli3lIIfgX0vBa4KfO/+Lnt+/QDUA+vEDARcxknRd4r882mgUjQIDAQABAoIBAAKlhdnKh/QJ534RMZ3vYox5yHVQgw9KflCp9sf2XlR5EL4RKNsBKuSCHaS9FNmJNL2Z/xpNgmKEkjVU3pxC46Xa5CEhfR4CBb8s/Ry7xQKcwr25WeiIWkML05i9pySXz5syl+SR7FYSko6FzvXatp+TYnP5aUdIAHjKPsykO82klHJhrWlwZy/KFfIZKyJdOkkqZ2O95bpeCtpiXY+vyWtsXoUeLs7KdP7J4WGg1UaSBlN5rAfNRafLFq+opvlOW+Bmmu8Il3gDqI86SsCzrfJ0UAETwDPfB2CDU6p9DTFg1V5qHobc3W1wJSAaVIGHqnoRhTiTV8p+vz22FT1C2tECgYEA54p9uAbFCxpWv8tqDHVQMutTy/ybwra88/QTLtENvFZF2ujp+UuWTLMmPoLiGj8kPvi60jTuIibxozjklcsNoMK920d96JUS+1B+fAExG2uPzredsqvyzloLWKBbD9Cl8SzJYHPMvd3UsG86vQDx9wx77qhP0f/sC/y9DAWR5jUCgYEAy0E4ODlxtTr8B84dk6Njkb9Q01eEZtcVH6WdlSJYCArREyHVtXVExxAokUwsPGPeSQcHyyystgylXsDk7+w3cCkxtzzsWLv48EAbzSbUBY/anqEjS3Hb7lqW7CbgFEeZjs7Md8Ok+qetnr92LL+N7JBxM8uIu3YrDc5Qnlr73/kCgYAV+op2hoX7yNNFP3zgk946ByacLW1nrUsAaUHM5uFD6HiXEBbtqCQrQbI2qtvlm+rH43pwa7/TFBlJ4iOoqG6xvOqivwWPp1725iwrfz13Gd27xKg+hp30wo1jK02Zkbe+zf6LnlJoD3+aA+Tyff158wBspiN1jrKMbMBRPNIdpQKBgQCKYiEwxouQZFRry/1/I2MLPEJ/YbqXeiC7pzqe6v8cmqC2cfba0Li+/1Vfd4k2pJgXCF5aClScrgd5d8NnRmCIUMO0ihowX+qhSESlbKY2Ezc5gWRSXXbr1Wmj5uPxar416L51XBLbRYaD6r8+wDyUr9Mi/JRfbZjqkWl6J2ANoQKBgHkevTc0oRYlBHH0lS3ZWHyGsAJi7MqcMhgD8UCxFJwpGnEuSLuQMqK73mmQ2dvPTliR61lvzDQRmMyc5+T1sJ4k/rBKkYYFKuOPKYkvDQxhvypY4yaNcYa09vduH0wdPKMKtyt+btUT+ciX7Ln5vNl93g0twGO1Q+0T1hXYxOF2\",\"pubKey\":\"MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAt9XOs/7YZ5wF9iWhEcXLmHEgc06m0IgZ0a2i6e1bh+hz5Nckz01v1fI0Dk5NfyIggeVUVEdGauRivN3uq53A+2XhWHxWJXtAxBmVSrgp43K7nuOL1tQf2pSXjmFgQNXcurqzH0npensVU0MDkktJvpyQt2HiEfvc8SLBKAe3yu4NI2JQD93sv3zDlGfOo3M7evC+16m38EBL4ab7lmz8d0gGdU0GhYGdqk/OudCVseD+b0bBDdiPgLa/bmpSSZtmDdjfEeuuJhNEH42p4kURSABTZ5xNpja5U0gZli3lIIfgX0vBa4KfO/+Lnt+/QDUA+vEDARcxknRd4r882mgUjQIDAQAB\",\"crt\":\"MIICnTCCAYUCBgGcZzdvozANBgkqhkiG9w0BAQsFADASMRAwDgYDVQQDDAdqZ3JvdXBzMB4XDTI2MDIxNjE2MDgzMloXDTI2MDQxNzE2MTAxMVowEjEQMA4GA1UEAwwHamdyb3VwczCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBALfVzrP+2GecBfYloRHFy5hxIHNOptCIGdGtountW4foc+TXJM9Nb9XyNA5OTX8iIIHlVFRHRmrkYrzd7qudwPtl4Vh8ViV7QMQZlUq4KeNyu57ji9bUH9qUl45hYEDV3Lq6sx9J6Xp7FVNDA5JLSb6ckLdh4hH73PEiwSgHt8ruDSNiUA/d7L98w5RnzqNzO3rwvtept/BAS+Gm+5Zs/HdIBnVNBoWBnapPzrnQlbHg/m9GwQ3Yj4C2v25qUkmbZg3Y3xHrriYTRB+NqeJFEUgAU2ecTaY2uVNIGZYt5SCH4F9LwWuCnzv/i57fv0A1APrxAwEXMZJ0XeK/PNpoFI0CAwEAATANBgkqhkiG9w0BAQsFAAOCAQEAJ/VXaE7x4ztmXOwMZ/BkQF6JpSTOb9LIE4TPULFHAhhBmbmNQOOfti2TmfQZoPyT1p/4nRCFavAtPq8+5NLhWjVaOfiNIIHATOy60eNJp0PBkWn1dpV6kaMVf1zlu7QTSgRqoZpOdZQR/12IiO1x1SR9g/3u5PhhBjfdom7p8kWZ5fZAZLSlqV7KcVicHHOY17/YbyZtS+ZW7v2llA4Cjo9fAvU0kIEldZwYpWPKkmtZpdylpfsCVvl7R5956tHZ4tHQvs7QmylYu7NZEgFLrwl2LLSc77bVkwhCoG5pTKNANpeyEDzNBJolweD2uO3eBmUdZVpxawklxhfjOH63hA==\",\"alias\":\"e23870e1-98fc-4651-876d-f4bae261448a\",\"generatedMillis\":1771258212706}',2),('JGROUPS_ADDRESS_SEQUENCE','99',99);
/*!40000 ALTER TABLE `SERVER_CONFIG` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `USER_ATTRIBUTE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `USER_ATTRIBUTE` (
  `NAME` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `VALUE` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `USER_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'sybase-needs-something-here',
  `LONG_VALUE_HASH` binary(64) DEFAULT NULL,
  `LONG_VALUE_HASH_LOWER_CASE` binary(64) DEFAULT NULL,
  `LONG_VALUE` longtext CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  PRIMARY KEY (`ID`),
  KEY `IDX_USER_ATTRIBUTE` (`USER_ID`),
  KEY `IDX_USER_ATTRIBUTE_NAME` (`NAME`,`VALUE`),
  KEY `USER_ATTR_LONG_VALUES` (`LONG_VALUE_HASH`,`NAME`),
  KEY `USER_ATTR_LONG_VALUES_LOWER_CASE` (`LONG_VALUE_HASH_LOWER_CASE`,`NAME`),
  CONSTRAINT `FK_5HRM2VLF9QL5FU043KQEPOVBR` FOREIGN KEY (`USER_ID`) REFERENCES `USER_ENTITY` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `USER_ATTRIBUTE` WRITE;
/*!40000 ALTER TABLE `USER_ATTRIBUTE` DISABLE KEYS */;
INSERT INTO `USER_ATTRIBUTE` VALUES ('is_temporary_admin','true','19b6d072-21e8-40f9-9d22-711e5814459f','146097e6-719b-458a-afb1-e936119b97a6',NULL,NULL,NULL);
/*!40000 ALTER TABLE `USER_ATTRIBUTE` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `USER_CONSENT`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `USER_CONSENT` (
  `ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `CLIENT_ID` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `USER_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `CREATED_DATE` bigint DEFAULT NULL,
  `LAST_UPDATED_DATE` bigint DEFAULT NULL,
  `CLIENT_STORAGE_PROVIDER` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `EXTERNAL_CLIENT_ID` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `UK_LOCAL_CONSENT` (`CLIENT_ID`,`USER_ID`),
  UNIQUE KEY `UK_EXTERNAL_CONSENT` (`CLIENT_STORAGE_PROVIDER`,`EXTERNAL_CLIENT_ID`,`USER_ID`),
  KEY `IDX_USER_CONSENT` (`USER_ID`),
  CONSTRAINT `FK_GRNTCSNT_USER` FOREIGN KEY (`USER_ID`) REFERENCES `USER_ENTITY` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `USER_CONSENT` WRITE;
/*!40000 ALTER TABLE `USER_CONSENT` DISABLE KEYS */;
/*!40000 ALTER TABLE `USER_CONSENT` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `USER_CONSENT_CLIENT_SCOPE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `USER_CONSENT_CLIENT_SCOPE` (
  `USER_CONSENT_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `SCOPE_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`USER_CONSENT_ID`,`SCOPE_ID`),
  KEY `IDX_USCONSENT_CLSCOPE` (`USER_CONSENT_ID`),
  KEY `IDX_USCONSENT_SCOPE_ID` (`SCOPE_ID`),
  CONSTRAINT `FK_GRNTCSNT_CLSC_USC` FOREIGN KEY (`USER_CONSENT_ID`) REFERENCES `USER_CONSENT` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `USER_CONSENT_CLIENT_SCOPE` WRITE;
/*!40000 ALTER TABLE `USER_CONSENT_CLIENT_SCOPE` DISABLE KEYS */;
/*!40000 ALTER TABLE `USER_CONSENT_CLIENT_SCOPE` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `USER_ENTITY`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `USER_ENTITY` (
  `ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `EMAIL` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `EMAIL_CONSTRAINT` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `EMAIL_VERIFIED` tinyint NOT NULL DEFAULT '0',
  `ENABLED` tinyint NOT NULL DEFAULT '0',
  `FEDERATION_LINK` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `FIRST_NAME` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `LAST_NAME` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `REALM_ID` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `USERNAME` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `CREATED_TIMESTAMP` bigint DEFAULT NULL,
  `SERVICE_ACCOUNT_CLIENT_LINK` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `NOT_BEFORE` int NOT NULL DEFAULT '0',
  PRIMARY KEY (`ID`),
  UNIQUE KEY `UK_DYKN684SL8UP1CRFEI6ECKHD7` (`REALM_ID`,`EMAIL_CONSTRAINT`),
  UNIQUE KEY `UK_RU8TT6T700S9V50BU18WS5HA6` (`REALM_ID`,`USERNAME`),
  KEY `IDX_USER_EMAIL` (`EMAIL`),
  KEY `IDX_USER_SERVICE_ACCOUNT` (`REALM_ID`,`SERVICE_ACCOUNT_CLIENT_LINK`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `USER_ENTITY` WRITE;
/*!40000 ALTER TABLE `USER_ENTITY` DISABLE KEYS */;
INSERT INTO `USER_ENTITY` VALUES ('19b6d072-21e8-40f9-9d22-711e5814459f','adolfo.agudelo6621@uco.net.co','adolfo.agudelo6621@uco.net.co',1,1,NULL,NULL,NULL,'5d56fd15-57bc-496f-8f3f-bf0d7a5c5581','motogo-admin',1765993612338,NULL,0),('79c7eea4-7f93-4b1e-ad06-b244ba9ca466','moto@yopmail.com','moto@yopmail.com',1,1,NULL,'usuario','cero uno','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','moto@yopmail.com',1773366116687,NULL,0),('d62612e2-e122-441d-8e71-7fee81f4c53c','admin@motogo.com','admin@motogo.com',1,1,NULL,'Esteban','Agudelo','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','admin',1768842410468,NULL,0),('fb214e14-d12c-451e-aca8-c62b09f77b11','taller@yopmail.com','taller@yopmail.com',1,1,NULL,'Taller','Cero uno','4c2886d8-e4c4-4b40-9064-445c7ab90a1c','taller@yopmail.com',1773366451869,NULL,0);
/*!40000 ALTER TABLE `USER_ENTITY` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `USER_FEDERATION_CONFIG`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `USER_FEDERATION_CONFIG` (
  `USER_FEDERATION_PROVIDER_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `VALUE` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `NAME` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`USER_FEDERATION_PROVIDER_ID`,`NAME`),
  CONSTRAINT `FK_T13HPU1J94R2EBPEKR39X5EU5` FOREIGN KEY (`USER_FEDERATION_PROVIDER_ID`) REFERENCES `USER_FEDERATION_PROVIDER` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `USER_FEDERATION_CONFIG` WRITE;
/*!40000 ALTER TABLE `USER_FEDERATION_CONFIG` DISABLE KEYS */;
/*!40000 ALTER TABLE `USER_FEDERATION_CONFIG` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `USER_FEDERATION_MAPPER`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `USER_FEDERATION_MAPPER` (
  `ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `NAME` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `FEDERATION_PROVIDER_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `FEDERATION_MAPPER_TYPE` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `REALM_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`ID`),
  KEY `IDX_USR_FED_MAP_FED_PRV` (`FEDERATION_PROVIDER_ID`),
  KEY `IDX_USR_FED_MAP_REALM` (`REALM_ID`),
  CONSTRAINT `FK_FEDMAPPERPM_FEDPRV` FOREIGN KEY (`FEDERATION_PROVIDER_ID`) REFERENCES `USER_FEDERATION_PROVIDER` (`ID`),
  CONSTRAINT `FK_FEDMAPPERPM_REALM` FOREIGN KEY (`REALM_ID`) REFERENCES `REALM` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `USER_FEDERATION_MAPPER` WRITE;
/*!40000 ALTER TABLE `USER_FEDERATION_MAPPER` DISABLE KEYS */;
/*!40000 ALTER TABLE `USER_FEDERATION_MAPPER` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `USER_FEDERATION_MAPPER_CONFIG`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `USER_FEDERATION_MAPPER_CONFIG` (
  `USER_FEDERATION_MAPPER_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `VALUE` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `NAME` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`USER_FEDERATION_MAPPER_ID`,`NAME`),
  CONSTRAINT `FK_FEDMAPPER_CFG` FOREIGN KEY (`USER_FEDERATION_MAPPER_ID`) REFERENCES `USER_FEDERATION_MAPPER` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `USER_FEDERATION_MAPPER_CONFIG` WRITE;
/*!40000 ALTER TABLE `USER_FEDERATION_MAPPER_CONFIG` DISABLE KEYS */;
/*!40000 ALTER TABLE `USER_FEDERATION_MAPPER_CONFIG` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `USER_FEDERATION_PROVIDER`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `USER_FEDERATION_PROVIDER` (
  `ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `CHANGED_SYNC_PERIOD` int DEFAULT NULL,
  `DISPLAY_NAME` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `FULL_SYNC_PERIOD` int DEFAULT NULL,
  `LAST_SYNC` int DEFAULT NULL,
  `PRIORITY` int DEFAULT NULL,
  `PROVIDER_NAME` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `REALM_ID` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`ID`),
  KEY `IDX_USR_FED_PRV_REALM` (`REALM_ID`),
  CONSTRAINT `FK_1FJ32F6PTOLW2QY60CD8N01E8` FOREIGN KEY (`REALM_ID`) REFERENCES `REALM` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `USER_FEDERATION_PROVIDER` WRITE;
/*!40000 ALTER TABLE `USER_FEDERATION_PROVIDER` DISABLE KEYS */;
/*!40000 ALTER TABLE `USER_FEDERATION_PROVIDER` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `USER_GROUP_MEMBERSHIP`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `USER_GROUP_MEMBERSHIP` (
  `GROUP_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `USER_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `MEMBERSHIP_TYPE` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`GROUP_ID`,`USER_ID`),
  KEY `IDX_USER_GROUP_MAPPING` (`USER_ID`),
  CONSTRAINT `FK_USER_GROUP_USER` FOREIGN KEY (`USER_ID`) REFERENCES `USER_ENTITY` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `USER_GROUP_MEMBERSHIP` WRITE;
/*!40000 ALTER TABLE `USER_GROUP_MEMBERSHIP` DISABLE KEYS */;
/*!40000 ALTER TABLE `USER_GROUP_MEMBERSHIP` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `USER_REQUIRED_ACTION`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `USER_REQUIRED_ACTION` (
  `USER_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `REQUIRED_ACTION` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT ' ',
  PRIMARY KEY (`REQUIRED_ACTION`,`USER_ID`),
  KEY `IDX_USER_REQACTIONS` (`USER_ID`),
  CONSTRAINT `FK_6QJ3W1JW9CVAFHE19BWSIUVMD` FOREIGN KEY (`USER_ID`) REFERENCES `USER_ENTITY` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `USER_REQUIRED_ACTION` WRITE;
/*!40000 ALTER TABLE `USER_REQUIRED_ACTION` DISABLE KEYS */;
/*!40000 ALTER TABLE `USER_REQUIRED_ACTION` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `USER_ROLE_MAPPING`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `USER_ROLE_MAPPING` (
  `ROLE_ID` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `USER_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`ROLE_ID`,`USER_ID`),
  KEY `IDX_USER_ROLE_MAPPING` (`USER_ID`),
  CONSTRAINT `FK_C4FQV34P1MBYLLOXANG7B1Q3L` FOREIGN KEY (`USER_ID`) REFERENCES `USER_ENTITY` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `USER_ROLE_MAPPING` WRITE;
/*!40000 ALTER TABLE `USER_ROLE_MAPPING` DISABLE KEYS */;
INSERT INTO `USER_ROLE_MAPPING` VALUES ('0c50da1d-eeb7-4074-b31e-f638245df45e','19b6d072-21e8-40f9-9d22-711e5814459f'),('a04cd05c-e599-4b52-8c7c-88e0ad734cae','19b6d072-21e8-40f9-9d22-711e5814459f'),('5125458c-9476-41c9-8fb6-0543f7083773','79c7eea4-7f93-4b1e-ad06-b244ba9ca466'),('da6b76cf-a450-4e8a-81f6-4286c9fa2ef7','79c7eea4-7f93-4b1e-ad06-b244ba9ca466'),('08e2bffc-779d-4dc5-b4fc-2b9c27ddf898','d62612e2-e122-441d-8e71-7fee81f4c53c'),('098e5980-f81f-4f7a-8eae-685dec71faa4','d62612e2-e122-441d-8e71-7fee81f4c53c'),('1302a0c9-cff5-4f38-a524-387fb61ecae5','d62612e2-e122-441d-8e71-7fee81f4c53c'),('1438fd35-067c-4657-a9ae-99394fcdddf4','d62612e2-e122-441d-8e71-7fee81f4c53c'),('14ba4893-70c6-4be5-9beb-225ca62ed5ec','d62612e2-e122-441d-8e71-7fee81f4c53c'),('15af1b7b-a8c8-46e7-b327-566cbfe05346','d62612e2-e122-441d-8e71-7fee81f4c53c'),('31f48475-6b51-42b4-8706-db2d938ed265','d62612e2-e122-441d-8e71-7fee81f4c53c'),('326e598f-eaae-4da1-acbf-daafb4789421','d62612e2-e122-441d-8e71-7fee81f4c53c'),('3dd13d47-55dd-4ac7-8e2a-c428fafc3739','d62612e2-e122-441d-8e71-7fee81f4c53c'),('49d17406-d5b7-4b46-801c-7aad314461cd','d62612e2-e122-441d-8e71-7fee81f4c53c'),('4fba4722-371b-47c9-a0ca-3a2b05b0f197','d62612e2-e122-441d-8e71-7fee81f4c53c'),('5125458c-9476-41c9-8fb6-0543f7083773','d62612e2-e122-441d-8e71-7fee81f4c53c'),('51fe34be-b8a7-4b86-8ff5-d036dcfae6a1','d62612e2-e122-441d-8e71-7fee81f4c53c'),('56de3a81-f2d8-4b89-ae11-f933ea5efcd8','d62612e2-e122-441d-8e71-7fee81f4c53c'),('5867419e-37e5-4166-ac97-92b7c53e6eb9','d62612e2-e122-441d-8e71-7fee81f4c53c'),('78eb646a-3134-4823-b2ce-daf30bbfcf9b','d62612e2-e122-441d-8e71-7fee81f4c53c'),('807232af-7777-454a-b1cf-82a74b4ddf4b','d62612e2-e122-441d-8e71-7fee81f4c53c'),('8868fe7b-aadb-4011-a14b-dea7ef1df1e1','d62612e2-e122-441d-8e71-7fee81f4c53c'),('9a839442-0a1f-42d6-8e91-2b8c6996e199','d62612e2-e122-441d-8e71-7fee81f4c53c'),('a3fa1093-7cb2-444c-82f9-dc04cf28ca42','d62612e2-e122-441d-8e71-7fee81f4c53c'),('b2a6016f-8d43-4ab7-be9d-7119c2b71b80','d62612e2-e122-441d-8e71-7fee81f4c53c'),('b723d2ac-460d-49a1-abbb-077d65b47e98','d62612e2-e122-441d-8e71-7fee81f4c53c'),('bc732c7b-1dda-4fcb-bbfd-d2c11138eebf','d62612e2-e122-441d-8e71-7fee81f4c53c'),('e7b39c0b-fef7-44a1-b209-df633fcbc8a5','d62612e2-e122-441d-8e71-7fee81f4c53c'),('eacb6d31-9ae2-4618-867c-2165b879f81d','d62612e2-e122-441d-8e71-7fee81f4c53c'),('ec508af1-4bcb-48f3-95af-d8e3f7efe47c','d62612e2-e122-441d-8e71-7fee81f4c53c'),('f13034bc-6106-42fb-950c-1467847ca692','d62612e2-e122-441d-8e71-7fee81f4c53c'),('f1b75db6-0157-406f-ade0-bc90cf859fda','d62612e2-e122-441d-8e71-7fee81f4c53c'),('ff7b90a1-7585-44e1-89c5-0a39963f236d','d62612e2-e122-441d-8e71-7fee81f4c53c'),('5125458c-9476-41c9-8fb6-0543f7083773','fb214e14-d12c-451e-aca8-c62b09f77b11'),('5ce52df1-959c-4528-9456-cec03ee9fb34','fb214e14-d12c-451e-aca8-c62b09f77b11');
/*!40000 ALTER TABLE `USER_ROLE_MAPPING` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `WEB_ORIGINS`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `WEB_ORIGINS` (
  `CLIENT_ID` varchar(36) COLLATE utf8mb4_unicode_ci NOT NULL,
  `VALUE` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`CLIENT_ID`,`VALUE`),
  KEY `IDX_WEB_ORIG_CLIENT` (`CLIENT_ID`),
  CONSTRAINT `FK_LOJPHO213XCX4WNKOG82SSRFY` FOREIGN KEY (`CLIENT_ID`) REFERENCES `CLIENT` (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `WEB_ORIGINS` WRITE;
/*!40000 ALTER TABLE `WEB_ORIGINS` DISABLE KEYS */;
INSERT INTO `WEB_ORIGINS` VALUES ('35ba46a4-b1e1-4dd0-8cff-fc98df4c2659','http://localhost:8085'),('3c6f327b-8018-44e9-8540-511170baeb2f','+'),('4c9063cc-607c-4f14-b55d-eb3d700e742c','+');
/*!40000 ALTER TABLE `WEB_ORIGINS` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `WORKFLOW_STATE`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `WORKFLOW_STATE` (
  `EXECUTION_ID` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `RESOURCE_ID` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `WORKFLOW_ID` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `WORKFLOW_PROVIDER_ID` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `RESOURCE_TYPE` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `SCHEDULED_STEP_ID` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `SCHEDULED_STEP_TIMESTAMP` bigint DEFAULT NULL,
  PRIMARY KEY (`EXECUTION_ID`),
  UNIQUE KEY `UQ_WORKFLOW_RESOURCE` (`WORKFLOW_ID`,`RESOURCE_ID`),
  KEY `IDX_WORKFLOW_STATE_STEP` (`WORKFLOW_ID`,`SCHEDULED_STEP_ID`),
  KEY `IDX_WORKFLOW_STATE_PROVIDER` (`RESOURCE_ID`,`WORKFLOW_PROVIDER_ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `WORKFLOW_STATE` WRITE;
/*!40000 ALTER TABLE `WORKFLOW_STATE` DISABLE KEYS */;
/*!40000 ALTER TABLE `WORKFLOW_STATE` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

