package util

import (
	"errors"
	"github.com/smallfish/simpleyaml"
  "strconv"
  "bytes"
)

var ArrayOfPaths = make([]string, 0)

func GetAllExistingPaths(y *simpleyaml.Yaml, PathSlice []string) ([]string, error) {
	if y.IsMap() {
		keys, err := y.GetMapKeys()
		if err != nil {
			return nil, errors.New("Retrieving map keys failed")
		}
		for k, _ := range keys {
			if k != 0 {
				PathSlice = PathSlice[:len(PathSlice)-1]
			}
			PathSlice = append(PathSlice, keys[k])
			GetAllExistingPaths(y.Get(keys[k]), PathSlice)
		}
	} else if y.IsArray() {
		arr, err := y.Array()
		if err != nil {
			return nil, errors.New("Retrieving array failed")
		}
		for k, _ := range arr {
			if k != 0 {
				PathSlice = PathSlice[:len(PathSlice)-1]
			}
			PathSlice = append(PathSlice, strconv.Itoa(k))
			GetAllExistingPaths(y.GetIndex(k), PathSlice)
		}
	} else {
		var buffer bytes.Buffer
		for k, _ := range PathSlice {
				if k == len(PathSlice)-1 {
				  buffer.WriteString(PathSlice[k])
				}else{
				  buffer.WriteString(PathSlice[k]+"/")
				}
		}
		ArrayOfPaths = append(ArrayOfPaths, buffer.String())
	}
	return ArrayOfPaths, nil
}

// GetAllPaths retrieves all possible paths in the YAML file
//
// Example:
//      util.GetAllPaths(*Yaml)
func GetAllPaths(y *simpleyaml.Yaml) ([]string, error)  {
	InitialPath := make([]string, 0)
	AllPaths, err :=  GetAllExistingPaths(y, InitialPath)
	if err != nil {
		return nil, errors.New("Retrieving paths failed")
	}
	return AllPaths, nil
}
