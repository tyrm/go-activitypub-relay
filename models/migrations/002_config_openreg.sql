-- +migrate Up
INSERT INTO "public"."config" (key, value)
VALUES ('open_registration', 'false');

-- +migrate Down
