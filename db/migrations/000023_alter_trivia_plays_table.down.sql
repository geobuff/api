ALTER TABLE triviaplays ALTER COLUMN triviaid SET NOT NULL;
ALTER TABLE triviaplays ADD CONSTRAINT triviaplays_triviaid_fkey FOREIGN KEY (triviaid) REFERENCES trivia(id);
