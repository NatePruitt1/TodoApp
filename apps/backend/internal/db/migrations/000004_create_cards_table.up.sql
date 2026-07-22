CREATE TABLE IF NOT EXISTS cards (
    id UUID PRIMARY KEY,
    title varchar(255) NOT NULL,
    content varchar,
    category_id UUID REFERENCES categories(id) ON DELETE CASCADE,
    finished BOOL NOT NULL
);