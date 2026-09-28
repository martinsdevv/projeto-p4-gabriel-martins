#!/bin/sh
cd "$(dirname "$0")" || exit 1
javac -encoding UTF-8 -d out -sourcepath src src/matchmaking/Main.java || exit 1
exec java -cp out matchmaking.Main "$@"
