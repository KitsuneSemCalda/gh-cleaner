package structures

// This struct contains a private login and private token builded in fabric method
// it's more secure than public field
type Login struct {
	login string
	token string
}

// Method Fabric to create the structure Login
func CreateLogin(l string, t string) Login {
	return Login{
		login: l,
		token: t,
	}
}

// Getter Method to catch the builded token
func (l *Login) GetToken() string {
	if (l.token) != "" {
		return l.token
	}

	return ""
}

// Getter Method to catch the builded login
func (l *Login) GetLogin() string {
	if (l.login) != "" {
		return l.login
	}

	return ""
}
