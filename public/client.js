const isServer = typeof module !== 'undefined';

if (isServer) {
  require('es6-promise');
  require('isomorphic-fetch');
}

const BASE = isServer ? 'http://localhost:8080/api' : '/api';

function basicAuth(username, password) {
	if (!username || !password) {
		return '';
	}
	const s = `${username}:${password}`;
	return 'Basic ' + (isServer ? new Buffer(s).toString('base64') : btoa(s));
}

function toJSON(res) {
	if (res.ok) {
		return res.json();
	}
	throw new Error(`http error: ${res.statusText}`);
}

function makeHeaders() {
	return {
		// Cookie: document.cookie,
	};
}

function makeAPI(api) {
	const collectionPath = `${BASE}/${api.collection}`;
	const resourcePath = id => `${BASE}/${api.resource}/${id}`;
	return {
		create(payload) {
			return fetch(collectionPath, {
				credentials: "same-origin",
				method: 'POST',
				body: JSON.stringify(payload),
				headers: makeHeaders(),
			}).then(toJSON);
		},
		getList() {
			return fetch(collectionPath, {
				credentials: "same-origin",
				headers: makeHeaders(),
			}).then(toJSON);
		},
		get(id) {
			return fetch(resourcePath(id), {
				credentials: "same-origin",
				headers: makeHeaders(),
			}).then(toJSON);
		},
		update(id, payload) {
			return fetch(resourcePath(id), {
				credentials: "same-origin",
				method: 'PUT',
				body: JSON.stringify(payload),
				headers: makeHeaders(),
			}).then(toJSON);
		},
		remove(id) {
			return fetch(resourcePath(id), {
				credentials: "same-origin",
				method: 'DELETE',
				headers: makeHeaders(),
			}).then(toJSON);
		},
	};
}

const API = {
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
	iamhere: function(payload) {
		return fetch(`${BASE}/iamhere`, {
			credentials: "same-origin",
			method: 'POST',
			body: JSON.stringify(payload),
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
	orgs: makeAPI({
		resource: 'org',
		collection: 'orgs',
	}),
	events: makeAPI({
		resource: 'event',
		collection: 'events',
	}),
	spectacles: makeAPI({
		resource: 'spectacle',
		collection: 'spectacles',
	}),
};

if (isServer) {
  module.exports = API;
} else if (typeof window !== 'undefined') {
	window.API = API;
}
