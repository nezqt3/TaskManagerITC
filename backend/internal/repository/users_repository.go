package repository

import "backend/internal/model"

func GetUsers(cfg *model.Config) ([]model.UserProfile, error) {
	db, err := openDB(cfg)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(`
		SELECT TelegramID, FirstName, LastName, Username, PhotoURL, FullName, DateOfBirthday, NumberOfPhone, Role, MayToOpen
		FROM users
		ORDER BY FullName
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]model.UserProfile, 0)
	for rows.Next() {
		var u model.UserProfile
		if err := rows.Scan(
			&u.TelegramID,
			&u.FirstName,
			&u.LastName,
			&u.Username,
			&u.PhotoURL,
			&u.FullName,
			&u.DateOfBirthday,
			&u.NumberOfPhone,
			&u.Role,
			&u.MayToOpen,
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func GetUserByTelegramID(cfg *model.Config, telegramID string) (*model.UserProfile, error) {
	db, err := openDB(cfg)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	u := &model.UserProfile{}
	err = db.QueryRow(`
		SELECT TelegramID, FirstName, LastName, Username, PhotoURL, FullName, DateOfBirthday, NumberOfPhone, Role, MayToOpen
		FROM users
		WHERE TelegramID = ?
	`, telegramID).Scan(
		&u.TelegramID,
		&u.FirstName,
		&u.LastName,
		&u.Username,
		&u.PhotoURL,
		&u.FullName,
		&u.DateOfBirthday,
		&u.NumberOfPhone,
		&u.Role,
		&u.MayToOpen,
	)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func GetUserByUsername(cfg *model.Config, normalized string, normalizedWithAt string) (*model.UserProfile, error) {
	db, err := openDB(cfg)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	u := &model.UserProfile{}
	err = db.QueryRow(`
		SELECT TelegramID, FirstName, LastName, Username, PhotoURL, FullName, DateOfBirthday, NumberOfPhone, Role, MayToOpen
		FROM users
		WHERE lower(trim(Username)) = ? OR lower(trim(Username)) = ?
	`, normalized, normalizedWithAt).Scan(
		&u.TelegramID,
		&u.FirstName,
		&u.LastName,
		&u.Username,
		&u.PhotoURL,
		&u.FullName,
		&u.DateOfBirthday,
		&u.NumberOfPhone,
		&u.Role,
		&u.MayToOpen,
	)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func GetUserByFullName(cfg *model.Config, normalized string) (*model.UserProfile, error) {
	db, err := openDB(cfg)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	u := &model.UserProfile{}
	err = db.QueryRow(`
		SELECT TelegramID, FirstName, LastName, Username, PhotoURL, FullName, DateOfBirthday, NumberOfPhone, Role, MayToOpen
		FROM users
		WHERE lower(trim(FullName)) = ? OR lower(trim(FullName)) = ?
	`, normalized).Scan(
		&u.TelegramID,
		&u.FirstName,
		&u.LastName,
		&u.Username,
		&u.PhotoURL,
		&u.FullName,
		&u.DateOfBirthday,
		&u.NumberOfPhone,
		&u.Role,
		&u.MayToOpen,
	)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func SearchUsersByFullName(cfg *model.Config, fullName string) ([]model.UserProfile, error) {
	db, err := openDB(cfg)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	users := []model.UserProfile{}

	rows, err := db.Query(`
		SELECT 
			TelegramID,
			FirstName,
			LastName,
			Username,
			PhotoURL,
			FullName,
			DateOfBirthday,
			NumberOfPhone,
			Role,
			MayToOpen
		FROM users
		WHERE trim(FullName) LIKE ?
		ORDER BY FullName
	`, "%"+fullName+"%")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var user model.UserProfile

		err := rows.Scan(
			&user.TelegramID,
			&user.FirstName,
			&user.LastName,
			&user.Username,
			&user.PhotoURL,
			&user.FullName,
			&user.DateOfBirthday,
			&user.NumberOfPhone,
			&user.Role,
			&user.MayToOpen,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func UpdateUser(cfg *model.Config, telegramID string, updates *model.UserProfile) error {
	db, err := openDB(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`
		UPDATE users
		SET FullName = ?, Username = ?, DateOfBirthday = ?, NumberOfPhone = ?, Role = ?, MayToOpen = ?
		WHERE TelegramID = ?
	`,
		updates.FullName,
		updates.Username,
		updates.DateOfBirthday,
		updates.NumberOfPhone,
		updates.Role,
		updates.MayToOpen,
		telegramID,
	)
	return err
}
