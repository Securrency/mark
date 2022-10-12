package headers

import (
	"bufio"
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var DisclaimerPath = ""
var HeadersPath = ""

func Headers_entrypoint(headersPath string, disclaimerPath string) {
	if disclaimerPath != "" {
		DisclaimerPath = disclaimerPath
	}
	if headersPath != "" {
		HeadersPath = headersPath
	}
	scan(".")

}
func readfile(file string) (string, error) {

	content, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(content), nil

}
func generateAttachmentHeaders(content string) (string, error) {
	scanner := bufio.NewScanner(strings.NewReader(content))

	result := ""
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "![") {
			attachmentPath := strings.Split(strings.Split(scanner.Text(), "(")[1], ")")[0]

			result += fmt.Sprintf("<!-- Attachment: %s -->\n", attachmentPath)
		} else {
			continue
		}
	}
	return result, nil

}
func generateTitle(path string) string {
	paths := strings.Split(path, "/")
	title := strings.Replace(paths[len(paths)-1], ".md", "", -1)

	return fmt.Sprintf("<!-- Title: %s -->\n", cases.Title(language.English).String(title))
}
func generateParentsHeaders(path string) ([]string, error) {
	parents := strings.Split(path, "/")
	parents = parents[0 : len(parents)-1]
	for i := 0; i < len(parents); i++ {
		parents[i] = fmt.Sprintf("<!-- Parent: ADR-%s -->\n", parents[i])
	}
	return parents, nil
}
func modifyHeader(path string) (string, error) {
	result := ""
	content, file_err := readfile(path)
	if file_err != nil {
		return "", file_err
	}
	scanner := bufio.NewScanner(strings.NewReader(content))

	default_headers_path := "./scripts/default-headers"
	if HeadersPath != "" {
		default_headers_path = HeadersPath
	}
	default_header, header_err := os.ReadFile(default_headers_path)
	if header_err != nil {
		return "", header_err
	}
	disclaimer_default_path := "./scripts/default-disclaimer"

	if DisclaimerPath != "" {
		disclaimer_default_path = DisclaimerPath
	}
	default_disclaimer, disclaimer_err := os.ReadFile(disclaimer_default_path)
	if disclaimer_err != nil {
		return "", disclaimer_err
	}

	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "<!--") || strings.HasSuffix(scanner.Text(), "-->") || strings.HasPrefix(scanner.Text(), ":box:") || strings.HasPrefix(scanner.Text(), "     Template:") || strings.HasPrefix(scanner.Text(), "     Icon:") || strings.HasPrefix(scanner.Text(), "     Name:") || strings.HasPrefix(scanner.Text(), "     Title:") {
			continue
		} else {
			result += scanner.Text() + "\n"
		}
	}

	attachments, attachments_err := generateAttachmentHeaders(result)
	if attachments_err != nil {
		return "", attachments_err
	}

	parents, parents_err := generateParentsHeaders(path)
	if parents_err != nil {
		return "", parents_err
	}

	title := generateTitle(path)

	content = string(default_header) + "\n" + strings.Title(strings.Join(parents, "")) + attachments + title + "\n" + string(default_disclaimer) + "\n" + result

	err := ioutil.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return "", err
	}

	return content, nil
}
func walkHandler(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if filepath.Ext(path) == ".md" {
		modifyHeader(path)
	}
	return nil
}
func scan(path string) {
	err := filepath.Walk(path, walkHandler)
	if err != nil {
	}
}
