DROP TABLE IF EXISTS tasks;
CREATE TABLE tasks(
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    completed BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO tasks(title, description, completed) VALUES ('Learn Go', 'Basic course', 'true'),
('Write api', 'write yourself', 'false'),
('Deploy app', 'by yourself', 'false');
