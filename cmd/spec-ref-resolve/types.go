package main

type openapiMapKey interface {
	int | string
}
type stringAnyMap = map[string]any
type intAnyMap = map[int]any

type anyMap2 interface {
	stringAnyMap | intAnyMap
}

type anyMap[T openapiMapKey] map[T]any

func (am *anyMap[T]) Set(key any, val any) {
	(*am)[key.(T)] = val
}
func (am *anyMap[T]) Del(key any) {
	delete(*am, key.(T))
}
func (am *anyMap[T]) Exists(key any) bool {
	_, exists := (*am)[key.(T)]
	return exists
}

func anymapToGeneric[T openapiMapKey](in map[any]any) anyMap[T] {
	var out anyMap[T] = make(anyMap[T], len(in))
	for k, v := range in {
		castedKey, casted := k.(T)
		if !casted {
			continue
		}
		out[castedKey] = v
	}
	return out
}
