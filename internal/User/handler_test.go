package User

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func Test_userHandler_login(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		u    *userHandler
		args args
	}{
		{name: "Test_userHandler_login", u: &userHandler{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.u.login(tt.args.c)
		})
	}
}
