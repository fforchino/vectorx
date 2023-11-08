## Welcome to vectorx...
This project is to develop new voice commands and features for Vector (VECTORX stands for VECTOR eXtended). 
It runs on top of Wirepod setups. The goals of this project:
- Add advanced voice commands to give Vector new capabilities. Much more than a feature friday!
- These commands are implemented via the GO SDK on the server where Wirepod runs. This means they do
not require any changes to Vector firmware. They work *with any production Vector that is able to be hooked
to Wirepod*. **You don't need OSKR.**

**NOTE:** VectorX has been tested with Vector 1.0 only. Vector 2.0 may not support all needed SDK commands. If you'd like
to fund Vector 2.0 development, contact me. I don't currently have any Vector 2.0 to try.

## Supported commands
Hey Vector...
- **Hello world!**
  - Simple test command. Vector will reply "Hello world" in the current STT language (the one Chipper is using to
  parse commands)
- **You are < NAME >**
  - Sets a special name for the robot. Vector will remember its name. Unfortunately this cannot be used to change 
  the wake-word, i.e. Vector will always react only to "Hey Vector"
- **What's your name?**
  - Vector will tell its name if any
- **Let's talk < LANGUAGE >**
  - The STT language will be switched to <LANGUAGE> (one of english, italian, spanish, french or german). 
  It also accidentally plays a very short clip of the chosen language anthem
- **How do you say < WORD > in < LANGUAGE >**
  - Vector will use Google translate to translate the word and speak the outcome. You can also try with simple sentences
  The supported languages are those used for STT
- **Roll a die!**
  - Vector will roll a D6 die and tell you the result
