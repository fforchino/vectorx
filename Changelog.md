RELEASE_12
- Add an uninstall script (uninstall.sh)

RELEASE_11
- Standard/Extended voice command speed improvement. Now VectorX code is compiled into binary code so it runs much faster. Also, the SDK is initialized only in case of a VectorX command, so standard wirepod voice commands are not delayed.
- Brand new UI for the VectorX control panel, with all VectorX features to try
- Fix OpenCV service not starting correctly
- Help on VectorX voice commands
- RPi4 Image is now updatable. "Check for updates" checks if updates are available and downloads them
- Bug fixes in Go SDK (for example, custom eye colors were not correctly handled when drawing monochrome graphics)

RELEASE_10
- Increase VectorX performance by building the go code to binary. Also initialize go SDK only if there is
  an intent match. Existing VectorX users should re-run setup.sh in order to have the code compiled.

RELEASE_09
- Added webserver for initial setup, runs on http://escapepod.local:8070 

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