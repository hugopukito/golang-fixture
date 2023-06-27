-- If the database does not exist, create it
CREATE DATABASE IF NOT EXISTS `{{.DatabaseName}}`;

-- If the database exists, drop and recreate it
DROP DATABASE IF EXISTS `{{.DatabaseName}}`;
CREATE DATABASE `{{.DatabaseName}}`;