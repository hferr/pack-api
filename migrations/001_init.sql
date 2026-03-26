-- +goose Up
-- +goose StatementBegin
CREATE TABLE packs (
    id UUID PRIMARY KEY,
    size INT UNIQUE NOT NULL
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS packs;
-- +goose StatementEnd
