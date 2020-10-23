package goarecovery

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	goa "goa.design/goa/v3/pkg"
)

// Recover is a middleware that recovers panics and maps them to errors.
func Recover() func(goa.Endpoint) goa.Endpoint {
	return func(e goa.Endpoint) goa.Endpoint {
		return goa.Endpoint(func(ctx context.Context, req interface{}) (res interface{}, err error) {
			defer func() {
				if r := recover(); r != nil {
					var msg string
					switch x := r.(type) {
					case string:
						msg = fmt.Sprintf("panic: %s", x)
					case error:
						msg = fmt.Sprintf("panic: %s", x)
					default:
						msg = "unknown panic"
					}
					const size = 64 << 10 // 64KB
					buf := make([]byte, size)
					buf = buf[:runtime.Stack(buf, false)]
					lines := strings.Split(string(buf), "\n")
					stack := lines[3:]
					err = fmt.Errorf("%s\n%s", msg, strings.Join(stack, "\n"))
				}
			}()
			return e(ctx, req)
		})
	}
}
