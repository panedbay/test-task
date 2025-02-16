--
-- PostgreSQL database dump
--

-- Dumped from database version 16.3
-- Dumped by pg_dump version 16.3

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
-- Name: merch; Type: SCHEMA; Schema: -; Owner: arch
--

CREATE SCHEMA merch;


ALTER SCHEMA merch OWNER TO arch;

--
-- Name: pgcrypto; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;


--
-- Name: EXTENSION pgcrypto; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION pgcrypto IS 'cryptographic functions';


--
-- Name: f_add_employee(text, text); Type: FUNCTION; Schema: merch; Owner: arch
--

CREATE FUNCTION merch.f_add_employee(p_username text, p_password text) RETURNS void
    LANGUAGE plpgsql
    AS $$
DECLARE
    hashed_password TEXT;
BEGIN
    hashed_password := crypt(p_password, gen_salt('bf'));

    INSERT INTO merch.t_employees (username, password) VALUES (p_username, hashed_password);
END;
$$;


ALTER FUNCTION merch.f_add_employee(p_username text, p_password text) OWNER TO arch;

--
-- Name: f_buy(text, text); Type: FUNCTION; Schema: merch; Owner: arch
--

CREATE FUNCTION merch.f_buy(p_item_name text, p_buyer_name text) RETURNS void
    LANGUAGE plpgsql
    AS $$

DECLARE
    v_item_price INT;
    v_buyer_coins INT;
    v_new_quantity INT;
BEGIN
    SELECT price INTO v_item_price FROM merch.t_items WHERE name = p_item_name;
    IF NOT FOUND THEN
        RAISE 'Error: Item not found.';
    END IF;

    SELECT coins INTO v_buyer_coins FROM merch.t_employees WHERE username = p_buyer_name;
    IF NOT FOUND THEN
        RAISE 'Error: Buyer not found.';
    END IF;

    IF v_buyer_coins < v_item_price THEN
        RAISE 'Error: Not enough coins.';
    END IF;
	
	UPDATE merch.t_employees 
    SET coins = coins - v_item_price 
    WHERE username = p_buyer_name;

    IF EXISTS (SELECT 1 FROM merch.t_buys WHERE item_name = p_item_name AND emp_name = p_buyer_name) THEN
        UPDATE merch.t_buys 
        SET quantity = quantity + 1 
        WHERE item_name = p_item_name AND emp_name = p_buyer_name;
    ELSE
        INSERT INTO merch.t_buys (item_name, emp_name, quantity) 
        VALUES (p_item_name, p_buyer_name, 1);
    END IF;
	RETURN;
	
END;
$$;


ALTER FUNCTION merch.f_buy(p_item_name text, p_buyer_name text) OWNER TO arch;

--
-- Name: f_check_employee_credentials(text, text); Type: FUNCTION; Schema: merch; Owner: arch
--

CREATE FUNCTION merch.f_check_employee_credentials(p_username text, p_password text) RETURNS boolean
    LANGUAGE plpgsql
    AS $$
DECLARE
    stored_password TEXT;
BEGIN
    SELECT e.password INTO stored_password
    FROM merch.t_employees e
    WHERE username = p_username;

    IF stored_password IS NULL THEN
        RETURN FALSE;
    END IF;

    IF crypt(p_password, stored_password) = stored_password THEN
        RETURN TRUE;
    ELSE
        RETURN FALSE;
    END IF;
END;
$$;


ALTER FUNCTION merch.f_check_employee_credentials(p_username text, p_password text) OWNER TO arch;

--
-- Name: f_employee_exists(text); Type: FUNCTION; Schema: merch; Owner: arch
--

CREATE FUNCTION merch.f_employee_exists(p_emp_username text) RETURNS boolean
    LANGUAGE plpgsql
    AS $$
DECLARE
    exists_flag BOOLEAN;
BEGIN
    SELECT EXISTS (
        SELECT 1 FROM merch.t_employees WHERE username = p_emp_username
    ) INTO exists_flag;

    RETURN exists_flag;
END;
$$;


ALTER FUNCTION merch.f_employee_exists(p_emp_username text) OWNER TO arch;

--
-- Name: f_get_employee_inventory(text); Type: FUNCTION; Schema: merch; Owner: arch
--

CREATE FUNCTION merch.f_get_employee_inventory(p_username text) RETURNS TABLE(item_name character varying, quantity bigint)
    LANGUAGE plpgsql
    AS $$
