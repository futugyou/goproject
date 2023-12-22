# install

```cmd
dotnet tool install --global Microsoft.OpenApi.Kiota
export PATH="$PATH:/home/codespace/.dotnet/tools"
```

```cmd
kiota generate -l go -c PostsClient -n kiota-samples/client -d ./posts-api.yml -o ./client
kiota generate -l go -c GraphApiClient -n kiota-samples/client2 -d ./get-me.yml -o ./client2
```
