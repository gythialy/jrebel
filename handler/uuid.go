package handler

import (
	"fmt"
	"net/http"

	"github.com/satori/go.uuid"
)

func UUID(w http.ResponseWriter, _ *http.Request) {
	u := uuid.NewV4().String()
	fmt.Printf("UUIDv4: %s\n", u)
	_, _ = w.Write([]byte(u))
}
