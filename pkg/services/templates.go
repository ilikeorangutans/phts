package services

import (
	"html/template"
	"time"

	"github.com/dustin/go-humanize"
)

const (
	DateFormat     = "2006-01-02"
	DateTimeFormat = "2006-01-02 15:04"
)

var templateFuncs = template.FuncMap{
	"humanizeTime": humanize.Time,
	"fullDateTime": func(t time.Time) string {
		return t.Format(time.RFC1123)
	},
}

var BaseTmpl = template.Must(template.ParseFiles("templates/services/internal/base.tmpl")).Funcs(templateFuncs)

var LoginPageTmpl = template.Must(template.Must(BaseTmpl.Clone()).ParseFiles("templates/services/internal/login_page.tmpl"))

var BaseUITmpl = template.Must(template.Must(BaseTmpl.Clone()).ParseFiles("templates/services/internal/base_ui.tmpl"))

var LandingPageTmpl = template.Must(template.Must(BaseUITmpl.Clone()).ParseFiles("templates/services/internal/landing_page.tmpl"))

var ServiceUsersPageTmpl = template.Must(template.Must(BaseUITmpl.Clone()).ParseFiles("templates/services/internal/service_users_page.tmpl"))

var UsersPageTmpl = template.Must(template.Must(BaseUITmpl.Clone()).ParseFiles("templates/services/internal/users_page.tmpl"))

var SmtpTestTmpl = template.Must(template.Must(BaseUITmpl.Clone()).ParseFiles("templates/services/internal/smtp_test.tmpl"))
