golang simple live reload for development

Install
-------

	$ go get github.com/yangxikun/glr

You will now find a `glr` binary in your `$GOPATH/bin` directory.

Usage
-----

	$ glr -main app

```
Usage of ./glr:
  -args string
    	args
  -build string
    	build flags
  -delay int
    	delay *ms before rebuild (default 1000)
  -main string
    	main package name
  -wd string
    	working directory
```
