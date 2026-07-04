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

func (l *Login) GetToken() string {
	return l.token
}

func (l *Login) GetLogin() string {
	return l.login
}
