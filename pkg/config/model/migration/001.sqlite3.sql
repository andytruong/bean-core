CREATE TABLE config_buckets
(
    id          character varying(26) NOT NULL PRIMARY KEY,
    version     character varying(26) NOT NULL UNIQUE,
    slug        character varying(64) NOT NULL UNIQUE,
    title       character varying(256),
    description text,
    host_id     character varying(26) NOT NULL,
    access      character varying(4)  NOT NULL,
    created_at  timestamp             NOT NULL,
    updated_at  timestamp             NOT NULL
);

CREATE INDEX config_buckets_fk ON config_buckets (host_id);

CREATE TABLE config_variables
(
    id          character varying(26) NOT NULL PRIMARY KEY,
    version     character varying(26) NOT NULL UNIQUE,
    bucket_id   character varying(26) NOT NULL,
    name        character varying(64) NOT NULL,
    description text,
    value       json                  NOT NULL,
    created_at  timestamp             NOT NULL,
    updated_at  timestamp             NOT NULL,
    CONSTRAINT config_config_unique UNIQUE (bucket_id, name)
);
