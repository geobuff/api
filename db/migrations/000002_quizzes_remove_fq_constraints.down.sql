ALTER TABLE quizzes
    ADD CONSTRAINT quizzes_badgeid_fkey FOREIGN KEY (badgeid) REFERENCES badges(id);

ALTER TABLE quizzes
    ADD CONSTRAINT quizzes_continentid_fkey FOREIGN KEY (continentid) REFERENCES continents(id);