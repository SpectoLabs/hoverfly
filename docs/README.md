# Hoverfly Docs
These are the docs that are used for https://hoverfly.readthedocs.io

## Building locally

```shell script
$ make html
$ open _builds/html/index.html
```

## Updating simulation examples

```shell script
$ ./pages/simulations/update.sh
```

## Pinning the dependencies

Run the following command after all python dependencies are installed:

```shell script
$ pip freeze > requirements.txt
```
