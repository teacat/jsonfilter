package jsonfilter

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testset = []byte(`{
	"id"      : 1,
	"username": "YamiOdymel",
	"friends" : [{
		"id"      : 2,
		"username": "Yan-K",
		"friends" : [{
			"id"      : 3,
			"username": "Angrybird"
		}]
	}, {
		"id"      : 4,
		"username": "Noodle",
		"friends" : [{
			"id"      : 5,
			"username": "Max"
		}]
	}],
	"information": {
		"address": "R.O.C (Taiwan)",
		"phones" : ["12345", "67890"],
		"skills" : [{
			"id"  : 1,
			"type": "eat"
		}, {
			"id"  : 2,
			"type": "drink"
		}]
	},
	"locations": [
		[100, 200, {"name": "Taiwan"}],
		[300, 400, {"name": "Hong Kong"}]
	]
}`)

func Check(a *assert.Assertions, expected string, data []byte) {
	var v interface{}
	err := json.Unmarshal(data, &v)
	a.NoError(err)

	b, err := json.MarshalIndent(v, "", "    ")
	a.NoError(err)

	a.Equal(expected, string(b))
}

func TestObjects(t *testing.T) {
	a := assert.New(t)
	expected := `{
    "friends": [
        {
            "friends": [
                {
                    "id": 3,
                    "username": "Angrybird"
                }
            ],
            "id": 2,
            "username": "Yan-K"
        },
        {
            "friends": [
                {
                    "id": 5,
                    "username": "Max"
                }
            ],
            "id": 4,
            "username": "Noodle"
        }
    ],
    "id": 1,
    "information": {
        "address": "R.O.C (Taiwan)",
        "phones": [
            "12345",
            "67890"
        ],
        "skills": [
            {
                "id": 1,
                "type": "eat"
            },
            {
                "id": 2,
                "type": "drink"
            }
        ]
    },
    "locations": [
        [
            100,
            200,
            {
                "name": "Taiwan"
            }
        ],
        [
            300,
            400,
            {
                "name": "Hong Kong"
            }
        ]
    ],
    "username": "YamiOdymel"
}`
	data, err := Filter(testset, "id,username,friends,information,locations")
	a.NoError(err)
	Check(a, expected, data)
}

func TestFields(t *testing.T) {
	a := assert.New(t)
	expected := `{
    "id": 1,
    "username": "YamiOdymel"
}`
	data, err := Filter(testset, "id,username")
	a.NoError(err)
	Check(a, expected, data)
}

func TestDeepFields(t *testing.T) {
	a := assert.New(t)
	expected := `{
    "friends": [
        {
            "friends": [
                {
                    "id": 3,
                    "username": "Angrybird"
                }
            ],
            "id": 2,
            "username": "Yan-K"
        },
        {
            "friends": [
                {
                    "id": 5,
                    "username": "Max"
                }
            ],
            "id": 4,
            "username": "Noodle"
        }
    ],
    "information": {
        "address": "R.O.C (Taiwan)",
        "phones": [
            "12345",
            "67890"
        ],
        "skills": [
            {
                "id": 1,
                "type": "eat"
            },
            {
                "id": 2,
                "type": "drink"
            }
        ]
    }
}`
	data, err := Filter(testset, "friends,information")
	a.NoError(err)
	Check(a, expected, data)
}

func TestGroup(t *testing.T) {
	a := assert.New(t)
	expected := `{
    "id": 1,
    "information": {
        "address": "R.O.C (Taiwan)",
        "phones": [
            "12345",
            "67890"
        ]
    }
}`
	// NOTE: `information/phones` overrides `information/address`
	// To get both of the fields, use `information(address,phones)`
	//
	// data, err := Filter(testset, "id,information/address,information/phones")
	// a.NoError(err)
	//
	// b, err := json.MarshalIndent(string(data), "", "    ")
	// a.NoError(err)
	//
	// a.Equal(expected, string(b))

	data, err := Filter(testset, "id,information(address,phones)")
	a.NoError(err)
	Check(a, expected, data)
}

func TestDeepGroup(t *testing.T) {
	a := assert.New(t)
	expected := `{
    "information": {
        "skills": [
            {
                "id": 1,
                "type": "eat"
            },
            {
                "id": 2,
                "type": "drink"
            }
        ]
    }
}`
	// NOTE: `information/skills/id` overrides `information/skills/type`
	// To get both of the fields, use `information/skills(type,id)`
	//
	// data, err := Filter(testset, "information/skills/type,information/skills/id")
	// a.NoError(err)
	//
	// b, err := json.MarshalIndent(string(data), "", "    ")
	// a.NoError(err)
	//
	// a.Equal(expected, string(b))

	data, err := Filter(testset, "information/skills(type,id)")
	a.NoError(err)
	Check(a, expected, data)
}

func TestGroupWithNestedObject(t *testing.T) {
	a := assert.New(t)
	expected := `{
    "information": {
        "address": "R.O.C (Taiwan)",
        "skills": [
            {
                "type": "eat"
            },
            {
                "type": "drink"
            }
        ]
    }
}`
	data, err := Filter(testset, "information(address,skills/type)")
	a.NoError(err)
	Check(a, expected, data)
}

func TestNestedGroup(t *testing.T) {
	a := assert.New(t)
	expected := `{
    "friends": [
        {
            "friends": [
                {
                    "id": 3
                }
            ],
            "id": 2
        },
        {
            "friends": [
                {
                    "id": 5
                }
            ],
            "id": 4
        }
    ],
    "id": 1
}`
	// NOTE: `friends/friends/id` overrides `friends/id`
	// To get both of the fields, use `friends(id,friends(id))`
	//
	// data, err := Filter(testset, "id,friends/id,friends/friends/id")
	// a.NoError(err)
	//
	// b, err := json.MarshalIndent(string(data), "", "    ")
	// a.NoError(err)
	//
	// a.Equal(expected, string(b))

	data, err := Filter(testset, "id,friends(id,friends(id))")
	a.NoError(err)
	Check(a, expected, data)

	data, err = Filter(testset, "id,friends(id,friends/id)")
	a.NoError(err)
	Check(a, expected, data)
}

func TestGroupWithArray(t *testing.T) {
	a := assert.New(t)
	expected := `{
    "friends": [
        {
            "friends": [
                {
                    "id": 3,
                    "username": "Angrybird"
                }
            ],
            "id": 2
        },
        {
            "friends": [
                {
                    "id": 5,
                    "username": "Max"
                }
            ],
            "id": 4
        }
    ],
    "id": 1
}`
	data, err := Filter(testset, "id,friends(id,friends)")
	a.NoError(err)
	Check(a, expected, data)
}

func TestDeepWildcard(t *testing.T) {
	a := assert.New(t)
	expected := `{
    "friends": [
        {
            "friends": [
                {
                    "id": 3
                }
            ]
        },
        {
            "friends": [
                {
                    "id": 5
                }
            ]
        }
    ]
}`
	data, err := Filter(testset, "friends/*/id")
	a.NoError(err)
	Check(a, expected, data)

	data, err = Filter(testset, "friends/*(id)")
	a.NoError(err)
	Check(a, expected, data)
}

func TestWildcardRoot(t *testing.T) {
	a := assert.New(t)
	// NOTE: Check https://github.com/nemtsov/json-mask/issues/11
	expected := `{
    "friends": [
        {
            "id": 2
        },
        {
            "id": 4
        }
    ],
    "information": {},
    "locations": [
        [],
        []
    ]
}`
	data, err := Filter(testset, "*/id")
	a.NoError(err)
	Check(a, expected, data)
}
