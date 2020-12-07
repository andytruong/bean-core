-- Make it possible for each space to have its own account.
-- Sender name can be customised
CREATE TABLE mailer_account
(
    id                   character varying(26)  NOT NULL PRIMARY KEY,
    version              character varying(26)  NOT NULL UNIQUE,
    space_id             character varying(26),
    status               boolean                NOT NULL, -- values: inactive (0), unverified (-1), active (1)
    created_at           timestamp              NOT NULL,
    updated_at           timestamp              NOT NULL,
    deleted_at           timestamp,
    sender_name          character varying(128) NOT NULL,
    sender_email         character varying(26)  NOT NULL,
    encrypted_connection character varying(26)  NOT NULL,
    FOREIGN KEY (space_id) REFERENCES spaces (id)
);

CREATE UNIQUE INDEX mailer_account ON mailer_account (space_id, sender_email);
CREATE INDEX mailer_account_status ON mailer_account (status);

-- Instead of sending full email message directly
-- Each send-out email must have configured template.
CREATE TABLE mailer_template
(
    id         character varying(26)  NOT NULL PRIMARY KEY,
    version    character varying(26)  NOT NULL UNIQUE,
    space_id   character varying(26)  NOT NULL,
    is_active  boolean                NOT NULL,
    created_at timestamp              NOT NULL,
    updated_at timestamp              NOT NULL,
    deleted_at timestamp,
    title      character varying(255) NOT NULL,
    language   character varying(16),
    body_html  character              NOT NULL,
    body_text  character,
    FOREIGN KEY (space_id) REFERENCES spaces (id),
    FOREIGN KEY (version) REFERENCES mailer_template_stream (id)
);

CREATE TABLE mailer_template_stream
(
    id          character varying(26) NOT NULL PRIMARY KEY,
    template_id character varying(26) NOT NULL,
    user_id     character varying(26) NOT NULL,
    key         character varying(26) NOT NULL, -- example: version 1
    payload     text                  NOT NULL,
    FOREIGN KEY (template_id) REFERENCES mailer_template (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);

-- The table to answer the question: Was my email sent, when, any error, any warning?
-- We don't store raw text here, just a hash.
CREATE TABLE mailer_audit
(
    id              character varying(26) NOT NULL PRIMARY KEY,
    account_id      character varying(26) NOT NULL,
    span_id         character varying(26) NOT NULL,
    template_id     character varying(26) NOT NULL,
    created_at      timestamp             NOT NULL,
    recipient_hash  character varying(32) NOT NULL,
    context_hash    character varying(32) NOT NULL,
    error_code      integer,
    error_message   character varying(255),
    warning_message character varying(255),
    FOREIGN KEY (account_id) REFERENCES mailer_account (id),
    FOREIGN KEY (template_id) REFERENCES mailer_template (id)
);
