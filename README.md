# Go web app with user service
Built following Jonathan Calhoun's "Web Development with Go".
Only currently working feature is a user database, and 
a stocklist service is being implemented.

Server is layered in a model-view-controller pattern. 
Homepage is served by a static controller and user-related pages 
are served by a users controller. 

GET requests are delegated by the static controller and the users controller 
to the view layer, which renders the base webpages 
using the bootstrap framework.

POST requests on the user-related pages are handled by the users model. 
The model is separated into a UserService layer which deals with 
authentication and hashing, and which wraps a lower UserDB layer 
that deals with database interaction. 
It works in conjunction with the services model which is responsible 
for the actual establishing of a connection with the database.

The Require Users middleware intercepts handlers which require a login 
to verify if a user is authenticated. 
If so, it adds user information to the request context and 
forwards the request to the corresponding handle function; 
if not, it redirects the client to the login page.
