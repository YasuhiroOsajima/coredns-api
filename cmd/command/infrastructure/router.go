package infrastructure

import (
	"encoding/json"
	"fmt"

	"coredns_api/pkg/interface/controllers"
)

type CommandContext struct{}

func (c *CommandContext) GetHeader(key string) string {
	return ""
}
func (c *CommandContext) ShouldBindJSON(obj interface{}) error {
	return nil
}
func (c *CommandContext) Param(key string) string {
	return ""
}
func (c *CommandContext) Bind(obj interface{}) error {
	return nil
}
func (c *CommandContext) Status(code int) {
	fmt.Println(code)
}
func (c *CommandContext) JSON(code int, obj interface{}) {
	data, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}

// Functions in this router has to be not processes with updating domain data.
// They can only read functions.
func Router() {
	tcntr := InitializeTenantController()
	var c controllers.Context = &CommandContext{}
	tcntr.List(c)
}
