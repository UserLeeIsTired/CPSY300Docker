DROP TABLE my_student CASCADE;

CREATE TABLE my_student (
    id VARCHAR(9) PRIMARY KEY,
    name VARCHAR(255),
    course VARCHAR(100),
    last_update TIMESTAMP
);

CREATE OR REPLACE FUNCTION insert_date_function()
RETURNS TRIGGER AS $$
BEGIN
    NEW.last_update := NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER insert_date_trigger
BEFORE INSERT ON my_student
FOR EACH ROW
EXECUTE FUNCTION insert_date_function();

CREATE OR REPLACE FUNCTION update_student_function()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.id = '' THEN
        NEW.id = OLD.id;
    END IF;

    IF NEW.name = '' THEN
        NEW.name = OLD.name;
    END IF;

    IF NEW.course = '' THEN
        NEW.course = OLD.course;
    END IF;

    NEW.last_update := NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER update_student_trigger
BEFORE UPDATE ON my_student
FOR EACH ROW
EXECUTE FUNCTION update_student_function();