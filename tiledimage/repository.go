package tiledimage

// Repository - key/value хранилище.
// Ключом является string, так как он гарантированно имеет хэш и может использоваться в качестве ключа map.
type Repository interface {
	Add(key string, value interface{})
	Get(key string) (interface{}, error)
	Delete(key string) error
}
