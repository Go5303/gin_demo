package model

// User maps to oa_user table
type User struct {
	ID       int    `gorm:"column:id;primaryKey" json:"id"`
	Username string `gorm:"column:username" json:"username"`
	Password string `gorm:"column:password" json:"-"`
	Nickname string `gorm:"column:nickname" json:"nickname"`
	Phone    string `gorm:"column:phone" json:"phone"`
	Status   int    `gorm:"column:status" json:"status"`
}

// TableName overrides the table name
func (User) TableName() string {
	return "oa_user"
}

// GetUserByUsername finds a user by username
func GetUserByUsername(username string) (*User, error) {
	var user User
	err := GetDB().Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID finds a user by id
func GetUserByID(id int) (*User, error) {
	var user User
	err := GetDB().Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
