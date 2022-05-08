package User

import (
	"reflect"
	"testing"

	"github.com/BatuhanSerin/final-project/internal/api"
	"github.com/BatuhanSerin/final-project/internal/models"
)

func Test_responseToUser(t *testing.T) {
	type args struct {
		u *api.User
	}
	tests := []struct {
		name string
		args args
		want *models.User
	}{

		{name: "Test_responseToUser", args: args{u: &api.User{Email: ""}}, want: &models.User{Email: ""}},
		{name: "Test_responseToUser", args: args{u: &api.User{Email: "a@hotmail.com"}}, want: &models.User{Email: "a@hotmail.com"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := responseToUser(tt.args.u); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("responseToUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserToResponse(t *testing.T) {
	type args struct {
		u *models.User
	}
	tests := []struct {
		name string
		args args
		want *api.User
	}{
		{name: "TestUserToResponse", args: args{u: &models.User{Email: ""}}, want: &api.User{Email: ""}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UserToResponse(tt.args.u); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserToResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}
