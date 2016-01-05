package common_test

import (
	"testing"
	"encoding/json"
	"github.com/letsgoli/common"
	"bytes"
	"github.com/stretchr/testify/assert"
)

func TestMarshall(t *testing.T) {
	content := []byte(`{
		"password": "test",
		"image": "no image",
		"email": "Jillian@example.com",
		"lastName": "Martha",
		"firstName": "Sharlene",
		"username": "Ingram",
		"id": "1"
	}`)

	user := new(common.User)
	err := json.NewDecoder(bytes.NewReader(content)).Decode(user)
	assert.Nil(t, err)
	assert.Equal(t,"test", user.Password)
	assert.Equal(t,"no image", user.Image)
	assert.Equal(t,"Jillian@example.com", user.Email)
	assert.Equal(t,"Martha", user.Lastname)
	assert.Equal(t,"Sharlene", user.Firstname)
	assert.Equal(t,"Ingram", user.Username)
	assert.Equal(t,"1", user.Id)
}