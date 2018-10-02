package gotten

import (
	"net/url"
)

type (
	UrlCtr struct {
		vars VarsConstructor
		base *url.URL
	}
)

func newUrlCtr(base *url.URL, vars VarsConstructor) *UrlCtr {
	return &UrlCtr{
		vars: vars,
		base: base,
	}
}

func (urlCtr *UrlCtr) getUrl() string {
	return urlCtr.base.ResolveReference(urlCtr.vars.getUrl()).String()
}
