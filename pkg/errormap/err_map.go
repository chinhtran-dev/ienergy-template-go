package errormap

import (
	"net/http"

	"ienergy-template-go/pkg/constant"
)

var (
	ErrorMapCode map[int]int
	ErrorMapMsg  map[int]string
)

func Initialize() error {
	return loadData()
}

// loadData loads data from database and save memcache
func loadData() error {
	ErrorMapCode = make(map[int]int)
	ErrorMapCode = map[int]int{
		http.StatusOK:                  constant.Success,
		http.StatusBadRequest:          constant.BadRequestErr,
		http.StatusInternalServerError: constant.InternalServerError,
		http.StatusNotFound:            constant.NotFound,
		http.StatusConflict:            constant.ConflictError,
	}

	ErrorMapMsg = make(map[int]string)
	ErrorMapMsg = map[int]string{
		http.StatusOK:                  constant.SuccessMess,
		http.StatusBadRequest:          constant.BadRequestErrMess,
		http.StatusInternalServerError: constant.InternalServerErrMess,
		http.StatusNotFound:            constant.NotFoundErrMess,
		http.StatusUnauthorized:        constant.UnAuthorizedErrMess,
	}
	return nil
}
