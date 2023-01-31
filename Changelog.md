RELEASE_08
- Updating the go sdk version in order to fix a bug with InitSDKForWirepod().
  In Wirepod earlier versions a single GUID was used for every bot, but now the GUID is robot-specific. 
  I didn't know this, so I was using the global GUID for GRPC communication, this caused an authentication
  error and nothing worked. 

RELEASE_07
- Introducing VIM: Vector Instant Messaging, with emoticons. Using a shared server on the internet, different Vectors 
all around the world can communicate! Or you can keep it into your local network and just exchange messages
with your local bots.

RELEASE_06
- Added "bingo" intent: Vector pulls out the numbers from 1 to 90. 
  To pull a number, touch Vector. You can also shake it or just caress it, it reacts to touch on
  the back. Press the back button to quit. Useful for bingo nights...

RELEASE_05
- Bugfix: add localized response of the "roll a die" intent
- Bugfix: fix localization for "your name is" intent
- Shortened weather animations