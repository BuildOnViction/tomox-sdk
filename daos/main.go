package daos

import "github.com/globalsign/mgo"

func existedIndex(indexName string, indexes []mgo.Index) bool {
	for _, index := range indexes {
		if index.Name == indexName {
			return true
		}
	}
	return false
}
