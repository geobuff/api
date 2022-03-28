ALTER TABLE avatars ADD column typeId INTEGER references avatarTypes(id) NOT NULL DEFAULT 1;
ALTER TABLE avatars ADD column countryCode TEXT NOT NULL DEFAULT 'nz';