package users

import (
	"crypto/md5"
	"io"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

//TODO: add tests for the various functions in user.go, as described in the assignment.
//use `go test -cover` to ensure that you are covering all or nearly all of your code paths.

func TestValidate(t *testing.T) {
	cases := []struct {
		CaseName     string
		Email        string
		Password     string
		PasswordConf string
		UserName     string
		FirstName    string
		LastName     string
		expectErr    bool
	}{
		{
			"valid person",
			"validEmail@gmail.com",
			"validPasswordLengthAndMatch",
			"validPasswordLengthAndMatch",
			"validUserName",
			"validFirstName",
			"validLastName",
			false,
		},
		{
			"invalid email format",
			"validEmail.com",
			"validPasswordLengthAndMatch",
			"validPasswordLengthAndMatch",
			"validUserName",
			"validFirstName",
			"validLastName",
			true,
		},
		{
			"invalid password length",
			"validEmail@gmail.com",
			"short",
			"short",
			"validUserName",
			"validFirstName",
			"validLastName",
			true,
		},
		{
			"invalid password match",
			"validEmail@gmail.com",
			"validPasswordLengthAndMatch",
			"validPasswordLengthAndMatch123",
			"validUserName",
			"validFirstName",
			"validLastName",
			true,
		},
		{
			"invalid no username",
			"validEmail@gmail.com",
			"validPasswordLengthAndMatch",
			"validPasswordLengthAndMatch",
			"",
			"validFirstName",
			"validLastName",
			true,
		},
	}
	// Dont forget to test bad inputs? Empty NewUser
	for _, c := range cases {
		nu := &NewUser{
			Email:        c.Email,
			Password:     c.Password,
			PasswordConf: c.PasswordConf,
			UserName:     c.UserName,
			FirstName:    c.FirstName,
			LastName:     c.LastName,
		}
		errorValidate := nu.Validate()

		if errorValidate != nil && !c.expectErr {
			t.Errorf("error when validating case:\"%s\". got error when none was expected. error was: %v", c.CaseName, errorValidate)
		}

		if c.expectErr && errorValidate == nil {
			t.Errorf("error when validating case:\"%s\". expected error but none was found.", c.CaseName)
		}
	}
}

// func TestSetPassword(t *testing.T) {
// 	user := &User{}
// 	user.SetPassword("password")
// 	userPasswordHash := user.PassHash
//
// }

func TestToUser(t *testing.T) {
	nu := &NewUser{
		Email:     "KyLeIwS@uw.edu",
		FirstName: "kyle",
		LastName:  "williams-smith",
		UserName:  "kylews",
		Password:  "its ya boy",
	}
	user, err := nu.ToUser()
	if err != nil {
		t.Errorf("error when calling ToUser on NewUser instance. error: %v\n", err)
	}
	transformedEmail := strings.ToLower(strings.Trim(nu.Email, " "))
	md5Hasher := md5.New()
	io.WriteString(md5Hasher, transformedEmail)
	gravatarURL := "https://www.gravatar.com/avatar/" + string(md5Hasher.Sum(nil))
	if user.PhotoURL != gravatarURL {
		t.Errorf("error when calling ToUser on NewUser instance. Gravatar photo is not being set correctly. Make sure you always force lower case.")
	}
	if len(user.ID) == 0 {
		t.Errorf("error when calling ToUser on NewUser instance. ID field should not be empty")
	}
}

func TestFullName(t *testing.T) {
	cases := []struct {
		CaseName  string
		FirstName string
		LastName  string
		expected  string
	}{
		{
			"valid names",
			"Ethan",
			"Brendan",
			"Ethan Brendan",
		},
		{
			"only first",
			"donald",
			"",
			"donald",
		},
		{
			"only last",
			"",
			"obama",
			"obama",
		},
		{
			"no name",
			"",
			"",
			"",
		},
	}
	for _, c := range cases {
		user := &User{
			FirstName: c.FirstName,
			LastName:  c.LastName,
		}
		actualFullName := user.FullName()
		if actualFullName != c.expected {
			t.Errorf("error when getting user fullname. Got %s but expected %s", actualFullName, c.expected)
		}
	}
}

func TestAutenticate(t *testing.T) {
	pass := "password"
	PassHashShouldBe, err := bcrypt.GenerateFromPassword([]byte(pass), bcryptCost)
	if err != nil {
		t.Errorf("error when bcrypting password when one wasn't expected. got error: %v", err)
	}
	user := &User{
		PassHash: PassHashShouldBe,
	}
	if err := user.Authenticate(pass); err != nil {
		t.Errorf("error with user password authentication. error when none was expected %v", err)
	}
}

func TestApplyUpdates(t *testing.T) {
	cases := []struct {
		CaseName      string
		FirstName     string
		LastName      string
		expectedError bool
	}{
		{
			"valid update",
			"Kyle",
			"Williams-Smith",
			false,
		},
		{
			"invalid firstname",
			"",
			"Williams-Smith",
			true,
		},
		{
			"invalid lastname",
			"Kyle",
			"",
			true,
		},
		{
			"no input",
			"",
			"",
			true,
		},
	}
	for _, c := range cases {
		user := &User{
			FirstName: "Kyleee",
			LastName:  "WS",
		}
		up := &Updates{
			FirstName: c.FirstName,
			LastName:  c.LastName,
		}
		err := user.ApplyUpdates(up)
		if err != nil && !c.expectedError {
			t.Errorf("error when applying name updates when none was expected: %v", err)
		}
		if err == nil && c.expectedError {
			t.Errorf("no error when one was expected. Case: %s", c.CaseName)
		}
	}

}
