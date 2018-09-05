#!/bin/bash
echo "`./sloccount.sh` SLOC (go) `date` pushed by `whoami`">>pushlog.txt
git add .
git commit -m'`echo $1`'
git push -u origin master ; git push -u github master
