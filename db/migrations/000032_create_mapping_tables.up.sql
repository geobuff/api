CREATE TABLE mappingGroups (
    id SERIAL PRIMARY KEY,
    key TEXT UNIQUE NOT NULL,
    label TEXT NOT NULL
);

CREATE TABLE mappingEntries (
    id SERIAL PRIMARY KEY,
    groupId INTEGER references mappingGroups(id) NOT NULL,
    name TEXT NOT NULL,
    code TEXT NOT NULL,
    svgName TEXT NOT NULL,
    alternativeNames TEXT[] NOT NULL,
    prefixes TEXT[] NOT NULL,
    grouping TEXT NOT NULL
);
