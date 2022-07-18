ALTER TABLE communityquizplays ALTER COLUMN communityquizid SET NOT NULL;
/* 
NOTE: Removed below command as not possible to reintroduce contraint after communityquiz that matches communityquizId has been removed from communityquizzes table. 
ALTER TABLE communityquizplays ADD CONSTRAINT communityquizplays_communityquizid_fkey FOREIGN KEY (communityquizid) REFERENCES communityquizzes(id);
*/
