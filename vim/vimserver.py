#!/usr/bin/python
"""
pip install requests-toolbelt
pip install numpy

curl -F "image=@test.jpg" localhost:8091

"""
import json
import datetime
import os
from urllib import parse
import cv2
import mediapipe as mp
from PIL import Image
import threading
from http.server import BaseHTTPRequestHandler, HTTPServer
from socketserver import ThreadingMixIn
from io import StringIO
import time
import numpy as np

class VIMHandler(BaseHTTPRequestHandler):
    def do_POST(self):
        length = int(self.headers.get('content-length'))
        field_data = self.rfile.read(length)
        fields = parse.parse_qs(str(field_data,"UTF-8"))
        if fields["op"]=="register":
            jsonData = self.registerChatBot("","")
        elif fields["op"]=="message":
            msgFrom = fields["from"]
            msgTo = fields["to"]
            msgBody = fields["body"]
            jsonData = self.sendChatMessage(msgFrom, msgTo, msgBody)
        else:
            jsonData = "{\"status\": \"error\" }"

        #Send response
        response = bytes(jsonData, 'utf-8')
        self.send_response(200) #create header
        self.send_header("Content-Length", str(len(response)))
        self.send_header('Content-type', 'application/json')
        self.end_headers()
        self.wfile.write(response)

    def registerChatBot(self, botName, botIP, botPort):
        return "{\"status\": \"ok\" }"

    def sendChatMessage(self, msgFrom, msgTo, msgBody):
        return "{\"status\": \"ok\" }"

class ThreadedHTTPServer(ThreadingMixIn, HTTPServer):
	"""Handle requests in a separate thread."""

def main():
	try:
		server = ThreadedHTTPServer(("0.0.0.0", 8090), VIMHandler)
		print("VIM server started")
		server.serve_forever()
	except KeyboardInterrupt:
		server.socket.close()

if __name__ == '__main__':
	main()

