CREATE TABLE namespaces
(
    id         character varying(26)  NOT NULL,
    version    character varying(26)  NOT NULL UNIQUE,
    is_active  boolean                NOT NULL,
    created_at timestamp              NOT NULL,
    updated_at timestamp              NOT NULL,
    deleted_at timestamp,
    title      character varying(255) NOT NULL,
    CONSTRAINT namespace_id PRIMARY KEY (id)
);

CREATE TABLE namespace_domains
(
    id           character varying(26)  NOT NULL PRIMARY KEY,
    namespace_id character varying(26)  NOT NULL,
    is_primary   boolean                NOT NULL,
    is_active    boolean                NOT NULL,
    is_verified  boolean                NOT NULL,
    created_at   timestamp              NOT NULL,
    updated_at   timestamp              NOT NULL,
    value        character varying(256) NOT NULL UNIQUE,
    FOREIGN KEY (namespace_id) REFERENCES namespaces (id)
);

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

CREATE TABLE namespace_memberships
(
    id           character varying(26) NOT NULL PRIMARY KEY,
    version      character varying(26) NOT NULL UNIQUE,
    namespace_id character varying(26) NOT NULL,
    user_id      character varying(26) NOT NULL,
    is_active    boolean               NOT NULL,
    created_at   timestamp             NOT NULL,
    updated_at   timestamp             NOT NULL,
    FOREIGN KEY (namespace_id) REFERENCES namespaces (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);
