-- Table: public.posts

-- DROP TABLE public.posts;

CREATE TABLE IF NOT EXISTS public.posts
(
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    username character varying(100) COLLATE pg_catalog."default" NOT NULL,
    image_id bigint NOT NULL,
    name character varying(300) COLLATE pg_catalog."default" NOT NULL,
    created_on timestamp(6) without time zone NOT NULL DEFAULT now(),
    CONSTRAINT post_pkey PRIMARY KEY (id),
    CONSTRAINT fk_image_id FOREIGN KEY (image_id)
        REFERENCES public.images (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT fk_username FOREIGN KEY (username)
        REFERENCES public.users (username) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT post_image_id_fkey FOREIGN KEY (image_id)
        REFERENCES public.images_h0 (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT post_image_id_fkey1 FOREIGN KEY (image_id)
        REFERENCES public.images_h1 (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT post_image_id_fkey2 FOREIGN KEY (image_id)
        REFERENCES public.images_h2 (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT post_image_id_fkey3 FOREIGN KEY (image_id)
        REFERENCES public.images_h3 (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT post_image_id_fkey4 FOREIGN KEY (image_id)
        REFERENCES public.images_h4 (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT post_image_id_fkey5 FOREIGN KEY (image_id)
        REFERENCES public.images_h5 (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT post_image_id_fkey6 FOREIGN KEY (image_id)
        REFERENCES public.images_h6 (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT post_image_id_fkey7 FOREIGN KEY (image_id)
        REFERENCES public.images_h7 (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT post_image_id_fkey8 FOREIGN KEY (image_id)
        REFERENCES public.images_h8 (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT post_image_id_fkey9 FOREIGN KEY (image_id)
        REFERENCES public.images_h9 (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)

TABLESPACE pg_default;

ALTER TABLE public.posts
    OWNER to postgres;