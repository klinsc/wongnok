-- +goose Up
-- +goose StatementBegin
ALTER TABLE ratings
ADD COLUMN IF NOT EXISTS user_id VARCHAR(100) NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE ratings
DROP COLUMN IF EXISTS user_id;
-- +goose StatementEnd
