# Run

## Docker
For a set of scenarios in `run/configs.json`,
```sh
docker-compose run -e FILES="foo bar" app
```

## Shell
Need to install go and python3.

For a set of scenarios in `run/configs.json`,
```sh
/bin/bash ./runset.sh foo bar
```

For one scenario in `run/config.json`,
```sh
/bin/bash ./run.sh start 2> foo
```
