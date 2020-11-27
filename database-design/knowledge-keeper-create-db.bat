@echo off
setlocal
set PGPASSWORD=postgres
psql -h localhost -U postgres -d postgres -f create-db.sql
set PGPASSWORD=knowledge-keeper
psql -h localhost -U knowledge-keeper -d knowledge-keeper -f knowledge-keeper-ddl.sql
endlocal
