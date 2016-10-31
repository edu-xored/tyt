require('es6-promise');
require('isomorphic-fetch');

// TODO allow to list teams, events

const API = require('./client');

API.users.getList().then(list => {
	console.log(list);
}, err => {
	console.log('error:', err);
});
