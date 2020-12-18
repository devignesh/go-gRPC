# go-gRPC
Goalng gRPC practise 


# gRPC 
In gRPC, a client application can directly call a method on a server application on a different machine as if it were a local object, making it easier for you to create distributed applications and services. As in many RPC systems, gRPC is based around the idea of defining a service, specifying the methods that can be called remotely with their parameters and return types. On the server side, the server implements this interface and runs a gRPC server to handle client calls. On the client side, the client has a stub (referred to as just a client in some languages) that provides the same methods as the server.

# service method:
    Unary RPCs :
        rpc SayHello(HelloRequest) returns (HelloResponse);
    Server streaming RPCs:
        rpc LotsOfReplies(HelloRequest) returns (stream HelloResponse);
    Client streaming RPCs:
        rpc LotsOfGreetings(stream HelloRequest) returns (HelloResponse);
    Bidirectional streaming RPCs:
        rpc BidiHello(stream HelloRequest) returns (stream HelloResponse);
