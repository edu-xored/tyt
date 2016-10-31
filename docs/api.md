Short API doc.

# Authentication
* POST /api/login - basic auth to get token
* POST /api/logout - close session
* GET /api/token - verifies whether auth token is still valid 

# Users
* GET /api/me - get current user
* GET /api/users - list all users
* GET /api/user/:id - get user by id
* PUT /api/user/:id - update user by id
* DELETE /api/user/:id - not needed now

# Teams
* POST /api/teams - create new team
* GET /api/teams - list all teams
* GET /api/team/:id - get team by id
* PUT /api/team/:id - update team by id
* DELETE /api/team/:id - delete specified team

# Events
* GET /api/events - list all events
* GET /api/event/:id - get event by id
* PUT /api/event/:id - update event by id
* DELETE /api/event/:id - delete event by id

# TODO Avatars
* GET /api/avatar/:user_id - get avatar image for specified user
* POST /api/avatar - uploads avatar image to server
* GET /api/myavatar - get avatar of current user
* POST /api/myavatar - change avatar of current user
