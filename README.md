

# flag server

This is a simple program for serving plain text content over TCP connection, often used during CTF challenges. Flag file can be either built-in, or mounted as a volume during run time.

## Defaults

By default `flagserver` will listen on port `9999` on `0.0.0.0`, and wil serve content of `/home/user/flag.txt` to each new TCP connection.

## Options

Options can be provided either via command line arguments, or environment variables. Currently environment variables take precedence over CLI arguments (yes, it's backwards)

* `-f` or `FLAG_SERVER_FILEPATH` - path to a file containing text flag. Defaults to `~/flag.txt`
* `-h` or `FLAG_SERVER_HOST` - host/IP used for listening to new traffic. Defaults to `0.0.0.0` (all traffic)
* `-p` or `FLAG_SERVER_PORT` - port to listen on. Defaults to a random port
* `-u` or `FLAG_SERVER_UDP` - use UDP instead of TCP. Defaults to `false`


## Sample `.env` options

These options can be stored in `.env` for use with `docker-compose`

```
CONTAINER_NAME=flagserver
FLAG_EXTERNAL_FILEPATH=./flag.txt # path on the docker host system for the flag file
FLAG_SERVER_FILEPATH=~/flag.txt # this is location of the flag file inside the container
FLAG_SERVER_HOST=0.0.0.0
FLAG_SERVER_PORT=9999
FLAG_SERVER_UDP=false
```
