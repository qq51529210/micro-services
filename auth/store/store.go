package store

var (
	_Store Store
)

type Store interface {
}

func SetStore(st Store) {
	_Store = st
}

func GetStore() Store {
	return _Store
}
