<h1 align="center">Go Fiber Starter</h1>

<p align="center">
  <img alt="Github top language" src="https://img.shields.io/github/languages/top/irsx/go-fiber-starter?color=56BEB8">

  <img alt="Github language count" src="https://img.shields.io/github/languages/count/irsx/go-fiber-starter?color=56BEB8">

  <img alt="Repository size" src="https://img.shields.io/github/repo-size/irsx/go-fiber-starter?color=56BEB8">

  <img alt="License" src="https://img.shields.io/github/license/irsx/go-fiber-starter?color=56BEB8">

  <img alt="Github issues" src="https://img.shields.io/github/issues/irsx/go-fiber-starter?color=56BEB8" />

  <img alt="Github forks" src="https://img.shields.io/github/forks/irsx/go-fiber-starter?color=56BEB8" />

  <img alt="Github stars" src="https://img.shields.io/github/stars/irsx/go-fiber-starter?color=56BEB8" />
</p>

<!-- Status -->

<p align="center">
  <a href="#dart-about">About</a> &#xa0; | &#xa0;
  <a href="#sparkles-features">Features</a> &#xa0; | &#xa0;
  <a href="#rocket-technologies">Technologies</a> &#xa0; | &#xa0;
  <a href="#white_check_mark-requirements">Requirements</a> &#xa0; | &#xa0;
  <a href="#checkered_flag-starting">Starting</a> &#xa0; | &#xa0;
  <a href="#memo-license">License</a> &#xa0; | &#xa0;
</p>

<br>

## :dart: About

Simple and scalable starter kit to build powerful and organized REST projects with Fiber.

## :sparkles: Features

-   [x] Repository Pattern
-   [x] Logging
-   [x] Live Reloading
-   [x] Redis Cache
-   [x] RabbitMQ Consumer & Publisher
-   [x] Server Side Event (SSE)
-   [x] ORM Database
-   [x] SQL Migration & Seeders
-   [x] Custom REST Client
-   [x] Image upload to CDN
-   [x] Excel Importer
-   [x] JWT Authentication
-   [x] CI\CD with Github Actions
-   [x] Docker

## :rocket: Technologies

The following tools were used in this project:

-   [Go](https://go.dev)
-   [Fiber](https://github.com/gofiber/fiber)
-   [Gorm](https://gorm.io)
-   [SQLMigrate](https://github.com/rubenv/sql-migrate)
-   [Zaplog](https://github.com/uber-go/zap)
-   [Air](https://github.com/cosmtrek/air)
-   [PostgreSQL](https://www.postgresql.org)
-   [Docker](https://www.docker.com/)

## :white_check_mark: Requirements

Before starting :checkered_flag:, you need to have [Git](https://git-scm.com), [Go](https://go.dev), [Docker](https://www.docker.com/) and [PostgreSQL](https://www.postgresql.org) installed.

## :checkered_flag: Starting

```bash
# Clone this project
$ git clone https://github.com/irsx/go-fiber-starter

# Access
$ cd go-fiber-starter

# Download dependencies
$ go get

# Run the project
$ go run main.go

# Run migrations and seeders
$ go run main.go --rollback --seed

# Run the project with live reloading
$ air

# The server will initialize in the <http://{host}:{port}>
```

## :memo: License

This project is under license from MIT. For more details, see the [LICENSE](LICENSE) file.

Made with :heart: by <a href="https://github.com/irsx" target="_blank">AVIANA DEV</a>

&#xa0;

<a href="#top">Back to top</a>
