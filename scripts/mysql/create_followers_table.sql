-- Table: public.followers

-- DROP TABLE public.followers;

CREATE TABLE IF NOT EXISTS public.followers
(
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    username character varying(100) COLLATE pg_catalog."default" NOT NULL,
    following character varying(100) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT followers_pkey PRIMARY KEY (id),
    CONSTRAINT fk_following FOREIGN KEY (username)
        REFERENCES public.users (username) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT fk_user FOREIGN KEY (username)
        REFERENCES public.users (username) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)

TABLESPACE pg_default;

ALTER TABLE public.followers
    OWNER to postgres;