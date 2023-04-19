# flagserver

This is a simple Go program for serving plain text content over TCP, or UDP, connection, often used during CTF challenges. Options can be supplied either via command line arguments or environment variables. Can be used as a stand-alone program, or via supplied docker images.


## Defaults

By default `flagserver` will listen on port `9999` on `0.0.0.0`, and wil serve content of `/home/user/flag.txt` to each new TCP connection.

## Options

Options can be provided either via command line arguments, or environment variables. Currently environment variables take precedence over CLI arguments (except for the content. And yes, it's backwards)

* `-c` or `FLAG_SERVER_CONTENT` - optional way to provide text flag via command line or an environment variable
* `-f` or `FLAG_SERVER_FILEPATH` - path to a file containing text flag. Defaults to `~/flag.txt`
* `-h` or `FLAG_SERVER_HOST` - host/IP used for listening to new traffic. Defaults to `0.0.0.0` (all traffic)
* `-p` or `FLAG_SERVER_PORT` - port to listen on. Defaults to a random port
* `-u` or `FLAG_SERVER_UDP` - use UDP instead of TCP. Defaults to `false`


## Sample `.env` options

These options can be stored in `.env` for use with `docker-compose`

```
# name of the container
CONTAINER_NAME=flagserver
# path on the docker host system for the flag file
FLAG_EXTERNAL_FILEPATH=./flag.txt
# this is location of the flag file inside the container
FLAG_SERVER_FILEPATH=/home/user/flag.txt
# address to listen to
FLAG_SERVER_HOST=0.0.0.0
# port to listen on
FLAG_SERVER_PORT=9999
# set to true to UDP isntead of default TCP
FLAG_SERVER_UDP=false

# instead of specifying paths for flags, one can use specify content via
FLAG_SERVER_CONTENT=supersecretword
```

## Example test run

```
$ cd main
$ go build && ./flagserver -c "foobar" -p 9999

2023/04/19 12:12:16 Starting new flag server instance
2023/04/19 12:12:16 Host: 0.0.0.0
2023/04/19 12:12:16 Port: 9999
2023/04/19 12:12:16 Filepath: ~/flag.txt
2023/04/19 12:12:16 Use UDP: false
2023/04/19 12:12:16 Using content from command line: foobar
2023/04/19 12:12:16 Listening for TCP connections on [::]:9999
```
