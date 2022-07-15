package utils

import (
	"log"
	"sort"
)

func CheckError(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func SortedMapKeys(stringsMap map[string]string) []string {
	keys := make([]string, 0, len(stringsMap))
	for k := range stringsMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
