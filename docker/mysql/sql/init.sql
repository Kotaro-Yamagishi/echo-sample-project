CREATE TABLE users (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(255) NOT NULL
);

INSERT INTO users (id, name) VALUES
  (1, 'Alice'),
  (2, 'Bob'),
  (3, 'Carol'),
  (4, 'Dave'),
  (5, 'Eve');