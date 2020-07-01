package scalar

type ContentType string

const (
	// reference https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types/Common_types
	ImagePNG          ContentType = "image/png"
	ImageJpeg         ContentType = "image/jpeg"
	ImageGif          ContentType = "image/gif"
	ApplicationJson   ContentType = "application/json"
	ApplicationJsonLD ContentType = "application/ld+json"
	ApplicationPdf    ContentType = "application/pdf"
	ApplicationGzip   ContentType = "application/gzip"
	ApplicationBzip   ContentType = "application/x-bzip"
	ApplicationZip    ContentType = "application/zip"
	AudioMpeg         ContentType = "audio/mpeg"
	VideoMpeg         ContentType = "video/mpeg"
	TextPlain         ContentType = "text/plain"
	TextCsv           ContentType = "text/csv"
)
