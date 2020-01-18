-- +migrate Up
CREATE TABLE "public"."instances" (
    id serial NOT NULL UNIQUE,
    hostname character varying NOT NULL,
    joined_at timestamp without time zone NOT NULL DEFAULT current_timestamp,
    approved_at timestamp without time zone,
    PRIMARY KEY ("id")
)
;

CREATE INDEX idx_instances_hostname ON instances(hostname);

-- +migrate Down
DROP INDEX idx_instances_hostname;
DROP TABLE "public"."instances";