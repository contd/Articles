# Articles REST Service

This is a rest api service for managing saved articles from sites like Pocket or Wallabag for personal use and to hopefully replace Wallbag for myself.

## Development

Besure to add to and run tests for development.  This will create a test database `articles_test`, if one does exist, it will drop it then create it. Its not removed after the test finish so you can open mongo and look at whats inside.  Run the following commands to get verbose testing results and then get basic code coverage:

```bash
go test -v .
go test -cover .
```

### Code Coverage

For getting code coverage information and a nice html output the following can and was run to generate both the `coverage.out` file and opens a browser showing the nice html coverage report:

```bash
go test -coverprofile coverage.out .
go tool cover  -html=coverage.out
```

Also, to add the frequency of lines executed add the `-covermode` flag to the test command:

```bash
go test -covermode=count -coverprofile coverage.out .
go tool cover  -html=coverage.out
```

### Docker Build

Once the tests all pass you can build your own docker image with the included docker file like so:

```bash
docker build -t contd/articles .
```
If you use `go install` the binary will expect the `config.toml` file to be in the same directory.  You can override this by passing an environment variable like so:

```bash
CONFIG_PATH=/some/other/path/config.toml articles
```

The default port is `3000` but can also be overridden by passing the following environment variable:

```bash
SERVER_PORT=3001 articles
```

This assumes your `PATH` includes `$GOPATH/bin` and you must include the file name of the `.toml` you want to use.

## Docker Running

To run this in docker and map the default port, use the following once you've created an image from the `Dockerfile`:

```bash
docker run --name goarticles -d -p 3000:3000 contd/articles

```

To run with a connected mongodb docker container:

```bash
docker run --link mymongo:mongo -d -p 3000:3000 contd/articles
```

Changing the Dockerfile to use the above mentioned environment variables will then allow you to override the defaults.  You can also add volume mapping to a local config file to use, so you can specify a different mongodb server like one on your host machine so your data is not destroyed with your container on rebuild/etc.

Then use the `docker-articles.service` file to make the service autostart on boot:

```bash
sudo cp $GOPATH/src/github.com/contd/articles/docker-articles.service /etc/systemd/system/
sudo systemctl enable docker-articles.service
sudo systemctl start docker-articles.service
```
