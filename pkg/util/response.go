package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"

	"strconv"

	"github.com/go-playground/validator/v10"

	"math"
	"strings"

	"github.com/google/uuid"
)

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type Meta struct {
	Message string          `json:"message"`
	Code    int             `json:"code"`
	Status  string          `json:"status"`
	Errors  []ErrorResponse `json:"errors,omitempty"`
}

type ErrorResponse struct {
	Key     string `json:"key,omitempty"`
	Message string `json:"message,omitempty"`
}

func (r *Response) WithError(errors ...ErrorResponse) Response {
	r.Meta.Errors = errors
	return *r
}

func FormatValidationError(err error) []string {
	var dataErrror []string
	var foo *json.UnmarshalTypeError
	if errors.As(err, &foo) {
		dataErrror = append(dataErrror, err.Error())
		return dataErrror
	}
	for _, e := range err.(validator.ValidationErrors) {
		dataErrror = append(dataErrror, e.Error())
	}

	return dataErrror
}

func APIResponse(message string, code int, status string, data interface{}) Response {
	meta := Meta{
		Message: message,
		Code:    code,
		Status:  status,
	}

	jsonResponse := Response{
		Meta: meta,
		Data: data,
	}

	return jsonResponse

}

// Helper function to get file extension
func getFileExtension(filename string) string {
	dotIndex := strings.LastIndex(filename, ".")
	if dotIndex == -1 {
		return ""
	}
	return filename[dotIndex:]
}

// Helper function to check if a slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func AddValueToTemplateString(templateText string, value any) (string, error) {
	temp, err := template.New(uuid.NewString()).Parse(templateText)
	if err != nil {
		return "", err
	}

	buff := new(bytes.Buffer)

	err = temp.Execute(buff, value)
	if err != nil {
		return "", err
	}

	return buff.String(), nil
}

func RoundUp(number float64) float64 {
	precision := 1.0
	// Shift decimal point to the right
	shifted := number * math.Pow(10, precision)

	// Apply ceiling
	rounded := math.Ceil(shifted)

	// Shift decimal point back to the left
	final := rounded / math.Pow(10, precision)
	return final
}

func RoundUpToNearestThousand(value int) int {
	return int(math.Ceil(float64(value)/1000) * 1000)
}

func TruncateToOneDecimalPlace(value float64) float64 {
	return math.Trunc(value*10) / 10
}

func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func GetFbTraceId(errString string) string {
	parts := strings.Split(errString, "&")

	var fbTraceID string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "fbtrace_id:") {
			fbTraceID = strings.TrimSpace(strings.TrimPrefix(part, "fbtrace_id:"))
			break
		}
	}

	return "trace_id:" + fbTraceID
}

// FormatIntWithCommas formats an integer to a string with comma as thousand separators
func FormatIntWithCommas(n int) string {
	s := strconv.Itoa(n)
	if len(s) <= 3 {
		return s
	}

	var result strings.Builder
	length := len(s)
	for i, c := range s {
		if i > 0 && (length-i)%3 == 0 {
			result.WriteRune('.')
		}
		result.WriteRune(c)
	}

	return result.String()
}

//
//func GetErrorResponse(err error) (ErrorResponse, int) {
//	var errMsg ErrorResponse
//	var code int
//
//	switch err {
//	case INSUFFICENTBALANCE:
//		errMsg = ErrorResponse{
//			Key:     "insufficient_balance",
//			Message: "Top up your balance first, then you can make changes right away!",
//		}
//		code = http.StatusBadRequest
//
//	case ADNOTACTIVE:
//		errMsg = ErrorResponse{
//			Key:     "ad_not_active",
//			Message: "Activate your ads first, then you can make changes right away!",
//		}
//		code = http.StatusBadRequest
//
//	case WHATSAPPNNOTLINKED:
//		errMsg = ErrorResponse{
//			Key:     "whatsapp_not_linked",
//			Message: "Nomor WhatsApp tidak terhubung ke Halaman Facebook. Ganti nomor dan terbitkan kembali.",
//		}
//		code = http.StatusBadRequest
//
//	default:
//		errMsg = ErrorResponse{
//			Key:     "server_error",
//			Message: "Oops, there was a problem on the server. Try again later, or contact our CS if you still have problems.",
//		}
//		code = http.StatusInternalServerError
//	}
//
//	return errMsg, code
//}
