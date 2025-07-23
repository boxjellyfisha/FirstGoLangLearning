package main

import (
	// "kzapp/exercise"

	"kzapp/webapi"
)

//	@title			First Server swagger
//	@version		1.0-test02
//	@description	the base server apis
//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html
//	@host			localhost:80
// schemes http
/* /Users/kellyhong/go/bin/swag init -g cmd/main.go -d . -o ./cmd/docs */

func main() {
	// exercise.RunLeetCode()
	// exercise.EnterLottery(100, 5)

	// webapi.RunMuxServer()
	webapi.RunGinServer()
}
