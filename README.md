<p align="center">
  <img src="wordmark.png" />
</p>

⚠️ This project is a work in progress ⚠️ 

## Features

- Provides a centralized repository for Nomad HCL templates
- Offers API for CD tools to automate deployment of new versions
- Simple Web UI for editing templates & monitoring deployments
- Supports Nomad instances behind Cloudflare Access Zero Trust network

Build with [Go](https://go.dev/), [Tailwind CSS](https://tailwindcss.com/) und [SQLite](https://sqlite.org/).

## Development

### Requirements

- Go 1.19+
- Yarn

### Go app

#### Install dependencies

```
go get
```

#### Build & hot reload

```shell
gow -e=go,html run .
```

### Frontend

#### Install dependencies

```
yarn install
```

#### Build & hot reload

```shell
yarn watch
```

## License

Released under the [MIT License](LICENSE.md).

## Authors

- [Roman Zipp](https://github.com/romanzipp)
