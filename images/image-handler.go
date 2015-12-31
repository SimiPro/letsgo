package images

type ImageHandler interface {
	Upload() Image
}

type Image struct {
	Id string
}
