CREATE TABLE "public"."user"
(
    id         uuid                   NOT NULL,
    version    uuid                   NOT NULL,
    is_active  boolean                NOT NULL,
    created_at timestamp              NOT NULL,
    updated_at timestamp              NOT NULL,
    deleted_at timestamp,
    avatar_uri character varying(255) NOT NULL,
    CONSTRAINT "user_id" PRIMARY KEY ("id"),
    CONSTRAINT "user_version" UNIQUE (version)
);

CREATE TABLE public.user_name
(
    user_id        uuid                  NOT NULL,
    first_name     character varying(64) NOT NULL,
    last_name      character varying(64) NOT NULL,
    preferred_name character varying(64) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES "user" (id)
);

CREATE INDEX "user_name_fk" ON public.user_name USING hash (user_id);

CREATE TABLE public.user_password
(
    user_id      uuid                   NOT NULL,
    is_active    boolean                NOT NULL,
    created_at   timestamp              NOT NULL,
    updated_at   timestamp              NOT NULL,
    algorithm    character varying(8)   NOT NULL,
    hashed_value character varying(255) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES "user" (id)
);

CREATE INDEX "user_pass_fk" ON public.user_password USING hash (user_id);
CREATE INDEX "user_pass" ON public.user_password USING hash (algorithm, hashed_value);
CREATE INDEX "user_pass_status" ON public.user_password USING hash (is_active);

CREATE TABLE public.user_email
(
    user_id    uuid                   NOT NULL,
    is_primary boolean                NOT NULL,
    is_active  boolean                NOT NULL,
    created_at timestamp              NOT NULL,
    updated_at timestamp              NOT NULL,
    value      character varying(128) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES "user" (id)
);

CREATE INDEX "user_email_fk" ON public.user_email USING hash (user_id);
CREATE INDEX "user_email" ON public.user_email USING hash (value);
