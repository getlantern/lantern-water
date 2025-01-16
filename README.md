# lantern-water

This library wraps the water implementation with the intention to make a easy integration between listener/client services (such as http-proxy and flashlight) and water. It also provides a way to fetch and manage WASM files locally so the listener or dialer just can use the provided version available.

```mermaid
C4Component
title Component diagram illustrating how lantern-water integration works

Container_Boundary(clientDialer, "Client dialer") {
    Component(versionControl, "Version control Manager", "go", "Allow the client to manage WASM files available locally. If the WASM file isn't available locally and was never loaded it'll try to download the file and load the file in memory. If the WASM file isn't used in the period of 7 days, it'll be deleted")
    Component(downloader0, "Downloader", "go", "Provide ways to download WASM files from given URLs. Currently it supports HTTPS URLs and magnet links.")
    Component(dialer, "Dialer", "go, water", "The dialer initializes a water dialer with the given WASM bytes and parameters.")
    Component(clientApp, "Client app", "go", "The example client app here (which is similar to what flashlight does) use the packages available to download and run a water dialer")
}

Rel(clientApp, versionControl, "Request WASM")
Rel(versionControl, downloader, "If WASM isn't available locally, try to download it")
Rel(clientApp, dialer, "When WASM is available in memory, initialize dialer")

Container_Boundary(appListener, "Listener") {
    Component(downloader1, "Downloader", "go", "Provide ways to download WASM files from given URLs. Currently it supports HTTPS URLs and magnet links.")
    Component(listener, "Listener", "go, water", "The dialer initializes a water dialer with the given WASM bytes and parameters.")
    Component(listenerApp, "Listener app", "go", "The example listener app here (which is similar to what http-proxy does) use the packages available to download and run a water listener")
}

Rel(listenerApp, downloader, "Request WASM")
Rel(listenerApp, listener, "Initialize listener when WASM is available in memory")
```

There are two clients available in this library: `Client` and `Listener`. The `Client` is a client dialer that can be used to connect to a listener. The `Listener` is a listener that can be used to accept connections from a client using the given WASM transport. They can be executed with the following commands:

```sh
go run cmd/listener/main.go
```

And in another terminal (to send a "hello world" message to the listener):

```sh
go run cmd/dialer/main.go
```
