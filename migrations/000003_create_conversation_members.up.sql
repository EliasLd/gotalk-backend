CREATE TABLE conversation_members (
	user_id UUID NOT NULL,
	conversation_id UUID NOT NULL,
	joined_at TIMESTAMP NOT NULL DEFAULT now(),

	PRIMARY KEY (user_id, conversation_id),
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
	FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE
)

