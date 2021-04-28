CREATE TABLE badges (
    id SERIAL PRIMARY KEY,
    name text NOT NULL,
    description text NOT NULL,
    icon text NOT NULL,
    total INTEGER NOT NULL
);

CREATE TABLE quiztype (
    id SERIAL PRIMARY KEY,
    name text UNIQUE NOT NULL
);

CREATE TABLE quizzes (
    id SERIAL PRIMARY KEY,
    type INTEGER references quiztype(id) NOT NULL,
    badgeGroup INTEGER references badges(id),
    name text NOT NULL,
    maxScore INTEGER NOT NULL,
    time INTEGER NOT NULL,
    mapSVG text NOT NULL,
    imageUrl text NOT NULL,
    verb text NOT NULL,
    apiPath text NOT NULL,
    route text NOT NULL,
    hasLeaderboard BOOLEAN NOT NULL,
    hasGrouping BOOLEAN NOT NULL,
    hasFlags BOOLEAN NOT NULL,
    enabled BOOLEAN NOT NULL
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY, 
    username text UNIQUE NOT NULL,
    email text NOT NULL,
    passwordHash text NOT NULL,
    countryCode text NOT NULL,
    xp BIGINT NOT NULL,
    isPremium BOOLEAN NOT NULL,
    isAdmin BOOLEAN NOT NULL
);

CREATE TABLE scores (
    id SERIAL PRIMARY KEY,
    userId INTEGER references users(id) NOT NULL,
    quizId INTEGER references quizzes(id) NOT NULL,
    score INTEGER NOT NULL,
    time INTEGER NOT NULL,
    added DATE NOT NULL
);

CREATE TABLE countries_leaderboard (
    id SERIAL PRIMARY KEY, 
    userId INTEGER references users(id) NOT NULL,
    score INTEGER NOT NULL, 
    time INTEGER NOT NULL,
    added DATE NOT NULL
);

CREATE TABLE capitals_leaderboard (
    id SERIAL PRIMARY KEY, 
    userId INTEGER references users(id) NOT NULL,
    score INTEGER NOT NULL, 
    time INTEGER NOT NULL,
    added DATE NOT NULL
);

