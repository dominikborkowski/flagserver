

# flag server

This is a simple program for serving plain text content over TCP connection, often used during CTF challenges. Flag file can be either built-in, or mounted as a volume during run time.

## Defaults

By default `flagserver` will listen on port 9999, and wil serve content of `/home/user/flag.txt` to each new TCP connection.


## Sample `.env` options

```
EXTERNAL_PORT=1234
CONTAINER_NAME=flagserver
FLAG_FILE_PATH=./some-other-flag.txt
```
