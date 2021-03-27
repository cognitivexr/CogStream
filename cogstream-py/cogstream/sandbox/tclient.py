import socket

from cv2 import cv2

from cogstream.engine.channel import JpegSendChannel, Frame
from cogstream.engine.io import SocketFrameWriter


def main():
    address = ('localhost', 54321)

    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    cam = cv2.VideoCapture(0)

    channel = JpegSendChannel(12345, SocketFrameWriter(sock))

    try:
        sock.connect(address)

        while True:
            check, frame = cam.read()
            if not check:
                break

            channel.send(Frame(frame))
            cv2.imshow("capture", frame)

            key = cv2.waitKey(1)
            if key == ord('q'):
                break

    except KeyboardInterrupt:
        pass
    finally:
        sock.close()


if __name__ == '__main__':
    main()
