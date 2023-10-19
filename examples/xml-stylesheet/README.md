XML-Stylesheet Example
======================

Input file `input.xml` contains an `xml-stylesheet` processing instruction.

Input file `input.html` is a html file without processing instructions.

The build pipeline in `build.yaml` contains the `template:inline` transformation.

When building like this: `gostatic build`, the XSL stylesheet is extracted from the input file and applied to its contents.

The input file can also be processed like this: `gostatic generate template:inline input.xml output.html`.

Or like this: `gostatic generate template:style.xsl input.html output.html`.


