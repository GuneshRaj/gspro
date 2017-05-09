# gspro
Go Server Pages Pre-processor

Go Server Pages - pre processor v0.1 (c) gunesh.raj@gmail.com
Converts *.gun files to Go files.
  init <outputDir>
    Creates skeleton code in a folder
  process <searchDir> <prefixDir> <outputDir>
    Processes *.gun files recursively & creates *.go files in outputDir

This project simplifies web development similar to JSP and PHP. The Templates are compiled into the Executable and technically performs better than the typical parsed templates. The result is a web server with route maps with GET and POST methods to the Template.

It uses the Echo Framework - https://github.com/labstack/echo
(Install Echo with $ go get -u github.com/labstack/echo)

The project is in its infancy stage. I am planing to add support for better Scriplet, Expressions and Directives sans the Expression language.
The second stage would be supporting seamless auto project building eg rebuilding on every change.

Im open for feedback and contribution.


Getting Started

1. Build the code. go build gspro.go, it creates an executable eg gspro.exe
2. Create a new project folder eg \demo\testproject\
3. Copy gspro.exe to the project folder or add the executable to the system path
4. Run gspro.exe init while in the project folder

Creating your Go Server Pages

1. Create a file with an extension *.gun in the project folder or any child folders within the project root folder
2. Use the Code <% and %> to embed Go Lang codes.
3. The codes are using Echo Framework, echo.Echo is represented as c
4. Use WriteString(s) and WriteInt(i) for writing to the buffer

Building your project

1. In the Project Root folder, run gspro.exe process <searchDir> <prefixDir> <outputDir>
2. Search Dir - Usually the Root Project folder that contains the Go Server Pages, GSPRO will recursively get all the files with the *.gun files
3. Prefix Dir - To simplify the filename and generated route methods, you could have the Prefix Dir and Search Dir as the same
4. Output Dir - The output of the files thats generated.
5. Generates a single Go file for each *.Gun file
6. Also generates routemap.go to map all the *.Gun files
7. Route Maps are generally using the filename, eg av.gun is mapped to /av and uses the method func av_gun(c echo.Context) error
8. Package name is main

Example

[av.gun]
```
<%
    s := "World!"

%>
Hello <% WriteStr(c, s) %>

<%
i := 4
%>
You are <% WriteInt(c, i) %> billion years old.
```

Compiles to

[av_gun.go]
```
package main

import (
	"net/http"
	"github.com/labstack/echo"
)


func av(c echo.Context) error {

    s := "World"


WriteStr(c, `
Hello `)
 WriteStr(c, s) 
WriteStr(c, `

`)

i := 4

WriteStr(c, `
I am `)
 WriteInt(c, i) 
WriteStr(c, ` billion years old.

`)


	c.Response().Flush()
	c.Response().WriteHeader(http.StatusOK)
	return nil
}
```

[routemap.go]
```
package main

import (
	"github.com/labstack/echo"
	"strconv"
)

func routeMap(e *echo.Echo) {
  e.GET("/av", av)
  e.POST("/av", av)
}


func WriteStr(c echo.Context, s string) {
	c.Response().Writer.Write([]byte(s))
}

func WriteInt(c echo.Context, s int) {
	c.Response().Writer.Write([]byte(strconv.Itoa(s)))
}

```

