--
-- PostgreSQL database dump
--

-- Dumped from database version 13.8
-- Dumped by pg_dump version 14.5 (Homebrew)

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
-- Name: hospital-mock; Type: DATABASE; Schema: -; Owner: postgres
--

CREATE DATABASE "hospital-mock" WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE = 'en_US.utf8';


ALTER DATABASE "hospital-mock" OWNER TO postgres;

\connect -reuse-previous=on "dbname='hospital-mock'"

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
-- Name: AppointmentStatus; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public."AppointmentStatus" AS ENUM (
    'SCHEDULED',
    'COMPLETED',
    'CANCELLED'
    );


ALTER TYPE public."AppointmentStatus" OWNER TO postgres;

--
-- Name: BloodType; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public."BloodType" AS ENUM (
    'O',
    'A',
    'B',
    'AB'
    );


ALTER TYPE public."BloodType" OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: Appointment; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Appointment" (
                                      id integer NOT NULL,
                                      "patientId" text NOT NULL,
                                      "doctorId" integer NOT NULL,
                                      "startDateTime" timestamp(3) without time zone NOT NULL,
                                      "endDateTime" timestamp(3) without time zone NOT NULL,
                                      detail text NOT NULL,
                                      "nextAppointment" timestamp(3) without time zone,
                                      status public."AppointmentStatus" DEFAULT 'SCHEDULED'::public."AppointmentStatus" NOT NULL,
                                      "createdAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
                                      "updatedAt" timestamp(3) without time zone NOT NULL
);


ALTER TABLE public."Appointment" OWNER TO postgres;

--
-- Name: Appointment_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."Appointment_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."Appointment_id_seq" OWNER TO postgres;

--
-- Name: Appointment_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."Appointment_id_seq" OWNED BY public."Appointment".id;


--
-- Name: Doctor; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Doctor" (
                                 id integer NOT NULL,
                                 initial_th text NOT NULL,
                                 firstname_th text NOT NULL,
                                 lastname_th text NOT NULL,
                                 initial_en text NOT NULL,
                                 firstname_en text NOT NULL,
                                 lastname_en text NOT NULL,
                                 "position" text NOT NULL,
                                 username text NOT NULL,
                                 password text NOT NULL,
                                 "profilePicURL" text NOT NULL,
                                 "createdAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
                                 "updatedAt" timestamp(3) without time zone NOT NULL
);


ALTER TABLE public."Doctor" OWNER TO postgres;

--
-- Name: Doctor_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."Doctor_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."Doctor_id_seq" OWNER TO postgres;

--
-- Name: Doctor_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."Doctor_id_seq" OWNED BY public."Doctor".id;


--
-- Name: Invoice; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Invoice" (
                                  id integer NOT NULL,
                                  "appointmentId" integer NOT NULL,
                                  paid boolean DEFAULT false NOT NULL,
                                  total double precision NOT NULL,
                                  "createdAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
                                  "updatedAt" timestamp(3) without time zone NOT NULL
);


ALTER TABLE public."Invoice" OWNER TO postgres;

--
-- Name: InvoiceDiscount; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."InvoiceDiscount" (
                                          id integer NOT NULL,
                                          name text NOT NULL,
                                          amount double precision NOT NULL,
                                          "invoiceId" integer NOT NULL,
                                          "createdAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
                                          "updatedAt" timestamp(3) without time zone NOT NULL
);


ALTER TABLE public."InvoiceDiscount" OWNER TO postgres;

--
-- Name: InvoiceDiscount_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."InvoiceDiscount_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."InvoiceDiscount_id_seq" OWNER TO postgres;

--
-- Name: InvoiceDiscount_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."InvoiceDiscount_id_seq" OWNED BY public."InvoiceDiscount".id;


--
-- Name: InvoiceItem; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."InvoiceItem" (
                                      id integer NOT NULL,
                                      name text NOT NULL,
                                      price double precision NOT NULL,
                                      quantity integer NOT NULL,
                                      "invoiceId" integer NOT NULL,
                                      "createdAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
                                      "updatedAt" timestamp(3) without time zone NOT NULL
);


ALTER TABLE public."InvoiceItem" OWNER TO postgres;

--
-- Name: InvoiceItem_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."InvoiceItem_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."InvoiceItem_id_seq" OWNER TO postgres;

--
-- Name: InvoiceItem_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."InvoiceItem_id_seq" OWNED BY public."InvoiceItem".id;


--
-- Name: Invoice_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."Invoice_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."Invoice_id_seq" OWNER TO postgres;

--
-- Name: Invoice_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."Invoice_id_seq" OWNED BY public."Invoice".id;


--
-- Name: Medicine; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Medicine" (
                                   id integer NOT NULL,
                                   name text NOT NULL,
                                   description text NOT NULL,
                                   "createdAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
                                   "updatedAt" timestamp(3) without time zone NOT NULL,
                                   "pictureURL" text NOT NULL
);


ALTER TABLE public."Medicine" OWNER TO postgres;

--
-- Name: Medicine_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."Medicine_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."Medicine_id_seq" OWNER TO postgres;

--
-- Name: Medicine_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."Medicine_id_seq" OWNED BY public."Medicine".id;


--
-- Name: Patient; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Patient" (
                                  id text NOT NULL,
                                  initial_th text NOT NULL,
                                  firstname_th text NOT NULL,
                                  lastname_th text NOT NULL,
                                  initial_en text NOT NULL,
                                  firstname_en text NOT NULL,
                                  lastname_en text NOT NULL,
                                  nationality text NOT NULL,
                                  "nationalId" text,
                                  "passportId" text,
                                  "phoneNumber" text NOT NULL,
                                  weight double precision NOT NULL,
                                  height double precision NOT NULL,
                                  "birthDate" timestamp(3) without time zone NOT NULL,
                                  "bloodType" public."BloodType" NOT NULL,
                                  "createdAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
                                  "updatedAt" timestamp(3) without time zone NOT NULL,
                                  "profilePicURL" text NOT NULL
);


ALTER TABLE public."Patient" OWNER TO postgres;

--
-- Name: Prescription; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Prescription" (
                                       id integer NOT NULL,
                                       "medicineId" integer NOT NULL,
                                       "appointmentId" integer NOT NULL,
                                       amount integer NOT NULL,
                                       "createdAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
                                       "updatedAt" timestamp(3) without time zone NOT NULL
);


ALTER TABLE public."Prescription" OWNER TO postgres;

--
-- Name: Prescription_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."Prescription_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public."Prescription_id_seq" OWNER TO postgres;

--
-- Name: Prescription_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."Prescription_id_seq" OWNED BY public."Prescription".id;


--
-- Name: _prisma_migrations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public._prisma_migrations (
                                           id character varying(36) NOT NULL,
                                           checksum character varying(64) NOT NULL,
                                           finished_at timestamp with time zone,
                                           migration_name character varying(255) NOT NULL,
                                           logs text,
                                           rolled_back_at timestamp with time zone,
                                           started_at timestamp with time zone DEFAULT now() NOT NULL,
                                           applied_steps_count integer DEFAULT 0 NOT NULL
);


ALTER TABLE public._prisma_migrations OWNER TO postgres;

--
-- Name: Appointment id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Appointment" ALTER COLUMN id SET DEFAULT nextval('public."Appointment_id_seq"'::regclass);


--
-- Name: Doctor id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Doctor" ALTER COLUMN id SET DEFAULT nextval('public."Doctor_id_seq"'::regclass);


--
-- Name: Invoice id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Invoice" ALTER COLUMN id SET DEFAULT nextval('public."Invoice_id_seq"'::regclass);


--
-- Name: InvoiceDiscount id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."InvoiceDiscount" ALTER COLUMN id SET DEFAULT nextval('public."InvoiceDiscount_id_seq"'::regclass);


--
-- Name: InvoiceItem id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."InvoiceItem" ALTER COLUMN id SET DEFAULT nextval('public."InvoiceItem_id_seq"'::regclass);


--
-- Name: Medicine id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Medicine" ALTER COLUMN id SET DEFAULT nextval('public."Medicine_id_seq"'::regclass);


--
-- Name: Prescription id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Prescription" ALTER COLUMN id SET DEFAULT nextval('public."Prescription_id_seq"'::regclass);


