package web

import (
	"fmt"
	"net/http"
)

func BuildNav(sections []Section, base string) Nav {
	return Nav{}
}

type Nav struct {
	Title    string
	Path     string
	Children []Nav
}

func (n Nav) Breadcrumbs(r *http.Request) []Nav {
	fmt.Printf("%s - %s", n.Path, r.RequestURI)
	return []Nav{}
}
