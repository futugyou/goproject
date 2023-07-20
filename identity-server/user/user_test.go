package user

import "testing"

func TestUser_GetType(t *testing.T) {
	type fields struct {
		ID       string
		Name     string
		Password string
		Email    string
		Birth    string
		Phone    string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := User{
				ID:       tt.fields.ID,
				Name:     tt.fields.Name,
				Password: tt.fields.Password,
				Email:    tt.fields.Email,
				Birth:    tt.fields.Birth,
				Phone:    tt.fields.Phone,
			}
			if got := u.GetType(); got != tt.want {
				t.Errorf("User.GetType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserLogin_GetType(t *testing.T) {
	type fields struct {
		ID        string
		UserID    string
		Timestamp int64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := UserLogin{
				ID:        tt.fields.ID,
				UserID:    tt.fields.UserID,
				Timestamp: tt.fields.Timestamp,
			}
			if got := u.GetType(); got != tt.want {
				t.Errorf("UserLogin.GetType() = %v, want %v", got, tt.want)
			}
		})
	}
}
