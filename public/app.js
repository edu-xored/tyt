function toJSON(res) {
 if (res.ok) {
  return res.json();
 }
 throw new Error(`http error: ${res.statusText}`);
}

function makeAPI(api) {
 const collectionPath = "/api/" + api.collection;
 const resourcePath = id => "/api/" + api.resource + "/" + id;
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

const API = {
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

function initUsersMenu() {
 const lastUserId = localStorage.getItem('last_user_id');

 API.users.getList().then(users => {
  const menu = $('#user-menu').empty();
  users.forEach(u => {
   const name = u.name || u.login;
   const a = $('<a href="#"></a>').text(name);

   const li = $('<li></li>');
   li.append(a);
   li.appendTo(menu);

   a.click(() => {
    $("#selected-user").text(name).attr('data-value', u.id);
    toggleButtonState();
    localStorage.setItem('last_user_id', u.id);
   });

   if (u.id === lastUserId) {
    $("#selected-user").text(name).attr('data-value', u.id);
    toggleButtonState();
   }
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

function makeEvent() {
 return {
  user_id: $("#selected-user").attr('data-value'),
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
 // keep user
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
 initDurationMenu();
 initUsersMenu();
 bindSubmitHandler();

 toggleButtonState();

 $("#message")
     .keyup(toggleButtonState)
     .change(toggleButtonState);
});
