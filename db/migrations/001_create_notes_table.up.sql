CREATE TABLE IF NOT EXISTS notes
(
    s_no serial PRIMARY KEY,
    user_id VARCHAR (50) NOT NULL,
    id VARCHAR ( 50 ) NOT NULL,
    title TEXT,
    body TEXT,
    created_at TIMESTAMP NOT NULL
);
