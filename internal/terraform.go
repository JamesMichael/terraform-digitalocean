package td

import (
	"io"
	"regexp"
)

type Terraformer interface {
	TerraformProvider(io.Writer) error
	TerraformResource(io.Writer, Resource) error
	TerraformResources(io.Writer, []Resource) error
	TerraformImports(io.Writer, []Resource) error
}

func TerraformResourceName(original string) (string, error) {
	r, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return "", err
	}

	filtered := r.ReplaceAllString(original, "_")
	return filtered, nil
}
