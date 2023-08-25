-- +goose Up
-- +goose StatementBegin
CREATE TYPE status_order_enum AS ENUM ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');
CREATE TABLE "content"."order" (
    "id" uuid NOT NULL DEFAULT gen_random_uuid(),
    "number" varchar NOT NULL,
    "status" status_order_enum NOT NULL DEFAULT 'NEW'::status_order_enum,
    "user_id" uuid NOT NULL,
    "bonuses" int4,
    "updated_at" timestamptz NOT NULL,
    "created_at" timestamptz NOT NULL,
    PRIMARY KEY ("id")
);
ALTER TABLE ONLY "content"."order"
    ADD CONSTRAINT "order_number_uniq" UNIQUE (number);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "content"."order";
DROP TYPE status_order_enum;
-- +goose StatementEnd
