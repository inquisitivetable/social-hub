ALTER TABLE posts DROP COLUMN title;

ALTER TABLE posts 
ADD COLUMN group_id INTEGER 
REFERENCES groups(id);


ALTER TABLE groups
ADD COLUMN image_path TEXT;


ALTER TABLE followers DROP COLUMN active;