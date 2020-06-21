CREATE TABLE s3_application
(
    id             character varying(26) NOT NULL PRIMARY KEY,
    version        character varying(26) NOT NULL UNIQUE,
    is_active      boolean               NOT NULL,
    created_at     timestamp             NOT NULL,
    updated_at     timestamp             NOT NULL,
    deleted_at     timestamp,
    credentials_id character varying(26) NOT NULL
);

CREATE TABLE s3_credentials
(
    id                character varying(26)  NOT NULL PRIMARY KEY,
    version           character varying(26)  NOT NULL UNIQUE,
    endpoint          character varying(256) NOT NULL,
    access_key_id     character varying(128) NOT NULL,
    secret_access_key character varying(128) NOT NULL, -- encrypted value.
    is_secure         boolean                NOT NULL
);

CREATE TABLE s3_application_policy
(
    id             character varying(26) NOT NULL PRIMARY KEY,
    version        character varying(26) NOT NULL UNIQUE,
    application_id character varying(26),
    is_active      boolean               NOT NULL,
    created_at     timestamp             NOT NULL,
    updated_at     timestamp             NOT NULL,
    kind           character varying(32) CHECK ( value IN ('file_extension', 'rate_limit')),
    -- examples:
    --      * file extension: pdf txt zip gz
    --      * rate limit: 1MB/user/hour
    value          text
);

CREATE TABLE s3_upload_token
(
    id           character varying(26)  NOT NULL PRIMARY KEY,
    namespace_id character varying(26)  NOT NULL,
    user_id      character varying(26)  NOT NULL,
    file_path    CHARACTER VARYING(128) NOT NULL,
    created_at   timestamp              NOT NULL
);

CREATE TABLE s3_file
(
    id             character varying(26)    NOT NULL PRIMARY KEY,
    version        character varying(26)    NOT NULL UNIQUE,
    application_id character varying(26)    NOT NULL,
    size           float CHECK (value >= 0) NOT NULL, -- in byte
    path           character varying(128)   NOT NULL,
    is_active      boolean                  NOT NULL,
    created_at     timestamp                NOT NULL
);
