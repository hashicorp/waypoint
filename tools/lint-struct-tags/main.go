package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <directory>\n", os.Args[0])
		os.Exit(1)
	}
	dir := os.Args[1]
	duplicatesFound := false
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, ".pb.go") &&
			!strings.HasSuffix(path, "_test.go") {
			dupes, err := analyzeFile(path)
			if err != nil {
				return err
			}
			if dupes {
				duplicatesFound = true
			}
		}
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if duplicatesFound {
		os.Exit(1)
	}
}

func analyzeFile(file string) (bool, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		return false, err
	}
	duplicatesFound := false
	ast.Inspect(node, func(n ast.Node) bool {
		s, ok := n.(*ast.StructType)
		if !ok {
			return true
		}
		seenTags := make(map[string]*token.Position)
		for _, field := range s.Fields.List {
			tag, found := getTag(field)
			if found {
				_, dup := seenTags[tag]
				if dup {
					fmt.Printf("Duplicate tag '%s' found in file '%s' at line %d\n",
						tag, file, fset.Position(field.Pos()).Line)
					duplicatesFound = true
				} else {
					seenTags[tag] = &token.Position{Line: fset.Position(field.Pos()).Line}
				}
			}
		}
		return false
	})
	return duplicatesFound, nil
}

func getTag(field *ast.Field) (string, bool) {
	if field.Tag == nil {
		return "", false
	}
	tag := field.Tag.Value
	tag = strings.Trim(tag, "`")
	tagParts := strings.Split(tag, " ")
	for _, part := range tagParts {
		match, _ := regexp.MatchString(`\w+:"[^"]*"`, part)
		if match {
			return strings.Trim(part, `"`), true
		}
	}
	return "", false
}
