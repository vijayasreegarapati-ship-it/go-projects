-- migrations/000001_init_schema.up.sql

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'employee', -- 'employee' or 'admin'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE resources (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    resource_type VARCHAR(50) NOT NULL, -- e.g., 'desk', 'meeting_room'
    status VARCHAR(50) NOT NULL DEFAULT 'active', -- 'active', 'maintenance'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);