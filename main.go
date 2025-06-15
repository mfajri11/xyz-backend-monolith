package main

import (
	"github.com/mfajri11/xyz-backend-monolith/app"
	"github.com/mfajri11/xyz-backend-monolith/util/log"
)

func main() {
	log.Fatal(app.Run(), "fail to start application")
}
