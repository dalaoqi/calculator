
go build -o ./bin/server

cd client/mean && go build -o ../../bin/clent1
cd -
cd client/median && go build -o ../../bin/clent2
cd -
cd client/mode && go build -o ../../bin/clent3