-- Database generated with pgModeler (PostgreSQL Database Modeler).
-- pgModeler  version: 0.9.1
-- PostgreSQL version: 10.0
-- Project Site: pgmodeler.io
-- Model Author: James Culbertson


-- Database creation must be done outside a multicommand file.
-- These commands were put in this file only as a convenience.
-- -- object: "knowledge-keepers" | type: DATABASE --
-- -- DROP DATABASE IF EXISTS "knowledge-keepers";
-- CREATE DATABASE "knowledge-keepers";
-- -- ddl-end --
-- 

-- object: public.categories | type: TABLE --
-- DROP TABLE IF EXISTS public.categories CASCADE;
CREATE TABLE public.categories(
	category_id bigserial NOT NULL,
	name varchar(50) NOT NULL,
	description varchar(500),
	created_by bigint NOT NULL,
	updated_by bigint,
	date_created timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP(),
	date_updated timestamp,
	CONSTRAINT categories_pk PRIMARY KEY (category_id)

);
-- ddl-end --

-- object: public.users | type: TABLE --
-- DROP TABLE IF EXISTS public.users CASCADE;
CREATE TABLE public.users(
	user_id bigserial NOT NULL,
	email varchar(255) NOT NULL,
	first_name varchar(50),
	last_name varchar(50),
	date_created timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP(),
	date_updated timestamp,
	CONSTRAINT users_pk PRIMARY KEY (user_id)

);
-- ddl-end --

-- object: users_email_unique_idx | type: INDEX --
-- DROP INDEX IF EXISTS public.users_email_unique_idx CASCADE;
CREATE UNIQUE INDEX users_email_unique_idx ON public.users
	USING btree
	(
	  (LOWER(email))
	);
-- ddl-end --

-- object: public.topics | type: TABLE --
-- DROP TABLE IF EXISTS public.topics CASCADE;
CREATE TABLE public.topics(
	topic_id bigserial NOT NULL,
	category_id bigint NOT NULL,
	title varchar(50) NOT NULL,
	description varchar(1000),
	created_by bigint NOT NULL,
	updated_by bigint,
	date_created timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP(),
	date_updated timestamp,
	CONSTRAINT topics_pk PRIMARY KEY (topic_id)

);
-- ddl-end --

-- object: public.tags | type: TABLE --
-- DROP TABLE IF EXISTS public.tags CASCADE;
CREATE TABLE public.tags(
	name varchar(50) NOT NULL,
	tag_id bigserial NOT NULL,
	created_by bigint NOT NULL,
	updated_by bigint,
	date_created timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP(),
	date_updated timestamp,
	CONSTRAINT tags_pk PRIMARY KEY (tag_id)

);
-- ddl-end --

-- object: public.topics_tags | type: TABLE --
-- DROP TABLE IF EXISTS public.topics_tags CASCADE;
CREATE TABLE public.topics_tags(
	topic_id_topics bigint NOT NULL,
	tag_id_tags bigint NOT NULL,
	CONSTRAINT topics_tags_pk PRIMARY KEY (topic_id_topics,tag_id_tags)

);
-- ddl-end --

-- object: topics_fk | type: CONSTRAINT --
-- ALTER TABLE public.topics_tags DROP CONSTRAINT IF EXISTS topics_fk CASCADE;
ALTER TABLE public.topics_tags ADD CONSTRAINT topics_fk FOREIGN KEY (topic_id_topics)
REFERENCES public.topics (topic_id) MATCH FULL
ON DELETE CASCADE ON UPDATE CASCADE;
-- ddl-end --

-- object: tags_fk | type: CONSTRAINT --
-- ALTER TABLE public.topics_tags DROP CONSTRAINT IF EXISTS tags_fk CASCADE;
ALTER TABLE public.topics_tags ADD CONSTRAINT tags_fk FOREIGN KEY (tag_id_tags)
REFERENCES public.tags (tag_id) MATCH FULL
ON DELETE CASCADE ON UPDATE CASCADE;
-- ddl-end --

-- object: public.media_type | type: TYPE --
-- DROP TYPE IF EXISTS public.media_type CASCADE;
CREATE TYPE public.media_type AS
 ENUM ('link','image','video');
-- ddl-end --

-- object: public.media | type: TABLE --
-- DROP TABLE IF EXISTS public.media CASCADE;
CREATE TABLE public.media(
	media_id bigserial NOT NULL,
	type public.media_type NOT NULL,
	title varchar(50) NOT NULL,
	description varchar(500),
	url varchar(255) NOT NULL,
	created_by bigint NOT NULL,
	updated_by bigint,
	date_created timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP(),
	date_updated timestamp,
	topic_id bigint,
	CONSTRAINT media_pk PRIMARY KEY (media_id)

);
-- ddl-end --

-- object: public.notes | type: TABLE --
-- DROP TABLE IF EXISTS public.notes CASCADE;
CREATE TABLE public.notes(
	note_id bigserial NOT NULL,
	title varchar(50) NOT NULL,
	description varchar(500),
	url varchar(255) NOT NULL,
	created_by bigint NOT NULL,
	updated_by bigint,
	date_created timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP(),
	date_updated timestamp,
	topic_id bigint,
	CONSTRAINT notes_pk PRIMARY KEY (note_id)

);
-- ddl-end --

