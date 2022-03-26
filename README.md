## vaultbin

[![MIT/Apache 2.0](https://img.shields.io/badge/license-MIT-blue.svg)](./LICENSE)

Vaultbin is a blazingly fast and secure alternative to Pastebin and Hastebin.

For each paste Vaultbin generates a random 128 bit AES key and encrypts the paste with it. 
To identify the encrypted paste in the database it's using a hashed version of that key.
The encryption key is then used in the URL to request and decrypt the paste.  
This way nobody with access to the database can read a paste unless they have the exact URL to that paste.

Vaultbin compiles to single binary and doesn't depend on an external database.
The binary includes everything you need to host an instance.

This project was initially inspired by [zer0b.in](https://github.com/zer0bin-dev/zer0bin).

## Public Instances

| URL                                            | Expiration | Max paste size | Location                            |
| ---------------------------------------------- | ---------- | -------------- | ----------------------------------- |
| [vaultb.in](https://vaultb.in)                 | 30 days    | 69,420 chars   | Germany                             |

## Installation

### Prebuilt Binaries

You can find prebuilt binaries for the most common operating systems [here](https://github.com/merlinfuchs/vaultbin/releases).

### Build from source

To build this project from source you need [node](https://nodejs.org/en/download/) and [rust](https://www.rust-lang.org/tools/install) installed.

Build the frontend:
```shell
npm run install
npm run build
```

Build the backend:
```shell
# just build
cargo build --release
# or build and install it in PATH
cargo install --path .
```

The frontend code is embedded into the binary during the build process. 

## Configuration

Vaultbin will look for a `Config.toml` file in the directory where you start it.  
The default config looks like this:
```toml
host = "127.0.0.1" # host to bind on
port = 8080 # port to bind on

max_paste_size = 69420 # max size of pastes in bytes
max_expiration = 2592000 # (30 days) seconds after a paste will be deleted

[database]
path = "./data" # path where data is stored
cache_size = 100000000 # (100MB) size of the database cache (recently used pastes will be kept in memory if possible)
```

In addition to creating a `Config.toml` file you can also override these values using environment variables:
For example:
```shell
VAULTBIN_HOST=0.0.0.0
VAULTBIN_DATABASE__CACHE_SIZE=9999999
```