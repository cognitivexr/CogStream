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

A packet consists of:

* 4 byte little endian unsigned integer to encode the packet length
* the image encoded as 8bit uint array
