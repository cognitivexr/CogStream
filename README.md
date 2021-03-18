CogStream
=========

The CogStream video analytics system consists of four components:

* `api`: shared API modules
* `mediator`: server component that mediates a handshake between a client and the streaming system
* `engines`: engines are services that perform video stream processing tasks
* `clients`: client code for the handshake and streaming protocol for various platforms

Build and run
-------------

### Mediator and engines

Running `make` in the root directory builds the mediator and the go-based engines as plugins.
The build creates a `dist` folder that should look like this:

    dist
    ├── engines
    │   ├── record.so
    │   └── record.so.spec.json
    └── mediator

Where `dist/engines` contains the engines as go plugins that are loaded by the mediator.
Run `(cd dist; ./mediator)` to start the mediator.


### Clients

To build and run the python client, run `cd clients/python` and execute `make venv`.
Run `source .venv/bin/activate` to start the virtual environment.
Then run

    python -m cogstream.cli.client --operation record

to start a client that starts recording a video using the recorder engine.
The videos are currently stored into `/tmp`.

Handshake
---------

TODO: document handshake

Streaming protocol
------------------

Streaming is initiated between a client and an engine by first sending the negotiated StreamSpec, serialized as UTF-8 encoded JSON.
The packet is prefixed with an uint32 (4 byte little endian) to indicate the string length.

The remaining packets on the connection are of type `FramePacket`, which are encoded as follows:

    +----------+------------------------+
    | Offset   | Field                  |
    +----------+------------------------+
    |        0 | Stream Id              | HEADER (little endian uint32 fields)
    |        4 | Frame Id               |
    |        8 | Unixtime seconds       |
    |       12 | Unixtime nanoseconds   |
    |       16 | Metadata length (L_m)  |
    |       20 | Data length (L_d)      |
    +----------+------------------------+
    |       24 | Metadata               | BODY
    | 24 + L_m | Data                   |
    +----------+------------------------+