- **What's the weather?**
- **What's the weather in < LOCATION >?**
- **What's the weather tomorrow?**
  - Vector checks weather forecast on openweathermap.org and reports the result. This is much better than the original 
  weather command, because: 
    - It is fully localized
    - Supports more weather conditions (for example, fog)
  Well, the animations are not as funny as the original ones. The icons used have been found on FlatIcon (https://www.flaticon.com/free-icons/weather),
  they are created by **kosonicon** and reworked a bit to give them the proper Vector look. Since graphics only (though animated) was a bit lame, I added also
  some sounds (thunderstorm, rain, wind...) Still unsure what's the sound of the fog (I thought about a wolf howling :) 
- **let's play a new game**
  - Play "rocket,paper,scissors!" with Vector. He will use the camera to track your hand and guess the number of fingers. Best played with right hand.
- **let's play a classic**
  - Play Pong against Vector. You need to use the standard Vector cube. Move it up and down in front of the camera to move the left paddle correspondly.
  The game ends when you or Vector reach 9 points.
- **let's play a classic**
- **bingo**
  - Vector will act as a bingo machine. Rub Vector's back to have him pull out a random number between 1 and 90. You can go on until all 90 numbers are pulled out, 
  or press the button on Vector's back to quit. 
All commands are fully localized in the 5 supported languages. Feel free to help with the trenslations, that can surely be improved! 

## Instant messaging
Since RELEASE_19, the redesigned VIM (Vector Instant Messaging) feature is enabled by default. 
The chat server uses websockets and runs on the local network (it is part of VectorX web server).
To use this feature, it is mandatory that you give your Vector a name. This name is used as chat username.

Hey Vector...
- **Chat with [username]**
  - Sets the chat target. For example, "Chat with Filippo". All messages will be sent to user Filippo. 
- **Who are you chatting with?**
    - Queries the chat target. In case you forgot with whom this Vector is chatting.
- **Say [sentence]**
  - Sends a text message with payload [sentence] to the current chat target

When a messge is received from another user, Vector will play a tune and then say:

**[username] says: [sentence]**

In the special case that an emoticon is received, Vector will show the emoticon.
Available emoticons: angel, angry, annoyed, blue, devil, disappointed, eheh, happy, heart,
hmm, hurt, inlove, kiss, lol, nooo, ok, panic, polemic, sad, shades, sick, smile, star, surprised,
tear, tongue, wink, wow, xxx

## Setup (on Raspberry Pi4)
If you plan to use VectorX on a Raspberry Pi, I recommend you to use the pre-built VectorX RPi4 Image. 
It currently ships with VectorX v. 20.

1. Download the image here: https://www.wondergarden.app/vectorx-rpi-images/getLatestImage.php
2. Flash the image with RPI Imager on a 16-32GB SD Card. The image is for RPI4 ARM64, 8GB RAM
2. Boot. The RPI will go into AP mode (HOTSPOT). From your PC/Mobile, look for the network "VectorX Setup" and connect 
   (no password)
3. Point your web browser to http://192.168.220.1:8080. It loads a page where you can choose your home wifi network and 
   input the password
4. After that, RPI shall turn to normal mode and connect to the given wifi network. All required services are started 
   automatically.
5. Reconnect your home network and go to http://escapepod.local:8070 to run VectorX setup procedure. 
6. Follow the instructions to onboard your robots. The procedure differs for OSKR and Production bots  
7. That's it. 

## Setup (on any machine) 
1. Install My own fork of Wirepod (https://github.com/fforchino/wire-pod) and install it. Note that VectorX and Wire-Pod 
   are two independent projects and I cannot guarantee that changes in the latter may break VectorX. That's why I ask you
   to use my own fork and not the main Wire-Pod repository (https://github.com/kercre123/wire-pod).
2. In the setup procedure, choose VOSK as a TTS
   engine since VectorX aims to fully support localization in different languages, and VOSK is currently the
   only TTS engine that enables Wirepod to do that.
3. Enable Wirepod as a system service and reboot.
4. Run setup.sh. It basically just asks you where you installed Wirepod and the openweathermap.org free 
   API key if you want to use the extended weather forecast (you'd better do). VectorX also installs Python and an 
   OpenCV server used for the computer vision games that need hand tracking.
5. That's it. To test that everything works fine, just try:
   - Hi Vector! Hello world!
   Vector should reply:
   - Hello world!
   The first time it may take a long time because the dependencies of the SDK have to be downloaded.
   
## How It works
VectorX works by injecting a special system intent into Wirepod's custom intents table. This special 
intent is not editable from Wirepod's front end and matches any utterance. Therefore, VectorX basically
acts as a pre-processor for any intent. If user's utterance matches a VectorX keyphrase, VectorX voice
command is executed. This way it is possible even to override the behavior of Wirepod default commands.

## Writing a new Voice Command using Vectorx
That is meant to be very simple, to inspire anyone to contribute to the project. The file
pkg/intents/voicecommand-helloworld.go is commented to guide you through. You just need to implement 3
simple steps.
1) Write a registration function like this:
```
func HelloWorld_Register(intentList *[]IntentDef) error {
	utterances := make(map[string][]string)
	utterances[LOCALE_ENGLISH] = []string{"hello world"}
	utterances[LOCALE_ITALIAN] = []string{"ciao mondo"}
	utterances[LOCALE_SPANISH] = []string{"hola Mundo"}
	utterances[LOCALE_FRENCH] = []string{"bonjour le monde"}
	utterances[LOCALE_GERMAN] = []string{"Hallo Welt"}

	var intent = IntentDef{
		IntentName: "extended_intent_hello_world",
		Utterances: utterances,
		Parameters: []string{},
		Handler:    helloWorld,
	}
	*intentList = append(*intentList, intent)
	addLocalizedString("STR_HELLO_WORLD", []string{"hello world!", "ciao mondo!", "hola mundo!", "bonjour le monde!", "hallo welt"})

	return nil
}
```   

Where you declare that you are adding a new intent named **extended_intent_hello_world** that is triggered
by the utterances:
**"hello world"** in English 
**"ciao mondo"** in Italian
**"hola Mundo"** in Spanish
**"bonjour le monde"** in French
**"Hallo Welt"** in German
don't worry about case since all comparisons are always done in lowercase, so no matter what capitalization you use, it will work.
The intent will be handled by a function called **helloWorld**. 
Also, you may add as many **addLocalizedString** calls you want to add the command-specific multi-language text resources 
to the localization engine. In this case we add just one, for a generic answer.

2) Next, you'we got to write this handler function:
```
func helloWorld(intent IntentDef, params IntentParams) string {
	returnIntent := STANDARD_INTENT_GREETING_HELLO
	sdk_wrapper.SayText(getText("STR_HELLO_WORLD"))
	return returnIntent
}
```   

Here you program what Vector does when the intent is matched. The Robot will simply say "Hello world" in the 
current STT laguage (the one Chipper is using right now) and return to Wirepod the standard intent **STANDARD_INTENT_GREETING_HELLO** 
to have Vector play its stock greeting animation.

3) Thid step, you go into pkg/intengts/intents.go and add your registration function to **RegisterIntents()**,
like this:

```
func RegisterIntents() {
    ...
    HelloWorld_Register(&intents)
}
```   

This way the new intent is added to the supported intent list. Next you can just try right away how it works!

## Data files
Data files are stored in the **vectorfs** directory. There are three different directories there:
1) **tmp** holds temporary data
2) **nvm** holds (r/w) bot-specific data (configuration files), it is organized into different subdirectories,
one for each robot serial number
3) **data** holds (read-only) data files. For an example of how to use data files, look at the "roll a die"
example. The resource files for the roll-a-die voice commands are located under vectorfs/data/images/dice.

## Upgrading
You can update VectorX as you would with any other git repository: "git pull". However, if you are not a developer
and thus you are not making changes to the source code, you can use the same update script that the web server does,
it will do everything for you in one shot:
```
./update.sh 
```   
The assumptions to run the script are:
1) You haven't done changes to the source code files. Else the script will fail to update VectorX git repository
2) You have run setup.sh at least once, and thus you have a valid source.sh file
3) You are connected to the Internet
 
The update script will do the following things:
1) Stop all services
2) Update VectorX and Wire-Pod (my own fork of it) repositories
3) Rebuild chipper, VectorX and its web server
4) Run again the VectorX setup in slient mode (no user input is required)
5) Start again all services

## Uninstalling
If you want to uninstall the hooks VectorX provides and remove all services, run as root the uninstall.sh 
script:
```
sudo ./uninstall.sh 
```   
You can always re-setup everything by launching setup.sh and starting over.
