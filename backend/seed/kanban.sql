DROP TABLE IF EXISTS cards;
DROP TABLE IF EXISTS columns;
DROP TABLE IF EXISTS boards;

CREATE TABLE boards (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE columns (
    id SERIAL PRIMARY KEY,
    board_id INTEGER REFERENCES boards(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    position INTEGER NOT NULL
);

CREATE TABLE cards (
    id SERIAL PRIMARY KEY,
    column_id INTEGER REFERENCES columns(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    position INTEGER NOT NULL
);

INSERT INTO boards (name) VALUES ('Project Alpha');

INSERT INTO columns (board_id, name, position) VALUES (1, 'To Do', 1);
INSERT INTO columns (board_id, name, position) VALUES (1, 'In Progress', 2);
INSERT INTO columns (board_id, name, position) VALUES (1, 'Done', 3);

INSERT INTO cards (column_id, title, description, position) VALUES (1, 'Set up project', 'Initialize repository and dependencies', 1);
INSERT INTO cards (column_id, title, description, position) VALUES (2, 'Develop feature X', 'Implement the main feature', 1);
INSERT INTO cards (column_id, title, description, position) VALUES (3, 'Deploy to production', 'Release the app', 1);