Markdown Example
================

Input file `input.html` contains an `<script type="text/markdown"></script>` element.

The build pipeline in `build.yaml` contains the `markdown` transformation.

When building like this: `gostatic build`, the script element is replaced with the rendered markdown source.

The input file can also be processed like this: `gostatic generate markdown input.html output.html`.

