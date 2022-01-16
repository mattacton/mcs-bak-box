# MCS-BAK-BOX
Minecraft Server Backup the Box is a Google Cloud pub/sub function built to backup a Minecraft server via FTP to a Google Cloud Storage bucket.

## Files backed up

* `/world.zip` - The function does not zip the world directory, this must be done before the function runs
* `/ops.json`
* `/server.properties`
* `whitelist.json`
* `banned-ips.json`
* `banned-players.json`

## Required Environment Variables

* `FTP_USER`
* `FTP_PWD`
* `FTP_HOST`
* `FTP_PORT`
* `BUCKET`

## Function Suggestions

* At least 2GB or more of function memory. The backup files are temporarily stored in memory during backup
* An extended function timeout. The backup can take several minutes based on world size and speed of the FTP server.

## Running Locally

The main package is within the `cmd` directory. To run the function locally, execute `go build` and `./cmd` from within the `cmd` directory.

The environment variables can either be set by normal means, or they can be hard coded within the `var` block of `main.go`. **NOTE:** *Only do this to run the function locally. Do not commit / push credentials or other sensitive information to the repository.*