package inertia

import (
	"fmt"
	"strings"

	"github.com/refiber/framework/vite"
)

// The goal of pre-rendering is to render the app before sending it to the user, similar to Server Side Rendering (SSR), but without requiring Node.js.
// How does it work? Refiber provides safe HTML data (HTML + inline JavaScript) to your app through PreRenderHandler. Your app then uses a headless browser to render the HTML and JavaScript. Tools like chromedp can help with this process, and you will need Chrome or Chromium installed on your machine.
// What about performance? I didn't use any fancy testing tools, just relied on Chrome DevTools to check request times in the Network tab. I'm using a MacBook Pro M1 with 16 GB of memory, so the results might be different on your machine. Using the default Refiber template v0.1.0-beta, I observed the following:
// - Without pre-render: 3ms request time.
// - With pre-render: 48ms request time.

type PreRenderInterface interface {
	GetSafeHTML() string
	GetCompailedJs() []byte
	GetProps() []byte
}

func newPreRender(html *string, sourceFilePath *string, scriptBuf []byte, vds *viewDataStruct, jsonProps []byte) *preRender {
	pr := preRender{
		base:        *html,
		odlJsScript: vite.CreateScriptTag(fmt.Sprintf(`/build/%s`, *sourceFilePath)),
		newJsScript: []byte(fmt.Sprintf(`<script type="module">%s</script>`, string(scriptBuf))),
		vds:         vds,
		jsonProps:   jsonProps,
		compailedJs: scriptBuf,
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
	jsonProps      []byte
	compailedJs    []byte
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

func (pr *preRender) GetProps() []byte {
	return pr.jsonProps
}

func (pr *preRender) GetCompailedJs() []byte {
	return pr.compailedJs
}

func (pr *preRender) createClientHTML() []byte {
	preRenderedBody := strings.ReplaceAll(
		extractBody(*pr.rendered),
		string(pr.newJsScript),
		pr.vds.ScriptTags,
	)

	return []byte(strings.Replace(strings.Replace(pr.injectedScript, pr.vds.DataPageDiv, "", 1), string(pr.newJsScript), preRenderedBody, 1))
}
