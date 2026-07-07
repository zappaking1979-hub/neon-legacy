CREATE TABLE crimes (
    id          SERIAL PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    nerve_cost  INT NOT NULL DEFAULT 1,
    min_level   INT NOT NULL DEFAULT 1,
    min_strength INT NOT NULL DEFAULT 0,
    min_defense INT NOT NULL DEFAULT 0,
    min_speed   INT NOT NULL DEFAULT 0,
    success_rate NUMERIC(5,2) NOT NULL DEFAULT 90.00,
    jail_chance NUMERIC(5,2) NOT NULL DEFAULT 2.00,
    exp_reward  INT NOT NULL DEFAULT 10,
    cash_reward_min INT NOT NULL DEFAULT 50,
    cash_reward_max INT NOT NULL DEFAULT 150,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO crimes (name, description, nerve_cost, min_level, exp_reward, cash_reward_min, cash_reward_max, success_rate, jail_chance) VALUES
    ('Mug Homeless Man', 'Target the weakest in society. Low risk, low reward.', 1, 1, 5, 10, 50, 95.00, 1.00),
    ('Shoplift', 'Finger some goods from the corner store.', 2, 1, 10, 25, 100, 92.00, 2.00),
    ('Pickpocket', 'Work the crowded streets. Stealth required.', 3, 2, 20, 50, 200, 88.00, 3.00),
    ('Break Into Car', 'Smash and grab. Quick cash if you are fast.', 4, 3, 35, 100, 400, 82.00, 4.00),
    ('Mug Businessman', 'Rich target, better payout. Higher risk.', 5, 4, 50, 200, 800, 78.00, 5.00),
    ('Bank Heist', 'Hit the bank vault. Big score, big risk.', 8, 6, 100, 500, 2000, 65.00, 8.00),
    ('Armored Truck Robbery', 'Military-grade heist. You will need friends in low places.', 10, 8, 150, 1000, 5000, 55.00, 12.00);
