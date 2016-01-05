package common_test
import (
	"testing"
	"github.com/letsgoli/common"
	"encoding/json"
	"bytes"
	"github.com/stretchr/testify/assert"
)


func TestDecodeStop(t *testing.T) {
	content := []byte(`{
		"id" : "1",
		"name" : "stopname",
		"Lat" : 47.3775448,
		"Long" : 8.5478786,
		"images" : []
	}`)

	stop := new(common.Stop)
	err := json.NewDecoder(bytes.NewReader(content)).Decode(stop)

	assert.Nil(t, err)
	assert.Equal(t, "1", stop.Id)
	assert.Equal(t, 47.3775448, stop.Lat)
	assert.Equal(t, 8.5478786, stop.Long)
	assert.Equal(t, "stopname", stop.Name)
}

func TestDecodeStopWithImages(t *testing.T) {
	content := []byte(`{
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
	}`)

	stop := new(common.Stop)
	err := json.NewDecoder(bytes.NewReader(content)).Decode(stop)

	assert.Nil(t, err)
	assert.Equal(t, "1", stop.Id)
	assert.Equal(t, 47.3775448, stop.Lat)
	assert.Equal(t, 8.5478786, stop.Long)
	assert.Equal(t, "stopname", stop.Name)

	assert.Equal(t, "1", stop.Images[0].Id)
	assert.Equal(t, "s3://blabla", stop.Images[0].Path)
	assert.Equal(t, "comment to my amazing stop there", stop.Images[0].Story)

	assert.Equal(t, "2", stop.Images[1].Id)
	assert.Equal(t, "s2://blabla", stop.Images[1].Path)
	assert.Equal(t, "commentli to my amazing stop there", stop.Images[1].Story)
}