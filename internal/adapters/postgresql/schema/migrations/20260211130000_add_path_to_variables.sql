-- +goose Up
ALTER TABLE variables ADD COLUMN path VARCHAR(500) NOT NULL DEFAULT '.env';

-- +goose Down
ALTER TABLE variables DROP COLUMN path;
