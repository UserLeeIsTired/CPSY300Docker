CREATE OR REPLACE FUNCTION insert_date_function()
RETURNS TRIGGER AS $$
BEGIN
    NEW.created_date := NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER insert_date_trigger
BEFORE INSERT ON my_student
FOR EACH ROW
EXECUTE FUNCTION insert_date_function();