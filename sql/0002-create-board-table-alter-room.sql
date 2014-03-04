CREATE SEQUENCE board_id_seq;
CREATE TABLE board
(
  id integer NOT NULL default nextval('board_id_seq'),
  slug character varying(3),
  description text,
  CONSTRAINT board_pk PRIMARY KEY (id)
);

INSERT INTO board (slug) VALUES ('nop');

ALTER TABLE room ADD COLUMN board_id integer;
UPDATE room SET board_id = 1;
ALTER TABLE room ALTER COLUMN board_id SET NOT NULL;

ALTER TABLE room ADD CONSTRAINT room_fk1 FOREIGN KEY (board_id) REFERENCES board(id);
