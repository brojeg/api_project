package models

import (
	"reflect"
	"testing"
)

func TestAccount_Validate(t *testing.T) {
	tests := []struct {
		name    string
		account *Account
		want    Response
	}{
		{
			name:    "Account not found test",
			account: &Account{Login: "testUser", Password: "123456"},
			want:    Response{ServerCode: 200, Message: "Requirement passed"},
		},
		{
			name:    "Short password test",
			account: &Account{Login: "testUser", Password: "1"},
			want:    Response{ServerCode: 400, Message: "Password is required"},
		},
		{
			name:    "No login test",
			account: &Account{Password: "1"},
			want:    Response{ServerCode: 400, Message: "Login is not valid"},
		},
		{
			name:    "No password test",
			account: &Account{Login: "TestUser"},
			want:    Response{ServerCode: 400, Message: "Password is required"},
		},
		{
			name:    "Existing login test",
			account: &Account{Login: "TestUser", Password: "123456"},
			want:    Response{ServerCode: 409, Message: "Email address already in use by another user."},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "Existing login test" {
				if resp := tt.account.Validate(); resp.ServerCode == 200 {
					tt.account.Create()
				}
			}
			if got := tt.account.Validate(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Account.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
