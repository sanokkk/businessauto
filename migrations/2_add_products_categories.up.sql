CREATE TABLE IF NOT EXISTS categories (
    id pg_catalog.uuid PRIMARY KEY,
    title text NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS products (
    id pg_catalog.uuid PRIMARY KEY,
    title text NOT NULL,
    price float NOT NULL,
    imageIds jsonb,
    maker TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS product_category (
    id pg_catalog.uuid PRIMARY KEY,
    productId pg_catalog.uuid NOT NULL,
    categoryId pg_catalog.uuid NOT NULL,
    FOREIGN KEY (productId) REFERENCES products (id) ON DELETE CASCADE,
    FOREIGN KEY (categoryId) REFERENCES categories (id) ON DELETE CASCADE

);
