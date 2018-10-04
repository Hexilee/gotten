package gotten

import (
	"net/url"
)

type (
	UrlCtr struct {
		vars VarsController
		base *url.URL
	}
)

func newUrlCtr(base *url.URL, vars VarsController) *UrlCtr {
	return &UrlCtr{
		vars: vars,
		base: base,
	}
}

func (urlCtr *UrlCtr) getUrl() (result *url.URL, err error) {
	result, err = urlCtr.vars.getUrl()
	if err == nil {
		result = urlCtr.base.ResolveReference(result)
	}
	return
}
