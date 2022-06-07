package storage

type Storage interface {
	Read()
	Write()
}
