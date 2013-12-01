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
f	2	\N	\N	2013-12-12
f	1	\N	\N	2013-12-19
t	\N	\N	Mr. Wu	2013-12-05
\.


--
-- Data for Name: vortraege; Type: TABLE DATA; Schema: public; Owner: mero
--

COPY vortraege (id, date, topic, abstract, speaker, infourl) FROM stdin;
4	\N	Hacking Yourself, Teil 4: Garbage Collection		Jan	
5	\N	Agile (Software-)Entwicklung mit Scrum		Hauro	
6	\N	Clojure		Hauro	
9	2013-11-28	Hacking Yourself, Teil 3: Atmung		Jan	https://www.noname-ev.de/wiki/uploads/8/84/Hacking_Yourself%2C_Teil_3-_Atmung_3.pdf
10	2013-11-14	Suche nach cLFV bei LHCb (oder: was Paul so den ganzen Tag lang macht)		Paul	https://www.noname-ev.de/w/File:Pseyfert_clfv_lhcb_20131114.pdf
11	2013-10-31	Wie exploited man einen buffer overflow auf dem stack? (Mit Anschließenden gemeinsamen Hackmes)		Merovius	https://www.noname-ev.de/w/File:Exploiting.tar.gz
12	2013-10-24	Amateurfunk		pennylane	
13	2013-10-17	Simulation einer Virusinfektion		cherti	
14	2013-10-10	Spezielle Relativitätstheorie		sPhErE	https://www.noname-ev.de/w/File:C14_SRT.pdf
15	2013-09-26	Merovius		Keysigning-Party mit Einführung	https://www.noname-ev.de/w/Keysigning#Keysigning-Party_am_2013-09-26
16	2013-09-19	DBus-Einführung		Merovius	https://www.noname-ev.de/wiki/index.php?title=Chaotische_Viertelstunde&action=submit#2013-09-19
17	2013-09-12	Hacking Yourself, Teil 2: Autosuggestion & Habituierung		Jan	https://www.noname-ev.de/wiki/uploads/6/68/Hacking_Yourself%2C_Teil_2-_Autosuggestion_%26_Habituierung_3.pdf
18	2013-08-29	Linux Sandboxing-Techniken		Merovius	http://youtu.be/wa4AUnUZ6GU
19	2013-08-22	Mailfiltersprache „Sieve“		xeen	
20	2013-08-15	Musikorganisation in 2 Zeichen		xeen	https://www.noname-ev.de/w/File:Musikorga-in-2-zeichen.pdf
21	2013-08-08	Crypto beyond encryption		Merovius	https://www.noname-ev.de/w/File:Crypto.pdf
22	2013-07-25	Recht und Robotik		ink	
23	2013-07-11	Mist, ich brauch dringend so ne Schutzhülle... oder: Warum ist mein Haustürschlüssel Mifare Classic?!?!		Emrys-Merlin	https://www.noname-ev.de/w/File:Mifare_Classic.pdf
24	2013-06-27	Cube-Fun-Facts		Koebi	
25	2013-06-20	git-internals (für Einsteiger und Fortgeschrittene)		Merovius	https://www.noname-ev.de/wiki/index.php?title=Chaotische_Viertelstunde&action=submit#2013-06-20
26	2013-06-13	Hacking Yourself, Teil 1: Ein Neuro-Crashkurs		Jan	https://www.noname-ev.de/wiki/index.php?title=Chaotische_Viertelstunde&action=submit#2013-06-13
27	2013-05-23	Bob der Baumeister und OSM (working title)		xeen	
28	2013-05-16	Life hacks		https://www.noname-ev.de/w/File:Life_hacks.pdf	Fabian
29	2013-05-09	tor and more		eeemsi	
30	2013-04-25	Modellierung biologischer Neuronen		Merovius	https://www.noname-ev.de/w/File:Neuro.pdf
31	2013-04-18	Noten lesen		Koebi	
32	2013-04-11	Einige Erfahrungen mit Abmahnanwälten		Merovius	https://www.noname-ev.de/w/File:Abmahnungen.pdf
33	2013-03-28	Human Enhancements		Alex	
34	2013-03-21	CHICKEN Scheme		sECuRE	https://www.noname-ev.de/w/File:C14h-chicken.pdf
35	2013-03-14	notmuch		sECuRE	https://www.noname-ev.de/w/File:Notmuch.pdf
36	2013-02-28	Building Android From Source		cradle	https://www.noname-ev.de/w/File:aosp_build.pdf
37	2013-02-21	rewriting the ingress intel map		xeen	
38	2013-02-14	Emacs Org-Mode		kungi	
39	2013-01-31	Plan9		sECuRE	https://www.noname-ev.de/w/File:Plan9.pdf
40	2013-01-24	IPv6		sECuRE	https://www.noname-ev.de/w/File:Ipv6-slides.7z.gz
41	2013-01-17	DRBD: Distributed Replicated Block Device		sECuRE	https://www.noname-ev.de/w/File:DRBD.pdf
42	2013-01-10	Profiling und Optimierung für C-Software oder „Writing your own regex-engine for fun, profit and slowness“		Merovius	
43	2012-12-20	Ingress		sECuRE, xeen	
44	2012-12-13	Roboterkunst		Alex	
45	2012-11-29	Kunst am Bau an der Uni Heidelberg		Merovius, inknoir	
46	2012-11-22	Geschichten über WLAN an Bildungseinrichtungen		sECuRE	
47	2012-11-15	Kalenderserver für Dummies		pennylane	
48	2012-11-08	Rufnummermitnahme bei der Telekom		xeen	
49	2012-10-25	the making of PIboard		the_nihilant	
50	2012-10-18	small-scale alerting		sECuRE	http://kanla.zekjur.net/
51	2012-09-27	Podcasts		sECuRE	https://www.noname-ev.de/w/Podcasts
52	2012-09-20	mitgliedsgebühren.pl		sECuRE	https://www.noname-ev.de/w/File:Mitgliedsgeb%C3%BChren.pdf
53	2012-09-13	Debian Code Search		sECuRE	
54	2012-08-30	Yubikey		sECuRE	
55	2012-08-23	buildbot		sECuRE	https://www.noname-ev.de/w/File:C14-Buildbot.pdf
56	2012-08-16	scaling the apache/php/postgresql-stack, a real-world example		sECuRE	https://www.noname-ev.de/w/File:C14h-Webscaling.pdf
57	2012-08-09	Recommender Systems		xeen	https://www.noname-ev.de/w/File:%22folien%22_recommender_systems.pdf
58	2012-07-26	Reverse Engineering, Teil 2		sECuRE	
59	2012-07-19	nVidia CUDA		sECuRE	https://www.noname-ev.de/w/File:Cuda.pdf
60	2012-07-12	Mathematik von Reed-solomon-codes		Merovius	https://www.noname-ev.de/w/File:Rs_codes.pdf
7	\N	TMVA, open source neuronales netz und mehr		Paul	
8	\N	How To Abspann		Koebi	
61	2012-06-28	Crackmes: Reverse-Engineering auf Assembler-Ebene		sECuRE	https://www.noname-ev.de/wiki/index.php?title=Chaotische_Viertelstunde&action=submit#2012-06-28
62	2012-06-21	Launchd - yet another (earlier) way to combine ALL the d's		eeemsi	
63	2012-06-14	Project Dantalion - Eine idee im pre-alpha stadium, die ich vorstellen moechte; diskussion erwuenscht		slowpoke	
64	2012-05-31	Go — eine moderne Programmiersprache		sECuRE	https://www.noname-ev.de/w/File:Gpn12-go.pdf
65	2012-05-24	SCSS		xeen	http://www.sass-lang.com/
66	2012-05-17	Google App Engine (+Go)		sECuRE	https://www.noname-ev.de/w/File:Appengine-go.pdf
67	2012-05-10	Esperanto		the_nihilant	
68	2012-04-26	Raspberry Pi		sECuRE	https://www.noname-ev.de/w/File:Raspberry-pi-slides.pdf
69	2012-04-19	Leaflet		xeen	https://shell.noname-ev.de/~xeen/leaflet-intro/
70	2012-04-12	Debian Packaging		sECuRE	%%%
71	2012-03-29	Linux Networking, Ninja-Style		sECuRE	https://www.noname-ev.de/w/File:Linux-Networking-Ninja-Style.pdf
72	2012-03-22	Protocol Buffers		sECuRE	https://www.noname-ev.de/w/File:Protobuf.pdf
73	2012-03-15	„AllKnowingDNS — Reverse DNS für 2^64 IPv6-Adressen“		sECuRE	https://www.noname-ev.de/w/File:All-knowing-dns.pdf
74	2012-03-08	Live-Coding, oder warum Testsuites toll sind, oder sECuRE macht sich zum Affen.		sECuRE	http://code.stapelberg.de/git/noname-wiki-updater/commit/?id=600bdd7a718ae78a43aa8328285cd347e93cce30
75	2012-02-23	Remote-upgrading microcontroller firmware like a boss		sECuRE	https://www.noname-ev.de/w/File:Firmware-pinpad.pdf
76	2011-08-18	systemd, ein init-replacement		sECuRE	https://www.noname-ev.de/w/File:C14h-Systemd.pdf
77	2011-05-26	Bitcoin, eine P2P Cryptocurrency		Moredread	http://www.bitcoin.org/
78	2011-04-07	X11-Visualisierung		sECuRE	http://www.x11vis.org/
79	2011-03-31	lolpizza – Making Of		sECuRE	https://www.noname-ev.de/w/File:Lolpizza-slides.pdf
80	2010-12-02	Vorbis, an open source and high performance audio codec. (seminar will be in english		the_nihilant	http://en.wikipedia.org/wiki/Vorbis
81	2010-11-11	Reactable		Mo	
82	2010-05-13	Das Usenet aus der Sicht eines Spätgeborenen, oder: Usenet damals und heute		Moredread	http://en.wikipedia.org/wiki/Usenet
83	2010-03-18	CouchDB, eine verteilte, dokumentbasierte Datenbank		sECuRE	https://www.noname-ev.de/w/File:Couchdb.pdf
84	2010-02-18	sup, ein Commandline-Mailclient auf Thread-Basis		sECuRE	http://sup.rubyforge.org/
85	2010-02-11	Der Sony PRS 600 Touch EBook Reader		Kungi	
86	2010-02-04	Objective Caml - eine kurze(?) Einführung		mxf	https://www.noname-ev.de/wiki/index.php?title=Chaotische_Viertelstunde&action=submit#2010-02-04
87	2010-01-14	dn42 - ein dynamisches VPN		sECuRE	https://www.noname-ev.de/w/File:Dn42.pdf
88	2009-12-17	Advanced Encryption Standard (AES)		Moredread	
89	2009-12-10	Mercurial aka hg		Moredread	http://mercurial.selenic.com/
90	2009-12-03	Einführung in go, eine moderne Programmiersprache		sECuRE	https://www.noname-ev.de/wiki/index.php?title=Chaotische_Viertelstunde&action=submit#2009-12-03
91	2009-10-15	Vorstellung Hackerspace Rhein-Neckar		sECuRE	https://www.noname-ev.de/w/File:Hackerspace-rn.pdf
92	2009-09-10	PSP - Geschichte der Custom Firmwares		Atsutane	https://www.noname-ev.de/w/File:Psp-cfws.pdf
93	2009-09-03	DCF77		jiska	
94	2009-08-20	Diskussion über 1 konkretes Projekt nächstes Jahr		TrickSTer	
95	2009-08-20	Warum modern perl toll ist		sECuRE	https://www.noname-ev.de/w/File:Modern_Perl.pdf
96	2009-08-13	Überraschungs-c1/4h		msi	
97	2009-07-30	Sehr chaotische Viertelstunde ueber pekuniaertransfermethoden		ch3ka	
98	2009-05-21	Cloud Computing-Infrastruktur mit Eucalyptus		sECuRE	http://pvs.informatik.uni-heidelberg.de/Teaching/CLCP-09/CLCP_SS2009_Michael_Stapelberg_Eucalyptus.pdf
99	2009-05-21	Nachrichten & Manipulation - Wie sauber arbeiten Redaktionen?		Dfg2	https://www.noname-ev.de/w/C14h:Nachrichten_und_Manipulation
100	2009-05-07	Vorstellung von Anon1984 UG (haftungsbeschränkt), Anbieter von VPN-Tunnels und Betreiber von Tor Exit-Node		Unbekannt	http://www.anon1984.de/
101	2009-04-20	Informationsspeicherung		shl	
102	2009-04-16	IPv6 – wie, wo, warum?		sECuRE	https://www.noname-ev.de/w/Image:Ipv6.pdf
103	2009-04-02	Vim 7 und seine Plugins (ein kleiner Überblick was es so alles gibt)		Kungi	https://www.noname-ev.de/w/C14h:Vim_7
104	2009-03-12	Hacking your own window manager		sECuRE	https://www.noname-ev.de/w/C14h:Hacking_your_own_window_manager
105	2009-02-26	Debuggen mit gdb		sECuRE	https://www.noname-ev.de/w/File:Gdb.pdf
106	2009-01-15	Computerlinguistik		Nicolas	https://www.noname-ev.de/w/C14h:Computerlinguistik
107	2008-10-09	SCADA-Hacking		TabascoEye	
108	2008-08-14	Perl GIF Quine		mxf	
109	2008-07-31	DNS Teil 2/3		shl	
110	2008-07-03	DNS Teil 1/3		shl	
111	2008-06-19	Adventures in Cryptography, induced by p2pdfs. Jetzt neu: die "fast ohne Mathematik"-Edition		ccount	https://www.noname-ev.de/w/C14h:Kryptographie_in_dustfs
112	2008-06-12	git -- the stupid content tracker. Wie arbeite ich mit git und warum ist es so viel cooler als subversion?		sECuRE	https://www.noname-ev.de/w/File:Git-120608.pdf
113	2008-06-05	Status von p2pdfs (peer-to-peer distributed filesystem), unserem bittorrent-gestützten FUSE-Filesystem - (Kurzform: Ja, da kommen Daten durch mittlerweile!)		sECuRE	https://www.noname-ev.de/w/File:P2pdfs-050608.pdf
114	2008-05-08	mxallowd		sECuRE	https://www.noname-ev.de/w/File:Mxallowd-08052008.pdf
115	2008-04-03	GNU-Screen, 50% Vortrag, danach Spielen am Gerät		TrirckSTer	https://www.noname-ev.de/w/C14h:GNU-Screen_von_TrickSTer
116	2008-03-27	Bacula-Vortrag, ebenfalls als 50% Vortrag/50%-Workshop		sECuRE	https://www.noname-ev.de/w/C14h:Bacula-Vortrag
117	2008-03-06	Chaotische 2h über Openstreetmap		sur5r	http://www.openstreetmap.org
118	2008-02-28	Lisp-Interpreter selbst bauen		PhilFry	https://www.noname-ev.de/w/C14h:Lisp_Interpreter_selbst_bauen
119	2008-02-22	Lisp-Vortrag, mit anschließendem Workshop für Interessierte		PhilFry, Kungi	https://www.noname-ev.de/w/C14h:Lisp_Vortrag
120	2007-04-12	Python auf Nokia Series 60 Telefonen eine kurze Einführung		Kungi	
121	2007-02-15	Visualisierung von Programmen		Benedikt	
122	2006-08-24	sECuRE und ch3ka stellen ihr (Python/PHP-)Script PIX vor, mit dem man Filme "galleryzen" kann.		sECuRE, ch3ka	
123	2006-08-17	Drmotte erzählt über die LPI (Linux Professional Institut) Zertifizierung, den Lernstoff, die Prüfungen usw.		Drmotte	
124	2006-08-12	Logikrätsel (Reverse-Engineering)		Benedikt	
125	2006-07-27	Mømø bringt euch die sehr simple und doch maechtige Pluginstruktur von s9y etwas näher.		Mømø	
126	2006-07-20	Alex hat spontan was über Lojban erzählt		Alex	
127	2006-07-06	k1w1 und sur5r erzählen über PICs und was man damit so machen kann.		k1w1, sur5r	
128	2006-07-06	TrickSTer erzählt kurz was zum Thema Elektronische Gesundheitskarte und den  konkreten Planungen.		TrickSTer	
129	2006-06-29	Toshiba schrauben		Craegga	
130	2006-06-22	Infon-Spiel von der GPN		Craegga	
131	2005-12-01	Routing-Protokolle im Internet		shl	https://www.noname-ev.de/w/C14h:Routing_Protokolle_im_Internet
132	2004-06-03	Barrierefreies Webdesign		Alex	https://www.noname-ev.de/w/C14h:Barrierefreies_Webdesign
133	2004-04-22	(La)TeX-Einführung		sur5r, Matthias	
\.


--
-- Name: vortraege_id_seq; Type: SEQUENCE SET; Schema: public; Owner: mero
--

SELECT pg_catalog.setval('vortraege_id_seq', 133, true);


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

