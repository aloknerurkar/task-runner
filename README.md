# task-runner

As the name suggests, this library provides concurrent task execution.

I recently found an article, explaining a design of Worker Queue implementation in Go.

http://nesv.github.io/golang/2014/02/25/worker-queues-in-go.html

I would suggest going through the article first. Also the article mentions about some videos you can watch,
which explain concurrency patterns in Go. I would suggest going through that as well.

I wanted to create a library out of this, so that I can use it everywhere.

For example usage you can go through the
/example/examples.go

Basically this TaskRunner provides two interfaces to run Tasks.

Task type is pretty much anything that has a Execute() routine. Tasks have to be structured this way.

The first interface is free is a simple function call interface.

The second is the channel interface.

Some basic tests can be found in the test/ dir as well as the main pkg.