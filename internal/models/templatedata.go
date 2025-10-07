package models

import "github.com/bayramovrahman/fastnet_vpn_bot/internal/forms"

type TemplateData struct {
	StringMap       map[string]string
	IntMap          map[string]int64
	FloatMap        map[string]float32
	Data            map[string]interface{}
	CsrfToken       string
	Form            *forms.Form
	Flash           string
	Warning         string
	Error           string
	IsAuthenticated int
}
