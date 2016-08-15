## Overview

A package to encode your structures in URL

## Installation

:warning: Make sure that Go is installed on your computer.

```shell
$ go get github.com/online-net/go-encodeUrl
```

Now, the package is ready to use.


## Examples

```go
import "github.com/online-net/go-encodeUrl"


type ID struct {
	Name        string `url:"name,ifStringIsNotEmpty"`
    //                        ^^        ^^^
    //              variable name  |  function
	DisplayName string `url:"display-name,ifStringIsNotEmpty"`
}


func main() {
	values, errs := encurl.Translate(&ID{"NotEmpty", ""})
	if errs != nil {
        fmt.Printf("errors %v", errs)
        return
	}
	fmt.Printf("https://example.com/?%v\n", values.Encode()) // https://example.com/?name=NotEmpty
}
```


## Functions

```console
ifStringIsNotEmpty
ifBoolIsFalse
ifBoolIsTrue
itoa
itoaIfNotNil
```

## Development

Feel free to contribute :smiley::beers:

## License

[MIT](https://github.com/online-net/go-encodeUrl/blob/master/LICENSE)
