CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    total_price DECIMAL(10,2) NOT NULL,
    status VARCHAR(20) CHECK (status IN ('pending', 'paid', 'shipped', 'delivered', 'cancelled')),
    created_at TIMESTAMP DEFAULT NOW()
);