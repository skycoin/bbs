# `bbswatcher`

This is a simple executable that launches and monitors a `bbsnode`, restarting the process on crash or exit.

To run, ensure that `bbsnode` is located within a location specified by the `$PATH` environment variable.

Flags provided when running `bbswatcher` are transferred to the `bbsnode` child process.

### Example

Here's an example running and monitoring `bbsnode` as master and without serving the GUI.

```bash
bbswatcher -master -http-gui=false

```