package Utils

import "fmt"

func URLJoin(s ...string) string {
	r := ""

	for _, v := range s {
		d := v
		if v[0] == '/' {
			d = d[1:]
		}

		if v[len(v)-1] == '/' {
			d = d[:len(v)-1]
		}

		if r == "" {
			r = d
			continue
		}
		r = fmt.Sprintf("%s/%s", r, d)
	}

	return r
}
