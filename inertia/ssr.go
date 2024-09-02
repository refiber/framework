package inertia

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type SSRInterface interface {
	GetBundledScripts() []byte
	GetScripts() map[string][]byte
}

var (
	textEncoderPolyfill = `function TextEncoder(){}TextEncoder.prototype.encode=function(string){var octets=[];var length=string.length;var i=0;while(i<length){var codePoint=string.codePointAt(i);var c=0;var bits=0;if(codePoint<=0x0000007F){c=0;bits=0x00}else if(codePoint<=0x000007FF){c=6;bits=0xC0}else if(codePoint<=0x0000FFFF){c=12;bits=0xE0}else if(codePoint<=0x001FFFFF){c=18;bits=0xF0}octets.push(bits|(codePoint>>c));c-=6;while(c>=0){octets.push(0x80|((codePoint>>c)&0x3F));c-=6}i+=codePoint>=0x10000?2:1}return octets};function TextDecoder(){}TextDecoder.prototype.decode=function(octets){var string="";var i=0;while(i<octets.length){var octet=octets[i];var bytesNeeded=0;var codePoint=0;if(octet<=0x7F){bytesNeeded=0;codePoint=octet&0xFF}else if(octet<=0xDF){bytesNeeded=1;codePoint=octet&0x1F}else if(octet<=0xEF){bytesNeeded=2;codePoint=octet&0x0F}else if(octet<=0xF4){bytesNeeded=3;codePoint=octet&0x07}if(octets.length-i-bytesNeeded>0){var k=0;while(k<bytesNeeded){octet=octets[i+k+1];codePoint=(codePoint<<6)|(octet&0x3F);k+=1}}else{codePoint=0xFFFD;bytesNeeded=octets.length-i}string+=String.fromCodePoint(codePoint);i+=bytesNeeded+1}return string};`
	processPolyfill     = `var process = {env: {NODE_ENV: "production"}};`
	consolePolyfill     = `var console = {log: function(){}};`
	ssrFilePath         = "./public/build/ssr/ssr.js"
)

func newSSR(html, props []byte, propsDivTag *string) (*ssr, error) {
	scriptBuf, err := os.ReadFile(ssrFilePath)
	if err != nil {
		return nil, err
	}

	scripts := make(map[string][]byte, 2)
	scripts["main.js"] = append([]byte(fmt.Sprintf(`%s`, textEncoderPolyfill+processPolyfill+consolePolyfill)), scriptBuf...)
	scripts["run.js"] = []byte(fmt.Sprintf(`(async () => JSON.stringify(await renderApp(%s)) )();`, string(props)))

	bundledScript := append(scripts["main.js"], scripts["run.js"]...)

	ssr := ssr{html: html, propsDivTag: propsDivTag, bundledScript: bundledScript, scripts: scripts}
	return &ssr, nil
}

type ssr struct {
	html          []byte
	bundledScript []byte
	results       *string
	propsDivTag   *string
	scripts       map[string][]byte
}

func (s *ssr) GetScripts() map[string][]byte {
	return s.scripts
}

func (s *ssr) GetBundledScripts() []byte {
	return s.bundledScript
}

func (s *ssr) createClientHTML() ([]byte, error) {
	if s.results == nil {
		return nil, fmt.Errorf("ssr results not found")
	}

	var resultData map[string]interface{}
	if err := json.Unmarshal([]byte(*s.results), &resultData); err != nil {
		return nil, err
	}

	body, exist := resultData["body"]
	if !exist {
		return nil, fmt.Errorf("ssr results body not found")
	}
	strBody, ok := body.(string)
	if !ok {
		return nil, fmt.Errorf("Failed when convert body interface to string")
	}

	results := strings.Replace(string(s.html), *s.propsDivTag, strBody, 1)
	return []byte(results), nil
}
