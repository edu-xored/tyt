API.events.getList()
	.then(events => {
		return API.users.getList().then(users => ({
			events: events.filter(e => e.type === 'status'),
			users
		}));
	})
	.then(({events, users}) => {
		const eventMap = new Map();
		events.forEach(event => {
			const userId = event.user_id;
			const {duration, message} = event;
			if (eventMap.has(userId)) {
				eventMap.set(userId, eventMap.get(userId).concat({duration, message}))
			} else {
				eventMap.set(userId, [{duration, message}])
			}
		});
		const data = users
			.filter(user => eventMap.has(user.id))
			.map(user => [user.name, eventMap.get(user.id) || []]);
		makeChart(data)
	});


function makeChart(data) {
	$('#users-statistics').highcharts({
		chart: {
			type: 'column',
			backgroundColor: 'transparent',
			style: {
				fontFamily: 'Helvetica, sans-serif'
			}
		},
		title: {
			text: 'График занятости',
			style: {
				color: '#ffffff'
			}
		},
		xAxis: {
			type: 'category',
			labels: {
				rotation: -45,
				style: {
					fontSize: '13px',
					fontFamily: 'Helvetica, sans-serif',
					color: '#ffffff'
				}
			}
		},
		yAxis: {
			min: 0,
			labels: {
				style: {
					color: '#ffffff'
				}
			},
			title: {
				text: 'Количество часов',
				style: {
					color: '#ffffff'
				}
			}
		},
		legend: {
			enabled: false
		},
		tooltip: {
			formatter: getFormattedTooltip(data)
		},
		series: [{
			name: 'users',
			data: getUserTotalDuration(data),
			dataLabels: {
				enabled: true,
				rotation: -90,
				color: '#FFFFFF',
				align: 'right',
				format: '{point.y:.1f}', // one decimal
				y: 10, // 10 pixels down from the top
				style: {
					fontSize: '13px',
					fontFamily: 'Helvetica, sans-serif',
					textShadow: false
				}
			}
		}
		]
	});
}

function getUserTotalDuration(data) {
	return data.map(i => [i[0], i[1].reduce((sum, current) => {
			return sum + current.duration;
		}, 0)
		]
	);
}

function getFormattedTooltip(data) {
	const events = new Map(data);
	return function () {
		if (!events.get(this.key).length) {
			return 'No Events';
		}
		return events.get(this.key).reduce((a, e) => a + e.duration, 0) + 'h';
	};
}
