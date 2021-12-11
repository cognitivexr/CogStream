import argparse
import logging
from collections import defaultdict

import cv2
from cogstream.api import OperationSpec, ClientFormatSpec, to_attributes
from cogstream.engine import EngineResult
from cogstream.engine.client import EngineClient, stream_camera
from cogstream.mediator.client import MediatorClient
from rich.console import Console
from rich.prompt import IntPrompt

console = Console()


def parser():
    p = argparse.ArgumentParser(description="CogStream Client for recording videos")
    p.add_argument(
        "--mediator-host", type=str, help="the mediator address", default="127.0.0.1"
    )
    p.add_argument(
        "--mediator-port",
        type=int,
        help="the mediator port (default 8191)",
        default=8191,
    )
    p.add_argument(
        "--capture-width", type=int, help="camera capture height", default=800
    )
    p.add_argument(
        "--capture-height", type=int, help="camera capture width", default=600
    )
    p.add_argument(
        "-ot",
        "--operation-type",
        type=str,
        help="engine operation type (record|expose|analyze)",
        default="analyze",
    )
    p.add_argument(
        "-a",
        "--attribute",
        help="key=value attribute for the operation spec",
        action="append",
    )
    return p


def main():
    logging.basicConfig(level=logging.INFO)

    args = parser().parse_args()

    # parse attributes
    attributes = defaultdict(list)
    if args.attribute:
        for kv in args.attribute:
            k, v = kv.split("=")
            attributes[k].append(v)

    mediator = MediatorClient(args.mediator_host, args.mediator_port)

    op_spec = OperationSpec(args.operation_type, to_attributes(attributes))
    available_engines = mediator.request_operation(op_spec)

    if not available_engines.engines:
        console.print(
            f"[red]error[/red]: no available engines for {args.operation_type} and attributes {dict(attributes)}"
        )
        exit(126)

    # create engine select prompt
    for i, engine in enumerate(available_engines.engines):
        console.print(f"{i}: {engine}")

    choice = IntPrompt.ask("which engine do you want to use?")
    print(choice)

    if choice >= len(available_engines.engines) or choice < 0:
        console.print(f"[red]error:[/red] index {choice} not in range")

    selection = available_engines.engines[choice]

    # create client format spec
    # TODO: get from available camera args
    client_format = ClientFormatSpec(
        selection.name,
        to_attributes(
            {
                "format.width": args.capture_width,
                "format.height": args.capture_height,
                "format.colorMode": "RGB",
            }
        ),
    )
    stream_spec = mediator.establish_format(client_format)
    mediator.close()

    engine_client = EngineClient(stream_spec)
    engine_client.open()

    if not engine_client.acknowledged:
        print("engine stream could not be initiated")
        engine_client.close()
        return

    cap = cv2.VideoCapture(0)

    try:
        cap.set(cv2.CAP_PROP_FRAME_WIDTH, args.capture_width)
        cap.set(cv2.CAP_PROP_FRAME_HEIGHT, args.capture_height)

        stream_camera(cap, engine_client, on_result=print_result)
    except KeyboardInterrupt:
        pass
    finally:
        cap.release()
        engine_client.close()


def print_result(ndarr, engine_result: EngineResult):
    console.log(engine_result)


if __name__ == "__main__":
    main()
