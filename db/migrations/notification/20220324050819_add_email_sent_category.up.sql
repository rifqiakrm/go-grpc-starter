BEGIN;

ALTER TABLE notification.email_sent
    ADD COLUMN "category" VARCHAR(255) null;

COMMIT;