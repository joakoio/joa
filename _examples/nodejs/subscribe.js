var zmq = require('zeromq')
var subscriber = zmq.socket('sub')

subscriber.on("message", function(filter, message) {
    console.log('Received message: ', message.toString());
})

subscriber.connect("tcp://localhost:5551")
subscriber.subscribe("A")

process.on('SIGINT', function() {
    subscriber.close()
    console.log('\nClosed')
})