--
-- PostgreSQL database dump
--

SET statement_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;

SET search_path = public, pg_catalog;

ALTER TABLE ONLY public.vortraege DROP CONSTRAINT vortraege_pkey;
ALTER TABLE ONLY public.zusagen DROP CONSTRAINT unique_nick;
ALTER TABLE ONLY public.termine DROP CONSTRAINT unique_date;
ALTER TABLE public.vortraege ALTER COLUMN id DROP DEFAULT;
DROP TABLE public.zusagen;
DROP SEQUENCE public.vortraege_id_seq;
DROP TABLE public.vortraege;
DROP TABLE public.termine;
DROP EXTENSION plpgsql;
DROP SCHEMA public;
--
-- Name: public; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA public;


ALTER SCHEMA public OWNER TO postgres;

--
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: postgres
--

COMMENT ON SCHEMA public IS 'standard public schema';


--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: termine; Type: TABLE; Schema: public; Owner: mero; Tablespace: 
--

CREATE TABLE termine (
    stammtisch boolean,
    vortrag integer,
    override text NOT NULL,
    location text DEFAULT ''::text NOT NULL,
    date date,
    override_long text DEFAULT ''::text NOT NULL
);


ALTER TABLE public.termine OWNER TO mero;

--
-- Name: vortraege; Type: TABLE; Schema: public; Owner: mero; Tablespace: 
--

CREATE TABLE vortraege (
    id integer NOT NULL,
    date date,
    topic text DEFAULT ''::text NOT NULL,
    abstract text DEFAULT ''::text NOT NULL,
    speaker text DEFAULT ''::text NOT NULL,
    infourl text DEFAULT ''::text NOT NULL,
    password text DEFAULT '':text NOT NULL
);


ALTER TABLE public.vortraege OWNER TO mero;

--
-- Name: vortraege_id_seq; Type: SEQUENCE; Schema: public; Owner: mero
--

CREATE SEQUENCE vortraege_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.vortraege_id_seq OWNER TO mero;

--
-- Name: vortraege_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mero
--

ALTER SEQUENCE vortraege_id_seq OWNED BY vortraege.id;


--
-- Name: zusagen; Type: TABLE; Schema: public; Owner: mero; Tablespace: 
--

CREATE TABLE zusagen (
    nick text,
    kommt boolean DEFAULT false NOT NULL,
    kommentar text DEFAULT ''::text NOT NULL
);


ALTER TABLE public.zusagen OWNER TO mero;

--
-- Name: id; Type: DEFAULT; Schema: public; Owner: mero
--

ALTER TABLE ONLY vortraege ALTER COLUMN id SET DEFAULT nextval('vortraege_id_seq'::regclass);


--
-- Name: unique_date; Type: CONSTRAINT; Schema: public; Owner: mero; Tablespace: 
--

ALTER TABLE ONLY termine
    ADD CONSTRAINT unique_date UNIQUE (date);


--
-- Name: unique_nick; Type: CONSTRAINT; Schema: public; Owner: mero; Tablespace: 
--

ALTER TABLE ONLY zusagen
    ADD CONSTRAINT unique_nick UNIQUE (nick);


--
-- Name: vortraege_pkey; Type: CONSTRAINT; Schema: public; Owner: mero; Tablespace: 
--

ALTER TABLE ONLY vortraege
    ADD CONSTRAINT vortraege_pkey PRIMARY KEY (id);


--
-- Name: public; Type: ACL; Schema: -; Owner: postgres
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM postgres;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- PostgreSQL database dump complete
--

