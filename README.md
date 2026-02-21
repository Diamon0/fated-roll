# Fated Rolls
A simple dice roller bot for Fluxer

## How do I set this up?
Either install it and run it, or run the container

### Running the container (recommended)
First, pull the image (For example, with Docker).
```bash
docker pull ghcr.io/diamon0/fated-roll:latest
```
Then, run it.
```bash
docker run -e BOT_TOKEN=YOUR_BOT_TOKEN --name fated-rolls ghcr.io/diamon0/fated-roll:latest
```
If you don't know where to get it, it's on your Fluxer profile when Developer Mode is enabled under Advanced Settings; then, at the very bottom on Applications, create an application, and get your bot token. Additionally, you can get its invite URL in the same place, selecting BOT and then manage messages and send messages permissions.

### Running directly
First, have Go (Golang) installed.
Then, download and build the binary with go.
```bash
go install github.com/Diamon0/fated-roll/cmd/fated-rolls@latest
```
Finally, set the environment variable BOT_TOKEN with your bot token, and run the fated-rolls binary from your go/bin directory. For example, assuming your go/bin folder is in PATH, in Bash you can do it like this:
```bash
BOT_TOKEN=YOUR_TOKEN fated-rolls
```
