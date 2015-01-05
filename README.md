Usage
-----

`configura.LoadEnv` will go through all the fields defined in a struct and load their values from system environment variables.

The environment variable name can set using struct tags. A default value can also be optionally set. For example:

```
var config = struct {
  LogPrefix   string `env:"LOG_PREFIX"`
  Port        int    `env:"PORT,8888"`
  Development bool   `env:"DEVELOPMENT"`
}{}

err = configura.LoadEnv(&config)
```


Background
----------

The Twelve-Factor methodology [recommends](http://12factor.net/config) recommends storing application configuration in system environment variables. On the other hand, environment variables don't provide an [explicit and fully-documented surface area] for an application, like command-line flags do.

Configura achieves both these goals by documenting the environment variables used by an application, allowing default values to be specified, and – in the words of the [original author](https://github.com/agonzalezro/configura) – keeping all this configuration-loading easy-peasy.
