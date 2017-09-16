-- CREATE user 'portal';
-- GRANT ALL PRIVILEGES ON *.* TO 'portal'@'%';
-- CREATE DATABASE IF NOT EXISTS portal;
use portal;

CREATE TABLE IF NOT EXISTS teams (
  `id`       int          NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `name`     varchar(128) NOT NULL UNIQUE,
  `password` varchar(32)  NOT NULL,
  `email`    varchar(128) NOT NULL UNIQUE,
  `instance` varchar(128)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS queues (
  `id`           int          NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `team_id`      int          NOT NULL,
  `msg_id`       varchar(128) UNIQUE
  `submitted_at` timestamp
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS scores (
  `id`           int         NOT NULL  AUTO_INCREMENT PRIMARY KEY,
  `team_id`      int         NOT NULL,
  `summary`      varchar(32) NOT NULL, -- success, fail
  `score`        int         NOT NULL,
  `submitted_at` timestamp,
  `json`         text
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS highscores (
  `team_id`      int NOT NULL PRIMARY KEY,
  `score`        int NOT NULL,
  `submitted_at` timestamp
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
