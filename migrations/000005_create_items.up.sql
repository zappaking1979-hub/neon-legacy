CREATE TABLE IF NOT EXISTS items (
    id          SERIAL PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    type        TEXT NOT NULL DEFAULT 'consumable',
    buy_price   BIGINT NOT NULL DEFAULT 0,
    sell_price  BIGINT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS player_items (
    id          SERIAL PRIMARY KEY,
    player_id   UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    item_id     INTEGER NOT NULL REFERENCES items(id),
    quantity    INTEGER NOT NULL DEFAULT 1,
    UNIQUE(player_id, item_id)
);

CREATE INDEX IF NOT EXISTS idx_player_items_player ON player_items(player_id);

INSERT INTO items (name, description, type, buy_price, sell_price) VALUES
    ('First Aid Kit',    'Restores 25 HP instantly.',                  'consumable', 500,   250),
    ('Energy Drink',     'Restores 20 Energy instantly.',              'consumable', 300,   150),
    ('Nerve Pills',      'Restores 10 Nerve instantly.',               'consumable', 400,   200),
    ('Coffee',           'Restores 15 Awake instantly.',               'consumable', 200,   100),
    ('Brass Knuckles',   'A simple weapon. +5 Strength in combat.',    'weapon',     2500,  1250),
    ('Kevlar Vest',      'Light body armor. +5 Defense in combat.',    'armor',      3000,  1500),
    ('Exp Booster',      'Double EXP for 30 minutes.',                 'special',    10000, 5000),
    ('Lucky Charm',      '+10% crime success rate for 30 minutes.',    'special',    8000,  4000);
