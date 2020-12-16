package api

import (
	"fmt"
	"io"
	"strconv"
)

type FileType string

const (
	FileTypePdf       FileType = "PDF"
	FileTypeTxt       FileType = "TXT"
	FileTypeJpeg      FileType = "JPEG"
	FileTypePng       FileType = "PNG"
	FileTypeMp3       FileType = "MP3"
	FileTypeMp4       FileType = "MP4"
	FileTypeWebmAudio FileType = "WEBM_AUDIO"
	FileTypeWebmVideo FileType = "WEBM_VIDEO"
	FileTypeZip       FileType = "ZIP"
	FileTypeGzip      FileType = "GZIP"
)

var AllFileType = []FileType{
	FileTypePdf,
	FileTypeTxt,
	FileTypeJpeg,
	FileTypePng,
	FileTypeMp3,
	FileTypeMp4,
	FileTypeWebmAudio,
	FileTypeWebmVideo,
	FileTypeZip,
	FileTypeGzip,
}

func (e FileType) IsValid() bool {
	switch e {
	case FileTypePdf, FileTypeTxt, FileTypeJpeg, FileTypePng, FileTypeMp3, FileTypeMp4, FileTypeWebmAudio, FileTypeWebmVideo, FileTypeZip, FileTypeGzip:
		return true
	}
	return false
}

func (e FileType) String() string {
	return string(e)
}

func (e *FileType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = FileType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid FileType", str)
	}
	return nil
}

func (e FileType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
