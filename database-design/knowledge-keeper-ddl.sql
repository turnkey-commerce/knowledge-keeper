-- Database generated with pgModeler (PostgreSQL Database Modeler).
-- pgModeler  version: 0.9.1
-- PostgreSQL version: 9.6
-- Project Site: pgmodeler.io
-- Model Author: James Culbertson

SET check_function_bodies = false;
-- ddl-end --


-- Database creation must be done outside a multicommand file.
-- These commands were put in this file only as a convenience.
-- -- object: "knowledge-keeper" | type: DATABASE --
-- -- DROP DATABASE IF EXISTS "knowledge-keeper";
-- CREATE DATABASE "knowledge-keeper";
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
	date_created timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	date_updated timestamp,
	CONSTRAINT categories_pk PRIMARY KEY (category_id)

);
-- ddl-end --

-- object: public.users | type: TABLE --
-- DROP TABLE IF EXISTS public.users CASCADE;
CREATE TABLE public.users(
	user_id bigserial NOT NULL,
	email varchar(255) NOT NULL,
	first_name varchar(50) NOT NULL,
	last_name varchar(50) NOT NULL,
	date_created timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
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
	date_created timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
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
	date_created timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	date_updated timestamp,
	CONSTRAINT tags_pk PRIMARY KEY (tag_id)

);
-- ddl-end --

-- object: public.topics_tags | type: TABLE --
-- DROP TABLE IF EXISTS public.topics_tags CASCADE;
CREATE TABLE public.topics_tags(
	topic_id bigint NOT NULL,
	tag_id bigint NOT NULL,
	CONSTRAINT topics_tags_pk PRIMARY KEY (topic_id,tag_id)

);
-- ddl-end --

-- object: topics_fk | type: CONSTRAINT --
-- ALTER TABLE public.topics_tags DROP CONSTRAINT IF EXISTS topics_fk CASCADE;
ALTER TABLE public.topics_tags ADD CONSTRAINT topics_fk FOREIGN KEY (topic_id)
REFERENCES public.topics (topic_id) MATCH FULL
ON DELETE CASCADE ON UPDATE CASCADE;
-- ddl-end --

-- object: tags_fk | type: CONSTRAINT --
-- ALTER TABLE public.topics_tags DROP CONSTRAINT IF EXISTS tags_fk CASCADE;
ALTER TABLE public.topics_tags ADD CONSTRAINT tags_fk FOREIGN KEY (tag_id)
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
	date_created timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
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
	created_by bigint NOT NULL,
	updated_by bigint,
	date_created timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
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

-- object: users_email_idx | type: INDEX --
-- DROP INDEX IF EXISTS public.users_email_idx CASCADE;
CREATE INDEX users_email_idx ON public.users
	USING btree
	(
	  email
	);
-- ddl-end --

-- object: categories_name_idx | type: INDEX --
-- DROP INDEX IF EXISTS public.categories_name_idx CASCADE;
CREATE INDEX categories_name_idx ON public.categories
	USING btree
	(
	  name
	);
-- ddl-end --

-- object: topics_title_idx | type: INDEX --
-- DROP INDEX IF EXISTS public.topics_title_idx CASCADE;
CREATE INDEX topics_title_idx ON public.topics
	USING btree
	(
	  title
	);
-- ddl-end --

-- object: tags_name_idx | type: INDEX --
-- DROP INDEX IF EXISTS public.tags_name_idx CASCADE;
CREATE INDEX tags_name_idx ON public.tags
	USING btree
	(
	  name
	);
-- ddl-end --

-- object: media_title_idx | type: INDEX --
-- DROP INDEX IF EXISTS public.media_title_idx CASCADE;
CREATE INDEX media_title_idx ON public.media
	USING btree
	(
	  title
	);
-- ddl-end --

-- object: notes_title_idx | type: INDEX --
-- DROP INDEX IF EXISTS public.notes_title_idx CASCADE;
CREATE INDEX notes_title_idx ON public.notes
	USING btree
	(
	  title
	);
-- ddl-end --

-- object: public.update_timestamp_column | type: FUNCTION --
-- DROP FUNCTION IF EXISTS public.update_timestamp_column() CASCADE;
CREATE FUNCTION public.update_timestamp_column ()
	RETURNS trigger
	LANGUAGE plpgsql
	VOLATILE 
	CALLED ON NULL INPUT
	SECURITY INVOKER
	COST 1
	AS $$
BEGIN
   NEW.date_updated = now(); 
   RETURN NEW;
END;
$$;
-- ddl-end --

-- object: update_date_updated_trigger | type: TRIGGER --
-- DROP TRIGGER IF EXISTS update_date_updated_trigger ON public.users CASCADE;
CREATE TRIGGER update_date_updated_trigger
	BEFORE UPDATE
	ON public.users
	FOR EACH ROW
	EXECUTE PROCEDURE public.update_timestamp_column();
-- ddl-end --

-- object: update_date_updated_trigger | type: TRIGGER --
-- DROP TRIGGER IF EXISTS update_date_updated_trigger ON public.categories CASCADE;
CREATE TRIGGER update_date_updated_trigger
	BEFORE UPDATE
	ON public.categories
	FOR EACH ROW
	EXECUTE PROCEDURE public.update_timestamp_column();
-- ddl-end --

-- object: update_date_updated_trigger | type: TRIGGER --
-- DROP TRIGGER IF EXISTS update_date_updated_trigger ON public.topics CASCADE;
CREATE TRIGGER update_date_updated_trigger
	BEFORE UPDATE
	ON public.topics
	FOR EACH ROW
	EXECUTE PROCEDURE public.update_timestamp_column();
-- ddl-end --

-- object: update_date_updated_trigger | type: TRIGGER --
-- DROP TRIGGER IF EXISTS update_date_updated_trigger ON public.tags CASCADE;
CREATE TRIGGER update_date_updated_trigger
	BEFORE UPDATE
	ON public.tags
	FOR EACH ROW
	EXECUTE PROCEDURE public.update_timestamp_column();
-- ddl-end --

-- object: update_date_updated_trigger | type: TRIGGER --
-- DROP TRIGGER IF EXISTS update_date_updated_trigger ON public.media CASCADE;
CREATE TRIGGER update_date_updated_trigger
	BEFORE UPDATE
	ON public.media
	FOR EACH ROW
	EXECUTE PROCEDURE public.update_timestamp_column();
-- ddl-end --

-- object: update_date_updated_trigger | type: TRIGGER --
-- DROP TRIGGER IF EXISTS update_date_updated_trigger ON public.notes CASCADE;
CREATE TRIGGER update_date_updated_trigger
	BEFORE UPDATE
	ON public.notes
	FOR EACH ROW
	EXECUTE PROCEDURE public.update_timestamp_column();
-- ddl-end --

-- object: public.tag_topics_view | type: VIEW --
-- DROP VIEW IF EXISTS public.tag_topics_view CASCADE;
CREATE VIEW public.tag_topics_view
AS 

SELECT tag_id, topics.*
    FROM topics_tags LEFT JOIN topics
      ON topics_tags.topic_id = topics.topic_id;
-- ddl-end --

-- object: public.topics_tags_view | type: VIEW --
-- DROP VIEW IF EXISTS public.topics_tags_view CASCADE;
CREATE VIEW public.topics_tags_view
AS 

SELECT topic_id, tags.*
    FROM topics_tags LEFT JOIN tags
      ON topics_tags.tag_id = tags.tag_id;
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
ALTER TABLE public.topics ADD CONSTRAINT topics_category_fk FOREIGN KEY (category_id)
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


