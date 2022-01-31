CogStream
=========

⚠️ CogStream is work in progress and subject to rapid change, as well as poor documentation.

---

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
    │   ├── record.so
    │   └── record.so.spec.json
    └── mediator

Where `dist/engines` contains the engines as go plugins that are loaded by the mediator.
Run `(cd dist; ./mediator)` to start the mediator.


#### Mediator using Docker

Build the go docker image:

    docker build -f build/mediator/Dockerfile -t cognitivexr/cogstream-mediator

Run the docker container using the container engines (currently only works reliably with host mode because of engine addressing)

    docker run --rm -it \
        --network=host \
        -v $(pwd)/engines/engines-docker:/cogstream/engines \
        -v /var/run/docker.sock:/var/run/docker.sock \
        cognitivexr/cogstream-mediator:latest

### Python engines

Python engines living in `engines/engines-py` that have a `.spec.json` file can also be loaded by the mediator.
To do that, navigate into a specific engine, create a virtual environment, and install the dependencies.

For example, to build the fermx engine, run:
```bash
cd engines/engines-py/fermx
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
```

Then, when running the mediator, you can add the directory as `-plugins` argument.

```bash
dist/mediator -plugins dist/engines -plugins engines/engines-py
```

### Docker Engines

We can package engines as Docker containers, and then use that as an abstraction to serve them.
For the python engines, there are instructions on how to build the container images.
There are both Dockerfiles to create normal images, as well as ones to create GPU-accelerated containers using [nvidia/cuda](https://hub.docker.com/r/nvidia/cuda) as base.

The engine plugin files to start engines as docker containers can be found in `engines/engines-docker`.

### Clients

To build and run the python client, run `cd clients/python` and execute `make venv`.
Run `source .venv/bin/activate` to start the virtual environment.

#### Record a video

    python -m cogstream.cli.client --operation record

to start a client that starts recording a video using the recorder engine.
The videos are currently stored into `/tmp`.

#### Interactive client

The interactive python client allows you to select an engine and automatically stream the default camera to it

    $ python -m interactive.main
    0: EngineSpec(name='debug', attributes={'format.colorMode': ['0'], 'format.height': ['0'], 'format.width': ['0']})
    1: EngineSpec(name='fermx', attributes={'format.colorMode': ['0'], 'format.height': ['0'], 'format.width': ['0']})
    2: EngineSpec(name='record', attributes={'format.colorMode': ['0'], 'format.height': ['0'], 'format.width': ['0']})
    which engine do you want to use?: 2


Handshake
---------

Each connection of a CogStream client with an engine starts with a websocket connection to the mediator.
Inside the websocket connection 2 message exchanges happen, thus performing the handshake.
More information about the handshake can be found [here](https://github.com/cognitivexr/CogStream/tree/main/mediator).

Streaming protocol
------------------

Streaming is initiated between a client and an engine by first sending the negotiated `StreamSpec`, serialized as UTF-8 encoded JSON.
The packet is prefixed with an uint32 (4 byte little endian) to indicate the string length.

### Frames

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

### Engine results

Results of analytics engines are transported in `ResultPacket` instances that are structured in the same way only without the metadata field:

    +----------+------------------------+
    | Offset   | Field                  |
    +----------+------------------------+
    |        0 | Stream Id              | HEADER (little endian uint32 fields)
    |        4 | Frame Id               |
    |        8 | Unixtime seconds       |
    |       12 | Unixtime nanoseconds   |
    |       16 | Data length (L_d)      |
    +----------+------------------------+
    |       20 | Data                   | BODY
    +----------+------------------------+

The engine results data will have different formats, currently they are JSON encoded documents.
