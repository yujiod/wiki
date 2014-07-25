package controllers

import "github.com/revel/revel"

type App struct {
	GormController
}

func (c App) Index() revel.Result {
	return c.Render()
}
