# goxml2json [![CircleCI](https://circleci.com/gh/basgys/goxml2json.svg?style=svg)](https://circleci.com/gh/basgys/goxml2json)

Go package that converts XML to JSON

### Install

	go get -u github.com/pschlump/goxml2json

### Importing

	import "github.com/pschlump/goxml2json"

### Usage

**Code example**

```go
	package main

	import (
		"fmt"
		"os"
		"strings"

		"github.com/pschlump/goxml2json"
	)

	func main() {
		// xml is an io.Reader
		xml := strings.NewReader(`<?xml version="1.0" encoding="UTF-8"?><hello>world</hello>`)
		json, err := xml2json.Convert(xml)
		if err != nil {
			fmt.Fprintf ( os.Stderr, "Error processing XML file: %s\n", err )
			os.Exit(1)
		}

		fmt.Println(json.String())
		// {"hello": "world"}
	}

```

**Input**

```xml
	<?xml version="1.0" encoding="UTF-8"?>
	<osm version="0.6" generator="CGImap 0.0.2">
	 <bounds minlat="54.0889580" minlon="12.2487570" maxlat="54.0913900" maxlon="12.2524800"/>
	 <foo>bar</foo>
	</osm>
```

**Output**

```json
	{
		"osm": {
			"-version": "0.6",
			"-generator": "CGImap 0.0.2",
			"bounds": {
				"-minlat": "54.0889580",
				"-minlon": "12.2487570",
				"-maxlat": "54.0913900",
				"-maxlon": "12.2524800"
			},
			"foo": "bar"
		}
	}
```

### Contributing
Feel free to contribute to this project if you want to fix/extend/improve it.

### Contributors

	- [DirectX](https://github.com/directx)
	- [samuelhug](https://github.com/samuelhug)
	- [pschlump](https://github.com/pschlump)

### TODO

	 * Extract data types in JSON (numbers, boolean, ...)
	 * Categorise errors
	 * Benchmark
