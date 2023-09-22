package main

type CustomError struct {
	message string
	trace   []int
}

func newCustomError(m string) *CustomError {
	return &CustomError{message: m}
}

func (cm *CustomError) addTrace(pos int) *CustomError {
	cm.trace = append(cm.trace, pos)
	return cm
}
