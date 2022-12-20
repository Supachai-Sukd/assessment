SET TIMEZONE="Asia/Bangkok";

CREATE TABLE IF NOT EXISTS expenses (
    id SERIAL PRIMARY KEY,
    title TEXT,
    amount FLOAT,
    note TEXT,
    tags TEXT[]
);

INSERT INTO "expenses" ("id", "title", "amount", "note", "tags") VALUES (1, 'test-title', 20, 'test-note test-note test-note', 'test-tags');