const API = require('../public/client');

API.login(process.argv[2], process.argv[3]).then(d => {
	console.log(d);
}, err => {
	console.log('error:', err);
});
