CogStream
=========

The CogStream video analytics system consists of three components:

* `mediator`: server component that mediates a handshake between a client and the streaming system
* `engines`: engines are services that perform video stream processing tasks
* `client`: client code for the handshake and streaming protocol

Handshake
---------

TODO: document handshake

Streaming protocol
------------------

A packet consists of:

* 4 byte little endian unsigned integer to encode the packet length
* the image encoded as 8bit uint array
