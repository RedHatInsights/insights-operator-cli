version=0.5
time=$(date)
go build -ldflags="-X 'main.BuildTime=$time' -X 'main.BuildVersion=$version'" .

if [ $? -eq 0 ]
then
    echo "Build ok"
    ./insights-operator-cli
else
    echo "Build failed"
fi
