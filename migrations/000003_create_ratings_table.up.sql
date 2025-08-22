CREATE TABLE ratings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_email VARCHAR(255) NOT NULL,
    article_id UUID NOT NULL,
    scores INT NOT NULL CHECK (scores >= 1 AND scores <= 5),
    tags TEXT[],
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE (user_email, article_id),

    CONSTRAINT fk_user
        FOREIGN KEY(user_email)
        REFERENCES users(email)
        ON DELETE CASCADE,
    CONSTRAINT fk_article
        FOREIGN KEY(article_id)
        REFERENCES articles(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_ratings_user_email ON ratings(user_email);
CREATE INDEX idx_ratings_article_id ON ratings(article_id);