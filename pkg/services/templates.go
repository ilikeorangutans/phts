package services

import "html/template"

var LoginPageTmpl = template.Must(template.New("").Parse(`<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <title>services/internal login</title>
  </head>
  <body>
	<form method="POST" action="/services/internal/sessions/create">
	  <label for="email">email</label><input type="email" name="email" id="email">
	  <label for="password">password</label><input type="password" name="password" id="password">
	  <button type="submit">Login</button>
	</form>
  </body>
</html>
`))

var LandingPageTmpl = template.Must(template.New("").Parse(`<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <title>services/internal</title>
  </head>
  <body>
	<a href="/services/internal/sessions/destroy">Logout</a>
	<h1>phts services/internal</h1>
  </body>
</html>
`))
