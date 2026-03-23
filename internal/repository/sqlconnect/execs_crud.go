package sqlconnect

import (
	"database/sql"
	"errors"

	"restapi/internal/models"
	"restapi/pkg/utils"
)

var (
	ErrExecNotFound      = errors.New("exec not found")
	ErrExecPasswordEmpty = errors.New("exec password must not be empty")
)

func ListExecs() ([]models.Exec, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "unable to connect listing execs to database")
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, first_name, last_name, email, username, role FROM execs")
	if err != nil {
		return nil, utils.ErrorHandler(err, "unable to query execs")
	}
	defer rows.Close()

	result := make([]models.Exec, 0)
	for rows.Next() {
		var e models.Exec
		if err := rows.Scan(&e.ID, &e.FirstName, &e.LastName, &e.Email, &e.Username, &e.Role); err != nil {
			return nil, utils.ErrorHandler(err, "unable to scan exec row")
		}
		result = append(result, e)
	}

	if err := rows.Err(); err != nil {
		return nil, utils.ErrorHandler(err, "error iterating exec rows")
	}
	return result, nil
}

func GetExecByID(id int) (*models.Exec, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "unable to connect retrieving exec to database")
	}
	defer db.Close()

	var exec models.Exec
	err = db.QueryRow(
		"SELECT id, first_name, last_name, email, username, password, role FROM execs WHERE id = ?",
		id,
	).
		Scan(&exec.ID, &exec.FirstName, &exec.LastName, &exec.Email, &exec.Username, &exec.Password, &exec.Role)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrExecNotFound
	}
	if err != nil {
		return nil, utils.ErrorHandler(err, "unable to scan exec")
	}
	return &exec, nil
}

func InsertExecs(execs []models.Exec) ([]models.Exec, error) {
	if len(execs) == 0 {
		return []models.Exec{}, nil
	}

	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "unable to connect inserting execs to database")
	}
	defer db.Close()

	query := "INSERT INTO execs (first_name, last_name, email, username, password, role) VALUES (?, ?, ?, ?, ?, ?)"
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, utils.ErrorHandler(err, "unable to prepare exec insert statement")
	}
	defer stmt.Close()

	inserted := make([]models.Exec, 0, len(execs))
	for _, exec := range execs {
		hashedPassword, err := utils.HashPassword(exec.Password)
		if err != nil {
			if errors.Is(err, utils.ErrPasswordEmpty) {
				return nil, ErrExecPasswordEmpty
			}
			return nil, err
		}
		exec.Password = hashedPassword

		result, err := stmt.Exec(exec.FirstName, exec.LastName, exec.Email, exec.Username, exec.Password, exec.Role)
		if err != nil {
			return nil, utils.ErrorHandler(err, "unable to insert exec")
		}
		lastID, err := result.LastInsertId()
		if err != nil {
			return nil, utils.ErrorHandler(err, "unable to read inserted exec id")
		}
		exec.ID = int(lastID)
		inserted = append(inserted, exec)
	}

	return inserted, nil
}

func UpdateExec(exec models.Exec) error {
	db, err := ConnectDB()
	if err != nil {
		return utils.ErrorHandler(err, "unable to connect updating exec to database")
	}
	defer db.Close()

	passwordValue := exec.Password
	if passwordValue != "" {
		passwordValue, err = utils.HashPassword(passwordValue)
		if err != nil {
			if errors.Is(err, utils.ErrPasswordEmpty) {
				return ErrExecPasswordEmpty
			}
			return err
		}
	} else {
		passwordValue, err = fetchExecPassword(db, exec.ID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ErrExecNotFound
			}
			return utils.ErrorHandler(err, "unable to load existing exec password")
		}
	}

	result, err := db.Exec(
		"UPDATE execs SET first_name = ?, last_name = ?, email = ?, username = ?, password = ?, role = ? WHERE id = ?",
		exec.FirstName, exec.LastName, exec.Email, exec.Username, passwordValue, exec.Role, exec.ID,
	)
	if err != nil {
		return utils.ErrorHandler(err, "unable to update exec")
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return utils.ErrorHandler(err, "unable to check exec update result")
	}
	if rows == 0 {
		return ErrExecNotFound
	}
	return nil
}

func DeleteExec(id int) error {
	db, err := ConnectDB()
	if err != nil {
		return utils.ErrorHandler(err, "unable to connect deleting exec to database")
	}
	defer db.Close()

	result, err := db.Exec("DELETE FROM execs WHERE id = ?", id)
	if err != nil {
		return utils.ErrorHandler(err, "unable to delete exec")
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return utils.ErrorHandler(err, "unable to check exec delete result")
	}
	if rows == 0 {
		return ErrExecNotFound
	}
	return nil
}

type ExecStore struct{}

func fetchExecPassword(db *sql.DB, id int) (string, error) {
	var password string
	err := db.QueryRow("SELECT password FROM execs WHERE id = ?", id).Scan(&password)
	if err != nil {
		return "", err
	}
	return password, nil
}

func GetExecByUsername(username string) (*models.Exec, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "unable to connect retrieving exec by username to database")
	}
	defer db.Close()

	var exec models.Exec
	err = db.QueryRow(
		"SELECT id, first_name, last_name, email, username, password, inactive_status, role FROM execs WHERE username = ?",
		username,
	).
		Scan(&exec.ID, &exec.FirstName, &exec.LastName, &exec.Email, &exec.Username, &exec.Password, &exec.InactiveStatus, &exec.Role)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrExecNotFound
	}
	if err != nil {
		return nil, utils.ErrorHandler(err, "unable to scan exec by username")
	}
	return &exec, nil
}
