# logger

Abstraction layer that takes care of some things from [logrus-papertrail-hook](https://github.com/polds/logrus-papertrail-hook) and [logrus](https://github.com/Sirupsen/logrus)

## Usage

```
package main

import (
    "gopkg.in/kyani-inc/logger.v2"
)

func main() {
    log := logger.New(logger.Config{
        Appname: "test",
        Host:    "logs.papertrailapp.com",
        Port:    "1337",
    })

    log.Info("Hello")
}
```