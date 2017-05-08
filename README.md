# gspro
Go Server Pages Pre-processor

Go Server Pages - pre processor v0.1 (c) gunesh.raj@gmail.com
Converts *.gun files to Go files.
  init <outputDir>
    Creates skeleton code in a folder
  process <searchDir> <prefixDir> <outputDir>
    Processes *.gun files recursively & creates *.go files in outputDir

This project simplifies web development similar to JSP and PHP. The Templates are compiled into the Executable and technically performs better than the typical parsed templates.
It uses the Echo Framework - https://github.com/labstack/echo

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



