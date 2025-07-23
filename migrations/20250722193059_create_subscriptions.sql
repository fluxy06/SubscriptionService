DROP TABLE IF EXISTS subscriptions;

CREATE TABLE subscriptions (
    id SERIAL PRIMARY KEY,
    service_name VARCHAR(100) NOT NULL,
    price INT NOT NULL CHECK (price > 0),
    user_id UUID NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);
