BEGIN;

CREATE TABLE IF NOT EXISTS users.users
(

    created_by   VARCHAR(200),
    updated_by   VARCHAR(200),
    deleted_by   VARCHAR(200),
    created_at   TIMESTAMP,
    updated_at   TIMESTAMP,
    deleted_at   TIMESTAMP,
    id           uuid PRIMARY KEY,
    username     VARCHAR(255) UNIQUE NULL,
    email        VARCHAR(255) UNIQUE NULL,
    phone_number VARCHAR(255) UNIQUE NULL,
    password     VARCHAR(255)        NULL
);

INSERT INTO "users"."users" ("created_by", "updated_by", "created_at", "updated_at", "id", "username", "email", "phone_number", "password") VALUES ('system', 'system', now(), now(), '0abc6437-bb96-4dc7-a8a1-04f4e3038741', 'rifqiakrm', 'rifqiakram57@gmail.com', '0895346419497', '$2a$12$xQxjo6oocTh1/vScQ/nA.eSkNa8C..ozVH7vJvN75OfXXhhXXQG3y');

COMMIT;
