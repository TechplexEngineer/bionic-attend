ALTER TABLE users
	ADD COLUMN hidden INTEGER
		NOT NULL
		DEFAULT FALSE
		CHECK (hidden IN (0, 1)); --0=false 1=true;