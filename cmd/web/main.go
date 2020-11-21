package main

import "coredns_api/cmd/web/infrastructure"

// @title CoreDNS API
// @version 1.0
// @description REST API to manage domain info on CoreDNS

// @termsOfService http://swagger.io/terms/

// @contact.name Hogehoge
// @contact.url http://hoge.hogehoge.hoge/support
// @contact.email support@hogehoge.hoge

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 172.28.21.40:8080
// @BasePath /v1
func main() {
	infrastructure.Router()
}
