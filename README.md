

### Build for Multiple Platforms
To build your application for different platforms, you can set the GOOS and GOARCH environment variables. Here’s how to build for Windows, Linux, and macOS:

For Windows:

```bash
GOOS=windows GOARCH=amd64 go build -o MyApp.exe main.go
```

For Linux:

```bash
GOOS=linux GOARCH=amd64 go build -o MyApp main.go
```

For macOS:

```bash
GOOS=darwin GOARCH=amd64 go build -o MyApp main.go
```