
go build -o ./bin/server

cd client/mean && go build -o ../../bin/client1
cd -
cd client/median && go build -o ../../bin/client2
cd -
cd client/mode && go build -o ../../bin/client3