# lantern-water

This library wraps the water implementation with the intention to make a easy integration between listener/client services (such as [http-proxy](https://github.com/getlantern/http-proxy) and [flashlight](https://github.com/getlantern/flashlight)) and water. It also provides a way to fetch and manage WASM files locally so the listener or dialer just can use the provided version available.

![A C4 Component diagram describing how the components interact with each other. In the diagram, a dialer app is responsible for initializing the version control manager (VC) and if the WASM file is not available locally, it'll try to download the WASM file from the given URL or magnet link. When the WASM file is available in memory, the dialer app will create a dialer and make a connection with the given listener address. The diagram also show the listener app perspective, the listener download the WASM file, load in memory and create a water listener. The listener also accept connections and read the received message/bytes.](./docs/component-diagram.png)

There are two programs available in this library: [`Dialer`](https://github.com/getlantern/lantern-water/blob/f07491cb2b423622182b08b71804f2ff6bdafbca/cmd/dialer/main.go) and [`Listener`](https://github.com/getlantern/lantern-water/blob/f07491cb2b423622182b08b71804f2ff6bdafbca/cmd/listener/main.go). The `Dialer` is used to create connections to a listener and send data through a given protocol implementation. The `Listener` is used for accepting connections from a client using the given WASM transport and in this case it only prints the connection message as an example. They can be executed with the following commands:

```sh
go run cmd/listener/main.go
```

And in another terminal (to send a "hello world" message to the listener):

```sh
go run cmd/dialer/main.go
```
