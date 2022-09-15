# Mokou

Importer between different 4chan SQL archive schemas

## Features

- Support for importing Asagi database dumps, fit for usage with both neofuuka and asagi dumps
- Supporting for importing genetically_enhanced_badger dumps

## Configuration

See config.example.toml and follow along with the comments.

Mokou will look for a configuration file:

- wherever the environment variable `MOKOU_CONFIG` says a file is
- at ./config.toml wherever you started mokou from

in that order.

