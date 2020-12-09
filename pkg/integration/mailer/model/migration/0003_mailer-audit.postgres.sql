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
