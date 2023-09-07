package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type FuncInfo struct {
	Name       string
	ReturnType string
}

var funcMap map[string]FuncInfo

func main() {
	// Check if there is at least one command-line argument
	if len(os.Args) < 2 {
		fmt.Println("Please provide a directory name")
		os.Exit(1)
	}
	// Get the last argument as the directory name
	dirName := os.Args[len(os.Args)-1]
	// Walk through the directory and its subdirectories
	err := filepath.Walk(dirName, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Check if the file has the .rgo extension
		if filepath.Ext(path) == ".rgo" {
			// Read the file content
			content, err := os.ReadFile(path)
			if err != nil {
				fmt.Println("Error reading file:", err)
				return err
			}
			code := string(content)
			funcMap = make(map[string]FuncInfo)
			findFunctions(code)
			// Add import statement for genericutils after the package <package> statement
			const genericUtilsPath = "github.com/acheong08/rusty-go/genericutils"
			// Get package name with regex
			packageRegex := regexp.MustCompile(`package (\w+)`)
			packageMatch := packageRegex.FindStringSubmatch(code)
			packageName := packageMatch[1]
			// Add import statement after package statement
			code = strings.Replace(code, "package "+packageName, "package "+packageName+"\n\nimport (\n\t\""+genericUtilsPath+"\"\n)", 1)

			code = replaceQuestionMarks(code)
			// Write the modified code to a new file with the .go extension
			err = os.WriteFile(path[:len(path)-4]+".go", []byte(code), 0644)
			if err != nil {
				fmt.Println("Error writing file:", err)
				return err
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error walking directory:", err)
		os.Exit(1)
	}
}
func findFunctions(code string) {
	funcRegex := regexp.MustCompile(`func (\w+)\((.*?)\) \((.*?)\)`)
	matches := funcRegex.FindAllStringSubmatch(code, -1)
	for _, match := range matches {
		funcName := match[1]
		returnType := match[3]
		funcInfo := FuncInfo{
			Name:       funcName,
			ReturnType: returnType,
		}
		funcMap[funcName] = funcInfo
	}
}
func replaceQuestionMarks(code string) string {
	callRegex := regexp.MustCompile(`(\w+), err := (\w+)\((.*?)\)\?`)
	matches := callRegex.FindAllStringSubmatchIndex(code, -1)
	for i := len(matches) - 1; i >= 0; i-- {
		match := matches[i]
		start := match[0]
		end := match[1]
		funcName := code[match[4]:match[5]]
		funcInfo, ok := funcMap[funcName]
		if !ok {
			fmt.Println("Error: function not found:", funcName)
			os.Exit(1)
		}
		returnType := funcInfo.ReturnType
		if returnType == "" {
			fmt.Println("Error: function has no return type:", funcName)
			os.Exit(1)
		}
		returnTypes := strings.Split(returnType, ",")
		for i, t := range returnTypes {
			returnTypes[i] = strings.TrimSpace(t)
		}
		defaultValues := make([]string, len(returnTypes))
		for i, t := range returnTypes {
			switch t {
			case "error":
				defaultValues[i] = "err"
			default:
				defaultValues[i] = fmt.Sprintf("genericutils.MakeGenericWithDefault[%s]()", t)
			}
		}
		defaultValue := strings.Join(defaultValues, ", ")
		replacement := fmt.Sprintf("%s, err := %s(%s)\n\tif err != nil {\n\t\treturn %s\n\t}", code[match[2]:match[3]], funcName, code[match[6]:match[7]], defaultValue)
		code = code[:start] + replacement + code[end:]
	}
	return code
}
