package handler

import (
	"fmt"
	"github.com/satori/go.uuid"
	"net/http"
)

func UUID(w http.ResponseWriter, req *http.Request) {
	u := uuid.NewV4().String()
	fmt.Printf("UUIDv4: %s\n", u)
	_, _ = w.Write([]byte(u))
}
