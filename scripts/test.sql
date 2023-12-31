

CREATE TABLE foo(
 `id` int,
 `name` VARCHAR(233),
  created bigint(13) NOT NULL DEFAULT (FLOOR(UNIX_TIMESTAMP(CURRENT_TIMESTAMP(3)) * 1000)),
  updated bigint(13) NOT NULL DEFAULT (FLOOR(UNIX_TIMESTAMP(CURRENT_TIMESTAMP(3)) * 1000))
);

INSERT into foo (id, name) VALUES (2, "lixin");
UPDATE foo SET name="lixin"  WHERE id = 2;

CREATE TRIGGER user_updated_at_column
BEFORE UPDATE ON foo
FOR EACH ROW
    SET NEW.updated = FLOOR(UNIX_TIMESTAMP(CURRENT_TIMESTAMP(3)) * 1000);

CREATE TRIGGER foo_updated_at_column
AFTER CREATE ON foo
FOR EACH ROW
    SET NEW.created = UNIX_TIMESTAMP(CURRENT_TIMESTAMP(3)) * 1000;