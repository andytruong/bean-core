CREATE TABLE applications
(
    id         character varying(26) NOT NULL PRIMARY KEY,
    version    character varying(26) NOT NULL UNIQUE,
    space_id   character varying(26) NOT NULL,
    is_active  boolean               NOT NULL,
    title      character varying(255),
    created_at timestamp             NOT NULL,
    updated_at timestamp             NOT NULL,
    deleted_at timestamp,
    FOREIGN KEY (space_id) REFERENCES spaces (id)
);

CREATE INDEX application_status ON applications (is_active);
CREATE INDEX application_deleted ON applications ((deleted_at IS NULL));
