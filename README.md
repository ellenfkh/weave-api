# weave-api
REST api for exploring a FS.

# API Reference

## List Files
    Returns the contents of a directory, or the stat and contents of a file .

* **URL:** `/:path`
* **Method:** `GET`
*  **URL Params:** `None`
* **Data Params:** `None`
* **Success Response:**
  * **Code:** 200 <br />
    **Content:** 
    If the path refers to a file: 
    ```javascript
    { 
        "Stat": {
                "Name": "filename",             // name of the file
                "Permissions": "-rw-r--r--",    // File permissions
                "Owner": "root",                // owner
                "Size": 11,                     // size
                "IsDir": false                  // Whether this is a directory
                },
        "Content": "contents of the file"
    }
    ```
    If the path refers to a directory: 

    ```javascript
    [
        {
        "Name": "filename",             // name of the file
        "Permissions": "-rw-r--r--",    // File permissions
        "Owner": "root",                // owner
        "Size": 11,                     // size
        "IsDir": false                  // Whether this is a directory
        },
    ]
    ```
 
* **Error Response:**

  * **Code:** 404 NOT FOUND <br />
    **Content:** `{ error : "File does not exists/is malformed" }`

# Development
Requirements: 

- go (tested on `go version go1.14 darwin/amd64`)
- docker (tested on `Docker version 18.06.1-ce, build e68fc7a`)
- make

## Quickstart

To get up and running immediately:
```bash
$ make run baseDir=$HOME
```
This will serve up your home directory on the host machine's port 8080.

## Building
Build the binary locally  with `go build`.

Build the docker image with `make build`. This will create a docker image
named `weave-api` tagged with the short name of the current git sha.

Note that docker will complain if you add a file in `./` that it can't read,
so if you're adding forbidden dirs and files for testing, don't do it until
after build.

## Testing
Run unit tests with `make unit`. These are not dockerized (for the sake of
iteration speed), so you'll have to install go first.

To run the same unit tests in a docker container, run `make test`.

## Run

### Locally
After running `go build`, run locally using `./weave-api`. 
```
$ ./weave-api --help
Usage of ./weave-api:
  -baseDir string
        The directory to mount the ls handler (default "./test")
  -port int
        The port to listen on (default 8080)
  -test
        Whether to artificially populate some files under ./test
```
The suggested command for getting up and running quickly with a generated local test dir is `./weave-api -test`, which will create a
directory to explore and point the fs handler at it. This will then be reachable at `localhost:8080`. 

### Dockerized
To run in docker, pick a path to mount and run the following command:
```bash
make run baseDir=<full-path-to-dir-you-want-to-mount>
```
This will start the image and expose the api on port 8080.


# Known Limitations and further work
- Because of the [behaviour](https://github.com/docker/for-mac/issues/2657)
of docker mounts on Mac, the owner will always be displayed as
container'suser, which will usually be `root`. This should not happen on
Linux systems, but I didn't have a convenient Linux env to test in. This
could probably be avoided by copying the files inside the container instead
of mounting, but I preferred to mount so that it could be extended to modify
local fs files directly without copying back out, and so that changes on the
host fs would be picked up seamlessly while the container was running.

- This repo doesn't have any integration tests. To be complete, it should
have a docker-compose target that spun up the `weave-api` and then another
container to hit the exposed endpoint.

- This repo lacks benchmark tests -- this api isn't particularly
designed for very high load, nor is it tested on such.

- There are no tests for forbidden files -- I manually tested against
`drw-------` dirs owned by root, but adding creation and deletion of test
dirs and files would require more permission massaging to get working on both
docker and local without littering, I left those un-unittested.

- Without a k8s cluster to run on or a docker registry to pull from, there wasn't much point making a helm template. However, on a cluster
deploy, I'd probably turn the arguments to the binary into environment
variables instead of arguments, then template this out into a bare-bones
combined `deployment` and `service` manifest template with the port/baseDir
baked directly in from values, probably not a `configmap` unless there were
more configs expected to be added in the future.