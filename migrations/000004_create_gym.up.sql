CREATE TABLE IF NOT EXISTS gym_exercises (
    id          SERIAL PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    stat        TEXT NOT NULL,
    energy_cost INTEGER NOT NULL DEFAULT 2,
    min_level   INTEGER NOT NULL DEFAULT 1,
    gain_min    INTEGER NOT NULL DEFAULT 1,
    gain_max    INTEGER NOT NULL DEFAULT 2
);

INSERT INTO gym_exercises (name, description, stat, energy_cost, min_level, gain_min, gain_max) VALUES
    ('Push-ups',        'Build upper body strength with classic push-ups.',              'strength', 2, 1, 1, 2),
    ('Sparring',        'Practice combat techniques with a training partner.',           'defense',  2, 1, 1, 2),
    ('Sprinting',       'Short bursts of maximum speed running.',                        'speed',    2, 1, 1, 2),
    ('Yoga',            'Flexibility and balance exercises.',                            'agility',  2, 1, 1, 2),
    ('Heavy Bag',       'Pound the heavy bag to build raw power.',                       'strength', 4, 5, 2, 4),
    ('Tactical Drill',  'Run through obstacle course to improve reflexes.',              'defense',  4, 5, 2, 4),
    ('Parkour',         'Navigate urban obstacles with fluid movement.',                 'agility',  4, 5, 2, 4),
    ('Reflex Training', 'React to randomized stimuli for faster response times.',        'speed',    4, 5, 2, 4);
