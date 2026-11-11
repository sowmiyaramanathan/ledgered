package apperrors

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

var (
	SSameAccoutErr        *AppError
	SInvalidSignAmountErr *AppError
	SEntryBalanceErr      *AppError
	SInvalidAmountErr     *AppError
	SAccountNotFoundErr   *AppError
	SCurrencyMismatchErr  *AppError

	DBError *AppError
)

func RegisterAppErrors() {
	SSameAccoutErr = &AppError{
		Title:      "Transaction Between Same Account",
		Detail:     "source and destination must differ",
		Code:       "3501",
		StatusCode: http.StatusBadRequest,
		IsInternal: false,
	}
	SInvalidSignAmountErr = &AppError{
		Title:      "Invalid Amount Sign",
		Detail:     "invalid sign for amount",
		Code:       "3502",
		StatusCode: http.StatusBadRequest,
		IsInternal: false,
	}
	SEntryBalanceErr = &AppError{
		Title:      "Entry not Balanced",
		Detail:     "Balance is not zero",
		Code:       "3503",
		StatusCode: http.StatusBadRequest,
		IsInternal: false,
	}
	SInvalidAmountErr = &AppError{
		Title:      "Invalid amount",
		Detail:     "Amount is not valid",
		Code:       "3504",
		StatusCode: http.StatusBadRequest,
		IsInternal: false,
	}
	SAccountNotFoundErr = &AppError{
		Title:      "Transaction not completed",
		Detail:     "Transaction not completed - Account not found",
		Code:       "3505",
		StatusCode: http.StatusBadRequest,
		IsInternal: false,
	}
	SCurrencyMismatchErr = &AppError{
		Title:      "Transaction not completed",
		Detail:     "Transaction not completed - Currency mismatch",
		Code:       "3506",
		StatusCode: http.StatusBadRequest,
		IsInternal: false,
	}

	DBError = &AppError{
		Title:      "DB Connection Issue",
		Detail:     "DB Error - %s",
		Code:       "4501",
		StatusCode: http.StatusInternalServerError,
		IsInternal: true,
	}
}

type AppError struct {
	Title      string
	Detail     string
	StatusCode int
	IsInternal bool
	Code       string
	Cause      error
}

func (c *AppError) Error() string {
	return c.Cause.Error()
}

func (a *AppError) SetDetailParameters(params ...any) *AppError {
	if strings.Count(a.Detail, "%") != len(params) {
		slog.Warn("Number of parameters does not match the number of format specifiers in the error message", "error", a)
	}
	e := *a
	e.Detail = fmt.Sprintf(a.Detail, params...)
	return &e
}

func (a *AppError) SetError(err error) *AppError {
	e := *a
	e.Cause = err
	return &e
}
