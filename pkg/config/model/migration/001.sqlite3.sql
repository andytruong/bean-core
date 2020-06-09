CREATE TABLE namespace_config
(
    id           character varying(26)  NOT NULL PRIMARY KEY,
    version      character varying(26)  NOT NULL UNIQUE,
    namespace_id character varying(26)  NOT NULL,
    bucket       character varying(128) NOT NULL,
    key          character varying(128) NOT NULL,
    value        json                   NOT NULL,
    created_at   timestamp              NOT NULL,
    updated_at   timestamp              NOT NULL,
    FOREIGN KEY (namespace_id) REFERENCES namespaces (id),
    CONSTRAINT namespace_config_unique UNIQUE (namespace_id, bucket, key)
);
