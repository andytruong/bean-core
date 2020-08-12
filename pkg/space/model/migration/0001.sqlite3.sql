CREATE TABLE spaces
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
    FOREIGN KEY (parent_id) REFERENCES spaces (id)
);

CREATE TABLE space_domains
(
    id           character varying(26)  NOT NULL PRIMARY KEY,
    space_id character varying(26)  NOT NULL,
    is_primary   boolean                NOT NULL,
    is_active    boolean                NOT NULL,
    is_verified  boolean                NOT NULL,
    created_at   timestamp              NOT NULL,
    updated_at   timestamp              NOT NULL,
    value        character varying(256) NOT NULL UNIQUE,
    FOREIGN KEY (space_id) REFERENCES spaces (id)
);

CREATE INDEX space_domains_value ON space_domains (value);

CREATE TABLE space_config
(
    id           character varying(26)  NOT NULL PRIMARY KEY,
    version      character varying(26)  NOT NULL UNIQUE,
    space_id character varying(26)  NOT NULL,
    bucket       character varying(128) NOT NULL,
    key          character varying(128) NOT NULL,
    value        json                   NOT NULL,
    created_at   timestamp              NOT NULL,
    updated_at   timestamp              NOT NULL,
    FOREIGN KEY (space_id) REFERENCES spaces (id),
    CONSTRAINT space_config_unique UNIQUE (space_id, bucket, key)
);

CREATE TABLE space_memberships
(
    id           character varying(26) NOT NULL PRIMARY KEY,
    version      character varying(26) NOT NULL UNIQUE,
    space_id character varying(26) NOT NULL,
    user_id      character varying(26) NOT NULL,
    is_active    boolean               NOT NULL,
    created_at   timestamp             NOT NULL,
    updated_at   timestamp             NOT NULL,
    logged_in_at timestamp,
    FOREIGN KEY (space_id) REFERENCES spaces (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE INDEX space_memberships_login ON space_memberships (logged_in_at);

CREATE TABLE space_manager_edge
(
    id                character varying(26) NOT NULL PRIMARY KEY,
    version           character varying(26) NOT NULL UNIQUE,
    user_member_id    character varying(26) NOT NULL,
    manager_member_id character varying(26) NOT NULL,
    is_active         BOOLEAN               NOT NULL,
    created_at        timestamp             NOT NULL,
    updated_at        timestamp             NOT NULL,
    FOREIGN KEY (user_member_id) REFERENCES space_memberships (id),
    FOREIGN KEY (manager_member_id) REFERENCES space_memberships (id)
);
