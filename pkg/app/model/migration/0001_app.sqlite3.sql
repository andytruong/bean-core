CREATE TABLE applications
(
    id         character varying(26) NOT NULL PRIMARY KEY,
    version    character varying(26) NOT NULL UNIQUE,
    is_active  boolean               NOT NULL,
    title      character varying(255),
    created_at timestamp             NOT NULL,
    updated_at timestamp             NOT NULL,
    deleted_at timestamp
);

-- TODO: Missing indexes
