# PuSu Engine daemon

PuSu Engine is a Pub-Sub system.

This project provides `pusud`, the relay server in the PuSu Engine system. The
clients for other languages are hosted in other repositories.


## Running

Go to [Golang.org](https://golang.org) and download Go for your platform.

After running the installer, make sure you set up your `PATH` and `GOPATH` appropriately. Basically `go` needs to be found in `PATH`, and `GOPATH` just needs to be set to e.g. `~/go`. Also make sure you have `git` installed for access from the CLI.

Fetch the code and dependencies, and build:

```
go get github.com/PuSuEngine/pusud
cd $GOPATH/src/github.com/PuSuEngine/pusud
go build pusud.go
```

On *nix systems you can then run it with:

```
./pusud
```

On Windows it's just:

```
pusud
```

During development you can simply build & run it in one command, which is not
recommended for normal use:

```
go run pusud.go
```


## Client libraries

To deliver messages, you will need to use a client library to send and receive them. The currently available officially supported libraries are listed here:

 * [GoPuSu](https://github.com/PuSuEngine/gopusu)


## Configuration

The configuration file uses the [YAML](http://yaml.org/) format.

Full example:

```
relays:
  - relay-1.example.com
  - relay-2.example.com
authenticator: MyPlugin
```


### Network

You define a network by defining the list of relays on every node, so they know how to connect to each other.

```
relays:
  - relay-1.example.com
  - relay-2.example.com
```


### Authentication plugin

You can define a requirement to authenticate to channels for reading, and writing, using the authentication system.

First, choose an authentication plugin you want to use from the `plugins/` -folder, and configure it as the authenticator.

```
authenticator: MyAuthenticator
```

To disable authentication (allow read and write for everyone, not recommended other than for purely internal use) use the `None` authenticator.

```
authenticator: None
```

Check the documentation for the specific authenticators for further information on configuration.


## Architecture description

The design is built based on the assumption that the workload is heavily read focused, having a lot of clients reading
relatively few messages written by a few sources.

Components:

 * `Source` - A client sending messages.
 * `Relay` - A server relaying the messages.
 * `Network` - Multiple relay servers connected together.
 * `Listener` - A client receiving messages.
 * `Channel` - A unique name for a target for **source**s to send messages to and **listener**s to listen to.

**Relay** servers connect to all other **relay** servers to build a *network*. **Listener** connects to a **relay** and
authenticates to **read** on the **channel**s they are interested in. **Source** connects to a **relay** and authenticates to **write** on the **channel**s they are interested in. **Source** sends a message to a **channel** on a **relay**,
which first distributes it to all other **relay** servers in the **network**, then sends it to all **listener**s connected to that **relay**
listening for messages on the **channel**.

```
Relay <-> Relay    # Build network

Listener -> Relay  # Auth READ for my.channel
Source -> Relay    # Auth WRITE for my.channel

Source -> Relay    # Write "foo" to my.channel
Relay -> Network   # Write "foo" to my.channel
Relay -> Listener  # Write "foo" to my.channel
```


## Security considerations

The relay network is assumed to be secured in a private network behind a firewall. The communications between the relays are not secured in any special way, if you give 3rd parties direct access (e.g. by not firewalling the servers properly) they can inject any messages they want, and probably read messages rather easily too.

`Client` and `source` connections are using websocket, by default connecting to TCP port `55000`, the `relay` servers communicate in the network using TCP port `55001`.


## Testing

When making changes to the code you might want to run the unit tests.

```
go test -v ./...
```

## License

Short version: MIT + New BSD.

Long version: Read the LICENSE.md -file.
