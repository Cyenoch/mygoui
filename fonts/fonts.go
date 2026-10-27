// Package fonts provides the optional Cal Sans font family.
package fonts

import (
	"embed"
	"errors"
	"fmt"
	"runtime"
	"sync"

	"github.com/egoist/mygo/ui"
)

// Family is the name to use when selecting Cal Sans in MyGo.
const Family = "Cal Sans"

//go:embed CalSansVF.ttf CalSansVF-Italic.ttf OFL.txt
var assets embed.FS

var (
	registerOnce sync.Once
	registerErr  error
)

// Register adds Cal Sans and its italic face to MyGo. Call it before drawing
// text with Family. Calls after the first return the same result.
//
// If the italic registration fails, the regular face remains registered.
// Browser builds return an error because MyGo cannot register these faces
// with the correct styles or report font-loading failures.
func Register() error {
	registerOnce.Do(func() {
		if runtime.GOOS == "js" {
			registerErr = errors.New("fonts: MyGo browser font registration cannot select Cal Sans styles or report load failures")
			return
		}
		for _, name := range []string{"CalSansVF.ttf", "CalSansVF-Italic.ttf"} {
			data, err := assets.ReadFile(name)
			if err != nil {
				registerErr = fmt.Errorf("fonts: read %s: %w", name, err)
				return
			}
			if err := ui.RegisterFont(data, Family); err != nil {
				registerErr = fmt.Errorf("fonts: register %s: %w", name, err)
				return
			}
		}
	})
	return registerErr
}
