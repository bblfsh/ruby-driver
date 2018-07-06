#!/usr/bin/env bash
# cd to script dir
cd ${0%/*}
# set ruby search path
export GEM_HOME=../gems
# run the ruby file
ruby ../native/native
