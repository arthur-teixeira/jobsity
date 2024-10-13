CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  email VARCHAR(255) NOT NULL,
  password bytea NOT NULL,
  salt bytea NOT NULL
);

CREATE TABLE IF NOT EXISTS tasks (
  id SERIAL PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  is_completed BIT NOT NULL DEFAULT 0::bit,
  user_id INTEGER NOT NULL REFERENCES users (id) 
);
