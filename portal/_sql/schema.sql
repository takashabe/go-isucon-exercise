DROP DATABASE IF EXISTS portal;
CREATE DATABASE IF NOT EXISTS portal;
use portal;


CREATE TABLE IF NOT EXISTS teams (
  `id` int NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `team` varchar(128) NOT NULL UNIQUE,
  `password` varchar(32) NOT NULL,
  `email` varchar(128) NOT NULL UNIQUE,
  `instance` varchar(128),
) DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS scores (
  `id` int NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `team_id` int NOT NULL,
  `summary` varchar(32) NOT NULL, -- success, fail
  `score` int NOT NULL,
  `submitted_at` timestamp,
  `json` text
) DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS highscores (
  `team_id` int NOT NULL PRIMARY KEY,
  `score` int NOT NULL,
  `submitted_at` timestamp
) DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS messages (
  `id` int NOT NULL AUTO_INCREMENT PRIMARY KEY,
  -- http://getbootstrap.com/components/#alerts
  `priority` varchar(16) DEFAULT 'alert-info', -- 'alert-success', 'alert-info', 'alert-warning', 'alert-danger'
  `content` TEXT NOT NULL,
  `show_at` timestamp NOT NULL,
  `hide_at` timestamp NOT NULL
) DEFAULT CHARSET=utf8mb4;
