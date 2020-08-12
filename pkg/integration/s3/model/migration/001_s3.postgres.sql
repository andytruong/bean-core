CREATE TABLE s3_application
(
    id         character varying(26) NOT NULL PRIMARY KEY,
    version    character varying(26) NOT NULL UNIQUE,
    is_active  boolean               NOT NULL,
    created_at timestamp             NOT NULL,
    updated_at timestamp             NOT NULL,
    deleted_at timestamp
);

CREATE TABLE s3_credentials
(
    id             character varying(26)  NOT NULL PRIMARY KEY,
    application_id character varying(26)  NOT NULL UNIQUE,
    endpoint       character varying(256) NOT NULL,
    bucket         character varying(64)  NOT NULL,
    access_key     character varying(128) NOT NULL,
    secret_key     character varying(128) NOT NULL,
    is_secure      boolean                NOT NULL,
    FOREIGN KEY (application_id) REFERENCES s3_application (id)
);

CREATE TABLE s3_application_policy
(
    id             character varying(26) NOT NULL PRIMARY KEY,
    application_id character varying(26),
    created_at     timestamp             NOT NULL,
    updated_at     timestamp             NOT NULL,
    kind           character varying(32) CHECK ( kind IN ('file_extensions', 'rate_limit')),
    -- examples:
    --      * file extension: pdf txt zip gz
    --      * rate limit: 1MB/user/hour
    value          text
);

CREATE TABLE s3_upload_token
(
    id         character varying(26)  NOT NULL PRIMARY KEY,
    space_id   character varying(26)  NOT NULL,
    user_id    character varying(26)  NOT NULL,
    file_path  CHARACTER VARYING(128) NOT NULL,
    created_at timestamp              NOT NULL
);

-- TODO: Remove file schema of S3 object
CREATE TABLE s3_file
(
    id             character varying(26)   NOT NULL PRIMARY KEY,
    version        character varying(26)   NOT NULL UNIQUE,
    application_id character varying(26)   NOT NULL,
    size           float CHECK (size >= 0) NOT NULL, -- in byte
    path           character varying(128)  NOT NULL,
    is_active      boolean                 NOT NULL,
    created_at     timestamp               NOT NULL,
    updated_at     timestamp               NOT NULL
);
