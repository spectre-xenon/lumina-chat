package util

// return a refernce to the variable
// useful when wanting to use constants as ref
func Of[E any](e E) *E {
	return &e
}
