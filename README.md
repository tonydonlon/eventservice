# eventservice
streaming event service - [interview assignment](assignment.pdf)

## HTTP Service
* `/healthcheck` http GET to return simple 200 and 'OK'
* `/ws?sessionId=<UUID>` Websocket connection to upload events
```
Websocket send JSON:
[
  {
    timestamp: 1624422533298,
    type: 'SESSION_START',
    session_id: '0d856b46-48cd-4fba-946b-53d15b7b43e0'
  },
  { timestamp: 1624422533298, type: 'EVENT', name: 'event2' },
  { timestamp: 1624422533298, type: 'EVENT', name: 'test' },
  { timestamp: 1624422533298, type: 'EVENT', name: 'event2' },
  {
    timestamp: 1624422533298,
    type: 'SESSION_END',
    session_id: '0d856b46-48cd-4fba-946b-53d15b7b43e0'
  }
]

```
* `/session/<UUID>` http GET endpoint to fetch a sessoin and it's set of events
```
Response:
{
    "type": "SESSION",
    "start": 1624418542901,
    "end": 1624418543902,
    "children": [
        {
            "type": "EVENT",
            "timestamp": 1624418542901,
            "name": "event1"
        },
        {
            "type": "EVENT",
            "timestamp": 1624418542901,
            "name": "event2"
        }
    ]
}
```
* `/database_events` http GET pub/sub endpoint for Server Sent Events to monitor all events from the DB

## Requirements
*  Docker or running locally: golang compiler plus postgresql server

## Setup
* `docker-compose build` then `docker-compose up`
* database schema is created () on postgres initialization
* database data is persisted to a local volume at `./postgres-data`
    * to restart with a fresh database, just `rm -rf ./postgres-data` with the container shutdown

## Testing
* k6 use for external API function tests: websocket, 
    *  K6 (https://k6.io/) can be used for function-tests and load-test, etc.
    * install k6 (https://k6.io/docs/getting-started/installation): `> brew install k6`
    * run a single endpoint test, for example: `> k6 run test/k6/websocket.js`
    * run the whole suite of tests: `./test/k5/run.sh`


## PosgreSQL Database Schema

```

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

```
