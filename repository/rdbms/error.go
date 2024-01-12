package rdbms

import (
	"errors"
	"fmt"
	"github.com/lib/pq"
	"io"
	"os"
)

const (
	uniqueViolationCode = "23505"
)

func showError(err error) {
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "close error:", err)
	}
}
func toPqError(err error) *pq.Error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr
	}
	return nil
}

func closeAndShowError(closable io.Closer) {
	showError(closable.Close())
}
func isErrorUniqueViolation(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == uniqueViolationCode
	}
	return false
}
