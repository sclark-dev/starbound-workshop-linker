# Starbound Workshop Linker
This is a small tool to help automate the process of setting up a modded Starbound server. Due to how the workshop is implemented, running the dedicated Starbound server will not automatically load the mods. You must manually place them in the mods folder, which can be very time consuming. To fix this, I created this program.

### Usage
To use this tool, just download the executable for your platform, and run it via command line. It will provide all the information you need to get it working.

This is an example that deploys symlinks between all workshop files, and the server mods folder.
```bash
starbound-workshop-linker symlink --workshop "/media/Game Drive/SteamLibrary/steamapps/workshop/content/211820/" --server "/home/admin/starbound_server/mods"
```

#### Symlink vs Copy
Symlink is the better option if you're running the server on the same machine you are playing Starbound on. It also provides the benefit of server mods stay updated.

Copying is the better option if you're running the server on a dedicated host, or through docker. Due to different configurations, there isn't an easy way to automate it. The best solution is to just copy all mods into the server's mods folder.

### Building
To build this from source, you'll need Go 1.18 installed.

### Attributions
Thanks to Ledhead on the Starbound Forums for SWEL, which provided inspiration for this project.

### FAQs
* If this was inspired by a batch script, why rewrite it in something like Go?

Well, the reasons are simple. I like Go, and I dislike Bash. Another benefit of Go, is the possibility of more complex logic to do more things. Go also gives the benefit of one codebase for multiple platforms which is the bane of all developers.

* Is there going to be a GUI version?
  
Maybe. I am considering making a crossplatform alternative to PenGUIn, and if I do, I will definitely be adding Steam workshop linking as a feature.