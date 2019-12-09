# koron-go/uneventicker

[![GoDoc](https://godoc.org/github.com/koron-go/uneventicker?status.svg)](https://godoc.org/github.com/koron-go/uneventicker)
[![Actions/Go](https://github.com/koron-go/uneventicker/workflows/Go/badge.svg)](https://github.com/koron-go/uneventicker/actions?query=workflow%3AGo)
[![CircleCI](https://img.shields.io/circleci/project/github/koron-go/uneventicker/master.svg)](https://circleci.com/gh/koron-go/uneventicker/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/koron-go/uneventicker)](https://goreportcard.com/report/github.com/koron-go/uneventicker)

Ticker with uneven durations.

## Getting started

This fires every 100msec.

```go
t := uneventicker.New(100*time.Millisecond)
for {
    println(<-t.C)
}
```

This fires after 50msec at first, then fires every 100msec.

```go
t := uneventicker.New(50*time.Millisecond, 100*time.Millisecond)
for {
    println(<-t.C)
}
```

This fires after 50msec at first, then after 200msec at second, then fires
every 75msec.

```go
t := uneventicker.New(50*time.Millisecond, 200*time.Millisecond, 75*time.Millisecond)
for {
    println(<-t.C)
}
```
