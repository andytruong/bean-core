CREATE TABLE namespaces
(
    id         character varying(26)  NOT NULL PRIMARY KEY,
    parent_id  character varying(26),
    version    character varying(26)  NOT NULL UNIQUE,
    kind       character varying(26)  NOT NULL,
    is_active  boolean                NOT NULL,
    created_at timestamp              NOT NULL,
    updated_at timestamp              NOT NULL,
    deleted_at timestamp,
    title      character varying(255) NOT NULL,
    language   character varying(16),
    FOREIGN KEY (parent_id) REFERENCES namespaces (id)
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
    logged_in_at timestamp,
    FOREIGN KEY (namespace_id) REFERENCES namespaces (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);

-- last login
CREATE INDEX namespace_memberships_login ON namespace_memberships USING btree (logged_in_at DESC NULLS LAST);

CREATE TABLE namespace_manager
(
    id                character varying(26) NOT NULL PRIMARY KEY,
    version           character varying(26) NOT NULL UNIQUE,
    user_member_id    character varying(26) NOT NULL,
    manager_member_id character varying(26) NOT NULL,
    is_active         BOOLEAN               NOT NULL,
    created_at        timestamp             NOT NULL,
    updated_at        timestamp             NOT NULL,
    FOREIGN KEY (user_member_id) REFERENCES namespace_memberships (id),
    FOREIGN KEY (manager_member_id) REFERENCES namespace_memberships (id)
);
