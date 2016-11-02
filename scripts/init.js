const teams = require('../data/teams.json');
const users = require('../data/users.json');
const API = require('../public/client');
const push = require('./queue');

function fromGithubTeam(t) {
	return {
		name: t.name,
		description: t.description,
		slug: t.slug,
	};
}

function findUser(githubLogin) {
	const githubURL = `https://github.com/${githubLogin}`.toLowerCase();
	return users.find(u => u.github.toLowerCase() == githubURL);
}

function fromGithubUser(u) {
	const githubURL = `https://github.com/${u.login}`;
	const user = findUser(u.login);
	return Object.assign(user || {}, {
		login: u.login,
		avatar_url: u.avatar_url,
		gravatar_id: u.gravatar_id,
		github: githubURL,
		web_url: u.html_url,
	});
}

function onError(err) {
	console.log('error:', err);
}

const registeredUsers = {};

function makeUser(user) {
	return API.users.create(user).then(u => {
		console.log('user %s created', u.login);
		return u;
	}, err => {
		console.log('fail to created user %s:', user.login, err);
	});
}

teams.forEach(githubTeam => {
	push(() => {
		API.teams
			.create(fromGithubTeam(githubTeam))
			.then(team => {
				console.log('team %s created', team.name);

				const promises = (githubTeam.members || []).map(member => {
					const u = fromGithubUser(member);
					u.team_id = team.id;
					return u;
				}).map(u => {
					const existingUser = registeredUsers[u.login];
					if (existingUser) {
						console.log('user already created:', u.login);
						return new Promise(resolve => resolve(existingUser));
					}
					return makeUser(u);
				});

				Promise.all(promises)
					.then(users => {

						users.forEach(u => {
							registeredUsers[u.login] = u;
						});

						console.log('team members created');

						team.members = users.map(u => u.id);

						API.teams
							.update(team.id, team)
							.then(() => {
								console.log('team %s members linked', team.name);
							}, onError);
					});
			}, onError);
	});
});
