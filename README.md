# td [![Build Status](https://travis-ci.org/mssola/td.svg?branch=master)](https://travis-ci.org/mssola/td)

This is a CLI application for my [To Do](https://github.com/mssola/todo) service.

## Usage

### Login & logout

First of all, we have to login.

    $ td login

It will ask for the full URL of the server that is running our [To
Do](https://github.com/mssola/todo) instance, the user name and its password.
If we're successful, the user will get authenticated and its topics will be
fetched automatically.

We don't have to do this again, this user will be logged in from now on.
However, you might want to delete the current session. For this you can
just use the `logout` command.

### Commands

After logging in, you can just perform the following command:

    $ td

This will fetch the topics from your server and open up your favorite editor
(through the `EDITOR` env. variable). When you are done, close your editor and
it will automatically push to the server the topics that have changed.

Besides editing, you can `create`, `delete` and `rename`. See:

    $ td create test
    $ td rename oldname newname
    $ td delete another

Finally, note that you don't have to open the editor to know the topics that
you have. You can just perform the `list` command for that. For more
information, just use the `help` command.

### Bash completion

This package includes a shell script that offers bash completion for this
application. The completion file is `scripts/tdcompletion.sh`. Install it
wherever you want and source it in your `.bashrc` file.

## License

Copyright &copy; 2014-2017 Miquel Sabaté Solà

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
