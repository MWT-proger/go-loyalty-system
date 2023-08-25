-- +goose Up
-- +goose StatementBegin
CREATE TABLE "content"."withdrawal" (
    "id" uuid NOT NULL,
    "number" varchar NOT NULL,
    "user_id" uuid NOT NULL,
    "bonuses" int4,
    "updated_at" timestamptz NOT NULL,
    "created_at" timestamptz NOT NULL,
    PRIMARY KEY ("id")
);
ALTER TABLE ONLY "content"."withdrawal"
    ADD CONSTRAINT "withdrawal_number_uniq" UNIQUE (number);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "content"."withdrawal";
-- +goose StatementEnd
