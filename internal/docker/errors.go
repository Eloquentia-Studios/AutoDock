package docker

type ErrConnectionFailed struct {
	err error
}

func (e ErrConnectionFailed) Error() string {
	return e.err.Error()
}

type ErrNotFound struct {
	err error
}

func (e ErrNotFound) Error() string {
	return e.err.Error()
}

type ErrVersionMismatch struct {
	err error
}

func (e ErrVersionMismatch) Error() string {
	return e.err.Error()
}
