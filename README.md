gocritic
========

gocritic is a [Go](https://golang.org) library for [CriticMarkup](http://criticmarkup.com).

**Note**: This project is a library only, if you are looking for a CLI tool, go to [gocritic-cli](https://github.com/dohzya/gocritic-cli).

Use
---

The main function is `gocritic.Critic`, here is its full signature:

```go
func gocritic.Critic(w io.Writer, r io.Reader, fopts ...func(*Options)) (int, error)
```

The `fopts` argument is used to specify options

A typical use is:

```go
if _, err := gocritic.Critic(out, int); err != nil {
  // handle the error
}
```

To generate only the original sources:

```go
if _, err := gocritic.Critic(out, int, gocritic.FilterOnlyOriginal); err != nil {
  // handle the error
}
```
