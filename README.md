# knowledge-keeper

A web frontend and API for storing and retrieving knowledge obtained about varieties of subjects.

## Build Tools

### Gendal

Allows rebuilding the data access layer after database changes.

Install:

```sh
go get github.com/turnkey-commerce/gendal
```

Build DB Access Layer:

```sh
cd API
gendal
```

### Fresh

Allows hot loading of echo API.

Install:

```sh
go get github.com/pilu/fresh
```

Startup:

```sh
cd API
fresh
```

### swag

Allows building Swagger 2.0 API documents.

```sh
go get github.com/swaggo/swag/cmd/swag
```

Build Swagger:

```sh
swag init
```
