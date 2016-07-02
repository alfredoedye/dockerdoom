#!/bin/bash
set -e

/dockerdoom& 

x11vnc -geometry 640x480 -forever -usepw -create


