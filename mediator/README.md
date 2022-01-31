# CogStream Mediator

The mediator is the part of CogStream which connects clients to corresponding engines.
The engine is selected via attributes given in a handshake protocol in which the client communicates with the mediator. 

## Mediator Handshake

Each connection of a CogStream client with an engine starts with a websocket connection to the mediator. 
Inside the websocket connection 2 message exchanges happen, thus performing the handshake.

1. RequestOperation:
    - The client requests an operation that needs to be performed on the data to be delivered. This operation is encoded as an `OperationSpec`.
    - The mediator replies with the list of engines under its control which can perform such an operation. The engines are encoded as an `EngineSpec`. 
2. EstablishFormat: 
    - After internal selection of the suitable engine, the client initiates the establishment of a common data format. The information coming from the client is encoded as a `ClientFormatSpec`.
    - Upon receiving the client format specification, the mediator sets up the Engine and internal data transformation and returns a `StreamSpec` with which the client is able to connect to the appropriate engine.

After the handshake the client continues to communicate with the engine, i.e. delivering stream data and receiving engine results.

![Mediator Handshake](/assets/mediator_handshake.png)

### Message Format

The messages are JSON-encoded and consist of the following fields:

- `OperationSpec`

    ```json
    {
        "code": <str>,
        "attributes": {
            <str>: [<str>, ...]
        }
    }
    ```

- `EngineSpec`

    ```json
    {
        "name": <str>,
        "attributes": {
            <str>: [<str>, ...]
        }
    }
    ```

- `AvailableEngines`

    ```json
    {
        "engines": [<EngineSpec>],
    }
    ```

- `ClientFormatSpec`

    ```json
    {
        "engine": <str>,
        "attributes": {
            <str>: [<str>, ...]
        }
    }
    ```

    Example: to indicate to the mediator the stream format, the following attributes must be specified:

    ```json
    {
        "engine": "fermx",
        "attributes": {
            "format.width": 800,
            "format.height": 600,
            "format.colorMode": 4,
            "format.orientation": 3,
        }
    }
    ```

    Where `colorMode` and `orientation` refer to the respective CogStream format objects defined [here](https://github.com/cognitivexr/CogStream/blob/master/api/format/core.go).

- `StreamSpec`

    ```json
    {
        "engineAddress": <str>,
        "attributes": {
            <str>: [<str>, ...]
        }
    }
    ```

### Example

Run the mediator and connect to it via a websocket client:

```bash
websocat ws://localhost:8191
> {"type":2,"content":{"code": "analyze", "attributes": {}}}
< {"type":3,"content":{...}}
> {"type":4,"content":{"engine": "fermx", "attributes": {}}}
< {"type":5,"content":{"engineAddress":"0.0.0.0:37597","attributes": {...}}}
```

### Alternative prototypical WebRTC-based connections

Alternatively to using the custom TCP-based protocol involving the handshake, we also have experimented with a WebRTC-based approach (see [webrtc-pipeline branch](/cognitivexr/CogStream/tree/webrtc-pipeline)).
Here, the entire handshake is redundant and replaced by the mechanisms in WebRTC.
We use the [WebRTC implementation from Pion](https://github.com/pion/webrtc), which also provides a [comprehensive explanation](https://webrtcforthecurious.com/) on how connections are established with this protocol/networking stack.

The prototype showed promise in first tests regarding performance, but it lacks some of the features provided by our own protocol, like the stream format specification and the engine metadata.
This, however, can be easily integrated with stream metadata features from WebRTC.  
Other aspects which need to be tested are the impact of the compression on the AI performance regarding quality of results as well as the impact of the compression on the energy consumption.
