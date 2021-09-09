package marika

import (
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

func TestParseObjects(t *testing.T) {
	var err error

	a := assert.New(t)
	_ = []byte(`{
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
	_, err = Filter(testset, "id, username, friends, information, locations")
	a.NoError(err)
	panic("")
}

func TestParseObjectsA(t *testing.T) {
	var err error

	a := assert.New(t)
	_ = []byte(`{
		"id"      : 1,
		"username": "YamiOdymel",
	}`)
	_, err = Filter(testset, "id, username")
	a.NoError(err)
	panic("")
}

func TestParseObjectsB(t *testing.T) {
	var err error
	a := assert.New(t)
	_ = []byte(`{
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
		}
	}`)
	_, err = Filter(testset, "friends, information")
	a.NoError(err)
	panic("")
}

func TestParseObjectsC(t *testing.T) {
	var err error
	a := assert.New(t)
	_ = []byte(`{
		"id"         : 1,
		"information": {
			"address": "R.O.C (Taiwan)",
			"phones" : ["12345", "67890"]
		}
	}`)
	_, err = Filter(testset, "id, information/address, information/phones")

	_ = []byte(`{
		"id"         : 1,
		"information": {
			"address": "R.O.C (Taiwan)",
			"phones" : ["12345", "67890"]
		}
	}`)
	_, err = Filter(testset, "id, information(address, phones)")
	a.NoError(err)
	panic("")
}

func TestParseObjectsD(t *testing.T) {
	var err error
	a := assert.New(t)
	_ = []byte(`{
		"information": {
			"skills" : [{
				"type": "eat"
			}, {
				"type": "drink"
			}]
		}
	}`)
	_, err = Filter(testset, "information/skills/type, information/skills/id")
	a.NoError(err)

	panic("")
}

func TestParseObjectsE(t *testing.T) {
	var err error
	a := assert.New(t)
	_ = []byte(`{
		"information": {
			"address": "R.O.C (Taiwan)",
			"skills" : [{
				"type": "eat"
			}, {
				"type": "drink"
			}]
		}
	}`)
	_, err = Filter(testset, "information(address, skills/type)")
	a.NoError(err)

	panic("")
}

func TestParseObjectsF(t *testing.T) {
	var err error
	a := assert.New(t)
	_ = []byte(`{
		"id"     : 1,
		"friends": [{
			"id"     : 2,
			"friends": [{
				"id": 3
			}]
		}, {
			"id"     : 4,
			"friends": [{
				"id": 5
			}]
		}]
	}`)
	_, err = Filter(testset, "id, friends/*/id")
	a.NoError(err)

	_ = []byte(`{
			"id"     : 1,
			"friends": [{
				"id"     : 2,
				"friends": [{
					"id": 3
				}]
			}, {
				"id"     : 4,
				"friends": [{
					"id": 5
				}]
			}]
		}`)
	_, err = Filter(testset, "id, friends(id, friends(id))")
	a.NoError(err)
	/*
		_ = []byte(`{
			"id"     : 1,
			"friends": [{
				"id"     : 2,
				"friends": [{
					"id": 3
				}]
			}, {
				"id"     : 4,
				"friends": [{
					"id": 5
				}]
			}]
		}`)
		_, err = Filter(testset, "id, friends(id, friends/id)")
		a.NoError(err)*/

	panic("")
}

func TestParseObjectsG(t *testing.T) {
	var err error
	a := assert.New(t)

	_ = []byte(`{
		"id"     : 1,
		"friends": [{
			"id"     : 2,
			"friends": [{
				"id"      : 3,
				"username": "Angrybird"
			}]
		}, {
			"id"     : 4,
			"friends": [{
				"id"      : 5,
				"username": "Max"
			}]
		}]
	}`)
	_, err = Filter(testset, "id, friends(id, friends)")
	a.NoError(err)

	panic("")
}

func TestParseObjectsH(t *testing.T) {
	var err error
	a := assert.New(t)

	_ = []byte(`{
		"friends": [{
			"friends": [{
				"id": 3
			}]
		}, {
			"friends": [{
				"id": 5
			}]
		}]
	}`)
	_, err = Filter(testset, "friends/*/id")

	_ = []byte(`{
		"friends": [{
			"friends": [{
				"id": 3
			}]
		}, {
			"friends": [{
				"id": 5
			}]
		}]
	}`)
	_, err = Filter(testset, "friends/*(id)")
	a.NoError(err)

	panic("")
}

func TestParseObjectsI(t *testing.T) {
	var err error
	a := assert.New(t)

	_ = []byte(`{
		"friends" : [{
			"id": 2
		}, {
			"id": 4
		}]
	}`)
	_, err = Filter(testset, "*/id")
	a.NoError(err)
}
