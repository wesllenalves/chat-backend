CREATE TABLE IF NOT EXISTS messages (
	id SERIAL PRIMARY KEY,
	sender TEXT NOT NULL,
	receiver TEXT NOT NULL,
	content TEXT NOT NULL,
	timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Adicionar índices para melhorar a performance
CREATE INDEX IF NOT EXISTS idx_sender ON messages(sender);
CREATE INDEX IF NOT EXISTS idx_receiver ON messages(receiver);