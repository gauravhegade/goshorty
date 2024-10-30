# GO Shorty - A URL Shortener Service

## Design is given in the image below
[](https://github.com/gauravhegade/goshorty/blob/main/goshorty_design.png)

## API Endpoints will be defined as follows

```/urls```
- GET - Get a list of all shortened URLs, returned as JSON (maybe changed later to gRPC)
- POST - Add a new URL from the request data sent as JSON

```/urls/:id```
- GET - Get a long URL by its shortened URL (id), and redirect the page to this long URL
