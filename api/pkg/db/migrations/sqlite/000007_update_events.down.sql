ALTER TABLE group_events 
ADD COLUMN timespan INTEGER;

ALTER TABLE group_events 
DROP COLUMN  event_end_time;