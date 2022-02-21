# Multiple Stack Support

`bossctl` supports multiple stack configurations in a single configuration
file. The format of the **v3** configuration file that supports multiple
stacks is of the format

```yaml
apiVersion: v3
currentStack: <stack-name-reference>
stacks:
  - name: <stack-name-reference>
    <values>
```

## `--stack` command line option

As part of the support for multiple stacks a new command line option,
`--stack` was added. This option can be used to specify which stack
configuration should be used for each command invocation.

## Configuration File Manipulations

A new `bossctl` subcommand, `stack` was added to support the listing,
deletion, and addition of stacks to the configuration file. In addtion
the `stack use` command was added to persistently set the current
stack name back to the configuration. _(see `bossctl stack --help`
for more information)_
