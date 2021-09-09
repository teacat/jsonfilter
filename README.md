# JSON Filter [台灣正體](./README-tw.md) [![GoDoc](https://godoc.org/github.com/teacat/jsonfilter?status.svg)](https://godoc.org/github.com/teacat/jsonfilter) [![Coverage Status](https://coveralls.io/repos/github/teacat/jsonfilter/badge.svg?branch=master)](https://coveralls.io/github/teacat/jsonfilter?branch=master) [![Build Status](https://travis-ci.org/teacat/jsonfilter.svg?branch=master)](https://travis-ci.com/teacat/jsonfilter) [![Go Report Card](https://goreportcard.com/badge/github.com/teacat/jsonfilter)](https://goreportcard.com/report/github.com/teacat/jsonfilter)

JSON Filter is a GraphQL-like but for REST API or common usage to filter the JSON data. With JSON Filter, you are able to filter the JSON with `users(username,nickname)` like syntax.

## Features

-   Filter the JSON data.
-   Easy to use.

## Why?

GraphQL solves the problem but it takes time to design a perfect schema/system. While REST is the most quickly and easily to implement, we might still wanted to get the same benifit as GraphQL as possible.

The REST API has it's own problems, like: Over-fetching, Under-fetching. You might get too many unnecessary data, or you might get too few data.

Google solves the problem with thier own APIs by providing [Partial Resources](https://cloud.google.com/compute/docs/api/how-tos/performance#partial) usage which you can filter the response with a `fields=kind,items(title,characteristics/length)` param.

Adobe does almost [the same thing](https://devdocs.magento.com/guides/v2.4/rest/retrieve-filtered-responses.html) with a bit different syntax: `?fields=items[name,qty,sku]`.

## Note

The current implementation is based on [nemtsov/json-mask](https://github.com/nemtsov/json-mask) with JavaScript interpreter (powered by [robertkrimen/otto](https://github.com/robertkrimen/otto)) so the performance might needs to be considered.

The full-Golang solution is WIP and won't be implmeneted recently.

## Installation

Install the package via `go get` command.

```bash
$ go get github.com/teacat/jsonfilter
```

## Usage

A complete example:

```go
package main

import (
	"fmt"
	"github.com/teacat/jsonfilter"
)

type User struct {
	ID       int      `json:"id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Phones   []string `json:"phones"`
}

func main() {
	data := []byte(`
{
    "id"      : 1,
    "username": "YamiOdymel",
    "email"   : "foobar@gmail.com",
    "phones"  : ["0912345678", "0987654321"]
}`)

	b, err := jsonfilter.Filter(data, "username,phones")
	if err != nil {
		panic(err)
	}

	var data MyData
	if err := json.Unmarshal(b, &data); err != nil {
		panic(err)
	}

	fmt.Printf("%+v", data) // Output: {ID:0 Username:YamiOdymel Email: Phones:[0912345678 0987654321]}
}
```

## Syntax

The syntax is loosely based on XPath:

-   `a,b,c` Comma-separated list will select multiple fields.
-   `a/b/c` Path will select a field from its parent.
-   `a(b,c)` Sub-selection will select many fields from a parent.
-   `a/*/c` The star \* wildcard will select all items in a field.

While trying to select the properties from the same parent,

use `a(b,c)` syntax instead of `a/b,a/c` because the `a/c` will override the results of `a/b` since the result won't be collected.

Also with wildcard syntax, you might get undesired results and this is an unknown issue from [nemtsov
/
json-mask](https://github.com/nemtsov/json-mask): [\* returns extra falsey values at the level of the wildcard #11](https://github.com/nemtsov/json-mask/issues/11).

## Examples

The following JSON structure as the data source.

```json
{
    "id": 1,
    "username": "YamiOdymel",
    "friends": [
        {
            "id": 2,
            "username": "Yan-K"
        }
    ],
    "information": {
        "address": "R.O.C (Taiwan)",
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
}
```

---

```json
"id,username"

{
    "id": 1,
    "username": "YamiOdymel"
}
```

---

```json
"friends(id,username)"

{
    "friends": [
        {
            "id": 2,
            "username": "Yan-K"
        }
    ]
}
```

---

```json
"information/skills(type)"
"information/skills/type"
"information/*/type"

{
    "information": {
        "skills": [
            {
                "type": "eat"
            },
            {
                "type": "drink"
            }
        ]
    }
}
```
