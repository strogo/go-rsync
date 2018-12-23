#!/bin/bash
#
# Example showing use of getopt detection and use of GNU enhanced getopt
# to handle arguments containing whitespace.
#
# Written in 2004 by Hoylen Sue <hoylen@hoylen.com>
# Modified in 2018 by Denis Dyakov <denis.dyakov@gmail.com>
#
# To the extent possible under law, the author(s) have dedicated all copyright and
# related and neighboring rights to this software to the public domain worldwide.
# This software is distributed without any warranty.
#
# You should have received a copy of the CC0 Public Domain Dedication along with this software.
# If not, see <http://creativecommons.org/publicdomain/zero/1.0/>.

PROG=$(basename $0)
VERSION=v0.1

# Define default values, if parameters not specified
RELEASE_TYPE="Release"
DEV_TYPE="Development"

# Remove this trap if you are doing your own error detection or don't care about errors
trap "echo $PROG: error encountered: aborted; exit 3" ERR

#----------------------------------------------------------------
# Process command line arguments

## Define options: trailing colon means has an argument (customize this: 1 of 3)

SHORT_OPTS=b:h
LONG_OPTS=buildtype:,version,help

SHORT_HELP="Usage: ${PROG} [options] arguments
Options:
  -b <build type>         Build type. Release type = [${RELEASE_TYPE}].
  -h                        Show this help message."

LONG_HELP="Usage: ${PROG} [options] arguments
Options:
  -b | --buildtype <build type>       Build type. Release type = [${RELEASE_TYPE}].
  -h | --help                         Show this help message.
  --version                           Show version information."

# Detect if GNU Enhanced getopt is available

HAS_GNU_ENHANCED_GETOPT=
if getopt -T >/dev/null; then :
else
  if [ $? -eq 4 ]; then
    HAS_GNU_ENHANCED_GETOPT=yes
  fi
fi

# Run getopt (runs getopt first in `if` so `trap ERR` does not interfere)

if [ -n "$HAS_GNU_ENHANCED_GETOPT" ]; then
  # Use GNU enhanced getopt
  if ! getopt --name "$PROG" --long $LONG_OPTS --options $SHORT_OPTS -- "$@" >/dev/null; then
    echo "$PROG: usage error (use -h or --help for help)" >&2
    exit 2
  fi
  ARGS=`getopt --name "$PROG" --long $LONG_OPTS --options $SHORT_OPTS -- "$@"`
else
  # Use original getopt (no long option names, no whitespace, no sorting)
  if ! getopt $SHORT_OPTS "$@" >/dev/null; then
    echo "$PROG: usage error (use -h for help)" >&2
    exit 2
  fi    
  ARGS=`getopt $SHORT_OPTS "$@"`
fi
eval set -- $ARGS

## Process parsed options (customize this: 2 of 3)
 
while [ $# -gt 0 ]; do
    case "$1" in
        -b | --buildtype)    BUILDTYPE="$2"; shift;;
        -v | --verbose)  VERBOSE=yes;;
        -h | --help)     if [ -n "$HAS_GNU_ENHANCED_GETOPT" ]
                         then echo "$LONG_HELP";
                         else echo "$SHORT_HELP";
                         fi;  exit 0;;
        --version)       echo "$PROG $VERSION"; exit 0;;
        --)              shift; break;; # end of options
    esac
    shift
done


if [ "$BUILDTYPE" == "$RELEASE_TYPE" ]; then
  echo "Release in progress..."
  go run data/generate/generate.go && mv ./assets_vfsdata.go ./data
  go build -v -ldflags="-X main.version=`head -1 version` -X main.buildnum=`date -u +%Y%m%d%H%M%S`" -tags gorsync_rel gorsync.go
else
  echo "Dev in progress..."
  go build -v -ldflags="-X main.version=`head -1 version` -X main.buildnum=`date -u +%Y%m%d%H%M%S`" gorsync.go
fi

