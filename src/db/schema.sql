CREATE TABLE users (
  userid     TEXT NOT NULL,
  first_name TEXT NOT NULL,
  last_name  TEXT NOT NULL,
  data       TEXT NOT NULL, -- json data
  UNIQUE(first_name, last_name)
);

CREATE TABLE attendance (
    userid TEXT NOT NULL,
    date   TEXT NOT NULL,
    UNIQUE(userid, date)
);