-- object: media_type_idx | type: INDEX --
-- DROP INDEX IF EXISTS public.media_type_idx CASCADE;
CREATE INDEX media_type_idx ON public.media
	USING btree
	(
	  type
	);
-- ddl-end --

-- object: topics_notes_fk | type: CONSTRAINT --
-- ALTER TABLE public.notes DROP CONSTRAINT IF EXISTS topics_notes_fk CASCADE;
ALTER TABLE public.notes ADD CONSTRAINT topics_notes_fk FOREIGN KEY (topic_id)
REFERENCES public.topics (topic_id) MATCH FULL
ON DELETE CASCADE ON UPDATE CASCADE;
-- ddl-end --

-- object: topics_media_fk | type: CONSTRAINT --
-- ALTER TABLE public.media DROP CONSTRAINT IF EXISTS topics_media_fk CASCADE;
ALTER TABLE public.media ADD CONSTRAINT topics_media_fk FOREIGN KEY (topic_id)
REFERENCES public.topics (topic_id) MATCH FULL
ON DELETE CASCADE ON UPDATE CASCADE;
-- ddl-end --

-- object: categories_created_by_fk | type: CONSTRAINT --
-- ALTER TABLE public.categories DROP CONSTRAINT IF EXISTS categories_created_by_fk CASCADE;
ALTER TABLE public.categories ADD CONSTRAINT categories_created_by_fk FOREIGN KEY (created_by)
REFERENCES public.users (user_id) MATCH FULL
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --

-- object: categories_updated_by_fk | type: CONSTRAINT --
-- ALTER TABLE public.categories DROP CONSTRAINT IF EXISTS categories_updated_by_fk CASCADE;
ALTER TABLE public.categories ADD CONSTRAINT categories_updated_by_fk FOREIGN KEY (updated_by)
REFERENCES public.users (user_id) MATCH FULL
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --

-- object: topics_category_fk | type: CONSTRAINT --
-- ALTER TABLE public.topics DROP CONSTRAINT IF EXISTS topics_category_fk CASCADE;
ALTER TABLE public.topics ADD CONSTRAINT topics_category_fk FOREIGN KEY (topic_id)
REFERENCES public.categories (category_id) MATCH FULL
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --

-- object: topics_created_by_fk | type: CONSTRAINT --
-- ALTER TABLE public.topics DROP CONSTRAINT IF EXISTS topics_created_by_fk CASCADE;
ALTER TABLE public.topics ADD CONSTRAINT topics_created_by_fk FOREIGN KEY (created_by)
REFERENCES public.users (user_id) MATCH FULL
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --

-- object: topics_updated_by_fk | type: CONSTRAINT --
-- ALTER TABLE public.topics DROP CONSTRAINT IF EXISTS topics_updated_by_fk CASCADE;
ALTER TABLE public.topics ADD CONSTRAINT topics_updated_by_fk FOREIGN KEY (updated_by)
REFERENCES public.users (user_id) MATCH FULL
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --

-- object: tags_created_by_fk | type: CONSTRAINT --
-- ALTER TABLE public.tags DROP CONSTRAINT IF EXISTS tags_created_by_fk CASCADE;
ALTER TABLE public.tags ADD CONSTRAINT tags_created_by_fk FOREIGN KEY (created_by)
REFERENCES public.users (user_id) MATCH FULL
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --

-- object: tags_updated_by_fk | type: CONSTRAINT --
-- ALTER TABLE public.tags DROP CONSTRAINT IF EXISTS tags_updated_by_fk CASCADE;
ALTER TABLE public.tags ADD CONSTRAINT tags_updated_by_fk FOREIGN KEY (updated_by)
REFERENCES public.users (user_id) MATCH FULL
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --

-- object: media_created_by_fk | type: CONSTRAINT --
-- ALTER TABLE public.media DROP CONSTRAINT IF EXISTS media_created_by_fk CASCADE;
ALTER TABLE public.media ADD CONSTRAINT media_created_by_fk FOREIGN KEY (created_by)
REFERENCES public.users (user_id) MATCH FULL
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --

-- object: media_updated_by_fk | type: CONSTRAINT --
-- ALTER TABLE public.media DROP CONSTRAINT IF EXISTS media_updated_by_fk CASCADE;
ALTER TABLE public.media ADD CONSTRAINT media_updated_by_fk FOREIGN KEY (updated_by)
REFERENCES public.users (user_id) MATCH FULL
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --

-- object: notes_created_by_fk | type: CONSTRAINT --
-- ALTER TABLE public.notes DROP CONSTRAINT IF EXISTS notes_created_by_fk CASCADE;
ALTER TABLE public.notes ADD CONSTRAINT notes_created_by_fk FOREIGN KEY (created_by)
REFERENCES public.users (user_id) MATCH FULL
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --

-- object: notes_updated_by_fk | type: CONSTRAINT --
-- ALTER TABLE public.notes DROP CONSTRAINT IF EXISTS notes_updated_by_fk CASCADE;
ALTER TABLE public.notes ADD CONSTRAINT notes_updated_by_fk FOREIGN KEY (updated_by)
REFERENCES public.users (user_id) MATCH FULL
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --


