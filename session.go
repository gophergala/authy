package authy

// The session object is used to store the CSRF token used by OAuth2
type Session interface {
	// Get a key from the session
	Get(key interface{}) interface{}
	// Set a key in the session
	Set(key interface{}, value interface{})
	// Unset a key from the session
	Delete(key interface{})
}
