DELETE FROM public.subscriptions
WHERE id=$1;

DELETE FROM public.users
WHERE id=$2 OR id=$3;
