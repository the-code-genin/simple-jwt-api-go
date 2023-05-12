CREATE TABLE IF NOT EXISTS service.users (
    id UUID NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
);

CREATE INDEX IF NOT EXISTS users_email_lower_index ON service.users ((LOWER(email)));