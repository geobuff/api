CREATE TABLE maps (
    id SERIAL PRIMARY KEY,
    key TEXT UNIQUE NOT NULL,
    className TEXT UNIQUE NOT NULL,
    label TEXT NOT NULL,
    viewBox TEXT NOT NULL
);

CREATE TABLE mapElementType (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE mapElements (
    id SERIAL PRIMARY KEY,
    mapId INTEGER references maps(id) NOT NULL,
    typeId INTEGER references mapElementType(id) NOT NULL,
    elementId TEXT NOT NULL,
    name TEXT NOT NULL,
    d TEXT NOT NULL,
    points TEXT NOT NULL,
    x TEXT NOT NULL,
    y TEXT NOT NULL,
    width TEXT NOT NULL,
    height TEXT NOT NULL,
    cx TEXT NOT NULL,
    cy TEXT NOT NULL,
    r TEXT NOT NULL,
    transform TEXT NOT NULL,
    xlinkHref TEXT NOT NULL,
    clipPath TEXT NOT NULL,
    clipPathId TEXT NOT NULL,
    x1 TEXT NOT NULL,
    y1 TEXT NOT NULL,
    x2 TEXT NOT NULL,
    y2 TEXT NOT NULL
);
