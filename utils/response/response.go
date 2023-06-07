package response

import "strings"

const (
	// 4xx
	NoTokenFound                           = 400001
	BearerTokenNotInProperFormat           = 400002
	TokenInvalid                           = 400003
	EmailAlreadyExists                     = 400004
	EmailNotExists                         = 400005
	FailedToGetStateToken                  = 400006
	IdInvalid                              = 400007
	ImageFileNameLimitOf100                = 400008
	ImageFileSizeLimitOf5MB                = 400009
	TokenDoesNotExistOrExpired             = 401001
	InvalidCredential                      = 401002
	TokenContainsAnInvalidNumberOfSegments = 401003
	FailedToLogout                         = 401004
	RecordNotFound                         = 401005
	TooManyRequests                        = 429001
	TokenBindingHasUnknownErrors           = 500001
	DataBindingHasUnknownErrors            = 500002
	DuplicateCreatedData                   = 500003
	DuplicatedTitle                        = 500004
)

var (
	Messages = map[int]string{
		// 4xx
		400001: "No token found.",
		400002: "Bearer token not in proper format.",
		400003: "Token invalid.",
		400004: "Email already exists.",
		400005: "Email not exists.",
		400006: "Failed to get state token.",
		400007: "ID Invalid.",
		400008: "Image file name limit of 100.",
		400009: "Image file size limit of 5 MB.",
		401001: "Token does not exist or expired.",
		401002: "Invalid credential.",
		401003: "Token contains an invalid number of segments.",
		401004: "Failed to logout.",
		401005: "Record not found.",
		429001: "Too many requests.",
		500001: "Token binding has unknown errors.",
		500002: "Data binding has unknown errors.",
		500003: "Duplicate created data.",
		500004: "Duplicated title.",
	}
)

// Create a new struct for the response data
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
	Data    interface{} `json:"data"`
}

// Create a new struct for the page response data
type PageResponse struct {
	Code        int         `json:"code"`
	Message     string      `json:"message"`
	CurrentPage int64       `json:"currentPage"`
	PageLimit   int64       `json:"pageLimit"`
	Total       int64       `json:"total"` // Data count
	Pages       int64       `json:"pages"` // Total page
	Errors      interface{} `json:"errors"`
	Data        interface{} `json:"data"`
}

// EmptyObj object is used when data doesnt want to be null on json
// type EmptyObject struct {
// }

// SuccessResponse returns a success response with the given data
func SuccessResponse(code int, message string, data interface{}) Response {
	return Response{
		Code:    code,
		Message: message,
		Errors:  nil,
		Data:    data,
	}
}

func SuccessPageResponse(code int, message string, currentPage int64, pageLimit int64, total int64, pages int64, data interface{}) PageResponse {
	return PageResponse{
		Code:        code,
		Message:     message,
		CurrentPage: currentPage,
		PageLimit:   pageLimit,
		Total:       total,
		Pages:       pages,
		Errors:      nil,
		Data:        data,
	}
}

// ErrorResponse returns an error response with the given data
func ErrorsResponse(code int, message string, err string, data interface{}) Response {
	splittedError := strings.Split(err, "\n")
	return Response{
		Code:    code,
		Message: message,
		Errors:  splittedError,
		Data:    data,
	}
}

// ErrorResponse Returns an error response with the message given by the code
func ErrorsResponseByCode(code int, message string, errCode int, data interface{}) Response {
	splittedError := strings.Split(Messages[errCode], "\n")
	return Response{
		Code:    code,
		Message: message,
		Errors:  splittedError,
		Data:    data,
	}
}
