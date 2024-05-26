package cache

type Cache interface {
	Set(key, otp string) error
	Get(key string) (string, error)
	Delete(id string) error
}
