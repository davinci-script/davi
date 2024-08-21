package interpreter

import (
	"encoding/json"
	"fmt"
	"github.com/DavinciScript/Davi/interpreter/functions"
	goParser "go/parser"
	goToken "go/token"
	"io"
	"os"
	"slices"
	"strings"
)

//
//var orderFunctionCategories = map[string]int{
//	"String":      1,
//	"Array":       2,
//	"Conversion":  3,
//	"System":      4,
//	"File System": 5,
//	"HTTP":        6,
//}

func GenerateDocs() {

	daviFileWithFunctions := `tests/functions_test.davi`
	markdownFileWithFunctions := `docs/docs/guide/functions.md`
	jsonFileReal := `docs/docs/.vuepress/dist/davi-details.json`

	daviFile, err := os.Create(daviFileWithFunctions)
	if err != nil {
		fmt.Println(err)
	}
	defer daviFile.Close()

	// open the file
	file, err := os.Create(markdownFileWithFunctions)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	// open the file
	jsonFile, err := os.Create(jsonFileReal)
	if err != nil {
		fmt.Println(err)
	}

	// write the davi content
	daviContent := "<?davi \n // DaVinci Script \n\n"

	// write the markdown content
	markdownContent := "# Functions \n\n"

	functionsCategories := []string{}
	functionDetails := GetFunctionsDetails()
	for _, f := range functionDetails {
		if f.category != "" {
			if !slices.Contains(functionsCategories, f.category) {
				functionsCategories = append(functionsCategories, f.category)
			}
		}
	}

	jsonContent := make(map[string]interface{})
	if len(functionsCategories) > 0 {

		for _, category := range functionsCategories {

			daviContent += "// Category:  " + category + "\n\n"
			markdownContent += "## " + category + "\n\n"

			jsonContent[category] = make(map[string]interface{})

			for _, f := range functionDetails {
				if f.category == category {

					// put f.functionName on jsonContent[category] on first level
					jsonContent[category].(map[string]interface{})[f.functionName] = map[string]interface{}{
						"args":        f.args,
						"returnValue": f.returnValue,
						"example":     f.example,
						"output":      f.output,
						"description": f.description,
						"title":       f.title,
						"category":    f.category,
					}

					markdownContent += "### " + f.title + "\n\n"
					markdownContent += "```php\n"
					markdownContent += f.functionName + "(" + f.args + ")\n"
					markdownContent += "```\n\n"
					markdownContent += f.description + "\n\n"
					markdownContent += "#### Example\n\n"
					markdownContent += "```php\n"
					markdownContent += f.example + "\n\n"
					markdownContent += "// output: " + f.output + "\n"
					markdownContent += "```\n\n"

					variable := "$call" + functions.UpWords(f.functionName)
					daviContent += "// " + f.title + "\n"
					daviContent += variable + " = " + f.example + "\n\n"
					daviContent += "// must output: " + f.output + "\n\n"
					daviContent += "echo(" + variable + ")\n\n"
				}
			}
		}
	}

	// write the json content
	outputJson, _ := json.Marshal(jsonContent)
	print(string(outputJson))
	_, err = jsonFile.WriteString(string(outputJson))
	if err != nil {
		fmt.Println(err)
	}

	_, err = file.WriteString(markdownContent)
	if err != nil {
		fmt.Println(err)
	}

	_, err = daviFile.WriteString(daviContent)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Documentation generated successfully.")

}

func GetFunctionsDetails() map[string]functionDetails {

	// all functions file path
	allFunctionsFile := "interpreter/functions.go"

	// open the file
	file, err := os.Open(allFunctionsFile)
	if err != nil {
		fmt.Println(err)
	}
	// get content of the file
	content, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	contentString := string(content)

	fs := goToken.NewFileSet()
	f, err := goParser.ParseFile(fs, "", contentString, goParser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	allFunctionsDetails := make(map[string]functionDetails)
	for _, c := range f.Comments {
		parsedComment := ParseComment(c.Text())
		allFunctionsDetails[parsedComment.functionName] = parsedComment
	}

	return allFunctionsDetails
}

type functionDetails struct {
	functionName string
	args         string
	returnValue  string
	example      string
	output       string
	description  string
	title        string
	category     string
}

func ParseComment(comment string) functionDetails {
	/**
	 * function: httpRegister
	 * args: pattern, handler
	 * return: nil
	 * example: httpRegister("/", func() { return "Hello, World!" })
	 * output: "Hello, World!"
	 * description: Register a handler function for a URL pattern.
	 * title: HTTP Register
	 * category: HTTP
	 */

	functionName := ""
	args := ""
	returnValue := ""
	example := ""
	output := ""
	description := ""
	title := ""
	category := ""

	// split the comment by new line
	commentLines := strings.Split(comment, "\n")
	for _, line := range commentLines {
		if strings.Contains(line, "function:") {
			functionName = strings.TrimSpace(strings.Split(line, ":")[1])
		}
		if strings.Contains(line, "args:") {
			args = strings.TrimSpace(strings.Split(line, ":")[1])
		}
		if strings.Contains(line, "return:") {
			returnValue = strings.TrimSpace(strings.Split(line, ":")[1])
		}
		if strings.Contains(line, "example:") {
			example = strings.TrimSpace(strings.Split(line, ":")[1])
		}
		if strings.Contains(line, "output:") {
			output = strings.TrimSpace(strings.Split(line, ":")[1])
		}
		if strings.Contains(line, "description:") {
			description = strings.TrimSpace(strings.Split(line, ":")[1])
		}
		if strings.Contains(line, "title:") {
			title = strings.TrimSpace(strings.Split(line, ":")[1])
		}
		if strings.Contains(line, "category:") {
			category = strings.TrimSpace(strings.Split(line, ":")[1])
		}
	}

	return functionDetails{
		functionName: functionName,
		args:         args,
		returnValue:  returnValue,
		example:      example,
		output:       output,
		description:  description,
		title:        title,
		category:     category,
	}
}
