#!/usr/bin/python3
import os
import errno


if __name__ == "__main__":
    IPC_FIFO_NAME = "mypipefile"

    try:
        os.mkfifo(IPC_FIFO_NAME)
    except OSError as oe:
        if oe.errno != errno.EEXIST:
            raise

    npipe = os.open(IPC_FIFO_NAME, os.O_WRONLY)
    try:
        while True:
            name = input("Enter a string:")
            content = f"{name}\n".encode("utf8")
            os.write(npipe, content)
    except KeyboardInterrupt:
        print("\nGoodbye!")
    finally:
        os.close(npipe)
