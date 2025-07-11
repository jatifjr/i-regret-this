-- Create function to generate plot_id
CREATE OR REPLACE FUNCTION generate_plot_id()
RETURNS TRIGGER AS $$
DECLARE
    date_str TEXT;
    order_num INTEGER;
    new_plot_id BIGINT;
BEGIN
    -- Get the date part (YYYYMMDD)
    date_str := TO_CHAR(NEW.date_time, 'YYYYMMDD');
    
    -- Get the count of schedules for this date
    SELECT COALESCE(COUNT(*), 0) + 1
    INTO order_num
    FROM schedules
    WHERE DATE(date_time) = DATE(NEW.date_time);
    
    -- Combine date and order number
    new_plot_id := (date_str || LPAD(order_num::TEXT, 2, '0'))::BIGINT;
    
    -- Set the plot_id
    NEW.plot_id := new_plot_id;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create function to handle quota updates
CREATE OR REPLACE FUNCTION handle_quota_update()
RETURNS TRIGGER AS $$
DECLARE
    used_slots INTEGER;
BEGIN
    -- Calculate used slots (quota - available)
    used_slots := OLD.quota - OLD.available;
    
    -- If new quota is less than used slots, prevent the update
    IF NEW.quota < used_slots THEN
        RAISE EXCEPTION 'Cannot reduce quota below used slots (currently % slots in use)', used_slots;
    END IF;
    
    -- Update available slots based on new quota
    NEW.available := NEW.quota - used_slots;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create the schedules table
CREATE TABLE IF NOT EXISTS schedules (
    id BIGSERIAL PRIMARY KEY,
    plot_id BIGINT,
    date_time TIMESTAMP WITH TIME ZONE NOT NULL,
    location VARCHAR(255) NOT NULL,
    quota INTEGER NOT NULL CHECK (quota > 0),
    available INTEGER NOT NULL DEFAULT 0 CHECK (available >= 0),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create trigger to auto-generate plot_id
CREATE TRIGGER set_plot_id
    BEFORE INSERT ON schedules
    FOR EACH ROW
    EXECUTE FUNCTION generate_plot_id();

-- Create trigger to set available equal to quota on insert
CREATE OR REPLACE FUNCTION set_available_to_quota()
RETURNS TRIGGER AS $$
BEGIN
    NEW.available := NEW.quota;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_available
    BEFORE INSERT ON schedules
    FOR EACH ROW
    EXECUTE FUNCTION set_available_to_quota();

-- Create trigger to handle quota updates
CREATE TRIGGER handle_quota_change
    BEFORE UPDATE ON schedules
    FOR EACH ROW
    WHEN (OLD.quota IS DISTINCT FROM NEW.quota)
    EXECUTE FUNCTION handle_quota_update();