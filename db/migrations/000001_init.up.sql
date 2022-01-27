CREATE TABLE continents (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE badgeType (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE badges (
    id SERIAL PRIMARY KEY,
    typeId INTEGER references badgeType(id) NOT NULL,
    continentId INTEGER references continents(id),
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    imageUrl TEXT NOT NULL,
    background TEXT NOT NULL,
    border TEXT NOT NULL
);

CREATE table avatars (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    primaryImageUrl TEXT NOT NULL,
    secondaryImageUrl TEXT NOT NULL
);

CREATE TABLE quiztype (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE quizzes (
    id SERIAL PRIMARY KEY,
    typeId INTEGER references quiztype(id) NOT NULL,
    badgeId INTEGER references badges(id),
    continentId INTEGER references continents(id),
    country TEXT NOT NULL,
    singular TEXT NOT NULL,
    name TEXT NOT NULL,
    maxScore INTEGER NOT NULL,
    time INTEGER NOT NULL,
    mapSVG TEXT NOT NULL,
    imageUrl TEXT NOT NULL,
    verb TEXT NOT NULL,
    apiPath TEXT NOT NULL,
    route TEXT NOT NULL,
    hasLeaderboard BOOLEAN NOT NULL,
    hasGrouping BOOLEAN NOT NULL,
    hasFlags BOOLEAN NOT NULL,
    enabled BOOLEAN NOT NULL
);

CREATE TABLE quizPlays (
    id SERIAL PRIMARY KEY,
    quizId INTEGER references quizzes(id) NOT NULL,
    plays INTEGER NOT NULL
);

CREATE TABLE trivia (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    date DATE NOT NULL
);

CREATE TABLE triviaQuestionType (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE triviaQuestions (
    id SERIAL PRIMARY KEY,
    triviaId INTEGER references trivia(id) NOT NULL,
    typeId INTEGER references triviaQuestionType(id) NOT NULL,
    question TEXT NOT NULL,
    map TEXT,
    highlighted TEXT,
    flagCode TEXT,
    imageUrl TEXT
);

CREATE TABLE triviaAnswers (
    id SERIAL PRIMARY KEY,
    triviaQuestionId INTEGER references triviaQuestions(id) NOT NULL,
    text TEXT NOT NULL,
    isCorrect BOOLEAN NOT NULL,
    flagCode TEXT
);

CREATE TABLE triviaPlays (
    id SERIAL PRIMARY KEY,
    triviaId INTEGER references trivia(id) NOT NULL,
    plays INTEGER NOT NULL
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    avatarId INTEGER references avatars(id) NOT NULL,
    username TEXT UNIQUE NOT NULL,
    email TEXT NOT NULL,
    passwordHash TEXT NOT NULL,
    countryCode TEXT NOT NULL,
    xp BIGINT NOT NULL,
    isPremium BOOLEAN NOT NULL,
    isAdmin BOOLEAN NOT NULL,
    passwordResetToken TEXT,
    passwordResetExpiry DATE,
    joined DATE
);

CREATE TABLE tempscores (
    id SERIAL PRIMARY KEY,
    score INTEGER NOT NULL,
    time INTEGER NOT NULL,
    results TEXT[] NOT NULL,
    recents TEXT[] NOT NULL,
    added DATE NOT NULL
);

CREATE TABLE leaderboard (
    id SERIAL PRIMARY KEY,
    quizId INTEGER references quizzes(id) NOT NULL,
    userId INTEGER references users(id) NOT NULL,
    score INTEGER NOT NULL, 
    time INTEGER NOT NULL,
    added DATE NOT NULL
);

CREATE TABLE merch (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    sizeGuideImageUrl TEXT,
    price DECIMAL(12,2),
    externalLink TEXT
);

CREATE TABLE merchSizes (
    id SERIAL PRIMARY KEY,
    merchId INTEGER references merch(id) NOT NULL,
    size TEXT NOT NULL,
    quantity INTEGER NOT NULL
);

CREATE TABLE merchImages (
    id SERIAL PRIMARY KEY,
    merchId INTEGER references merch(id) NOT NULL,
    imageUrl TEXT NOT NULL,
    isPrimary BOOLEAN NOT NULL
);

CREATE TABLE discounts (
    id SERIAL PRIMARY KEY,
    merchId INTEGER references merch(id),
    code TEXT NOT NULL,
    amount DECIMAL(12,2) NOT NULL
);

CREATE TABLE orderStatus (
    id SERIAL PRIMARY KEY,
    status TEXT NOT NULL
);

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    statusId INTEGER references orderStatus(id),
    email TEXT NOT NULL,
    firstName TEXT NOT NULL,
    lastName TEXT NOT NULL,
    address TEXT NOT NULL,
    suburb TEXT NOT NULL,
    city TEXT NOT NULL,
    postcode TEXT NOT NULL,
    added DATE NOT NULL,
    discount TEXT
);

CREATE TABLE orderItems (
    id SERIAL PRIMARY KEY,
    orderId INTEGER references orders(id),
    merchId INTEGER references merch(id),
    sizeId INTEGER references merchSizes(id),
    quantity INTEGER NOT NULL
);