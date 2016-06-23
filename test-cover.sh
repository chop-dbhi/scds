#!/bin/sh

# Run test coverage on each subdirectories and merge the coverage profile.
echo "mode: count" > profile.cov

RETVAL=0

# Standard go tooling behavior is to ignore dirs with leading underscores
for dir in $(find . -maxdepth 10 -not -path './vendor*' -not -path './.git*' -not -path '*/_*' -type d); do
    if find $dir -type f | grep \.go$ > /dev/null; then
        go test -v -covermode=count -coverprofile=$dir/profile.tmp $dir || RETVAL=1

        if [ -f $dir/profile.tmp ]; then
            cat $dir/profile.tmp | tail -n +2 >> profile.cov
            rm $dir/profile.tmp
        fi
    fi
done

go tool cover -func profile.cov

exit $RETVAL
