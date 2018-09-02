#!/bin/bash
echo "`./sloccount.sh` SLOC `date` pushed by `whoami`">>pushlog.txt
git add .
git commit -m'$1'
git push -u origin master ; git push -u github master
