DROP TABLE IF EXISTS weather;

DROP TABLE IF EXISTS cities;

CREATE TABLE cities (
    id SERIAL NOT NULL,
    name TEXT,
    country text,
    lat numeric,
    lon numeric,
    PRIMARY KEY(id)
);

CREATE TABLE weather (
    id SERIAL NOT NULL,
    temp real,
    f_date DATE,
    full_info JSONB,
    city_id INTEGER,
    PRIMARY KEY(id),
    CONSTRAINT fk_city
        FOREIGN KEY(city_id)
            REFERENCES cities(id)
);

INSERT INTO cities(name, country, lat, lon) VALUES ('Chelyabinsk', 'RU', 55.159841, 61.402555), 
('Izhevsk', 'RU', 56.866557, 53.209417), 
('Kazan', 'RU', 55.782355, 49.124227), 
('Krasnodar', 'RU', 45.035272, 38.976481), 
('Krasnoyarsk', 'RU', 56.009097, 92.872515), 
('Moscow', 'RU', 55.750446, 37.617494), 
('Nizhny Novgorod', 'RU', 56.326482, 44.005139), 
('Novosibirsk', 'RU', 55.028217, 82.923451), 
('Omsk', 'RU', 54.991375, 73.371529), 
('Perm', 'RU', 58.02148705, 56.23076652679421), 
('Rostov-on-Don', 'RU', 47.221386, 39.711420), 
('Saint Petersburg', 'RU', 59.938732, 30.316229), 
('Samara', 'RU', 53.198627, 50.113987), 
('Saratov', 'RU', 51.530018, 46.034683), 
('Tolyatti', 'RU', 53.514950, 49.407574), 
('Tyumen', 'RU', 57.153534, 65.542274), 
('Ufa', 'RU', 54.726141, 55.947499), 
('Volgograd', 'RU', 48.708191, 44.515335), 
('Voronezh', 'RU', 51.660598, 39.200586), 
('Yekaterinburg', 'RU', 56.839104, 60.608250);