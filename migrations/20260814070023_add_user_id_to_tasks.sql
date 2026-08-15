-- +goose Up
ALTER TABLE tasks ADD COLUMN user_id BIGINT;
UPDATE tasks SET user_id = 1 WHERE user_id IS NULL;
ALTER TABLE tasks ALTER COLUMN user_id SET NOT NULL;
CREATE INDEX tasks_user_id_idx ON tasks (user_id);

-- +goose Down
DROP INDEX tasks_user_id_idx;
ALTER TABLE tasks DROP COLUMN user_id;

