package setup

const (
	CreateUsersTable = `CREATE TABLE IF NOT EXISTS public.users
						(
							id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
							username character varying(100) COLLATE pg_catalog."default" NOT NULL,
							hashed_password text COLLATE pg_catalog."default" NOT NULL,
							role character varying(50) COLLATE pg_catalog."default",
							email character varying(350) COLLATE pg_catalog."default",
							CONSTRAINT users_pkey PRIMARY KEY (id),
							CONSTRAINT username UNIQUE (username),
							CONSTRAINT users_email_key UNIQUE (email)
						)
						
						TABLESPACE pg_default;
						
						ALTER TABLE public.users
						OWNER to postgres;`
	CreateFollowersTable = `CREATE TABLE IF NOT EXISTS public.followers
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
								OWNER to postgres;`
	DropUsersTable      = `DROP TABLE users;`
	DropFollowersTable  = `DROP TABLE followers;`
	CreateTestUser      = `INSERT INTO users(username, hashed_password, email, role) VALUES ('TestUser', 'TestPassword', 'test@email.com', 'test');`
	CreateTestFollower  = `INSERT INTO users(username, hashed_password, email, role) VALUES ('TestFollower', 'TestFollowerPassword', 'test_follower@email.com', 'test');`
	CreateTestFollowing = `INSERT INTO users(username, hashed_password, email, role) VALUES ('TestFollowing', 'TestFollowingPassword', 'test_following@email.com', 'test');`
	FollowTestUser      = `INSERT INTO followers(username, following) VALUES (TestFollower, TestUser);`
	FollowTestFollowing = `INSERT INTO followers(username, following) VALUES (TestUser, TestFollowing);`
	UnfollowTestUser    = `DELETE FROM followers WHERE username = 'TestFollower' AND following = 'TestUser';`
)
