CREATE TABLE config_bucket
(
    id         character varying(26) NOT NULL PRIMARY KEY,
    version    character varying(26) NOT NULL UNIQUE,
    slug       character varying(64) NOT NULL UNIQUE,
    host_id    character varying(26) NOT NULL,
    user_id    character varying(26) NOT NULL,
    access     character varying(4)  NOT NULL,
    created_at timestamp             NOT NULL,
    updated_at timestamp             NOT NULL
);

CREATE TABLE config_record
(
    id         character varying(26) NOT NULL PRIMARY KEY,
    version    character varying(26) NOT NULL UNIQUE,
    bucket_id  character varying(26) NOT NULL,
    key        character varying(64) NOT NULL,
    value      json                  NOT NULL,
    created_at timestamp             NOT NULL,
    updated_at timestamp             NOT NULL,
    CONSTRAINT config_config_unique UNIQUE (bucket_id, key)
);
