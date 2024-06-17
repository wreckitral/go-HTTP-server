# Go HTTP Server
simple HTTP/1.1 server built using go. Deep dive into low level how tcp servers, and HTTP protocols works.

# How to run
- `make build`, build the program
for specifying the directory for endpoint that handles file
- `./bin/server --directory <yourdirname>`
for the other endpoints
- `./bin/server`

# Endpoints
- `/`
- `/echo/<randomString>` -> can be specified with Accept-Encoding: gzip
- `/files/{filename}` -> GET to read inside the file in the speficied directory and return it in response, POST to write file in the specified directory with request body inside of it
