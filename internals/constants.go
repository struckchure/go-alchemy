package internals

import (
	"os"

	"github.com/samber/lo"
)

var baseUrlFromEnv string = os.Getenv("ALCHEMY_BASE_DIR")
var defaultBaseUrl string = "https://raw.githubusercontent.com/struckchure/go-alchemy/refs/heads/main/"

var ALCHEMY_BASE_DIR string = lo.Ternary(baseUrlFromEnv != "", baseUrlFromEnv, defaultBaseUrl)
