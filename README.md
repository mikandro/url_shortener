# URL Shortener

A simple URL shortener service written in Go (Golang) that provides a RESTful API for creating and retrieving shortened URLs.

## Features

- **Create Short URLs**: Generate shortened URLs that redirect to the original addresses.
- **Custom Short Codes**: Optionally specify custom codes for your URLs.
- **Redis Integration**: Uses Redis for fast storage and retrieval.
- **Docker Support**: Easily deployable with Docker.
- **Unit Tests**: Includes tests to ensure reliability.

## Installation

### Prerequisites

- [Go](https://golang.org/dl/) (version 1.16 or higher)
- [Redis](https://redis.io/download)
- [Docker](https://www.docker.com/get-started) (optional, for containerization)

### Clone the Repository

```bash
git clone https://github.com/mikandro/url_shortener.git
cd url_shortener
```

### Configure Environment Variables

Create a `.env` file in the root directory and set the necessary environment variables:

```env
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
SERVER_PORT=8080
```

### Run the Application

#### Using Go

```bash
go run cmd/url_shortener/main.go
```

#### Using Docker

```bash
docker build -t url_shortener .
docker run -p 8080:8080 --env-file .env url_shortener
```

## Usage

### API Endpoints

- **POST** `/shorten`

  Create a shortened URL.

  **Request Body:**

  ```json
  {
    "url": "https://www.example.com",
    "custom_code": "example"
  }
  ```

  *Note: The `custom_code` field is optional.*

  **Response:**

  ```json
  {
    "short_url": "http://localhost:8080/{code}"
  }
  ```

- **GET** `/{code}`

  Redirects to the original URL associated with the given code.

### Example

#### Create a Short URL

```bash
curl -X POST -H "Content-Type: application/json" -d '{"url":"https://www.example.com"}' http://localhost:8080/shorten
```

**Response:**

```json
{
  "short_url": "http://localhost:8080/abc123"
}
```

#### Redirect Using Short URL

Navigate to `http://localhost:8080/abc123` in your browser to be redirected to `https://www.example.com`.

## Dependencies

- [Go Redis Client](https://github.com/go-redis/redis)
- [Chi Router](https://github.com/go-chi/chi)
- [Testify](https://github.com/stretchr/testify) (for testing)

## Testing

Run the unit tests using:

```bash
go test ./...
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contact

Created by [mikandro](https://github.com/mikandro). Feel free to reach out!
