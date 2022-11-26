package message

func Text(s string) Segment {
	return Segment{
		Type: TypeText,
		Data: s,
	}
}

func Record(name, url string) Segment {
	return Segment{
		Type: TypeRecord,
		Ref: Ref{
			Name: name,
			Url:  url,
		},
	}
}

func Image(name, url string) Segment {
	return Segment{
		Type: TypeImage,
		Ref: Ref{
			Name: name,
			Url:  url,
		},
	}
}

func Video(name, url string) Segment {
	return Segment{
		Type: TypeVideo,
		Ref: Ref{
			Name: name,
			Url:  url,
		},
	}
}

func File(name, url string) Segment {
	return Segment{
		Type: TypeFile,
		Ref: Ref{
			Name: name,
			Url:  url,
		},
	}
}
