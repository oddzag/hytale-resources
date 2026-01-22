I was re-setting up my server and wrote myself a guide for reference, thought maybe it might be useful. It's on a headless Debian 13 server, just for my usage. Docker is probably much easier, but I wanted to setup the Nitrado Webserver/Query plugin, as well as set up automatic backups and updates.

  
### Required
- Java
- Unzip
- UFW (for public access)


### Setup

We'll be working out of `/opt`, specifically `/opt/java` for the Java SDK and `/opt/hytale` for the server files, downloader, mods, etc. You can deviate however you see fit, this is just a general way of managing/configuring your server files. I prefer to keep the downloader separate, and then copy

```
sudo apt update
```

Start by installing Java. [HyPixel recommends](https://support.hytale.com/hc/en-us/articles/45326769420827-Hytale-Server-Manual#installing-java-25) using [Adoptium](https://adoptium.net/temurin/releases). As of writing this, it's JDK 25.0.1_8
```
sudo mkdir /opt/java
cd /opt/java
wget https://github.com/adoptium/temurin25-binaries/releases/download/jdk-25.0.1%2B8/OpenJDK25U-jdk_x64_linux_hotspot_25.0.1_8.tar.gz
sudo tar -xzf OpenJDK25U-jdk_x64_linux_hotspot_25.0.1_8.tar.gz -C .
```

Java can now be accessed at `/opt/java/jdk-25.0.1+8/bin/java`


Now create a directory for the [updater](https://support.hytale.com/hc/en-us/articles/45326769420827-Hytale-Server-Manual#hytale-downloader-cli) and download it.
```
sudo mkdir -r /opt/hytale/downloader
cd /opt/hytale/downloader
sudo wget https://downloader.hytale.com/hytale-downloader.zip
# if you need unzip
# sudo apt install unzip
sudo unzip *.zip
# remove the unneeded .exe
sudo rm -R *.exe
ls
hytale-downloader-linux-amd64  hytale-downloader.zip  QUICKSTART.md
# rename the downloader for readability
sudo mv hytale-downloader-linux-amd64 hytale-downloader
```

Now you can automatically download the server files with [`hytale-downloader`](https://support.hytale.com/hc/en-us/articles/45326769420827-Hytale-Server-Manual#hytale-downloader-cli). 

Start by making the downloader executible, then create a directory specifically for this new server, we'll call it Homeworld. Then, download the latest server/asset files for this new server and unzip them. As of writing this, they're about 1.4GB in size.

```
cd /opt/hytale/downloader
sudo chmod +x ./hytale-downloader
sudo mkdir /opt/hytale/Homeworld
sudo ./hytale-downloader -download-path /opt/hytale/Homeworld/Homeworld.zip
cd /opt/hytale/Homeworld
sudo unzip *.zip
sudo rm -R Homeworld.zip
ls
Assets.zip  Server
```

Refer to the [dedicated server manual](https://support.hytale.com/hc/en-us/articles/45326769420827-Hytale-Server-Manual) for more info on these files, but our main focus is `/opt/hytale/Homeworld/Server/HytaleServer.jar`. Once we run the server with that `jar`, it will generate our world. We'll also need to authenticate our account, and persist it.

Before we do that and create a bunch of new files, let's set up our user that will own/run the server so that it's not running as root, since we're going to configure the server as a system service. Start by creating the user and then assigning them ownership of the server files in Homeworld, as well as the java SDKs

```
# Create system user (no login, no home directory)
sudo adduser --system --group --no-create-home hytale
sudo chown -R hytale:hytale /opt/hytale/Homeworld
sudo chown -R hytale:hytale /opt/java
```

Create the system service for the Hytale server and then reload unit files, **but don't start the service yet**.

```
sudo nano /etc/systemd/system/hytale-Homeworld.service
```
```
[Unit]
Description=Hytale - Homeworld
After=network.target

[Service]
Type=simple
User=hytale
WorkingDirectory=/opt/hytale/Homeworld

# Main process
ExecStart=/opt/java/jdk-25.0.1+8/bin/java -XX:AOTCache=HytaleServer.aot -jar Server/HytaleServer.jar --assets Assets.zip

[Install]
WantedBy=multi-user.target
```
```
sudo systemctl daemon-reload
```

Before starting the service, let's run the server once to generate the world, and also set up persistent authentication for the server. You'll need to switch to the `hytale` user before starting the server
```
sudo -u hytale -s
cd /opt/hytale/Homeworld
/opt/java/jdk-25.0.1+8/bin/java -XX:AOTCache=HytaleServer.aot -jar Server/HytaleServer.jar --assets Assets.zip
```
Let this run for a second, then you should see:
```
[2026/01/21 21:37:32   WARN]                   [HytaleServer] No server tokens configured. Use /auth login to authenticate.
```
Type in `/auth login device` and hit enter. It will then generate a section with a link to authenticate your account to this server. Copy paste the link/code into your browser. Once you authenticate, the server will continue creating the session. You should now see:
```
Authentication successful! Use '/auth status' to view details.
WARNING: Credentials stored in memory only - they will be lost on restart!
To persist credentials, run: /auth persistence <type>
Available types: Memory, Encrypted
```

Type in the following to store your credentials.
```
/auth persistence Encrypted
```

First switch back to your elevated user, then we can enable and then start the service, it should authenticate automatically, and it will run in the background so you don't need the terminal opened/connected to the server. Then check the server logs with `status`
```
sudo systemctl enable hytale-Homeworld.service
sudo systemctl start hytale-Homeworld.service
sudo systemctl status hytale-Homeworld.service
```

You should see
```
Jan 21 16:47:51 hytale java[3111]: [2026/01/21 21:47:51   INFO]                   [HytaleServer] ===============================================================================================
Jan 21 16:47:51 hytale java[3111]: [2026/01/21 21:47:51   INFO]                   [HytaleServer]          Hytale Server Booted! [Multiplayer] took 11sec 819ms 394us 977ns
Jan 21 16:47:51 hytale java[3111]: [2026/01/21 21:47:51   INFO]                   [HytaleServer] ===============================================================================================
```

Since this guide is targeting a headless Debian server, I assume you're not also running the client on the server. As a result, you'll need to allow access to port 5520, that's the default server port. Install `ufw`, enable it, and allow the port (including 22 for `ssh`)
```
sudo apt install ufw
sudo apt enable ufw
sudo ufw allow 22
sudo ufw allow 5520/udp
```

Use `ip a` to get your server's IP. You should now be able to connect to the server.
