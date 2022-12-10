package atmt

func Error(resp *Message, code StatusCode) {
	b := MessageBuilder[DataResponse]{
		Type: MessageResponse,
		Data: DataResponse{
			StatusCode: code,
		},
	}
	_ = b.Write(resp)
}

func BadRequest(resp *Message) {
	Error(resp, StatusBadRequest)
}

func InternalServerError(resp *Message) {
	Error(resp, StatusInternalServerError)
}
