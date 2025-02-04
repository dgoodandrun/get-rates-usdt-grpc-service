package errors

import "fmt"

var (
	ApiError      = fmt.Errorf("API unavailable")
	StorageError  = fmt.Errorf("storage unavailable")
	SaveRateError = fmt.Errorf("failed to save rate")
)
