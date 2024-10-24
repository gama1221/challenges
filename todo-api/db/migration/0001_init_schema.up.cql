-- Create the keyspace
CREATE KEYSPACE IF NOT EXISTS todoapp WITH REPLICATION = { 'class': 'SimpleStrategy', 'replication_factor': 3 };

-- Use the keyspace
USE todoapp;

-- Create the users table
CREATE TABLE IF NOT EXISTS users (
    user_id TEXT PRIMARY KEY,
    username TEXT,
    email TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- Create the todos table
CREATE TABLE IF NOT EXISTS todos (
    id TEXT PRIMARY KEY,
    user_id TEXT,
    title TEXT,
    description TEXT,
    status TEXT,
    created TIMESTAMP,
    updated TIMESTAMP
);
