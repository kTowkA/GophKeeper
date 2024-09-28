BEGIN;
CREATE TABLE IF NOT EXISTS keep_folder (
    folder_id uuid,
    user_id uuid,
    title text,
    description text,
    create_at timestamp,
    update_at timestamp,
    PRIMARY KEY(folder_id),
    UNIQUE(user_id,title)
);
CREATE TABLE IF NOT EXISTS keep_value (
    value_id uuid,
    folder_id uuid,
    title text,
    description text,
    value bytea,
    create_at timestamp,
    PRIMARY KEY(value_id),
    UNIQUE(folder_id,title)
);
COMMIT;