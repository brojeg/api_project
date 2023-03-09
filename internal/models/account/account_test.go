package models

import (
	db "diploma/go-musthave-diploma-tpl/internal/models/database"
	server "diploma/go-musthave-diploma-tpl/internal/models/server"
	"fmt"
	"reflect"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/suite"
)

var loginToTest = "some_username"
var passwordToTest = "some_password"

func TestAccount_Validate(t *testing.T) {

	tests := []struct {
		name    string
		account *Account
		want    server.Response
	}{
		{
			name:    "// Test case 1: Verify that the account.Validate() function returns the expected 200 code for non-existing user.",
			account: &Account{Login: loginToTest, Password: passwordToTest},
			want:    server.Response{Message: "Requirement passed", ServerCode: 200},
		},
		{
			name:    "// Test case 2: Verify that the account.Validate() function returns the expected 400 code for the account with missing Login field.",
			account: &Account{Password: passwordToTest},
			want:    server.Response{Message: "Login is not valid", ServerCode: 400},
		},
		{
			name:    "// Test case 3: Verify that the account.Validate() function returns the expected 400 code for the account with missing Password field.",
			account: &Account{Login: loginToTest},
			want:    server.Response{Message: "Valid password is required", ServerCode: 400},
		},
		{
			name:    "// Test case 4: Verify that the account.Validate() function returns the expected 400 code for the empty request.",
			account: &Account{},
			want:    server.Response{Message: "Login is not valid", ServerCode: 400},
		},
		{
			name:    "// Test case 5: Verify that the account.Validate() function returns the expected 400 code for the short Login field.",
			account: &Account{Login: "1", Password: passwordToTest},
			want:    server.Response{Message: "Login is not valid", ServerCode: 400},
		},
		{
			name:    "// Test case 6: Verify that the account.Validate() function returns the expected 400 code for the short Password field.",
			account: &Account{Login: loginToTest, Password: "pass"},
			want:    server.Response{Message: "Valid password is required", ServerCode: 400},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.account.Validate(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Account.Validate() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestAccount_Create(t *testing.T) {
	testAccount := Account{Login: loginToTest, Password: passwordToTest}
	tests := []struct {
		name    string
		account *Account
		want    server.Response
	}{
		{
			name:    "// Test case 1: Verify that the account.Create() function returns the expected 200 code for non-existing user.",
			account: &testAccount,
			want:    server.Response{ServerCode: 200},
		},
		{
			name:    "// Test case 2: Verify that the account.Create() function returns the expected 400 code for an existing user.",
			account: &testAccount,
			want:    server.Response{ServerCode: 400},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.account.Create(); !reflect.DeepEqual(got.ServerCode, tt.want.ServerCode) {
				t.Errorf("Account.Create() = %v, want %v", got.ServerCode, tt.want.ServerCode)
			}
		})
	}
	// Clear the artifact after the testing is complete
	testAccount.Delete()
}

func TestLogin(t *testing.T) {
	testAccount := Account{Login: loginToTest, Password: passwordToTest}
	testAccount.Create()
	type args struct {
		email    string
		password string
	}
	tests := []struct {
		name string
		args args
		want server.Response
	}{
		{
			name: "// Test case 1: Verify that the account.Login() function returns the expected 200 code for the existing user.",
			args: args{email: loginToTest, password: passwordToTest},
			want: server.Response{ServerCode: 200},
		},
		{
			name: "// Test case 2: Verify that the account.Login() function returns the expected 500 code for the non-existing user.",
			args: args{email: "other_login", password: passwordToTest},
			want: server.Response{ServerCode: 500},
		},
		{
			name: "// Test case 3: Verify that the account.Login() function returns the expected 401 code for the existing user and wrong password",
			args: args{email: loginToTest, password: ""},
			want: server.Response{ServerCode: 401},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Login(tt.args.email, tt.args.password); !reflect.DeepEqual(got.ServerCode, tt.want.ServerCode) {
				t.Errorf("Login() = %v, want %v", got, tt.want)
			}
		})
	}
	// Clear the artifact after the testing is complete
	testAccount.Delete()

}
func (account *Account) Delete() server.Response {

	db.Get().Where("login = ?", account.Login).Delete(&Account{})
	return server.Response{Message: "Account was deleted", ServerCode: 200}
}

type TestSuiteEnv struct {
	suite.Suite
	db      *gorm.DB
	Account Account
}

// Tests are run before they start
func (suite *TestSuiteEnv) SetupSuite() {
	db.InitDBConnectionString("host=localhost user=test_user dbname=test sslmode=disable password=111")
	suite.db = db.Get()
	suite.Account = Account{Login: loginToTest, Password: passwordToTest}
}

// Running after each test
func (suite *TestSuiteEnv) TearDownTest() {
	suite.Account.Delete()
}

// Running after all tests are completed
func (suite *TestSuiteEnv) TearDownSuite() {
	suite.db.Close()
}

// This gets run automatically by `go test` so we call `suite.Run` inside it
func TestSuite(t *testing.T) {

	suite.Run(t, new(TestSuiteEnv))
}

func (suite *TestSuiteEnv) TestAccount_Validate() {
	tests := []struct {
		name    string
		account *Account
		want    server.Response
	}{
		{
			name:    "// Test case 1: Verify that the account.Validate() function returns the expected 200 code for non-existing user.",
			account: &Account{Login: loginToTest, Password: passwordToTest},
			want:    server.Response{Message: "Requirement passed", ServerCode: 200},
		},
		{
			name:    "// Test case 2: Verify that the account.Validate() function returns the expected 400 code for the account with missing Login field.",
			account: &Account{Password: passwordToTest},
			want:    server.Response{Message: "Login is not valid", ServerCode: 400},
		},
		{
			name:    "// Test case 3: Verify that the account.Validate() function returns the expected 400 code for the account with missing Password field.",
			account: &Account{Login: loginToTest},
			want:    server.Response{Message: "Valid password is required", ServerCode: 400},
		},
		{
			name:    "// Test case 4: Verify that the account.Validate() function returns the expected 400 code for the empty request.",
			account: &Account{},
			want:    server.Response{Message: "Login is not valid", ServerCode: 400},
		},
		{
			name:    "// Test case 5: Verify that the account.Validate() function returns the expected 400 code for the short Login field.",
			account: &Account{Login: "1", Password: passwordToTest},
			want:    server.Response{Message: "Login is not valid", ServerCode: 400},
		},
		{
			name:    "// Test case 6: Verify that the account.Validate() function returns the expected 400 code for the short Password field.",
			account: &Account{Login: loginToTest, Password: "pass"},
			want:    server.Response{Message: "Valid password is required", ServerCode: 400},
		},
	}
	for _, tt := range tests {
		a := suite.Assert()
		got := tt.account.Validate()
		a.Equal(got, tt.want, fmt.Sprintf("Account.Validate() = %v, want %v", got, tt.want))
		// if got := tt.account.Validate(); !reflect.DeepEqual(got, tt.want) {
		// 	t.Errorf("Account.Validate() = %v, want %v", got, tt.want)
		// }

	}
}
