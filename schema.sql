
-- SCHEMA
CREATE SEQUENCE room_id_seq;
CREATE TABLE room
(
  id integer NOT NULL default nextval('room_id_seq'),
  CONSTRAINT room_pk PRIMARY KEY (id)
);

CREATE SEQUENCE message_id_seq;
CREATE TABLE message
(
  id integer NOT NULL default nextval('message_id_seq'),
  room_id integer NOT NULL,
  text text,
  create_time timestamp without time zone default NOW(),
  CONSTRAINT message_pk PRIMARY KEY (id)
);
ALTER TABLE message ADD CONSTRAINT message_fk1 FOREIGN KEY (room_id) REFERENCES room(id);


CREATE LANGUAGE plpgsql;
CREATE OR REPLACE FUNCTION notify_new_message() RETURNS trigger AS $$
DECLARE
BEGIN
  PERFORM pg_notify('new_message', CAST(NEW.id as text) || '|' || CAST(NEW.room_id as text) || '|' || NEW.text);
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER notify_on_insert_message
    AFTER INSERT ON message
    FOR EACH ROW
    EXECUTE PROCEDURE notify_new_message();


-- DEFAULT DATA

INSERT INTO room VALUES (1);
INSERT INTO message (room_id, text) VALUES (1, 'Hello');
INSERT INTO message (room_id, text) VALUES (1, 'World');
INSERT INTO message (room_id, text) VALUES (1, 'Test');
