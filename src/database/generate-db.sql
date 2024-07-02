-- Execute on MySQL using source /path/to/generate-db.sql

-- Create the database
DROP DATABASE IF EXISTS ehealth;
CREATE DATABASE ehealth;
USE ehealth;

-- Create the tables

-- Users table
CREATE TABLE IF NOT EXISTS user (
    id INT AUTO_INCREMENT,
    category ENUM('patient', 'practitioner', 'admin') NOT NULL,
    password_hash CHAR(100) NOT NULL,
    user_hash CHAR(100) NOT NULL,
    permissions JSON NOT NULL,
    PRIMARY KEY (id)
);

-- Patient table
CREATE TABLE IF NOT EXISTS patient (
    id INT AUTO_INCREMENT,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    date_of_birth DATE NOT NULL,
    street_address VARCHAR(100) NOT NULL,
    contact_number VARCHAR(20) NOT NULL,
    email VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

-- Practitioner table
CREATE TABLE IF NOT EXISTS practitioner (
    id INT AUTO_INCREMENT,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    date_of_birth DATE NOT NULL,
    street_address VARCHAR(100) NOT NULL,
    contact_number VARCHAR(20) NOT NULL,
    email VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

-- Location table
CREATE TABLE IF NOT EXISTS location (
    id INT AUTO_INCREMENT,
    location_name VARCHAR(100) NOT NULL,
    street_address VARCHAR(100) NOT NULL,
    contact_number VARCHAR(20) NOT NULL,
    email VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

-- Record table
CREATE TABLE IF NOT EXISTS record (
    patient_id INT,
    record_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    practitioner_id INT,
    location_id INT,
    record_type ENUM('check-up', 'annual', 'blood-work', 'vaccination', 'emergency'),
    notes TEXT,
    code_cpt VARCHAR(200),
    code_icd VARCHAR(200),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (patient_id, record_date)
    -- FOREIGN KEY (patient_id) REFERENCES patient(id),
    -- FOREIGN KEY (practitioner_id) REFERENCES practitioner(id),
    -- FOREIGN KEY (location_id) REFERENCES location(id)
);

-- Client table
CREATE TABLE IF NOT EXISTS client (
    patient_id INT,
    practitioner_id INT,
    client_status ENUM('active', 'inactive') NOT NULL,
    visits INT NOT NULL,
    first_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (patient_id, practitioner_id),
    FOREIGN KEY (patient_id) REFERENCES patient(id),
    FOREIGN KEY (practitioner_id) REFERENCES practitioner(id)
);