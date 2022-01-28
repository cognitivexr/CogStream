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

* `OperationSpec`
    ```json
        {
            "code": <str>,
            "attributes": {
                <str>: [<str>, ...]
            }
        }
    ```
* `EngineSpec`
    ```json
        {
            "name": <str>,
            "attributes": {
                <str>: [<str>, ...]
            }
        }
    ```
* `AvailableEngines`
    ```json
        {
            "engines": [<EngineSpec>],
        }
    ```
* `ClientFormatSpec`
    ```json
        {
            "engine": <str>,
            "attributes": {
                <str>: [<str>, ...]
            }
        }
    ```
* `StreamSpec`
    ```json
        {
            "engineAdress": <str>,
            "attributes": {
                <str>: [<str>, ...]
            }
        }
    ```
 

### Alternative prototypical WebRTC-based connection establishment

TODO: briefly describe webrtc connection establishment

## Mediator Cluster Management