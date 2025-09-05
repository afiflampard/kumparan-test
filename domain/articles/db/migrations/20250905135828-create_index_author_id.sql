
-- +migrate Up
CREATE INDEX idx_articles_author_id ON articles (author_id);

-- +migrate Down
DROP INDEX idx_articles_author_id;
