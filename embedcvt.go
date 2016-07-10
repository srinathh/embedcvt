// Command embedcvt creates a Go program file that embeds the contents of the
// specified asset file(s) simple base64 constant that can be retrieved
// by calling a getter function unique to the asset. This is
// meant for simple data embedding use cases like inlining test data,
// inserting a few assets etc and does not do any compression, file
// system emulation, asset discovery etc. For more complex use cases,
// see https://github.com/jteeuwen/go-bindata.
//
// Usage: embedcvt [-p packageName] asset1 [asset2]...
package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

const usageString = `Usage: embedcvt [-p packageName] asset1 [asset2]...
embedcvt creates a Go program file that embeds the contents of the
specified asset file(s) as simple base64 constant that can be retrieved
by calling a getter function unique to the asset. This is 
meant for simple data embedding use cases like inlining test data,
inserting a few assets etc and does not do any compression, file
system emulation, asset discovery etc. For more complex use cases,
see https://github.com/jteeuwen/go-bindata.
Options:
`

var namifyRe = regexp.MustCompile("[A-Za-z0-9]+")

func namify(pathname string) string {
	texts := namifyRe.FindAllString(pathname, -1)

	if len(texts) == 0 {
		panic("could not create a name")
	}

	ret := strings.ToLower(texts[0])

	if len(texts) > 1 {
		for _, text := range texts[1:] {
			ret = ret + strings.ToUpper(string(text[0])) + strings.ToLower(text[1:])
		}
	}
	return ret
}

func main() {
	pkgName := flag.String("p", "main", "package name for the generated Go program file.")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usageString)
		flag.PrintDefaults()
	}

	flag.Parse()
	assetNames := flag.Args()
	if len(assetNames) == 0 {
		fmt.Fprintf(os.Stderr, "Error: need to specify at least one file to embed.\n")
		os.Exit(1)
	}

	buf := bytes.Buffer{}
	fmt.Fprintf(&buf, "package %s\n\n", *pkgName)
	fmt.Fprintf(&buf, "import \"encoding/base64\"\n\n")

	for _, assetName := range assetNames {
		contents, err := ioutil.ReadFile(assetName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: reading asset %s: %s\n", assetName, err)
			os.Exit(1)
		}

		nm := namify(assetName)

		fmt.Fprintf(&buf, "const %s = \"%s\"\n", nm, base64.StdEncoding.EncodeToString(contents))
		fmt.Fprintf(&buf, "func get%s() []byte{\n ret, _ := base64.StdEncoding.DecodeString(%s)\n return ret\n }", nm, nm)
		/*
				fmt.Fprintf(&buf, "var %s = []byte{", namify(assetName))
				for j, byt := range contents {
					fmt.Fprintf(&buf, " %d", byt)
					if j != len(contents)-1 {
						fmt.Fprint(&buf, ",")
					}
				}
			fmt.Fprint(&buf, "}\n\n")
		*/
	}
	ret, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: running gofmt on the generated source: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("%s", ret)
}
