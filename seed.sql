INSERT INTO continents (name) values
('Africa'),
('Asia'),
('Europe'),
('North America'),
('South America'),
('Oceania');

INSERT INTO badgetype (name) values
('Leaderboard Submit'),
('World'),
('Continent'),
('Community Quiz');

INSERT INTO badges (typeid, continentId, name, description, imageUrl, background, border) values
(1, null, 'Competitor', 'Submit a leaderboard entry.', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f3c6.svg', '#FFF1CE', '#C1694F'),
(2, null, 'International Traveler', 'Complete all world quizzes.', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f9f3.svg', '#A8D9FF', '#4289C1'),
(3, 1, 'SaharanBuff', 'Complete all Africa quizzes.', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f418.svg', '#DDDDDD', '#66757F'),
(3, 2, 'OrientalBuff', 'Complete all Asia quizzes.', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f338.svg', '#FFE4EA', '#EA596E'),
(3, 3, 'EuropaBuff', 'Complete all Europe quizzes.', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f3f0.svg', '#ECF3F9', '#226699'),
(3, 4, 'RockiesBuff', 'Complete all North America quizzes.', 'https://twemoji.maxcdn.com/v/13.0.1/svg/26f0.svg', '#E9E9E9', '#4B545D'),
(3, 5, 'AmazonBuff', 'Complete all South America quizzes.', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f483.svg', '#FFBFC7', '#A0041E'),
(3, 6, 'PacificBuff', 'Complete all Oceania quizzes.', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f3dd.svg', '#D3ECFF', '#F4900C'),
(4, null, 'Sharing is Caring', 'Create a community quiz.', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f91d.svg', '#276f86', '#000000');

INSERT INTO avatarTypes (name) values
('Commando'),
('Traveller'),
('Researcher');

INSERT INTO avatars (typeId, countryCode, name, description, primaryImageUrl, secondaryImageUrl, gridPlacement) values
(1, 'nz', 'Sarg', 'After a 20 year stint in the NZSAS, Sarg spent some time smuggling precious stones into all those places that are on the tip of your tongue and rather difficult to spell. Cool, calm, collected and very handy with a zipwire; I keep my distance once this man''s on a Countries of the World roll. I advise you do the same...', '/avatars/commando-one-primary.svg', '/avatars/commando-one-secondary.svg', 4),
(1, 'us', 'Cypher', 'Bona fide hacker and rare NFT collector Cypher cut her teeth toying with Wall St. bankers before becoming the head of operations at the FBI. Word on the street is she lost her eye at 16 when a pirated mint of the original CryptoPunks went horribly wrong. Get a little bit too close to her score on the leaderboard and watch it get ''adjusted'' to zero in no time.', '/avatars/commando-two-primary.svg', '/avatars/commando-two-secondary.svg', 7),
(2, 'au', 'RoboGazza', 'Is he a man or is he a robot? No-one living knows. With a particular aptitude for shooting things with a bow and bombing hills on his bike, Gazza isn''t a man you want on your tail. Well travelled and well-seasoned; keep this in mind before going head-to-head with him on Flags of the World.', '/avatars/traveller-one-primary.svg', '/avatars/traveller-one-secondary.svg', 8),
(2, 'gb-eng', 'Kathy2.0', 'A Yorkshire clone resulting from an experiment run by GeoBuff Research Institute at an undisclosed location off the British coastline. The intent was to genetically increase VO2 max and number of steps per day, and increase it we did. Kathy has walked the Pennines start to finish 7 times this year. Speaking of, has anyone seen the original Kathy?', '/avatars/traveller-two-primary.svg', '/avatars/traveller-two-secondary.svg', 2),
(3, 'zm', 'Professor Lungu', 'Top of his class at MIT and hell-bent on being the only academic to publish a paper for the esteemed school AND score max points on Countries of the World in under 5 minutes. Sharp as a tack and quick as a whippet; don''t underestimate the power of a nerd on a mission.', '/avatars/researcher-one-primary.svg', '/avatars/researcher-one-secondary.svg', 3),
(3, 'mx', 'Sanchez', 'Ex-champion Jarabe dancer in her home Mexico City, Sanchez was taken under Prof. Lungu''s wing as a field researcher after he discovered she has a photographic memory and is the only person alive that can spell Kyrgyzstan correctly on the first try. Dance with Miss Sanchez, and well, it may just be your last...', '/avatars/researcher-two-primary.svg', '/avatars/researcher-two-secondary.svg', 6),
(1, 'fj', 'Willy', 'Born and raised on the Coral Coast, Willy represented Fiji internationally in Sevens before hanging up his boots to pursue a captaincy on the Sea Shepherd chasing pirates all across the Pacific. Strong as an ox and deeply knowledgable of the oceanography of the seven seas, give this man a hearty Bula next time you see him, otherwise it may be time to walk the plank...', '/avatars/commando-three-primary.svg', '/avatars/commando-three-secondary.svg', 1),
(2, 'cn', 'Ling', 'As an orphan who spent her formattive years stowed away on the Orient Express, Ling developed a particular talent for pickpocketing and mastering foreign languages. Forever on the move sharing her spoils of war with the street children of Shanghai through to Cairo, keep your eyes peeled for a train conductor donning a Jade exterior. It may be the last you see of that expensive wristwatch...', '/avatars/traveller-three-primary.svg', '/avatars/traveller-three-secondary.svg', 5),
(3, 'gb-sct', 'John John', 'They say good things come in small packages but JJ McClay is a reminder that sometime the old adages don''t always prove fruitful (It seems Dunkeld haggis is particularly vitamin-heavy). This 6''8 giant-of-a-man rounds out Prof Lungu''s team of researchers as the formitive expert on Climatology. When he''s not out on assignment in a faraway land, you''ll find him bird-watching in the Cairngorms or serving up bass-thumping tunes in the underground clubs of Edinburgh or Glasgow.', '/avatars/researcher-three-primary.svg', '/avatars/researcher-three-secondary.svg', 9);

INSERT INTO quiztype (name) values
('Map'),
('Flag');

INSERT INTO quizzes (typeId, badgeId, continentId, country, singular, name, maxScore, time, mapSVG, imageUrl, plural, apiPath, route, hasLeaderboard, hasGrouping, hasFlags, enabled) values
(1, 2, null, '', 'country', 'Countries of the World', 197, 900, 'WorldCountries', '/headers/world-map-header.svg', 'countries', 'world-countries', 'countries-of-the-world', TRUE, TRUE, TRUE, TRUE),
(1, 2, null, '', 'capital', 'Capitals of the World', 197, 900, 'WorldCapitals', '/headers/world-map-header.svg', 'capitals', 'world-capitals', 'capitals-of-the-world', TRUE, TRUE, TRUE, TRUE),
(1, 3, 1, '', 'country', 'Countries of Africa', 54, 300, 'AfricaCountries', '/headers/africa-countries-header.svg', 'countries', 'africa-countries', 'countries-of-africa', TRUE, FALSE, TRUE, TRUE),
(1, 4, 2, '', 'country', 'Countries of Asia', 50, 300, 'AsiaCountries', '/headers/asia-countries-header.svg', 'countries', 'asia-countries', 'countries-of-asia', TRUE, FALSE, TRUE, TRUE),
(1, 5, 3, '', 'country', 'Countries of Europe', 51, 300, 'EuropeCountries', '/headers/europe-countries-header.svg', 'countries', 'europe-countries', 'countries-of-europe', TRUE, FALSE, TRUE, TRUE),
(1, 6, 4, '', 'country', 'Countries of North America', 23, 300, 'NorthAmericaCountries', '/headers/north-america-countries-header.svg', 'countries', 'north-america-countries', 'countries-of-north-america', TRUE, FALSE, TRUE, TRUE),
(1, 7, 5, '', 'country', 'Countries of South America', 12, 300, 'SouthAmericaCountries', '/headers/south-america-countries-header.svg', 'countries', 'south-america-countries', 'countries-of-south-america', TRUE, FALSE, TRUE, TRUE),
(1, 8, 6, '', 'country', 'Countries of Oceania', 15, 300, 'OceaniaCountries', '/headers/oceania-countries-header.svg', 'countries', 'oceania-countries', 'countries-of-oceania', TRUE, FALSE, TRUE, TRUE),
(1, 7, 5, 'Argentina', 'province', 'Provinces of Argentina', 24, 300, 'ArgentinaProvinces', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e6-1f1f7.svg', 'provinces', 'argentina-provinces', 'provinces-of-argentina', TRUE, FALSE, TRUE, TRUE),
(1, 8, 6, 'Australia', 'state or territory', 'States and Territories of Australia', 9, 300, 'AustraliaStatesAndTerritories', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e6-1f1fa.svg', 'states and territories', 'australia-states-and-territories', 'states-and-territories-of-australia', TRUE, FALSE, TRUE, TRUE),
(1, 7, 5, 'Brazil', 'state', 'States of Brazil', 27, 300, 'BrazilStates', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e7-1f1f7.svg', 'states', 'brazil-states', 'states-of-brazil', TRUE, FALSE, TRUE, TRUE),
(1, 6, 4, 'Canada', 'province or territory', 'Provinces and Territories of Canada', 13, 300, 'CanadaProvincesAndTerritories', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e8-1f1e6.svg', 'provinces and territories', 'canada-provinces-and-territories', 'provinces-and-territories-of-canada', TRUE, FALSE, TRUE, TRUE),
(1, 4, 2, 'China', 'administrative division', 'Administrative Divisions of China', 33, 300, 'ChinaAdministrativeDivisions', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e8-1f1f3.svg', 'administrative divisions', 'china-administrative-divisions', 'administrative-divisions-of-china', TRUE, FALSE, FALSE, TRUE),
(1, 7, 5, 'Colombia', 'department', 'Departments of Colombia', 33, 300, 'ColombiaDepartments', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e8-1f1f4.svg', 'departments', 'colombia-departments', 'departments-of-colombia', TRUE, FALSE, TRUE, TRUE),
(1, 5, 3, 'France', 'region', 'Regions of France', 13, 300, 'FranceRegions', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1eb-1f1f7.svg', 'regions', 'france-regions', 'regions-of-france', TRUE, FALSE, TRUE, TRUE),
(1, 5, 3, 'Germany', 'state', 'States of Germany', 16, 300, 'GermanyStates', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e9-1f1ea.svg', 'states', 'germany-states', 'states-of-germany', TRUE, FALSE, TRUE, TRUE),
(1, 4, 2, 'India', 'state or union territory', 'States and Union Territories of India', 36, 300, 'IndiaStatesAndUnionTerritories', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ee-1f1f3.svg', 'states and union territories', 'india-states-and-union-territories', 'states-and-union-territories-of-india', TRUE, FALSE, FALSE, TRUE),
(1, 5, 3, 'Italy', 'region', 'Regions of Italy', 20, 300, 'ItalyRegions', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ee-1f1f9.svg', 'regions', 'italy-regions', 'regions-of-italy', TRUE, FALSE, TRUE, TRUE),
(1, 4, 2, 'Japan', 'prefecture', 'Prefectures of Japan', 47, 300, 'JapanPrefectures', 'https://upload.wikimedia.org/wikipedia/en/9/9e/Flag_of_Japan.svg', 'prefectures', 'japan-prefectures', 'prefectures-of-japan', TRUE, FALSE, TRUE, TRUE),
(1, 6, 4, 'Mexico', 'state', 'States of Mexico', 32, 300, 'MexicoStates', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f2-1f1fd.svg', 'states', 'mexico-states', 'states-of-mexico', TRUE, FALSE, FALSE, TRUE),
(1, 8, 6, 'New Zealand', 'region', 'Regions of New Zealand', 16, 300, 'NewZealandRegions', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f3-1f1ff.svg', 'regions', 'new-zealand-regions', 'regions-of-new-zealand', TRUE, FALSE, FALSE, TRUE),
(1, 3, 1, 'Nigeria', 'state', 'States of Nigeria', 37, 300, 'NigeriaStates', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f3-1f1ec.svg', 'states', 'nigeria-states', 'states-of-nigeria', TRUE, FALSE, FALSE, TRUE),
(1, 4, 2, 'Pakistan', 'administrative unit', 'Administrative Units of Pakistan', 8, 300, 'PakistanAdministrativeUnits', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f5-1f1f0.svg', 'administrative units', 'pakistan-administrative-units', 'administrative-units-of-pakistan', TRUE, FALSE, FALSE, TRUE),
(1, 4, 2, 'Russia', 'federal subject', 'Federal Subjects of Russia', 83, 600, 'RussiaFederalSubjects', 'https://upload.wikimedia.org/wikipedia/en/f/f3/Flag_of_Russia.svg', 'federal subjects', 'russia-federal-subjects', 'federal-subjects-of-russia', TRUE, FALSE, TRUE, TRUE),
(1, 2, null, '', 'mountain', 'Seven Summits', 7, 300, 'SevenSummits', '/headers/world-map-header.svg', 'mountains', 'seven-summits', 'seven-summits', FALSE, FALSE, TRUE, TRUE),
(1, 4, 2, 'South Korea', 'province', 'Provinces of South Korea', 17, 300, 'SouthKoreaProvinces', 'https://upload.wikimedia.org/wikipedia/commons/0/09/Flag_of_South_Korea.svg', 'provinces', 'south-korea-provinces', 'provinces-of-south-korea', TRUE, FALSE, TRUE, TRUE),
(1, 5, 3, 'Spain', 'province', 'Provinces of Spain', 52, 300, 'SpainProvinces', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ea-1f1f8.svg', 'provinces', 'spain-provinces', 'provinces-of-spain', TRUE, FALSE, TRUE, TRUE),
(1, 4, 2, 'Turkey', 'province', 'Provinces of Turkey', 81, 600, 'TurkeyProvinces', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1f9-1f1f7.svg', 'provinces', 'turkey-provinces', 'provinces-of-turkey', TRUE, FALSE, FALSE, TRUE),
(1, 3, 1, 'Uganda', 'district', 'Districts of Uganda', 112, 600, 'UgandaDistricts', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1fa-1f1ec.svg', 'districts', 'uganda-districts', 'districts-of-uganda', TRUE, FALSE, FALSE, TRUE),
(1, 5, 3, 'Ukraine', 'oblast', 'Oblasts of Ukraine', 25, 300, 'UkraineOblasts', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1fa-1f1e6.svg', 'oblasts', 'ukraine-oblasts', 'oblasts-of-ukraine', TRUE, FALSE, FALSE, TRUE),
(1, 5, 3, 'United Kingdom', 'country', 'Countries of the UK', 3, 300, 'UkCountries', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ec-1f1e7.svg', 'countries', 'uk-countries', 'countries-of-the-uk', FALSE, FALSE, TRUE, TRUE),
(1, 5, 3, 'United Kingdom', 'capital', 'Capitals of the UK', 3, 300, 'UkCapitals', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ec-1f1e7.svg', 'capitals', 'uk-capitals', 'capitals-of-the-uk', FALSE, FALSE, TRUE, TRUE),
(1, 5, 3, 'United Kingdom', 'sea', 'Seas of the UK', 4, 300, 'UkSeas', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ec-1f1e7.svg', 'seas', 'uk-seas', 'seas-of-the-uk', FALSE, FALSE, FALSE, TRUE),
(1, 5, 3, 'United Kingdom', 'county', 'Counties of the UK', 87, 900, 'UkCounties', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ec-1f1e7.svg', 'counties', 'uk-counties', 'counties-of-the-uk', TRUE, FALSE, FALSE, TRUE),
(1, 5, 3, 'United Kingdom', 'city', 'Major Cities of the UK', 106, 900, 'UkMajorCities', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ec-1f1e7.svg', 'cities', 'uk-major-cities', 'major-cities-of-the-uk', TRUE, FALSE, FALSE, TRUE),
(1, 6, 4, 'United States', 'state', 'US States', 50, 300, 'UsStates', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1fa-1f1f8.svg', 'states', 'us-states', 'us-states', TRUE, FALSE, TRUE, TRUE),
(1, 3, 1, 'Zambia', 'province', 'Provinces of Zambia', 10, 300, 'ZambiaProvinces', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ff-1f1f2.svg', 'provinces', 'zambia-provinces', 'provinces-of-zambia', TRUE, FALSE, FALSE, TRUE),
(2, 2, null, '', 'flag', 'Flags of the World', 197, 900, '', '/headers/world-map-header.svg', 'flags', 'world-countries', 'flags-of-the-world', TRUE, TRUE, TRUE, TRUE),
(2, 8, 6, 'Australia', 'flag', 'Flags of Australia', 8, 300, '', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e6-1f1fa.svg', 'flags', 'australia-states-and-territories-flags', 'flags-of-australia', TRUE, FALSE, TRUE, TRUE),
(2, 7, 5, 'Argentina', 'flag', 'Flags of Argentina', 24, 300, '', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e6-1f1f7.svg', 'flags', 'argentina-provinces', 'flags-of-argentina', TRUE, FALSE, TRUE, TRUE),
(2, 7, 5, 'Brazil', 'flag', 'Flags of Brazil', 27, 300, '', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e7-1f1f7.svg', 'flags', 'brazil-states', 'flags-of-brazil', TRUE, FALSE, TRUE, TRUE),
(2, 6, 4, 'Canada', 'flag', 'Flags of Canada', 13, 300, '', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e8-1f1e6.svg', 'flags', 'canada-provinces-and-territories', 'flags-of-canada', TRUE, FALSE, TRUE, TRUE),
(2, 7, 5, 'Colombia', 'flag' ,'Flags of Colombia', 33, 300, '', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e8-1f1f4.svg', 'flags', 'colombia-departments', 'flags-of-colombia', TRUE, FALSE, TRUE, TRUE),
(2, 5, 3, 'France', 'flag', 'Flags of France', 13, 300, '', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1eb-1f1f7.svg', 'flags', 'france-regions', 'flags-of-france', TRUE, FALSE, TRUE, TRUE),
(2, 5, 3, 'Germany', 'flag', 'Flags of Germany', 16, 300, '', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1e9-1f1ea.svg', 'flags', 'germany-states', 'flags-of-germany', TRUE, FALSE, TRUE, TRUE),
(2, 5, 3, 'Italy', 'flag', 'Flags of Italy', 20, 300, '', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ee-1f1f9.svg', 'flags', 'italy-regions', 'flags-of-italy', TRUE, FALSE, TRUE, TRUE),
(2, 4, 2, 'Japan', 'flag', 'Flags of Japan', 47, 300, '', 'https://upload.wikimedia.org/wikipedia/en/9/9e/Flag_of_Japan.svg', 'flags', 'japan-prefectures', 'flags-of-japan', TRUE, FALSE, TRUE, TRUE),
(2, 4, 2, 'Russia', 'flag', 'Flags of Russia', 83, 600, '', 'https://upload.wikimedia.org/wikipedia/en/f/f3/Flag_of_Russia.svg', 'flags', 'russia-federal-subjects', 'flags-of-russia', TRUE, FALSE, TRUE, TRUE),
(2, 4, 2, 'South Korea', 'flag', 'Flags of South Korea', 17, 300, '', 'https://upload.wikimedia.org/wikipedia/commons/0/09/Flag_of_South_Korea.svg', 'flags', 'south-korea-provinces', 'flags-of-south-korea', TRUE, FALSE, TRUE, TRUE),
(2, 5, 3, 'Spain', 'flag', 'Flags of Spain', 52, 300, '', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1ea-1f1f8.svg', 'flags', 'spain-provinces', 'flags-of-spain', TRUE, FALSE, TRUE, TRUE),
(2, 5, 3, 'Ukraine', 'flag', 'Flags of Ukraine', 25, 300, '', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1fa-1f1e6.svg', 'flags', 'ukraine-oblasts', 'flags-of-ukraine', TRUE, FALSE, TRUE, TRUE),
(2, 6, 4, 'United States', 'flag', 'Flags of the US', 50, 300, '', 'https://twemoji.maxcdn.com/v/13.0.1/svg/1f1fa-1f1f8.svg', 'flags', 'us-states', 'flags-of-the-us', TRUE, FALSE, TRUE, TRUE);

INSERT INTO triviaQuestionType (name) values
('Text'),
('Image'),
('Flag'),
('Map');

INSERT INTO merch (name, description, sizeGuideImageUrl, price, externalLink, route) values
('Tee', 'With summer just around the corner, we teamed up with the only carbon neutral clothing company in New Zealand (Koa Goods) to bring you guys the freshest eco-friendly tee to let the squad know you''re ready to drop those countries of the world at a moments notice. By copping one of these OG pieces you''re directly contributing to hosting costs and helping us keep this thing afloat. Kia Kaha!', '/merch/tee/size-guide.png', 49.99, null, 'tee'),
('Socks', 'With summer just around the corner, we teamed up with the only carbon neutral clothing company in New Zealand (Koa Goods) to bring you guys the freshest eco-friendly hoof covers to let the squad know you don''t mess around when it comes to capital cities. By copping one of these OG pieces you''re directly contributing to hosting costs and helping us keep this thing afloat. Kia Kaha!', null, 11.99, null, 'socks'),
('Poster Combo', 'Just moved flats and the bedroom walls are looking bare, boring and barren? The folks at GeoBuff HQ have got you covered. We''ve teamed up with the goodfella''s at The Big Picture to spruce up that decor and let the homies know that when the flags come out, you mean business. By copping one of these OG pieces you''re directly contributing to hosting costs and helping us keep this thing afloat. Kia Kaha!', null, 29.99, null, 'poster-combo'),
('Sticker Pack', 'Rear window on the wagon covered in dust and Raglan Roast have run out of stickers? The folks at GeoBuff HQ have got you covered. We''ve teamed up with the goodfella''s at The Big Picture to spice up that rear window and let the geezers in the slow lane know that you get your geoflex on. By copping one of these OG pieces you''re directly contributing to hosting costs and helping us keep this thing afloat. Kia Kaha!', null, 24.99, null, 'sticker-pack');

INSERT INTO merchSizes (merchId, size, quantity) values
(1, 'S', 0),
(1, 'M', 0),
(1, 'L', 0),
(1, 'XL', 0),
(1, '2XL', 0),
(2, 'One Size Fits All', 0),
(3, 'A2', 0),
(4, 'A4', 0);

INSERT INTO merchImages (merchId, imageUrl, isPrimary) values
(1, '/merch/tee/ash-front.jpg', TRUE),
(1, '/merch/tee/pearl-front.jpg', FALSE),
(1, '/merch/tee/dan-back.jpg', FALSE),
(1, '/merch/tee/koa-label.jpg', FALSE),
(2, '/merch/socks/holding.jpg', TRUE),
(2, '/merch/socks/full-leg.jpg', FALSE),
(2, '/merch/socks/close-up.jpg', FALSE),
(3, '/merch/posters/combo.jpg', TRUE),
(3, '/merch/posters/solo-green.jpg', FALSE),
(3, '/merch/posters/solo-blue.jpg', FALSE),
(3, '/merch/posters/mockup.png', FALSE),
(4, '/merch/stickers/pack-flat.jpg', TRUE),
(4, '/merch/stickers/pack-angle.jpg', FALSE),
(4, '/merch/stickers/laptop.jpg', FALSE),
(4, '/merch/stickers/van.jpg', FALSE);

INSERT INTO discounts (merchId, code, amount) values
(null, 'NOSHIP', 5.99);

INSERT INTO orderStatus (status) values
('Pending'),
('Payment Received'),
('Shipped');

INSERT INTO shippingOptions (name, description, price, imageUrl) values
('NZ-Wide Standard Shipping', 'Expect delivery in 5-7 days', 5.99, 'https://upload.wikimedia.org/wikipedia/commons/1/1f/NZ_Post_logo.png');

INSERT INTO communityQuizStatus (name) values
('Pending'),
('Approved');

INSERT INTO triviaQuestionCategory (name) values
('Landmarks'),
('Currency'),
('Mountains'),
('Rivers'),
('Borders'),
('Community Quiz'),
('Misc');
