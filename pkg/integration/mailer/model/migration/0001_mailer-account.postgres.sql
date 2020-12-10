-- Make it possible for each space to have its own account.
-- Sender name can be customised
CREATE TABLE mailer_account
(
    id                            character varying(26)  NOT NULL PRIMARY KEY,
    version                       character varying(26)  NOT NULL UNIQUE,
    space_id                      character varying(26),
    status                        boolean                NOT NULL, -- values: inactive (0), unverified (-1), active (1)
    created_at                    timestamp              NOT NULL,
    updated_at                    timestamp              NOT NULL,
    deleted_at                    timestamp,
    sender_name                   character varying(128) NOT NULL,
    sender_email                  character varying(26)  NOT NULL,
    encrypted_connection          character varying(26)  NOT NULL, -- # example: encrypt('smtp.sendgrid.net:587?username=test@bean.qa&password=bean')
    attachment_size_limit         int CHECK ( attachment_size_limit > 0 ),
    attachment_size_limit_each    int CHECK ( mailer_account.attachment_size_limit_each > 0 ),
    attachment_allowed_file_types int,
    FOREIGN KEY (space_id) REFERENCES spaces (id)
);

CREATE UNIQUE INDEX mailer_account_unique_sender ON mailer_account (space_id, sender_email);
CREATE INDEX mailer_account_status ON mailer_account (status);

CREATE TABLE mailer_account_stream
(
    id         character varying(26) NOT NULL PRIMARY KEY,
    account_id character varying(26) NOT NULL,
    user_id    character varying(26) NOT NULL,
    key        character varying(26) NOT NULL, -- example: version 1
    payload    text                  NOT NULL,
    FOREIGN KEY (account_id) REFERENCES mailer_account (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);

ALTER TABLE mailer_account
    ADD FOREIGN KEY (version)
        REFERENCES mailer_account_stream (id);
