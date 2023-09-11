## vaultbin

[![MIT](https://img.shields.io/badge/license-MIT-blue.svg)](./LICENSE)
[![Release CI](https://github.com/merlinfuchs/vaultbin/actions/workflows/release.yml/badge.svg)](https://github.com/merlinfuchs/vaultbin/releases)

Vaultbin is a blazingly fast and secure alternative to Pastebin and Hastebin.

For each paste Vaultbin generates a random 128 bit AES key and encrypts the paste with it.
To identify the encrypted paste in the database it's using a hashed version of that key.
The encryption key is then used in the URL to request and decrypt the paste.  
This way nobody with access to the database can read a paste unless they have the exact URL to that paste.

Vaultbin compiles to a single binary and doesn't depend on an external database.
The binary includes everything you need to host an instance.

This project was initially inspired by [zer0b.in](https://github.com/zer0bin-dev/zer0bin).

## Public Instances

| URL                            | Expiration | Max paste size | Location |
| ------------------------------ | ---------- | -------------- | -------- |
| [vaultb.in](https://vaultb.in) | 30 days    | 69,420 chars   | Germany  |

## API Routes

| Route                    | Method | Description                   | Parameters                          |
| ------------------------ | ------ | ----------------------------- | ----------------------------------- |
| `/api/pastes`            | `POST` | Create a paste                | `language`, `content`, `expiration` |
| `/api/pastes/{paste_id}` | `GET`  | Get information about a paste | None                                |

## Installation

### Prebuilt Binaries

You can find prebuilt binaries for the most common operating systems [here](https://github.com/merlinfuchs/vaultbin/releases).

### Build from source

To build this project from source you need [Go](https://go.dev/dl/) installed.

```shell
# Install from git directly (recommended)
go install github.com/merlinfuchs/vaultbin

# Clone and build locally
git clone https://github.com/merlinfuchs/vaultbin
cd vaultbin
go build
```

## Configuration

Vaultbin will look for a `vaultbin.toml` file in the directory where you start it.  
The default config looks like this:

```toml
port = "8080"
host = "localhost"

paste_max_size = 69420 # max size of pastes in bytes
paste_ttl = 2592000 # (30 seconds) seconds after a paste will be deleted

[database]
path = "vaultbin.db" # path where data is stored

[ratelimit] # 5 request / 5 seconds
burst_size = 5 # number of request before subsequent request are block
per_second = 5 # seconds it takes to refill one request
reverse_proxy = false # if the backend is deployed behind a revers proxy -> this changes the way the peers IP is retrieved
```
