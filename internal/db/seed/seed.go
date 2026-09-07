package seed

import (
	_ "embed"
	"encoding/json"

	"github.com/boracomet/ai-resume-builder/internal/models"
)

//go:embed example-profile.json
var exampleProfileJSON []byte

func ExampleProfile() *models.CVProfile {
	var profile models.CVProfile
	if err := json.Unmarshal(exampleProfileJSON, &profile); err != nil {
		panic("seed profile: " + err.Error())
	}
	profile.Normalize()
	return &profile
}
