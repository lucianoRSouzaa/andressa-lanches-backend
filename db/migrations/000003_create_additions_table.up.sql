CREATE TABLE IF NOT EXISTS additions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    price NUMERIC(10, 2) NOT NULL
);
