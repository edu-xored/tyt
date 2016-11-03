function durationLabel(val) {
	switch (val) {
		case 0.5:
			return 'полчаса';
		case 1:
			return '1 час'
		case 2:
		case 3:
		case 4: 
			return val + " часа";
		default:
			return val + " часов";
	}
}

function initDurationMenu() {
	[0.5, 1, 2, 3, 4, 5, 6, 7, 8].forEach(val => {
		const label = durationLabel(val);
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
		.text(durationLabel(1))
		.attr('data-value', 1);
}

function bindSubmitHandler() {
	$("#form").submit(e => {
		e.preventDefault();
		send();
	});
}

function fetchCurrentUser() {
	API.me().then(user => {
		currentUser = user;
		if (user) {
			$(".greeting").text("Привет, " + user.name + "!");
		}
	}, err => {
		console.log("api error:", err);
	});
}

$(function() {
	initDurationMenu();
	bindSubmitHandler();
	toggleButtonState();

	$("#message")
		.keyup(toggleButtonState)
		.change(toggleButtonState);

	fetchCurrentUser();
});
