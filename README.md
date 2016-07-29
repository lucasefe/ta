# ta
Tmux automation

## Automate your Tmux setup

`ta` reads your tmux window/pane definition from the .ta file located in the current directory. 

## Motivation

I'll leave that to [pote](https://github.com/pote) to explain it better. 

## Example // Syntax

Check the file located in [examples/.ta](https://github.com/lucasefe/ta/blob/master/examples/.ta)

## Inspiration

Mainly, [tmuxify](https://github.com/tonchis/tmuxify) from [Tonchis](https://github.com/tonchis)

## .ta example file

```
win1 c vim
win1 v foreman start
win1 h bash
win2 c tail -f log/*.log
win2 h htop
```
