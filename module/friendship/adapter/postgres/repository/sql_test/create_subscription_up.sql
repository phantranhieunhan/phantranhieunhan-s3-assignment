INSERT INTO public.users
(id, email, created_at, updated_at)
VALUES($1, $2, 'now()', 'now()'),
($3, $4, 'now()', 'now()')
