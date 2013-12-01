--
-- PostgreSQL database dump
--

SET statement_timeout = 0;
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
-- Name: termine; Type: TABLE; Schema: public; Owner: mero; Tablespace: 
--

CREATE TABLE termine (
    stammtisch boolean,
    vortrag integer,
    override text,
    location text,
    date date
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
    infourl text DEFAULT ''::text NOT NULL
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
-- Name: id; Type: DEFAULT; Schema: public; Owner: mero
--

ALTER TABLE ONLY vortraege ALTER COLUMN id SET DEFAULT nextval('vortraege_id_seq'::regclass);


--
-- Data for Name: termine; Type: TABLE DATA; Schema: public; Owner: mero
--

COPY termine (stammtisch, vortrag, override, location, date) FROM stdin;
t	\N	\N	Mister Wu	2013-12-05
f	2	\N	\N	2013-12-12
f	1	\N	\N	2013-12-19
\.


--
-- Data for Name: vortraege; Type: TABLE DATA; Schema: public; Owner: mero
--

COPY vortraege (id, date, topic, abstract, speaker, infourl) FROM stdin;
1	2013-12-19	Enlarge your penis	REAL Doctors, REAL Science, REAL Results!\n\nDr. MaxMan was created by George Acuilar, M.D, a Board Certified Urologist who has treated over 70,000 patients with erectile problems. He is a member of\nboth the College of Urology and the Society of Urology, and the director of 46 Urologists. He is also the past president of his state society of\nUrologists. \n\nAfter over seven years of research and testing in the area of erectile dysfunction, Dr Acuilar and his team came up with the breakthrough herbal formula\nthat is now known as Dr MaxMan : a 100% natural, powerful male enhancement formula.\n\nNot only do men report AMAZING increases in penis length, width and stamina, but they are also equally delighted by the sheer intensity and concentrated\npower of their orgasms!	Dr. MaxMan	http://qren.kosed.ru
3	2013-01-02	Secure saugt	foo	Merovius	fuuuu you
2	2013-12-12	This stock is about to go through the roof!	This Company could have a Rumble after the Tumble! WE ARE UP AND \nGOING HIGHER.\n\nLong Term Target Price: .21\nToday Price: .02\nCompany: Registered Express Corp\nDate: Monday, November 25th\nSym: R-G T-X\n\nCould have a Massive Run This Week! A Technical Report with a \nBullish Chart Setup!!!	smiley465@waspu.ru	
\.


--
-- Name: vortraege_id_seq; Type: SEQUENCE SET; Schema: public; Owner: mero
--

SELECT pg_catalog.setval('vortraege_id_seq', 3, true);


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

