package user

type User struct {
	ID       string `bson:"_id" json:"id"`
	Name     string `bson:"name" json:"name"`
	Password string `bson:"password" json:"password"`
	Email    string `bson:"email" json:"email"`
	Birth    string `bson:"brth" json:"brth"`
	Phone    string `bson:"phone" json:"phone"`
}

type UserLogin struct {
	ID        string `bson:"_id" json:"id"`
	UserID    string `bson:"userid" json:"userid"`
	Timestamp int64  `bson:"timestamp" json:"timestamp"`
}

func (User) GetType() string {
	return "user"
}

func (UserLogin) GetType() string {
	return "userlogin"
}
