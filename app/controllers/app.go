package controllers

import "github.com/revel/revel"

type App struct {
	GormController
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) ForwardAction(methodName string) {
	c.MethodName = methodName
	c.MethodType = c.Type.Method(c.MethodName)
	c.Action = c.Name + "." + c.MethodName
}
