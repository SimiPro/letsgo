package images

import (
	""
)

type Image struct {
}

type ImageS3Handler struct {
}

func (i *ImageS3Handler) Register(container *restful.Container) {
	ws := new(restful.WebService)
}

func (i *ImageS3Handler) UploadObjectToS3(fileName string) {
	
}
