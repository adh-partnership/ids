# SignalR Package

This is a partial SignalR implementation for the ADH IDS purposes...

## Limitations

- We will only support the WebSockets transport. We will not support SSE or Long Polling
- We will only support the transport format `Text`
- We will rely on the authentication to be handled externally, not using connection id/token