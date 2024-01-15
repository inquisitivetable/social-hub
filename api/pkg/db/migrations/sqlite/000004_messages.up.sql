CREATE TABLE IF NOT EXISTS messages(
	id INTEGER PRIMARY KEY,
	sender_id INTEGER NOT NULL,
	recipient_id INTEGER,
	group_id INTEGER,
	content TEXT NOT NULL,
	image_path TEXT,
	sent_at DATETIME NOT NULL,
	read_at DATETIME,
	FOREIGN KEY (sender_id) 
		REFERENCES users (id)
	FOREIGN KEY (recipient_id) 
		REFERENCES users (id)
	FOREIGN KEY (group_id) 
		REFERENCES groups (id)
);