package service_test

import (
	"atc/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	// Call the constructor
	s := service.NewService()

	// Verify that the service is not nil
	assert.NotNil(t, s)

	// Verify that the service's fields are initialized correctly
	assert.NotNil(t, s.Config)
	assert.NotNil(t, s.Log)
	assert.NotNil(t, s.Backend)
	assert.NotNil(t, s.Web.Handle)
}
