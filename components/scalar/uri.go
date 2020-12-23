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

func (this *Uri) UnmarshalGQL(v interface{}) error {
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

		*this = Uri(in)
	}

	return nil
}

func (this Uri) MarshalGQL(w io.Writer) {
	fmt.Fprintf(w, `"%s"`, this)
}
