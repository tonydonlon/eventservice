
BEGIN;
CREATE OR REPLACE FUNCTION customer_events_notify()
	RETURNS trigger AS
$$
BEGIN
	PERFORM pg_notify('customer_events', row_to_json(NEW)::text);
	RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER events_status
	AFTER INSERT
	ON events
	FOR EACH ROW
EXECUTE PROCEDURE customer_events_notify();
COMMIT;