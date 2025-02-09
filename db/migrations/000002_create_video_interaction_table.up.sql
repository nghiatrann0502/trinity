CREATE TABLE IF NOT EXISTS `video_interaction` (
  `id` int NOT NULL AUTO_INCREMENT,
  `video_id` int DEFAULT NULL,
  `action` varchar(10) DEFAULT NULL,
  `value` float DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `video_interaction_video_id_idx` (`video_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
