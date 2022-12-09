package entity

import "time"

// SMSPayload is the payload for sending email
type SMSPayload struct {
	To       string `json:"to"`
	Content  string `json:"content"`
	Category string `json:"category"`
}

// SMSCallback is the callback response from wave cell
type SMSCallback struct {
	Namespace   string             `json:"namespace"`
	EventType   string             `json:"eventType"`
	Description string             `json:"description"`
	Payload     SMSCallbackPayload `json:"payload"`
}

// SMSCallbackPayload is the payload callback response from wave cell
type SMSCallbackPayload struct {
	Umid            string                   `json:"umid"`
	ClientMessageID string                   `json:"clientMessageId"`
	SubAccountID    string                   `json:"subAccountId"`
	Source          string                   `json:"source"`
	Destination     string                   `json:"destination"`
	Status          SMSCallbackPayloadStatus `json:"status"`
	Price           SMSCallbackPayloadPrice  `json:"price"`
	SmsCount        int                      `json:"smsCount"`
}

// SMSCallbackPayloadStatus is the payload status callback response from wave cell
type SMSCallbackPayloadStatus struct {
	State     string    `json:"state"`
	Detail    string    `json:"detail"`
	Timestamp time.Time `json:"timestamp"`
}

// SMSCallbackPayloadPrice is the payload price callback response from wave cell
type SMSCallbackPayloadPrice struct {
	Total    int    `json:"total"`
	PerSms   int    `json:"perSms"`
	Currency string `json:"currency"`
}

// SMSBodyRequest define struct for sms body request to wave call
type SMSBodyRequest struct {
	Destination     string `json:"destination"`
	Country         string `json:"country"`
	Text            string `json:"text"`
	Source          string `json:"source"`
	ClientMessageID string `json:"clientMessageId"`
	Encoding        string `json:"encoding"`
	DlrCallbackURL  string `json:"dlrCallbackUrl"`
}

// SMSSuccessResponse define struct for sms success response from wave call
type SMSSuccessResponse struct {
	Umid            string            `json:"umid"`
	ClientMessageID string            `json:"clientMessageId"`
	Destination     string            `json:"destination"`
	Encoding        string            `json:"encoding"`
	Status          SMSStatusResponse `json:"status"`
}

// SMSStatusResponse define struct for sms status response from wave call
type SMSStatusResponse struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

// SMSErrorResponse define struct for sms error response from wave call
type SMSErrorResponse struct {
	Code      int       `json:"code"`
	Message   string    `json:"message"`
	ErrorID   string    `json:"errorId"`
	Timestamp time.Time `json:"timestamp"`
}
