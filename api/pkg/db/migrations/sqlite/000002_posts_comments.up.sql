CREATE TABLE IF NOT EXISTS posts (
	id INTEGER PRIMARY KEY,
	user_id INTEGER NOT NULL,
	privacy_type_id INTEGER NOT NULL,
	title TEXT NOT NULL,
	content TEXT NOT NULL,
	created_at DATETIME NOT NULL,
	image_path TEXT,
		FOREIGN KEY (user_id) 
			REFERENCES users(id)
		FOREIGN KEY (privacy_type_id) 
			REFERENCES privacy_type(id)
);


CREATE TABLE IF NOT EXISTS privacy_type (
	id integer PRIMARY KEY NOT NULL,
	name text UNIQUE NOT NULL,
	description text
);

INSERT INTO privacy_type (id, name)
VALUES 
(1, "public"),
(2, "private"),
(3, "sub-private");


CREATE TABLE IF NOT EXISTS allowed_private_posts (
	id INTEGER PRIMARY KEY,
	post_id INTEGER NOT NULL,
	user_id INTEGER NOT NULL,
	FOREIGN KEY (user_id) 
		REFERENCES users(id)
	FOREIGN KEY (post_id) 
		REFERENCES posts(id)
);


CREATE TABLE IF NOT EXISTS comments (
	id INTEGER PRIMARY KEY,
	post_id INTEGER NOT NULL,
	user_id INTEGER NOT NULL,
	content TEXT NOT NULL,
	image_path TEXT,
	created_at DATETIME NOT NULL,
	FOREIGN KEY (user_id) 
		REFERENCES users (id)
	FOREIGN KEY (post_id) 
		REFERENCES posts (id)
);