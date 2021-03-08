package customerr


type BaseErr struct {
	Op  string
	Msg string
	Err error
}

func (c BaseErr) Error() string {
	return c.Err.Error()
}
