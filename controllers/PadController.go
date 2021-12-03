package controllers

import (
	"net/http"

	"go-pad/models"
)

type PadController struct {
	pm *models.PadModel
}

func NewPadController(pm *models.PadModel) *PadController {
	return &PadController{
		pm: pm,
	}
}

func (pc *PadController) ServeHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "home.html")
}


func (pc *PadController) ServePad(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "pad.html")
}