--
-- Data for Name: Appointment; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (12, 'HN-209464', 22, '2023-09-06 06:07:20.472', '2023-09-06 06:37:20.472', 'Eius quibusdam commodi. Voluptas maxime quo ut sint et qui. Eligendi qui sapiente laborum distinctio numquam. Id reprehenderit et. Aperiam id quis corrupti vel consequatur dolores. Similique fugiat quaerat ipsa soluta reprehenderit ullam voluptatem deleniti id.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (13, 'HN-853857', 16, '2023-06-07 02:05:30.575', '2023-06-07 02:35:30.575', 'Vel soluta dolor eos adipisci ipsum ipsum. Harum voluptatem et. Et repellendus nihil veritatis ea ea laudantium provident ut. Corporis qui laudantium aut deleniti placeat animi excepturi itaque.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (14, 'HN-618756', 16, '2023-01-24 22:55:21.63', '2023-01-24 23:25:21.63', 'Fuga voluptatem nobis sint ratione. Quidem fugit fugit repudiandae harum sapiente vitae vero inventore. Vero non excepturi sit repudiandae. Ab perferendis non et. Non quod provident repellendus ut architecto rerum deserunt ut. Rerum consequuntur et.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (15, 'HN-725244', 28, '2022-09-08 01:14:46.857', '2022-09-08 01:44:46.857', 'Vel est alias omnis. Facilis numquam porro aut. Repellat asperiores ad. Sit ut doloribus vero sunt ea aut. Natus cumque ex a adipisci est ut quidem magnam nam. Provident odio repudiandae iusto.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.951');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (16, 'HN-254739', 27, '2022-08-05 21:15:52.633', '2022-08-05 21:45:52.633', 'Aut dolorum eum est assumenda. Quae voluptatem in voluptas ut ex voluptatem numquam et qui. Atque et qui vitae rerum blanditiis iusto vel et. Voluptatem recusandae voluptates tenetur magni ea itaque et. Officia non recusandae consequatur quia fuga. Sunt consequatur qui commodi consequatur.', NULL, 'CANCELLED', '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.951');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (17, 'HN-779528', 7, '2022-09-06 13:03:23.757', '2022-09-06 13:33:23.757', 'Repellendus ducimus rem distinctio fuga qui. Eos et quia doloribus atque non est. Cum adipisci quibusdam est autem consequatur. Quam itaque quidem. Atque harum vel et dolorem a.', NULL, 'CANCELLED', '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.951');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (18, 'HN-254739', 18, '2022-06-01 18:19:31.728', '2022-06-01 18:49:31.728', 'Reprehenderit deleniti enim. Doloribus ipsa enim ut. Quas a enim tenetur voluptatem veniam. Impedit adipisci quo. Animi dolorum reiciendis sed non at suscipit dolor suscipit harum.', NULL, 'CANCELLED', '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.951');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (19, 'HN-997298', 22, '2022-09-07 11:03:45.607', '2022-09-07 11:33:45.607', 'Itaque ex nihil tempore accusantium occaecati incidunt. Beatae amet quam nobis minus placeat. Odit impedit culpa et doloribus. Cupiditate aliquam tempore.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.951');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (20, 'HN-159601', 24, '2022-09-06 19:33:40.989', '2022-09-06 20:03:40.989', 'Culpa occaecati minima velit repellat nisi quam. Quas quia ad labore. Et facilis itaque deleniti architecto at. A est sit magni. Sed et neque minus voluptatem molestiae aut fugiat aut. Consequatur saepe consequatur eaque corporis.', NULL, 'CANCELLED', '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.951');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (21, 'HN-843948', 8, '2022-07-03 05:17:41.951', '2022-07-03 05:47:41.951', 'Accusantium et temporibus magni. Autem recusandae corrupti. Illum earum quia quae ut itaque ipsa officiis.', NULL, 'CANCELLED', '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.951');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (22, 'HN-126977', 5, '2022-09-07 12:38:43.086', '2022-09-07 13:08:43.086', 'Quis excepturi vel aut perferendis minus et quam dolores dignissimos. Nisi quo dolorum et et doloremque dolores culpa aut necessitatibus. Et autem dolor.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.951');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (23, 'HN-618756', 22, '2021-12-12 02:08:11.96', '2021-12-12 02:38:11.96', 'Et fugiat dolorum enim voluptatibus delectus id ipsa explicabo placeat. Ut laborum doloremque. Adipisci sequi tempore cumque molestiae. Corporis tenetur dolor qui corrupti.', '2023-01-30 02:36:50.537', 'COMPLETED', '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (24, 'HN-693045', 2, '2023-05-09 09:17:21.151', '2023-05-09 09:47:21.151', 'Consequatur dignissimos nostrum dolor perspiciatis repudiandae molestiae quo maiores. Quia perspiciatis ipsa dolor atque. Et reiciendis qui sed tenetur. Ut nihil pariatur molestias. Nisi inventore illum dolor repellat voluptatem et voluptas rem.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (25, 'HN-314696', 29, '2022-09-06 23:47:25.171', '2022-09-07 00:17:25.171', 'Modi dolore saepe aliquid. Ab eveniet possimus sed id nemo eos enim amet quas. Veritatis quod vero facilis ratione voluptatem hic optio dignissimos sunt. Temporibus quas eligendi sequi optio aut. Saepe vel nihil sit quo id nihil culpa aut cum. Ipsum quasi nostrum a et voluptates voluptas.', NULL, 'CANCELLED', '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (26, 'HN-897438', 14, '2022-09-07 14:02:00.678', '2022-09-07 14:32:00.678', 'Blanditiis sunt numquam dolor dolorem. Necessitatibus a eveniet sapiente voluptas qui temporibus veritatis. Vero voluptas qui eligendi hic fugiat saepe aliquam. Eum modi atque in laborum nulla magnam tempora animi magnam.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.951');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (27, 'HN-519312', 5, '2021-11-07 05:53:41.031', '2021-11-07 06:23:41.031', 'Quidem quis commodi et nesciunt. Perferendis nihil quia adipisci vel repellendus. Neque neque ea quisquam quam sit est labore. Voluptatum tempora totam velit.', NULL, 'CANCELLED', '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (28, 'HN-188305', 3, '2022-04-23 10:08:02.616', '2022-04-23 10:38:02.616', 'Eveniet alias qui quidem quod velit sed consectetur voluptatem. Veritatis ut recusandae et reprehenderit cupiditate in et laborum voluptatem. Sequi mollitia qui dolorem.', '2023-04-03 13:29:11.101', 'COMPLETED', '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (29, 'HN-919225', 1, '2022-09-07 21:43:50.288', '2022-09-07 22:13:50.288', 'Quos magni voluptatem ratione explicabo perferendis sunt quis quis. Voluptatem et cumque ea. Alias aut quo. Quia eligendi rerum at.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.951');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (32, 'HN-897438', 19, '2022-12-08 08:45:32.172', '2022-12-08 09:15:32.172', 'Fugiat tempore aut qui voluptatem quae vel vitae sint tempora. Sint quasi quo magni sit architecto culpa. Impedit architecto quia ipsum molestiae et. Ut aut ut ea sed aliquam aut saepe.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (43, 'HN-853857', 19, '2023-05-13 05:59:48.796', '2023-05-13 06:29:48.796', 'In commodi nesciunt magnam eius vitae eos ea adipisci. Eum explicabo corrupti dignissimos facilis. Totam reiciendis impedit maxime. Numquam porro iusto eius quaerat est omnis sint iusto. Adipisci et rerum culpa voluptas temporibus rerum. Consequatur impedit accusantium similique omnis.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.952');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (52, 'HN-126977', 25, '2023-03-20 11:31:26.924', '2023-03-20 12:01:26.924', 'Ut aliquid delectus iste neque dolorum porro eum consequatur. Saepe illum eius. Possimus molestiae molestiae harum tempora vero quos et. Et tempore eos quae sequi.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.952');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (62, 'HN-881131', 30, '2022-09-07 12:05:49.489', '2022-09-07 12:35:49.489', 'Eius qui enim. Aliquid odit in rerum. Dolorem aut facere. Natus occaecati excepturi quae placeat voluptates quisquam hic natus exercitationem. Enim minus et.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (67, 'HN-629664', 15, '2022-03-29 15:24:54.407', '2022-03-29 15:54:54.407', 'Eos non libero aliquid sit dolorem veniam odit nisi. Expedita id sint molestias eum cumque veniam sunt. Corrupti enim voluptatibus amet blanditiis vel architecto et quod.', NULL, 'CANCELLED', '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.953');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (76, 'HN-784077', 17, '2022-09-07 09:32:51.088', '2022-09-07 10:02:51.088', 'Praesentium fugit placeat non eum et et. Id quos qui. Enim suscipit autem ducimus quia sapiente voluptas. Assumenda sit laboriosam architecto.', NULL, 'CANCELLED', '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (90, 'HN-179590', 5, '2022-09-06 18:32:35.477', '2022-09-06 19:02:35.477', 'Ut odit qui itaque facere. Velit recusandae vero dicta id id. Dignissimos earum consequatur ratione aperiam voluptatem eveniet earum excepturi non.', NULL, 'CANCELLED', '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.95');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (30, 'HN-129512', 29, '2022-12-31 05:32:12.722', '2022-12-31 06:02:12.722', 'Dolorum accusamus est sequi dolores sed. Omnis et quod officiis impedit fugit dolores architecto velit. Veritatis est et voluptas quod aut nesciunt est sit. Harum commodi qui consequatur aliquam.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (31, 'HN-119713', 10, '2022-07-13 08:10:44.169', '2022-07-13 08:40:44.169', 'Et molestias autem voluptas voluptatem numquam iure corporis quo error. Qui qui assumenda quo expedita fuga. Voluptatem at numquam cumque nihil. Provident error aperiam molestiae voluptates.', NULL, 'COMPLETED', '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.951');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (11, 'HN-757364', 20, '2021-10-22 04:51:30.537', '2021-10-22 05:21:30.537', 'Quo ut accusamus consequatur sint vel rem hic tempora esse. Itaque error repellat reiciendis ab expedita quia iure. Repellat qui debitis expedita.', NULL, 'COMPLETED', '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (33, 'HN-519312', 1, '2023-07-14 02:54:44.048', '2023-07-14 03:24:44.048', 'Magnam eius molestiae soluta. Voluptatem culpa accusamus eaque est. Optio explicabo dolor quae est deserunt. Consequatur est praesentium et cum quo laudantium quidem est. Dolore distinctio et nulla aut nam fugit harum exercitationem. Odio accusamus perspiciatis quod et aut dolorem eos.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (44, 'HN-861496', 27, '2022-06-06 01:06:20.497', '2022-06-06 01:36:20.497', 'Quod architecto asperiores ipsum ea quis ut tempora distinctio officiis. Tenetur iusto maxime non autem eveniet odio magnam repellendus quidem. Asperiores repudiandae rerum officiis velit.', NULL, 'CANCELLED', '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (47, 'HN-853857', 16, '2023-08-10 17:25:34.98', '2023-08-10 17:55:34.98', 'Quia molestiae animi. Animi et aliquam ut minima consequatur. Aliquam nulla repudiandae nostrum odio sapiente. Incidunt ut eius eos. Et voluptatem quae quia ex. Illum facere accusantium fugit voluptatem eligendi excepturi.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (53, 'HN-910978', 11, '2022-09-06 23:45:54.808', '2022-09-07 00:15:54.808', 'Porro occaecati ullam voluptas quos impedit quia earum. Sit perferendis in quibusdam consequatur quia. Et et laboriosam cumque excepturi provident mollitia quo dolore.', NULL, 'COMPLETED', '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (58, 'HN-162462', 6, '2022-01-03 01:54:49.11', '2022-01-03 02:24:49.11', 'Maiores ipsum officia rerum. Voluptate omnis voluptatem rem ut. Harum iusto dicta et et ut.', NULL, 'COMPLETED', '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (66, 'HN-919225', 7, '2023-06-17 13:46:25.884', '2023-06-17 14:16:25.884', 'Inventore magnam eum odio quam necessitatibus. Molestiae ea sunt dolore aut voluptatibus consequatur quia cumque. Quia sit aut fugiat dicta.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.953');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (68, 'HN-323662', 10, '2023-05-06 16:48:03.439', '2023-05-06 17:18:03.439', 'Veniam iure sint doloribus architecto molestiae aut id tempora. Amet iste sit perferendis velit mollitia ullam. Voluptatem laboriosam nobis est id quam quaerat non. Et id soluta non. Et quae ea similique assumenda minima quia voluptates voluptatem. Sit necessitatibus quisquam.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.953');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (77, 'HN-633981', 2, '2023-04-06 03:36:27.46', '2023-04-06 04:06:27.46', 'Quo aut omnis quod aut accusamus sed rerum. Officiis non optio adipisci ut expedita. Pariatur qui animi enim a consequatur consectetur cum. Consequatur saepe est quo ex. Temporibus officiis et optio ab corporis corrupti iure provident odit.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (78, 'HN-883377', 10, '2023-09-03 04:17:21.949', '2023-09-03 04:47:21.949', 'Porro voluptas dolorem ad explicabo corporis quidem. Sit incidunt corporis. Quia ut impedit est et error quas blanditiis odio. Aperiam quo a perferendis aut dolor. Eveniet laudantium fuga voluptatem.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (88, 'HN-873016', 20, '2022-09-08 09:18:26.159', '2022-09-08 09:48:26.159', 'Et atque possimus nam. Fuga numquam quas molestiae sunt rerum. Omnis fuga et qui molestias. Minus quaerat dolores magnam voluptatem non esse qui. Sed quis ea ea pariatur.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.954', '2022-09-07 10:52:46.954');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (98, 'HN-285237', 26, '2021-12-31 07:07:35.921', '2021-12-31 07:37:35.921', 'Est illo recusandae mollitia qui. Maiores non omnis. Atque deserunt cumque autem. Exercitationem molestiae earum quibusdam. Dolor qui minus vel occaecati. Minus ut consequatur voluptatem eveniet ab et ut iure.', '2023-01-29 06:21:18.349', 'COMPLETED', '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (34, 'HN-113321', 2, '2022-09-14 16:39:23.036', '2022-09-14 17:09:23.036', 'Soluta illo vero sapiente voluptatem et quo dolorum a. Adipisci labore illum molestiae dolores. Quae quis sequi eum at ratione qui voluptatem quae. Veniam modi dolores earum sequi sit voluptas facere. Expedita eaque voluptatibus architecto debitis id autem veniam repellat suscipit.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (45, 'HN-561154', 24, '2022-09-07 10:24:57.746', '2022-09-07 10:54:57.746', 'Quibusdam eos pariatur modi sunt incidunt quisquam quia nulla modi. Veniam delectus deleniti voluptas quis excepturi facere qui amet. Quaerat magnam saepe consequuntur illo eos. Quibusdam sed dolorem eos illo facere voluptatem sed. Corrupti et numquam dolores dignissimos sequi corporis odio omnis. Accusantium laboriosam voluptatem similique nihil eius enim ut officiis quam.', '2023-07-18 20:25:25.047', 'COMPLETED', '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (55, 'HN-137318', 5, '2022-09-07 23:48:00.437', '2022-09-08 00:18:00.437', 'Voluptates consequuntur et. Quibusdam accusantium omnis. Facilis sit alias et aut nam cum voluptas. Dolorem quos ratione non qui.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (71, 'HN-200755', 6, '2022-09-08 02:14:18.583', '2022-09-08 02:44:18.583', 'Et cum debitis distinctio accusamus vel et qui et numquam. Ex vitae sit voluptas est vel commodi harum et iusto. Omnis magni dolorem accusamus repudiandae voluptate mollitia in consequatur aut.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.953');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (84, 'HN-932425', 13, '2022-07-09 17:49:14.073', '2022-07-09 18:19:14.073', 'Ex laborum nam dolores nulla perspiciatis accusamus facilis. Mollitia omnis nam. Et unde sequi perferendis mollitia fugiat maiores nam et.', '2022-09-11 22:14:02.591', 'COMPLETED', '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (35, 'HN-698271', 12, '2023-03-01 14:52:00.743', '2023-03-01 15:22:00.743', 'Alias nemo tenetur quidem eligendi incidunt quaerat. Et amet a quia repudiandae eligendi ex nisi aliquam aut. Facere vero nemo sed quidem a rerum in. Excepturi fuga placeat debitis dolore sapiente.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (42, 'HN-618756', 4, '2022-09-08 02:59:00.286', '2022-09-08 03:29:00.286', 'Molestiae dolorum vel laborum possimus. Dolorem quia laudantium velit quis. Quis dignissimos hic excepturi placeat expedita sed. Animi officiis quia quasi nam quia enim dignissimos eligendi aut.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.952');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (49, 'HN-740157', 4, '2022-09-06 14:58:21.154', '2022-09-06 15:28:21.154', 'Laudantium accusamus numquam adipisci. Blanditiis error cupiditate. Minus facere enim quia a accusantium dolor pariatur corporis accusamus.', NULL, 'COMPLETED', '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (69, 'HN-209014', 15, '2022-09-07 21:33:03.268', '2022-09-07 22:03:03.268', 'Architecto et molestiae omnis fugiat totam. Repudiandae at qui nesciunt ducimus sit unde in et possimus. Autem tempora et explicabo praesentium et voluptatum sint quidem eos. Aut aliquam ut. Doloribus ex itaque sed consequuntur consectetur qui sint veniam.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.953');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (81, 'HN-932425', 2, '2023-03-12 11:59:49.753', '2023-03-12 12:29:49.753', 'Provident quia et ut sed numquam libero. Est voluptatem ea nisi quasi voluptatem maiores vero expedita et. Sunt et repellat.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (92, 'HN-162462', 17, '2023-01-07 20:55:20.356', '2023-01-07 21:25:20.356', 'Sed quo ipsum tenetur quod et sed tenetur. Earum ea culpa veritatis. Non qui consequuntur voluptates perspiciatis magnam. Nulla perferendis molestiae impedit quia numquam repudiandae unde. Reiciendis voluptates rerum.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.954', '2022-09-07 10:52:46.954');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (36, 'HN-209014', 9, '2022-09-08 10:26:00.093', '2022-09-08 10:56:00.093', 'Laboriosam suscipit maxime reiciendis necessitatibus voluptatem. Non ducimus repellendus hic minima repellendus fuga. Aut odio dolor minus sequi ut repudiandae velit. Voluptates eveniet deleniti vel aut corrupti impedit unde rerum. Tempore ut est in consequatur repudiandae laboriosam. Iusto aut eum temporibus tempora corporis.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (41, 'HN-910978', 8, '2022-12-13 10:01:51.768', '2022-12-13 10:31:51.768', 'Placeat nihil ullam praesentium. In ab fugit qui aut tenetur non. Corporis consequatur unde aut minus rerum expedita voluptatem beatae temporibus.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.952');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (54, 'HN-651394', 1, '2022-11-17 15:12:51.276', '2022-11-17 15:42:51.276', 'Cumque id nihil corrupti quo id voluptas nihil. Voluptatem non et totam odit adipisci qui. Sed laborum praesentium nobis aut possimus repellendus dolores. Ut deleniti ipsa mollitia maiores non voluptatum fugit itaque libero. Consectetur sunt explicabo et odit maiores aut.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.952');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (73, 'HN-421398', 2, '2021-10-19 04:41:07.012', '2021-10-19 05:11:07.012', 'Reprehenderit omnis rerum sit neque et eaque excepturi dolor. Possimus eaque fuga aut qui. Et explicabo officia. Nihil natus voluptatem ut totam omnis itaque sit voluptas aspernatur. Eos impedit et dolores quibusdam nobis sequi esse.', NULL, 'COMPLETED', '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (60, 'HN-881131', 19, '2022-09-07 16:53:01.227', '2022-09-07 17:23:01.227', 'Reiciendis aut veritatis atque quis sed consequatur enim enim. Iusto blanditiis ut repellendus ab qui fuga. Eum eligendi enim error aut rerum labore adipisci voluptatem temporibus. Error nobis quam et est possimus aut et. Quasi sit odio labore assumenda consectetur enim consequatur ut qui. Eum dicta aspernatur.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (72, 'HN-934746', 24, '2023-03-31 02:47:11.009', '2023-03-31 03:17:11.009', 'Inventore voluptates pariatur quae magnam totam aut voluptatem. Animi ratione laboriosam. Ut doloribus laboriosam rerum modi voluptas exercitationem repellat provident assumenda. Eveniet fugiat ut.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.953');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (83, 'HN-517246', 21, '2022-09-07 02:42:19.376', '2022-09-07 03:12:19.376', 'Et explicabo sunt illo quaerat. Odit inventore est quia sed amet reprehenderit rerum. Facilis nobis voluptatibus sed. Veritatis dolor at fuga esse. Mollitia animi ab non esse. Quo aliquam et necessitatibus rerum dolore voluptatem modi nulla eius.', NULL, 'CANCELLED', '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (94, 'HN-113321', 2, '2023-04-07 15:19:18.252', '2023-04-07 15:49:18.252', 'Velit adipisci modi dolore non aut eveniet adipisci vel. Ut qui enim eius aut autem doloremque. Non eligendi est aliquid similique quo. Qui accusamus rerum nobis rerum veritatis asperiores fugit.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (37, 'HN-562380', 12, '2022-09-07 01:17:23.081', '2022-09-07 01:47:23.081', 'Illum nam in sunt iure et. In et quaerat qui suscipit fugit qui dolorem nostrum. Totam distinctio aspernatur aut nulla sit consequuntur debitis nemo eaque. Quaerat qui repellendus et molestias. Iste sunt et ullam distinctio et. Excepturi sed laborum aspernatur inventore.', '2023-06-16 15:58:38.764', 'COMPLETED', '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (59, 'HN-963617', 14, '2022-09-07 02:41:12.199', '2022-09-07 03:11:12.199', 'Saepe et et dolorem. Deleniti quis alias quis magni perferendis cupiditate. Est molestias autem consequatur. Qui praesentium commodi beatae nostrum assumenda sint quidem dignissimos sequi.', NULL, 'COMPLETED', '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (75, 'HN-853857', 20, '2022-09-07 19:06:21.417', '2022-09-07 19:36:21.417', 'Impedit possimus aut. Aut doloribus ex et aut aperiam. Numquam molestiae est at inventore atque.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.953');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (93, 'HN-126977', 4, '2022-08-16 02:33:52.253', '2022-08-16 03:03:52.253', 'Aliquam aut dolor. Dolores ut tempore aliquam. Non facilis adipisci vel qui aut velit harum.', '2022-10-17 20:35:45.695', 'COMPLETED', '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (38, 'HN-883377', 12, '2022-09-07 05:08:19.726', '2022-09-07 05:38:19.726', 'Mollitia facere fugit cumque odio. Fugit occaecati ratione. Totam facilis provident autem maiores dolores cum quia est nesciunt.', '2023-06-28 09:50:45.821', 'COMPLETED', '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (39, 'HN-740755', 28, '2021-12-23 12:22:15.373', '2021-12-23 12:52:15.373', 'Ratione saepe aut vero non fugiat harum sed doloribus. Veniam neque dolorum iusto quos quis. Nesciunt dignissimos exercitationem soluta labore aut.', NULL, 'COMPLETED', '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (50, 'HN-285237', 11, '2022-04-28 01:13:21.323', '2022-04-28 01:43:21.323', 'Laudantium ut sit. Est autem illo voluptas ea autem praesentium aut magnam voluptatibus. Incidunt eius et beatae aut impedit necessitatibus nemo.', NULL, 'COMPLETED', '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.952');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (63, 'HN-443354', 6, '2022-09-07 08:21:02.519', '2022-09-07 08:51:02.519', 'Ipsum aspernatur et corrupti in assumenda quis doloremque. Vel sit voluptates alias repudiandae minima laudantium sint. Deleniti cumque et.', NULL, 'COMPLETED', '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (70, 'HN-297559', 20, '2022-08-07 09:20:40.931', '2022-08-07 09:50:40.931', 'Ducimus aspernatur eius dolor reiciendis et quo. Unde aut maxime autem nisi. Omnis eos rerum velit hic. In minima aut quod nemo voluptatem dolores soluta aperiam et. Tempore ea et earum suscipit sunt eum. Quia voluptatem repellat distinctio adipisci quia qui.', '2023-06-15 04:46:57.383', 'COMPLETED', '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (82, 'HN-629664', 18, '2022-07-01 07:17:52.166', '2022-07-01 07:47:52.166', 'Quam ut deserunt porro. Rerum quam rem consequatur impedit deserunt. Iste recusandae impedit tenetur minima ut.', NULL, 'CANCELLED', '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (86, 'HN-844945', 23, '2023-03-03 10:45:45.945', '2023-03-03 11:15:45.945', 'Animi officiis hic aut distinctio doloribus commodi. Quo quis quia explicabo dolores voluptas quis quaerat et aut. Ipsa aspernatur est non. Facilis consectetur magnam porro debitis sed. Excepturi consequuntur voluptatem vitae amet. Et labore eveniet sit hic officiis.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.953');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (91, 'HN-179590', 25, '2022-09-07 02:28:12.616', '2022-09-07 02:58:12.616', 'Et assumenda dolor eveniet et neque quia praesentium veniam. Eum et aut doloremque quisquam consequatur veniam suscipit. Alias beatae vero assumenda ut. Pariatur ducimus quo doloremque aspernatur aliquam eaque numquam. Rerum rerum voluptate soluta est illo ut. Et aut atque et sunt omnis.', NULL, 'CANCELLED', '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.95');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (99, 'HN-517289', 10, '2022-08-20 14:49:13.55', '2022-08-20 15:19:13.55', 'Et a aliquid eum velit sapiente blanditiis voluptas. Odio perferendis qui voluptatem et dignissimos alias reiciendis. Earum officiis et modi fugiat sint. Id ducimus est et et ad nemo iste. Nesciunt incidunt asperiores error et deserunt voluptas eaque nostrum.', '2023-02-22 05:42:17.211', 'COMPLETED', '2022-09-07 10:52:46.954', '2022-09-07 10:52:46.954');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (40, 'HN-803959', 6, '2022-09-06 16:56:09.01', '2022-09-06 17:26:09.01', 'Non voluptatem ut molestiae vel officiis sequi odit. Hic incidunt distinctio et. Repellendus sed exercitationem ea deserunt sed molestias incidunt aspernatur. Praesentium consequuntur repellat quidem velit quo vero soluta. Magnam commodi illo dolor tenetur animi.', NULL, 'CANCELLED', '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (46, 'HN-347868', 18, '2022-09-07 09:09:31.443', '2022-09-07 09:39:31.443', 'Qui ipsum blanditiis voluptas qui vel unde delectus est exercitationem. Ab labore aut. Dolores ut cum est aliquam ea. Laborum earum labore omnis voluptatem. Et non ut ea amet consectetur optio in possimus.', NULL, 'CANCELLED', '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (51, 'HN-209464', 15, '2022-09-06 11:04:29.767', '2022-09-06 11:34:29.767', 'Blanditiis delectus qui maxime. Velit non vero. Magni cumque mollitia et quia ex quia. Laborum hic sequi. Nihil sequi ut amet eius necessitatibus quia sit ducimus non.', NULL, 'CANCELLED', '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.952');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (56, 'HN-450439', 23, '2022-09-07 09:27:40.095', '2022-09-07 09:57:40.095', 'Corporis consectetur rerum. Exercitationem deleniti eos sed quod dolor. Et commodi aut corrupti quasi vel itaque est laboriosam ut. Deleniti sint omnis delectus quia ipsam magni et impedit. Quasi expedita odit. Dolore pariatur omnis rerum aut molestias.', NULL, 'COMPLETED', '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (57, 'HN-835237', 19, '2022-09-08 04:40:44.362', '2022-09-08 05:10:44.362', 'Velit hic dolores deleniti doloremque dignissimos totam. Earum eligendi nulla pariatur sed quaerat voluptatem itaque facilis soluta. Ut doloribus aut sed ea mollitia deserunt et veritatis. Neque numquam recusandae dolor autem at. Qui velit non iure. Repudiandae eaque ab qui quia ut assumenda.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.952');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (61, 'HN-129512', 9, '2022-04-16 01:47:25.487', '2022-04-16 02:17:25.487', 'Sapiente fugit deleniti est eaque nam odit ut. Nulla incidunt optio in molestiae nemo. Quia nisi earum ea quisquam in amet est et.', NULL, 'CANCELLED', '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (64, 'HN-221203', 22, '2022-09-07 01:53:41.308', '2022-09-07 02:23:41.308', 'Non excepturi omnis ut. Minima est et accusantium ea similique voluptatibus aliquid ea. Modi molestiae sed doloremque qui mollitia amet consequatur officiis.', '2023-06-18 11:45:42.049', 'COMPLETED', '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (1, 'HN-693045', 24, '2023-08-09 14:09:30.449', '2023-08-09 14:39:30.449', 'Assumenda quia et aspernatur sit quidem maxime qui. Perferendis aut nihil eveniet. Quo accusamus sint totam dolorem. Modi sed minima eius autem occaecati exercitationem voluptatum. Impedit assumenda sit ut quia. Tempora ut error ex vero dicta voluptatem.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.95');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (2, 'HN-414878', 23, '2022-04-09 04:24:38.009', '2022-04-09 04:54:38.009', 'Officia totam rerum et. Dolorum reiciendis veniam et et sit qui tenetur. Deleniti voluptas consectetur ut natus officia odio debitis sapiente soluta. Placeat et iure. Corporis dolor fugiat animi.', NULL, 'COMPLETED', '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (3, 'HN-126977', 27, '2022-09-08 03:19:48.698', '2022-09-08 03:49:48.698', 'Molestiae voluptatum aut. Autem reprehenderit ab dignissimos ipsum vitae deleniti vel quis. Odio quasi molestiae.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (4, 'HN-796550', 30, '2021-10-15 08:57:21.519', '2021-10-15 09:27:21.519', 'Quos quas dolorem provident quia aut. Voluptas numquam accusantium quia rerum ad non magni aut aliquam. Qui vel explicabo accusantium velit quos voluptas voluptatem officiis quia.', '2023-07-26 02:29:29.449', 'COMPLETED', '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (5, 'HN-874621', 12, '2022-09-07 19:23:43.144', '2022-09-07 19:53:43.144', 'Magnam laboriosam quia sit ullam nobis corrupti repudiandae voluptatem. Incidunt quos excepturi error ipsam vero. Aut tempore cupiditate quia. Quibusdam debitis animi aut repellendus molestias voluptas et perferendis. Omnis perferendis sit est qui iure qui dolorem rerum aperiam. Tempora tenetur excepturi.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (6, 'HN-634304', 20, '2022-03-31 04:23:17.078', '2022-03-31 04:53:17.078', 'Quae et sint similique temporibus. Officia assumenda animi laborum sequi quia. Illum cumque ab. Velit ipsam facere qui provident enim sunt sit.', NULL, 'CANCELLED', '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (7, 'HN-835237', 10, '2023-02-25 00:48:02.597', '2023-02-25 01:18:02.597', 'Magnam odio odit enim dolorem. Voluptatem aperiam quis laboriosam quam eos quod culpa nihil consequatur. Fuga dolore nesciunt eveniet.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (8, 'HN-414878', 3, '2022-09-06 02:45:24.798', '2022-09-06 03:15:24.798', 'Et rerum sed eaque deleniti et eveniet rerum voluptatibus est. Ratione quos qui id itaque reiciendis. Aut commodi rerum dolorum est vel ipsum. Vel impedit error nobis dicta sit.', NULL, 'CANCELLED', '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (9, 'HN-190919', 29, '2021-12-07 15:12:03.825', '2021-12-07 15:42:03.825', 'Amet illum sit aut rem qui aut deleniti ullam magni. Eum est nesciunt aliquid consectetur reiciendis rerum. Id deserunt molestiae. Iure possimus libero iste quibusdam delectus voluptas odio. Est fuga velit.', NULL, 'COMPLETED', '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (10, 'HN-162951', 28, '2022-09-08 08:06:46.307', '2022-09-08 08:36:46.307', 'Excepturi delectus similique ea laudantium incidunt repudiandae eum. Pariatur similique optio alias eius quis. Dolorem perspiciatis quis fugit.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.951');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (65, 'HN-853857', 18, '2022-09-08 00:03:37.478', '2022-09-08 00:33:37.478', 'Molestiae velit aut nam voluptatem aut voluptatem. Vitae omnis non tempora praesentium aut autem veritatis maiores aliquam. Ipsum tempora blanditiis debitis veritatis ea ut repellendus similique.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.953');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (48, 'HN-873016', 11, '2023-07-18 06:26:45.882', '2023-07-18 06:56:45.882', 'Et ab tempore nemo sunt qui. Rem eum sed eum sint autem et distinctio illo sunt. Perspiciatis itaque at consequatur praesentium et et blanditiis. Numquam sed alias quia eaque.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.952');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (74, 'HN-725244', 4, '2022-09-06 23:40:12.482', '2022-09-07 00:10:12.482', 'Perferendis ipsa minus illum consequuntur tenetur blanditiis consectetur. Veritatis soluta aut molestiae sit quia aut saepe repudiandae quia. Non enim voluptatibus aperiam ullam adipisci. Ad temporibus cum consectetur.', NULL, 'CANCELLED', '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.953');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (79, 'HN-993787', 3, '2022-09-07 21:29:14.744', '2022-09-07 21:59:14.744', 'Accusamus iste odit fugiat autem odit omnis aut quae. Aut deserunt tenetur aut. Est ipsum voluptatum blanditiis eius velit est illo a voluptas. Minus voluptas itaque. Et ullam fugiat praesentium quas.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.953');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (80, 'HN-297559', 1, '2023-08-28 12:37:40.244', '2023-08-28 13:07:40.244', 'In aut qui nostrum. Sed id aliquam aut quas rerum voluptatem. Cum nobis reprehenderit cumque provident tempora eligendi. Optio dolorem est odit.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (85, 'HN-209464', 9, '2022-09-07 06:41:54.9', '2022-09-07 07:11:54.9', 'Et magnam beatae quibusdam quia. Mollitia voluptatem eius omnis voluptatem quis. Distinctio fugit consequatur.', '2022-10-14 08:58:42.13', 'COMPLETED', '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (87, 'HN-159601', 9, '2022-09-07 04:00:57.836', '2022-09-07 04:30:57.836', 'Placeat autem ducimus ea tempore quibusdam. Autem ipsum nesciunt debitis voluptatibus unde et doloremque. Delectus et laudantium laudantium incidunt ipsa molestiae. Aut quia provident quibusdam nostrum omnis.', NULL, 'CANCELLED', '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (89, 'HN-262701', 24, '2022-02-07 08:49:46.417', '2022-02-07 09:19:46.417', 'Explicabo molestias ex deleniti molestias. Placeat aut architecto culpa praesentium. Voluptatum maxime sed est magnam mollitia assumenda. Nihil aut vero. Est aperiam est molestiae dolorem rerum illo nihil voluptatum praesentium. Ex dolores tempore omnis sunt.', NULL, 'COMPLETED', '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (95, 'HN-951062', 27, '2022-09-08 05:05:51.765', '2022-09-08 05:35:51.765', 'Temporibus mollitia quis. Assumenda quis atque deserunt dolore ut. Repellendus ratione deserunt accusamus nemo similique cumque dolore. Sed unde aut saepe enim voluptas et. Temporibus impedit libero earum assumenda. Accusantium facere quidem aliquam maxime.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (96, 'HN-163878', 11, '2021-09-10 15:22:03.598', '2021-09-10 15:52:03.598', 'Ut ullam laudantium sit expedita et. Possimus nihil nihil. Veritatis et eos quidem omnis itaque ut a libero corporis. Et sed non vel molestiae natus facilis aspernatur libero.', NULL, 'COMPLETED', '2022-09-07 10:52:46.954', '2022-09-07 10:52:46.955');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (97, 'HN-846103', 1, '2022-09-07 09:38:25.587', '2022-09-07 10:08:25.587', 'Similique eaque cumque doloremque fuga blanditiis placeat. Qui nesciunt aut consequatur. Minus id magni repellendus. Quis officiis optio dicta sed.', '2022-09-07 19:00:11.081', 'COMPLETED', '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.953');
INSERT INTO public."Appointment" (id, "patientId", "doctorId", "startDateTime", "endDateTime", detail, "nextAppointment", status, "createdAt", "updatedAt") VALUES (100, 'HN-162462', 9, '2022-09-07 16:35:10.023', '2022-09-07 17:05:10.023', 'Minima quidem nostrum quo mollitia. Qui ad qui quo. Placeat ut similique. Commodi occaecati temporibus deserunt aut ipsa. Corporis culpa rerum dolore quia.', NULL, 'SCHEDULED', '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');


--
-- Data for Name: Doctor; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public."Doctor" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, "position", username, password, "profilePicURL", "createdAt", "updatedAt") VALUES (1, 'Mr._ไทย', 'Glen_ไทย', 'Medhurst_ไทย', 'Mr.', 'Glen', 'Medhurst', 'Neurologist', 'Zack_Metz72', '$2b$10$PxEZMSF13RVS9jE5aMQnNunK9judFgW49qrbJPakCA8hMEfCS6DoC', 'https://cloudflare-ipfs.com/ipfs/Qmd3W5DuhgHirLHGVixi6V76LhCkZUz6pnFt5AJBiyvHye/avatar/1000.jpg', '2022-09-07 10:52:46.89', '2022-09-07 10:52:46.891');
INSERT INTO public."Doctor" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, "position", username, password, "profilePicURL", "createdAt", "updatedAt") VALUES (2, 'Dr._ไทย', 'Danika_ไทย', 'Ryan_ไทย', 'Dr.', 'Danika', 'Ryan', 'Orthopedist', 'Colin84', '$2b$10$TNZUvuL2fZ6TEV5/RRcnKeXVi6PCeRAz8KF0hgawPBfoMMSNQVaPa', 'https://cloudflare-ipfs.com/ipfs/Qmd3W5DuhgHirLHGVixi6V76LhCkZUz6pnFt5AJBiyvHye/avatar/522.jpg', '2022-09-07 10:52:46.89', '2022-09-07 10:52:46.891');
INSERT INTO public."Doctor" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, "position", username, password, "profilePicURL", "createdAt", "updatedAt") VALUES (3, 'Mrs._ไทย', 'Arianna_ไทย', 'Gislason_ไทย', 'Mrs.', 'Arianna', 'Gislason', 'Orthopedist', 'Hattie.Thompson46', '$2b$10$oRgQUk0f8laCILlhQBu4I.HzljUTvZhuSQM0hT8OFhxcoJ6sKXIoC', 'https://cloudflare-ipfs.com/ipfs/Qmd3W5DuhgHirLHGVixi6V76LhCkZUz6pnFt5AJBiyvHye/avatar/239.jpg', '2022-09-07 10:52:46.89', '2022-09-07 10:52:46.891');
INSERT INTO public."Doctor" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, "position", username, password, "profilePicURL", "createdAt", "updatedAt") VALUES (4, 'Mrs._ไทย', 'Brando_ไทย', 'Franey_ไทย', 'Mrs.', 'Brando', 'Franey', 'Psychiatrist', 'Christophe.Boehm', '$2b$10$aL9lu9TvwChG4R/QiSgi/eNH96lGGTAEnC7X.IRSIU1vIIqhi/djq', 'https://cloudflare-ipfs.com/ipfs/Qmd3W5DuhgHirLHGVixi6V76LhCkZUz6pnFt5AJBiyvHye/avatar/438.jpg', '2022-09-07 10:52:46.89', '2022-09-07 10:52:46.891');
INSERT INTO public."Doctor" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, "position", username, password, "profilePicURL", "createdAt", "updatedAt") VALUES (5, 'Mrs._ไทย', 'Faustino_ไทย', 'Doyle_ไทย', 'Mrs.', 'Faustino', 'Doyle', 'Orthopedist', 'Elias_Wolf', '$2b$10$JTpxlMzpF3MZWTH2HkpWGu0JZhTxeHBw7mtAwDbj.r6uvleOLEEAu', 'https://cloudflare-ipfs.com/ipfs/Qmd3W5DuhgHirLHGVixi6V76LhCkZUz6pnFt5AJBiyvHye/avatar/1033.jpg', '2022-09-07 10:52:46.89', '2022-09-07 10:52:46.891');
INSERT INTO public."Doctor" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, "position", username, password, "profilePicURL", "createdAt", "updatedAt") VALUES (6, 'Miss_ไทย', 'Tod_ไทย', 'Zemlak_ไทย', 'Miss', 'Tod', 'Zemlak', 'Podiatrist', 'Jane_Lueilwitz', '$2b$10$uVFz.Rb6HX.zJGmmbPRLw.GL1MBp8MGWavBUpbfsbly9OSIlKMATq', 'https://cloudflare-ipfs.com/ipfs/Qmd3W5DuhgHirLHGVixi6V76LhCkZUz6pnFt5AJBiyvHye/avatar/1237.jpg', '2022-09-07 10:52:46.89', '2022-09-07 10:52:46.891');
INSERT INTO public."Doctor" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, "position", username, password, "profilePicURL", "createdAt", "updatedAt") VALUES (7, 'Ms._ไทย', 'Erwin_ไทย', 'Leannon_ไทย', 'Ms.', 'Erwin', 'Leannon', 'Podiatrist', 'Greyson27', '$2b$10$wLL1SywRUaZEO8QpWob6..f96pr9BTl.DUPicLJydUnJNqaewd/d2', 'https://cloudflare-ipfs.com/ipfs/Qmd3W5DuhgHirLHGVixi6V76LhCkZUz6pnFt5AJBiyvHye/avatar/1081.jpg', '2022-09-07 10:52:46.89', '2022-09-07 10:52:46.891');
INSERT INTO public."Doctor" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, "position", username, password, "profilePicURL", "createdAt", "updatedAt") VALUES (8, 'Miss_ไทย', 'Darrion_ไทย', 'DuBuque_ไทย', 'Miss', 'Darrion', 'DuBuque', 'Psychiatrist', 'Martin.Bartell33', '$2b$10$BujzVDBOuJtCpFbBC.CtweBJZUhqS.Qg0LgW0sAq3MRS8ol2xuB/i', 'https://cloudflare-ipfs.com/ipfs/Qmd3W5DuhgHirLHGVixi6V76LhCkZUz6pnFt5AJBiyvHye/avatar/584.jpg', '2022-09-07 10:52:46.89', '2022-09-07 10:52:46.891');
INSERT INTO public."Doctor" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, "position", username, password, "profilePicURL", "createdAt", "updatedAt") VALUES (9, 'Miss_ไทย', 'Mustafa_ไทย', 'Nikolaus_ไทย', 'Miss', 'Mustafa', 'Nikolaus', 'Orthopedist', 'Grady_Hane', '$2b$10$pj5cqeVHId1RJfjlhv0al.xqSTheUnmy7vT0RxAC.LP8DAJrT3dVy', 'https://cloudflare-ipfs.com/ipfs/Qmd3W5DuhgHirLHGVixi6V76LhCkZUz6pnFt5AJBiyvHye/avatar/357.jpg', '2022-09-07 10:52:46.89', '2022-09-07 10:52:46.891');
INSERT INTO public."Doctor" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, "position", username, password, "profilePicURL", "createdAt", "updatedAt") VALUES (10, 'Ms._ไทย', 'Tristin_ไทย', 'Kuvalis_ไทย', 'Ms.', 'Tristin', 'Kuvalis', 'Cardiologist', 'Camila86', '$2b$10$P7lgNB0cKcn20UoVcrSwL.E.DQho6VUALcqSNhL4QFr4r9sIUgcv6', 'https://cloudflare-ipfs.com/ipfs/Qmd3W5DuhgHirLHGVixi6V76LhCkZUz6pnFt5AJBiyvHye/avatar/916.jpg', '2022-09-07 10:52:46.89', '2022-09-07 10:52:46.891');
INSERT INTO public."Doctor" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, "position", username, password, "profilePicURL", "createdAt", "updatedAt") VALUES (11, 'Mrs._ไทย', 'D''angelo_ไทย', 'Berge_ไทย', 'Mrs.', 'D''angelo', 'Berge', 'Orthopedist', 'Sibyl.Kuhlman51', '$2b$10$BFJZq5aGkvthNBY4dbyjau7Ihokkv.MJRU8deVzeHiYE9TfnsS3Nm', 'https://cloudflare-ipfs.com/ipfs/Qmd3W5DuhgHirLHGVixi6V76LhCkZUz6pnFt5AJBiyvHye/avatar/727.jpg', '2022-09-07 10:52:46.89', '2022-09-07 10:52:46.891');
INSERT INTO public."Doctor" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, "position", username, password, "profilePicURL", "createdAt", "updatedAt") VALUES (12, 'Mrs._ไทย', 'Seth_ไทย', 'Sauer_ไทย', 'Mrs.', 'Seth', 'Sauer', 'Cardiologist', 'Zoila.Kemmer', '$2b$10$HSGimnSpsm6LauNmslLsBulWHz9IoSV16AmoxKsHkvYtaMRYN8iuO', 'https://cloudflare-ipfs.com/ipfs/Qmd3W5DuhgHirLHGVixi6V76LhCkZUz6pnFt5AJBiyvHye/avatar/5.jpg', '2022-09-07 10:52:46.89', '2022-09-07 10:52:46.891');
INSERT INTO public."Doctor" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, "position", username, password, "profilePicURL", "createdAt", "updatedAt") VALUES (13, 'Mr._ไทย', 'Jovan_ไทย', 'Feeney_ไทย', 'Mr.', 'Jovan', 'Feeney', 'Neurologist', 'Brook_Turner63', '$2b$10$fvYXFLnH2mt5eGpDq2YOWO5eKBcFROWe3QTGfYI.Rd.sirIS8Mw4i', 'https://cloudflare-ipfs.com/ipfs/Qmd3W5DuhgHirLHGVixi6V76LhCkZUz6pnFt5AJBiyvHye/avatar/629.jpg', '2022-09-07 10:52:46.89', '2022-09-07 10:52:46.891');
INSERT INTO public."Doctor" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, "position", username, password, "profilePicURL", "createdAt", "updatedAt") VALUES (14, 'Dr._ไทย', 'Sylvia_ไทย', 'Corwin_ไทย', 'Dr.', 'Sylvia', 'Corwin', 'Cardiologist', 'Rowland.Veum25', '$2b$10$O.6KGsDd2iE769wxL/G6q.kFt0kbIWIkogp85wr8orq1IP67FX7Gy', 'https://cloudflare-ipfs.com/ipfs/Qmd3W5DuhgHirLHGVixi6V76LhCkZUz6pnFt5AJBiyvHye/avatar/652.jpg', '2022-09-07 10:52:46.89', '2022-09-07 10:52:46.891');
INSERT INTO public."Doctor" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, "position", username, password, "profilePicURL", "createdAt", "updatedAt") VALUES (15, 'Dr._ไทย', 'Remington_ไทย', 'Schroeder_ไทย', 'Dr.', 'Remington', 'Schroeder', 'Orthopedist', 'Berta.Gusikowski76', '$2b$10$OC/Lm5fQlRK.otq6TRkgRujCNV8htxMbimRH76kgMF7azjTgNwPN6', 'https://cloudflare-ipfs.com/ipfs/Qmd3W5DuhgHirLHGVixi6V76LhCkZUz6pnFt5AJBiyvHye/avatar/1224.jpg', '2022-09-07 10:52:46.89', '2022-09-07 10:52:46.891');
INSERT INTO public."Doctor" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, "position", username, password, "profilePicURL", "createdAt", "updatedAt") VALUES (16, 'Miss_ไทย', 'Jody_ไทย', 'Botsford_ไทย', 'Miss', 'Jody', 'Botsford', 'Psychiatrist', 'Ashtyn.Rath0', '$2b$10$w02EKd.1Y2BpvgvHODoW3OrsVBBLbob4tACTW6PYVfSq2Exg8HwVi', 'https://cloudflare-ipfs.com/ipfs/Qmd3W5DuhgHirLHGVixi6V76LhCkZUz6pnFt5AJBiyvHye/avatar/937.jpg', '2022-09-07 10:52:46.89', '2022-09-07 10:52:46.891');
INSERT INTO public."Doctor" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, "position", username, password, "profilePicURL", "createdAt", "updatedAt") VALUES (17, 'Ms._ไทย', 'Lawrence_ไทย', 'Murphy_ไทย', 'Ms.', 'Lawrence', 'Murphy', 'Orthopedist', 'Daisy.Larson6', '$2b$10$w.IY6k2.w5ToXEImzDBBMu4TxpvL1.7mQPFvrDUZjQ9pyy0CG2szy', 'https://cloudflare-ipfs.com/ipfs/Qmd3W5DuhgHirLHGVixi6V76LhCkZUz6pnFt5AJBiyvHye/avatar/820.jpg', '2022-09-07 10:52:46.89', '2022-09-07 10:52:46.891');
INSERT INTO public."Doctor" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, "position", username, password, "profilePicURL", "createdAt", "updatedAt") VALUES (18, 'Mrs._ไทย', 'Prudence_ไทย', 'McCullough_ไทย', 'Mrs.', 'Prudence', 'McCullough', 'Cardiologist', 'Rachelle50', '$2b$10$0RnLXLvcQH97JOGb5u6OV.jMMVrQ2bchjecpDKt9hUTgU1TqjpSTa', 'https://cloudflare-ipfs.com/ipfs/Qmd3W5DuhgHirLHGVixi6V76LhCkZUz6pnFt5AJBiyvHye/avatar/240.jpg', '2022-09-07 10:52:46.89', '2022-09-07 10:52:46.891');
INSERT INTO public."Doctor" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, "position", username, password, "profilePicURL", "createdAt", "updatedAt") VALUES (19, 'Miss_ไทย', 'Darren_ไทย', 'Hamill_ไทย', 'Miss', 'Darren', 'Hamill', 'Psychiatrist', 'Jarred_Schmidt', '$2b$10$sitkcUIYXy/0q6etftpWseb7dVVocEj6mw8bj/MM3kYXd/NkxRvQe', 'https://cloudflare-ipfs.com/ipfs/Qmd3W5DuhgHirLHGVixi6V76LhCkZUz6pnFt5AJBiyvHye/avatar/622.jpg', '2022-09-07 10:52:46.89', '2022-09-07 10:52:46.891');
INSERT INTO public."Doctor" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, "position", username, password, "profilePicURL", "createdAt", "updatedAt") VALUES (20, 'Miss_ไทย', 'Wade_ไทย', 'Tremblay_ไทย', 'Miss', 'Wade', 'Tremblay', 'Neurologist', 'Gia.Gutmann', '$2b$10$XLL.IfKROXl3oHn4qmd0a.4Ws2Tif3T2XKQHuhA2q8jtrvf/qrGOG', 'https://cloudflare-ipfs.com/ipfs/Qmd3W5DuhgHirLHGVixi6V76LhCkZUz6pnFt5AJBiyvHye/avatar/235.jpg', '2022-09-07 10:52:46.89', '2022-09-07 10:52:46.891');
INSERT INTO public."Doctor" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, "position", username, password, "profilePicURL", "createdAt", "updatedAt") VALUES (21, 'Miss_ไทย', 'Yasmine_ไทย', 'Bogisich_ไทย', 'Miss', 'Yasmine', 'Bogisich', 'Orthopedist', 'Mervin51', '$2b$10$vF0Ee4w3HaM2TRA8295pt.70v/GKO3.AmS21fweWX9hVdfLKMdmlm', 'https://cloudflare-ipfs.com/ipfs/Qmd3W5DuhgHirLHGVixi6V76LhCkZUz6pnFt5AJBiyvHye/avatar/577.jpg', '2022-09-07 10:52:46.89', '2022-09-07 10:52:46.891');
INSERT INTO public."Doctor" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, "position", username, password, "profilePicURL", "createdAt", "updatedAt") VALUES (22, 'Ms._ไทย', 'Cristobal_ไทย', 'Fadel_ไทย', 'Ms.', 'Cristobal', 'Fadel', 'Cardiologist', 'Amalia69', '$2b$10$f1wk6/46OngITAb0xDDtgeDyloGBpnC47BIAU0wD8eRuFxogFDGHu', 'https://cloudflare-ipfs.com/ipfs/Qmd3W5DuhgHirLHGVixi6V76LhCkZUz6pnFt5AJBiyvHye/avatar/1206.jpg', '2022-09-07 10:52:46.89', '2022-09-07 10:52:46.891');
INSERT INTO public."Doctor" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, "position", username, password, "profilePicURL", "createdAt", "updatedAt") VALUES (23, 'Mr._ไทย', 'Amara_ไทย', 'Roberts_ไทย', 'Mr.', 'Amara', 'Roberts', 'Cardiologist', 'Pedro2', '$2b$10$HLQfhRui9VwCkoPICP56YuMy86YSiGUEt4N76UHtgs5nChb7pElmu', 'https://cloudflare-ipfs.com/ipfs/Qmd3W5DuhgHirLHGVixi6V76LhCkZUz6pnFt5AJBiyvHye/avatar/604.jpg', '2022-09-07 10:52:46.89', '2022-09-07 10:52:46.891');
INSERT INTO public."Doctor" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, "position", username, password, "profilePicURL", "createdAt", "updatedAt") VALUES (24, 'Mrs._ไทย', 'Camille_ไทย', 'Gaylord_ไทย', 'Mrs.', 'Camille', 'Gaylord', 'Neurologist', 'Hayley13', '$2b$10$UHSSgMbuhKp6BVhD5Kqgtun92BWSzlzIXrws/fMMvqWvARkKuhPn.', 'https://cloudflare-ipfs.com/ipfs/Qmd3W5DuhgHirLHGVixi6V76LhCkZUz6pnFt5AJBiyvHye/avatar/1245.jpg', '2022-09-07 10:52:46.89', '2022-09-07 10:52:46.891');
INSERT INTO public."Doctor" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, "position", username, password, "profilePicURL", "createdAt", "updatedAt") VALUES (25, 'Miss_ไทย', 'Bryce_ไทย', 'O''Keefe_ไทย', 'Miss', 'Bryce', 'O''Keefe', 'Podiatrist', 'Itzel74', '$2b$10$F100RgEI0DBR5lUbTcmeC.HWh478fh7IcKZn.6ySxT5Oh84kIkf2i', 'https://cloudflare-ipfs.com/ipfs/Qmd3W5DuhgHirLHGVixi6V76LhCkZUz6pnFt5AJBiyvHye/avatar/779.jpg', '2022-09-07 10:52:46.89', '2022-09-07 10:52:46.891');
INSERT INTO public."Doctor" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, "position", username, password, "profilePicURL", "createdAt", "updatedAt") VALUES (26, 'Miss_ไทย', 'Cary_ไทย', 'Dickens_ไทย', 'Miss', 'Cary', 'Dickens', 'Neurologist', 'Adolfo_Corwin', '$2b$10$hWKGlhk/KLyMOwGAsGwb5.yEWgff2RKkqY0P7tU6FTYMzXcTPYVtm', 'https://cloudflare-ipfs.com/ipfs/Qmd3W5DuhgHirLHGVixi6V76LhCkZUz6pnFt5AJBiyvHye/avatar/182.jpg', '2022-09-07 10:52:46.89', '2022-09-07 10:52:46.891');
INSERT INTO public."Doctor" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, "position", username, password, "profilePicURL", "createdAt", "updatedAt") VALUES (27, 'Ms._ไทย', 'Maeve_ไทย', 'Streich_ไทย', 'Ms.', 'Maeve', 'Streich', 'Neurologist', 'Piper76', '$2b$10$F5qua9dtvZo2ZhsyGoJ.Yulv1DUFOD.y.8Hu9QphpRkmKJUqLtdWq', 'https://cloudflare-ipfs.com/ipfs/Qmd3W5DuhgHirLHGVixi6V76LhCkZUz6pnFt5AJBiyvHye/avatar/626.jpg', '2022-09-07 10:52:46.89', '2022-09-07 10:52:46.891');
INSERT INTO public."Doctor" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, "position", username, password, "profilePicURL", "createdAt", "updatedAt") VALUES (28, 'Mrs._ไทย', 'Isaac_ไทย', 'Treutel_ไทย', 'Mrs.', 'Isaac', 'Treutel', 'Psychiatrist', 'Yasmine_Kub18', '$2b$10$d.kS92ZHSwRzhu8Uw4eCu.ap8CTRYNV6kqfpF1ogC8AaX00JPklVi', 'https://cloudflare-ipfs.com/ipfs/Qmd3W5DuhgHirLHGVixi6V76LhCkZUz6pnFt5AJBiyvHye/avatar/336.jpg', '2022-09-07 10:52:46.89', '2022-09-07 10:52:46.891');
INSERT INTO public."Doctor" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, "position", username, password, "profilePicURL", "createdAt", "updatedAt") VALUES (29, 'Mr._ไทย', 'Shyann_ไทย', 'Mitchell_ไทย', 'Mr.', 'Shyann', 'Mitchell', 'Podiatrist', 'Jose6', '$2b$10$hUoxOBBv8DMW11a.Egjkt.g4YPIevBu6Z8nJ7ANZzRSOYgb6beAyi', 'https://cloudflare-ipfs.com/ipfs/Qmd3W5DuhgHirLHGVixi6V76LhCkZUz6pnFt5AJBiyvHye/avatar/458.jpg', '2022-09-07 10:52:46.89', '2022-09-07 10:52:46.891');
INSERT INTO public."Doctor" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, "position", username, password, "profilePicURL", "createdAt", "updatedAt") VALUES (30, 'Mr._ไทย', 'Stefan_ไทย', 'Roberts_ไทย', 'Mr.', 'Stefan', 'Roberts', 'Orthopedist', 'Anthony23', '$2b$10$L8Ty/8h4lNj4s12TEdNCyO2faA2Rct3BwPhJ2FX.EG8JQEYEM.wB6', 'https://cloudflare-ipfs.com/ipfs/Qmd3W5DuhgHirLHGVixi6V76LhCkZUz6pnFt5AJBiyvHye/avatar/745.jpg', '2022-09-07 10:52:46.89', '2022-09-07 10:52:46.891');


