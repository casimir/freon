package auth

import "fmt"

type InvalidTokenError struct {
	Token  string
	Source error
}

func (e *InvalidTokenError) Error() string {
	if e.Source != nil {
		return fmt.Sprintf("invalid token: %q: %v", e.Token, e.Source)
	}
	return fmt.Sprintf("invalid token: %q", e.Token)
}

type UnknownTokenError struct {
	Token string
}

func (e *UnknownTokenError) Error() string {
	return fmt.Sprintf("unknow token: %q", e.Token)
}

type InvalidSessionError struct {
	Session string
	Source  error
}

func (e *InvalidSessionError) Error() string {
	if e.Source != nil {
		return fmt.Sprintf("invalid session: %q: %v", e.Session, e.Source)
	}
	return fmt.Sprintf("invalid session: %q", e.Session)
}

type UnknownUserError struct {
	Value string
}

func (e *UnknownUserError) Error() string {
	return fmt.Sprintf("unknown user: %q", e.Value)
}

type IncorrectPasswordError struct{}

func (e IncorrectPasswordError) Error() string {
	return "incorrect password"
}
