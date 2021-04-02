CREATE OR REPLACE FUNCTION trigger_set_timestamp()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE users
(
    id                 BIGSERIAL PRIMARY KEY,
    email              VARCHAR(320)  NOT NULL UNIQUE,
    username           VARCHAR(50)   NOT NULL UNIQUE,
    first_name         VARCHAR(255)  NOT NULL,
    last_name          VARCHAR(255)  NOT NULL,
    password           VARCHAR(255)  NOT NULL,
    created_at         TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    is_email_confirmed BOOLEAN       NOT NULL DEFAULT FALSE,
    gender             INTEGER       NOT NULL DEFAULT 0,
    sexual_preferences INTEGER       NOT NULL DEFAULT 0,
    biography          VARCHAR(1000) NOT NULL DEFAULT '',
    tags               TEXT[]        NOT NULL DEFAULT ARRAY []::TEXT[],
    avatar_url         TEXT          NOT NULL DEFAULT '',
    pictures_url       TEXT[]        NOT NULL DEFAULT ARRAY []::TEXT[],
    likes_num          INTEGER       NOT NULL DEFAULT 0,
    views_num          INTEGER       NOT NULL DEFAULT 0,
    gps_position       TEXT          NOT NULL DEFAULT ''
);

CREATE TRIGGER set_timestamp
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp();

CREATE INDEX idx_users_tags ON users USING GIN (tags);
