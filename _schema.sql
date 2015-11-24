--
-- PostgreSQL database dump
--

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;

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
-- Name: termine; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE termine (
    stammtisch boolean,
    vortrag integer,
    override text NOT NULL,
    location text DEFAULT ''::text NOT NULL,
    date date,
    override_long text DEFAULT ''::text NOT NULL
);


ALTER TABLE termine OWNER TO postgres;

--
-- Name: vortraege; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE vortraege (
    id integer NOT NULL,
    date date,
    topic text DEFAULT ''::text NOT NULL,
    abstract text DEFAULT ''::text NOT NULL,
    speaker text DEFAULT ''::text NOT NULL,
    infourl text DEFAULT ''::text NOT NULL,
    password text DEFAULT ''::text NOT NULL
);


ALTER TABLE vortraege OWNER TO postgres;

--
-- Name: vortraege_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE vortraege_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE vortraege_id_seq OWNER TO postgres;

--
-- Name: vortraege_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE vortraege_id_seq OWNED BY vortraege.id;


--
-- Name: vortrag_links; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE vortrag_links (
    id integer NOT NULL,
    vortrag integer,
    kind text,
    url text
);


ALTER TABLE vortrag_links OWNER TO postgres;

--
-- Name: vortrag_links_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE vortrag_links_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE vortrag_links_id_seq OWNER TO postgres;

--
-- Name: vortrag_links_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE vortrag_links_id_seq OWNED BY vortrag_links.id;


--
-- Name: zusagen; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE zusagen (
    nick text,
    kommt boolean DEFAULT false NOT NULL,
    kommentar text DEFAULT ''::text NOT NULL,
    registered timestamp with time zone
);


ALTER TABLE zusagen OWNER TO postgres;

--
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY vortraege ALTER COLUMN id SET DEFAULT nextval('vortraege_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY vortrag_links ALTER COLUMN id SET DEFAULT nextval('vortrag_links_id_seq'::regclass);


--
-- Name: unique_date; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY termine
    ADD CONSTRAINT unique_date UNIQUE (date);


--
-- Name: vortraege_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY vortraege
    ADD CONSTRAINT vortraege_pkey PRIMARY KEY (id);


--
-- Name: vortrag_links_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY vortrag_links
    ADD CONSTRAINT vortrag_links_pkey PRIMARY KEY (id);


--
-- Name: vortrag_links_vortrag_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY vortrag_links
    ADD CONSTRAINT vortrag_links_vortrag_fkey FOREIGN KEY (vortrag) REFERENCES vortraege(id);


--
-- Name: public; Type: ACL; Schema: -; Owner: postgres
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM postgres;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- Name: termine; Type: ACL; Schema: public; Owner: postgres
--

REVOKE ALL ON TABLE termine FROM PUBLIC;
REVOKE ALL ON TABLE termine FROM postgres;
GRANT ALL ON TABLE termine TO postgres;
GRANT SELECT ON TABLE termine TO PUBLIC;
GRANT UPDATE ON TABLE termine TO nnweb;
GRANT ALL ON TABLE termine TO anon;


--
-- Name: vortraege; Type: ACL; Schema: public; Owner: postgres
--

REVOKE ALL ON TABLE vortraege FROM PUBLIC;
REVOKE ALL ON TABLE vortraege FROM postgres;
GRANT ALL ON TABLE vortraege TO postgres;
GRANT SELECT ON TABLE vortraege TO PUBLIC;
GRANT INSERT,DELETE,UPDATE ON TABLE vortraege TO nnweb;
GRANT ALL ON TABLE vortraege TO anon;


--
-- Name: vortraege_id_seq; Type: ACL; Schema: public; Owner: postgres
--

REVOKE ALL ON SEQUENCE vortraege_id_seq FROM PUBLIC;
REVOKE ALL ON SEQUENCE vortraege_id_seq FROM postgres;
GRANT ALL ON SEQUENCE vortraege_id_seq TO postgres;
GRANT ALL ON SEQUENCE vortraege_id_seq TO nnweb;
GRANT ALL ON SEQUENCE vortraege_id_seq TO anon;


--
-- Name: zusagen; Type: ACL; Schema: public; Owner: postgres
--

REVOKE ALL ON TABLE zusagen FROM PUBLIC;
REVOKE ALL ON TABLE zusagen FROM postgres;
GRANT ALL ON TABLE zusagen TO postgres;
GRANT SELECT ON TABLE zusagen TO PUBLIC;
GRANT INSERT,DELETE,UPDATE ON TABLE zusagen TO nnweb;
GRANT ALL ON TABLE zusagen TO anon;


--
-- PostgreSQL database dump complete
--
