go build

if [ $? -eq 0 ]
then
    echo "Build ok"
    ./insights-operator-cli
else
    echo "Build failed"
fi
