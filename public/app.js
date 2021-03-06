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
var currentLecture;

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

function extractFirstName(fullName) {
	const a = (fullName || '').split(/\s/g).map(t => t.trim()).filter(t => !!t);
	return a.length >= 2 ? a[1] : fullName;
}

function fetchCurrentUser() {
	API.me().then(user => {
		currentUser = user;
		if (user) {
			$(".greeting").text("Привет, " + extractFirstName(user.name) + "!");
		}
	}, err => {
		console.log("api error:", err);
	});
}

function initLectureUI() {
	moment.locale('ru');
	const now = moment();
	API.spectacles.getList().then(allLectures => {
		allLectures.sort((a, b) => new Date(a.start).getTime() - new Date(b.start).getTime());
		const lectures = allLectures.filter(t => moment(new Date(t.start)).add(t.duration, 'hours').isAfter(now));
		const lecture = lectures[0];
		if (lecture) {
			$('.lecture-info').show();

			const start = moment(new Date(lecture.start));
			const isCurrent = start.isSameOrBefore(now);

			$('.lecture-label').text(isCurrent ? 'Текущий доклад' : 'Следующий доклад');
			$('.lecture-title').text(lecture.title);

			const btnDate = $('.btn-lecture-date').text(start.format('DD MMMM, [18:10]'));

			if (isCurrent) {
				currentLecture = lecture;
				btnDate.hide();
				initHereButton();
			} else {
				btnDate.show();
				$('.btn-here').hide();
			}
		} else {
			$('.lecture-info').hide();
		}
	});
}

function initHereButton() {
	// hide here button if the user is already reported presense
	API.events.getList().then(events => {
		const alreadyReported = events.some(e => e.type === 'presence' && e.spectacle_id === currentLecture.id);
		if (alreadyReported) {
			$('.btn-here').hide();
		} else {
			$('.btn-here').show();
		}
	});

	$('.btn-here').click(() => {
		if (!currentUser || !currentLecture) {
			return;
		}
		API.iamhere({
			spectacle_id: currentLecture.id,
		}).then(() => {
			$('.btn-here').hide();
		}, err => {
			// TODO use sweet alert
			alert(err);
		});
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
	initLectureUI();
});
