## Welcome to vectorx...
This project is to develop new voice commands and features for Vector (VECTORX stands for VECTOR eXtended). 
It runs on top of Wirepod setups. The goals of this project:
- Add advanced voice commands to give Vector new capabilities. Much more than a feature friday!
- These commands are implemented via the GO SDK on the server where Wirepod runs. This means they do
not require any changes to Vector firmware. They work *with any production Vector that is able to be hooked
to Wirepod*. **You don't need OSKR.**

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
  they are created by **kosonicon** and reworked a bit to give them the proper Vector look.
  
All commands are fully localized in the 5 supported languages. Feel free to help with the trenslations, that can surely be improved! 

## Setup
1. Install Wirepod (https://github.com/kercre123/wire-pod). In the setup procedure, choose VOSK as a TTS
   engine since VectorX aims to fully support localization in different languages, and VOSK is currently the
   only TTS engine that enables Wirepod to do that.
2. Enable Wirepod as a system service and reboot.
   3. Run setup.sh. It basically just asks you where you installed Wirepod and the openweathermap.org free 
   API key if you want to use the extended weather forecast (you'd better do)
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
