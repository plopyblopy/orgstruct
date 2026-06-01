-- +goose Up


-- chk_department_cycle проверяет id и parent_id новой записи - они не должны совпадать.
-- Проверяет наличие цикла в цепочке departments.
-- В случае ошибки одинаковых значений вернет: error_msg: self_reference, SQLSTATE: P0002.
-- В случае ошибки цикла вернет: error_msg: cycle, SQLSTATE: P0003.

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION chk_department_cycle()
RETURNS trigger AS $$
DECLARE
    current_id INT;
BEGIN
    IF NEW.parent_id = NEW.id THEN
        RAISE EXCEPTION SQLSTATE 'P0002'
            USING MESSAGE = 'A department cannot be a parent of itself';
    END IF;

    current_id := NEW.parent_id;
    WHILE current_id IS NOT NULL LOOP
        IF current_id = NEW.id THEN
            RAISE EXCEPTION SQLSTATE 'P0003' 
                USING MESSAGE = 'A cycle has been detected: department ' || NEW.id ||' is already the parent';
        END IF;
        SELECT parent_id INTO current_id
        FROM departments
        WHERE id = current_id;
    END LOOP;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- trg_departments_no_cycle триггер для проверки цикла.
CREATE TRIGGER trg_departments_no_cycle
BEFORE INSERT OR UPDATE ON departments
FOR EACH ROW
WHEN (NEW.parent_id IS NOT NULL ) -- Если не корень.
EXECUTE FUNCTION chk_department_cycle();

-- +goose Down
DROP TRIGGER IF EXISTS trg_departments_no_cycle ON departments;
DROP FUNCTION IF EXISTS chk_department_cycle;
