const API = require('../public/client');
const push = require('./queue');

const users = API.users;

function getUsers() {
	push(users.getList);
}

function createUser(user) {
	push(() => users.create(user));
}

function createUsers() {
	for (let i = 0; i < 10; i++) {
		createUser({
			name: `user${i + 1}`,
			email: `user${i + 1}@example.com`,
		})
	}
}

function print(d) {
	console.log(d);
	return d;
}

function testCrud() {
	users.create({
		name: 'test',
		email: 'test@example.com',
	}).then(print, print)
		.then(user => users.get(user.id), print)
		.then(print, print)
		.then(user => {
			user.name = 'newname';
			return users.update(user.id, user);
		}, print)
		.then(print, print)
		.then(user => users.remove(user.id))
		.then(print, print)
	;
}

// createUsers();
// getUsers();
testCrud();
