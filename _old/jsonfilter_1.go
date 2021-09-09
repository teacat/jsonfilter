package marika

import (
	"encoding/json"
	"log"

	"github.com/imdario/mergo"
	"github.com/teacat/marika/parser"
	// "github.com/imdario/mergo"
	// "github.com/alecthomas/repr"
)

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
	var vv map[string]interface{}
	for _, v := range results {
		err := mergo.Merge(&vv, v)
		if err != nil {
			panic(err)
		}
	}
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

func Resolve(src interface{}, prop *parser.Prop) interface{} {
	switch data := src.(type) {
	// {}
	case map[string]interface{}:
		if prop.Object.Prop != nil {
			return Object(prop.Object.Name, Resolve(data[prop.Object.Name], prop.Object.Prop))
		} else {
			return Object(prop.Object.Name, data[prop.Object.Name])
		}

	// [{}], [""]
	case []interface{}:
		var results []interface{}
		for _, j := range data {
			results = append(results, Resolve(j, prop))
		}
		return results
	default:
		return data

	}
}
