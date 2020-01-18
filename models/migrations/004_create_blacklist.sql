-- +migrate Up
CREATE TABLE "public"."blacklist" (
    id serial NOT NULL UNIQUE,
    hostname character varying NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT current_timestamp,
    PRIMARY KEY ("id")
)
;

CREATE INDEX idx_blacklist_hostname ON blacklist(hostname);

-- +migrate Down
DROP INDEX idx_blacklist_hostname;
DROP TABLE "public"."blacklist";