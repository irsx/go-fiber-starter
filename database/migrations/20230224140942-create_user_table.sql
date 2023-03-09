
-- +migrate Up
-- Table Definition ----------------------------------------------
CREATE TABLE "user" (
    guid uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    name character varying,
    email character varying,
    password character varying,
    phone_number character varying,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);

-- Indices -------------------------------------------------------
CREATE UNIQUE INDEX IF NOT EXISTS user_pkey ON "user"(guid uuid_ops);
CREATE UNIQUE INDEX IF NOT EXISTS user_ukey ON "user"(email text_ops);


-- +migrate Down
DROP TABLE "user";
