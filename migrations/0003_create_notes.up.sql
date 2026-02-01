CREATE TABLE IF NOT EXISTS Notes (
    id UUID PRIMARY KEY,
    daily_entry_id UUID,
    note TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (daily_entry_id) REFERENCES DailyEntries(id)
);