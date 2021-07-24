-- Table: public.images

-- DROP TABLE public.images;

CREATE TABLE IF NOT EXISTS public.images
(
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    username character varying(100) COLLATE pg_catalog."default" NOT NULL,
    path text COLLATE pg_catalog."default" NOT NULL,
    created_on timestamp(6) without time zone NOT NULL DEFAULT now(),
    /*CONSTRAINT images_hash_p_pkey PRIMARY KEY (id),*/
    CONSTRAINT fk_username FOREIGN KEY (username)
        REFERENCES public.users (username) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
) PARTITION BY HASH (id);

ALTER TABLE public.images
    OWNER to postgres;

-- Partitions SQL

CREATE TABLE IF NOT EXISTS public.images_h0 PARTITION OF public.images
    FOR VALUES WITH (modulus 10, remainder 0);

ALTER TABLE public.images_h0
    OWNER to postgres;
CREATE TABLE IF NOT EXISTS public.images_h1 PARTITION OF public.images
    FOR VALUES WITH (modulus 10, remainder 1);

ALTER TABLE public.images_h1
    OWNER to postgres;
CREATE TABLE IF NOT EXISTS public.images_h2 PARTITION OF public.images
    FOR VALUES WITH (modulus 10, remainder 2);

ALTER TABLE public.images_h2
    OWNER to postgres;
CREATE TABLE IF NOT EXISTS public.images_h3 PARTITION OF public.images
    FOR VALUES WITH (modulus 10, remainder 3);

ALTER TABLE public.images_h3
    OWNER to postgres;
CREATE TABLE IF NOT EXISTS public.images_h4 PARTITION OF public.images
    FOR VALUES WITH (modulus 10, remainder 4);

ALTER TABLE public.images_h4
    OWNER to postgres;
CREATE TABLE IF NOT EXISTS public.images_h5 PARTITION OF public.images
    FOR VALUES WITH (modulus 10, remainder 5);

ALTER TABLE public.images_h5
    OWNER to postgres;
CREATE TABLE IF NOT EXISTS public.images_h6 PARTITION OF public.images
    FOR VALUES WITH (modulus 10, remainder 6);

ALTER TABLE public.images_h6
    OWNER to postgres;
CREATE TABLE IF NOT EXISTS public.images_h7 PARTITION OF public.images
    FOR VALUES WITH (modulus 10, remainder 7);

ALTER TABLE public.images_h7
    OWNER to postgres;
CREATE TABLE IF NOT EXISTS public.images_h8 PARTITION OF public.images
    FOR VALUES WITH (modulus 10, remainder 8);

ALTER TABLE public.images_h8
    OWNER to postgres;
CREATE TABLE IF NOT EXISTS public.images_h9 PARTITION OF public.images
    FOR VALUES WITH (modulus 10, remainder 9);

ALTER TABLE public.images_h9
    OWNER to postgres;