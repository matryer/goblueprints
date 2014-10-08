package meander

// Facade represents objects that provide a public view of
// themselves.
type Facade interface {
	Public() interface{}
}

// Public gets the public representation of the specified object
// if it implements the Facade interface. Otherwise returns the
// object untouched.
func Public(o interface{}) interface{} {
	if p, ok := o.(Facade); ok {
		return p.Public()
	}
	return o
}
