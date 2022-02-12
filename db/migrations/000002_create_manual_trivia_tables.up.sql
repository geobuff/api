CREATE TABLE manualTriviaQuestions (
    id SERIAL PRIMARY KEY,
    typeId INTEGER references triviaQuestionType(id) NOT NULL,
    question TEXT NOT NULL,
    map TEXT,
    highlighted TEXT,
    flagCode TEXT,
    imageUrl TEXT
);

CREATE TABLE manualTriviaAnswers (
    id SERIAL PRIMARY KEY,
    manualTriviaQuestionId INTEGER references manualTriviaQuestions(id) NOT NULL,
    text TEXT NOT NULL,
    isCorrect BOOLEAN NOT NULL,
    flagCode TEXT
);