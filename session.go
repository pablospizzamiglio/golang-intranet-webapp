package main

// Session struct to hold information about the session
type Session struct {
	Values map[interface{}]interface{}
}

// Store baseline
type Store interface {
	Get(string) (Session, error)
	Set(string, Session) error
}

// User struct to hold information retrieved from LDAP
type User struct {
	DisplayName string
	Email       string
	Username    string
}
