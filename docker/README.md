# docker

Example Forgejo installation hosted by [`docker-compose`](https://github.com/docker/compose).

## Usage

Install `docker-compose` if not already done:

```command
# curl -LO https://github.com/docker/compose/releases/download/v2.29.7/docker-compose-linux-x86_64 /usr/local/bin/docker-compose
# chmod +x $_
```

Run the following commands if you're using **Podman**:

```command
$ systemctl --user start podman.socket
$ export DOCKER_HOST=unix:///run/user/1000/docker.sock
$ export CONTAINER_HOST=unix:/run/user/1000/podman/podman.sock
```

Start the stack:

```command
$ docker-compose up --no-start
$ docker-compose start
```

Forgejo listens on [port 3000](http://localhost:3000).

In the initial configuration wizard, configure the database:

- Database type: MySQL
- Host: `forgejo_db`
- Username: `forgejo`
- Password: `password`
- Database name: `forgejo`
- SSH server port: 2222

Create an administrator account:

- Administrator username: <<< username >>>
- Email address: <<< username >>>@localhost
- Password: <<< password >>>
- Confirm password: <<< password >>>

Finally, click **Install Forgejo**.

## Generate API Token for API Usage

> https://forgejo.org/docs/latest/user/api-usage/#generating-and-listing-api-tokens

```command
$ curl -H "Content-Type: application/json" -d '{"name":"test"}' -u <<< username >>>:<<< password >>> http://localhost:3000/api/v1/users/<<< username >>>/tokens
{"id":1,"name":"test","sha1":"9fcb1158165773dd010fca5f0cf7174316c3e37d","token_last_eight":"16c3e37d"}
```

Save the token (sha1), it will not be shown again!

## Swagger

Swagger is available at <http://localhost:3000/api/swagger>.
