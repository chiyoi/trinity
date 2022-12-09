package atmt

func Error(resp *Message, code StatusCode) (err error) {
	b := MessageBuilder[DataResponse]{
		Type: MessageResponse,
		Data: DataResponse{
			StatusCode: code,
		},
	}
	return b.Write(resp)
}
