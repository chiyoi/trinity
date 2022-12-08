package atmt

func Error(resp *Message, code StatusCode) (err error) {
	*resp, err = (&MessageBuilder[DataResponse]{
		Type: MessageResponse,
		Data: DataResponse{
			StatusCode: code,
		},
	}).Message()
	return
}
