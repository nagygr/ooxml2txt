# ooxml2txt

Reader library for the OOXML formats

## Example

Simple text retrieval from a `docx` file:

```go
package main

import (
	"fmt"
	"github.com/nagygr/ooxml2txt/format"
)

func main() {
	doc, _ := format.MakeDocx("example.docx")
	fmt.Printf("%s\n", doc.Text())
}
```
