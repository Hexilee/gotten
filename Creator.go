package gotten

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"github.com/Hexilee/gotten/headers"
	"github.com/Hexilee/unhtml"
	"io"
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

func (builder *Builder) AddReadUnmarshalFunc(unmarshaler ReadUnmarshalFunc, checker Checker) *Builder {
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

// func(*params) (*http.Request, error) ||
// func(*params) (gotten.Response, error)
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
					fieldType.Out(1) == ErrorType {
					paramsType := fieldType.In(0)
					varsParser, parseErr := newVarsParser(fieldTag.Get(KeyPath))
					if err = parseErr; err == nil {
						err = varsParser.parse(paramsType)
						if err == nil {
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
								switch fieldType.Out(0) {
								case ResponseType:
									fieldValue.Set(reflect.MakeFunc(fieldType, creator.getCompleteFunc(varsParser, method)))
								case RequestType:
									fieldValue.Set(reflect.MakeFunc(fieldType, creator.getRequestFunc(varsParser, method)))
								default:
								}

							default:
								err = UnrecognizedHTTPMethodError(method)
							}
						}
					}

					if err != nil {
						break
					}
				}
			}
		}
	}
	return
}

// for func(*params) (*http.Request, error)
func (creator Creator) getRequestFunc(varsParser *VarsParser, method string) func([]reflect.Value) []reflect.Value {
	return func(values []reflect.Value) []reflect.Value {
		results := []reflect.Value{
			reflect.New(RequestType).Elem(),
			reflect.New(ErrorType).Elem(),
		}
		varsCtr := varsParser.Build()
		setValuesErr := varsCtr.setValues(values[0])

		if setValuesErr != nil {
			results[1].Set(reflect.ValueOf(setValuesErr).Convert(ErrorType))
			return results
		}

		finalUrl, err := newUrlCtr(creator.baseUrl, varsCtr).getUrl()
		// err always be nil if all test pass
		//if err != nil {
		//	results[1].Set(reflect.ValueOf(err).Convert(ErrorType))
		//	return results
		//}

		var body io.Reader
		contentType := varsCtr.getContentType()

		if contentType != ZeroStr {
			body, err = varsCtr.getBody()
			if err != nil {
				results[1].Set(reflect.ValueOf(err).Convert(ErrorType))
				return results
			}
		}

		req, err := http.NewRequest(method, finalUrl.String(), body)
		// err always be nil with checked method and URL
		//if err != nil {
		//	results[1].Set(reflect.ValueOf(err).Convert(ErrorType))
		//	return results
		//}

		for key, values := range creator.headers {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		// cover header of creator
		for key, values := range varsCtr.getHeader() {
			for _, value := range values {
				req.Header.Set(key, value)
			}
		}

		// cover all ContentType
		if contentType != ZeroStr {
			req.Header.Set(headers.HeaderContentType, contentType)
		}

		// add cookie of creator
		for _, cookie := range creator.cookies {
			req.AddCookie(cookie)
		}

		// add cookie of VarsCtr
		for _, cookie := range varsCtr.getCookies() {
			req.AddCookie(cookie)
		}

		results[0].Set(reflect.ValueOf(req).Convert(RequestType))
		return results
	}
}

// for func(*params) (gotten.Response, error)
func (creator Creator) getCompleteFunc(varsParser *VarsParser, method string) func([]reflect.Value) []reflect.Value {
	return func(values []reflect.Value) []reflect.Value {
		results := creator.getRequestFunc(varsParser, method)(values)
		req := results[0].Interface().(*http.Request)
		results[0] = reflect.New(ResponseType).Elem()
		if results[1].IsNil() {
			resp, err := creator.client.Do(req)

			if err != nil {
				results[1].Set(reflect.ValueOf(err).Convert(ErrorType))
				return results
			}

			readUnmarshaler, exist := creator.unmarshalers.Check(resp)
			if !exist {
				results[1].Set(reflect.ValueOf(NoUnmarshalerFoundForResponseError(resp)).Convert(ErrorType))
				results[0].Set(reflect.ValueOf(newResponse(resp, nil)).Convert(ResponseType))
				return results
			}

			results[0].Set(reflect.ValueOf(newResponse(resp, readUnmarshaler)).Convert(ResponseType))
		}
		return results
	}
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