BEGIN
    RETURN QUERY
    SELECT b.item_name, b.quantity
    FROM merch.t_buys b
    WHERE b.emp_name = p_username;
END;
$$;


ALTER FUNCTION merch.f_get_employee_inventory(p_username text) OWNER TO arch;

--
-- Name: f_get_transfers_receiver(text); Type: FUNCTION; Schema: merch; Owner: arch
--

CREATE FUNCTION merch.f_get_transfers_receiver(p_username text) RETURNS TABLE(fromuser text, amount bigint)
    LANGUAGE plpgsql
    AS $$
BEGIN
    RETURN QUERY
    SELECT th.sender_name, th.amount
    FROM merch.t_transfer_history th
    WHERE th.receiver_name = p_username;
END;
$$;


ALTER FUNCTION merch.f_get_transfers_receiver(p_username text) OWNER TO arch;

--
-- Name: f_get_transfers_sender(text); Type: FUNCTION; Schema: merch; Owner: arch
--

CREATE FUNCTION merch.f_get_transfers_sender(p_username text) RETURNS TABLE(touser text, amount bigint)
    LANGUAGE plpgsql
    AS $$
BEGIN
    RETURN QUERY
    SELECT th.receiver_name, th.amount
    FROM merch.t_transfer_history th
    WHERE th.sender_name = p_username;
END;
$$;


ALTER FUNCTION merch.f_get_transfers_sender(p_username text) OWNER TO arch;

--
-- Name: f_get_user_coins(text); Type: FUNCTION; Schema: merch; Owner: arch
--

CREATE FUNCTION merch.f_get_user_coins(p_username text) RETURNS integer
    LANGUAGE plpgsql
    AS $$
DECLARE
    user_coins INT;
BEGIN
    SELECT coins INTO user_coins
    FROM merch.t_employees
    WHERE username = p_username;

    IF user_coins IS NULL THEN
        RAISE EXCEPTION 'User does not exist';
    END IF;

    RETURN user_coins;
END;
$$;


ALTER FUNCTION merch.f_get_user_coins(p_username text) OWNER TO arch;

--
-- Name: f_transfer_coins(text, text, bigint); Type: FUNCTION; Schema: merch; Owner: arch
--

CREATE FUNCTION merch.f_transfer_coins(sender_username text, receiver_username text, amount bigint) RETURNS void
    LANGUAGE plpgsql
    AS $$
DECLARE
    sender_balance BIGINT;
BEGIN
    IF amount <= 0 THEN
        RAISE EXCEPTION 'Transfer amount must be greater than zero';
    END IF;

    SELECT coins INTO sender_balance FROM merch.t_employees WHERE username = sender_username;

	IF sender_username = receiver_username THEN
		RAISE EXCEPTION 'Cannot send yorself';
	END IF;

    IF sender_username IS NULL THEN
        RAISE EXCEPTION 'Sender does not exist';
    END IF;

    IF receiver_username IS NULL THEN
        RAISE EXCEPTION 'Receiver does not exist';
    END IF;

    IF sender_balance < amount THEN
        RAISE EXCEPTION 'Insufficient balance';
    END IF;

    UPDATE merch.t_employees SET coins = coins - amount WHERE username = sender_username;
    UPDATE merch.t_employees SET coins = coins + amount WHERE username = receiver_username;

    INSERT INTO merch.t_transfer_history (sender_name, receiver_name, amount) 
    VALUES (sender_username, receiver_username, amount);

END;
$$;


ALTER FUNCTION merch.f_transfer_coins(sender_username text, receiver_username text, amount bigint) OWNER TO arch;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: t_employees; Type: TABLE; Schema: merch; Owner: arch
--

CREATE TABLE merch.t_employees (
    id bigint NOT NULL,
    username character varying NOT NULL,
    password character varying NOT NULL,
    coins bigint DEFAULT 1000 NOT NULL
);


ALTER TABLE merch.t_employees OWNER TO arch;

--
-- Name: TABLE t_employees; Type: COMMENT; Schema: merch; Owner: arch
--

COMMENT ON TABLE merch.t_employees IS 'Information about employees';


--
-- Name: employees_id_seq; Type: SEQUENCE; Schema: merch; Owner: arch
--

CREATE SEQUENCE merch.employees_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE merch.employees_id_seq OWNER TO arch;

--
-- Name: employees_id_seq; Type: SEQUENCE OWNED BY; Schema: merch; Owner: arch
--

