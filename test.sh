#! /bin/sh

# Run tests.  This script assumes that the current directory is the one in which
# it lives.  With no argument, run all tests.  With argument "unit" run just the
# unit tests.  With argument "int" run just the integration tests.

testcmd='go test -test.v'
if test ! -z $1
then
	case $1 in
	unit )
		testcmd="$testcmd -run='^TestUnit'";;
	int )
		testcmd="$testcmd -run='^TestInt'";;
	* )
		echo "first argument must be unit or integration" >&2
		exit -1
		;;
	esac
fi

startDir=`pwd`

. ./setenv.sh

# Build

# cd src
# go install -v -gcflags "-N -l" ./...
go build github.com/goblimey/films

# Test

dir='github.com/goblimey/films/models/person'
echo ${dir}
cd ${startDir}/src/$dir
${testcmd}

dir='github.com/goblimey/films/models/person/gorpmysql'
echo ${dir}
cd ${startDir}/src/$dir
${testcmd}

dir='github.com/goblimey/films/forms/people'
echo ${dir}
cd ${startDir}/src/$dir
${testcmd}

dir='github.com/goblimey/films/daos/people'
echo ${dir}
cd ${startDir}/src/$dir
${testcmd}
