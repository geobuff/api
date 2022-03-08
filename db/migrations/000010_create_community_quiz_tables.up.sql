CREATE TABLE communityQuizzes (
    id SERIAL PRIMARY KEY,
    userId INTEGER references users(id) NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    maxScore INTEGER NOT NULL,
    added DATE NOT NULL
);

CREATE TABLE communityQuizQuestions (
    id SERIAL PRIMARY KEY,
    communityQuizId INTEGER references communityQuizzes(id) NOT NULL,
    typeId INTEGER references triviaQuestionType(id) NOT NULL,
    question TEXT NOT NULL,
    map TEXT,
    highlighted TEXT,
    flagCode TEXT,
    imageUrl TEXT
);

CREATE TABLE communityQuizAnswers (
    id SERIAL PRIMARY KEY,
    communityQuizQuestionId INTEGER references communityquizquestions(id) NOT NULL,
    text TEXT NOT NULL,
    isCorrect BOOLEAN NOT NULL,
    flagCode TEXT
);