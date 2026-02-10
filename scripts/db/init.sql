CREATE TABLE IF NOT EXISTS "Book" (
    "bookId" VARCHAR(255) PRIMARY KEY,
    "title" VARCHAR(255) NOT NULL,
    "priceAmount" NUMERIC(10, 2) NOT NULL
);

CREATE TABLE IF NOT EXISTS "Stock" (
    "stockId" VARCHAR(255) PRIMARY KEY,
    "bookId" VARCHAR(255) NOT NULL,
    "quantityAvailable" INTEGER NOT NULL,
    "status" VARCHAR(50) NOT NULL,
    CONSTRAINT fk_book
        FOREIGN KEY ("bookId")
        REFERENCES "Book" ("bookId")
        ON DELETE CASCADE
);
