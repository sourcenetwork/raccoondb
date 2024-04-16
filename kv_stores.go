package raccoon

func NewMemKV() KVStore {
	return NewMemDB()
}
