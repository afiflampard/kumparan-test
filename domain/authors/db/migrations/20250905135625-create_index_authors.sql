
-- +migrate Up
CREATE UNIQUE INDEX uq_authors_email ON authors (email);

-- +migrate Down
DROP INDEX uq_authors_email;

