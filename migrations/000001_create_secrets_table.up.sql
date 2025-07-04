-- Create secrets table
CREATE TABLE IF NOT EXISTS secrets (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
  value TEXT NOT NULL,
  description TEXT
);
