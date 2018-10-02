package gotten

import (
	"net/url"
)

type (
	UrlConstructor interface {
		getUrl() *url.URL
	}

	UrlBuilder struct {
		path VarsConstructor
		base *url.URL
	}

	UrlCtr struct {
		path    VarsConstructor
		queries map[string]*QueryField
		base    *url.URL
	}
)

func newUrlBuilder(base *url.URL, constructor VarsConstructor) *UrlBuilder {
	return &UrlBuilder{
		path: constructor,
		base: base,
	}
}

func (urlCtr *UrlCtr) getUrl() *url.URL {
	// TODO: *UrlCtr.getUrl
	return nil
}
