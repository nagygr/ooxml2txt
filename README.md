# ooxml2txt

Reader library for the OOXML formats

## Example

Simple text retrieval from a `docx`, `pptx` and an `xlsx` file:

```go
package main

import (
	"fmt"
	"github.com/nagygr/ooxml2txt/format"
)

func main() {
	/* DOCX file: */
	doc, _ := format.MakeDocx("example.docx")
	fmt.Printf("%s\n", doc.Text)

	/* PPTX file: */
	ppt, _ := format.MakePptx("example.pptx")

	for n, slide := range ppt.Text {
		fmt.Printf("Slide %d: %s\n", n, slide)
	}

	/* XLSX file: */
	xls, _ := format.MakeXlsx("example.xlsx")

	for _, str := range xls.Text {
		fmt.Printf("%s\n", str)
	}
}
```
