CREATE TABLE namespaces
(
    id         character varying(26)  NOT NULL,
    version    character varying(26)  NOT NULL,
    is_active  boolean                NOT NULL,
    created_at timestamp              NOT NULL,
    updated_at timestamp              NOT NULL,
    deleted_at timestamp,
    avatar_uri character varying(255) NOT NULL,
    CONSTRAINT namespace_id PRIMARY KEY (id),
    CONSTRAINT namespace_version UNIQUE (version)
);

CREATE TABLE namespace_domains
(
    id           character varying(26)  NOT NULL PRIMARY KEY,
    namespace_id character varying(26)  NOT NULL,
    is_primary   boolean                NOT NULL,
    is_active    boolean                NOT NULL,
    created_at   timestamp              NOT NULL,
    updated_at   timestamp              NOT NULL,
    value        character varying(256) NOT NULL,
    FOREIGN KEY (namespace_id) REFERENCES namespaces (id),
    CONSTRAINT user_unique_email UNIQUE (value)
)

-- domain
-- membership
