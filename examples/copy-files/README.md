Copy Files Example
==================

Input file `input.txt` is just a text file.

The `build.yaml` contains no transformations. Input files with file extensions other than
`.html`and `.xml` are copied without applying any transformations.

When building like this: `gostatic build`, the input file is copied to the output path.

The output path is `output.txt`.

The `input.txt` file can also be copied like this: `gostatic generate input.txt output.txt`.

