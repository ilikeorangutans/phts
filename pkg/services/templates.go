package services

import "html/template"

var BaseTmpl = template.Must(template.ParseFiles("templates/services/internal/base.tmpl"))

var LoginPageTmpl = template.Must(template.Must(BaseTmpl.Clone()).ParseFiles("templates/services/internal/login_page.tmpl"))

var BaseUITmpl = template.Must(template.Must(BaseTmpl.Clone()).ParseFiles("templates/services/internal/base_ui.tmpl"))

var LandingPageTmpl = template.Must(template.Must(BaseUITmpl.Clone()).ParseFiles("templates/services/internal/landing_page.tmpl"))

var ServiceUsersPageTmpl = template.Must(template.Must(BaseUITmpl.Clone()).ParseFiles("templates/services/internal/service_users_page.tmpl"))

var UsersPageTmpl = template.Must(template.Must(BaseUITmpl.Clone()).ParseFiles("templates/services/internal/users_page.tmpl"))

var SmtpTestTmpl = template.Must(template.Must(BaseUITmpl.Clone()).ParseFiles("templates/services/internal/smtp_test.tmpl"))
