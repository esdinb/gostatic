XML-Stylesheet Example
======================

Input file `input.html` contains an `xml-stylesheet` processing instruction.

The build pipeline in `build.yaml` contains the `template:inline` transformation.

When building like this: `gostatic build`, the XSL stylesheet is extracted from the input file and applied to its contents.

The input file can also be processed like this: `gostatic generate template:inline input.html output.html`.

