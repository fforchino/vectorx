#!/usr/bin/python
"""
pip install requests-toolbelt
pip install numpy

curl -F "image=@test.jpg" localhost:8090

"""
import json
import datetime
import os
import cv2
import mediapipe as mp
from PIL import Image
import threading
from http.server import BaseHTTPRequestHandler, HTTPServer
from socketserver import ThreadingMixIn
from io import StringIO
import time
import numpy as np
from requests_toolbelt.multipart import decoder

class VideoHandler(BaseHTTPRequestHandler):
    def do_POST(self):
            content_length = int(self.headers['Content-Length'])

            #get data content bytes
            file_content = self.rfile.read(content_length)

            #Use multipart parser to strip boundary
            multipart_data = decoder.MultipartDecoder(file_content, self.headers['Content-Type']).parts
            image_byte = multipart_data[0].content
            #Read image using cv2
            image_numpy = np.frombuffer(image_byte, np.int8)
            img = cv2.imdecode(image_numpy, cv2.IMREAD_UNCHANGED)
            """
            pref = "IMAGE" + datetime.datetime.now().strftime("%y%m%d_%H%M%S") + ".jpg"
            filedir = "/tmp"
            filename = os.path.join(filedir, pref)
            cv2.imwrite(filename, img)
            """
            jsonData = self.detectFingers(img)      
            print(jsonData)

            #Send response
            response = bytes(jsonData, 'utf-8')
            self.send_response(200) #create header
            self.send_header("Content-Length", str(len(response)))
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            self.wfile.write(response)
    
    def detectFingers(self, image):
        debug = False
        mp_hands = mp.solutions.hands
        handFacing = ""
        handLabel = ""
        upCount = index_x = index_y = -1

        with mp_hands.Hands(static_image_mode=True, max_num_hands=1, min_detection_confidence=0.5) as hands:
            # Tke the image, flip it around y-axis for correct handedness output
            image = cv2.flip(image, 1)
            image_height, image_width, _ = image.shape
            # Convert the BGR image to RGB before processing.
            results = hands.process(cv2.cvtColor(image, cv2.COLOR_BGR2RGB))

            if (debug) : 
                print('Handedness:', results.multi_handedness)
            
            if results.multi_hand_landmarks != None :
                hand_landmarks = results.multi_hand_landmarks[0]                
                if (debug) : 
                    print('hand_landmarks:', hand_landmarks)
        
                if results.multi_hand_landmarks:
                    handList = []
                    upCount = 0 
                    for handLms in results.multi_hand_landmarks:
                        for idx, lm in enumerate(handLms.landmark):
                            h, w, c = image.shape
                            cx, cy = int(lm.x * w), int(lm.y * h)
                            handList.append((cx, cy))
                        
                    finger_Coord = [(8, 5), (12, 9), (16, 13), (20, 17)]
                    thumb_Coord = (4,1)
                    
                    i = 0
                    upCount = 0
                    for coordinate in finger_Coord:
                        #Calculate the minimum finger length we are allowing to call it "raised"
                        end = coordinate[1]
                        start = end+1
                        minFingerLen = (handList[end][1] - handList[start][1])*80/100
                        fingerLen = handList[start][1] - handList[coordinate[0]][1]
                        if (debug) :
                            print("Start("+str(start)+"): "+str(handList[start][1])+",end ("+str(end)+"):"+str(handList[end][1])) 
                            print("Min. finger length for finger "+str(i)+": "+str(minFingerLen)+", found:"+str(fingerLen))
                        if minFingerLen>0 and fingerLen>minFingerLen:
                            if (debug) :
                                print("   ->"+str(i)+" finger up: "+str(handList[start][1])+"<"+str(handList[coordinate[0]][1])+", fingerlen = "+str(fingerLen)+" (min:"+str(minFingerLen)+")")
                            upCount += 1
                        i = i+1
                    
                    #Thumb processing
                    end = thumb_Coord[1]
                    start = end+1
                    
                    handIndex = results.multi_hand_landmarks.index(hand_landmarks)
                    handLabel = results.multi_handedness[handIndex].classification[0].label

                    #TODO: We have also to understand whether the hand is reversed
                    handFacing = "reverse"
                    if (handList[17][0]>handList[5][0]):
                        handFacing = "front"
                        
                    
                    if handLabel == "Left":
                        minThumbLen = (handList[end][0] - handList[start][0])*80/100
                        if handFacing == "front":
                            thumbLen = handList[4][0] - handList[2][0]
                        else:
                            thumbLen = handList[2][0] - handList[4][0]
                    else:  
                        minThumbLen = (handList[start][0]-handList[end][0])*80/100
                        if handFacing == "front":
                            thumbLen = handList[2][0] - handList[4][0]
                        else: 
                            thumbLen = handList[4][0] - handList[2][0]
                            
                    if (debug) :
                        print("Hand: "+handLabel+" facing: "+handFacing)
                        print("Start("+str(start)+"): "+str(handList[start][0])+",end ("+str(end)+"):"+str(handList[end][0])) 
                        print("Min. thumb size: "+str(minThumbLen)+" found:"+str(thumbLen))

                    if thumbLen>minThumbLen:
                        if (debug) :
                            print("Thumb up!")
                        upCount += 1
                    
                    #Index finger coordinates
                    #index_x = hand_landmarks.landmark[mp_hands.HandLandmark.INDEX_FINGER_TIP].x * image_width
                    #index_y = hand_landmarks.landmark[mp_hands.HandLandmark.INDEX_FINGER_TIP].y * image_height
                    index_x = handList[8][0]
                    index_y = handList[8][1]
                    #print("IMAGE: "+str(image_width)+"x"+str(image_height))
                    #print("INDEX: "+str(index_x)+"x"+str(index_y))
                
        jsonData =  "{ \"handedness\": \""+str(handLabel)+"\", "
        jsonData += "\"facing\": \""+str(handFacing)+"\", "
        jsonData += "\"raisedfingers\": "+str(upCount) + ", "
        jsonData += "\"index_x\": "+str(index_x) + ", "
        jsonData += "\"index_y\": "+str(index_y) + ", "
        jsonData += "\"index_z\": "+str(upCount) + " }"
        return jsonData  


class ThreadedHTTPServer(ThreadingMixIn, HTTPServer):
	"""Handle requests in a separate thread."""

def main():
	global capture
	global img
	try:
		server = ThreadedHTTPServer(("0.0.0.0", 8090), VideoHandler)
		print("server started")
		server.serve_forever()
	except KeyboardInterrupt:
		server.socket.close()

if __name__ == '__main__':
	main()

