-- +goose UP
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pingpong_counter (
    id BIGSERIAL PRIMARY KEY,
    count INTEGER NOT NULL
);
--- +goose StatementEnd

--- +goose Down
--- +goose StatementBegin
DROP TABLE pingpong;
--- +goose StatementEnd
