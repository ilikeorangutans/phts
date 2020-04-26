package services

import "html/template"

var BaseTmpl = template.Must(template.ParseFiles("pkg/services/base.tmpl"))

var LoginPageTmpl = template.Must(template.Must(BaseTmpl.Clone()).ParseFiles("pkg/services/login_page.tmpl"))

var BaseUITmpl = template.Must(template.Must(BaseTmpl.Clone()).ParseFiles("pkg/services/base_ui.tmpl"))

var LandingPageTmpl = template.Must(template.Must(BaseUITmpl.Clone()).ParseFiles("pkg/services/landing_page.tmpl"))

var ServiceUsersPageTmpl = template.Must(template.Must(BaseUITmpl.Clone()).ParseFiles("pkg/services/service_users_page.tmpl"))

var UsersPageTmpl = template.Must(template.Must(BaseUITmpl.Clone()).ParseFiles("pkg/services/users_page.tmpl"))

var SmtpTestTmpl = template.Must(template.Must(BaseUITmpl.Clone()).ParseFiles("pkg/services/smtp_test.tmpl"))
