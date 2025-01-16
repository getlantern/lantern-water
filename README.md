# lantern-water

This library wraps the water implementation with the intention to make a easy integration between listener/client services (such as http-proxy and flashlight) and water. It also provides a way to fetch and manage WASM files locally so the listener or dialer just can use the provided version available.

![A C4 Component diagram describing how the components interact with each other. In the diagram, a client app is responsible for initializing the version control manager (VC) and if the WASM file is not available locally, it'll try to download the WASM file from the given URL or magnet link. When the WASM file is available in memory, the client app will create a dialer and make a connection with the given listener address. The diagram also show the listener app perspective, the listener download the WASM file, load in memory and create a water listener. The listener also accept connections and read the received message/bytes.](./docs/component-diagram.png)

There are two clients available in this library: `Client` and `Listener`. The `Client` is a client dialer that can be used to connect to a listener. The `Listener` is a listener that can be used to accept connections from a client using the given WASM transport. They can be executed with the following commands:

```sh
go run cmd/listener/main.go
```

And in another terminal (to send a "hello world" message to the listener):

```sh
go run cmd/dialer/main.go
```
