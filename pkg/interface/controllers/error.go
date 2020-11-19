package controllers

// type Error struct {
// 	Message string
// }
//
// func NewError(err error) *Error {
// 	return &Error{
// 		Message: err.Error(),
// 	}
// }

func NewError(ctx Context, status int, err error) {
	er := HTTPError{
		Code:    status,
		Message: err.Error(),
	}
	ctx.JSON(status, er)
}

// HTTPError example
type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}
