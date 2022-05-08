package User

import (
	"reflect"
	"testing"

	"github.com/BatuhanSerin/final-project/internal/models"
)

func TestUserRepository_createUser(t *testing.T) {
	type args struct {
		user *models.User
	}
	tests := []struct {
		name    string
		u       *UserRepository
		args    args
		want    *models.User
		wantErr bool
	}{
		{name: "TestUserRepository_createUser", u: &UserRepository{}, args: args{user: &models.User{Phone: "12"}}, want: &models.User{}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.u.createUser(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.createUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserRepository.createUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserRepository_getUser(t *testing.T) {
	type args struct {
		email    string
		password string
	}
	tests := []struct {
		name    string
		u       *UserRepository
		args    args
		want    *models.User
		wantErr bool
	}{
		{name: "TestUserRepository_getUser", u: &UserRepository{}, args: args{email: "admin@gmail.com", password: "1"}, want: &models.User{}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.u.getUser(&tt.args.email, &tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.getUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserRepository.getUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
