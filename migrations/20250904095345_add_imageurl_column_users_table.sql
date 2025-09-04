-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN image_url VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN IF EXISTS image_url;
-- +goose StatementEnd
