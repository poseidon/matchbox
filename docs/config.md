
# Config: Flags and Variables

`bootcfg` arguments can be provided as flags or as environment variables.

| flag | varaible | example |
|------|----------|---------|
| -address | BOOTCFG_ADDRESS | 0.0.0.0:8080 |
| -data-path | BOOTCFG_DATA_PATH | ./data |
| -images-path | BOOTCFG_IMAGES_PATH | ./images, ./static |
| -log-level | BOOTCFG_LOG_LEVEL | critical, error, warning, notice, info, debug |

## Examples

Binary

    ./run -address=0.0.0.0:8080 -data-path=./data -images-path=./images -log-level=info

Container

    docker run -p 8080:8080 --name=bootcfg --rm -v $PWD/data:/data -v $PWD/images:/images coreos/bootcfg:latest -address=0.0.0.0:8080 -data-path=./data -images-path=./images

