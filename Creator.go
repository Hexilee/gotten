package gotten

import (
	"errors"
	"net/http"
	"net/url"
	"reflect"
)

type (
	Builder struct {
		baseUrl      string
		cookies      []*http.Cookie
		headers      http.Header
		client       Client
		unmarshalers []*ConditionalUnmarshaler
	}

	Creator struct {
		baseUrl      *url.URL
		cookies      []*http.Cookie
		headers      http.Header
		client       Client
		unmarshalers []*ConditionalUnmarshaler
	}

	ConditionalUnmarshaler struct {
		Checker
		Unmarshaler
	}
)

func NewBuilder() *Builder {
	return &Builder{
		cookies:      make([]*http.Cookie, 0),
		headers:      make(http.Header),
		unmarshalers: make([]*ConditionalUnmarshaler, 0),
	}
}

func (builder *Builder) SetBaseUrl(url string) *Builder {
	builder.baseUrl = url
	return builder
}

func (builder *Builder) AddCookie(cookie *http.Cookie) *Builder {
	builder.cookies = append(builder.cookies, cookie)
	return builder
}

func (builder *Builder) AddCookies(cookies []*http.Cookie) *Builder {
	builder.cookies = append(builder.cookies, cookies...)
	return builder
}

func (builder *Builder) SetHeader(key, value string) *Builder {
	builder.headers.Set(key, value)
	return builder
}

func (builder *Builder) AddHeader(key, value string) *Builder {
	builder.headers.Add(key, value)
	return builder
}

func (builder *Builder) AddUnmarshaler(unmarshaler Unmarshaler, checker Checker) *Builder {
	builder.unmarshalers = append(builder.unmarshalers, &ConditionalUnmarshaler{checker, unmarshaler})
	return builder
}

func (builder *Builder) AddUnmarshalFunc(unmarshaler UnmarshalFunc, checker Checker) *Builder {
	return builder.AddUnmarshaler(unmarshaler, checker)
}

func (builder *Builder) SetClient(client Client) *Builder {
	builder.client = client
	return builder
}

func (builder *Builder) Build() (creator *Creator, err error) {
	if builder.baseUrl == "" {
		err = errors.New(BaseUrlCannotBeEmpty)
	}

	if err == nil {
		var baseUrl *url.URL
		baseUrl, err = url.Parse(builder.baseUrl)
		if err == nil {
			if builder.client == nil {
				builder.client = &http.Client{}
			}
			creator = &Creator{
				baseUrl:      baseUrl,
				cookies:      builder.cookies,
				headers:      builder.headers,
				client:       builder.client,
				unmarshalers: builder.unmarshalers,
			}
		}
	}
	return
}

func (creator *Creator) Impl(service interface{}) (err error) {
	serviceVal := reflect.ValueOf(service)
	if serviceVal.Type().Kind() != reflect.Ptr {
		err = MustPassPtrToImplError(serviceVal.Type())
	}

	if err == nil {
		serviceVal = serviceVal.Elem()
		serviceType := serviceVal.Type()
		if serviceType.Kind() != reflect.Struct {
			err = ServiceMustBeStructError(serviceType)
		}

		if err == nil {
			for i := 0; i < serviceType.NumField(); i++ {
				field := serviceType.Field(i)
				fieldType := field.Type
				fieldTag := field.Tag
				fieldValue := serviceVal.Field(i)
				if fieldType.Kind() == reflect.Func &&
					fieldValue.CanSet() &&
					fieldType.NumIn() == 1 &&
					fieldType.NumOut() == 2 &&
					fieldType.Out(1) == reflect.TypeOf(errors.New("")) {
					method := fieldTag.Get(KeyMethod)
					switch method {
					case "":
						method = http.MethodGet
						fallthrough
					case http.MethodGet:
						fallthrough
					case http.MethodHead:
						fallthrough
					case http.MethodPost:
						fallthrough
					case http.MethodPut:
						fallthrough
					case http.MethodPatch:
						fallthrough
					case http.MethodDelete:
						fallthrough
					case http.MethodConnect:
						fallthrough
					case http.MethodOptions:
						fallthrough
					case http.MethodTrace:
						rawFunc := func([]reflect.Value) []reflect.Value {

							return []reflect.Value{}
						}
						fieldValue.Set(reflect.MakeFunc(fieldType, rawFunc))
					default:
						err = UnrecognizedHTTPMethodError(method)
					}
				}
			}
		}
	}
	return
}
