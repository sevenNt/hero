package hero

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

var replacer *regexp.Regexp

func init() {
	replacer, _ = regexp.Compile("\\s")
}

func TestWriteToFile(t *testing.T) {
	path := "/tmp/hero.test"
	content := "hello, hero"

	buffer := bytes.NewBufferString(content)
	writeToFile(path, buffer)

	defer os.Remove(path)

	if _, err := os.Stat(path); err != nil {
		t.Fail()
	}

	if c, err := ioutil.ReadFile(path); err != nil || string(c) != content {
		t.Fail()
	}
}

func TestGenAbsPath(t *testing.T) {
	dir, _ := filepath.Abs("./")

	parts := strings.Split(dir, "/")
	parent := strings.Join(parts[:len(parts)-1], "/")

	cases := []struct {
		in  string
		out string
	}{
		{in: "/", out: "/"},
		{in: ".", out: dir},
		{in: "../", out: parent},
	}

	for _, c := range cases {
		if genAbsPath(c.in) != c.out {
			t.Fail()
		}
	}
}

func TestGenerate(t *testing.T) {
	Generate(rootDir, rootDir, "template", "dtpl")

	cases := []struct {
		file string
		code string
	}{
		{file: "index.html.go", code: `
// Code generated by hero.
// source: /tmp/gohero/index.html
// DO NOT EDIT!
package template
`},
		{file: "item.html.go", code: `
// Code generated by hero.
// source: /tmp/gohero/item.html
// DO NOT EDIT!
package template
`},
		{file: "list.html.go", code: `
// Code generated by hero.
// source: /tmp/gohero/list.html
// DO NOT EDIT!
package template

import (
	"bytes"
	"html"

	"github.com/sevenNt/hero"
)

func Add(a, b int) int {
	return a + b
}

func UserList(userList []string) *bytes.Buffer {
	_buffer := hero.GetBuffer()
	_buffer.WriteString(` + "`" + `
<!DOCTYPE html>
<html>
  <head>
  </head>
  <body>
    ` + "`" + `)
	for _, user := range userList {
		_buffer.WriteString(` + "`" + `
<div>
    <a href="/user/` + "`" + `)
		_buffer.WriteString(html.EscapeString(user))
		_buffer.WriteString(` + "`" + `">
        ` + "`" + `)
		_buffer.WriteString(user)
		_buffer.WriteString(` + "`" + `
    </a>
</div>
` + "`" + `)

	}

	_buffer.WriteString(` + "`" + `
  </body>
</html>
` + "`" + `)

	return _buffer
}
		`},
	}

	for _, c := range cases {
		content, err := ioutil.ReadFile(filepath.Join(rootDir, c.file))
		if err != nil || !reflect.DeepEqual(
			replacer.ReplaceAll(content, nil),
			[]byte(replacer.ReplaceAllString(c.code, "")),
		) {
			t.Fail()
		}
	}
}

func TestGen(t *testing.T) {
	root := parseFile(rootDir, "list.html")
	buffer := new(bytes.Buffer)

	gen(root, buffer)

	if buffer.String() == replacer.ReplaceAllString(
		`for _, user := range userList {}`, "") {
		t.Fail()
	}
}
