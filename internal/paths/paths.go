package paths

import (
	"fmt"

	"github.com/adrg/xdg"
	"github.com/reynn/notifier/internal/constants"
)

func ConfigFile(ext string) string {
	configFile, err := xdg.ConfigFile(fmt.Sprintf("%s/%s.config.%s", constants.AppName, constants.AppModule, ext))
	if err != nil {
		panic(err)
	}
	return configFile
}
