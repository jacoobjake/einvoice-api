-- Create trigger function for updated_at
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger function for soft delete
CREATE OR REPLACE FUNCTION soft_delete()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE ONLY users
    SET deleted_at = NOW()
    WHERE id = OLD.id;

    RETURN NULL; -- stop the actual DELETE
END;
$$ LANGUAGE plpgsql;