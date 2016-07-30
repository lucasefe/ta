# ta
Tmux automation

## Automate your Tmux setup

`ta` reads your tmux window/pane definition from the .ta file located in the current directory.

## Motivation

We've found that having a tmux session open per project greatly reduces the cost of context switching: one regularly needs to open a text editor, set environment configuration, run tests, run servers, and so on, why not automate that process?

There are multiple other tools that accomplish this, but they are often very complex, while designing `ta` we purposely attempted to do away with all that accidental complexity: one shouldn't think about `ta` as a layout manager but as a simple way to write and execute a set of tmux commands, basic tmux knowledge is assumed, as it does most of the work.

##  Configuration file

When executed, `ta` will look for a `.ta` file in your current directory and set up a tmux session described by it, the `.ta` file will have one command per line on the following form:

```
<window name> [target]<action> [commands]
```

`<window name>` is the name of the window in which the tmux command will be executed, if it doesn't exist you'll need to create it first. More on that later.

`[target]` is an optional tmux pane number, when provided the command will be executed in the specified pane, by default it'll be done in whatever pane is active at the time.

`<action>` is the kind of command you'll be translating the line to, supported actions are:

* `c` - stands for `create`, you need to call this the first time you reference a window or the command will fail.
* `v` - stands for `vertical split`
* `h` - stands for `horizontal split`
* `a` - stands for `active pane` on a given window

`[command]` is an optional set of bash commands to execute in the **resulting** pane that is created by a given line.


Take some time to check the [example .ta file](https://github.com/lucasefe/ta/blob/master/examples/.ta) and to play around with different configurations.

## Inspiration

Mainly, [tmuxify](https://github.com/tonchis/tmuxify) from [Tonchis](https://github.com/tonchis)

## .ta example file

```
win1 c vim
win1 v foreman start
win1 h bash
win2 c tail -f log/*.log
win2 h htop
win1 a echo window win 1, will be active after all
```
