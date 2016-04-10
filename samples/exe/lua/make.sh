#!/bin/sh

# brew install lua lua51

luac-5.2 -o example1-52.luac example1.lua
luac-5.2 -o example2-52.luac example2.lua

luac-5.1 -o example1-51.luac example1.lua
luac-5.1 -o example2-51.luac example2.lua
