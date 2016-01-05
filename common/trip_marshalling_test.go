package common_test

import (
	"testing"
	"github.com/letsgoli/common"
	"encoding/json"
	"bytes"
	"github.com/stretchr/testify/assert"
)


func TestDecodeWithStop(t *testing.T) {
	content := []byte(`{
		"id": "1",
		"user": {
			"password": "test",
			"image": "no image",
			"email": "Jillian@example.com",
			"lastName": "Martha",
			"firstName": "Sharlene",
			"username": "Ingram",
			"id": "1"
		},
		"stops": [
			{
					"id" : "1",
					"name" : "stopname",
					"Lat" : 47.3775448,
					"Long" : 8.5478786,
					"images" : [
						{
							"id": "1",
							"path":"s3://blabla",
							"story":"comment to my amazing stop there"
						},
						{
							"id": "2",
							"path":"s2://blabla",
							"story":"commentli to my amazing stop there"
						}
					]
			},
			{
					"id" : "2",
					"name" : "stopname",
					"Lat" : 1.12312,
					"Long" : 8.5478786,
					"images" : [
						{
							"id": "1",
							"path":"s3://blabla",
							"story":"comment to my amazing stop there"
						},
						{
							"id": "2",
							"path":"s2://blabla",
							"story":"commentli to my amazing stop there"
						}
					]
			}


		]
	}`)

	trip := new(common.Trip)
	err := json.NewDecoder(bytes.NewReader(content)).Decode(trip)
	assert.Nil(t, err)
	assert.Equal(t, "1", trip.Id)
	assert.Equal(t,"test", trip.User.Password)
	assert.Equal(t,"no image", trip.User.Image)
	assert.Equal(t,"Jillian@example.com", trip.User.Email)
	assert.Equal(t,"Martha", trip.User.Lastname)
	assert.Equal(t,"Sharlene", trip.User.Firstname)
	assert.Equal(t,"Ingram", trip.User.Username)
	assert.Equal(t,"1", trip.User.Id)

	assert.Equal(t, "1", trip.Stops[0].Id)
	assert.Equal(t, 47.3775448, trip.Stops[0].Lat)
	assert.Equal(t, 8.5478786, trip.Stops[0].Long)
	assert.Equal(t, "stopname", trip.Stops[0].Name)

	assert.Equal(t, "1", trip.Stops[0].Images[0].Id)
	assert.Equal(t, "s3://blabla", trip.Stops[0].Images[0].Path)
	assert.Equal(t, "comment to my amazing stop there", trip.Stops[0].Images[0].Story)

	assert.Equal(t, "2", trip.Stops[0].Images[1].Id)
	assert.Equal(t, "s2://blabla", trip.Stops[0].Images[1].Path)
	assert.Equal(t, "commentli to my amazing stop there", trip.Stops[0].Images[1].Story)

	assert.Equal(t, "2", trip.Stops[1].Id)
	assert.Equal(t, 1.12312, trip.Stops[1].Lat)
	assert.Equal(t, 8.5478786, trip.Stops[1].Long)
	assert.Equal(t, "stopname", trip.Stops[1].Name)

	assert.Equal(t, "1", trip.Stops[1].Images[0].Id)
	assert.Equal(t, "s3://blabla", trip.Stops[1].Images[0].Path)
	assert.Equal(t, "comment to my amazing stop there", trip.Stops[0].Images[0].Story)

	assert.Equal(t, "2", trip.Stops[1].Images[1].Id)
	assert.Equal(t, "s2://blabla", trip.Stops[1].Images[1].Path)
	assert.Equal(t, "commentli to my amazing stop there", trip.Stops[1].Images[1].Story)
}



func TestDecodeUser(t *testing.T) {
	content := []byte(`{
		"id": "1",
		"user": {
			"password": "test",
			"image": "no image",
			"email": "Jillian@example.com",
			"lastName": "Martha",
			"firstName": "Sharlene",
			"username": "Ingram",
			"id": "1"
		},
		"stops": []
	}`)

	trip := new(common.Trip)
	err := json.NewDecoder(bytes.NewReader(content)).Decode(trip)
	assert.Nil(t, err)
	assert.Equal(t, "1", trip.Id)
	assert.Equal(t,"test", trip.User.Password)
	assert.Equal(t,"no image", trip.User.Image)
	assert.Equal(t,"Jillian@example.com", trip.User.Email)
	assert.Equal(t,"Martha", trip.User.Lastname)
	assert.Equal(t,"Sharlene", trip.User.Firstname)
	assert.Equal(t,"Ingram", trip.User.Username)
	assert.Equal(t,"1", trip.User.Id)
}

func TestDecodeEmptyTrip(t *testing.T) {
	content := []byte(`{
		"id": "1",
		"user": {},
		"stops": []
	}`)

	trip := new(common.Trip)
	err := json.NewDecoder(bytes.NewReader(content)).Decode(trip)
	assert.Nil(t, err)
	assert.Equal(t, "1", trip.Id)
}

