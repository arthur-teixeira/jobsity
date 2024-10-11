CREATE TABLE IF NOT EXISTS tasks (
  id SERIAL PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  is_completed BIT NOT NULL DEFAULT 0::bit,
  ts tsvector GENERATED ALWAYS AS (to_tsvector('english', title)) STORED
)

CREATE INDEX ts_idx ON tasks USING GIN (ts);