--
-- Data for Name: Invoice; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public."Invoice" (id, "appointmentId", paid, total, "createdAt", "updatedAt") VALUES (1, 2, true, 3922590, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."Invoice" (id, "appointmentId", paid, total, "createdAt", "updatedAt") VALUES (2, 4, false, 2748108, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."Invoice" (id, "appointmentId", paid, total, "createdAt", "updatedAt") VALUES (3, 9, true, 1049944, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."Invoice" (id, "appointmentId", paid, total, "createdAt", "updatedAt") VALUES (4, 11, false, 2434877, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."Invoice" (id, "appointmentId", paid, total, "createdAt", "updatedAt") VALUES (5, 23, true, 1184758, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."Invoice" (id, "appointmentId", paid, total, "createdAt", "updatedAt") VALUES (6, 28, false, 2992140, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Invoice" (id, "appointmentId", paid, total, "createdAt", "updatedAt") VALUES (7, 37, true, 401219, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Invoice" (id, "appointmentId", paid, total, "createdAt", "updatedAt") VALUES (8, 31, true, 442596, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.951');
INSERT INTO public."Invoice" (id, "appointmentId", paid, total, "createdAt", "updatedAt") VALUES (9, 38, true, 65895, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Invoice" (id, "appointmentId", paid, total, "createdAt", "updatedAt") VALUES (10, 39, true, 461052, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Invoice" (id, "appointmentId", paid, total, "createdAt", "updatedAt") VALUES (11, 45, true, 1945318, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Invoice" (id, "appointmentId", paid, total, "createdAt", "updatedAt") VALUES (12, 50, true, 251713, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.952');
INSERT INTO public."Invoice" (id, "appointmentId", paid, total, "createdAt", "updatedAt") VALUES (13, 53, false, 2875989, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Invoice" (id, "appointmentId", paid, total, "createdAt", "updatedAt") VALUES (14, 49, false, 2114117, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Invoice" (id, "appointmentId", paid, total, "createdAt", "updatedAt") VALUES (15, 56, false, 2110424, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Invoice" (id, "appointmentId", paid, total, "createdAt", "updatedAt") VALUES (16, 58, true, 449162, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Invoice" (id, "appointmentId", paid, total, "createdAt", "updatedAt") VALUES (17, 59, false, 626257, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Invoice" (id, "appointmentId", paid, total, "createdAt", "updatedAt") VALUES (18, 63, false, 887464, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Invoice" (id, "appointmentId", paid, total, "createdAt", "updatedAt") VALUES (19, 64, false, 1263247, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Invoice" (id, "appointmentId", paid, total, "createdAt", "updatedAt") VALUES (20, 70, false, 435635, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Invoice" (id, "appointmentId", paid, total, "createdAt", "updatedAt") VALUES (21, 73, true, 2086461, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Invoice" (id, "appointmentId", paid, total, "createdAt", "updatedAt") VALUES (22, 84, true, 1776716, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Invoice" (id, "appointmentId", paid, total, "createdAt", "updatedAt") VALUES (23, 85, true, 851947, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Invoice" (id, "appointmentId", paid, total, "createdAt", "updatedAt") VALUES (24, 89, false, 1049953, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Invoice" (id, "appointmentId", paid, total, "createdAt", "updatedAt") VALUES (25, 93, false, 1470436, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Invoice" (id, "appointmentId", paid, total, "createdAt", "updatedAt") VALUES (26, 98, false, 1420466, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Invoice" (id, "appointmentId", paid, total, "createdAt", "updatedAt") VALUES (28, 99, false, 1666940, '2022-09-07 10:52:46.954', '2022-09-07 10:52:46.954');
INSERT INTO public."Invoice" (id, "appointmentId", paid, total, "createdAt", "updatedAt") VALUES (27, 96, true, 3436462, '2022-09-07 10:52:46.954', '2022-09-07 10:52:46.955');
INSERT INTO public."Invoice" (id, "appointmentId", paid, total, "createdAt", "updatedAt") VALUES (29, 97, true, 1061326, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.953');


--
-- Data for Name: InvoiceDiscount; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public."InvoiceDiscount" (id, name, amount, "invoiceId", "createdAt", "updatedAt") VALUES (1, 'Social Security', 50000, 2, '2022-09-27 08:50:21.102', '2022-09-27 08:50:21.102');


--
-- Data for Name: InvoiceItem; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (1, 'Incredible Frozen Sausages', 21272, 5, 1, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (2, 'Fantastic Cotton Tuna', 47997, 20, 1, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (3, 'Unbranded Fresh Computer', 47002, 6, 1, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (4, 'Unbranded Soft Bacon', 33498, 14, 1, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (5, 'Electronic Fresh Tuna', 29655, 10, 1, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (6, 'Handmade Fresh Fish', 49024, 12, 1, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (7, 'Handmade Steel Tuna', 7885, 15, 1, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (8, 'Refined Concrete Fish', 25966, 13, 1, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (9, 'Fantastic Granite Bike', 28655, 13, 1, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (10, 'Small Fresh Soap', 49015, 8, 1, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (11, 'Generic Bronze Shirt', 32446, 16, 2, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (12, 'Oriental Plastic Shirt', 38438, 15, 2, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (13, 'Refined Metal Pants', 28038, 11, 2, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (14, 'Unbranded Soft Gloves', 15495, 17, 2, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (15, 'Oriental Wooden Keyboard', 22103, 2, 2, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (16, 'Generic Granite Keyboard', 21707, 2, 2, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (17, 'Awesome Plastic Chair', 37873, 3, 2, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (18, 'Electronic Steel Mouse', 38720, 18, 2, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (19, 'Elegant Wooden Chair', 18237, 10, 2, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (20, 'Elegant Bronze Chair', 22402, 16, 3, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (21, 'Luxurious Cotton Computer', 14944, 15, 3, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (22, 'Oriental Steel Chicken', 18990, 10, 3, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (23, 'Elegant Soft Hat', 30828, 9, 3, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (24, 'Luxurious Concrete Pants', 34007, 10, 4, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (25, 'Fantastic Granite Car', 24742, 5, 4, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (26, 'Gorgeous Frozen Pizza', 38647, 17, 4, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (27, 'Ergonomic Bronze Bike', 45781, 11, 4, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (28, 'Refined Frozen Towels', 35364, 8, 4, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (29, 'Refined Cotton Shirt', 47096, 8, 4, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (30, 'Unbranded Metal Table', 8182, 13, 4, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (31, 'Elegant Rubber Sausages', 2701, 1, 4, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (32, 'Generic Concrete Pizza', 2784, 15, 4, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (33, 'Luxurious Granite Bacon', 26159, 20, 5, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (34, 'Electronic Frozen Shirt', 35247, 6, 5, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (35, 'Electronic Fresh Salad', 37508, 12, 5, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (36, 'Rustic Granite Shoes', 40506, 3, 6, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (37, 'Incredible Plastic Sausages', 46835, 7, 6, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (38, 'Electronic Frozen Cheese', 35621, 19, 6, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (39, 'Bespoke Plastic Chips', 12339, 5, 6, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (40, 'Recycled Soft Keyboard', 15781, 19, 6, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (41, 'Fantastic Wooden Ball', 26472, 6, 6, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (42, 'Refined Wooden Shirt', 23402, 20, 6, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (43, 'Awesome Granite Keyboard', 18468, 13, 6, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (44, 'Modern Wooden Pants', 33552, 19, 6, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (45, 'Incredible Concrete Ball', 19867, 20, 7, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (46, 'Small Cotton Tuna', 431, 9, 7, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (47, 'Practical Cotton Chips', 47571, 1, 9, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (48, 'Handcrafted Bronze Tuna', 9162, 2, 9, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (49, 'Handmade Soft Computer', 13421, 3, 8, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.951');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (50, 'Electronic Plastic Hat', 35677, 10, 8, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.951');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (51, 'Refined Cotton Cheese', 45563, 1, 8, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.951');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (52, 'Recycled Wooden Tuna', 25614, 18, 10, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (53, 'Luxurious Cotton Sausages', 38618, 6, 11, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (54, 'Elegant Metal Fish', 46307, 17, 11, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (55, 'Electronic Steel Pizza', 42861, 5, 11, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (56, 'Modern Bronze Salad', 31734, 5, 11, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (57, 'Licensed Fresh Pants', 46118, 12, 11, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (58, 'Tasty Granite Computer', 10675, 4, 13, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (59, 'Tasty Plastic Chips', 48482, 5, 13, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (60, 'Electronic Rubber Bacon', 15950, 18, 13, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (61, 'Ergonomic Soft Cheese', 24561, 13, 13, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (62, 'Incredible Fresh Tuna', 18780, 17, 13, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (63, 'Licensed Granite Table', 48365, 18, 13, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (64, 'Gorgeous Granite Keyboard', 32699, 6, 13, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (65, 'Unbranded Rubber Soap', 13118, 10, 13, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (66, 'Handmade Plastic Bike', 12626, 9, 13, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (67, 'Luxurious Steel Shirt', 44206, 8, 13, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (68, 'Awesome Soft Pizza', 23870, 10, 12, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.952');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (69, 'Elegant Rubber Cheese', 1859, 7, 12, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.952');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (70, 'Handmade Frozen Sausages', 19593, 15, 14, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (71, 'Rustic Soft Pants', 35248, 20, 14, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (72, 'Unbranded Granite Bacon', 43862, 12, 14, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (73, 'Gorgeous Fresh Gloves', 19246, 1, 14, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (74, 'Elegant Concrete Car', 40950, 2, 14, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (75, 'Rustic Concrete Keyboard', 25372, 11, 14, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (76, 'Refined Soft Mouse', 17390, 12, 14, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (77, 'Gorgeous Rubber Gloves', 37376, 19, 15, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (78, 'Generic Metal Computer', 31807, 3, 15, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (79, 'Fantastic Fresh Shirt', 3568, 18, 15, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (80, 'Rustic Metal Soap', 28717, 2, 15, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (81, 'Elegant Rubber Keyboard', 5849, 2, 15, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (82, 'Handcrafted Granite Chair', 34516, 12, 15, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (83, 'Modern Cotton Salad', 5638, 2, 15, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (84, 'Recycled Metal Ball', 25974, 4, 15, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (85, 'Awesome Granite Gloves', 36113, 17, 15, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (86, 'Modern Wooden Bike', 14109, 2, 15, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (87, 'Generic Cotton Shoes', 6760, 4, 17, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (88, 'Electronic Wooden Pizza', 20439, 8, 17, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (89, 'Fantastic Fresh Pizza', 13612, 13, 17, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (90, 'Rustic Frozen Fish', 17695, 8, 17, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (93, 'Sleek Frozen Gloves', 37566, 9, 18, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (94, 'Ergonomic Wooden Chicken', 4336, 15, 18, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (95, 'Generic Cotton Soap', 44030, 11, 18, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (92, 'Electronic Frozen Pants', 37241, 7, 16, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (96, 'Refined Rubber Bacon', 26925, 7, 16, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (131, 'Modern Bronze Table', 40358, 7, 28, '2022-09-07 10:52:46.954', '2022-09-07 10:52:46.954');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (132, 'Elegant Rubber Pants', 25932, 13, 28, '2022-09-07 10:52:46.954', '2022-09-07 10:52:46.954');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (133, 'Gorgeous Granite Keyboard', 38315, 10, 28, '2022-09-07 10:52:46.954', '2022-09-07 10:52:46.954');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (134, 'Bespoke Fresh Towels', 9012, 10, 28, '2022-09-07 10:52:46.954', '2022-09-07 10:52:46.954');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (135, 'Awesome Cotton Salad', 12560, 12, 28, '2022-09-07 10:52:46.954', '2022-09-07 10:52:46.954');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (136, 'Elegant Granite Pants', 21661, 9, 28, '2022-09-07 10:52:46.954', '2022-09-07 10:52:46.954');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (137, 'Practical Bronze Chicken', 28565, 7, 28, '2022-09-07 10:52:46.954', '2022-09-07 10:52:46.954');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (138, 'Electronic Soft Tuna', 7106, 4, 28, '2022-09-07 10:52:46.954', '2022-09-07 10:52:46.954');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (97, 'Intelligent Fresh Table', 33481, 6, 19, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (98, 'Gorgeous Rubber Keyboard', 49917, 2, 19, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (99, 'Practical Soft Shoes', 6845, 17, 19, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (100, 'Refined Concrete Gloves', 14592, 5, 19, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (101, 'Practical Soft Bacon', 25555, 17, 19, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (102, 'Unbranded Concrete Gloves', 23785, 1, 19, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (103, 'Gorgeous Metal Gloves', 12350, 12, 19, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (104, 'Handcrafted Plastic Bike', 27797, 6, 19, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (105, 'Gorgeous Wooden Sausages', 48363, 5, 20, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (106, 'Practical Soft Tuna', 9691, 20, 20, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (107, 'Handmade Bronze Shoes', 17647, 5, 21, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (108, 'Unbranded Soft Shirt', 42298, 19, 21, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (109, 'Handcrafted Wooden Keyboard', 909, 8, 21, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (110, 'Ergonomic Concrete Mouse', 3084, 17, 21, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (111, 'Intelligent Fresh Mouse', 19356, 10, 21, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (112, 'Unbranded Frozen Towels', 29362, 12, 21, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (113, 'Electronic Cotton Table', 45660, 4, 21, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (114, 'Refined Steel Soap', 27088, 15, 21, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (126, 'Awesome Wooden Hat', 5703, 5, 24, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (127, 'Unbranded Frozen Car', 554, 2, 24, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (128, 'Licensed Frozen Shirt', 10294, 6, 24, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (129, 'Luxurious Granite Table', 20469, 10, 24, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (130, 'Sleek Granite Soap', 41882, 18, 24, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (115, 'Handcrafted Rubber Towels', 42402, 15, 22, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (116, 'Gorgeous Metal Pizza', 14723, 20, 22, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (117, 'Handmade Frozen Tuna', 49778, 17, 22, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (118, 'Elegant Granite Sausages', 34415, 5, 23, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (119, 'Small Frozen Sausages', 18587, 16, 23, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (120, 'Sleek Frozen Chips', 19124, 20, 23, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (140, 'Rustic Concrete Shoes', 6042, 18, 26, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (141, 'Ergonomic Rubber Bacon', 24373, 10, 26, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (142, 'Practical Frozen Cheese', 41948, 7, 26, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (143, 'Oriental Rubber Shoes', 22341, 15, 26, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (144, 'Electronic Metal Bike', 31368, 3, 26, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (145, 'Tasty Cotton Bacon', 41695, 5, 26, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (146, 'Generic Concrete Computer', 22775, 6, 26, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (139, 'Unbranded Rubber Bacon', 25332, 4, 27, '2022-09-07 10:52:46.954', '2022-09-07 10:52:46.955');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (147, 'Handcrafted Steel Shoes', 40845, 20, 27, '2022-09-07 10:52:46.954', '2022-09-07 10:52:46.955');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (148, 'Generic Wooden Table', 2499, 16, 27, '2022-09-07 10:52:46.954', '2022-09-07 10:52:46.955');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (149, 'Recycled Plastic Towels', 10050, 3, 27, '2022-09-07 10:52:46.954', '2022-09-07 10:52:46.955');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (150, 'Gorgeous Steel Bacon', 40714, 18, 27, '2022-09-07 10:52:46.954', '2022-09-07 10:52:46.955');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (151, 'Gorgeous Bronze Cheese', 2880, 11, 27, '2022-09-07 10:52:46.954', '2022-09-07 10:52:46.955');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (152, 'Unbranded Wooden Bike', 12499, 13, 27, '2022-09-07 10:52:46.954', '2022-09-07 10:52:46.955');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (153, 'Refined Steel Shirt', 44717, 19, 27, '2022-09-07 10:52:46.954', '2022-09-07 10:52:46.955');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (154, 'Unbranded Concrete Hat', 42644, 11, 27, '2022-09-07 10:52:46.954', '2022-09-07 10:52:46.955');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (155, 'Recycled Metal Tuna', 11243, 18, 27, '2022-09-07 10:52:46.954', '2022-09-07 10:52:46.955');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (156, 'Sleek Rubber Ball', 6760, 13, 29, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.953');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (157, 'Intelligent Metal Chicken', 37112, 13, 29, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.953');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (158, 'Unbranded Rubber Towels', 49099, 10, 29, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.953');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (91, 'Handmade Concrete Bacon', 13021, 9, 17, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (121, 'Licensed Bronze Bike', 45807, 9, 25, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (122, 'Luxurious Fresh Chair', 39041, 17, 25, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (123, 'Intelligent Fresh Car', 16588, 4, 25, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (124, 'Handcrafted Fresh Keyboard', 518, 16, 25, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."InvoiceItem" (id, name, price, quantity, "invoiceId", "createdAt", "updatedAt") VALUES (125, 'Fantastic Concrete Pants', 29076, 11, 25, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');


--
-- Data for Name: Medicine; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public."Medicine" (id, name, description, "createdAt", "updatedAt", "pictureURL") VALUES (1, 'nostrum', 'Illo et aut. Dolores ut sed at fuga dolorem. Sunt corporis molestias doloribus quod nisi. Qui quibusdam doloremque quia labore ea reprehenderit aut. Repellat unde ut saepe et dolorem eveniet molestiae. Et quam pariatur architecto architecto eum quidem saepe officia.', '2022-09-07 10:52:45.188', '2022-09-07 10:52:45.189', 'https://hips.hearstapps.com/hmg-prod.s3.amazonaws.com/images/many-colorful-pills-royalty-free-image-504370555-1542820898.jpg?crop=0.670xw:1.00xh;0,0&resize=480:*');
INSERT INTO public."Medicine" (id, name, description, "createdAt", "updatedAt", "pictureURL") VALUES (2, 'qui', 'Iusto rem officiis occaecati ut quis quibusdam. Soluta saepe eos voluptatum ut nesciunt excepturi. Aperiam eveniet sed sint amet voluptatem dolorum sint. Laudantium alias qui architecto possimus aliquid aut et. Consequatur aut tempora tempore vel. Et accusamus suscipit error autem blanditiis voluptas.', '2022-09-07 10:52:45.188', '2022-09-07 10:52:45.189', 'https://hips.hearstapps.com/hmg-prod.s3.amazonaws.com/images/many-colorful-pills-royalty-free-image-504370555-1542820898.jpg?crop=0.670xw:1.00xh;0,0&resize=480:*');
INSERT INTO public."Medicine" (id, name, description, "createdAt", "updatedAt", "pictureURL") VALUES (3, 'incidunt', 'Et commodi animi velit et nihil incidunt et. Impedit voluptate deleniti nobis quaerat ipsa in ex deserunt. Labore eligendi numquam veritatis. Ut ut ad voluptas aut aut at est.', '2022-09-07 10:52:45.188', '2022-09-07 10:52:45.189', 'https://hips.hearstapps.com/hmg-prod.s3.amazonaws.com/images/many-colorful-pills-royalty-free-image-504370555-1542820898.jpg?crop=0.670xw:1.00xh;0,0&resize=480:*');
INSERT INTO public."Medicine" (id, name, description, "createdAt", "updatedAt", "pictureURL") VALUES (4, 'necessitatibus', 'Cum illum aperiam autem dignissimos harum fugit in unde amet. Magni maxime eius harum quia unde facere voluptatem. Maxime nulla provident rerum praesentium ratione eum ipsum. Minus ea occaecati et tempore.', '2022-09-07 10:52:45.188', '2022-09-07 10:52:45.189', 'https://hips.hearstapps.com/hmg-prod.s3.amazonaws.com/images/many-colorful-pills-royalty-free-image-504370555-1542820898.jpg?crop=0.670xw:1.00xh;0,0&resize=480:*');
INSERT INTO public."Medicine" (id, name, description, "createdAt", "updatedAt", "pictureURL") VALUES (5, 'quisquam', 'Quas nemo et non molestiae cumque ipsa aut accusamus quidem. Quia amet voluptatem. Voluptatem id nihil excepturi autem fugit velit ut magnam deleniti. Cumque corrupti sint similique est velit eveniet.', '2022-09-07 10:52:45.188', '2022-09-07 10:52:45.189', 'https://hips.hearstapps.com/hmg-prod.s3.amazonaws.com/images/many-colorful-pills-royalty-free-image-504370555-1542820898.jpg?crop=0.670xw:1.00xh;0,0&resize=480:*');
INSERT INTO public."Medicine" (id, name, description, "createdAt", "updatedAt", "pictureURL") VALUES (6, 'atque', 'Consequatur ipsum labore voluptas commodi harum culpa et dolorem et. Facere voluptas sed. Id voluptatem sint. Nulla deserunt consequatur tenetur quibusdam sint pariatur esse. Et omnis tempora laborum.', '2022-09-07 10:52:45.188', '2022-09-07 10:52:45.189', 'https://hips.hearstapps.com/hmg-prod.s3.amazonaws.com/images/many-colorful-pills-royalty-free-image-504370555-1542820898.jpg?crop=0.670xw:1.00xh;0,0&resize=480:*');
INSERT INTO public."Medicine" (id, name, description, "createdAt", "updatedAt", "pictureURL") VALUES (7, 'enim', 'Neque minus reprehenderit illum voluptatibus esse unde voluptatem. Cupiditate asperiores totam molestiae aspernatur. Nam ut at explicabo ratione nihil ut. Rem voluptatem dolore placeat pariatur ut sequi ea ut quas. Amet ut dolores nulla quaerat officia suscipit.', '2022-09-07 10:52:45.188', '2022-09-07 10:52:45.189', 'https://hips.hearstapps.com/hmg-prod.s3.amazonaws.com/images/many-colorful-pills-royalty-free-image-504370555-1542820898.jpg?crop=0.670xw:1.00xh;0,0&resize=480:*');
INSERT INTO public."Medicine" (id, name, description, "createdAt", "updatedAt", "pictureURL") VALUES (8, 'facilis', 'Eum nihil cum omnis dicta sunt at at voluptatum. Nisi non magni asperiores iure eos esse blanditiis adipisci explicabo. Quisquam debitis eum et at. Accusantium sunt eius ea vero odio fugit at voluptatem. Autem omnis rerum sed corporis animi. Sit qui voluptatem nihil aut perferendis consequatur et beatae.', '2022-09-07 10:52:45.188', '2022-09-07 10:52:45.189', 'https://hips.hearstapps.com/hmg-prod.s3.amazonaws.com/images/many-colorful-pills-royalty-free-image-504370555-1542820898.jpg?crop=0.670xw:1.00xh;0,0&resize=480:*');
INSERT INTO public."Medicine" (id, name, description, "createdAt", "updatedAt", "pictureURL") VALUES (9, 'cum', 'Harum quia tempora dolores nisi est. Similique modi quae at voluptate voluptas velit. Asperiores ut et corrupti dolores sed omnis. Vel ut aut voluptatem fugit id excepturi sit. Est necessitatibus sed quia eos quia quis facere minima. Aut est quis facere eos eos.', '2022-09-07 10:52:45.188', '2022-09-07 10:52:45.189', 'https://hips.hearstapps.com/hmg-prod.s3.amazonaws.com/images/many-colorful-pills-royalty-free-image-504370555-1542820898.jpg?crop=0.670xw:1.00xh;0,0&resize=480:*');
INSERT INTO public."Medicine" (id, name, description, "createdAt", "updatedAt", "pictureURL") VALUES (10, 'eius', 'Vero eveniet ea est ut qui iste. Repudiandae atque totam ipsam occaecati exercitationem omnis ut facilis. Et perferendis ut et dolor voluptas labore. Culpa cum dicta asperiores ipsum dicta libero cum. Modi voluptates molestiae quia repellat qui.', '2022-09-07 10:52:45.188', '2022-09-07 10:52:45.189', 'https://hips.hearstapps.com/hmg-prod.s3.amazonaws.com/images/many-colorful-pills-royalty-free-image-504370555-1542820898.jpg?crop=0.670xw:1.00xh;0,0&resize=480:*');


--
-- Data for Name: Patient; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-328759', 'Miss_ไทย', 'Joshuah_ไทย', 'Connelly_ไทย', 'Miss', 'Joshuah', 'Connelly', 'British', NULL, 'YV919588', '0008225725', 121, 152, '1967-10-07 14:52:14.485', 'B', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-370705', 'Ms._ไทย', 'Tillman_ไทย', 'Ratke_ไทย', 'Ms.', 'Tillman', 'Ratke', 'Lao', NULL, 'AT674281', '0489048801', 135, 183, '1996-08-07 08:06:07.73', 'O', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-779528', 'Miss_ไทย', 'Zola_ไทย', 'Lind_ไทย', 'Miss', 'Zola', 'Lind', 'American', NULL, 'FI779720', '0930094946', 47, 196, '1989-10-09 14:36:35.622', 'B', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-162462', 'Dr._ไทย', 'Leora_ไทย', 'Kessler_ไทย', 'Dr.', 'Leora', 'Kessler', 'American', NULL, 'JN848321', '0110165178', 98, 197, '1964-10-29 18:35:01.23', 'B', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-740755', 'Mrs._ไทย', 'Dorian_ไทย', 'Mills_ไทย', 'Mrs.', 'Dorian', 'Mills', 'Lao', NULL, 'QZ968648', '0548324590', 75, 185, '1944-11-25 04:18:14.082', 'AB', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-963617', 'Dr._ไทย', 'Frida_ไทย', 'Quitzon_ไทย', 'Dr.', 'Frida', 'Quitzon', 'British', NULL, 'QS526935', '0424822267', 148, 155, '1953-08-13 09:48:33.389', 'A', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-443354', 'Dr._ไทย', 'Arlie_ไทย', 'Harris_ไทย', 'Dr.', 'Arlie', 'Harris', 'Finnish', NULL, 'YH383117', '0818312547', 120, 175, '1999-09-02 01:58:43.907', 'A', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-562380', 'Miss_ไทย', 'Candelario_ไทย', 'Kris_ไทย', 'Miss', 'Candelario', 'Kris', 'Lao', NULL, 'UP122527', '0501389238', 72, 174, '1999-11-12 09:11:17.17', 'O', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-359266', 'Mr._ไทย', 'Orrin_ไทย', 'Purdy_ไทย', 'Mr.', 'Orrin', 'Purdy', 'Lao', NULL, 'UN538501', '0005832542', 111, 192, '1964-03-05 03:47:14.304', 'B', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-919225', 'Mrs._ไทย', 'Hettie_ไทย', 'Pagac_ไทย', 'Mrs.', 'Hettie', 'Pagac', 'Australian', NULL, 'CL577071', '0267710735', 132, 160, '1993-04-04 15:14:07.123', 'B', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-209014', 'Miss_ไทย', 'Dagmar_ไทย', 'Orn_ไทย', 'Miss', 'Dagmar', 'Orn', 'Lao', NULL, 'LH380017', '0226294999', 46, 185, '1990-01-09 09:18:30.493', 'AB', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-745945', 'Ms._ไทย', 'Amanda_ไทย', 'Brown_ไทย', 'Ms.', 'Amanda', 'Brown', 'Dutch', NULL, 'WD668765', '0548644104', 32, 160, '1963-08-16 00:26:42.88', 'A', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-179590', 'Miss_ไทย', 'Shayna_ไทย', 'Waelchi_ไทย', 'Miss', 'Shayna', 'Waelchi', 'Finnish', NULL, 'LQ792489', '0943118295', 43, 139, '1997-04-07 11:14:45.483', 'AB', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-427845', 'Dr._ไทย', 'Robbie_ไทย', 'Murphy_ไทย', 'Dr.', 'Robbie', 'Murphy', 'Thai', '4671253551800', NULL, '0863848430', 81, 175, '1955-11-28 09:26:03.322', 'A', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-762894', 'Ms._ไทย', 'Ludie_ไทย', 'Schultz_ไทย', 'Ms.', 'Ludie', 'Schultz', 'British', NULL, 'PW230748', '0636966348', 43, 198, '1982-02-23 03:55:55.727', 'A', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-323662', 'Mrs._ไทย', 'Keeley_ไทย', 'Nitzsche_ไทย', 'Mrs.', 'Keeley', 'Nitzsche', 'Thai', '6144057883297', NULL, '0196599324', 57, 182, '1994-09-24 01:52:52.639', 'O', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-873016', 'Mr._ไทย', 'Davion_ไทย', 'Beer_ไทย', 'Mr.', 'Davion', 'Beer', 'British', NULL, 'SI747264', '0318635447', 67, 193, '1965-03-08 19:37:38.002', 'AB', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-173595', 'Miss_ไทย', 'Xavier_ไทย', 'Harber_ไทย', 'Miss', 'Xavier', 'Harber', 'British', NULL, 'WT199926', '0698114146', 107, 144, '1958-03-29 14:59:45.594', 'A', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-352218', 'Ms._ไทย', 'Sabryna_ไทย', 'Von_ไทย', 'Ms.', 'Sabryna', 'Von', 'British', NULL, 'TO414356', '0028009139', 96, 135, '1942-05-27 06:52:09.677', 'B', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-379566', 'Mr._ไทย', 'Nathanial_ไทย', 'Olson_ไทย', 'Mr.', 'Nathanial', 'Olson', 'Thai', '5363483669607', NULL, '0626139789', 97, 180, '1981-06-11 09:58:45.577', 'B', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-427797', 'Ms._ไทย', 'Davon_ไทย', 'Kertzmann_ไทย', 'Ms.', 'Davon', 'Kertzmann', 'American', NULL, 'IQ878160', '0264820394', 66, 181, '2000-07-19 18:45:24.317', 'O', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-254739', 'Mrs._ไทย', 'Janessa_ไทย', 'Greenholt_ไทย', 'Mrs.', 'Janessa', 'Greenholt', 'Thai', '7995428213939', NULL, '0311074482', 100, 183, '1957-08-29 22:21:07.153', 'B', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-297559', 'Dr._ไทย', 'Waldo_ไทย', 'Runte_ไทย', 'Dr.', 'Waldo', 'Runte', 'Thai', '3384069894640', NULL, '0515471189', 105, 162, '1977-01-13 17:36:13.016', 'O', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-190919', 'Mr._ไทย', 'Al_ไทย', 'Schmitt_ไทย', 'Mr.', 'Al', 'Schmitt', 'Australian', NULL, 'PM190288', '0437853017', 40, 188, '1997-05-08 17:04:03.364', 'AB', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-517246', 'Mrs._ไทย', 'Ryleigh_ไทย', 'Hahn_ไทย', 'Mrs.', 'Ryleigh', 'Hahn', 'Thai', '4154429102483', NULL, '0344433248', 115, 152, '1960-10-12 13:02:55.954', 'AB', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-465636', 'Ms._ไทย', 'Aleen_ไทย', 'Harber_ไทย', 'Ms.', 'Aleen', 'Harber', 'Australian', NULL, 'QI847902', '0573425196', 137, 128, '1960-09-08 11:46:47.926', 'O', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-629664', 'Ms._ไทย', 'Remington_ไทย', 'Hamill_ไทย', 'Ms.', 'Remington', 'Hamill', 'American', NULL, 'SY618943', '0091455362', 132, 185, '1967-08-22 17:40:41.948', 'AB', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-693045', 'Ms._ไทย', 'Sammie_ไทย', 'Mertz_ไทย', 'Ms.', 'Sammie', 'Mertz', 'Thai', '9877076115898', NULL, '0504010256', 129, 123, '1980-06-24 21:50:05.699', 'AB', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-838578', 'Miss_ไทย', 'Kylie_ไทย', 'Batz_ไทย', 'Miss', 'Kylie', 'Batz', 'British', NULL, 'YS382368', '0871945370', 82, 153, '1957-06-27 01:56:52.31', 'B', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-200755', 'Mrs._ไทย', 'Abbigail_ไทย', 'Greenfelder_ไทย', 'Mrs.', 'Abbigail', 'Greenfelder', 'Dutch', NULL, 'UZ289259', '0434024632', 124, 186, '1980-07-11 00:44:15.531', 'AB', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-951062', 'Mrs._ไทย', 'Rosemary_ไทย', 'Hackett_ไทย', 'Mrs.', 'Rosemary', 'Hackett', 'Australian', NULL, 'RF114575', '0491486893', 120, 143, '1982-10-01 16:35:28.376', 'AB', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-932425', 'Dr._ไทย', 'Lyda_ไทย', 'Hintz_ไทย', 'Dr.', 'Lyda', 'Hintz', 'Thai', '6414456328781', NULL, '0375611752', 105, 163, '1942-09-02 00:26:29.74', 'B', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-209464', 'Dr._ไทย', 'Craig_ไทย', 'Purdy_ไทย', 'Dr.', 'Craig', 'Purdy', 'Dutch', NULL, 'ME399466', '0551613890', 100, 154, '1963-12-18 21:11:18.464', 'A', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-796550', 'Miss_ไทย', 'Rebeca_ไทย', 'Von_ไทย', 'Miss', 'Rebeca', 'Von', 'British', NULL, 'ZR451303', '0959679238', 122, 126, '1962-05-03 22:46:47.356', 'A', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-835237', 'Dr._ไทย', 'Dayana_ไทย', 'Legros_ไทย', 'Dr.', 'Dayana', 'Legros', 'American', NULL, 'YM317791', '0004223215', 137, 130, '1974-10-11 22:50:43.342', 'O', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-414878', 'Miss_ไทย', 'Gladys_ไทย', 'D''Amore_ไทย', 'Miss', 'Gladys', 'D''Amore', 'American', NULL, 'JZ354997', '0903495864', 91, 131, '1986-02-06 03:49:31.257', 'B', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-576415', 'Dr._ไทย', 'Efren_ไทย', 'Turner_ไทย', 'Dr.', 'Efren', 'Turner', 'Thai', '9025343428386', NULL, '0830073524', 60, 196, '2001-09-17 10:29:51.804', 'O', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-993787', 'Miss_ไทย', 'Maiya_ไทย', 'Schaden_ไทย', 'Miss', 'Maiya', 'Schaden', 'Finnish', NULL, 'VU328689', '0844356699', 39, 135, '1942-01-08 14:05:41.055', 'A', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-740157', 'Mr._ไทย', 'Andres_ไทย', 'McGlynn_ไทย', 'Mr.', 'Andres', 'McGlynn', 'British', NULL, 'TP664495', '0099226719', 42, 157, '1952-04-05 10:00:12.528', 'B', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-228848', 'Miss_ไทย', 'Kelley_ไทย', 'Rath_ไทย', 'Miss', 'Kelley', 'Rath', 'Dutch', NULL, 'VL497056', '0195337477', 117, 136, '1959-01-17 20:53:17.187', 'O', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-651394', 'Mrs._ไทย', 'Maddison_ไทย', 'Mills_ไทย', 'Mrs.', 'Maddison', 'Mills', 'American', NULL, 'LJ451606', '0321248243', 60, 177, '1992-07-13 20:09:41.098', 'AB', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-633981', 'Mr._ไทย', 'Dion_ไทย', 'Wolf_ไทย', 'Mr.', 'Dion', 'Wolf', 'Thai', '8284903599520', NULL, '0493931240', 112, 178, '1991-08-27 00:29:43.405', 'A', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-893772', 'Ms._ไทย', 'Freeman_ไทย', 'DuBuque_ไทย', 'Ms.', 'Freeman', 'DuBuque', 'Lao', NULL, 'WV699153', '0182865247', 56, 191, '1969-06-17 16:55:03.802', 'A', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-733025', 'Miss_ไทย', 'Emmie_ไทย', 'Hegmann_ไทย', 'Miss', 'Emmie', 'Hegmann', 'British', NULL, 'XR613693', '0995181298', 82, 182, '1946-12-28 10:14:10.875', 'O', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-137318', 'Dr._ไทย', 'Tabitha_ไทย', 'Thiel_ไทย', 'Dr.', 'Tabitha', 'Thiel', 'Thai', '2856253163126', NULL, '0348118259', 139, 121, '1960-02-14 10:26:13.645', 'A', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-853857', 'Mr._ไทย', 'Isabelle_ไทย', 'Waters_ไทย', 'Mr.', 'Isabelle', 'Waters', 'Dutch', NULL, 'XZ849789', '0981706912', 90, 166, '1962-08-06 16:57:50.24', 'B', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-129512', 'Miss_ไทย', 'Francisco_ไทย', 'Carter_ไทย', 'Miss', 'Francisco', 'Carter', 'Finnish', NULL, 'WJ771038', '0573741466', 36, 132, '1963-04-29 07:00:44.337', 'O', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-844945', 'Mrs._ไทย', 'Margarita_ไทย', 'Zboncak_ไทย', 'Mrs.', 'Margarita', 'Zboncak', 'Dutch', NULL, 'GK628982', '0332425790', 73, 176, '1972-12-10 21:42:40.074', 'AB', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-910978', 'Mrs._ไทย', 'Gudrun_ไทย', 'Murazik_ไทย', 'Mrs.', 'Gudrun', 'Murazik', 'Finnish', NULL, 'JV376547', '0827741380', 45, 197, '1973-05-19 15:04:28.348', 'O', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-803959', 'Mrs._ไทย', 'Keith_ไทย', 'Nicolas_ไทย', 'Mrs.', 'Keith', 'Nicolas', 'American', NULL, 'DR183706', '0062064868', 140, 169, '1947-12-25 15:12:27.148', 'A', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-336491', 'Dr._ไทย', 'Brando_ไทย', 'Miller_ไทย', 'Dr.', 'Brando', 'Miller', 'Dutch', NULL, 'KG742329', '0611739851', 94, 149, '1999-05-04 08:14:13.779', 'B', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-698271', 'Dr._ไทย', 'Tristian_ไทย', 'Brekke_ไทย', 'Dr.', 'Tristian', 'Brekke', 'British', NULL, 'RC144058', '0060873582', 48, 121, '1962-10-18 19:25:33.731', 'O', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-159601', 'Miss_ไทย', 'Clifton_ไทย', 'Rowe_ไทย', 'Miss', 'Clifton', 'Rowe', 'American', NULL, 'VJ152819', '0659995642', 32, 145, '1955-01-14 05:52:49.758', 'A', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-285237', 'Miss_ไทย', 'Jade_ไทย', 'Blanda_ไทย', 'Miss', 'Jade', 'Blanda', 'Lao', NULL, 'UV500486', '0141653332', 102, 145, '1981-01-09 05:38:00.206', 'O', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-284139', 'Dr._ไทย', 'Avis_ไทย', 'Jenkins_ไทย', 'Dr.', 'Avis', 'Jenkins', 'Australian', NULL, 'OY474751', '0466953079', 115, 128, '1990-03-10 11:55:12.789', 'O', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-561154', 'Miss_ไทย', 'Sasha_ไทย', 'Herzog_ไทย', 'Miss', 'Sasha', 'Herzog', 'Australian', NULL, 'ZN660021', '0462590732', 136, 199, '1944-09-15 18:52:41.526', 'A', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-450439', 'Ms._ไทย', 'Leola_ไทย', 'Buckridge_ไทย', 'Ms.', 'Leola', 'Buckridge', 'Lao', NULL, 'FF435529', '0258069592', 36, 196, '1996-09-03 19:45:03.514', 'AB', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-673284', 'Dr._ไทย', 'Matilda_ไทย', 'Yundt_ไทย', 'Dr.', 'Matilda', 'Yundt', 'American', NULL, 'DL274169', '0282966250', 136, 151, '1997-04-13 17:02:21.757', 'B', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-168114', 'Dr._ไทย', 'Eula_ไทย', 'Gutkowski_ไทย', 'Dr.', 'Eula', 'Gutkowski', 'Lao', NULL, 'BN815386', '0676245424', 115, 140, '1948-08-04 19:56:50.426', 'A', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-421398', 'Ms._ไทย', 'Omer_ไทย', 'Huels_ไทย', 'Ms.', 'Omer', 'Huels', 'Finnish', NULL, 'XG944839', '0484442628', 50, 192, '1957-08-27 06:45:45.043', 'AB', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-415803', 'Mr._ไทย', 'Christ_ไทย', 'Rolfson_ไทย', 'Mr.', 'Christ', 'Rolfson', 'Thai', '2263204321556', NULL, '0075097828', 47, 193, '1948-06-13 14:00:22.73', 'O', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-862140', 'Miss_ไทย', 'Niko_ไทย', 'Will_ไทย', 'Miss', 'Niko', 'Will', 'Dutch', NULL, 'ML783404', '0741337680', 84, 188, '1953-05-10 07:34:46.316', 'O', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-113321', 'Ms._ไทย', 'Shaina_ไทย', 'Lind_ไทย', 'Ms.', 'Shaina', 'Lind', 'American', NULL, 'PT662483', '0632577608', 68, 137, '1977-09-23 22:48:34.853', 'AB', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-188305', 'Mr._ไทย', 'Ella_ไทย', 'Heidenreich_ไทย', 'Mr.', 'Ella', 'Heidenreich', 'Dutch', NULL, 'JP674942', '0985975042', 103, 162, '1989-10-25 14:48:53.179', 'B', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-262701', 'Mr._ไทย', 'Isac_ไทย', 'McCullough_ไทย', 'Mr.', 'Isac', 'McCullough', 'American', NULL, 'NE579992', '0593046521', 89, 171, '2000-03-06 18:01:20.631', 'AB', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-774386', 'Mrs._ไทย', 'Chester_ไทย', 'Morar_ไทย', 'Mrs.', 'Chester', 'Morar', 'Australian', NULL, 'GD685378', '0628938501', 45, 187, '1984-05-14 16:19:42.185', 'B', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-162951', 'Dr._ไทย', 'Isidro_ไทย', 'Nienow_ไทย', 'Dr.', 'Isidro', 'Nienow', 'Finnish', NULL, 'UX452495', '0550437688', 44, 170, '1959-03-28 05:47:29.114', 'B', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-618756', 'Mrs._ไทย', 'Kendrick_ไทย', 'Spinka_ไทย', 'Mrs.', 'Kendrick', 'Spinka', 'Finnish', NULL, 'AJ227956', '0299825548', 138, 181, '1955-09-04 22:32:00.072', 'O', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-734281', 'Ms._ไทย', 'Freida_ไทย', 'Leannon_ไทย', 'Ms.', 'Freida', 'Leannon', 'Australian', NULL, 'PN553121', '0284591504', 36, 148, '1974-03-07 12:37:22.216', 'B', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-680193', 'Ms._ไทย', 'Audie_ไทย', 'Jacobi_ไทย', 'Ms.', 'Audie', 'Jacobi', 'Thai', '8203073395562', NULL, '0437670889', 118, 146, '2003-11-19 11:23:16.228', 'O', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-757364', 'Ms._ไทย', 'Sylvester_ไทย', 'Koss_ไทย', 'Ms.', 'Sylvester', 'Koss', 'Dutch', NULL, 'SH191284', '0793880455', 118, 126, '1975-01-18 22:39:00.236', 'O', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-293338', 'Ms._ไทย', 'Telly_ไทย', 'Metz_ไทย', 'Ms.', 'Telly', 'Metz', 'Finnish', NULL, 'SQ580400', '0503759044', 108, 130, '1974-10-26 20:50:49.073', 'AB', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-883377', 'Ms._ไทย', 'Skylar_ไทย', 'Pouros_ไทย', 'Ms.', 'Skylar', 'Pouros', 'Lao', NULL, 'GG892527', '0183716247', 46, 141, '1980-03-07 23:02:18.314', 'B', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-874621', 'Ms._ไทย', 'Flo_ไทย', 'Zemlak_ไทย', 'Ms.', 'Flo', 'Zemlak', 'Finnish', NULL, 'ZP855320', '0733979280', 46, 123, '1976-08-25 21:43:56.177', 'AB', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-686201', 'Ms._ไทย', 'Kavon_ไทย', 'Mayert_ไทย', 'Ms.', 'Kavon', 'Mayert', 'Thai', '4478018125737', NULL, '0582303052', 135, 124, '1995-02-15 23:31:50.529', 'A', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-552487', 'Mrs._ไทย', 'Haley_ไทย', 'Windler_ไทย', 'Mrs.', 'Haley', 'Windler', 'Australian', NULL, 'NK471460', '0089007090', 68, 141, '2004-11-05 02:37:04.314', 'A', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-194183', 'Mr._ไทย', 'Robin_ไทย', 'O''Connell_ไทย', 'Mr.', 'Robin', 'O''Connell', 'American', NULL, 'KD173328', '0769942937', 79, 160, '1999-11-25 02:46:48.142', 'A', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-753461', 'Mrs._ไทย', 'Delilah_ไทย', 'Conroy_ไทย', 'Mrs.', 'Delilah', 'Conroy', 'Dutch', NULL, 'QL646697', '0180425336', 41, 143, '1991-02-01 15:20:11.211', 'B', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-347868', 'Miss_ไทย', 'Ethyl_ไทย', 'Shields_ไทย', 'Miss', 'Ethyl', 'Shields', 'Dutch', NULL, 'KU718641', '0544493228', 131, 167, '1942-11-01 11:19:18.33', 'O', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-843948', 'Mr._ไทย', 'Jevon_ไทย', 'Ruecker_ไทย', 'Mr.', 'Jevon', 'Ruecker', 'Thai', '3471236873503', NULL, '0015506424', 135, 137, '1957-01-27 00:30:05.156', 'B', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-725244', 'Mrs._ไทย', 'Maida_ไทย', 'Oberbrunner_ไทย', 'Mrs.', 'Maida', 'Oberbrunner', 'Lao', NULL, 'MW814458', '0573512265', 89, 190, '1989-12-17 19:09:18.595', 'AB', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-615472', 'Mrs._ไทย', 'Jarod_ไทย', 'Rosenbaum_ไทย', 'Mrs.', 'Jarod', 'Rosenbaum', 'Finnish', NULL, 'EU309219', '0634245669', 52, 170, '1951-05-11 09:35:51.809', 'A', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-784077', 'Mrs._ไทย', 'Marcelino_ไทย', 'Waters_ไทย', 'Mrs.', 'Marcelino', 'Waters', 'Finnish', NULL, 'SW884063', '0116632391', 100, 172, '1976-10-21 04:05:29.445', 'B', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-634304', 'Miss_ไทย', 'Rudolph_ไทย', 'Reynolds_ไทย', 'Miss', 'Rudolph', 'Reynolds', 'British', NULL, 'FV317423', '0681589886', 119, 144, '1960-08-05 18:35:38.311', 'AB', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-861496', 'Ms._ไทย', 'Peyton_ไทย', 'Schultz_ไทย', 'Ms.', 'Peyton', 'Schultz', 'British', NULL, 'CR224892', '0558124382', 37, 137, '1977-11-08 03:33:14.162', 'B', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-881131', 'Mrs._ไทย', 'Willis_ไทย', 'Bernhard_ไทย', 'Mrs.', 'Willis', 'Bernhard', 'Finnish', NULL, 'YQ121943', '0075904440', 112, 124, '1998-02-08 19:18:08.5', 'AB', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-934746', 'Mrs._ไทย', 'Trevion_ไทย', 'Bogisich_ไทย', 'Mrs.', 'Trevion', 'Bogisich', 'Thai', '7553151235864', NULL, '0945014359', 149, 140, '1994-06-20 21:53:57.493', 'AB', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-163878', 'Mr._ไทย', 'Coy_ไทย', 'Gleason_ไทย', 'Mr.', 'Coy', 'Gleason', 'Dutch', NULL, 'WV830922', '0534047401', 126, 173, '1983-09-13 15:38:53.77', 'AB', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-897438', 'Ms._ไทย', 'Devante_ไทย', 'O''Reilly_ไทย', 'Ms.', 'Devante', 'O''Reilly', 'Thai', '1566300176814', NULL, '0544831378', 139, 179, '1977-12-31 22:09:58.279', 'A', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-221203', 'Mrs._ไทย', 'Carlee_ไทย', 'Predovic_ไทย', 'Mrs.', 'Carlee', 'Predovic', 'Thai', '5951742041675', NULL, '0886326883', 126, 190, '1990-09-11 09:45:08.769', 'O', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-519312', 'Miss_ไทย', 'Gene_ไทย', 'Kirlin_ไทย', 'Miss', 'Gene', 'Kirlin', 'Thai', '9279435417889', NULL, '0744354569', 67, 134, '1996-05-01 17:00:23.392', 'B', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-517289', 'Mrs._ไทย', 'Kale_ไทย', 'Raynor_ไทย', 'Mrs.', 'Kale', 'Raynor', 'Lao', NULL, 'JL530874', '0017206543', 39, 164, '1979-06-24 02:45:34.325', 'AB', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-997298', 'Miss_ไทย', 'Alicia_ไทย', 'Schowalter_ไทย', 'Miss', 'Alicia', 'Schowalter', 'British', NULL, 'WZ575421', '0081618227', 56, 174, '2002-03-30 18:07:32.636', 'A', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-374325', 'Miss_ไทย', 'Jameson_ไทย', 'Abshire_ไทย', 'Miss', 'Jameson', 'Abshire', 'Australian', NULL, 'LX489527', '0541627966', 91, 137, '1998-11-11 20:29:28.416', 'AB', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-126977', 'Mrs._ไทย', 'Ona_ไทย', 'Hodkiewicz_ไทย', 'Mrs.', 'Ona', 'Hodkiewicz', 'Dutch', NULL, 'AX189457', '0845700300', 91, 153, '1976-03-25 19:29:59.928', 'B', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-314696', 'Ms._ไทย', 'Wallace_ไทย', 'Mayert_ไทย', 'Ms.', 'Wallace', 'Mayert', 'American', NULL, 'GT685731', '0596704436', 101, 162, '1944-03-14 19:27:03.065', 'B', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-175152', 'Mrs._ไทย', 'Vivienne_ไทย', 'Denesik_ไทย', 'Mrs.', 'Vivienne', 'Denesik', 'Finnish', NULL, 'YN232430', '0660568004', 38, 196, '1986-03-16 07:21:11.827', 'O', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-846103', 'Mrs._ไทย', 'Princess_ไทย', 'Koch_ไทย', 'Mrs.', 'Princess', 'Koch', 'Thai', '1108182787046', NULL, '0109053804', 66, 183, '1994-01-25 02:24:43.475', 'B', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-119713', 'Ms._ไทย', 'Mariane_ไทย', 'Stoltenberg_ไทย', 'Ms.', 'Mariane', 'Stoltenberg', 'Finnish', NULL, 'FQ529605', '0673415653', 77, 137, '1964-12-08 20:02:37.192', 'A', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');
INSERT INTO public."Patient" (id, initial_th, firstname_th, lastname_th, initial_en, firstname_en, lastname_en, nationality, "nationalId", "passportId", "phoneNumber", weight, height, "birthDate", "bloodType", "createdAt", "updatedAt", "profilePicURL") VALUES ('HN-322420', 'Ms._ไทย', 'Arianna_ไทย', 'Hahn_ไทย', 'Ms.', 'Arianna', 'Hahn', 'Finnish', NULL, 'NG839687', '0959640148', 59, 164, '1972-03-09 06:20:23.133', 'A', '2022-09-07 10:52:45.209', '2022-09-07 10:52:45.211', 'https://img.cscms.me/X53hRfCNpi9bGzEGH2gn.png');


--
-- Data for Name: Prescription; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (1, 7, 2, 90, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (2, 3, 2, 40, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (3, 5, 2, 80, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (4, 1, 2, 10, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (5, 9, 4, 70, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (6, 3, 9, 80, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (7, 1, 9, 30, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (8, 2, 9, 50, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (9, 10, 9, 10, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (10, 5, 11, 10, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (11, 8, 11, 80, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (12, 9, 11, 60, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (13, 8, 11, 10, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (14, 2, 11, 50, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (15, 1, 11, 40, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (16, 6, 11, 70, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (17, 1, 11, 70, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (18, 3, 11, 50, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (19, 8, 23, 10, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (20, 2, 23, 10, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (21, 5, 23, 60, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (22, 1, 23, 90, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (23, 8, 23, 90, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (24, 5, 23, 70, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (25, 6, 23, 70, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (26, 3, 23, 70, '2022-09-07 10:52:46.95', '2022-09-07 10:52:46.951');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (27, 3, 28, 80, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (28, 9, 28, 10, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (29, 3, 28, 80, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (30, 7, 28, 20, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (31, 7, 37, 70, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (32, 2, 38, 40, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (33, 7, 38, 40, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (34, 10, 31, 50, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.951');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (35, 3, 31, 90, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.951');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (36, 9, 31, 10, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.951');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (37, 8, 31, 50, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.951');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (38, 5, 31, 80, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.951');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (39, 2, 39, 30, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (40, 5, 39, 70, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (41, 2, 39, 50, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (42, 1, 39, 70, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (43, 7, 45, 40, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (44, 1, 45, 70, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (45, 7, 45, 70, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (46, 5, 45, 30, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (47, 4, 53, 60, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (48, 9, 53, 20, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (49, 4, 50, 80, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.952');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (50, 4, 50, 80, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.952');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (51, 3, 50, 10, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.952');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (52, 8, 50, 40, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.952');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (53, 3, 50, 10, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.952');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (54, 7, 50, 70, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.952');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (55, 7, 49, 60, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (56, 7, 49, 20, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (57, 3, 49, 80, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (58, 2, 49, 60, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (59, 6, 49, 10, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (60, 3, 49, 50, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (61, 9, 49, 90, '2022-09-07 10:52:46.951', '2022-09-07 10:52:46.952');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (62, 6, 56, 10, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (63, 6, 59, 90, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (64, 4, 59, 20, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (65, 10, 59, 40, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (66, 1, 59, 60, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (67, 8, 59, 50, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (68, 10, 59, 60, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (69, 4, 59, 40, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (70, 4, 59, 30, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (71, 8, 59, 80, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (72, 2, 58, 30, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (73, 2, 63, 80, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (74, 8, 58, 70, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (75, 10, 63, 20, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (76, 2, 58, 20, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (77, 9, 63, 60, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (78, 6, 58, 20, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (79, 6, 63, 10, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (80, 2, 63, 70, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (81, 4, 63, 20, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (82, 9, 63, 70, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (83, 7, 63, 50, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (84, 4, 58, 10, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (85, 6, 58, 60, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (86, 6, 58, 70, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (87, 7, 58, 70, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (88, 10, 58, 90, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (89, 3, 58, 20, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (90, 7, 70, 80, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (91, 6, 70, 20, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (92, 2, 64, 40, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (93, 4, 64, 60, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (94, 3, 64, 60, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (95, 10, 64, 80, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (96, 8, 64, 80, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (97, 1, 64, 50, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (98, 8, 64, 90, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (99, 3, 64, 20, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (100, 2, 64, 40, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (101, 4, 73, 20, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (102, 6, 73, 40, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (103, 2, 73, 60, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (104, 7, 73, 90, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (105, 6, 73, 60, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (106, 1, 73, 10, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (107, 7, 73, 60, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (108, 9, 73, 90, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (109, 9, 73, 60, '2022-09-07 10:52:46.952', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (110, 5, 84, 20, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (111, 2, 84, 60, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (112, 9, 84, 90, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (113, 5, 84, 80, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (114, 5, 84, 90, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (115, 2, 84, 70, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (116, 1, 84, 70, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (117, 4, 84, 20, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (118, 4, 85, 10, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (119, 10, 85, 70, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (120, 2, 93, 30, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (121, 4, 93, 20, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (122, 7, 93, 10, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (123, 10, 93, 20, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (124, 5, 93, 90, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (125, 4, 93, 30, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (126, 1, 89, 40, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (127, 6, 89, 80, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (128, 7, 89, 50, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (129, 8, 89, 10, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (130, 9, 89, 40, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (131, 2, 89, 90, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (132, 3, 89, 60, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (133, 9, 89, 10, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (134, 5, 99, 60, '2022-09-07 10:52:46.954', '2022-09-07 10:52:46.954');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (135, 3, 99, 60, '2022-09-07 10:52:46.954', '2022-09-07 10:52:46.954');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (136, 10, 99, 30, '2022-09-07 10:52:46.954', '2022-09-07 10:52:46.954');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (137, 2, 99, 20, '2022-09-07 10:52:46.954', '2022-09-07 10:52:46.954');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (138, 4, 98, 70, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (139, 1, 98, 50, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (140, 5, 98, 30, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (141, 7, 98, 90, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (142, 7, 98, 10, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (143, 9, 98, 40, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.954');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (144, 8, 96, 70, '2022-09-07 10:52:46.954', '2022-09-07 10:52:46.955');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (145, 2, 96, 80, '2022-09-07 10:52:46.954', '2022-09-07 10:52:46.955');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (146, 4, 96, 60, '2022-09-07 10:52:46.954', '2022-09-07 10:52:46.955');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (147, 6, 96, 60, '2022-09-07 10:52:46.954', '2022-09-07 10:52:46.955');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (148, 7, 96, 70, '2022-09-07 10:52:46.954', '2022-09-07 10:52:46.955');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (149, 8, 96, 50, '2022-09-07 10:52:46.954', '2022-09-07 10:52:46.955');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (150, 9, 97, 40, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (151, 9, 97, 10, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (152, 5, 97, 90, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (153, 7, 97, 90, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (154, 4, 97, 60, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (155, 10, 97, 90, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (156, 3, 97, 30, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (157, 1, 97, 80, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (158, 3, 97, 90, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.953');
INSERT INTO public."Prescription" (id, "medicineId", "appointmentId", amount, "createdAt", "updatedAt") VALUES (159, 4, 97, 70, '2022-09-07 10:52:46.953', '2022-09-07 10:52:46.953');


--
-- Data for Name: _prisma_migrations; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public._prisma_migrations (id, checksum, finished_at, migration_name, logs, rolled_back_at, started_at, applied_steps_count) VALUES ('df490715-f51c-48cb-abf9-5febd686d8be', '75a8aec5f84334978dc6e95b0a6fbb930a9eb0448d04da25fb4b584167831024', '2022-09-27 08:39:42.163213+00', '20220927075430_init', NULL, NULL, '2022-09-27 08:39:42.124804+00', 1);
INSERT INTO public._prisma_migrations (id, checksum, finished_at, migration_name, logs, rolled_back_at, started_at, applied_steps_count) VALUES ('889f14d5-e5a5-4cfa-af1f-cee4407605a3', '3879844f0c9d9235ce2512bba45dc322e4e7bb254826f4293106aa05f80f05a4', '2022-09-27 08:39:49.25402+00', '20220927083949_invoice_discount', NULL, NULL, '2022-09-27 08:39:49.243081+00', 1);
INSERT INTO public._prisma_migrations (id, checksum, finished_at, migration_name, logs, rolled_back_at, started_at, applied_steps_count) VALUES ('614ba9fa-f46e-4d7a-b4f4-4d5e1e4a2ebf', '006ff93640be69a1b7d7bc4fefdf9a74766f3194beb766d75e78c4eb8312aee6', '2022-09-27 15:36:56.715377+00', '20220927153359_add_medicine_picture', NULL, NULL, '2022-09-27 15:36:56.706172+00', 1);
INSERT INTO public._prisma_migrations (id, checksum, finished_at, migration_name, logs, rolled_back_at, started_at, applied_steps_count) VALUES ('cf1c3e76-83ee-458c-8be0-ba28d978f9ea', 'aa539a968aacecc0e3185fa860ab82290f645d753ab86918415e8545c987d844', '2022-09-27 15:37:10.798269+00', '20220927153710_add_medicine_picture', NULL, NULL, '2022-09-27 15:37:10.792696+00', 1);
INSERT INTO public._prisma_migrations (id, checksum, finished_at, migration_name, logs, rolled_back_at, started_at, applied_steps_count) VALUES ('2d817aea-b99d-4cfb-a87b-ef899f1b64a5', 'f3ad4fdd17f81b120b0655a023b3c0499c41b76244f857b15fadda3484d3db95', '2022-10-26 10:26:12.098338+00', '20221026090657_add_patient_profile_pic_url', NULL, NULL, '2022-10-26 10:26:12.093527+00', 1);
INSERT INTO public._prisma_migrations (id, checksum, finished_at, migration_name, logs, rolled_back_at, started_at, applied_steps_count) VALUES ('a219fe50-08ba-43b8-9ff2-fc89176ae7b5', '474ebc0237fa478f1196ca19c3113cff3e3efa92e8dce4d203bf614ecdd428ac', '2022-10-26 10:26:12.103157+00', '20221026092334_add_patient_profile_pic_url', NULL, NULL, '2022-10-26 10:26:12.099543+00', 1);


--
-- Name: Appointment_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."Appointment_id_seq"', 1, false);


--
-- Name: Doctor_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."Doctor_id_seq"', 1, false);


--
-- Name: InvoiceDiscount_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."InvoiceDiscount_id_seq"', 1, true);


--
-- Name: InvoiceItem_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."InvoiceItem_id_seq"', 1, false);


--
-- Name: Invoice_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."Invoice_id_seq"', 1, false);


--
-- Name: Medicine_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."Medicine_id_seq"', 1, false);


--
-- Name: Prescription_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."Prescription_id_seq"', 1, false);


--
-- Name: Appointment Appointment_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Appointment"
    ADD CONSTRAINT "Appointment_pkey" PRIMARY KEY (id);


--
-- Name: Doctor Doctor_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Doctor"
    ADD CONSTRAINT "Doctor_pkey" PRIMARY KEY (id);


--
-- Name: InvoiceDiscount InvoiceDiscount_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."InvoiceDiscount"
    ADD CONSTRAINT "InvoiceDiscount_pkey" PRIMARY KEY (id);


--
-- Name: InvoiceItem InvoiceItem_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."InvoiceItem"
    ADD CONSTRAINT "InvoiceItem_pkey" PRIMARY KEY (id);


--
-- Name: Invoice Invoice_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Invoice"
    ADD CONSTRAINT "Invoice_pkey" PRIMARY KEY (id);


--
-- Name: Medicine Medicine_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Medicine"
    ADD CONSTRAINT "Medicine_pkey" PRIMARY KEY (id);


--
-- Name: Patient Patient_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Patient"
    ADD CONSTRAINT "Patient_pkey" PRIMARY KEY (id);


--
-- Name: Prescription Prescription_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Prescription"
    ADD CONSTRAINT "Prescription_pkey" PRIMARY KEY (id);


--
-- Name: _prisma_migrations _prisma_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public._prisma_migrations
    ADD CONSTRAINT _prisma_migrations_pkey PRIMARY KEY (id);


--
-- Name: Doctor_username_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX "Doctor_username_key" ON public."Doctor" USING btree (username);


--
-- Name: Invoice_appointmentId_key; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX "Invoice_appointmentId_key" ON public."Invoice" USING btree ("appointmentId");


--
-- Name: Appointment Appointment_doctorId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Appointment"
    ADD CONSTRAINT "Appointment_doctorId_fkey" FOREIGN KEY ("doctorId") REFERENCES public."Doctor"(id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: Appointment Appointment_patientId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Appointment"
    ADD CONSTRAINT "Appointment_patientId_fkey" FOREIGN KEY ("patientId") REFERENCES public."Patient"(id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: InvoiceDiscount InvoiceDiscount_invoiceId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."InvoiceDiscount"
    ADD CONSTRAINT "InvoiceDiscount_invoiceId_fkey" FOREIGN KEY ("invoiceId") REFERENCES public."Invoice"(id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: InvoiceItem InvoiceItem_invoiceId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."InvoiceItem"
    ADD CONSTRAINT "InvoiceItem_invoiceId_fkey" FOREIGN KEY ("invoiceId") REFERENCES public."Invoice"(id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: Invoice Invoice_appointmentId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Invoice"
    ADD CONSTRAINT "Invoice_appointmentId_fkey" FOREIGN KEY ("appointmentId") REFERENCES public."Appointment"(id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: Prescription Prescription_appointmentId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Prescription"
    ADD CONSTRAINT "Prescription_appointmentId_fkey" FOREIGN KEY ("appointmentId") REFERENCES public."Appointment"(id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- Name: Prescription Prescription_medicineId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Prescription"
    ADD CONSTRAINT "Prescription_medicineId_fkey" FOREIGN KEY ("medicineId") REFERENCES public."Medicine"(id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- PostgreSQL database dump complete
--

