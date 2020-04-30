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

func BaseTmpl() *template.Template {
	return template.Must(template.ParseFiles("templates/services/internal/base.tmpl")).Funcs(templateFuncs)
}

func LoginPageTmpl() *template.Template {
	return template.Must(template.Must(BaseTmpl().Clone()).ParseFiles("templates/services/internal/login_page.tmpl"))
}

func BaseUITmpl() *template.Template {
	return template.Must(template.Must(BaseTmpl().Clone()).ParseFiles("templates/services/internal/base_ui.tmpl"))
}

func LandingPageTmpl() *template.Template {
	return template.Must(template.Must(BaseUITmpl().Clone()).ParseFiles("templates/services/internal/landing_page.tmpl"))
}

func ServiceUsersPageTmpl() *template.Template {
	return template.Must(template.Must(BaseUITmpl().Clone()).ParseFiles("templates/services/internal/service_users_page.tmpl"))
}

func UsersPageTmpl() *template.Template {
	return template.Must(template.Must(BaseUITmpl().Clone()).ParseFiles("templates/services/internal/users_page.tmpl"))
}

func SmtpTestTmpl() *template.Template {
	return template.Must(template.Must(BaseUITmpl().Clone()).ParseFiles("templates/services/internal/smtp_test.tmpl"))
}

func UserInviteEmailTmpl() *template.Template {
	return template.Must(template.ParseFiles("templates/services/internal/user_invite_email.tmpl"))
}
