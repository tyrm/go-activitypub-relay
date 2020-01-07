-- +migrate Up
CREATE TABLE "public"."config" (
    id serial NOT NULL UNIQUE,
    key character varying NOT NULL UNIQUE,
    value character varying,
    created_at timestamp without time zone NOT NULL DEFAULT current_timestamp,
    updated_at timestamp without time zone NOT NULL DEFAULT current_timestamp,
    PRIMARY KEY ("id")
)
;

-- +migrate Down
DROP TABLE "public"."config";