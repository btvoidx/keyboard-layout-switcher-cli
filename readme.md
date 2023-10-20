# Windows keyboard layout switcher CLI
## The problem
I love [helix](https://helix-editor.com). It is truly one of the editors of all time.
But it has a slight problem: the controls are *english-only*. As a bilingual person, this irritates me, as I have to press
shift+alt like 3 times a second when trying to edit some non-english text.

## The solution
Just make a simple to use CLI to switch keyboard layout and integrate it with
[helix](https://helix-editor.com) by overriding mode switching keybindings to include
calls to said CLI, like so:
```toml
[keys.insert]
esc = [
  "normal_mode",
  ":sh lang get -f C:/Users/%username%/.hxlang",
  ":sh echo 00000409 | lang set",
]

[keys.normal]
i = [":sh lang set -f C:/Users/%username%/.hxlang", "insert_mode"]
a = [":sh lang set -f C:/Users/%username%/.hxlang", "append_mode"]
o = [":sh lang set -f C:/Users/%username%/.hxlang", "open_below"]
O = [":sh lang set -f C:/Users/%username%/.hxlang", "open_above"]
```
> Overriding select mode keybindings is an exercise for the reader.

## I want it too!
Grab your favourite version of Go since 1.18 and run a
```
go install github.com/btvoidx/keyboard-layout-switcher-cli
```
The CLI will be installed as `lang.exe` into your **%GOPATH%/bin**.

Note that the tool is windows only, as windows is the OS I use.
It will not compile for any other OS.
