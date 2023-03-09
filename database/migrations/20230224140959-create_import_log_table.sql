
-- +migrate Up
-- Table Definition ----------------------------------------------
CREATE TABLE "import_log" (
    guid uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    user_guid uuid REFERENCES "user"(guid),
    exec_time numeric,
    total_success bigint,
    total_error bigint,
    errors text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);

-- Indices -------------------------------------------------------
CREATE UNIQUE INDEX IF NOT EXISTS import_log_pkey ON "import_log"(guid uuid_ops);


-- +migrate Down
DROP TABLE "import_log";
