-- Create the database
CREATE DATABASE IF NOT EXISTS s276_TelemetryData;
USE s276_TelemetryData;

-- Table: sensor_types
CREATE TABLE IF NOT EXISTS sensor_types (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL
);

-- Table: sensor_data (for scalar values like temperature, humidity, pressure)
CREATE TABLE IF NOT EXISTS sensor_data (
    id INT AUTO_INCREMENT PRIMARY KEY,
    sensor_type_id INT NOT NULL,
    value DOUBLE NOT NULL,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sensor_type_id) REFERENCES sensor_types(id)
);

-- Table: vector_data (for x, y, z vectors like acceleration, gyroscope, gps)
CREATE TABLE IF NOT EXISTS vector_data (
    id INT AUTO_INCREMENT PRIMARY KEY,
    sensor_type_id INT NOT NULL,
    x DOUBLE NOT NULL,
    y DOUBLE NOT NULL,
    z DOUBLE NOT NULL,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sensor_type_id) REFERENCES sensor_types(id)
);

-- Pre-fill sensor_types
INSERT IGNORE INTO sensor_types (name) VALUES 
('temperature'),
('humidity'),
('pressure'),
('acceleration'),
('gyroscope'),
('gps');

