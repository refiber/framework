package inertia

import "strings"

func extractBody(html string) string {
	// As much as possible, try to avoid using regexp, use strings instead
	startTag := "<body>"
	endTag := "</body>"

	startIndex := strings.Index(html, startTag)
	if startIndex == -1 {
		return ""
	}

	startIndex += len(startTag)
	endIndex := strings.Index(html[startIndex:], endTag)
	if endIndex == -1 {
		return ""
	}

	endIndex += startIndex
	return html[startIndex:endIndex]
}
