package main

import (
	"flag"
	"log"
	"strings"
)

var (
	flagFilepath    *string = flag.String(`f`, ``, `path to spec file`)
	flagOutFilepath *string = flag.String(`o`, `/dev/stdout`, `path to out spec file`)
)

func main() {
	flag.Parse()
	spec := readYml(*flagFilepath)
	resolveRefs(&spec, &spec)
	writeYaml(&spec)
	return
}

const (
	keyComponents    string = "components"
	keySchemas       string = "schemas"
	keyRequestBodies string = "requestBodies"
	responses        string = "requestBodies"

	keyRef     string = `$ref`
	symSharp   string = `#`
	skipPrefix string = `x-`
)

var createdPaths map[string]bool = make(map[string]bool)

func createPath[T openapiMapKey](dict *anyMap[T], keyPath string, value anyMap[string]) {
	if _, exist := createdPaths[keyPath]; exist {
		return
	}
	lastDictPtr := dict
	nextDictPtr := dict
	var key string
	for _, key = range strings.Split(keyPath, "/") {
		lastDictPtr = nextDictPtr
		if !nextDictPtr.Exists(key) {
			val := make(anyMap[T])
			nextDictPtr.Set(key, val)
		}
		content := (*nextDictPtr)[any(key).(T)]
		val, casted := content.(anyMap[T])
		if !casted {
			anymap, casted := content.(map[any]any)
			if !casted {
				log.Fatalf("key %s is not fomething mappable (%T)", key, content)
			}
			val = anymapToGeneric[T](anymap)
		}
		nextDictPtr = &val
	}
	lastDictPtr.Set(key, value)
	createdPaths[key] = true
}
