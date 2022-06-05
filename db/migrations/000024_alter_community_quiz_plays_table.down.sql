ALTER TABLE communityquizplays ALTER COLUMN communityquizid SET NOT NULL;
ALTER TABLE communityquizplays ADD CONSTRAINT communityquizplays_communityquizid_fkey FOREIGN KEY (communityquizid) REFERENCES communityquizzes(id);
