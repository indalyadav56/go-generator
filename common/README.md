# Guide

## Adding submodule to a repository

To add a submodule to a repository, follow these steps

1. Run `git submodule init`
2. Run `git submodule add <submodule_url>`

## Cloning a repository

To clone a repository which already has submodules, follow these steps

1. Run `git clone <repository_url>`
2. Run `git submodule init`
3. Run `git submodule update --recursive`

## NOTE:

Never run the command `git submodule update --remote`. It will not work as expected for us

## Updating a submodule

To update version of a submodule in a repository, follow these steps

1. CD into the submodule direcotry: `cd <path/to/submodule/directory>`. E.g. `cd paydoh-commons`
2. Run `git fetch --tags`
3. Checkout to the required version: `git checkout <tagname>`
4. CD to repository's main directory. E.g. `cd ..`
5. Commit the changes and push to the server
