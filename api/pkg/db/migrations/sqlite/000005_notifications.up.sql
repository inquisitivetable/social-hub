CREATE TABLE IF NOT EXISTS notification_types(
	id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    entity TEXT NOT NULL
);

INSERT INTO notification_types (id, name, entity)
VALUES 
(0, "follow_request", "followers"),
(1, "group_request", "groups"),
(2, "group_invite", "groups"),
(3, "event_invite", "group_events");



CREATE TABLE IF NOT EXISTS notification_details(
	id INTEGER PRIMARY KEY,
    sender_id INTEGER NOT NULL,
    notification_type_id INTEGER NOT NULL,
    entity_id INTEGER NOT NULL,
	created_at DATETIME NOT NULL,
	FOREIGN KEY (sender_id) 
		REFERENCES users (id)
	FOREIGN KEY (notification_type_id) 
		REFERENCES notification_types (id)
);


CREATE TABLE IF NOT EXISTS notifications(
	id INTEGER PRIMARY KEY,
    receiver_id INTEGER NOT NULL,
    notification_details_id INTEGER NOT NULL,
    seen_at DATETIME,
    reaction BOOL,
	FOREIGN KEY (receiver_id) 
		REFERENCES users (id)
	FOREIGN KEY (notification_details_id) 
		REFERENCES notification_details (id)
);