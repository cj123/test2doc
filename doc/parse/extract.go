package parse

import (
	"regexp"
	"strings"

	"github.com/fatih/camelcase"
)

// GetTitle extracts a title from the function name,
// where longFnName is of the form:
// github.com/cj123/test2doc/example.GetWidget
// and the out title would be:
// Handle Get Widget
func GetTitle(longFnName string) string {
	shortFnName := getShortFnName(longFnName)

	words := camelcase.Split(shortFnName)

	for i := range words {
		words[i] = strings.Title(words[i])
	}

	return strings.Join(words, " ")
}

func GetDescription(longFnName string) (desc string) {
	shortFnName := getShortFnName(longFnName)

	doc := funcsMap[shortFnName]
	if doc != nil {
		desc = strings.TrimPrefix(doc.Doc, shortFnName+" ")
	}

	return
}

// IsFuncInPkg checks if this func belongs to the package
func IsFuncInPkg(longFnName string) bool {
	shortFnName := getShortFnName(longFnName)
	doc := funcsMap[shortFnName]
	return doc != nil
}

// getShortFnName returns the name of the function
// without the package name so:
//   github.com/user/project/package.method
// becomes
//   method
// and
//   github.com/user/project/package.(*type).method
// becomes
//   type.method
func getShortFnName(longFnName string) string {
	methodRE := regexp.MustCompile(`/(.*)\.(.*)\.(.*)`)
	funcRE := regexp.MustCompile(`/(.*)\.(.*)`)

	matches := methodRE.FindStringSubmatch(longFnName)
	if len(matches) > 0 {
		fnName := strings.Join(matches[len(matches)-2:], ".")
		fnName = strings.Replace(fnName, "(*", "", -1)
		return strings.Replace(fnName, ")", "", -1)
	}

	matches = funcRE.FindStringSubmatch(longFnName)
	if len(matches) > 0 {
		return matches[len(matches)-1]
	}

	return ""
}
