package scalar

import (
	"fmt"
	"io"
	"net/url"
)

func NilUri(in Uri) *Uri {
	return &in
}

type Uri string

func (value *Uri) UnmarshalGQL(v interface{}) error {
	if in, ok := v.(string); !ok {
		return fmt.Errorf("URI must be strings")
	} else {
		uri, err := url.ParseRequestURI(in)
		if nil != err {
			return fmt.Errorf("invalid URI for request")
		}

		if uri.RequestURI() == "" {
			return fmt.Errorf("missing request URI")
		}

		*value = Uri(in)
	}

	return nil
}

func (value Uri) MarshalGQL(w io.Writer) {
	fmt.Fprintf(w, `"%s"`, value)
}
