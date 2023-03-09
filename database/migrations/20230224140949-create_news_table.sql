
-- +migrate Up
-- Table Definition ----------------------------------------------
CREATE TABLE "news" (
    guid uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    user_guid uuid REFERENCES "user"(guid),
    title character varying,
    description text,
    image character varying,
    hyper_link character varying,
    status smallint DEFAULT '1'::smallint,
    created_at timestamp with time zone,
    updated_at timestamp with time zone

);

-- Indices -------------------------------------------------------
CREATE UNIQUE INDEX IF NOT EXISTS news_pkey ON news(guid uuid_ops);


-- +migrate Down
DROP TABLE "news";
