CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE conversations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    is_public BOOLEAN NOT NULL,
    name TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);