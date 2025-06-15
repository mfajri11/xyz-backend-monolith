package repository

const (
	queryUpdateUserById = `UPDATE xyz.user
SET 
    national_id = COALESCE(?, national_id), 
    full_name = COALESCE(?, full_name), 
    legal_name = COALESCE(?, legal_name), 
    birth_of_place = COALESCE(?, birth_of_place), 
    birth_of_date = COALESCE(?, birth_of_date), 
    salary = COALESCE(?, salary), 
    nation_id_photo = COALESCE(?, nation_id_photo), 
    user_photo = COALESCE(?, user_photo), 
    is_nid_valid = COALESCE(?, is_nid_valid), 
    is_photo_valid = COALESCE(?, is_photo_valid), 
    created_by = COALESCE(?, created_by), 
    updated_by = COALESCE(?, updated_by)
WHERE id = ?`

	getUserByNationalID = `SELECT id, national_id, full_name, legal_name, is_nid_valid, is_photo_valid
FROM user
WHERE national_id = ?
`
)
