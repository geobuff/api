CREATE TABLE badges (
    id SERIAL PRIMARY KEY,
    name text NOT NULL,
    description text NOT NULL,
    icon text NOT NULL,
    total INTEGER NOT NULL
);

CREATE TABLE quiztype (
    id SERIAL PRIMARY KEY,
    name text UNIQUE NOT NULL
);

CREATE TABLE quizzes (
    id SERIAL PRIMARY KEY,
    type INTEGER references quiztype(id) NOT NULL,
    badgeGroup INTEGER references badges(id),
    name text NOT NULL,
    maxScore INTEGER NOT NULL,
    time INTEGER NOT NULL,
    mapSVG text NOT NULL,
    imageUrl text NOT NULL,
    verb text NOT NULL,
    apiPath text NOT NULL,
    route text NOT NULL,
    hasLeaderboard BOOLEAN NOT NULL,
    hasGrouping BOOLEAN NOT NULL,
    hasFlags BOOLEAN NOT NULL,
    enabled BOOLEAN NOT NULL
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY, 
    username text UNIQUE NOT NULL,
    countryCode text,
    xp BIGINT
);

CREATE TABLE keys (
    id SERIAL PRIMARY KEY,
    name text NOT NULL,
    key text UNIQUE NOT NULL
);

CREATE TABLE scores (
    id SERIAL PRIMARY KEY,
    userId INTEGER references users(id) NOT NULL,
    quizId INTEGER references quizzes(id) NOT NULL,
    score INTEGER NOT NULL,
    time INTEGER NOT NULL,
    added DATE NOT NULL
);

CREATE TABLE countries_leaderboard (
    id SERIAL PRIMARY KEY, 
    userId INTEGER references users(id) NOT NULL,
    score INTEGER NOT NULL, 
    time INTEGER NOT NULL,
    added DATE NOT NULL
);

CREATE TABLE capitals_leaderboard (
    id SERIAL PRIMARY KEY, 
    userId INTEGER references users(id) NOT NULL,
    score INTEGER NOT NULL, 
    time INTEGER NOT NULL,
    added DATE NOT NULL
);