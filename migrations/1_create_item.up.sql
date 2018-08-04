CREATE TABLE Item
(
   uid serial NOT NULL,
   name character varying(100) NOT NULL,
   Created date,
   CONSTRAINT userinfo_pkey PRIMARY KEY (uid)
)
