package web

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildNav(t *testing.T) {
	sections := []Section{
		{
			Path: "foo",
			Routes: []Route{
				{
					Path: "/",
				},
			},
		},
	}

	nav := BuildNav(sections, "/")

	req := httptest.NewRequest("GET", "/foo/bar/blargh", nil)
	breadcrumbs := nav.Breadcrumbs(req)
	assert.Equal(t, []Nav{}, breadcrumbs)
}
