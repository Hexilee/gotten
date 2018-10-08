package gotten

import (
	"net/url"
	"strings"
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
		base := *urlCtr.base
		base.Path = strings.TrimRight(base.Path, "/") + "/" + strings.TrimLeft(result.Path, "/")
		base.RawQuery = result.RawQuery
		result = &base
	}
	return
}
