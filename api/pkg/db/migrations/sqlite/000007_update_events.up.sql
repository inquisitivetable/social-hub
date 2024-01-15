ALTER TABLE group_events DROP COLUMN timespan;

ALTER TABLE group_events 
ADD COLUMN event_end_time DATETIME NOT NULL;