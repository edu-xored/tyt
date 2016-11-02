const API = require('../public/client');

// TODO allow to list teams, events

API.users.getList().then(list => {
	console.log(list);
}, err => {
	console.log('error:', err);
});
