package td

import (
	"context"
)

type Provider interface {
	Resources(context.Context) []Resource
}

type Resource interface {
}
