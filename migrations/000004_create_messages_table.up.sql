CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE messages (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
	sender_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	--Text only for now
	content TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT now()
);
