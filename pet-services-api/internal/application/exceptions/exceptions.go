package exceptions

import "github.com/gin-gonic/gin"

const (
	RFC200_CODE = 200
	RFC201_CODE = 201
	RFC204_CODE = 204
	RFC400_CODE = 400
	RFC401_CODE = 401
	RFC403_CODE = 403
	RFC404_CODE = 404
	RFC409_CODE = 409
	RFC429_CODE = 429
	RFC500_CODE = 500
	RFC503_CODE = 503
)

const (
	RFC200 = "https://datatracker.ietf.org/doc/html/rfc7231#section-6.3.1"
	RFC201 = "https://datatracker.ietf.org/doc/html/rfc7231#section-6.3.2"
	RFC204 = "https://datatracker.ietf.org/doc/html/rfc7231#section-6.3.5"
	RFC400 = "https://datatracker.ietf.org/doc/html/rfc7231#section-6.5.1"
	RFC401 = "https://datatracker.ietf.org/doc/html/rfc7235#section-3.1"
	RFC403 = "https://datatracker.ietf.org/doc/html/rfc7231#section-6.5.3"
	RFC404 = "https://datatracker.ietf.org/doc/html/rfc7231#section-6.5.4"
	RFC409 = "https://datatracker.ietf.org/doc/html/rfc7231#section-6.5.8"
	RFC422 = "https://datatracker.ietf.org/doc/html/rfc4918#section-11.2"
	RFC429 = "https://datatracker.ietf.org/doc/html/rfc6585#section-4"
	RFC500 = "https://datatracker.ietf.org/doc/html/rfc7231#section-6.6.1"
	RFC503 = "https://datatracker.ietf.org/doc/html/rfc7231#section-6.6.4"
)

type ErrorType string

const (
	BadRequest          ErrorType = "Bad Request"
	Unauthorized        ErrorType = "Unauthorized"
	Forbidden           ErrorType = "Forbidden"
	NotFound            ErrorType = "Not Found"
	Conflict            ErrorType = "Conflict"
	UnprocessableEntity ErrorType = "Unprocessable Entity"
	TooManyRequests     ErrorType = "Too Many Requests"
	InternalServerError ErrorType = "Internal Server Error"
	ServiceUnavailable  ErrorType = "Service Unavailable"
)

type ErrorMetadata struct {
	Status   int
	Instance string
}

var errorMetadataMap = map[ErrorType]ErrorMetadata{
	BadRequest:          {400, RFC400},
	Unauthorized:        {401, RFC401},
	Forbidden:           {403, RFC403},
	NotFound:            {404, RFC404},
	Conflict:            {409, RFC409},
	UnprocessableEntity: {422, RFC422},
	TooManyRequests:     {429, RFC429},
	InternalServerError: {500, RFC500},
	ServiceUnavailable:  {503, RFC503},
}

func (e ErrorType) Metadata() ErrorMetadata {
	if meta, ok := errorMetadataMap[e]; ok {
		return meta
	}
	return errorMetadataMap[InternalServerError]
}

type ErrorMessage struct {
	Title  string
	Detail string
}

type ProblemDetails struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance,omitempty"`
}

type ProblemDetailsOutputDTO struct {
	ProblemDetails []ProblemDetails `json:"problem_details"`
}

func NewProblemDetails(errorType ErrorType, msg ErrorMessage) ProblemDetails {
	meta := errorType.Metadata()

	return ProblemDetails{
		Type:     string(errorType),
		Title:    msg.Title,
		Status:   meta.Status,
		Detail:   msg.Detail,
		Instance: meta.Instance,
	}
}

func HandleErrors(c *gin.Context, errs []ProblemDetails) {
	if len(errs) == 0 {
		return
	}

	status := errs[0].Status
	c.JSON(status, gin.H{"errors": errs})
}
