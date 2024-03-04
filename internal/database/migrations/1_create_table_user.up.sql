CREATE TABLE IF NOT EXISTS "user" (
  id SERIAL PRIMARY KEY,
  first_name VARCHAR(50),
  surname VARCHAR(50),
  email VARCHAR(50) UNIQUE,
  password TEXT NOT NULL,
  api_key TEXT UNIQUE,
  role VARCHAR(10),
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);
