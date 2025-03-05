# Render Template
A simple program to render go templates using one or more value files.   [Sprig](https://masterminds.github.io/sprig/)
and other customer functions are included.  See [golang template](https://pkg.go.dev/text/template#pkg-overview)
documentation for template syntax help.  

## Build

```
go mod tidy
go mod download
go build -o render-tpl
```

## Usage Help
Run ```./render-tpl -h``` to get the latest help information.  

### Example

```
./render-tpl -v ./test/value.yaml -t ./test/template.tpl -s
```
