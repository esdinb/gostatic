Bundle Example
==============

Input file `input.html` contains an `script` element that references a javascript module.

The build pipeline in `build.yaml` contains the `bundle` transformation.

When building like this: `gostatic build`, the javascript module is bundled and placed inside the `script` element.

Output is written to `output.html`.

The input file can also be processed like this: `gostatic generate bundle input.html output.html`.

