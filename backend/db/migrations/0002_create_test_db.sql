-- Create a separate database for tests and load the same schema/seeds
-- This runs only on first container init (empty data dir)

-- Create test database owned by the same user
CREATE DATABASE learnlang_test;

-- Connect to the test database and apply the initial schema/seeds
\connect learnlang_test
\i /docker-entrypoint-initdb.d/0001_init.sql
