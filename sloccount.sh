find $1|grep '\.go$'|xargs cat|grep -v '^\s*//.*'|grep -v '^\s*$'|wc -l
