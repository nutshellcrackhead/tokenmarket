--
-- PostgreSQL database dump
--

-- Dumped from database version 9.6.2
-- Dumped by pg_dump version 9.6.2

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

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
-- Name: token_sessions; Type: TABLE; Schema: public; Owner: token
--

CREATE TABLE token_sessions (
    id integer NOT NULL,
    username integer,
    valid_till integer,
    token text
);

ALTER TABLE token_sessions OWNER TO token;

--
-- Name: token_sessions_id_seq; Type: SEQUENCE; Schema: public; Owner: token
--

CREATE SEQUENCE token_sessions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE token_sessions_id_seq OWNER TO token;

--
-- Name: token_sessions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: token
--

ALTER SEQUENCE token_sessions_id_seq OWNED BY token_sessions.id;


--
-- Name: token_users; Type: TABLE; Schema: public; Owner: token
--

CREATE TABLE token_users (
    id integer NOT NULL,
    name text,
    location text,
    skype text,
    phone text,
    email text,
    salt character varying DEFAULT ''::character varying,
    password character varying DEFAULT ''::character varying,
    muted timestamp without time zone,
    avatar text
);


ALTER TABLE token_users OWNER TO token;

--
-- Name: token_users_id_seq; Type: SEQUENCE; Schema: public; Owner: token
--

CREATE SEQUENCE token_users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE token_users_id_seq OWNER TO token;

--
-- Name: token_users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: token
--

ALTER SEQUENCE token_users_id_seq OWNED BY token_users.id;



--
-- Name: token_sessions id; Type: DEFAULT; Schema: public; Owner: token
--

ALTER TABLE ONLY token_sessions ALTER COLUMN id SET DEFAULT nextval('token_sessions_id_seq'::regclass);


--
-- Name: token_users id; Type: DEFAULT; Schema: public; Owner: token
--

ALTER TABLE ONLY token_users ALTER COLUMN id SET DEFAULT nextval('token_users_id_seq'::regclass);

--
-- Name: token_sessions token_sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: token
--

ALTER TABLE ONLY token_sessions
    ADD CONSTRAINT token_sessions_pkey PRIMARY KEY (id);


--
-- Name: token_users token_users_pkey; Type: CONSTRAINT; Schema: public; Owner: token
--

ALTER TABLE ONLY token_users
    ADD CONSTRAINT token_users_pkey PRIMARY KEY (id);


--
-- Name: token_users token_users_skype_key; Type: CONSTRAINT; Schema: public; Owner: token
--

ALTER TABLE ONLY token_users
    ADD CONSTRAINT token_users_skype_key UNIQUE (skype);

--
-- Name: token_sessions_token; Type: INDEX; Schema: public; Owner: token
--

CREATE INDEX token_sessions_token ON token_sessions USING btree (token);


--
-- PostgreSQL database dump complete
--

