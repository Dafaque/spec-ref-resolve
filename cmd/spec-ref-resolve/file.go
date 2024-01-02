package main

import (
	"log"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

func readYml(filePath string) anyMap[string] {
	file, errOpenFile := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
	if errOpenFile != nil {
		log.Fatal(errOpenFile)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	var m anyMap[string]
	if err := decoder.Decode(&m); err != nil {
		log.Fatal(err)
	}
	return m
}
func writeYaml(data *anyMap[string]) {
	if err := os.MkdirAll(path.Dir(*flagOutFilepath), os.ModePerm); err != nil {
		log.Fatal(err)
	}
	f, err := os.OpenFile(*flagOutFilepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	encoder := yaml.NewEncoder(f)
	encoder.SetIndent(2)
	defer encoder.Close()
	if err := encoder.Encode(data); err != nil {
		log.Fatal(err)
	}
}
