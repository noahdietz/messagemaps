# protoc-gen-messagemaps

Collects all fields of type `map` with a value type of kind `Message`. Use flag
`resources_only` to limit the findings to those fields belonging to Resource
messages. Multiple runs using the same `out_file` appends the results of the
subsequent runs to the existing contents.

Run `./analyze.sh` in the projects directory to get results. This script is
specifically meant to analyze all proto packages in `google/cloud` of
[googleapis][].

Clone [googleapis] and export the variable `GOOGLEAPIS` in your shell to avoid
repetitive downloads. The script will download and set the variable itself
if unset.

[googleapis]: https://github.com/googleapis
