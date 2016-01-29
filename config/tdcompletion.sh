# Copyright (C) 2014-2016 Miquel Sabaté Solà <mikisabate@gmail.com>
# This file is licensed under the MIT license.
# See the LICENSE file.

# Auto-complete the td command in Bash.

# Configurable values.
_DIR="$HOME"
_FILE=".td/config.json"


# Do the completion work.
__tdcomp()
{
    local all c s=$'\n' IFS=' '$'\t'$'\n'
    local cur="${COMP_WORDS[COMP_CWORD]}"

    for c in $1; do
        case "$c$4" in
        --*=*) all="$all$c$4$s" ;;
        *.)    all="$all$c$4$s" ;;
        *)     all="$all$c$4 $s" ;;
        esac
    done
    IFS=$s
    COMPREPLY=($(compgen -P "$2" -W "$all" -- "$cur"))
    return
}

# Main function for the completion of the td command.
_td()
{
    # First of all, check whether the current user is logged in or not.
    cmds="login"
    if [ -f "$_DIR/$_FILE" ]; then
        # Maybe it exists but it's empty.
        contents=`cat $_DIR/$_FILE`
        if [ "$contents" != "" ]; then
            cmds="create delete fetch list logout push rename status"
        fi
    fi

    local c=1 command
    while [ $c -lt $COMP_CWORD ]; do
        command="${COMP_WORDS[c]}"
        c=$((++c))
    done

    # Complete a command.
    if [ $c -eq $COMP_CWORD -a -z "$command" ]; then
        case "${COMP_WORDS[COMP_CWORD]}" in
        -*|--*) __tdcomp "--help --version" ;;
        *)      __tdcomp "$cmds" ;;
        esac
        return
    fi

    # If we reach this point, then the command has already been written.
    # Therefore, we only have to check for commands that accept a known
    # parameter.

    topics=$(td list | xargs)

    case "$command" in
    rename|delete)  __tdcomp "${topics}" ;;
    *) COMPREPLY=() ;;
    esac
}

complete -o default -o nospace -F _td td

