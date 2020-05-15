CREATE TABLE users
(
    id         character varying(26)  NOT NULL,
    version    character varying(26)  NOT NULL,
    is_active  boolean                NOT NULL,
    created_at timestamp              NOT NULL,
    updated_at timestamp              NOT NULL,
    deleted_at timestamp,
    avatar_uri character varying(255) NOT NULL,
    CONSTRAINT user_id PRIMARY KEY ("id"),
    CONSTRAINT user_version UNIQUE (version)
);

CREATE TABLE user_names
(
    id             character varying(26) NOT NULL PRIMARY KEY,
    user_id        character varying(26) NOT NULL,
    first_name     character varying(64) NOT NULL,
    last_name      character varying(64) NOT NULL,
    preferred_name character varying(64) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE INDEX user_name_fk ON user_names USING hash (user_id);

CREATE TABLE user_password
(
    id           character varying(26)  NOT NULL PRIMARY KEY,
    user_id      character varying(26)  NOT NULL,
    is_active    boolean                NOT NULL,
    created_at   timestamp              NOT NULL,
    updated_at   timestamp              NOT NULL,
    algorithm    character varying(8)   NOT NULL,
    hashed_value character varying(255) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE INDEX user_pass_fk ON user_password USING hash (user_id);
CREATE INDEX user_pass ON user_password USING btree (algorithm, hashed_value);
CREATE INDEX user_pass_status ON user_password USING hash (is_active);

CREATE TABLE user_emails
(
    id         character varying(26)  NOT NULL PRIMARY KEY,
    user_id    character varying(26)  NOT NULL,
    is_primary boolean                NOT NULL,
    is_active  boolean                NOT NULL,
    created_at timestamp              NOT NULL,
    updated_at timestamp              NOT NULL,
    value      character varying(128) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE INDEX user_email_fk ON user_emails USING hash (user_id);
CREATE INDEX user_email_lookup ON user_emails USING hash (value);

CREATE TABLE user_unverified_emails
(
    id         character varying(26)  NOT NULL PRIMARY KEY,
    user_id    character varying(26)  NOT NULL,
    is_primary boolean                NOT NULL,
    is_active  boolean                NOT NULL,
    created_at timestamp              NOT NULL,
    updated_at timestamp              NOT NULL,
    value      character varying(128) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE INDEX user_un_email_fk ON user_unverified_emails USING hash (user_id);
CREATE INDEX user_un_email_lookup ON user_unverified_emails USING hash (value);
