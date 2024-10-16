# Description

A simple url shortener with the goal of learning Go language and principles, builded to be production ready with Docker, Nginx (reverse proxy) and Redis.

# Usage

To run this project, you might have installed:
- [Docker](https://docs.docker.com/engine/install/) (required)
- [Go](https://go.dev/doc/install) (for development usage)

### Clone the project

```bash
git clone https://github.com/brunopstephan/url-shortener.git
```

### Setup enviroment

In the project root, you'll see a `.env.example` file, you must rename it to `.env`.

This file have all enviroment variables that the project will need to function correctly, so it's very important to you to follow this step.

### Run containers

```bash
docker compose up -d
```

If everthing goes well, you will be able to make requests at `http://localhost:9000`

### Endpoints

#### Public:
- `GET /api/{code}` - redirect to the code's url (`json=true` query param will bring the url data in JSON format);
- `POST /api/shorten` - create a shortened url, `URL` body is required with a url and it reponse with the shortened code;

##### Protected:
These endpoints are protected with **Basic Auth**, the default in `.env.example` is `admin:admin`, so transform it into Base64 and pass a `Authorization` header in the request with value like: ``Basic myCredentialsToBase64``

- `GET /admin` - get all shortened urls;
- `DELETE /admin/{code}` - delete a shortened url;
- `PUT /admin/{code}` - update the url of shortened url;



