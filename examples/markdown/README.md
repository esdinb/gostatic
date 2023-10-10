Markdown Example
================

Input file `input.html` contains an `<section is="markdown-element"></section>` element.

The build pipeline in `build.yaml` contains the `markdown` transformation.

When building like this: `gostatic build`, the section element content is replaced with the rendered markdown source.

The input file can also be processed like this: `gostatic generate markdown input.html output.html`.

