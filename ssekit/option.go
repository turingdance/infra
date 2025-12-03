package ssekit

type SSEOption func(*SseCtrl)

func Base64Wraper() SSEOption {
	return func(sc *SseCtrl) {
		sc.b64wraper = true
	}
}

func EncodeWith(fmt DataFormat) SSEOption {
	return func(sc *SseCtrl) {
		sc.format = fmt
	}
}
