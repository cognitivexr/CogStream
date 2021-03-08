import struct


def send_packet(sock, data: bytes):
    """
    Sends the given data over the socket and prefixes the packet with an appropriate packet header.

    :param sock: the socket to send the data over
    :param data: the data to send
    :return:
    """
    # Prefix each message with a 4-byte length (little endian)
    arr = struct.pack('<i', len(data)) + data
    sock.sendall(arr)


def serialize_frame(img: 'numpy.ndarray'):
    return img.tobytes()
