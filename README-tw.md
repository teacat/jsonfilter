# JSON Filter [台灣正體](./README-tw.md) [![GoDoc](https://godoc.org/github.com/teacat/jsonfilter?status.svg)](https://godoc.org/github.com/teacat/jsonfilter) [![Coverage Status](https://coveralls.io/repos/github/teacat/jsonfilter/badge.svg?branch=master)](https://coveralls.io/github/teacat/jsonfilter?branch=master) [![Build Status](https://travis-ci.org/teacat/jsonfilter.svg?branch=master)](https://travis-ci.com/teacat/jsonfilter) [![Go Report Card](https://goreportcard.com/badge/github.com/teacat/jsonfilter)](https://goreportcard.com/report/github.com/teacat/jsonfilter)

JSON Filter 是一個為 REST API 或一般用途所設計的類 GraphQL JSON 資料過濾篩選函式庫。你能夠透過 `users(username,nickname)` 這種用法篩選 JSON 資料。

## 特色

-   過濾 JSON 資料。
-   簡單易用。

## 為什麼？

GraphQL 雖然解決的資料撈取的問題，但你需要花費許多心思在設計完美的結構、系統上。而 REST 設計方式總是最簡單、快速的。即便如此，我們還是會想要盡可能地擁有類似 GraphQL 那樣的特點。

在這裡不得不提 REST API 會有一些問題：過度撈取、撈取不足。例如：你可能只需要一點資料，但單個 API 卻回傳過多的複雜結構；又或者你可能需要呼叫並彙整多個 API 才能得到你想要的資料。

Google 提出 [Partial Resources](https://cloud.google.com/compute/docs/api/how-tos/performance#partial) 解決了這部份的問題，呼叫 Google API 時可以傳入 `fields=kind,items(title,characteristics/length)` 參數過濾回應資料。

而 Adobe 也做出了[類似的舉動](https://devdocs.magento.com/guides/v2.4/rest/retrieve-filtered-responses.html)但用法有點不同：`?fields=items[name,qty,sku]`。

## 備註

目前的實作方式是基於 [nemtsov/json-mask](https://github.com/nemtsov/json-mask) 並以 JavaScript 直譯器（由 [robertkrimen/otto](https://github.com/robertkrimen/otto) 驅動）執行，所以效能上可能需要多加顧慮。

完整的 Golang 解決方案尚未完成且不會於近日推出。

## 安裝方式

透過 `go get` 指令安裝套件。

```bash
$ go get github.com/teacat/jsonfilter
```

## 使用方式

一個完整的範例如下：

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

	fmt.Printf("%+v", data) // 輸出：{ID:0 Username:YamiOdymel Email: Phones:[0912345678 0987654321]}
}
```

## 語法

語法某方面近似於 XPath：

-   `a,b,c` 透過逗號分隔能夠選擇不同的欄位。
-   `a/b/c` 以路徑的方式下去選取子欄位。
-   `a(b,c)` 在某個欄位中批次選取其子欄位。
-   `a/*/c` 萬用符號 \* 會選擇該階層中的所有欄位。

位於同個階層欲選擇不同欄位時，請使用 `a(b,c)` 而不是 `a/b,a/c`。因為 `a/c` 會將 `a/b` 的選取結果覆蓋而不是彙整。

而目前萬用符號可能會取得到非預期的結果，這是目前已知的問題：[nemtsov
/
json-mask](https://github.com/nemtsov/json-mask): [\* returns extra falsey values at the level of the wildcard #11](https://github.com/nemtsov/json-mask/issues/11).

## 範例

接下來的範例都基於這個 JSON 資料結構。

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
