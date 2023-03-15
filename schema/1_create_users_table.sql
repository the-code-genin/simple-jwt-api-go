-- up
CREATE TABLE IF NOT EXISTS users (
    id SERIAL NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,

    PRIMARY KEY (id)
);

-- down
DROP TABLE IF EXISTS users;