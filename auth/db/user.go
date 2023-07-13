package db

// User 表示用户
type User struct {
	ID string `gorm:"type:varchar(40);primayKey"`
	// 账号
	Account string `gorm:"type:varchar(40);uniqueIndex;not null"`
	// 密码，SHA1 格式
	Password *string `gorm:"type:varchar(40);not null"`
}

// GetUser 查询单个
func GetUser(id string) (*User, error) {
	m := new(User)
	err := _db.
		Where("`ID` = ?", id).
		First(m).Error
	if err != nil {
		return nil, err
	}
	return m, nil
}

// GetUserByAccount 查询单个
func GetUserByAccount(account string) (*User, error) {
	m := new(User)
	err := _db.
		Where("`Account` = ?", account).
		First(m).Error
	if err != nil {
		return nil, err
	}
	return m, nil
}

// AddUser 添加单个
func AddUser(m *User) (int64, error) {
	db := _db.Create(m)
	return db.RowsAffected, db.Error
}

// UpdateUser 修改单个
func UpdateUser(m *User) (int64, error) {
	db := _db.
		Where("`ID` = ?", m.ID).
		Updates(m)
	return db.RowsAffected, db.Error
}

// DeleteUser 删除单个
func DeleteUser(id string) (int64, error) {
	db := _db.
		Delete(&App{
			ID: id,
		})
	return db.RowsAffected, db.Error
}
