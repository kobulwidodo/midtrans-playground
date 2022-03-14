package utils

import (
	"go-midtrans/domain"
	"net/http"
)

func GetErrCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch err {
	case domain.ErrEmailConflict:
		return http.StatusConflict
	case domain.ErrInternalServer:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrPassNotMatch:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
