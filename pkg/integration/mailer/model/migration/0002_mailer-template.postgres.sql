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
