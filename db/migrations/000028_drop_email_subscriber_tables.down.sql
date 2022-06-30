CREATE TABLE triviaReminderSubscribers (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL,
    added DATE NOT NULL
);

CREATE TABLE newsletterSubscribers (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL,
    added DATE NOT NULL
);
