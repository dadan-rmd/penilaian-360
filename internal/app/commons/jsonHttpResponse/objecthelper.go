package jsonHttpResponse

import "github.com/gin-gonic/gin"

//type helpers. Used as an example value in swagger documentation

const (
	SuccessStatus   = "SUCCESS"
	FailedStatus    = "FAILED"
	DefaultPassword = "DLMS1234"
)

type SuccessResponse struct {
	Status       string      `json:"status" example:"success"`
	ResponseCode string      `json:"response_code" example:"98"`
	Message      string      `json:"message" example:"process done successfully"`
	Data         interface{} `json:"data"`
}

type FailedResponse struct {
	Status       string      `json:"status" example:"failed"`
	ResponseCode string      `json:"response_code" example:"98"`
	Message      interface{} `json:"message"`
}

//responses helper
type FailedUnauthorizedResponse struct {
	Status  string `json:"status" example:"failed"`
	Message string `json:"message" example:"invalid user token"`
}

type FailedBadRequestResponse struct {
	Status  string `json:"status" example:"failed"`
	Message string `json:"message" example:"bad request"`
}

type FailedUnprocessableEntityResponse struct {
	Status  string `json:"status" example:"failed"`
	Message string `json:"message" example:"unprocessable entity"`
}

type FailedInternalServerErrorResponse struct {
	Status  string `json:"status" example:"failed"`
	Message string `json:"message" example:"internal server error"`
}

type FailedNotFoundResponse struct {
	Status  string `json:"status" example:"failed"`
	Message string `json:"message" example:"data not found"`
}

type ConflictResponse struct {
	Status  string `json:"status" example:"failed"`
	Message string `json:"message" example:"conflict"`
}

func NewSuccessResponse(payload interface{}) SuccessResponse {
	return SuccessResponse{Status: "success", Data: payload}
}

func NewFailedResponse(message interface{}) FailedResponse {
	return FailedResponse{Status: "failed", Message: message}
}

func NewSuccessfulOKResponse(c *gin.Context, payload interface{}) {
	OK(c, NewSuccessResponse(payload))
}

func NewFailedBadRequestResponse(c *gin.Context, message interface{}) {
	BadRequest(c, NewFailedResponse(message))
}

func NewNotFoundResponse(c *gin.Context, message interface{}) {
	NotFound(c, NewFailedResponse(message))
}

func NewFailedUnauthorizedResponse(c *gin.Context, message interface{}) {
	Unauthorized(c, NewFailedResponse(message))
}

func NewFailedUnprocessableResponse(c *gin.Context, message interface{}) {
	Unprocessable(c, NewFailedResponse(message))
}

func NewFailedInternalServerResponse(c *gin.Context, message interface{}) {
	InternalServerError(c, NewFailedResponse(message))
}

func NewFailedConflictResponse(c *gin.Context, message interface{}) {
	Conflict(c, NewFailedResponse(message))
}

func NewTooManyRequest(c *gin.Context, message interface{}) {
	TooManyRequest(c, NewFailedResponse(message))
}

type MissingRequiredField struct {
	Message string      `json:"message"`
	Fields  interface{} `json:"fields"`
}

func NewFailedMissingRequiredFieldResponse(c *gin.Context, message interface{}) {
	missingFieldMessage := MissingRequiredField{Message: "please input the required field", Fields: message}
	BadRequest(c, NewFailedResponse(missingFieldMessage))
}
