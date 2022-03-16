#!/bin/bash -xe
git remote add upstream https://github.com/upstream/repo.git
git pull --rebase upstream master
git push --force-with-lease origin master
git checkout master
git reset --hard upstream/master