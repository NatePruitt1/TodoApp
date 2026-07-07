CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY,
    project_id UUID REFERENCES projects(id),
    category_name varchar(255) NOT NULL,
    index INT NOT NULL
);