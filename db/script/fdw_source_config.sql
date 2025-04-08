-- use lendbook_analytics
create extension if not exists postgres_fdw;
select * from pg_extension;

create server lendbook_fdw FOREIGN DATA WRAPPER
postgres_fdw options (host '127.0.0.1',port '5432',dbname 'lendbook');

select * from pg_foreign_server;

create user mapping for postgres server lendbook_fdw
options (user 'postgres',password 'postgres');

select * from pg_user_mapping;

select * from pg_user_mappings;

GRANT CONNECT ON DATABASE lendbook TO fdwUser;
grant connect on database lendbook_analytics to fdwUser;

import FOREIGN SCHEMA public from server lendbook_fdw into public;

