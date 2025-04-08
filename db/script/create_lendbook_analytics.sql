-- Database: lendbook_analytics

-- DROP DATABASE IF EXISTS lendbook_analytics;

CREATE DATABASE lendbook_analytics
    WITH
    OWNER = postgres
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.UTF-8'
    LC_CTYPE = 'en_US.UTF-8'
    LOCALE_PROVIDER = 'libc'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1
    IS_TEMPLATE = False;

GRANT TEMPORARY, CONNECT ON DATABASE lendbook_analytics TO PUBLIC;

GRANT CONNECT ON DATABASE lendbook_analytics TO fdwuser;

GRANT ALL ON DATABASE lendbook_analytics TO postgres;