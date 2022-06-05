ALTER TABLE communityquizplays DROP CONSTRAINT communityquizplays_communityquizid_fkey;
ALTER TABLE communityquizplays ALTER COLUMN communityquizid DROP NOT NULL;
