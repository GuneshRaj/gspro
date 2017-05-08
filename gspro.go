/*
(c) Gunesh Raj - gunesh.raj@gmail.com
version 0.1
 */
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"io/ioutil"
	"regexp"
	"bytes"
)

func main() {

	searchDir := "" //"/Users/ice/GoglandProjects/gtest1/"
	prefixDir := "" //"\\Users\\ice\\GoglandProjects\\"
	outputDir := "" //"/Users/ice/GoglandProjects/gtest1/"

	argsWithProg := os.Args

	if len(argsWithProg) > 1 {
		arg2 := argsWithProg[1]

		if arg2 == "init" {
			if len(argsWithProg) == 3 {
				outputDir = argsWithProg[2]
				createInitial(outputDir)
			} else {
				// error
				fmt.Println("init needs 1 parameter. init <outputDir>")
			}
		} else
		if arg2 == "process" {
			if len(argsWithProg) == 5 {
				searchDir = argsWithProg[2]
				prefixDir = argsWithProg[3]
				outputDir = argsWithProg[4]
				processPages(searchDir, prefixDir, outputDir)
			} else {
				fmt.Println("process needs 3 parameter. process <searchDir> <prefixDir> <outputDir>")
			}
		} else {

		}


	} else {
		// print Help
		fmt.Println(`Go Server Pages - pre processor v0.1 (c) gunesh.raj@gmail.com
Converts *.gun files to Go files.
  init <outputDir>
    Creates skeleton code in a folder
  process <searchDir> <prefixDir> <outputDir>
    Processes *.gun files recursively & creates *.go files in outputDir
`)

	}






}


func processPages(searchDir string, prefixDir string, outputDir string) {
	fileList := getFileList(searchDir)
	methodList := []string{}
	for _, file := range fileList {
		_, methodNameRet, _ := parseFile(file, outputDir, prefixDir)
		methodList = append(methodList, methodNameRet)
		//fmt.Println(file)
	}
	createRouteMapFile(outputDir, methodList)
}


func createRouteMapFile(outputDir string, methodName []string) {

	header := `// Auto generated Code. Do no modify below this.
package main

import (
	"github.com/labstack/echo"
	"strconv"
)

func routeMap(e *echo.Echo) {

`

	footer := `

}


func WriteStr(c echo.Context, s string) {
	c.Response().Writer.Write([]byte(s))
}

func WriteInt(c echo.Context, s int) {
	c.Response().Writer.Write([]byte(strconv.Itoa(s)))
}

// End of Auto generated Code. Do not modify above this.
`


	var buffer bytes.Buffer
	for _, mr := range methodName {
		// e.GET("/123", the_method_123_gun)
		// e.POST("/123", the_method_123_gun)
		buffer.WriteString("e.GET(\"/" + mr + "\", " + mr + ")\n")
		buffer.WriteString("e.POST(\"/" + mr + "\", " + mr + ")\n")
	}



	outStr := header + buffer.String() + footer
	writeStringToFile(outputDir + "routemap.go", outStr)
}


func getFileList(searchDir string) []string {
	fileList := []string{}
	err := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		if (strings.HasSuffix(path, ".gun") == true) { // .gun
			fileList = append(fileList, path)
		}
		return nil
	})
	if (err != nil) {
		//log.Fatal(err)
	}
	return fileList
}

func parseFile(path string, outputPath string, prefixDir string) (fileNameRet string, methodNameRet string, outputFileNameRet string) {
	fileName := ""
	methodName := ""
	outputFileName := ""

	//fmt.Println("-",path, prefixDir )
	fileName = strings.TrimPrefix(path, prefixDir)
	fileName = strings.Replace(fileName, "\\", "_", -1)
	methodName = strings.TrimSuffix(fileName, ".gun") // .gun
	fileName = strings.Replace(fileName, ".", "_", -1)
	methodName = strings.Replace(methodName, ".", "_", -1)
	fileName = fileName + ".go"
	outputFileName = outputPath + fileName
	fmt.Println(fileName, ":", methodName, ":", outputFileName)

	str := readFileToString(path)
	//fmt.Println(str) // print the content as a 'string'

	outStr := generateCode(str, methodName)
	//fmt.Println(outStr)
	//str := readFileToString("/Users/ice/GoglandProjects/gtest1/av.gun")
	//generateCode(str, "the_method_123_gun")
	writeStringToFile(outputFileName, outStr)
	return fileName, methodName, outputFileName
}

func generateCode(dx string, methodName string) string {
	outStr := ""

	dataR := [][]string{}

	aORb := regexp.MustCompile("<%|%>")
	matches := aORb.FindAllStringIndex(dx, -1)
	//fmt.Println(matches)
	iLen := len(matches)

	for i := 0; i < iLen; i++ {
		mx := matches[i]
		if mx[0] > 0 && i == 0 {
			ds := []string{"X", dx[0: mx[0]]}
			dataR = append(dataR, ds)
		} else {
			if i > 0 {
				sp := ""
				mz := matches[i-1]
				if dx[mx[0]:mx[1]] == "<%" {
					sp = "O"
				}
				if dx[mx[0]:mx[1]] == "%>" {
					sp = "X"
				}
				ds := []string{sp, dx[mz[1]: mx[0]]}
				dataR = append(dataR, ds)
			}
		}
		if i == (iLen - 1) {
			if mx[1] != len(dx) {
				ds := []string{"O", dx[mx[1]:]}
				dataR = append(dataR, ds)
			}
		}

	}

	var buffer bytes.Buffer
	for _, dr := range dataR {
		//fmt.Println("[", dr[0], ":", dr[1], "]")
		if dr[0] == "X" {
			buffer.WriteString(dr[1] + "\n")
		}
		if dr[0] == "O" {
			buffer.WriteString("WriteStr(c, `" + dr[1] + "`)\n")
		}

	}
	header := `// Auto generated Code. Do no modify below this.
package main

import (
	"net/http"
	"github.com/labstack/echo"
)

`

	footer := `

	c.Response().Flush()
	c.Response().WriteHeader(http.StatusOK)
	return nil
}

// End of Auto generated Code. Do not modify above this.
`

	outStr = header + "\nfunc " + methodName + "(c echo.Context) error {\n" + buffer.String() + footer
	return outStr
}



func readFileToString(path string) string {
	b, err := ioutil.ReadFile(path) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	str := string(b) // convert content to a 'string'
	return str
}

func writeStringToFile(path string, content string) {
	ioutil.WriteFile(path, []byte(content), 0644)
}


func createInitial(path string) {
	// generate initial files & templates in current folder.
	mainPageText := `
// Auto generated code
package main

import (
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	// Start: Enter your Routes here.

	// End: Enter your Routes here.
	routeMap(e)
	e.Logger.Fatal(e.Start(":1323"))
}

`


	routePageText := `
// Auto generated Code. Do no modify below this.
package main

import (
	"github.com/labstack/echo"
	"strconv"
)

func routeMap(e *echo.Echo) {

}


func WriteStr(c echo.Context, s string) {
	c.Response().Writer.Write([]byte(s))
}

func WriteInt(c echo.Context, s int) {
	c.Response().Writer.Write([]byte(strconv.Itoa(s)))
}

// End of Auto generated Code. Do not modify above this.

`

	writeStringToFile(path + "main.go", mainPageText)
	writeStringToFile(path + "routemap.go", routePageText)



}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
