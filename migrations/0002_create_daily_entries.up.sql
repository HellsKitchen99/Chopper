CREATE TABLE IF NOT EXISTS DailyEntries (
    id UUID PRIMARY KEY,
    user_id UUID,
    date DATE NOT NULL, 
    UNIQUE(user_id, date),
    mood SMALLINT NOT NULL,
    sleep_hours NUMERIC(2, 1) NOT NULL,
    load SMALLINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES Users(id)
);