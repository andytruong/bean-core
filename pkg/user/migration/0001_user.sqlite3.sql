CREATE TABLE "user"
(
    id         uuid                   NOT NULL,
    version    uuid                   NOT NULL,
    is_active  boolean                NOT NULL,
    created_at timestamp              NOT NULL,
    updated_at timestamp              NOT NULL,
    deleted_at timestamp,
    avatar_uri character varying(255) NOT NULL,
    CONSTRAINT user_id PRIMARY KEY ("id"),
    CONSTRAINT user_version UNIQUE (version)
);

CREATE TABLE user_name
(
    user_id        uuid                  NOT NULL,
    first_name     character varying(64) NOT NULL,
    last_name      character varying(64) NOT NULL,
    preferred_name character varying(64) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES "user" (id)
);

CREATE INDEX user_name_fk ON user_name (user_id);

CREATE TABLE user_password
(
    user_id      uuid                   NOT NULL,
    is_active    boolean                NOT NULL,
    created_at   timestamp              NOT NULL,
    updated_at   timestamp              NOT NULL,
    algorithm    character varying(8)   NOT NULL,
    hashed_value character varying(255) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES "user" (id)
);

CREATE INDEX user_pass_fk ON user_password (user_id);
CREATE INDEX user_pass ON user_password (algorithm, hashed_value);
CREATE INDEX user_pass_status ON user_password (is_active);

CREATE TABLE user_email
(
    user_id    uuid                   NOT NULL,
    is_primary boolean                NOT NULL,
    is_active  boolean                NOT NULL,
    created_at timestamp              NOT NULL,
    updated_at timestamp              NOT NULL,
    value      character varying(128) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES "user" (id)
);

CREATE INDEX user_email_fk ON user_email (user_id);
CREATE INDEX user_email_lookup ON user_email (value);
