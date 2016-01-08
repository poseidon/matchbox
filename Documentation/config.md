
# Config: Flags and Variables

`bootcfg` arguments can be provided as flags or as environment variables.

| flag | variable | example |
|------|----------|---------|
| -address | BOOTCFG_ADDRESS | 0.0.0.0:8080 |
| -config | BOOTCFG_CONFIG | ./data/config.yaml |
| -data-path | BOOTCFG_DATA_PATH | ./data |
| -images-path | BOOTCFG_IMAGES_PATH | ./images, ./static |
| -log-level | BOOTCFG_LOG_LEVEL | critical, error, warning, notice, info, debug |

## Examples

Binary

    ./run -address=0.0.0.0:8080 -data-path=./examples/dev -config=./examples/dev/config.yaml -images-path=./images -log-level=debug

Container

    docker run -p 8080:8080 --name=bootcfg --rm -v $PWD/examples/dev:/data:Z -v $PWD/images:/images:Z coreos/bootcfg:latest -address=0.0.0.0:8080 -data-path=./data -config=./data/config.yaml -images-path=./images -log-level=debug

