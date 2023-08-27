-- +goose Up
-- +goose StatementBegin
CREATE TABLE "content"."account" (
    "id" uuid NOT NULL DEFAULT gen_random_uuid(),
    "user_id" uuid NOT NULL,
    "current" int4,
    "withdrawn" int4,
    "updated_at" timestamptz NOT NULL,
    "created_at" timestamptz NOT NULL,
    PRIMARY KEY ("id")
);
ALTER TABLE ONLY "content"."account"
    ADD CONSTRAINT "account_user_id_uniq" UNIQUE (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "content"."account";
-- +goose StatementEnd
