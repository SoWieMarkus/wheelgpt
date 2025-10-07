# WheelGPT-Plugin

WheelGPT is a Twitch bot that was originally developed for the Twitch channel of GranaDyy - the wheel guy. With the bot, viewers can ask for the current map of the streamer. They can also guess the next personal best time of the streamer. To automate the evaluation of the guesses and the updating of the map, this plugin was developed.

## Features

### Token

To use this plugin, you need a token. We are currently working on a system with which you can automatically register for the bot via Twitch and receive a token. As long as this system isn't finished, you are welcome to contact me on discord (sowiemarkus). The token assigned to you is sent with every request to our server and identifies your Twitch channel. So don't share this token with anyone. Otherwise, others can use the plugin under your name and, for example, send best times to your Twitch channel that you have never driven.

The token assigned to you must be entered in the `Settings` -> `WheelGPT-Plugin` -> `Info/Token`.

### Personal Best

With the Twitch Bot, your viewers can guess what your next best time will be with `!guess <time>`. The format of the time is `hours:minutes:seconds.milliseconds`. `Hours` and `minutes` are optional. Example: `!guess 12.345` or `!guess 01:23.456`.

If the plugin is activated, whenever you set a new best time on a map, this time will be sent to our server. Afterwards all received guesses of your users will be checked and in the bot will announce in the chat which viewer was closest with the guess.

If you want to deactivate that your best times are sent to our server, you can deactivate this under `Settings` -> `WheelGPT-Plugin` -> ` Send PBs to Server`. To deny users to use the `!guess` command you can disable the time module with `!wgpt-disable time` in your Twitch chat.

### Map

With the Twitch bot, viewers can use the command `!map` to find out which map you are currently playing. 
If the plugin is activated, whenever you enter a new map, this map is sent to our server. If you use the `Champion Medals` plugin, the Champion Medal time will also be sent to the server.

If you want to disable this feature, you can deactivate this under `Settings`->`WheelGPT-Plugin`->` Send Maps to Server`. To deny users to use the `!map` command you can disable the map module with `!wgpt-disable map` in your Twitch chat.

## Any Questions? Feedback?

If you have any questions regarding the plugin or the Twitch Bot, feel free to contact me!

- Twitter: https://twitter.com/sowiemarkus
- Discord: sowiemarkus
