# flagserver

This is a simple Go program for serving plain text content over TCP, or UDP, connection, often used during CTF challenges. Options can be supplied either via command line arguments or environment variables. Can be used as a stand-alone program, or via supplied docker images.


## Defaults

By default `flagserver` will listen on port `9999` on `0.0.0.0`, and wil serve content of `/home/user/flag.txt` to each new TCP connection.

## Options

Options can be provided either via command line arguments, or environment variables. Currently environment variables take precedence over CLI arguments (except for the content. And yes, it's backwards)

* `--content` or `FLAG_SERVER_CONTENT` - optional way to provide text flag via command line or an environment variable
* `--file_path` or `FLAG_SERVER_FILE_PATH` - path to a file containing text flag. Defaults to `~/flag.txt`
* `--host` or `FLAG_SERVER_HOST` - host/IP used for listening to new traffic. Defaults to `0.0.0.0` (all traffic)
* `--port` or `FLAG_SERVER_PORT` - port to listen on. Defaults to a random port
* `--protocol` or `FLAG_SERVER_PROTOCOL` - use `tcp`, `udp`, or `http`. Defaults to `tcp`
* `--http-path` or `FLAG_SERVER_HTTP_PATH` - specify what HTTP path to listen on. Used only for `http` protocol, defaults to `/`


## Sample `.env` options

These options can be stored in `.env` for use with `docker-compose`

```
# name of the container
CONTAINER_NAME=flagserver
# path on the docker host system for the flag file
FLAG_EXTERNAL_FILE_PATH=./flag.txt
# this is location of the flag file inside the container
FLAG_SERVER_FILE_PATH=/home/user/flag.txt
# address to listen to
FLAG_SERVER_HOST=0.0.0.0
# port to listen on
FLAG_SERVER_PORT=9999
# specify which protocol to use: TCP (default), UDP, or HTTP
FLAG_SERVER_PROTOCOL=http
# specify what path should HTTP requests reply to, defaults to /
FLAG_SERVER_HTTP_PATH=/content

# instead of specifying paths for flags, one can use specify content via
FLAG_SERVER_CONTENT=supersecretword
```

## Example test run

```
$ cd main
$ go build && ./flagserver --content "foobar" --port 9999 --protocol http
2023/04/21 22:23:40 Starting new flag server instance
2023/04/21 22:23:40 Host: 0.0.0.0
2023/04/21 22:23:40 Port: 9999
2023/04/21 22:23:40 File path: ~/flag.txt
2023/04/21 22:23:40 Protocol: http
2023/04/21 22:23:40 Using content from command line: foobar
```
