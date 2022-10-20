health check: /ping

function route: /dial

environments:

+ PORT: server port, default value is `8080`.
+ CONFIG_PATH: resolved config, default value is `/data/config.toml`. example is `./config.toml`

build:

+ make build.image
+ CONFIG_PATH={} PORT={} make build.image

run:

+ docker run -d -p 8080:8080 -v $PWD/config.toml:/data/config.toml server:release