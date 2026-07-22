CREATE TABLE IF NOT EXISTS projects (
    id UUID PRIMARY KEY,
    project_name varchar(255) NOT NULL,
    description varchar,
    owner_id UUID REFERENCES users(id) ON DELETE CASCADE
);