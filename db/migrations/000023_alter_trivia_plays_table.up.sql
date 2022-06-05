ALTER TABLE triviaplays DROP CONSTRAINT triviaplays_triviaid_fkey;
ALTER TABLE triviaplays ALTER COLUMN triviaid DROP NOT NULL;
