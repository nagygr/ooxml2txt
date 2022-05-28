# ooxml2txt

Reader library for the OOXML (Office Open XML) formats (i.e. `docx`, `pptx` and
`xlsx`). The library extracts the textual content from documents thus making it
possible to perform simple text search on them.

## Example

Simple text retrieval from a `docx`, `pptx` and an `xlsx` file:

```go
package main

import (
	"fmt"
	"github.com/nagygr/ooxml2txt/pkg/format"
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

	/* DOCX file from URL: */
	docUrl, _ := format.MakeDocxFromUrl(
		"https://github.com/nagygr/ooxml2txt/raw/main/test_data/example.docx",
	)

	for _, str := range docUrl.Text {
		fmt.Printf("%s\n", str)
	}
}
```

## Formats

There's a dedicated type for each document format the library recognizes. They
all reside in the `format` package and are called: `Docx`, `Pptx` and `Xlsx`
respectively.

They reflect the structure of the given format and thus each needs to be
handled in a different way.

What's common in them is that they are bare structs without methods: once they
are created successfully, they contain valid information in their data members.

If something goes wrong (the given document path doesn't exist, the document's
structure doesn't conform to the format recognized by the library, etc.) then
an `error` is returned. Although errors are not handled in the examples above,
they should always be handled in real life applications.

Each format handler can be instantiated for a local file and also for a URL. In
the latter case, the document is loaded directly into memory without the need
to save it to the filesystem first. The functions creating the format handler
from a URL end with *"FromUrl"*.

### Docx

`Docx` represents text documents. It has the following public members:

```go
type Docx struct {
	Text      string
	Links     []string
	Footnotes []string
	Headers   []string
	Footers   []string
	// ...
}
```

-	`Text`: contains the document text
-	`Links`: contains the links within the document (`Text` contains references
	to the links)
-	`Footnotes`: contains the footnotes of the document
-	`Headers`: contains the headers of the document
-	`Footers`: contains the footers of the document

### Pptx

`Pptx` represents presentations. It has the following public member:

```go
type Pptx struct {
	Text      []string
	// ...
}
```

The `Text` slice contains the text of the slides, each as a separate string.

### Xlsx

`Xlsx` represents spreadsheet documents. It is a bit special in that it doesn't
contain everything from the spreadsheet. As this library is targeted towards
text search in OOXML documents, gathering the formulas and numeric values would
be of little benefit. The `xlsx` format contains the unique strings from every
worksheet in a single file and this library only reads that file. As a result,
only strings containing text are returned (no numbers, dates, etc.) and each
text fragment is only returned once no matter how many times it appears in the
document.

The `Xlsx` struct has one public member:

```go
type Xlsx struct {
	Text      []string
	// ...
}
```

The `Text` slice contains the unique strings from the given document.
