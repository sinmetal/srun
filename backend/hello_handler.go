package backend

import (
	"fmt"
	"net/http"
	"time"
)

func (ah *AppHandlers) HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World : %s", time.Now())
}