ALTER SEQUENCE merch.employees_id_seq OWNED BY merch.t_employees.id;


--
-- Name: t_buys; Type: TABLE; Schema: merch; Owner: arch
--

CREATE TABLE merch.t_buys (
    id bigint NOT NULL,
    item_name character varying NOT NULL,
    emp_name character varying NOT NULL,
    quantity bigint NOT NULL
);


ALTER TABLE merch.t_buys OWNER TO arch;

--
-- Name: TABLE t_buys; Type: COMMENT; Schema: merch; Owner: arch
--

COMMENT ON TABLE merch.t_buys IS 'Table of employee merch buys';


--
-- Name: t_buys_id_seq; Type: SEQUENCE; Schema: merch; Owner: arch
--

CREATE SEQUENCE merch.t_buys_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE merch.t_buys_id_seq OWNER TO arch;

--
-- Name: t_buys_id_seq; Type: SEQUENCE OWNED BY; Schema: merch; Owner: arch
--

ALTER SEQUENCE merch.t_buys_id_seq OWNED BY merch.t_buys.id;


--
-- Name: t_items; Type: TABLE; Schema: merch; Owner: arch
--

CREATE TABLE merch.t_items (
    id bigint NOT NULL,
    name character varying NOT NULL,
    price bigint NOT NULL
);


ALTER TABLE merch.t_items OWNER TO arch;

--
-- Name: TABLE t_items; Type: COMMENT; Schema: merch; Owner: arch
--

COMMENT ON TABLE merch.t_items IS 'Collection of items';


--
-- Name: t_items_id_seq; Type: SEQUENCE; Schema: merch; Owner: arch
--

CREATE SEQUENCE merch.t_items_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE merch.t_items_id_seq OWNER TO arch;

--
-- Name: t_items_id_seq; Type: SEQUENCE OWNED BY; Schema: merch; Owner: arch
--

ALTER SEQUENCE merch.t_items_id_seq OWNED BY merch.t_items.id;


--
-- Name: t_transfer_history; Type: TABLE; Schema: merch; Owner: arch
--

CREATE TABLE merch.t_transfer_history (
    id bigint NOT NULL,
    sender_name text NOT NULL,
    receiver_name text NOT NULL,
    amount bigint NOT NULL
);


ALTER TABLE merch.t_transfer_history OWNER TO arch;

--
-- Name: TABLE t_transfer_history; Type: COMMENT; Schema: merch; Owner: arch
--

COMMENT ON TABLE merch.t_transfer_history IS 'History of money transfers';


--
-- Name: t_transfer_history_id_seq; Type: SEQUENCE; Schema: merch; Owner: arch
--

CREATE SEQUENCE merch.t_transfer_history_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE merch.t_transfer_history_id_seq OWNER TO arch;

--
-- Name: t_transfer_history_id_seq; Type: SEQUENCE OWNED BY; Schema: merch; Owner: arch
--

ALTER SEQUENCE merch.t_transfer_history_id_seq OWNED BY merch.t_transfer_history.id;


--
-- Name: t_buys id; Type: DEFAULT; Schema: merch; Owner: arch
--

ALTER TABLE ONLY merch.t_buys ALTER COLUMN id SET DEFAULT nextval('merch.t_buys_id_seq'::regclass);


--
-- Name: t_employees id; Type: DEFAULT; Schema: merch; Owner: arch
--

ALTER TABLE ONLY merch.t_employees ALTER COLUMN id SET DEFAULT nextval('merch.employees_id_seq'::regclass);


--
-- Name: t_items id; Type: DEFAULT; Schema: merch; Owner: arch
--

ALTER TABLE ONLY merch.t_items ALTER COLUMN id SET DEFAULT nextval('merch.t_items_id_seq'::regclass);


--
-- Name: t_transfer_history id; Type: DEFAULT; Schema: merch; Owner: arch
--

ALTER TABLE ONLY merch.t_transfer_history ALTER COLUMN id SET DEFAULT nextval('merch.t_transfer_history_id_seq'::regclass);


--
-- Data for Name: t_buys; Type: TABLE DATA; Schema: merch; Owner: arch
--

COPY merch.t_buys (id, item_name, emp_name, quantity) FROM stdin;
\.


--
-- Data for Name: t_employees; Type: TABLE DATA; Schema: merch; Owner: arch
--

COPY merch.t_employees (id, username, password, coins) FROM stdin;
\.


