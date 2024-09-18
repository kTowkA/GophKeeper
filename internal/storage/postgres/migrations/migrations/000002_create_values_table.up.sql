BEGIN;
CREATE TABLE IF NOT EXISTS keep_element (
    element_id uuid,
    user_id uuid,
    title text,
    description text,
    adding_at timestamp,
    PRIMARY KEY(element_id),
    UNIQUE(user_id,title)
);
CREATE TABLE IF NOT EXISTS keep_value (
    value_id uuid,
    element_id uuid,
    title text,
    description text,
    value bytea,
    adding_at timestamp,
    PRIMARY KEY(value_id),
    UNIQUE(element_id,title)
);
COMMIT;