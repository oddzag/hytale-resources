# How to create a server plugin

This guide is based on [this server setup](https://github.com/oddzag/hytale-resources/blob/main/guides/how-to-dedicated-server-debian-13.md). It is a simple server plugin that creates a usable in-game command, in this case `/test`. Start by making a working directory, then clone the template, and configure it for your project.
```
mkdir -p /opt/hytale/dev/Hytale-Example-Project
cd /opt/hytale/dev/Hytale-Example-Project
git clone https://github.com/jrilez/Hytale-Example-Project .
```
Next, make `gradlew` executable and then build the plugin
```
chmod +x /opt/hytale/dev/Hytale-Example-Project/gradlew
./gradlew shadowJar
```

Now copy the built jar file to your `mods` directory and restart your server
```
cp /opt/hytale/dev/Hytale-Example-Project/build/libs/ExamplePlugin-0.0.2.jar /opt/hytale/Homeworld/mods
sudo systemctl restart hytale-Homeworld.service
```

You can monitor the server's console output to look for the plugin's start message
```
sudo journalctl -u hytale-Homeworld.service -f
```
```
Jan 24 00:32:58 hytale java[8202]: [2026/01/24 05:32:58   INFO]                [ExamplePlugin] Setting up plugin Example:ExamplePlugin
```

And finally, if you log into the server, you can use the `/test` command and you'll get a response:
```
Hello from the Example:ExamplePlugin v0.0.2 plugin!
```

### Configuration

For some basic personalization: set your `version`, `maven_group`, and `hytale_home` path in `gradle.properties`


`/opt/hytale/dev/Hytale-Example-Project/gradle.properties`
```
# The current version of your project. Please use semantic versioning!
version=0.0.2

# The group ID used for maven publishing. Usually the same as your package name
# but not the same as your plugin group!
maven_group=org.oddzag

# The version of Java used by your plugin. The game is built on Java 21 but
# actually runs on Java 25.
java_version=25

# Determines if your plugin should also be loaded as an asset pack. If your
# pack contains assets, or you intend to use the in-game asset editor, you
# want this to be true.
includes_pack=true

# The release channel your plugin should be built and ran against. This is
# usually release or pre-release. You can verify your settings in the
# official launcher.
patchline=release

# Determines if the development server should also load mods from the user's
# standard mods folder. This lets you test mods by installing them where a
# normal player would, instead of adding them as dependencies or adding them
# to the development server manually.
load_user_mods=false

# If Hytale was installed to a custom location, you must set the home path
# manually. You may also want to use a custom path if you are building in
# a non-standard environment like a build server. The home path should
# the folder that contains the install and UserData folder.
hytale_home=/opt/hytale/Homeworld
```

Now update your manifest.json to point to the correct main class
`/opt/hytale/dev/Hytale-Example-Project/src/main/resources/manifest.json`
```
{
    "Group": "Example",
    "Name": "ExamplePlugin",
    "Version": "0.0.2",
    "Description": "An example plugin for HyTale!",
    "Authors": [
        {
            "Name": "It's you!"
        }
    ],
    "Website": "example.org",
    "ServerVersion": "*",
    "Dependencies": {
        
    },
    "OptionalDependencies": {
        
    },
    "DisabledByDefault": false,
    "Main": "org.example",
    "IncludesAssetPack": true
}
```
