# Redactrus <img src="http://i.imgur.com/nHsZvo9.png" width="40" height="40" alt=":walrus:" class="emoji" title=":walrus:"/>&nbsp; [![Build Status](https://travis-ci.org/whuang8/redactrus.svg?branch=master)](https://travis-ci.org/whuang8/redactrus)&nbsp;[![Coverage Status](https://coveralls.io/repos/github/whuang8/redactrus/badge.svg)](https://coveralls.io/github/whuang8/redactrus)&nbsp;[![GoDoc](https://godoc.org/github.com/whuang8/redactrus?status.svg)](https://godoc.org/github.com/whuang8/redactrus)

Redactrus is a [Logrus hook](https://github.com/sirupsen/logrus#hooks) that redacts specified text from your logs.

Easy redaction of log data:

![Redacted](https://i.imgur.com/7bWHxKq.png)

##### Example

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

### Usage
You can create a Redactrus hook by initializing the `redactrus.Hook` struct with defined `LogLevels` and `RedactionList` fields. Then, simply use the [`logrus.AddHook`](https://godoc.org/github.com/sirupsen/logrus#AddHook) function to attach the hook to your logrus logger. Refer to the example above to see this in action.

### RedactionList
The `RedactionList` is a slice of strings that defines the patterns you want to redact from your logs. There are many ways to define a string in the `RedactionList` to allow for flexibility in the patterns of text you want to redact.

##### Explicit Words
By defining words explicitly in the `RedactionList` , Redactrus will try to find a key that is the same as the defined word and redact the associated value in the log entry’s data fields.

```go
rh := &redactrus.Hook{
  AcceptedLevels: log.AllLevels,
  RedactionList: []string{"word"},
}

// ...

log.WithFields(log.Fields{
  "word": "bird",
}).Info("Bah bah bah bird bird bird")
```

![explicit1](https://i.imgur.com/nLPFBAt.png)

Redactrus will also redact any occurrences of the word in the log entry’s message, as well as the values in the log entry’s data fields.

```go
rh := &redactrus.Hook{
  AcceptedLevels: log.AllLevels,
  RedactionList: []string{"bird"},
}

// ...

log.WithFields(log.Fields{
  "song": "surfin' bird",
}).Info("A-well-a bird bird bird, well-a bird is the word")
```

![explicit2](https://i.imgur.com/pjUfsGI.png)

##### Regular Expressions
You can also define [regular expressions](https://gobyexample.com/regular-expressions) in the `RedactionList` to match a specific pattern instead of an explicit word.
```go
rh := &redactrus.Hook{
  AcceptedLevels: log.AllLevels,
  RedactionList: []string{"[0-9]{3}-[0-9]{2}-[0-9]{4}"},
}

// ...

log.WithFields(log.Fields{
  "ssn": "111-22-3333",
}).Info("A new customer with ssn: 111-22-3333 has been registered")
```

![regex1](https://i.imgur.com/YXMw28u.png)

Sometimes, redacting all of the text that matches a regex pattern is not wanted. If a logrus logger logs the following string: `/api/v1/endpoint?ssn=111223333&account_id=123456789&pin_number=987654321`, how do we only redact the SSN number and yield `ssn=[REDACTED]&`? Additionally, how can we avoid redacting the other 9-digit numbers?

For sections of a regular expression that you **do not** want to redact, simply wrap it in **parenthesis**. Redactrus will know to match on the entire expression provided, but will not redact the parts of the text that match the regex defined in parenthesis.

```go
rh := &redactrus.Hook{
  AcceptedLevels: log.AllLevels,
  RedactionList: []string{"(ssn=)[0-9]{9}(&)"},
}

// ...

log.WithFields(log.Fields{
  "path": "/api/v1/endpoint?ssn=111223333&account_id=123456789&pin_number=987654321",
}).Info("Request received")
```
![regex2](https://i.imgur.com/8u9CWh3.png)
