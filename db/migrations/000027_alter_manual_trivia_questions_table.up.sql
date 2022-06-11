ALTER TABLE manualtriviaquestions ADD column categoryId INTEGER references triviaquestioncategory(id) NOT NULL DEFAULT 1;
