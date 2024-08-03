package inertia

import (
	"fmt"
	"strings"

	"github.com/refiber/framework/vite"
)

// As much as possible, try to avoid using regexp, use strings instead

type PreRenderInterface interface {
	GetSafeHTML() string
}

func extractBody(html string) string {
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

func newPreRender(html *string, sourceFilePath *string, scriptBuf []byte, vds *viewDataStruct) *preRender {
	pr := preRender{
		base:        *html,
		odlJsScript: vite.CreateScriptTag(fmt.Sprintf(`/build/%s`, *sourceFilePath)),
		newJsScript: []byte(fmt.Sprintf(`<script type="module">%s</script>`, string(scriptBuf))),
		vds:         vds,
	}

	pr.injectedScript = strings.Replace(
		*html,
		vds.AllTags,
		fmt.Sprintf(`
      %s 
      %s
    `, string(pr.newJsScript), vds.DataPageDiv),
		1)

	return &pr
}

// base: The HTML response from the controller.
// injectedScript: The JavaScript script tag from the base HTML, replaced by compiled JavaScript (inlined).
// rendered: The injectedScript rendered by a headless browser.
type preRender struct {
	base           string
	injectedScript string
	odlJsScript    string
	newJsScript    []byte
	rendered       *string
	vds            *viewDataStruct
}

// Will return raw html for pre-render
// TODO: get html tag and title tag from template (preRender.base)
func (pr *preRender) GetSafeHTML() string {
	body := fmt.Sprintf(`
    %s
    %s
  `, string(pr.newJsScript), pr.vds.DataPageDiv)

	html := fmt.Sprintf(`
    <!DOCTYPE html>
    <html lang="en">
    <head>
      <meta charset="UTF-8">
      <meta name="viewport" content="width=device-width, initial-scale=1.0">
      <title>Refiber</title>
    </head>
    <body>
      %s
    </body>
    </html>
  `, body)

	return html
}

func (pr *preRender) createClientHTML() string {
	preRenderedBody := strings.ReplaceAll(
		extractBody(*pr.rendered),
		string(pr.newJsScript),
		pr.vds.ScriptTags,
	)

	return strings.Replace(strings.Replace(pr.injectedScript, pr.vds.DataPageDiv, "", 1), string(pr.newJsScript), preRenderedBody, 1)
}
