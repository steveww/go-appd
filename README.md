# go-appd
Go wrapper around the AppDynamics C SDK

Only tested on Linux

You will need to set some environmental variables before you can use this.

    APPD_SDK_HOME=/path/to/the/sdk_lib
    CGO_CFLAGS="-I $APPD_SDK_HOME"
    CGO_LDFLAGS="-L $APPD_SDK_HOME/lib -l appdynamics_native_sdk"
    
Now just install

    go get github.com/steveww/go-appd
