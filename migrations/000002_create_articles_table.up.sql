CREATE TABLE articles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_email VARCHAR(255) NOT NULL,
    url VARCHAR(2048) NOT NULL,
    title VARCHAR(255) DEFAULT '',
    description TEXT DEFAULT '',
    image_url VARCHAR(2048) DEFAULT '',
    scrape_status VARCHAR(20) NOT NULL DEFAULT 'pending',
    retry_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT fk_user
        FOREIGN KEY(user_email) 
        REFERENCES users(email)
        ON DELETE CASCADE
);

CREATE INDEX idx_articles_user_email ON articles(user_email);
CREATE INDEX idx_articles_status_retry ON articles(scrape_status, retry_count);