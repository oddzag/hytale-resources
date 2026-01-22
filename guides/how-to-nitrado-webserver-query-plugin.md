# How to use Nitrado's Hytale Webserver and Query plugin

Nitrado has created a [Webserver](https://github.com/nitrado/hytale-plugin-webserver) and [Query](https://github.com/nitrado/hytale-plugin-query) plugin that allows for web-based interaction with your Hytale server. This is being setup on a dedicated server, using [this server setup guide](https://github.com/oddzag/hytale-resources/blob/1428ee76ddaf65119ce02480e63bd9aae63365e7/guides/how-to-nitrado-webserver-query-plugin.md) as reference.

Once the server is started, the webserver starts automatically, running on the server's default port `5520` + 3 i.e if you make no changes, the webserver will run on port `5523`. By default, TLS is enabled, so any attempts to connect must use SSL. For my personal testing purposes, I'm disabling TLS and manually specifying the host IP. If you want to leave it enabled, you'll need to trust the certificates in your browser or however you're accessing the endpoints.

Start by stopping the server and downloading both JAR files into your server's `mods` directory. Then create the Webserver's directory (or just start/stop the server to let the app create it automatically) and then the Webserver config file. You don't **need** to create the config file, this is just if you need to make changes like the ones described above (disabling TLS/specifying host IP):

```
sudo systemctl stop hytale-Homeworld.service
cd /opt/hytale/Homeworld/mods
wget https://github.com/nitrado/hytale-plugin-webserver/releases/download/v1.0.0/nitrado-webserver-1.0.0.jar
wget https://github.com/nitrado/hytale-plugin-query/releases/download/v1.0.1/nitrado-query-1.0.1.jar
mkdir Nitrado_WebServer
# Now optionally create the Webserver's config file 
nano Nitrado_WebServer/config.json
```

This will bind the host to specific IP and disable TLS so you can make insecure GET requests to /Nitrado/Query. I'm doing this on my local network with no external access, I do not recommend doing this on a public facing server.

`/opt/hytale/Homeworld/mods/Nitrado_WebServer/config.json`
```
{
  "BindHost": "192.168.1.10",
  "Tls": {
    "Insecure": true
  }
}
```

Now start the server and you can access the WebServer UI via http://192.168.1.10:5523. If you left TLS enabled, you'd simply go to https://192.168.1.10:5523 and just accept the risk of an unrecognized self-signed certificate.
```
sudo systemctl start hytale-Homeworld.service
```

View the server's live console output
```
sudo journalctl -u hytale-Homeworld.service -f
```

Above the "Hytale Server Booted!" banner, you'll see output from the WebServer plugin that it's started and listening on port 5523:
```
Jan 22 13:38:58 hytale java[4707]: [2026/01/22 18:38:58   INFO]         [WebServer|P][WebServer] WebServer listening on 192.168.1.10:5523
```

Log into your Hytale server and use the command `/webserver code create`, it'll issue you a temporary 8-digit code that you can use to access the WebServer's GUI at http://192.168.1.10:5523. Once logged in, you can click Nitrado > Query in the left window, you'll then see the JSON payload from the `/Nitrado/Query` endpoint, which has a huge list of server information.

### Service accounts

The problem so far is that in order to access the /Nitrado/Query endpoint, you need to authenticate via the web GUI with that temporary login-code. However the WebServer plugin allows for the creation of service accounts that "are intended for processes that automatically interact with the server through HTTP APIs". I like using Go, so I'd like to simply be able to GET that endpoint and manage the JSON payload myself. 

So let's create a temporary service account for testing, and then we'll use it to access the /Nitrado/Query in the browser.
```
cd /opt/hytale/Homeworld/mods/Nitrado_WebServer/provisioning
nano bumzag.serviceaccount.json
```
```
{
  "Enabled": true,
  "Name": "serviceaccount.bumzag",
  "PasswordHash": "$2y$10$3qjK4mnB3qFo7DMRJG8ZyeOK7PcEvuBWSE1IlQWFIW.eI2AyNNzze",
  "Groups": ["OP"],
  "Permissions": ["nitrado.query.web.read.players"]
}
```

Based on Nitrado's WebServer [authentication flow](https://github.com/nitrado/hytale-plugin-webserver/blob/fddbaa09d3ec08f9a6a0ce831ee2bedc1652d46f/src/main/java/net/nitrado/hytale/plugins/webserver/authentication/store/JsonPasswordStore.java#L67), passwords are hashed using Bcrypt with a cost factor of 10. In my example above, my password is "ThisIsATestPassword". To test this yourself, try generating a test password using [this Bcrypt Hash generator](https://bcrypt.online/) and using it in your own service account config.

Save that config file and restart the server. This time, if you monitor the server's console output, you'll see:
```
Jan 22 14:12:06 hytale java[4843]: [2026/01/22 19:12:06   INFO]                    [WebServer|P] Importing service account file bumzag.serviceaccount.json
```

Now you can hit that endpoint with your desired stack. Test it with `curl` and you should get that default `JSON` payload.
```
curl -u serviceAccount.Bumzag:ThisIsATestPassword http://192.168.1.10:5523/Nitrado/Query
```
