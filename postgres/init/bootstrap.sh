#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 <<-EOSQL
  -- Runtime user
  CREATE ROLE appuser
    LOGIN
    PASSWORD '${APP_DB_PASSWORD}'
    NOSUPERUSER
    NOCREATEDB
    NOCREATEROLE;

  -- Migration user
  CREATE ROLE migrator
    LOGIN
    PASSWORD '${MIGRATOR_DB_PASSWORD}'
    NOSUPERUSER
    CREATEDB
    CREATEROLE;

  CREATE DATABASE cspotlight OWNER migrator;

  \connect cspotlight

  -- Schema ownership
  ALTER SCHEMA public OWNER TO migrator;

  -- Runtime permissions
  GRANT CONNECT ON DATABASE cspotlight TO appuser;
  GRANT USAGE ON SCHEMA public TO appuser;
  GRANT CREATE ON SCHEMA public TO appuser;

  -- Future tables (created by migrator)
  ALTER DEFAULT PRIVILEGES FOR ROLE migrator
    IN SCHEMA public
    GRANT SELECT, INSERT, UPDATE, DELETE
    ON TABLES TO appuser;
EOSQL
