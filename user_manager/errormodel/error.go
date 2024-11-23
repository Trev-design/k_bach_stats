package errormodel

import "fmt"

type ErrorModel struct {
	ErrorMesssage    error  `json:"error_message"`
	Entity           string `json:"entity"`
	Timestamp        string `json:"timestamp"`
	Line             int    `json:"line"`
	ErrorType        int    `json:"error_type"`
	ErrorDescription string `json:"error_description"`
	Serverity        string `json:"serverity"`
	FunctionName     string `json:"function_name"`
	Stacktrace       string `json:"stacktrace"`
}

func (err *ErrorModel) Error() string {
	return fmt.Sprintf(
		"%s error occured in line %d\n%s\nin function %s\n\nstacktrace: %s\n\n%s\n%s\n%d\nentity: %s",
		err.Timestamp,
		err.Line,
		err.ErrorMesssage.Error(),
		err.FunctionName,
		err.Stacktrace,
		err.ErrorDescription,
		err.Serverity,
		err.ErrorType,
		err.Entity,
	)
}
