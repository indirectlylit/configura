// configura.LoadEnv will go through all the fields
// defined in the struct and load their values from
// system environment variables.
//
// The variable name can set using struct tags.
// A default value can also be optionally set.
//
// For example:
//
// var config = struct {
// 	LogPrefix   string `env:"LOG_PREFIX"`
// 	Port        int    `env:"PORT,8888"`
// 	Development bool   `env:"DEVELOPMENT"`
// }{}
//
// err = LoadEnv(&config)
package configura
