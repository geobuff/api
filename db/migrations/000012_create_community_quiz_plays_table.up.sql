CREATE TABLE communityquizplays (
    id SERIAL PRIMARY KEY,
    communityQuizId INTEGER references communityQuizzes(id) NOT NULL,
    plays INTEGER NOT NULL
);