package services

import "html/template"

var BaseTmpl = template.Must(template.ParseFiles("pkg/services/base.tmpl"))

var LoginPageTmpl = template.Must(template.Must(BaseTmpl.Clone()).ParseFiles("pkg/services/login_page.tmpl"))

var LandingPageTmpl = template.Must(template.Must(BaseTmpl.Clone()).ParseFiles("pkg/services/landing_page.tmpl"))

var ServiceUsersPageTmpl = template.Must(template.Must(BaseTmpl.Clone()).ParseFiles("pkg/services/service_users_page.tmpl"))
