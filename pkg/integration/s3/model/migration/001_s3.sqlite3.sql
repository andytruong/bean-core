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
