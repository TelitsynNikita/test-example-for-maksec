package main

import "github.com/TelitsynNikita/test-example-for-maksec/internal/app"

// @title           Script Monitor API
// @version         1.0
// @description     Service for remote script management via SSH
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@script-monitor.io

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the token.

func main() {
	app := app.New()
	app.Run()
}