--
-- Data for Name: t_items; Type: TABLE DATA; Schema: merch; Owner: arch
--

COPY merch.t_items (id, name, price) FROM stdin;
1	t-shirt	80
2	cup	20
3	book	50
4	pen	10
5	powerbank	200
6	hoody	300
7	umbrella	200
8	socks	10
9	wallet	50
10	pink-hoody	500
\.


--
-- Data for Name: t_transfer_history; Type: TABLE DATA; Schema: merch; Owner: arch
--

COPY merch.t_transfer_history (id, sender_name, receiver_name, amount) FROM stdin;
\.


--
-- Name: employees_id_seq; Type: SEQUENCE SET; Schema: merch; Owner: arch
--

SELECT pg_catalog.setval('merch.employees_id_seq', 5, true);


--
-- Name: t_buys_id_seq; Type: SEQUENCE SET; Schema: merch; Owner: arch
--

SELECT pg_catalog.setval('merch.t_buys_id_seq', 3, true);


--
-- Name: t_items_id_seq; Type: SEQUENCE SET; Schema: merch; Owner: arch
--

SELECT pg_catalog.setval('merch.t_items_id_seq', 10, true);


--
-- Name: t_transfer_history_id_seq; Type: SEQUENCE SET; Schema: merch; Owner: arch
--

SELECT pg_catalog.setval('merch.t_transfer_history_id_seq', 3, true);


--
-- Name: t_employees employees_pkey; Type: CONSTRAINT; Schema: merch; Owner: arch
--

ALTER TABLE ONLY merch.t_employees
    ADD CONSTRAINT employees_pkey PRIMARY KEY (id);


--
-- Name: t_buys t_buys_pkey; Type: CONSTRAINT; Schema: merch; Owner: arch
--

ALTER TABLE ONLY merch.t_buys
    ADD CONSTRAINT t_buys_pkey PRIMARY KEY (id);


--
-- Name: t_employees t_employees_username_key; Type: CONSTRAINT; Schema: merch; Owner: arch
--

ALTER TABLE ONLY merch.t_employees
    ADD CONSTRAINT t_employees_username_key UNIQUE (username);


--
-- Name: t_items t_items_name_key; Type: CONSTRAINT; Schema: merch; Owner: arch
--

ALTER TABLE ONLY merch.t_items
    ADD CONSTRAINT t_items_name_key UNIQUE (name);


--
-- Name: t_items t_items_pkey; Type: CONSTRAINT; Schema: merch; Owner: arch
--

ALTER TABLE ONLY merch.t_items
    ADD CONSTRAINT t_items_pkey PRIMARY KEY (id);


--
-- Name: t_transfer_history t_transfer_history_pkey; Type: CONSTRAINT; Schema: merch; Owner: arch
--

ALTER TABLE ONLY merch.t_transfer_history
    ADD CONSTRAINT t_transfer_history_pkey PRIMARY KEY (id);


--
-- Name: t_buys t_buys_emp_name_fkey; Type: FK CONSTRAINT; Schema: merch; Owner: arch
--

ALTER TABLE ONLY merch.t_buys
    ADD CONSTRAINT t_buys_emp_name_fkey FOREIGN KEY (emp_name) REFERENCES merch.t_employees(username) NOT VALID;


--
-- Name: t_buys t_buys_item_name_fkey; Type: FK CONSTRAINT; Schema: merch; Owner: arch
--

ALTER TABLE ONLY merch.t_buys
    ADD CONSTRAINT t_buys_item_name_fkey FOREIGN KEY (item_name) REFERENCES merch.t_items(name) NOT VALID;


--
-- Name: t_transfer_history t_transfer_history_receiver_name_fkey; Type: FK CONSTRAINT; Schema: merch; Owner: arch
--

ALTER TABLE ONLY merch.t_transfer_history
    ADD CONSTRAINT t_transfer_history_receiver_name_fkey FOREIGN KEY (receiver_name) REFERENCES merch.t_employees(username) NOT VALID;


--
-- Name: t_transfer_history t_transfer_history_sender_name_fkey; Type: FK CONSTRAINT; Schema: merch; Owner: arch
--

ALTER TABLE ONLY merch.t_transfer_history
    ADD CONSTRAINT t_transfer_history_sender_name_fkey FOREIGN KEY (sender_name) REFERENCES merch.t_employees(username);


--
-- PostgreSQL database dump complete
--

