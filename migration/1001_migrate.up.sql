CREATE SCHEMA IF NOT EXISTS public;

CREATE TABLE public.users(
    id text NOT NULL,
    username text NOT NULL,
    password text NOT NULL,
    email text not null, 
    created_at timestamp with time zone NOT null,
    updated_at timestamp with time zone NOT null,
    CONSTRAINT users_pk PRIMARY KEY (id)
);

CREATE TABLE public.friendships(
	id text not null,
	user_id text not null,
	friend_id text not null,
	status int not null default 0,
    created_at timestamp with time zone NOT null,
    updated_at timestamp with time zone NOT null,
    CONSTRAINT friendships_pk PRIMARY KEY (id)
);

CREATE TABLE public.followers(
	id text not null,
	user_id text not null,
	following_id text not null,
	status int not null default 0,
	created_at timestamp with time zone not null,
    updated_at timestamp with time zone not null,
    CONSTRAINT followers_pk PRIMARY KEY (id)
);

INSERT INTO public.users
(id, username, "password", email, created_at, updated_at)
VALUES('cd2543cd-6566-4661-a122-2c963fc16b7c', 'andy', 'encrypted-password', 'andy@example', 'now()', 'now()'),
('b44ca9eb-5d0f-41be-9ecd-dd0158e72e2c', 'john', 'encrypted-password', 'john@example', 'now()', 'now()'),
('afed6e29-07d1-443a-a0c7-38d77ef8f332', 'lisa', 'encrypted-password', 'lisa@example', 'now()', 'now()'),
('6bf98bcf-dd9a-4fd8-b43b-b96ea5f5fe7f', 'kate', 'encrypted-password', 'kate@example', 'now()', 'now()'),
('a46cef8e-ef3d-46e0-9f06-a7bb0d32b310', 'common', 'encrypted-password', 'common@example', 'now()', 'now()')
