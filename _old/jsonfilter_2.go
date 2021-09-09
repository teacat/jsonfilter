package marika

import (
	"encoding/json"
	"log"

	"github.com/teacat/marika/parser"
	// "github.com/imdario/mergo"
	// "github.com/alecthomas/repr"
)

//
// a,b,c comma-separated list will select multiple fields
// a/b/c path will select a field from its parent
// a(b,c) sub-selection will select many fields from a parent
// a/*/c the star * wildcard will select all items in a field

func Filter(b []byte, filter string) ([]byte, error) {
	r, err := parser.Parse(filter)
	if err != nil {
		return []byte(``), err
	}
	var data map[string]interface{}
	err = json.Unmarshal(b, &data)
	if err != nil {
		return []byte(``), err
	}
	var results []interface{}
	for _, v := range r.Props {
		results = append(results, Resolve(data, v))
	}
	//newresult := Merge(results)
	/*var vv map[string]interface{}
	for _, v := range results {
		err := mergo.Merge(&vv, v)
		if err != nil {
			panic(err)
		}
	}*/
	a, _ := json.MarshalIndent(results, "", "    ")
	log.Println(string(a))

	return json.Marshal([]byte(`{}`))
}

func Object(k string, v interface{}) map[string]interface{} {
	return map[string]interface{}{
		k: v,
	}
}

func Debug(v interface{}) string {
	a, _ := json.MarshalIndent(v, "", "    ")
	return string(a)
}

func Merge(slice []interface{}) interface{} {
	results := make(map[string]interface{})
	for _, j := range slice {
		switch data := j.(type) {
		case map[string]interface{}:
			for k, v := range data {
				results[k] = v
			}
		case []interface{}:
			return data
		}
	}
	return results
}

func ResolveGroup(obj interface{}, mask *parser.Group) interface{} {

	switch data := obj.(type) {
	// {}
	case map[string]interface{}:
		var results []interface{}
		for _, v := range mask.Props {
			results = append(results, Resolve(data[mask.Name], v))
		}
		return Object(mask.Name, Merge(results))

	// []
	case []interface{}:
		var results []interface{}
		for _, v := range mask.Props {

			results = append(results, Resolve(data, v))
		}
		return Object(mask.Name, results)

	default:
		panic("panik")
	}

}

func Resolve(obj interface{}, mask *parser.Prop) interface{} {
	if mask.Object != nil {
		return ResolveObject(obj, mask)
	}
	return ResolveGroup(obj, mask.Group)
}

func ResolveObject(obj interface{}, mask *parser.Prop) interface{} {
	switch data := obj.(type) {
	// {}
	case map[string]interface{}:
		if mask.Object.Prop != nil {
			return Object(mask.Object.Name, Resolve(data[mask.Object.Name], mask.Object.Prop))
		} else {

			return Object(mask.Object.Name, data[mask.Object.Name])
		}

	// []
	case []interface{}:
		var values []interface{}
		for _, v := range data {
			switch j := v.(type) {
			case map[string]interface{}:
				values = append(values, Resolve(j[mask.Object.Name], mask))
			default:
				panic("yeee")
			}

		}
		return values

	}
	// string, int, float, bool
	return Object(mask.Object.Name, obj)
}
