require('es6-promise');
require('isomorphic-fetch');

const BASE = 'http://localhost:8080/api';

function toJSON(res) {
	if (res.ok) {
		return res.json();
	}
	throw new Error(`http error: ${res.statusText}`);
}

function makeAPI(api) {
	const collectionPath = `${BASE}/${api.collection}`;
	const resourcePath = id => `${BASE}/${api.resource}/${id}`;
	return {
		create(payload) {
			return fetch(collectionPath, {
				method: 'POST',
				body: JSON.stringify(payload)
			}).then(toJSON);
		},
		getList() {
			return fetch(collectionPath).then(toJSON);
		},
		get(id) {
			return fetch(resourcePath(id)).then(toJSON);
		},
		update(id, payload) {
			return fetch(resourcePath(id), {
				method: 'PUT',
				body: JSON.stringify(payload)
			}).then(toJSON);
		},
		remove(id) {
			return fetch(resourcePath(id), { method: 'DELETE' }).then(toJSON);
		},
	};
}

function basicAuth(username, password) {
	if (!username || !password) {
		return '';
	}
	return 'Basic ' + new Buffer(`${username}:${password}`).toString('base64');
}

module.exports = {
	login: function(username, password) {
		return fetch(`${BASE}/login`, {
			credentials: "same-origin",
			method: 'POST',
			headers: {
				Authorization: basicAuth(username, password),
			},
		}).then(toJSON);
	},
	me: function() {
		return fetch(`${BASE}/me`, {
			credentials: "same-origin",
			headers: makeHeaders(),
		}).then(toJSON);
	},
	users: makeAPI({
		resource: 'user',
		collection: 'users',
	}),
	teams: makeAPI({
		resource: 'team',
		collection: 'teams',
	}),
	events: makeAPI({
		resource: 'event',
		collection: 'events',
	}),
};
