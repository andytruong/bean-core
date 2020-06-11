CREATE TABLE config_buckets
(
    id           character varying(26) NOT NULL PRIMARY KEY,
    version      character varying(26) NOT NULL UNIQUE,
    slug         character varying(64) NOT NULL UNIQUE,
    title        character varying(256),
    description  text,
    host_id      character varying(26) NOT NULL,

    -- ---------------------
    -- Access mode that apply to variables inside.
    --
    -- Means:
    --    100: Private read-only
    --    111: Public read
    --    400: Private writable variables
    access       character varying(4)  NOT NULL,

    -- ---------------------
    -- User can change this from off to one
    --  but can't change from on to off
    --
    -- When off: user can't create variable inside the bucket
    -- When on:  schema can not changed.
    is_published boolean               NOT NULL,

    -- ---------------------
    -- JSON schema
    schema       json                  NOT NULL,
    created_at   timestamp             NOT NULL,
    updated_at   timestamp             NOT NULL
);

CREATE INDEX config_buckets_fk ON config_buckets USING hash (host_id);

CREATE TABLE config_variables
(
    id          character varying(26) NOT NULL PRIMARY KEY,
    version     character varying(26) NOT NULL UNIQUE,
    bucket_id   character varying(26) NOT NULL,
    name        character varying(64) NOT NULL,
    description text,
    value       json                  NOT NULL,

    -- ---------------------
    -- User can't update variable if it's locked
    is_locked   boolean               NOT NULL,
    created_at  timestamp             NOT NULL,
    updated_at  timestamp             NOT NULL,
    CONSTRAINT config_config_unique UNIQUE (bucket_id, name),
    FOREIGN KEY (bucket_id) REFERENCES config_buckets (id)
);
