package errors

/* Not Found */
type NotFound struct {
	Err error
}

func (e *NotFound) Error() string {
	return e.Err.Error()
}

/* InvalidPayload */
type InvalidPayload struct {
	Err error
}

func (e *InvalidPayload) Error() string {
	return e.Err.Error()
}

/* InternalServerError */
type InternalServerError struct {
	Err error
}

func (e *InternalServerError) Error() string {
	return e.Err.Error()
}

/* BadRequestError*/
type BadRequest struct {
	Err error
}
func (e *BadRequest) Error() string {
	return e.Err.Error()
}

/* Conflict */
type Conflict struct {
	Err error
}

func (e *Conflict) Error() string {
	return e.Err.Error()
}

/* Unauthorized */
type Unauthorized struct {
	Err error
}

func (e *Unauthorized) Error() string {
	return e.Err.Error()
}
