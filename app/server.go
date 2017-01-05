package madmin

import (
	"fmt"
	"net/http"
)

type DalekHandler struct {}

func (f DalekHandler) ServeHTTP (rw ResponseWriter, req *Request) {}