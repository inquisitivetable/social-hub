CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY,
	forname TEXT NOT NULL,
	surname TEXT NOT NULL,
	email TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL,
	birthday DATETIME,
	nickname TEXT,
	about TEXT,
	image_path TEXT,
	created_at DATETIME,
	is_public BOOL NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS user_sessions(
	id INTEGER PRIMARY KEY,
	user_id INTEGER NOT NULL,
	token TEXT NOT NUll,
	created_at DATETIME NOT NULL,
		FOREIGN KEY (user_id)
			REFERENCES users(id)

);


CREATE TABLE IF NOT EXISTS followers(
	id INTEGER PRIMARY KEY,
	following_id INTEGER NOT NULL,
	follower_id INTEGER NOT NULL,
	accepted bool,
	active bool NOT NULL,
		FOREIGN KEY (following_id)
			REFERENCES users(id)	
		FOREIGN KEY (follower_id)
			REFERENCES users(id)
);
