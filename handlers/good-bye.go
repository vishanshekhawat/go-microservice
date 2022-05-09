package handlers

import (
	"log"
	"net/http"
)

type GoodBye struct {
	logger *log.Logger
}

func NewGoodBye(l *log.Logger) *GoodBye {
	return &GoodBye{
		logger: l,
	}
}
func (g *GoodBye) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	g.logger.Println("Good Bye")
	rw.Write([]byte("Good Bye"))
}
