#!/bin/sh

usage() {
    cat << EOF
Usage:
    $0 webserver    Start the webserver.
    $0 -h | --help  Show this help message.

Options:
    -h --help       Show this help message.
EOF
}

if [ -z "$1" ] || [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
    usage
    exit 0
fi

COMMAND=$1

./manage.py migrate -v 0
./manage.py ensureadminuser --skip-on-missing-env
 
case "$COMMAND" in
    "webserver")
        python -m granian --interface asginl --access-log freon_server/asgi.py:application
        ;;
    *)
        echo "error: unknown command: $COMMAND\n"
        usage
        exit 1
        ;;
esac