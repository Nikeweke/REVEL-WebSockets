#!/bin/sh
SCRIPTPATH=$(cd "$(dirname "$0")"; pwd)
"$SCRIPTPATH/ws_app" -importPath ws_app -srcPath "$SCRIPTPATH/src" -runMode dev
