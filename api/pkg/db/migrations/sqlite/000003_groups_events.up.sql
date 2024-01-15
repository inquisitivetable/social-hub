CREATE TABLE IF NOT EXISTS groups(
	id INTEGER PRIMARY KEY,
	creator_id INTEGER NOT NULL,
	title TEXT NOT NULL,
	description TEXT,
	created_at DATETIME NOT NULL,
		FOREIGN KEY (creator_id)
			REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS user_groups (
	id INTEGER PRIMARY KEY,
	user_id INTEGER NOT NULL,
	group_id INTEGER NOT NULL,
	joined_at DATETIME NOT NULL,
	accepted BOOL NOT NULL,
	FOREIGN KEY (user_id) 
		REFERENCES users (id)
	FOREIGN KEY (group_id) 
		REFERENCES groups (id)
);

CREATE TABLE IF NOT EXISTS group_events (
	id INTEGER PRIMARY KEY,
	group_id INTEGER NOT NULL,
	user_id INTEGER NOT NULL,
	created_at DATETIME NOT NULL,
	event_time DATETIME NOT NULL,
	timespan INTEGER,
	title TEXT NOT NULL,
	description TEXT NOT NULL,
	FOREIGN KEY (user_id) 
		REFERENCES users (id)
	FOREIGN KEY (group_id) 
		REFERENCES groups (id)
);

CREATE TABLE IF NOT EXISTS group_event_attendance (
	id INTEGER PRIMARY KEY,
	user_id INTEGER NOT NULL,
	event_id INTEGER NOT NULL,
	is_attending BOOL,
	FOREIGN KEY (user_id) 
		REFERENCES users (id)
	FOREIGN KEY (event_id) 
		REFERENCES group_events (id)
);