INSERT INTO badges (name, description, icon, total) values
('Cartographer', 'Set country in user profile.', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f4cd.svg', 1),
('Competitor', 'Submit a leaderboard entry.', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f3c6.svg', 1),
('International Traveler', 'Complete all world quizzes.', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f9f3.svg', 2),
('SaharanBuff', 'Complete all Africa quizzes.', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f418.svg', 4),
('OrientalBuff', 'Complete all Asia quizzes.', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f338.svg', 4),
('EuropaBuff', 'Complete all Europe quizzes.', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f3f0.svg', 4),
('RockiesBuff', 'Complete all North America quizzes.', 'https://twemoji.maxcdn.com/v/13.0.1/svg/26f0.svg', 4),
('AmazonBuff', 'Complete all South America quizzes.', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f483.svg', 3),
('PacificBuff', 'Complete all Oceania quizzes.', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f3dd.svg', 3);

INSERT INTO quiztype (name) values ('MAP'), ('FLAG');

INSERT INTO quizzes (type, badgeGroup, name, maxScore, time, mapSVG, imageUrl, verb, apiPath, route, hasLeaderboard, hasGrouping, hasFlags, enabled) values
(1, 3, 'Countries of the World', 197, 900, 'WorldCountries', '/world-map-header.svg', 'countries', 'world-countries', 'countries-of-the-world', TRUE, TRUE, TRUE, TRUE),
(1, 3, 'Capitals of the World', 197, 900, 'WorldCapitals', '/world-map-header.svg', 'capitals', 'world-capitals', 'capitals-of-the-world', TRUE, TRUE, TRUE, TRUE),
(1, 4, 'Countries of Africa', 54, 300, 'AfricaCountries', '/africa-countries-header.svg', 'countries', 'africa-countries', 'countries-of-africa', FALSE, FALSE, TRUE, TRUE),
(1, 5, 'Countries of Asia', 50, 300, 'AsiaCountries', '/asia-countries-header.svg', 'countries', 'asia-countries', 'countries-of-asia', FALSE, FALSE, TRUE, TRUE),
(1, 6, 'Countries of Europe', 51, 300, 'EuropeCountries', '/europe-countries-header.svg', 'countries', 'europe-countries', 'countries-of-europe', FALSE, FALSE, TRUE, TRUE),
(1, 7, 'Countries of North America', 23, 300, 'NorthAmericaCountries', '/north-america-countries-header.svg', 'countries', 'north-america-countries', 'countries-of-north-america', FALSE, FALSE, TRUE, TRUE),
(1, 8, 'Countries of South America', 12, 300, 'SouthAmericaCountries', '/south-america-countries-header.svg', 'countries', 'south-america-countries', 'countries-of-south-america', FALSE, FALSE, TRUE, TRUE),
(1, 9, 'Countries of Oceania', 15, 300, 'OceaniaCountries', '/oceania-countries-header.svg', 'countries', 'oceania-countries', 'countries-of-oceania', FALSE, FALSE, TRUE, TRUE),
(1, 8, 'Provinces of Argentina', 24, 300, 'ArgentinaProvinces', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e6-1f1f7.svg', 'provinces', 'argentina-provinces', 'provinces-of-argentina', FALSE, FALSE, TRUE, TRUE),
(1, 9, 'States and Territories of Australia', 9, 300, 'AustraliaStatesAndTerritories', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e6-1f1fa.svg', 'states and territories', 'australia-states-and-territories', 'states-and-territories-of-australia', FALSE, FALSE, TRUE, TRUE),
(1, 8, 'States of Brazil', 27, 300, 'BrazilStates', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e7-1f1f7.svg', 'states', 'brazil-states', 'states-of-brazil', FALSE, FALSE, TRUE, TRUE),
(1, 7, 'Provinces and Territories of Canada', 13, 300, 'CanadaProvincesAndTerritories', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e8-1f1e6.svg', 'provinces and territories', 'canada-provinces-and-territories', 'provinces-and-territories-of-canada', FALSE, FALSE, TRUE, TRUE),
(1, 6, 'Regions of France', 13, 300, 'FranceRegions', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1eb-1f1f7.svg', 'regions', 'france-regions', 'regions-of-france', FALSE, FALSE, TRUE, TRUE),
(1, 6, 'States of Germany', 16, 300, 'GermanyStates', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e9-1f1ea.svg', 'states', 'germany-states', 'states-of-germany', FALSE, FALSE, TRUE, TRUE),
(1, 5, 'States and Union Territories of India', 36, 300, 'IndiaStatesAndUnionTerritories', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ee-1f1f3.svg', 'states and union territories', 'india-states-and-union-territories', 'states-and-union-territories-of-india', FALSE, FALSE, FALSE, TRUE),
(1, 5, 'Prefectures of Japan', 47, 300, 'JapanPrefectures', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ef-1f1f5.svg', 'prefectures', 'japan-prefectures', 'prefectures-of-japan', FALSE, FALSE, TRUE, TRUE),
(1, 7, 'States of Mexico', 32, 300, 'MexicoStates', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f2-1f1fd.svg', 'states', 'mexico-states', 'states-of-mexico', FALSE, FALSE, FALSE, TRUE),
(1, 9, 'Regions of New Zealand', 16, 300, 'NewZealandRegions', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f3-1f1ff.svg', 'regions', 'new-zealand-regions', 'regions-of-new-zealand', FALSE, FALSE, FALSE, TRUE),
(1, 4, 'States of Nigeria', 37, 300, 'NigeriaStates', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f3-1f1ec.svg', 'states', 'nigeria-states', 'states-of-nigeria', FALSE, FALSE, FALSE, TRUE),
(1, 5, 'Provinces of Turkey', 81, 600, 'TurkeyProvinces', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f9-1f1f7.svg', 'provinces', 'turkey-provinces', 'provinces-of-turkey', FALSE, FALSE, FALSE, TRUE),
(1, 7, 'US States', 50, 300, 'UsStates', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1fa-1f1f8.svg', 'states', 'us-states', 'us-states', FALSE, FALSE, TRUE, TRUE),
(1, 4, 'Districts of Uganda', 112, 600, 'UgandaDistricts', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1fa-1f1ec.svg', 'districts', 'uganda-districts', 'districts-of-uganda', FALSE, FALSE, FALSE, TRUE),
(1, 4, 'Provinces of Zambia', 10, 300, 'ZambiaProvinces', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ff-1f1f2.svg', 'provinces', 'zambia-provinces', 'provinces-of-zambia', FALSE, FALSE, FALSE, TRUE),
(2, 3, 'Flags of the World', 197, 900, '', '/world-map-header.svg', 'flags', 'world-countries', 'flags-of-the-world', TRUE, TRUE, TRUE, TRUE),
(2, 9, 'Flags of Australia', 9, 300, '', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e6-1f1fa.svg', 'flags', 'australia-states-and-territories', 'flags-of-australia', FALSE, FALSE, TRUE, TRUE),
(2, 8, 'Flags of Argentina', 24, 300, '', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e6-1f1f7.svg', 'flags', 'argentina-provinces', 'flags-of-argentina', FALSE, FALSE, TRUE, TRUE),
(2, 8, 'Flags of Brazil', 27, 300, '', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e7-1f1f7.svg', 'flags', 'brazil-states', 'flags-of-brazil', FALSE, FALSE, TRUE, TRUE),
(2, 7, 'Flags of Canada', 13, 300, '', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e8-1f1e6.svg', 'flags', 'canada-provinces-and-territories', 'flags-of-canada', FALSE, FALSE, TRUE, TRUE),
(2, 6, 'Flags of France', 13, 300, '', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1eb-1f1f7.svg', 'flags', 'france-regions', 'flags-of-france', FALSE, FALSE, TRUE, TRUE),
(2, 6, 'Flags of Germany', 16, 300, '', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e9-1f1ea.svg', 'flags', 'germany-states', 'flags-of-germany', FALSE, FALSE, TRUE, TRUE),
(2, 5, 'Flags of Japan', 47, 300, '', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ef-1f1f5.svg', 'flags', 'japan-prefectures', 'flags-of-japan', FALSE, FALSE, TRUE, TRUE),
(2, 7, 'Flags of the US', 50, 300, '', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1fa-1f1f8.svg', 'flags', 'us-states', 'flags-of-the-us', FALSE, FALSE, TRUE, TRUE);