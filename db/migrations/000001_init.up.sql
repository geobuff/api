CREATE TABLE badges (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    total INTEGER NOT NULL,
    imageUrl TEXT NOT NULL,
    background TEXT NOT NULL,
    border TEXT NOT NULL
);

CREATE table avatars (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    imageUrl TEXT NOT NULL,
    background TEXT NOT NULL,
    border TEXT NOT NULL
);

CREATE TABLE quiztype (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE quizzes (
    id SERIAL PRIMARY KEY,
    type INTEGER references quiztype(id) NOT NULL,
    badgeGroup INTEGER references badges(id),
    name TEXT NOT NULL,
    maxScore INTEGER NOT NULL,
    time INTEGER NOT NULL,
    mapSVG TEXT NOT NULL,
    imageUrl TEXT NOT NULL,
    verb TEXT NOT NULL,
    apiPath TEXT NOT NULL,
    route TEXT NOT NULL,
    hasLeaderboard BOOLEAN NOT NULL,
    hasGrouping BOOLEAN NOT NULL,
    hasFlags BOOLEAN NOT NULL,
    enabled BOOLEAN NOT NULL
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    avatarId INTEGER references avatars(id) NOT NULL,
    username TEXT UNIQUE NOT NULL,
    email TEXT NOT NULL,
    passwordHash TEXT NOT NULL,
    countryCode TEXT NOT NULL,
    xp BIGINT NOT NULL,
    isPremium BOOLEAN NOT NULL,
    stripeSessionId TEXT,
    isAdmin BOOLEAN NOT NULL,
    passwordResetToken TEXT,
    passwordResetExpiry DATE
);

CREATE TABLE tempscores (
    id SERIAL PRIMARY KEY,
    score INTEGER NOT NULL,
    time INTEGER NOT NULL,
    added DATE NOT NULL
);

CREATE TABLE leaderboard (
    id SERIAL PRIMARY KEY,
    quizId INTEGER references quizzes(id) NOT NULL,
    userId INTEGER references users(id) NOT NULL,
    score INTEGER NOT NULL, 
    time INTEGER NOT NULL,
    added DATE NOT NULL
);

CREATE TABLE plays (
    id SERIAL PRIMARY KEY,
    quizId INTEGER references quizzes(id) NOT NULL,
    value INTEGER NOT NULL
);

CREATE TABLE merch (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    price DECIMAL NOT NULL,
    disabled BOOLEAN NOT NULL
);

CREATE TABLE merchSizes (
    id SERIAL PRIMARY KEY,
    merchId INTEGER references merch(id) NOT NULL,
    size TEXT NOT NULL,
    soldOut BOOLEAN NOT NULL
);

CREATE TABLE merchImages (
    id SERIAL PRIMARY KEY,
    merchId INTEGER references merch(id) NOT NULL,
    imageUrl TEXT NOT NULL,
    isPrimary BOOLEAN NOT NULL
);

INSERT INTO badges (name, description, total, imageUrl, background, border) values
('Competitor', 'Submit a leaderboard entry.', 1, 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f3c6.svg', '#FFF1CE', '#C1694F'),
('International Traveler', 'Complete all world quizzes.', 3, 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f9f3.svg', '#A8D9FF', '#4289C1'),
('SaharanBuff', 'Complete all Africa quizzes.', 4, 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f418.svg', '#DDDDDD', '#66757F'),
('OrientalBuff', 'Complete all Asia quizzes.', 5, 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f338.svg', '#FFE4EA', '#EA596E'),
('EuropaBuff', 'Complete all Europe quizzes.', 5, 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f3f0.svg', '#ECF3F9', '#226699'),
('RockiesBuff', 'Complete all North America quizzes.', 6, 'https://twemoji.maxcdn.com/v/13.0.1/svg/26f0.svg', '#E9E9E9', '#4B545D'),
('AmazonBuff', 'Complete all South America quizzes.', 5, 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f483.svg', '#FFBFC7', '#A0041E'),
('PacificBuff', 'Complete all Oceania quizzes.', 4, 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f3dd.svg', '#D3ECFF', '#F4900C');

INSERT INTO avatars (name, imageUrl, background, border) values
('Mr. Pumpkin', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f383.svg', '#FFE0B7', '#F79D27'),
('Bolts', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f916.svg', '#E0F2FF', '#3B88C3'),
('GeoKitty', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f431.svg', '#FFE8AD', '#F18F26'),
('Cosmo', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f47d.svg', '#E6F2FB', '#CCD6DD'),
('Prof. Doo Doo', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f4a9.svg', '#EDCBC1', '#BF6952');

INSERT INTO quiztype (name) values ('Map'), ('Flag');

INSERT INTO quizzes (type, badgeGroup, name, maxScore, time, mapSVG, imageUrl, verb, apiPath, route, hasLeaderboard, hasGrouping, hasFlags, enabled) values
(1, 2, 'Countries of the World', 197, 900, 'WorldCountries', '/world-map-header.svg', 'countries', 'world-countries', 'countries-of-the-world', TRUE, TRUE, TRUE, TRUE),
(1, 2, 'Capitals of the World', 197, 900, 'WorldCapitals', '/world-map-header.svg', 'capitals', 'world-capitals', 'capitals-of-the-world', TRUE, TRUE, TRUE, TRUE),
(1, 3, 'Countries of Africa', 54, 300, 'AfricaCountries', '/africa-countries-header.svg', 'countries', 'africa-countries', 'countries-of-africa', TRUE, FALSE, TRUE, TRUE),
(1, 4, 'Countries of Asia', 50, 300, 'AsiaCountries', '/asia-countries-header.svg', 'countries', 'asia-countries', 'countries-of-asia', TRUE, FALSE, TRUE, TRUE),
(1, 5, 'Countries of Europe', 51, 300, 'EuropeCountries', '/europe-countries-header.svg', 'countries', 'europe-countries', 'countries-of-europe', TRUE, FALSE, TRUE, TRUE),
(1, 6, 'Countries of North America', 23, 300, 'NorthAmericaCountries', '/north-america-countries-header.svg', 'countries', 'north-america-countries', 'countries-of-north-america', TRUE, FALSE, TRUE, TRUE),
(1, 7, 'Countries of South America', 12, 300, 'SouthAmericaCountries', '/south-america-countries-header.svg', 'countries', 'south-america-countries', 'countries-of-south-america', TRUE, FALSE, TRUE, TRUE),
(1, 8, 'Countries of Oceania', 15, 300, 'OceaniaCountries', '/oceania-countries-header.svg', 'countries', 'oceania-countries', 'countries-of-oceania', TRUE, FALSE, TRUE, TRUE),
(1, 7, 'Provinces of Argentina', 24, 300, 'ArgentinaProvinces', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e6-1f1f7.svg', 'provinces', 'argentina-provinces', 'provinces-of-argentina', TRUE, FALSE, TRUE, TRUE),
(1, 8, 'States and Territories of Australia', 9, 300, 'AustraliaStatesAndTerritories', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e6-1f1fa.svg', 'states and territories', 'australia-states-and-territories', 'states-and-territories-of-australia', TRUE, FALSE, TRUE, TRUE),
(1, 7, 'States of Brazil', 27, 300, 'BrazilStates', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e7-1f1f7.svg', 'states', 'brazil-states', 'states-of-brazil', TRUE, FALSE, TRUE, TRUE),
(1, 6, 'Provinces and Territories of Canada', 13, 300, 'CanadaProvincesAndTerritories', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e8-1f1e6.svg', 'provinces and territories', 'canada-provinces-and-territories', 'provinces-and-territories-of-canada', TRUE, FALSE, TRUE, TRUE),
(1, 4, 'Administrative Divisions of China', 34, 300, 'ChinaAdministrativeDivisions', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e8-1f1f3.svg', 'administrative divisions', 'china-administrative-divisions', 'administrative-divisions-of-china', TRUE, FALSE, FALSE, TRUE),
(1, 5, 'Regions of France', 13, 300, 'FranceRegions', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1eb-1f1f7.svg', 'regions', 'france-regions', 'regions-of-france', TRUE, FALSE, TRUE, TRUE),
(1, 5, 'States of Germany', 16, 300, 'GermanyStates', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e9-1f1ea.svg', 'states', 'germany-states', 'states-of-germany', TRUE, FALSE, TRUE, TRUE),
(1, 4, 'States and Union Territories of India', 36, 300, 'IndiaStatesAndUnionTerritories', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ee-1f1f3.svg', 'states and union territories', 'india-states-and-union-territories', 'states-and-union-territories-of-india', TRUE, FALSE, FALSE, TRUE),
(1, 4, 'Prefectures of Japan', 47, 300, 'JapanPrefectures', 'https://upload.wikimedia.org/wikipedia/en/9/9e/Flag_of_Japan.svg', 'prefectures', 'japan-prefectures', 'prefectures-of-japan', TRUE, FALSE, TRUE, TRUE),
(1, 6, 'States of Mexico', 32, 300, 'MexicoStates', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f2-1f1fd.svg', 'states', 'mexico-states', 'states-of-mexico', TRUE, FALSE, FALSE, TRUE),
(1, 8, 'Regions of New Zealand', 16, 300, 'NewZealandRegions', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f3-1f1ff.svg', 'regions', 'new-zealand-regions', 'regions-of-new-zealand', TRUE, FALSE, FALSE, TRUE),
(1, 3, 'States of Nigeria', 37, 300, 'NigeriaStates', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f3-1f1ec.svg', 'states', 'nigeria-states', 'states-of-nigeria', TRUE, FALSE, FALSE, TRUE),
(1, 4, 'Provinces of South Korea', 17, 300, 'SouthKoreaProvinces', 'https://upload.wikimedia.org/wikipedia/commons/0/09/Flag_of_South_Korea.svg', 'provinces', 'south-korea-provinces', 'provinces-of-south-korea', TRUE, FALSE, TRUE, TRUE),
(1, 5, 'Provinces of Spain', 52, 300, 'SpainProvinces', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ea-1f1f8.svg', 'provinces', 'spain-provinces', 'provinces-of-spain', TRUE, FALSE, TRUE, TRUE),
(1, 4, 'Provinces of Turkey', 81, 600, 'TurkeyProvinces', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f9-1f1f7.svg', 'provinces', 'turkey-provinces', 'provinces-of-turkey', TRUE, FALSE, FALSE, TRUE),
(1, 6, 'US States', 50, 300, 'UsStates', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1fa-1f1f8.svg', 'states', 'us-states', 'us-states', TRUE, FALSE, TRUE, TRUE),
(1, 3, 'Districts of Uganda', 112, 600, 'UgandaDistricts', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1fa-1f1ec.svg', 'districts', 'uganda-districts', 'districts-of-uganda', TRUE, FALSE, FALSE, TRUE),
(1, 3, 'Provinces of Zambia', 10, 300, 'ZambiaProvinces', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ff-1f1f2.svg', 'provinces', 'zambia-provinces', 'provinces-of-zambia', TRUE, FALSE, FALSE, TRUE),
(2, 2, 'Flags of the World', 197, 900, '', '/world-map-header.svg', 'flags', 'world-countries', 'flags-of-the-world', TRUE, TRUE, TRUE, TRUE),
(2, 8, 'Flags of Australia', 9, 300, '', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e6-1f1fa.svg', 'flags', 'australia-states-and-territories', 'flags-of-australia', TRUE, FALSE, TRUE, TRUE),
(2, 7, 'Flags of Argentina', 24, 300, '', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e6-1f1f7.svg', 'flags', 'argentina-provinces', 'flags-of-argentina', TRUE, FALSE, TRUE, TRUE),
(2, 7, 'Flags of Brazil', 27, 300, '', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e7-1f1f7.svg', 'flags', 'brazil-states', 'flags-of-brazil', TRUE, FALSE, TRUE, TRUE),
(2, 6, 'Flags of Canada', 13, 300, '', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e8-1f1e6.svg', 'flags', 'canada-provinces-and-territories', 'flags-of-canada', TRUE, FALSE, TRUE, TRUE),
(2, 5, 'Flags of France', 13, 300, '', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1eb-1f1f7.svg', 'flags', 'france-regions', 'flags-of-france', TRUE, FALSE, TRUE, TRUE),
(2, 5, 'Flags of Germany', 16, 300, '', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e9-1f1ea.svg', 'flags', 'germany-states', 'flags-of-germany', TRUE, FALSE, TRUE, TRUE),
(2, 4, 'Flags of Japan', 47, 300, '', 'https://upload.wikimedia.org/wikipedia/en/9/9e/Flag_of_Japan.svg', 'flags', 'japan-prefectures', 'flags-of-japan', TRUE, FALSE, TRUE, TRUE),
(2, 4, 'Flags of South Korea', 17, 300, '', 'https://upload.wikimedia.org/wikipedia/commons/0/09/Flag_of_South_Korea.svg', 'flags', 'south-korea-provinces', 'flags-of-south-korea', TRUE, FALSE, TRUE, TRUE),
(2, 5, 'Flags of Spain', 52, 300, '', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ea-1f1f8.svg', 'flags', 'spain-provinces', 'flags-of-spain', TRUE, FALSE, TRUE, TRUE),
(2, 6, 'Flags of the US', 50, 300, '', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1fa-1f1f8.svg', 'flags', 'us-states', 'flags-of-the-us', TRUE, FALSE, TRUE, TRUE);

INSERT INTO plays (quizId, value) values
(1, 0),
(2, 0),
(3, 0),
(4, 0),
(5, 0),
(6, 0),
(7, 0),
(8, 0),
(9, 0),
(10, 0),
(11, 0),
(12, 0),
(13, 0),
(14, 0),
(15, 0),
(16, 0),
(17, 0),
(18, 0),
(19, 0),
(20, 0),
(21, 0),
(22, 0),
(23, 0),
(24, 0),
(25, 0),
(26, 0),
(27, 0),
(28, 0),
(29, 0),
(30, 0),
(31, 0),
(32, 0),
(33, 0),
(34, 0),
(35, 0),
(36, 0);

INSERT INTO merch (name, description, price, disabled) values
('GeoTee', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.', 69.99, FALSE),
('GeoCap', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.', 39.99, FALSE);

INSERT INTO merchSizes (merchId, size, soldOut) values
(1, 'S', FALSE),
(1, 'M', FALSE),
(1, 'L', FALSE),
(1, 'XL', FALSE),
(2, 'One Size Fits All', FALSE);

INSERT INTO merchImages (merchId, imageUrl, isPrimary) values
(1, '/geotee.png', TRUE),
(2, '/geocap.png', TRUE);