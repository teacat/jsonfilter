package jsonfilter

import (
	"github.com/robertkrimen/otto"
)

// Filter filters the object by the specified masks.
// The syntax is loosely based on XPath:
//
//  a,b,c comma-separated list will select multiple fields
//  a/b/c path will select a field from its parent
//  a(b,c) sub-selection will select many fields from a parent
//  a/*/c the star * wildcard will select all items in a field
func Filter(data []byte, mask string) ([]byte, error) {
	if mask == "" {
		return data, nil
	}

	vm := otto.New()
	if err := vm.Set("inputObj", string(data)); err != nil {
		return nil, err
	}
	if err := vm.Set("inputMask", mask); err != nil {
		return nil, err
	}
	if _, err := vm.Run(js); err != nil {
		return nil, err
	}
	value, err := vm.Get("result")
	if err != nil {
		return nil, err
	}
	// inter, err := value.Export()
	// if err != nil {
	// 	return nil, err
	// }
	return []byte(value.String()), nil
}
