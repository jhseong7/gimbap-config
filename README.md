# GIMBAP Config provider

This is a simple configuration provider to use with GIMBAP. It reads configurations from a file and provides them to GIMBAP.

This provider uses the famous Viper under the hood so any configurations that is compatitable with Viper can be used with this provider.

## Installation

```bash
go get github.com/jhseong7/gimbap-config
```

## Usage

import the `ConfigModule` provided with this package to the GIMBAP application you intend to use.

```go
import (
  "github.com/jhseong7/gimbap"
  config "github.com/jhseong7/gimbap-config"
)

var AppModule = gimbap.DefineModule(gimbap.ModuleConfig{
  Name: "AppModule",
  SubModules: []*gimbap.Module{
    config.ConfigModule("config.yaml"),
  },
  // other configurations
})
```

Then you must inject the default options provider to the appmodule.

```go
func main() {
	app := gimbap.CreateApp(gimbap.AppOption{
		AppModule: app.AppModule,
	})

	// Provider the root config for the config service
	app.UseInjection(config.ConfigOption{ConfigFilePathList: []string{"./.env"}})

  app.Run()
}
```

then inject the `ConfigService` to the module you want to use the configurations.

## Third-party libraries

- **Library Name:** spf13/viper
  - **Purpose:** Used for reading configurations
  - **License:** MIT License. [Link](https://github.com/spf13/viper/blob/master/LICENSE)
