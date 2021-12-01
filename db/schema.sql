CREATE TABLE users (
  userid     TEXT NOT NULL,
  first_name TEXT NOT NULL,
  last_name  TEXT NOT NULL,
  data       TEXT NOT NULL, -- json data
  hidden     INTEGER NOT NULL CHECK (hidden IN (0, 1)), -- 0=false 1=true
  UNIQUE(first_name, last_name)
);

CREATE TABLE attendance (
    userid TEXT NOT NULL,
    date   TEXT NOT NULL,
    UNIQUE(userid, date)
);