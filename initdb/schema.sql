CREATE TABLE IF NOT EXISTS messages (
	id SERIAL PRIMARY KEY,
	sender TEXT NOT NULL,
	receiver TEXT NOT NULL,
	content TEXT NOT NULL,
	timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS groups (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS group_members (
  group_id INTEGER REFERENCES groups(id) ON DELETE CASCADE,
  user_id TEXT NOT NULL,
  PRIMARY KEY (group_id, user_id)
);

CREATE TABLE IF NOT EXISTS group_messages (
  id SERIAL PRIMARY KEY,
  group_id INTEGER REFERENCES groups(id) ON DELETE CASCADE,
  sender TEXT NOT NULL,
  content TEXT NOT NULL,
  timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Adicionar índices para melhorar a performance
CREATE INDEX IF NOT EXISTS idx_sender ON messages(sender);
CREATE INDEX IF NOT EXISTS idx_receiver ON messages(receiver);


INSERT INTO groups (id, name) VALUES (1, 'Grupo Teste');
INSERT INTO group_members (group_id, user_id) VALUES (1, 'joao'), (1, 'maria');