CREATE TABLE access_session
(
    id           character varying(26) PRIMARY KEY NOT NULL,
    version      character varying(26)             NOT NULL,
    user_id      character varying(26)             NOT NULL,
    namespace_id character varying(26)             NOT NULL,
    hashed_token character varying(128)            NOT NULL,
    scopes       character varying(256),
    -- TODO: context
    is_active    boolean                           NOT NULL,
    created_at   timestamp                         NOT NULL,
    updated_at   timestamp                         NOT NULL,
    expired_at   timestamp                         NOT NULL,

    CONSTRAINT session_version UNIQUE (version),
    FOREIGN KEY (user_id) REFERENCES users (user_id),
    FOREIGN KEY (namespace_id) REFERENCES namespaces (namespace_id)
);