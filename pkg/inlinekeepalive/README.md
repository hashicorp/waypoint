# inlinekeepalive

inlinekeepalive is a package that sends "keepalive" messages over
existing grpc streams. See `doc.go` for full package description


## Diagrams

These diagrams demonstrate the intended usage of inline keepalives within Waypoint.

### ClientStream

```mermaid
sequenceDiagram
autonumber

actor Waypoint Runner
participant Client
participant ClientInterceptor
participant ALB
participant ServerInterceptor
participant RunnerConfig Server Handler
participant Server


Waypoint Runner ->> +Client: RunnerConfig

%% Unnecessary detail
%% Client ->> ServerInterceptor: Open connection for RunnerConfig
%% ServerInterceptor -x ClientInterceptor: Not a ServerStream - DOES NOT send keepalives

Client ->> ClientInterceptor: Create the GRPC client handler

ClientInterceptor ->> +Server: GetVersionInfo

Server -->> -ClientInterceptor: Has feature inlinekeepalives

loop Async send inline keepalives for duration of RPC
    activate ClientInterceptor
    ClientInterceptor -) +ServerInterceptor: SendMsg(InlineKeepalive)

    %% Critical
    ServerInterceptor -x RunnerConfig Server Handler: Recognizes InlineKeepalive. DOES NOT forward.
   
    deactivate ServerInterceptor
    deactivate ClientInterceptor
end

Waypoint Runner ->> Client: Send(event)
Client ->> ClientInterceptor: Send(event) to interceptor
ClientInterceptor ->> ServerInterceptor: Passthrough event to server

%% Critical
ServerInterceptor ->> RunnerConfig Server Handler: Not a keepalive, passthrough

deactivate Client

```


### ServerStream

```mermaid
sequenceDiagram
autonumber

actor User
participant Client
participant ClientInterceptor
participant ALB
participant ServerInterceptor
participant GetLogStream Server Handler
actor CEB


User ->> +Client: GetLogStream

Client ->> +ServerInterceptor: Open connection for RunnerConfig, including inline-keepalive GRPC metadata

loop Async send inline keepalives for duration of RPC
    ServerInterceptor -) +ClientInterceptor: SendMsg(InlineKeepalive)
    
    %% Critical
    ClientInterceptor -x Client: Recognizes InlineKeepalive. DOES NOT forward.
    
    deactivate ClientInterceptor
end

Client ->> ClientInterceptor: Create the GRPC client handler

CEB ->> GetLogStream Server Handler: Send(Log)
GetLogStream Server Handler ->> ServerInterceptor: Send(Log)
ServerInterceptor ->> ClientInterceptor: Send(event) to interceptor

%% Critical
ClientInterceptor ->> Client: Not a keepalive, passthrough

Client ->> User: 

deactivate Client
```
