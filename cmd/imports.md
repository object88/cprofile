# Imports

The imports command takes a path to a package, and returns a list of imported packages, in alphabetical order.

## Syntax

`cprofile imports path-to-package [optional flags]`

## Optional flags

### `--astDepth [value]`

Determines how deep into the AST to locate imports.  There are 5 levels:

* `shallow` (default)
* `deep`
* `local`
* `wide`
* `complete`

Consider the very simplified source tree for the `cprofile` app.  The application's main entry point is in `main`, which imports `cmds`, which in turn imports `cprofile`, `cobra`, and `pflags`, etc., etc.

```text
[GOPATH]
└─src
  └─github.com
    ├─object88
    │ ├─cprofile [imports context, logging]
    │ │ ├─cmds [imports cprofile, cobra, pflags]
    │ │ └─main [imports cmds, os]
    │ └─logging [imports fmt]
    └─spf13
      ├─cobra [imports pflags]
      └─pflags

[GOROOT]
└─src
  ├─context
  ├─fmt
  └─os
```

The output from the `imports` command will report a different list of packages for each possible setting.  Presuming the provided package is `main`, then...

* The `shallow` setting will list only the packages imported directly by `cprofile/main`: `cprofile/cmds`.
* The `deep` setting will list the packages imported by the specified by the provided package and any imported packages which fall under the provided package.  In this example, that is still only `cprofile/cmds`, but a more complete example will be available elsewhere.
* The `local` setting will list the packages imported by the specified package, all packages that it imports, as long as it's under the same organization (i.e., github.com/object88).  In this case, that would include `cprofile/main`, `cprofile/cmds`, `cprofile`, and `logging`.
  * If there is a circular reference from one organization back to the organization of the specified package, it will not be included.
* The `wide` setting will include everything that is not part of the standard Go library.  In this case, it's everything reported by `local`, plus `spf13/cobra` and `spf13/pflags`.
* The `complete` setting will report every package in the tree.

### `--verbose`