# leakybucket
Leaky bucket pattern implementation in Go

Server periodically performs function which outcome might be success or failure.
In case when number of failures exceeds a given threshold value for a given period of time,
the server is marked as unhealthy and temporarily disabled. 
The budget of errors is recovered continually at a specified rate (the "leaky bucket" part), and when it achieves a given value, the server is enabled again. 