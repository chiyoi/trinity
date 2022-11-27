package message

func Text(txt string) Segment {
	return Segment{
		Type: TypeText,
		Data: map[string]string{
			"text": txt,
		},
	}
}

func Image(url string) Segment {
	return Segment{
		Type: TypeImage,
		Data: map[string]string{
			"file": url,
		},
	}
}

func Record(url string) Segment {
	return Segment{
		Type: TypeRecord,
		Data: map[string]string{
			"file": url,
		},
	}
}
