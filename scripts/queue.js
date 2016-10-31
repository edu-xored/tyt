const Queue = require('promise-queue');
const queue = new Queue(1, 10000);

function push(makeRequest) {
	queue.add(() => {
		makeRequest().then(d => {
			console.log(d);
		}, err => {
			console.log('error', err);
		});
	});
}

module.exports = push;
