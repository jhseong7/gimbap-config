package config

import (
	"github.com/jhseong7/gimbap"
)

var ConfigModule = gimbap.DefineModule(gimbap.ModuleOption{
	Name: "ConfigModule",
	Providers: []*gimbap.Provider{
		ConfigServiceProvider,
	},
})
