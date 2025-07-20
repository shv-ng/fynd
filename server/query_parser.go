package server

import (
	"strconv"
	"strings"

	"github.com/ShivangSrivastava/fynd/app"
	"github.com/ShivangSrivastava/fynd/indexer"
)

type QueryOptions struct {
	Top   int
	Ext   []string
	Query []string
}

func ParseQuery(ctx app.Context, input string, opts QueryOptions) QueryOptions {
	input = strings.ToLower(input)
	var remaining []string
	if opts.Top == 0 {
		opts.Top = ctx.Setting.Top
	}
	if opts.Top == 0 {
		opts.Top = -1
	}
	parts := strings.SplitSeq(input, ";")
	for part := range parts {
		part = strings.TrimSpace(part)
		switch {
		case strings.HasPrefix(part, "top:"):
			val := strings.TrimPrefix(part, "top:")
			if n, err := strconv.Atoi(val); err == nil {
				opts.Top = n
			}
		case strings.HasPrefix(part, "ext:"):
			val := strings.TrimPrefix(part, "ext:")
			opts.Ext = strings.Split(val, ",")
		default:
			remaining = append(remaining, part)
		}
	}

	// sanatize query
	opts.Query = indexer.Sanatize(strings.Join(remaining, "_"))
	return opts
}
