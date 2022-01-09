package controllers

import "gastb.ar/views"

func NewStatic() *Static {
	return &Static{
		Home:     views.NewView(
			  "bootstrap", "static/home"),
		Profile:  views.NewView(
			  "bootstrap", "static/profile"),
		}
	}

type Static struct {
	Home    *views.View
	Profile *views.View
}
