CREATE TABLE IF NOT EXISTS messages(
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    content TEXT,
    created_at TIMESTAMP DEFAULT now(),
    is_anon BOOLEAN DEFAULT FALSE
);