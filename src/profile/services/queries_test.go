package services_test

import (
	"github.com/stretchr/testify/assert"
	"profile/services"
	"testing"
)

func TestGetQueries(t *testing.T) {
	t.Run("should retrieve queries map", func(t *testing.T) {
		var queries map[string]string = services.GetQueries()
		localAssert := assert.New(t)

		localAssert.NotEmpty(queries["GetProfile"])
		localAssert.NotEmpty(queries["UpdateProfileData"])
	})
}
