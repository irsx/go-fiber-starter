
-- +migrate Up
-- Table Definition ----------------------------------------------
CREATE TABLE "product" (
    guid uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    user_guid uuid REFERENCES "user"(guid),
    barcode_id character varying,
    sku character varying,
    name character varying,
    description text,
    image character varying,
    stock integer DEFAULT 0,
    sell_price numeric,
    buy_price numeric,
    expired_at date,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);

-- Indices -------------------------------------------------------
CREATE UNIQUE INDEX IF NOT EXISTS product_pkey ON product(guid uuid_ops);
CREATE UNIQUE INDEX IF NOT EXISTS product_ukey ON product(barcode_id text_ops,sku text_ops);

-- +migrate Down
DROP TABLE "product";
