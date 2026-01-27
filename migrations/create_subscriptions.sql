CREATE TABLE IF NOT EXISTS subscriptions (
    id SERIAL PRIMARY KEY DEFAULT,
    service_name VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL CHECK (price>0),
    user_id UUID NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE 
);

CREATE INDEX IF NOT EXISTS idx_sub_user_id ON subscriptions(user_id);
CREATE INDEX IF NOT EXISTS idx_sub_service_name ON subscriptions(service_name);
CREATE INDEX IF NOT EXISTS idx_sub_start_date ON subscriptions(start_date);