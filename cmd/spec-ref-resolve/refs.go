package main

import (
	"log"
	"path"
	"strings"
)

func resolveRefs[T openapiMapKey](root *anyMap[string], dict *anyMap[T]) {
	for k, v := range *dict {
		keyStr, keyCasted := any(k).(string)
		ref, valueCasted := any(v).(string)
		if keyCasted &&
			keyStr == keyRef &&
			valueCasted &&
			!strings.HasPrefix(ref, symSharp) {
			resolveRef(ref, root, dict)
			continue
		}
		if keyCasted && strings.HasPrefix(keyStr, skipPrefix) {
			dict.Del(k)
			continue
		}
		if nextDict, casted := v.(anyMap[string]); casted {
			resolveRefs(root, &nextDict)
			continue
		}
		if nextDict, casted := v.(map[any]any); casted {
			// @todo ymal integer map keys casted as interface{}
			(*dict)[k] = anymapToGeneric[int](nextDict)
			nextDictCasted := (*dict)[k].(anyMap[int])
			resolveRefs(root, &nextDictCasted)
			continue
		}
	}
}

var resolvedRefs map[string]string = make(map[string]string, 0)

func resolveRef[T openapiMapKey](ref string, root *anyMap[string], this *anyMap[T]) {
	if rootRefNodePath, exists := resolvedRefs[ref]; exists {
		this.Set(keyRef, createLocalRefPath(rootRefNodePath))
		return
	}
	var filePath, nodePath string
	{
		components := strings.Split(ref, symSharp)
		if len(components) != 2 || len(components[1]) < 2 {
			log.Fatalf(`ref %s is not valid`, ref)
		}
		filePath = components[0]
		nodePath = components[1][1:]
	}
	filePath = path.Join(path.Dir(*flagFilepath), filePath)

	refData := readYml(filePath)
	var key string
	for _, key = range strings.Split(nodePath, `/`) {
		val, exists := refData[key]
		if !exists {
			log.Fatalf("ref %s does not contains key %s", filePath, key)
		}
		typed, casted := val.(anyMap[string])
		if !casted {
			log.Fatalf("ref %s does not contains map at key %s", filePath, key)
		}
		refData = typed
	}
	var dirname string
	{
		components := strings.Split(path.Dir(filePath), "/")
		if len(components) == 0 || components[len(components)-1] != keySchemas {
			log.Fatalf(`invalid directory name for ref %s; must contain "schemas"`, ref)
			// @todo think about responses and requestBodies
		}
		dirname = components[len(components)-1]
	}

	var rootRefNodePath string = strings.Join([]string{keyComponents, dirname, key}, "/")
	createPath[string](root, rootRefNodePath, refData)
	this.Set(keyRef, createLocalRefPath(rootRefNodePath))
	resolvedRefs[ref] = rootRefNodePath
	return
}

func createLocalRefPath(refPath string) string {
	return "#/" + refPath
}
