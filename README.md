# Iceberg 🧊

![](https://i.imgur.com/bn7PSv0.png)

Iceberg is a Discord bot created for the BlockHermit server. It has music, name verification, tools related to the BlockHermit Minecraft server and more fun stuff.

Note: this is the rewrite and is still in development. Feel free to contribute as below.

## Contributing ✨
Any contributions are welcome. From grammar corrections to entirely new features or just prettifying some commands.

> If you plan to add new features (commands, subcommands, something small, etc.) **PLEASE** [create an issue](https://github.com/lieuweberg/iceberg/issues/new) first. Give it a fitting title and a nice description explaining what it will do and why it would be good to add. It would be a bummer if you make something but then it's not even added.

## Running the bot 🚒
 1. Firstly, install [Go](https://golang.org/dl/). Make sure you download a version that is fairly new (1.14+) just so it all works.
 2. Make sure you have a Java version installed that is higher than 1.11. There currently is no way to run the bot without music, and thus without Java.
 3. Also install a GCC compiler if you don't have one already (try running `gcc` in a command prompt).
    - You can use [MinGW](https://sourceforge.net/projects/tdm-gcc/) on Windows for easy installation management. In MinGW, select recommended for C/C++. Hit install. You can later return here to easily delete it later if you want. Just look for "MinGW" in the start menu.
 4. Open a command line, navigate to a desired folder (e.g. desktop) and clone this repo:
 ```
 git clone https://github.com/lieuweberg/iceberg.git && cd iceberg
 ```
 4. Install dependencies:
 ```
 go get
 ```
 5. Rename `config.example.json` to `config.json` and add a bot token. See [this](https://raw.githubusercontent.com/denverquane/amongusdiscord/master/BOT_README.md) for how to obtain one and invite a bot to a server. Make sure to do this in a private/test server.
 6. Install Lavalink, the bot's music player, from [here](https://ci.fredboat.com/viewLog.html?buildId=lastSuccessful&buildTypeId=Lavalink_Build&tab=artifacts&guest=1#). Click `Lavalink.jar` on the left and save it to the `.lavalink` folder. Run
 ```
 java -jar ./lavalink/Lavalink.jar
 ```
 7. Open a new command line and run the bot (always start Lavalink before the bot):
 ```
 go run main.go
 ```

 ## License 🔑
 See the `LICENSE` file
