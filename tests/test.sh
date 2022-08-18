#!/bin/bash -e

# dialects=("sqlite" "mysql" "postgres" "sqlserver")
dialects=("mysql")

if [[ $(pwd) == *"gen/tests"* ]]; then
  cd ..
fi

if [ -d tests ]
then
  cd tests
  go get -u -t ./...
  go mod download
  go mod tidy
  cd ..
fi

# SqlServer for Mac M1
if [[ -z $GITHUB_ACTION ]]; then
  if [ -d tests ]
  then
    cd tests
    if [[ $(uname -a) == *" arm64" ]]; then
      MSSQL_IMAGE=mcr.microsoft.com/azure-sql-edge docker-compose start || true
      go install github.com/microsoft/go-sqlcmd/cmd/sqlcmd@latest || true
      SQLCMDPASSWORD=LoremIpsum86 sqlcmd -U sa -S localhost:9930 -Q "IF DB_ID('gen') IS NULL CREATE DATABASE gen" > /dev/null || true
      SQLCMDPASSWORD=LoremIpsum86 sqlcmd -U sa -S localhost:9930 -Q "IF SUSER_ID (N'gen') IS NULL CREATE LOGIN gen WITH PASSWORD = 'LoremIpsum86';" > /dev/null || true
      SQLCMDPASSWORD=LoremIpsum86 sqlcmd -U sa -S localhost:9930 -Q "IF USER_ID (N'gen') IS NULL CREATE USER gen FROM LOGIN gen; ALTER SERVER ROLE sysadmin ADD MEMBER [gen];" > /dev/null || true
    else
      docker-compose start
    fi
    cd ..
  fi
fi

for dialect in "${dialects[@]}" ; do
  if [ "$GORM_DIALECT" = "" ] || [ "$GORM_DIALECT" = "${dialect}" ]
  then
    echo "testing ${dialect}..."

    if [ "$GEN_VERBOSE" = "" ]
    then
      if [ -d tests ]
      then
        cd tests
        GORM_DIALECT=${dialect} go test -race -count=1 ./...
        cd ..
      fi
    else
      if [ -d tests ]
      then
        cd tests
        GORM_DIALECT=${dialect} go test -race -count=1 -v ./...
        cd ..
      fi
    fi
  fi
done
