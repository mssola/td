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

In this section I'm assuming that you've already logged in. At this point there
are a lot of things that can be done. You can see the full list of commands by
executing the following command:

    $ td --help

First of all, let's see how we can interact with the server. There are
basically two commands for this: `fetch` and `push`. The `fetch` command pulls
everything from the server and saves it to our local setup. If there are
some local changes that have not been pushed yet, then it will ask for
permission before fetching the topics. The `push` command will grab all
our local editions and push it to the server.

One can create and delete topics with the `create` and `delete` commands,
respectively. Note that both commands expect the name of the topic to be
created/deleted. Moreover, a topic can be renamed with the `rename` command.
This command expects two parameters: the old and the new name, in this order.
As an example:

    $ td create test
    $ td rename test another
    $ td delete another

The most important thing is to view the "To do" list and edit it. For this you
don't have to pass any parameter to the `td` executable. If you do that, then
the following will happen:

1. The topics will be fetched automatically from the server (equivalent to `td
   fetch`).
2. Your preferred editor (`$EDITOR`) will be openned in the directory where all
   the topics are being stored. Each topic is represented as a file, and the
   name of the topic is the name of its file.
3. The topics that you have modified will be pushed to the server (equivalent
   to `td push`, but only pushing the modified topics).

The previous points are important because it means that you will rarely use
the `push` and `fetch` commands: the `td` command alone will do the right thing
always.

Finally, note that you don't have to open the editor to know the topics that
you have. You can just perform the `list` command for that.

### Bash completion

This package includes a shell script that offers bash completion for this
application. The completion file is `config/tdcompletion.sh`. Install it
wherever you want and source it in your `.bashrc` file.

## License

Copyright &copy; 2014-2016 Miquel Sabaté Solà

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

