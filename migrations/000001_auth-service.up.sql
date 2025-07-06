CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS  TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
    username VARCHAR(255) NOT NULL UNIQUE,         
    password VARCHAR(255) NOT NULL,                
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP 
);

CREATE EXTENSION IF NOT EXISTS  TABLE accounts (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    refresh_token TEXT NOT NULL,
    user_agent TEXT NOT NULL,
    x_forwarded_for TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id)
        ON DELETE CASCADE
);