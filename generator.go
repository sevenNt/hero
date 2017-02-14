package hero

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type genRes struct {
	useUnsafe bool
}

func writeToFile(path string, buffer *bytes.Buffer) {
	err := ioutil.WriteFile(path, buffer.Bytes(), os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func genAbsPath(path string) string {
	if !filepath.IsAbs(path) {
		path, _ = filepath.Abs(path)
	}
	return path
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func gen(n *node, buffer *bytes.Buffer) (res *genRes) {
	res = new(genRes)
	for _, child := range n.children {
		switch child.t {
		case TypeCode:
			buffer.Write(child.chunk.Bytes())
		case TypeHTML:
			buffer.WriteString(fmt.Sprintf(
				"_buffer.WriteString(`%s`)",
				child.chunk.String(),
			))
		case TypeRawValue, TypeEscapedValue:
			var format string

			if child.subtype == Bytes {
				if child.t == TypeRawValue {
					buffer.WriteString(fmt.Sprintf(
						"_buffer.Write(%s)",
						child.chunk.String(),
					))
					goto WriteBreakLine
				} else {
					format = "(*string)(*unsafe.Pointer(&%s))"
					res.useUnsafe = true
				}
			} else {
				switch child.subtype {
				case Int:
					format = "strconv.FormatInt(int64(%s), 10)"
				case Uint:
					format = "strconv.FormatUint(uint64(%s), 10)"
				case Float:
					format = "strconv.FormatFloat(float64(%s), 'f', -1, 64)"
				case Bool:
					format = "strconv.FormatBool(%s)"
				case String:
					format = "%s"
				case Interface:
					format = "fmt.Sprintf(\"%%v\", %s)"
				default:
					log.Fatal("unknown type")
				}
			}

			if child.t == TypeEscapedValue {
				format = fmt.Sprintf("html.EscapeString(%s)", format)
			}

			buffer.WriteString(fmt.Sprintf(
				fmt.Sprintf("_buffer.WriteString(%s)", format),
				child.chunk.String()),
			)
		case TypeBlock, TypeInclude:
			gen(child, buffer)
		default:
			continue
		}

	WriteBreakLine:
		buffer.WriteByte(BreakLine)
	}

	return
}

// Generate generates go code from source to test. pkgName represents the
// package name of the generated code.
func Generate(source, dest, pkgName string) {
	defer cleanGlobal()

	source, dest = genAbsPath(source), genAbsPath(dest)

	stat, err := os.Stat(source)
	checkError(err)

	fmt.Println("Parsing...")
	if stat.IsDir() {
		parseDir(source)
	} else {
		source, file := filepath.Split(source)
		parseFile(source, file)
	}

	stat, err = os.Stat(dest)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(dest, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	} else if !stat.IsDir() {
		log.Fatal(dest + " is not dir")
	} else if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Generating...")

	var wg sync.WaitGroup
	for path, n := range parsedNodes {
		wg.Add(1)

		fileName := filepath.Join(dest, fmt.Sprintf(
			"%s.go",
			strings.Join(strings.Split(path[len(source)+1:], "/"), "_"),
		))

		go func(n *node, source, fileName string) {
			defer wg.Done()

			buffer := bytes.NewBufferString(`
				// Code generated by hero.
			`)
			buffer.WriteString(fmt.Sprintf("// source: %s", source))
			buffer.WriteString(`
				// DO NOT EDIT!
			`)
			buffer.WriteString(fmt.Sprintf("package %s\n", pkgName))
			buffer.WriteString(`
				import "html"
				import "strconv"
				import "bytes"

				import "github.com/sevenNt/hero"
			`)

			imports := n.childrenByType(TypeImport)
			for _, item := range imports {
				buffer.Write(item.chunk.Bytes())
			}

			definitions := n.childrenByType(TypeDefinition)
			if len(definitions) == 0 {
				writeToFile(fileName, buffer)
				return
			}

			buffer.Write(definitions[0].chunk.Bytes())
			buffer.WriteString(`{
				_buffer := hero.GetBuffer()
			`)
			genRes := gen(n, buffer)
			buffer.WriteString(`
				return _buffer
			}`)

			if genRes.useUnsafe {
				buffer.WriteString(`
					import "unsafe"
				`)
			}
			writeToFile(fileName, buffer)
		}(n, path, fileName)
	}
	wg.Wait()

	fmt.Println("Executing goimports...")
	execCommand("goimports -w " + dest)

	fmt.Println("Executing go fmt...")
	execCommand("go fmt " + dest + "/*.go")

	fmt.Println("Executing go vet...")
	execCommand("go tool vet -v " + dest)
}
