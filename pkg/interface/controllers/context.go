package controllers

type Context interface {
	GetHeader(key string) string
	ShouldBindJSON(obj interface{}) error
	Param(string) string
	Bind(interface{}) error
	Status(int)
	JSON(int, interface{})
}

