ALTER TABLE posts ADD COLUMN title;

ALTER TABLE posts DROP COLUMN group_id;

ALTER TABLE groups DROP COLUMN image_path;

ALTER TABLE followers ADD COLUMN active BOOL;