package modeltest

import (
	"testing"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/test/integration"
)

func TestFoo(t *testing.T) {
	integration.RunTestInDB(t, func(dbx db.DB) {
		createTestCollection(t, dbx)

	})

}
