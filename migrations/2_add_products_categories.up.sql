CREATE TABLE IF NOT EXISTS categories (
    id VARCHAR(36) PRIMARY KEY,
    title text NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS products (
    id VARCHAR(36) PRIMARY KEY,
    title text NOT NULL,
    price float NOT NULL,
    imageIds jsonb,
    maker TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS product_category (
    id VARCHAR(36) PRIMARY KEY,
    productId VARCHAR(36) NOT NULL,
    categoryId VARCHAR(36) NOT NULL,
    FOREIGN KEY (productId) REFERENCES products (id) ON DELETE CASCADE,
    FOREIGN KEY (categoryId) REFERENCES categories (id) ON DELETE CASCADE

);
