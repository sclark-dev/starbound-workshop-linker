# Starbound Workshop Linker
This is a small tool to help automate the process of setting up a modded Starbound server. Due to how the workshop is implemented, running the dedicated Starbound server will not automatically load the mods. You must manually place them in the mods folder, which can be very time consuming. To fix this, I created this program.

### Usage
To use this tool, just download the executable for your platform and run is via terminal. All features are documented via the `-h` argument, however the most common uses are listed below. The examples given are for Linux and MacOS. Windows is done largly the same, just with a different executable name and using backslashes.

#### Symlink
Symlinking is the preferred method if the server will be running on the same machine as you play on. This method provides automatic updates for all mods.

```bash
starbound-workshop-linker symlink --workshop "/media/Game Drive/SteamLibrary/steamapps/workshop/content/211820/" --server "/home/admin/starbound_server/mods"
```

#### Copy
Copying is the preferred method is you're running the server on a dedicated host, or in a docker/k8s container. Due to the more varied configurations, a one size fits all solution isn't possible. This method requires you to unlink and copy all files to update mods.

```bash
starbound-workshop-linker copy --workshop "/media/Game Drive/SteamLibrary/steamapps/workshop/content/211820/" --server "/home/admin/starbound_server/mods"
```

#### Unlink
Unlinking removes all symlinks or copies of Starbound mods found in the specified folder. This is required when adding new mods.

```bash
starbound-workshop-linker unlink --workshop "/media/Game Drive/SteamLibrary/steamapps/workshop/content/211820/" --server "/home/admin/starbound_server/mods"
```

### Building
To build this from source, you'll need Go 1.18 installed.

### Attributions
Thanks to Ledhead on the Starbound Forums for SWEL, which provided inspiration for this project.

### FAQs
* If this was inspired by a batch script, why rewrite it in something like Go?

Well, the reasons are simple. I like Go, and I dislike Bash. Another benefit of Go, is the possibility of more complex logic to do more things. Go also gives the benefit of one codebase for multiple platforms which is the bane of all developers.

* Is there going to be a GUI version?
  
Maybe. I am considering making a crossplatform alternative to PenGUIn, and if I do, I will definitely be adding Steam workshop linking as a feature.

* I have an idea!

Ideas are welcome. If you can write Go, feel free to open a Pull Request. If you cannot you can always open an issue and suggest your idea there.