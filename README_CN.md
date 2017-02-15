# Hero
(访问源分支请前往[shiyanhui/hero](github.com/shiyanhui/hero))

Hero是一个高性能、强大并且易用的go模板引擎，工作原理是把模板预编译为go代码。Hero目前已经在[bthub.io](http://bthub.io)的线上环境上使用。

[![GoDoc](https://godoc.org/github.com/shiyanhui/hero?status.svg)](https://godoc.org/github.com/shiyanhui/hero)
[![Go Report Card](https://goreportcard.com/badge/github.com/shiyanhui/hero)](https://goreportcard.com/report/github.com/shiyanhui/hero)

- [Features](#features)
- [Install](#install)
- [Usage](#usage)
- [Quick Start](#quick-start)
- [Template Syntax](#template-syntax)
- [License](#license)

## Features

- 非常易用.
- 功能强大，支持模板继承和模板include.
- 高性能.
- 自动编译.

## Install

    go get github.com/sevenNt/hero
    go install github.com/sevenNt/hero/hero

## Usage

```shell
hero [options]

options:
	- s: 模板目录，默认为当前目录
	- d: 生成的go代码的目录，如果没有设置的话，和source一样
	- p: 生成的go代码包的名称，默认为template
	- w: 是否监控模板文件改动并自动编译
	- e: 指定编译源文件后缀，默认为dtpl

example:
	hero -source="./"
	hero -source="$GOPATH/src/app/template" -w
```

## Quick Start

假设我们现在要渲染一个用户列表模板`userlist.dtpl`, 它继承自`index.dtpl`, 并且一个用户的模板是`user.dtpl`. 我们还假设所有的模板都在`$GOPATH/src/app/template`目录下。

### index.dtpl

```html
<!DOCTYPE html>
<html>
    <head>
        <meta charset="utf-8">
    </head>

    <body>
        <%@ body { %>
        <% } %>
    </body>
</html>
```

### users.dtpl

```html
<%: func UserList(userList []string) []byte %>

<%~ "index.dtpl" %>

<%@ body { %>
    <% for _, user := range userList { %>
        <ul>
            <%+ "user.dtpl" %>
        </ul>
    <% } %>
<% } %>
```

### user.dtpl

```html
<li>
    <%= user %>
</li>
```

然后我们编译这些模板:

```shell
hero -source="$GOPATH/src/app/template"
```

编译后，我们将在同一个目录下得到三个go文件，分别是`index.dtpl.go`, `user.dtpl.go` and `userlist.dtpl.go`, 然后我们在http server里边去调用模板：

### main.go

```go
package main

import (
	"app/template"
	"net/http"
)

func main() {
	http.HandleFunc("/users", func(w http.ResponseWriter, req *http.Request) {
		var userList = []string {
          	"Alice",
			"Bob",
			"Tom",
		}
		w.Write(template.UserList(userList))
	})

	http.ListenAndServe(":8080", nil)
}
```

最后，运行这个http server，访问`http://localhost:8080/users`，我们就能得到我们期待的结果了！

## Template syntax

Hero总共有九种语句，他们分别是：

- 函数定义语句 `<%: func define %>`
  - 该语句定义了该模板所对应的函数，如果一个模板中没有函数定义语句，那么最终结果不会生成对应的函数。
  - 该函数必须返回一个`[]byte`参数。
  - 例:`<%: func UserList(userList []string) []byte %>`

- 模板继承语句 `<%~ "parent template" %>`
  - 该语句声明要继承的模板。
  - 例: `<%~ "index.dtpl" >`

- 模板include语句 `<%+ "sub template" %>`
  - 该语句把要include的模板加载进该模板，工作原理和`C++`中的`#include`有点类似。
  - 例: `<%+ "user.dtpl" >`

- 包导入语句 `<%! go code %>`
  - 该语句用来声明所有在函数外的代码，包括依赖包导入、全局变量、const等。

  - 该语句不会被子模板所继承

  - 例:

    ```go
    <%!
    	import (
          	"fmt"
        	"strings"
        )

    	var a int

    	const b = "hello, world"

    	func Add(a, b int) int {
        	return a + b
    	}

    	type S struct {
        	Name string
    	}

    	func (s S) String() string {
        	return s.Name
    	}
    %>
    ```

- 块语句 `<%@ blockName { %> <% } %>`

  - 块语句是用来在子模板中重写父模中的同名块，进而实现模板的继承。

  - 例:

    ```html
    <!DOCTYPE html>
    <html>
        <head>
            <meta charset="utf-8">
        </head>

        <body>
            <%@ body { %>
            <% } %>
        </body>
    </html>
    ```

- Go代码语句 `<% go code %>`

  - 该语句定义了函数内部的代码部分。

  - 例:

    ```go
    <% for _, user := userList { %>
        <% if user != "Alice" { %>
        	<%= user %>
        <% } %>
    <% } %>

    <%
    	a, b := 1, 2
    	c := Add(a, b)
    %>
    ```

- 原生值语句 `<%== statement %>`

  - 该语句把变量转换为string。

  - 例:

    ```go
    <%== a %>
    <%== a + b %>
    <%== Add(a, b) %>
    <%== user.Name %>
    ```

- 转义值语句 `<%= statement %>`

  - 该语句把变量转换为string后，又通过`html.EscapesString`记性转义。

  - 例:

    ```go
    <%= a %>
    <%= a + b %>
    <%= Add(a, b) %>
    <%= user.Name %>
    ```

- 注释语句 `<%# note %>`

  - 该语句注释相关模板，注释不会被生成到go代码里边去。
  - 例: `<# 这是一个注释 >`.

## License

Hero is licensed under the Apache License.
