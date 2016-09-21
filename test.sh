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


# Build mocks
mkdir -p ${startDir}/src/github.com/goblimey/films/mocks/gomock
dir='github.com/goblimey/films/mocks/gomock'
echo ${dir}
cd ${startDir}/src/$dir
mockgen --package gomock net/http ResponseWriter >mock_response_writer.go
mockgen --package gomock github.com/goblimey/films/retrofit/template Template >mock_template.go
mockgen --package gomock github.com/goblimey/films/repositories/people Repository >mock_people_repository.go

mkdir -p ${startDir}/src/github.com/goblimey/films/mocks/pegomock
dir='github.com/goblimey/films/mocks/pegomock'
echo ${dir}
cd ${startDir}/src/$dir
pegomock generate --package pegomock --output=mock_template.go github.com/goblimey/films/retrofit/template Template
pegomock generate --package pegomock --output=mock_response_writer.go net/http ResponseWriter

# Build

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

dir='github.com/goblimey/films/repositories/people'
echo ${dir}
cd ${startDir}/src/$dir
${testcmd}

dir='github.com/goblimey/films/controllers/people'
echo ${dir}
cd ${startDir}/src/$dir
${testcmd}
