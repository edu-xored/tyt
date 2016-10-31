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
	const collectionPath = "/api/" + api.collection;
	const resourcePath = id => "/api/" + api.resource + "/" + id;
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
	me: function() {
		return fetch('/api/me', {
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

function initDurationMenu() {
	[0.5, 1, 2, 3, 4, 5, 6, 7, 8].forEach(val => {
		const label = val + " hour";
		const a = $('<a href="#"></a>').text(label);
		const li = $('<li></li>');
		li.append(a);
		li.appendTo($("#duration-menu"));

		a.click(() => {
			$("#duration").text(label).attr('data-value', val);
			toggleButtonState();
		});
	});
}

function toggleButtonState() {
	const valid = isValidEvent(makeEvent());
	$("#btn-send").attr('disabled', !valid);
}

function isValidEvent(event) {
	return !!event.user_id && !!(event.message || '').trim();
}

var currentUser = {};

function makeEvent() {
	return {
		user_id: currentUser.id,
		type: 'status',
		message: $("#message").val(),
		duration: parseInt($("#duration").attr('data-value')),
	};
}

function send() {
	const event = makeEvent();
	if (!isValidEvent(event)) {
		return;
	}
	API.events.create(event).then(event => {
		// TODO show notification that status is reported successfully
		reset();
	}, err => {
		alert(err);
	});
}

function reset() {
	$("#message").val('');
	// reset duratoin
	$("#duration")
		.text('1 hour')
		.attr('data-value', 1);
}

function bindSubmitHandler() {
	$("#form").submit(e => {
		e.preventDefault();
		send();
	});
}

$(function() {
	API.me().then(user => {
		currentUser = user;
		if (user) {
			$(".greeting").text("Hey, " + user.name);
		}
	});

	initDurationMenu();
	bindSubmitHandler();

	toggleButtonState();

	$("#message")
		.keyup(toggleButtonState)
		.change(toggleButtonState);
});
