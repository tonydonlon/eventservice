# eventservice
streaming event service


## PosgreSQL Database Schema

                    List of relations
 Schema |           Name           |   Type   |  Owner   
--------+--------------------------+----------+----------
 public | event_names              | table    | postgres
 public | event_names_event_id_seq | sequence | postgres
 public | events                   | table    | postgres
 public | events_event_id_seq      | sequence | postgres
 public | sessions                 | table    | postgres
(5 rows)

postgres=# \d events;
                                            Table "public.events"
     Column      |           Type           | Collation | Nullable |                 Default                  
-----------------+--------------------------+-----------+----------+------------------------------------------
 event_id        | bigint                   |           | not null | nextval('events_event_id_seq'::regclass)
 event_timestamp | timestamp with time zone |           | not null | 
 event_name_id   | smallint                 |           | not null | 
 session_id      | uuid                     |           | not null | 
Indexes:
    "events_pkey" PRIMARY KEY, btree (event_id)
Foreign-key constraints:
    "fk_event_name_id" FOREIGN KEY (event_name_id) REFERENCES event_names(event_name_id)
    "fk_session_id" FOREIGN KEY (session_id) REFERENCES sessions(session_id)

postgres=# \d event_names;
                                          Table "public.event_names"
    Column     |          Type          | Collation | Nullable |                    Default                    
---------------+------------------------+-----------+----------+-----------------------------------------------
 event_name_id | smallint               |           | not null | nextval('event_names_event_id_seq'::regclass)
 event_name    | character varying(255) |           | not null | 
Indexes:
    "event_names_pkey" PRIMARY KEY, btree (event_name_id)
Referenced by:
    TABLE "events" CONSTRAINT "fk_event_name_id" FOREIGN KEY (event_name_id) REFERENCES event_names(event_name_id)

postgres=# \d sessions;
                          Table "public.sessions"
    Column     |           Type           | Collation | Nullable | Default 
---------------+--------------------------+-----------+----------+---------
 session_id    | uuid                     |           | not null | 
 session_start | timestamp with time zone |           |          | 
 session_end   | timestamp with time zone |           |          | 
Indexes:
    "sessions_pkey" PRIMARY KEY, btree (session_id)
Referenced by:
    TABLE "events" CONSTRAINT "fk_session_id" FOREIGN KEY (session_id) REFERENCES sessions(session_id)

