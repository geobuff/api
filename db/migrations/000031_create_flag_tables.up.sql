CREATE TABLE flagGroups (
    id SERIAL PRIMARY KEY,
    key TEXT UNIQUE NOT NULL,
    label TEXT NOT NULL
);

CREATE TABLE flagEntries (
    id SERIAL PRIMARY KEY,
    groupId INTEGER references flagGroups(id) NOT NULL,
    code TEXT UNIQUE NOT NULL,
    url TEXT NOT NULL
);
