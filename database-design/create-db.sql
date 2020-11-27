DROP DATABASE "knowledge-keeper";
DROP USER "knowledge-keeper";
CREATE DATABASE "knowledge-keeper";
CREATE USER "knowledge-keeper" WITH ENCRYPTED PASSWORD 'knowledge-keeper';
GRANT ALL PRIVILEGES ON DATABASE "knowledge-keeper" TO "knowledge-keeper";