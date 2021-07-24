-- Table: public.images

-- DROP TABLE public.images;

CREATE TABLE IF NOT EXISTS public.images
(
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    username character varying(100) COLLATE pg_catalog."default" NOT NULL,
    path text COLLATE pg_catalog."default" NOT NULL,
    created_on timestamp(6) without time zone NOT NULL DEFAULT now(),
    CONSTRAINT images_hash_p_pkey PRIMARY KEY (id),
    CONSTRAINT fk_username FOREIGN KEY (username)
        REFERENCES public.users (username) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);