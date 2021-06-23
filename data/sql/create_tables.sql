--
-- PostgreSQL database dump
--

-- Dumped from database version 13.3
-- Dumped by pg_dump version 13.3

-- Started on 2021-06-18 01:00:51 CDT

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- TOC entry 2 (class 3079 OID 16384)
-- Name: adminpack; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS adminpack WITH SCHEMA pg_catalog;


--
-- TOC entry 3278 (class 0 OID 0)
-- Dependencies: 2
-- Name: EXTENSION adminpack; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION adminpack IS 'administrative functions for PostgreSQL';


--
-- TOC entry 205 (class 1259 OID 16469)
-- Name: event_names_event_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.event_names_event_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.event_names_event_id_seq OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 201 (class 1259 OID 16399)
-- Name: event_names; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.event_names (
    event_name_id smallint DEFAULT nextval('public.event_names_event_id_seq'::regclass) NOT NULL,
    event_name character varying(255) NOT NULL
);


ALTER TABLE public.event_names OWNER TO postgres;

--
-- TOC entry 204 (class 1259 OID 16442)
-- Name: events; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.events (
    event_id bigint NOT NULL,
    event_timestamp timestamp with time zone NOT NULL,
    event_name_id smallint NOT NULL,
    session_id uuid NOT NULL
);


ALTER TABLE public.events OWNER TO postgres;

--
-- TOC entry 203 (class 1259 OID 16440)
-- Name: events_event_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.events_event_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.events_event_id_seq OWNER TO postgres;

--
-- TOC entry 3279 (class 0 OID 0)
-- Dependencies: 203
-- Name: events_event_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.events_event_id_seq OWNED BY public.events.event_id;


--
-- TOC entry 202 (class 1259 OID 16407)
-- Name: sessions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.sessions (
    session_id uuid NOT NULL,
    session_start timestamp with time zone,
    session_end timestamp with time zone
);


ALTER TABLE public.sessions OWNER TO postgres;

--
-- TOC entry 3129 (class 2604 OID 16458)
-- Name: events event_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.events ALTER COLUMN event_id SET DEFAULT nextval('public.events_event_id_seq'::regclass);


--
-- TOC entry 3268 (class 0 OID 16399)
-- Dependencies: 201
-- Data for Name: event_names; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.event_names (event_name_id, event_name) FROM stdin;
1	test
2	event1
3	event2
\.


--
-- TOC entry 3271 (class 0 OID 16442)
-- Dependencies: 204
-- Data for Name: events; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.events (event_id, event_timestamp, event_name_id, session_id) FROM stdin;
3	2021-06-17 23:55:01.531278-05	1	38586dcd-4ce4-486c-8770-6b2f87fed6bf
4	2021-06-17 23:59:03-05	1	38586dcd-4ce4-486c-8770-6b2f87fed6bf
5	2021-06-17 23:59:04-05	1	38586dcd-4ce4-486c-8770-6b2f87fed6bf
6	2021-06-17 23:59:04-05	2	38586dcd-4ce4-486c-8770-6b2f87fed6bf
\.


--
-- TOC entry 3269 (class 0 OID 16407)
-- Dependencies: 202
-- Data for Name: sessions; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.sessions (session_id, session_start, session_end) FROM stdin;
38586dcd-4ce4-486c-8770-6b2f87fed6bf	2021-06-17 23:21:47.525826-05	2021-06-17 23:59:59.525826-05
\.


--
-- TOC entry 3280 (class 0 OID 0)
-- Dependencies: 205
-- Name: event_names_event_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.event_names_event_id_seq', 3, true);


--
-- TOC entry 3281 (class 0 OID 0)
-- Dependencies: 203
-- Name: events_event_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.events_event_id_seq', 6, true);


--
-- TOC entry 3131 (class 2606 OID 16406)
-- Name: event_names event_names_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.event_names
    ADD CONSTRAINT event_names_pkey PRIMARY KEY (event_name_id);


--
-- TOC entry 3135 (class 2606 OID 16460)
-- Name: events events_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.events
    ADD CONSTRAINT events_pkey PRIMARY KEY (event_id);


--
-- TOC entry 3133 (class 2606 OID 16411)
-- Name: sessions sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_pkey PRIMARY KEY (session_id);


--
-- TOC entry 3136 (class 2606 OID 16448)
-- Name: events fk_event_name_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.events
    ADD CONSTRAINT fk_event_name_id FOREIGN KEY (event_name_id) REFERENCES public.event_names(event_name_id);


--
-- TOC entry 206 (class 1255 OID 16472)
-- Name: customer_events_notify(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.customer_events_notify() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
	PERFORM pg_notify('customer_events', row_to_json(NEW)::text);
	RETURN NEW;
END;
$$;


ALTER FUNCTION public.customer_events_notify() OWNER TO postgres;

--
-- TOC entry 3139 (class 2620 OID 16478)
-- Name: events events_status; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER events_status AFTER INSERT ON public.events FOR EACH ROW EXECUTE FUNCTION public.customer_events_notify();


--
-- TOC entry 3137 (class 2606 OID 16453)
-- Name: events fk_session_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.events
    ADD CONSTRAINT fk_session_id FOREIGN KEY (session_id) REFERENCES public.sessions(session_id);


-- Completed on 2021-06-18 01:00:51 CDT

--
-- PostgreSQL database dump complete
--

