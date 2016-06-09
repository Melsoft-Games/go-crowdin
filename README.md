# go-crowdin
Crowdin API in Go - https://crowdin.com/page/api

> Note: In progress of supporting all API endpoints

#### Install

`go get github.com/medisafe/go-crowdin`

- [Initialize](#initialize)
- [API](#api)
- [Debug](#debug)
- [App Engine](#app-engine)

##### Initialize

``` Go
api := crowdin.New("token", "project-name")
```

##### API

:blue_book: Check the doc - [Documentation](https://godoc.org/github.com/medisafe/go-crowdin)

> Example:
``` Go
// get language status
files, err := api.GetLanguageStatus("ru")
```

##### Debug

You can print the internal errors by enabling debug to true

``` Go
api.SetDebug(true, nil)
```

You can also define your own `io.Writer` in case you want to persist the logs somewhere.
For example keeping the errors on file

``` Go
logFile, err := os.Create("crowdin.log")
api.SetDebug(true, logFile)
```

##### App Engine

Initialize app engine client and continue as usual

``` Go
c := appengine.NewContext(r)
client := urlfetch.Client(c)

api := crowdin.New("token", "project-name")
api.SetClient(client)
```

[Documentation](https://godoc.org/github.com/medisafe/go-crowdin)

##### Author

Roman Kushnarenko - [sromku](https://github.com/sromku)

##### License 

Apache License 2.0