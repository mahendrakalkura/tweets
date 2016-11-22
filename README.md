How to install?
===============

```
$ mkdir tweets
$ cd tweets
$ git clone --recursive git@bitbucket.org:tweettv/tweets.git .
$ cp --archive settings.toml.sample settings.toml
$ go get
```

How to run?
===========

```
$ cd tweets
$ go run tweets.go
```
