import socket
import json
import os

os.unlink("/tmp/example.sock")

def main():
    sock = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
    sock.bind("/tmp/example.sock")
    sock.listen(1)
    while True:
        conn, addr = sock.accept()
        try:
            f = conn.makefile('rw')
            for ln in f:
                obj = json.loads(ln)
                print(obj)
        except ValueError, e:
            print("ValueError", e)
        finally:
            conn.close()

main()
