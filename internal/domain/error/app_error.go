package domainerrors

type AppError struct {
	code     *int
	message  string
	canShow  bool
	entity   string
	function string
	needTr   bool
}

func (e *AppError) Error() string {
	return e.message
}

func (e *AppError) Code() *int {
	return e.code
}

func (e *AppError) Message() string {
	return e.message
}

func (e *AppError) CanShow() bool {
	return e.canShow
}

func (e *AppError) Entity() string {
	return e.entity
}

func (e *AppError) Function() string {
	return e.function
}
func New(code *int, msg string, entity, fn string) *AppError {
	return &AppError{
		code:     code,
		message:  msg,
		canShow:  true,
		entity:   entity,
		function: fn,
		needTr:   true,
	}
}
func NewWithCode(code int, msg string, entity, fn string) *AppError {
	return &AppError{
		code:     &code,
		message:  msg,
		canShow:  true,
		entity:   entity,
		function: fn,
		needTr:   true,
	}
}
func NewWithCodeNoTr(code int, msg string, entity, fn string) *AppError {
	return &AppError{
		code:     &code,
		message:  msg,
		canShow:  true,
		entity:   entity,
		function: fn,
		needTr:   false,
	}
}
func Internal(entity, fn string) *AppError {
	code := 500
	return &AppError{
		code:     &code,
		message:  "Internal server error",
		canShow:  false,
		entity:   entity,
		function: fn,
		needTr:   true,
	}
}
