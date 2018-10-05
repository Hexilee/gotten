package gotten

const (
// support types: fmt.Stringer, int, string

)

//type Header string
//
//func (header Header) Value() string {
//	return string(header)
//}
//
//func GetKey(t reflect.Type) (key string, exist bool) {
//	key, exist = HeaderTable[t]
//	return
//}
//
//// generated
//type (
//	HeaderAccept interface {
//		Value() string
//	}
//	HeaderAcceptEncoding interface {
//		Value() string
//	}
//	HeaderAllow interface {
//		Value() string
//	}
//	HeaderAuthorization interface {
//		Value() string
//	}
//	HeaderContentDisposition interface {
//		Value() string
//	}
//	HeaderContentEncoding interface {
//		Value() string
//	}
//	HeaderContentLength interface {
//		Value() string
//	}
//	HeaderContentType interface {
//		Value() string
//	}
//	HeaderCookie interface {
//		Value() string
//	}
//	HeaderSetCookie interface {
//		Value() string
//	}
//	HeaderIfModifiedSince interface {
//		Value() string
//	}
//	HeaderLastModified interface {
//		Value() string
//	}
//	HeaderLocation interface {
//		Value() string
//	}
//	HeaderUpgrade interface {
//		Value() string
//	}
//	HeaderVary interface {
//		Value() string
//	}
//	HeaderWWWAuthenticate interface {
//		Value() string
//	}
//	HeaderXForwardedFor interface {
//		Value() string
//	}
//	HeaderXForwardedProto interface {
//		Value() string
//	}
//	HeaderXForwardedProtocol interface {
//		Value() string
//	}
//	HeaderXForwardedSsl interface {
//		Value() string
//	}
//	HeaderXUrlScheme interface {
//		Value() string
//	}
//	HeaderXHTTPMethodOverride interface {
//		Value() string
//	}
//	HeaderXRealIP interface {
//		Value() string
//	}
//	HeaderXRequestID interface {
//		Value() string
//	}
//	HeaderServer interface {
//		Value() string
//	}
//	HeaderOrigin interface {
//		Value() string
//	}
//	// Access control
//	HeaderAccessControlRequestMethod interface {
//		Value() string
//	}
//	HeaderAccessControlRequestHeaders interface {
//		Value() string
//	}
//	HeaderAccessControlAllowOrigin interface {
//		Value() string
//	}
//	HeaderAccessControlAllowMethods interface {
//		Value() string
//	}
//	HeaderAccessControlAllowHeaders interface {
//		Value() string
//	}
//	HeaderAccessControlAllowCredentials interface {
//		Value() string
//	}
//	HeaderAccessControlExposeHeaders interface {
//		Value() string
//	}
//	HeaderAccessControlMaxAge interface {
//		Value() string
//	}
//	// Security
//	HeaderStrictTransportSecurity interface {
//		Value() string
//	}
//	HeaderXContentTypeOptions interface {
//		Value() string
//	}
//	HeaderXXSSProtection interface {
//		Value() string
//	}
//	HeaderXFrameOptions interface {
//		Value() string
//	}
//	HeaderContentSecurityPolicy interface {
//		Value() string
//	}
//	HeaderXCSRFToken interface {
//		Value() string
//	}
//)
//
//var (
//	headerAccept              HeaderAccept
//	headerAcceptEncoding      HeaderAcceptEncoding
//	headerAllow               HeaderAllow
//	headerAuthorization       HeaderAuthorization
//	headerContentDisposition  HeaderContentDisposition
//	headerContentEncoding     HeaderContentEncoding
//	headerContentLength       HeaderContentLength
//	headerContentType         HeaderContentType
//	headerCookie              HeaderCookie
//	headerSetCookie           HeaderSetCookie
//	headerIfModifiedSince     HeaderIfModifiedSince
//	headerLastModified        HeaderLastModified
//	headerLocation            HeaderLocation
//	headerUpgrade             HeaderUpgrade
//	headerVary                HeaderVary
//	headerWWWAuthenticate     HeaderWWWAuthenticate
//	headerXForwardedFor       HeaderXForwardedFor
//	headerXForwardedProto     HeaderXForwardedProto
//	headerXForwardedProtocol  HeaderXForwardedProtocol
//	headerXForwardedSsl       HeaderXForwardedSsl
//	headerXUrlScheme          HeaderXUrlScheme
//	headerXHTTPMethodOverride HeaderXHTTPMethodOverride
//	headerXRealIP             HeaderXRealIP
//	headerXRequestID          HeaderXRequestID
//	headerServer              HeaderServer
//	headerOrigin              HeaderOrigin
//	// Access control
//	headerAccessControlRequestMethod    HeaderAccessControlRequestMethod
//	headerAccessControlRequestHeaders   HeaderAccessControlRequestHeaders
//	headerAccessControlAllowOrigin      HeaderAccessControlAllowOrigin
//	headerAccessControlAllowMethods     HeaderAccessControlAllowMethods
//	headerAccessControlAllowHeaders     HeaderAccessControlAllowHeaders
//	headerAccessControlAllowCredentials HeaderAccessControlAllowCredentials
//	headerAccessControlExposeHeaders    HeaderAccessControlExposeHeaders
//	headerAccessControlMaxAge           HeaderAccessControlMaxAge
//	// Security
//	headerStrictTransportSecurity HeaderStrictTransportSecurity
//	headerXContentTypeOptions     HeaderXContentTypeOptions
//	headerXXSSProtection          HeaderXXSSProtection
//	headerXFrameOptions           HeaderXFrameOptions
//	headerContentSecurityPolicy   HeaderContentSecurityPolicy
//	headerXCSRFToken              HeaderXCSRFToken
//)
//
//var (
//	HeaderAcceptType                        = reflect.TypeOf(headerAccept)
//	HeaderAcceptEncodingType                = reflect.TypeOf(headerAcceptEncoding)
//	HeaderAllowType                         = reflect.TypeOf(headerAllow)
//	HeaderAuthorizationType                 = reflect.TypeOf(headerAuthorization)
//	HeaderContentDispositionType            = reflect.TypeOf(headerContentDisposition)
//	HeaderContentEncodingType               = reflect.TypeOf(headerContentEncoding)
//	HeaderContentLengthType                 = reflect.TypeOf(headerContentLength)
//	HeaderContentTypeType                   = reflect.TypeOf(headerContentType)
//	HeaderCookieType                        = reflect.TypeOf(headerCookie)
//	HeaderSetCookieType                     = reflect.TypeOf(headerSetCookie)
//	HeaderIfModifiedSinceType               = reflect.TypeOf(headerIfModifiedSince)
//	HeaderLastModifiedType                  = reflect.TypeOf(headerLastModified)
//	HeaderLocationType                      = reflect.TypeOf(headerLocation)
//	HeaderUpgradeType                       = reflect.TypeOf(headerUpgrade)
//	HeaderVaryType                          = reflect.TypeOf(headerVary)
//	HeaderWWWAuthenticateType               = reflect.TypeOf(headerWWWAuthenticate)
//	HeaderXForwardedForType                 = reflect.TypeOf(headerXForwardedFor)
//	HeaderXForwardedProtoType               = reflect.TypeOf(headerXForwardedProto)
//	HeaderXForwardedProtocolType            = reflect.TypeOf(headerXForwardedProtocol)
//	HeaderXForwardedSslType                 = reflect.TypeOf(headerXForwardedSsl)
//	HeaderXUrlSchemeType                    = reflect.TypeOf(headerXUrlScheme)
//	HeaderXHTTPMethodOverrideType           = reflect.TypeOf(headerXHTTPMethodOverride)
//	HeaderXRealIPType                       = reflect.TypeOf(headerXRealIP)
//	HeaderXRequestIDType                    = reflect.TypeOf(headerXRequestID)
//	HeaderServerType                        = reflect.TypeOf(headerServer)
//	HeaderOriginType                        = reflect.TypeOf(headerOrigin)
//	HeaderAccessControlRequestMethodType    = reflect.TypeOf(headerAccessControlRequestMethod)
//	HeaderAccessControlRequestHeadersType   = reflect.TypeOf(headerAccessControlRequestHeaders)
//	HeaderAccessControlAllowOriginType      = reflect.TypeOf(headerAccessControlAllowOrigin)
//	HeaderAccessControlAllowMethodsType     = reflect.TypeOf(headerAccessControlAllowMethods)
//	HeaderAccessControlAllowHeadersType     = reflect.TypeOf(headerAccessControlAllowHeaders)
//	HeaderAccessControlAllowCredentialsType = reflect.TypeOf(headerAccessControlAllowCredentials)
//	HeaderAccessControlExposeHeadersType    = reflect.TypeOf(headerAccessControlExposeHeaders)
//	HeaderAccessControlMaxAgeType           = reflect.TypeOf(headerAccessControlMaxAge)
//	HeaderStrictTransportSecurityType       = reflect.TypeOf(headerStrictTransportSecurity)
//	HeaderXContentTypeOptionsType           = reflect.TypeOf(headerXContentTypeOptions)
//	HeaderXXSSProtectionType                = reflect.TypeOf(headerXXSSProtection)
//	HeaderXFrameOptionsType                 = reflect.TypeOf(headerXFrameOptions)
//	HeaderContentSecurityPolicyType         = reflect.TypeOf(headerContentSecurityPolicy)
//	HeaderXCSRFTokenType                    = reflect.TypeOf(headerXCSRFToken)
//
//	HeaderTable = map[reflect.Type]string{
//		HeaderAcceptType:                        headers.HeaderAccept,
//		HeaderAcceptEncodingType:                headers.HeaderAcceptEncoding,
//		HeaderAllowType:                         headers.HeaderAllow,
//		HeaderAuthorizationType:                 headers.HeaderAuthorization,
//		HeaderContentDispositionType:            headers.HeaderContentDisposition,
//		HeaderContentEncodingType:               headers.HeaderContentEncoding,
//		HeaderContentLengthType:                 headers.HeaderContentLength,
//		HeaderContentTypeType:                   headers.HeaderContentType,
//		HeaderCookieType:                        headers.HeaderCookie,
//		HeaderSetCookieType:                     headers.HeaderSetCookie,
//		HeaderIfModifiedSinceType:               headers.HeaderIfModifiedSince,
//		HeaderLastModifiedType:                  headers.HeaderLastModified,
//		HeaderLocationType:                      headers.HeaderLocation,
//		HeaderUpgradeType:                       headers.HeaderUpgrade,
//		HeaderVaryType:                          headers.HeaderVary,
//		HeaderWWWAuthenticateType:               headers.HeaderWWWAuthenticate,
//		HeaderXForwardedForType:                 headers.HeaderXForwardedFor,
//		HeaderXForwardedProtoType:               headers.HeaderXForwardedProto,
//		HeaderXForwardedProtocolType:            headers.HeaderXForwardedProtocol,
//		HeaderXForwardedSslType:                 headers.HeaderXForwardedSsl,
//		HeaderXUrlSchemeType:                    headers.HeaderXUrlScheme,
//		HeaderXHTTPMethodOverrideType:           headers.HeaderXHTTPMethodOverride,
//		HeaderXRealIPType:                       headers.HeaderXRealIP,
//		HeaderXRequestIDType:                    headers.HeaderXRequestID,
//		HeaderServerType:                        headers.HeaderServer,
//		HeaderOriginType:                        headers.HeaderOrigin,
//		HeaderAccessControlRequestMethodType:    headers.HeaderAccessControlRequestMethod,
//		HeaderAccessControlRequestHeadersType:   headers.HeaderAccessControlRequestHeaders,
//		HeaderAccessControlAllowOriginType:      headers.HeaderAccessControlAllowOrigin,
//		HeaderAccessControlAllowMethodsType:     headers.HeaderAccessControlAllowMethods,
//		HeaderAccessControlAllowHeadersType:     headers.HeaderAccessControlAllowHeaders,
//		HeaderAccessControlAllowCredentialsType: headers.HeaderAccessControlAllowCredentials,
//		HeaderAccessControlExposeHeadersType:    headers.HeaderAccessControlExposeHeaders,
//		HeaderAccessControlMaxAgeType:           headers.HeaderAccessControlMaxAge,
//		HeaderStrictTransportSecurityType:       headers.HeaderStrictTransportSecurity,
//		HeaderXContentTypeOptionsType:           headers.HeaderXContentTypeOptions,
//		HeaderXXSSProtectionType:                headers.HeaderXXSSProtection,
//		HeaderXFrameOptionsType:                 headers.HeaderXFrameOptions,
//		HeaderContentSecurityPolicyType:         headers.HeaderContentSecurityPolicy,
//		HeaderXCSRFTokenType:                    headers.HeaderXCSRFToken,
//	}
//)
