package cqrs

type Errors []error

// Error conforms the error interface
func (e Errors) Error() string {

}

// Add will add a error to the stack
func (e Errors) Add( err error) {
	&e = append(&e, err)
}
