# Redactrus <img src="http://i.imgur.com/nHsZvo9.png" width="40" height="40" alt=":walrus:" class="emoji" title=":walrus:"/>&nbsp; [![Build Status](https://travis-ci.org/whuang8/redactrus.svg?branch=master)](https://travis-ci.org/whuang8/redactrus)&nbsp;[![Coverage Status](https://coveralls.io/repos/github/whuang8/redactrus/badge.svg)](https://coveralls.io/github/whuang8/redactrus)&nbsp;[![GoDoc](https://godoc.org/github.com/whuang8/redactrus?status.svg)](https://godoc.org/github.com/whuang8/redactrus)

Redactrus is a [Logrus hook](https://github.com/sirupsen/logrus#hooks) that redacts specified text from your logs.

Easy redaction of log data:

![Redacted](https://i.imgur.com/7bWHxKq.png)

#### Example

To use Redactrus, simply add the Redactrus hook to your Logrus logger:

```go
package main

import (
  log "github.com/sirupsen/logrus"
  "github.com/whuang8/redactrus"
)

func init() {
  // Create Redactrus hook that is triggered
  // for every logger level and redacts any
  // Logrus fields with the key as 'password'
  rh := &redactrus.Hook{
      AcceptedLevels: log.AllLevels,
      RedactionList: []string{"password"},
  }
  
  log.AddHook(rh)
}

func main() {
  log.WithFields(log.Fields{
    "walrusName": "Walrein",
    "password": "iloveclams<3",
  }).Info("A walrus attempted to log in.")
}
```

![example1](https://i.imgur.com/4LwOMr2.png)
