# Dark Souls Autosave
This is a command line autosaver for Dark Souls  I/II/III <br />
When launched, this program stores a backup of every save done by the game.<br />
You then get the ability to revert to previous state of you character,
especially useful in case of his death.
You must exit your game before loading the desired save file.<br />
In case the file you loaded didn't satisfy you, you may undo your load.<br />

## Installation
If you're familiar with Go, you can compile the latest version using
```console
go build
```
Otherwise, just get the latest version from the /dist directory

## Usage
Place the executable file inside Dark Souls 2 saves directory.<br />
By default:
```console
C:\Users\<Username>\AppData\Roaming\DarkSoulsII\<SomeRandomString>
```
Then launch it.

After your character dies, you the game will also save.<br />
You should load 3-4 saves back, if you wish to recover you status.

## Licence
MIT

