create user fdwUser with password 'postgres';
grant usage on schema public to fdwUser;
grant pg_read_all_data TO fdwUser;