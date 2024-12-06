package errormodel

import (
	"fmt"
)

type ErrorModel struct {
	ErrorMesssage    string `json:"error_message"`
	Entity           string `json:"entity"`
	Timestamp        string `json:"timestamp"`
	Line             int    `json:"line"`
	ErrorType        int    `json:"error_type"`
	ErrorDescription string `json:"error_description"`
	Serverity        string `json:"serverity"`
	FunctionName     string `json:"function_name"`
	Stacktrace       string `json:"stacktrace"`
}

func NewErrorMessage(
	errMsg,
	entity,
	timestamp,
	description,
	serverity,
	funcName,
	stacktrace string,
	line,
	errorType int,
) *ErrorModel {
	return &ErrorModel{
		ErrorMesssage:    errMsg,
		Entity:           entity,
		Timestamp:        timestamp,
		ErrorDescription: description,
		Serverity:        serverity,
		FunctionName:     funcName,
		Stacktrace:       stacktrace,
		Line:             line,
		ErrorType:        errorType,
	}
}

func (err *ErrorModel) Error() string {
	return fmt.Sprintf(
		"%s error occured in line %d\n%s\nin function %s\n\nstacktrace: %s\n\n%s\n%s\n%d\nentity: %s",
		err.Timestamp,
		err.Line,
		err.ErrorMesssage,
		err.FunctionName,
		err.Stacktrace,
		err.ErrorDescription,
		err.Serverity,
		err.ErrorType,
		err.Entity,
	)
}
