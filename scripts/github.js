require('es6-promise');
require('isomorphic-fetch');

const jsonfile = require('jsonfile');
const fs = require('fs');
const push = require('./queue');

const token = process.env.TOKEN;

const fetchOptions = {
	headers: {
		Accept: 'application/vnd.github.v3+json',
		Authorization: `token ${token}`,
	}
};

function print(d) {
	// console.log(d.length);
	console.log(d);
	return d;
}

function get(path) {
	const url = `https://api.github.com${path}`;
	return fetch(url, fetchOptions).then(res => res.json());
}

// push(() => get(`/orgs/edu-xored/members`));
// push(() => get(`/orgs/edu-xored/teams`));

// dump members
// get(`/orgs/edu-xored/members`).then(d => {
// 	jsonfile.writeFile('members.json', d, {spaces: 2});
// });

// dump teams
get(`/orgs/edu-xored/teams`).then(teams => {
	// jsonfile.writeFile('teams.json', teams, {spaces: 2});
	const members = teams.map(t => fetch(t.members_url.replace('{/member}', ''), fetchOptions).then(r => r.json()))
	return Promise.all(members).then(data => {
		teams.forEach((t, i) => {
			t.members = data[i];
		});
		jsonfile.writeFile('teams.json', teams, {spaces: 2});
	});
});
