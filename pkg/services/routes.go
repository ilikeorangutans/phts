package services

import (
	"net/http"

	"github.com/ilikeorangutans/phts/pkg/model"
	"github.com/ilikeorangutans/phts/pkg/session"
	"github.com/ilikeorangutans/phts/pkg/smtp"
	"github.com/ilikeorangutans/phts/web"

	"github.com/jmoiron/sqlx"
)

func SetupServices(sessions session.Storage, db *sqlx.DB, emailer *smtp.Email, adminEmail, adminPassword, serverURL string) []web.Section {
	serviceUsersRepo := NewServiceUsersRepo(db)
	usersRepo := model.NewUserRepo(db)

	return []web.Section{
		{
			Path: "/services",
			Routes: []web.Route{
				{
					Path:    "/ping",
					Handler: PingHandler,
				},
				{
					Path:    "/version",
					Handler: VersionHandler,
				},
				{
					Path:    "/internal/login",
					Handler: LoginHandler,
					Methods: []string{"POST", "GET"},
				},
				{
					Path:    "/internal/sessions/create",
					Handler: AuthenticationHandler(sessions, serviceUsersRepo),
					Methods: []string{"POST"},
				},
				// TODO add /internal/sessions/refresh and /internal/sessions/check
			},
			Sections: []web.Section{
				{
					Path: "/internal",
					Middleware: []func(http.Handler) http.Handler{
						RequiresAuthentication(sessions),
					},
					Routes: []web.Route{
						{
							Path:    "/",
							Handler: LandingPageHandler,
						},
						{
							Path:    "/sessions/destroy",
							Handler: LogoutHandler(sessions, serviceUsersRepo),
							Methods: []string{"GET", "POST", "DELETE"},
						},
						{
							Path:    "/service_users",
							Handler: ServiceUsersListHandler(serviceUsersRepo),
						},
						{
							Path:    "/users",
							Handler: UsersListHandler(usersRepo),
						},
						{
							Path:    "/users/invite",
							Handler: UsersInviteHandler(usersRepo, emailer, serverURL),
							Methods: []string{"POST"},
						},
						{
							Path:    "/smtp_test",
							Handler: SmtpTestHandler(emailer, serverURL),
							Methods: []string{"GET", "POST"},
						},
					},
				},
			},
		},
	}
}
