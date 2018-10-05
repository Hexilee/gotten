package gotten

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"github.com/Hexilee/gotten/headers"
	"github.com/Hexilee/unhtml"
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
		unmarshalers ConditionalUnmarshalers
	}

	Creator struct {
		baseUrl      *url.URL
		cookies      []*http.Cookie
		headers      http.Header
		client       Client
		unmarshalers ConditionalUnmarshalers
	}

	ConditionalUnmarshaler struct {
		unmarshaler ReadUnmarshaler
		checker     Checker
	}

	ConditionalUnmarshalers []*ConditionalUnmarshaler
)

var (
	DefaultUnmarshalers = []*ConditionalUnmarshaler{
		{
			NewReaderAdapter(UnmarshalAdapter(json.Unmarshal)),
			new(CheckerFactory).WhenContentType(
				headers.MIMEApplicationJSON,
				headers.MIMEApplicationJSONCharsetUTF8,
			).Create(),
		},
		{
			NewReaderAdapter(UnmarshalAdapter(xml.Unmarshal)),
			new(CheckerFactory).WhenContentType(
				headers.MIMEApplicationXML,
				headers.MIMEApplicationXMLCharsetUTF8,
				headers.MIMETextXML,
				headers.MIMETextXMLCharsetUTF8,
			).Create(),
		},
		{
			NewReaderAdapter(UnmarshalAdapter(unhtml.Unmarshal)),
			new(CheckerFactory).WhenContentType(
				headers.MIMETextHTML,
				headers.MIMETextHTMLCharsetUTF8,
			).Create(),
		},
	}
)

func NewBuilder() *Builder {
	return &Builder{
		cookies:      make([]*http.Cookie, 0),
		headers:      make(http.Header),
		unmarshalers: make(ConditionalUnmarshalers, 0),
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
	return builder.AddReaderUnmarshaler(NewReaderAdapter(unmarshaler), checker)
}

func (builder *Builder) AddUnmarshalFunc(unmarshaler UnmarshalFunc, checker Checker) *Builder {
	return builder.AddUnmarshaler(unmarshaler, checker)
}

func (builder *Builder) AddReaderUnmarshaler(unmarshaler ReadUnmarshaler, checker Checker) *Builder {
	builder.unmarshalers = append(builder.unmarshalers, &ConditionalUnmarshaler{unmarshaler, checker})
	return builder
}

func (builder *Builder) AddReaderUnmarshalerFunc(unmarshaler ReadUnmarshalFunc, checker Checker) *Builder {
	return builder.AddReaderUnmarshaler(unmarshaler, checker)
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
				unmarshalers: append(builder.unmarshalers, DefaultUnmarshalers...),
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
		FieldCycle:
			for i := 0; i < serviceType.NumField(); i++ {
				field := serviceType.Field(i)
				fieldType := field.Type
				fieldTag := field.Tag
				fieldValue := serviceVal.Field(i)
				if fieldType.Kind() == reflect.Func &&
					fieldValue.CanSet() &&
					fieldType.NumIn() == 1 &&
					fieldType.NumOut() == 2 &&
					fieldType.Out(0) == ResponseType &&
					fieldType.Out(1) == ErrorType {
					paramsType := fieldType.In(0)
					varsParser, parseErr := newVarsParser(fieldTag.Get(KeyPath))
					if err = parseErr; err != nil {
						break FieldCycle
					}

					err = varsParser.parse(paramsType)
					if err != nil {
						break FieldCycle
					}

					method := fieldTag.Get(KeyMethod)

					// TODO: add body check for different methods
					switch method {
					case "": // "" means "GET" in standard library
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
						rawFunc := func(values []reflect.Value) []reflect.Value {
							results := make([]reflect.Value, 2, 2)
							varsCtr := varsParser.Build()
							setValuesErr := varsCtr.setValues(values[0])

							if setValuesErr != nil {
								results[1] = reflect.ValueOf(setValuesErr)
								return results
							}

							finalUrl, err := newUrlCtr(creator.baseUrl, varsCtr).getUrl()
							if err != nil {
								results[1] = reflect.ValueOf(err)
								return results
							}

							req, err := http.NewRequest(method, finalUrl.String(), nil)
							if err != nil {
								results[1] = reflect.ValueOf(err)
								return results
							}

							resp, err := creator.client.Do(req)

							if err != nil {
								results[1] = reflect.ValueOf(err)
								return results
							}

							readUnmarshaler, exist := creator.unmarshalers.Check(resp)
							if !exist {
								results[1] = reflect.ValueOf(NoUnmarshalerFoundForResponseError(resp))
								return results
							}

							results[0] = reflect.ValueOf(newResponse(resp, readUnmarshaler))
							return results
						}
						fieldValue.Set(reflect.MakeFunc(fieldType, rawFunc))
					default:
						err = UnrecognizedHTTPMethodError(method)
						break FieldCycle
					}
				}
			}
		}
	}
	return
}

func (unmarshalers ConditionalUnmarshalers) Check(response *http.Response) (unmarshaler ReadUnmarshaler, exist bool) {
	for _, conditional := range unmarshalers {
		if conditional.checker.Check(response) {
			unmarshaler = conditional.unmarshaler
			exist = true
			break
		}
	}
	return
}
