package serve

import (
	"log"
	"net/http"
)

func Serve(port string) {
	log.Fatal(http.ListenAndServe(":"+port, http.FileServer(http.Dir("out/"))))
}
