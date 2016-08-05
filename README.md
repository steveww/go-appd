# go-appd
Go wrapper around the AppDynamics C SDK
See official documentation of AppD C/C++ SDK [Guide](https://docs.appdynamics.com/pages/viewpage.action?pageId=35233447) [Reference](https://docs.appdynamics.com/pages/viewpage.action?pageId=35233445)

Only tested on Linux

You will need to install the AppDynamics C/C++ SDK first get it from [here](http://download.appdynamics.com/) You will need an account to login in to this portal, see your AppDynamics account manager or sign up for a free trial.

Take a look at the included example and note how the configuration is set. The **access key** and **account name** is unique to your account, you can get this from the License page in the Controller UI. When the SDK proxy is running the logs subdirectory will provide useful diagnostic data. Note how any backends are defined first before they are used.

You will need to set some environmental variables before you can use this.

    APPD_SDK_HOME=/path/to/the/sdk_lib
    CGO_CFLAGS="-I $APPD_SDK_HOME"
    CGO_LDFLAGS="-L $APPD_SDK_HOME/lib -l appdynamics_native_sdk"
    LD_LIBRARY_PATH=$APPD_SDK_HOME/lib
    
Now just install

    go get github.com/steveww/go-appd
