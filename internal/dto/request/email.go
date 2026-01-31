package request

type Attachment struct {
	FileName    string
	FileByte    []byte
	ContentType string
}

type EmailRequest struct {
	From        string
	To          string
	Subject     string
	Body        string
	Attachments []Attachment
}