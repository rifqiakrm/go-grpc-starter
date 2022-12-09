BEGIN;
CREATE TABLE IF NOT EXISTS notification.email_sent
(
    created_by VARCHAR(200) NOT NULL,
    updated_by VARCHAR(200) NOT NULL,
    deleted_by VARCHAR(200),
    created_at TIMESTAMPTZ  NOT NULL,
    updated_at TIMESTAMPTZ  NOT NULL,
    deleted_at TIMESTAMPTZ,
    id         BIGSERIAL PRIMARY KEY,
    "m_id"     VARCHAR(200) NULL,
    "from"     VARCHAR(200) NOT NULL,
    "to"       VARCHAR(200) NOT NULL,
    "subject"  VARCHAR(200) NOT NULL,
    content    TEXT         NOT NULL,
    "status"   VARCHAR(100) NOT NULL
);
COMMIT;