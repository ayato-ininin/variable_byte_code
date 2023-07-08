package formatByte

import (
	"fmt"
	"strings"
)

func FormatBytes(data []byte) string {
	var builder strings.Builder
	builder.WriteString("[]byte{")
	for i, b := range data {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(fmt.Sprintf("0x%02x", b))
	}
	builder.WriteString("}")
	return builder.String()
}
