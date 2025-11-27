# File Serving API

This is a simple file serving API built with Go and Gin.

## Features

- List files
- Upload files
- Download files
- Swagger Documentation

## Running Locally

1. Install dependencies:

   ```bash
   go mod download
   ```

2. Run the application:

   ```bash
   go run main.go
   ```

3. Access Swagger UI at `http://localhost:8080/swagger/index.html`.

## Running with Docker

1. Build the image:

   ```bash
   docker build -t file-api .
   ```

2. Run the container with a bind mount for storage:
   ```bash
   docker run -p 8080:8080 -v $(pwd)/storage:/app/storage file-api
   ```
   Replace `$(pwd)/storage` with your desired local storage path.

## API Endpoints

- `GET /files`: List all files.
- `POST /files`: Upload a file (multipart/form-data, key: `file`).
- `GET /files/:filename`: Download a file.

## Updating Swagger Documentation

To regenerate and update the Swagger documentation after you change API comments or handlers, follow one of the approaches below.

Windows (PowerShell):

1. Install the `swag` CLI (if you haven't already):

   ```powershell
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

   If `swag` is not on your `PATH`, you can run it directly from your Go bin directory:

   ```powershell
   & "$env:USERPROFILE\go\bin\swag.exe" init
   ```

2. From the project root run:

   ```powershell
   & "$env:USERPROFILE\go\bin\swag.exe" init
   ```

   This will regenerate `docs/swagger.json`, `docs/swagger.yaml` and `docs/docs.go` based on the `// @` annotations in your source files.

Unix / macOS:

1. Install `swag`:

   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

2. Ensure `$GOBIN` or `$GOPATH/bin` is in your `PATH`, then run from the project root:

   ```bash
   swag init
   ```

Notes and tips:

- Run `swag init` from the project root so it finds `main.go` and your package imports (the generator looks for comment annotations).
- Ensure your `main.go` (or router package) includes the import for generated docs, e.g. `import _ "github.com/your/module/docs"` so the generated `docs` package is compiled into the binary.
- If you change annotations, rerun `swag init` and then rebuild your binary or Docker image.

Rebuild and redeploy Docker image:

```bash
docker build -t file-api .
docker stop file-api || true
docker rm file-api || true
docker run -d --name file-api -p 8080:8080 -v /path/to/local/storage:/app/storage file-api
```

Replace `/path/to/local/storage` with your local storage directory. On PowerShell you can use `$(pwd)/storage` or an absolute path.
