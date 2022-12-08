# vectorx
This project is to develop new voice commands and features for Vector. It runs on top of Wirepod setups.
The goals of this project:
- Add advanced voice commands to give Vector new capabilities. Much more than a feature friday!
- These commands are implemented via the GO SDK on the server where Wirepod runs. This means they do
not require any changes to Vector firmware. They work *with any production Vector that is able to be hooked
to Wirepod*. **You don't need OSKR.**

## Setup
1. Install Wirepod (https://github.com/kercre123/wire-pod). In the setup procedure, choose VOSK as a TTS
   engine since VectorX aims to fully support localization in different languages, and VOSK is currently the
   only TTS engine that enables Wirepod to do that.
2. Enable Wirepod as a system service and reboot.
3. Run setup.sh. It basically just asks you where you installed Wirepod.
4. That's it. To test that everything works fine, just try:
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

   var intent = IntentDef{
   IntentName: "extended_intent_hello_world",
   Utterances: utterances,
   Parameters: []string{},
   Handler:    helloWorld,
   }
   *intentList = append(*intentList, intent)
   return nil 
}
```   

Where you declare that you are adding a new intent named **extended_intent_hello_world** that is triggered
by the utterances **"hello world"** in English and **"ciao mondo"** in Italian. The intent will be handled
by a function called **helloWorld**. 

2) Next, you'we got to write this handler function:
```
func helloWorld(intent IntentDef, params IntentParams) string {
	returnIntent := STANDARD_INTENT_GREETING_HELLO
	sdk_wrapper.SayText("Hello world!")
	return returnIntent
}
```   

Here you program what Vector does when the intent is matched. The Robot will simply say "Hello world" and
return to Wirepod the standard intent **STANDARD_INTENT_GREETING_HELLO** to have Vector play its stock 
greeting animation.

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

