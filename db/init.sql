-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS expenses_id_seq;

-- Table Definition
CREATE TABLE IF NOT EXISTS expenses (
        id SERIAL PRIMARY KEY,
        title TEXT,
        amount FLOAT,
        note TEXT,
        tags TEXT[]
);

INSERT INTO "expenses" ("id", "title", "amount", "note", "tags") VALUES (1, 'test-title', 99, 'test-note', '{test_tag1, test_tag2}');