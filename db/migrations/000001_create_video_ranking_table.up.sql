CREATE TABLE IF NOT EXISTS `video_ranking` (
  `id` int NOT NULL AUTO_INCREMENT,
  `video_id` int NOT NULL,
  `score` float NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_video_id` (`video_id`),
  KEY `idx_score_desc` (`score` DESC),
  KEY `idx_score_video` (`score` DESC,`video_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
