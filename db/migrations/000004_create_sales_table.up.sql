CREATE TABLE IF NOT EXISTS sales (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    date TIMESTAMP NOT NULL,
    total_amount NUMERIC(10, 2) NOT NULL,
    discount NUMERIC(10, 2),
    additional_charges NUMERIC(10, 2)
);
