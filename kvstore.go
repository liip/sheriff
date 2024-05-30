package sheriff

// kvStore is the default implementation of the KVStore interface that sheriff converts a struct into.
// It is the fastest option, but does result in a re-ordering of the final JSON properties.
type kvStore map[string]interface{}

// Set inserts the value into the map at the given key.
func (m kvStore) Set(k string, v interface{}) {
	m[k] = v
}

// Each applies the callback function to each element in the map.
func (m kvStore) Each(f func(k string, v interface{})) {
	for k, v := range m {
		f(k, v)
	}
}
