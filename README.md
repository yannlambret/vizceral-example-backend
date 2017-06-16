# vizceral-example-backend

This is a minimalist go backend to test Netflix's vizceral-example project by using websockets.

The backend periodically reloads a JSON file provided by the original project and push resulting
data towards the javascript frontend. I mainly did it in order to test my own version of the
frontend.

## Usage

```bash
go run main.go
```

Then you have to run the vizceral-exemple project on the same host, and connect with a browser
to ```http://localhost:8080```.

You can then edit the JSON file and get an updated view within 5 seconds.
