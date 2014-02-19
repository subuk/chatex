ALTER TABLE message ALTER create_time SET DEFAULT (NOW() at time zone 'utc');
