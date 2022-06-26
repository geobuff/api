ALTER TABLE triviaplays ALTER COLUMN triviaid SET NOT NULL;
/* 
NOTE: Removed below command as not possible to reintroduce contraint after trivia that matches triviaId has been removed from trivia table.  
ALTER TABLE triviaplays ADD CONSTRAINT triviaplays_triviaid_fkey FOREIGN KEY (triviaid) REFERENCES trivia(id);
*/
